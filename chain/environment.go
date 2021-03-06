// Copyright © 2017 The Blocknet Developers
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package chain

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/BlocknetDX/dxregress/containers"
	"github.com/BlocknetDX/dxregress/util"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Environment interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type EnvConfig struct {
	ConfigPath          string
	ContainerPrefix     string
	DefaultImage        string
	ContainerFilter     string
	ContainerFilterFunc func(filter string) string
	DockerFileName      string
	Activator           Node
	Nodes               []Node
	XWallets            []XWallet
}

// TestEnv is the default implementation for a test environment.
type TestEnv struct {
	config *EnvConfig
	docker *client.Client
	xwalletNodes []Node
}

// Start the environment.
func (env *TestEnv) Start(ctx context.Context) error {
	// Write test blocknetdx.conf file
	testLocalenvDir := path.Dir(TestBlocknetConfFile(env.config.ConfigPath))
	if err := os.MkdirAll(testLocalenvDir, 0775); err != nil {
		return errors.Wrapf(err, "Failed to create directory %s", testLocalenvDir)
	}
	if err := ioutil.WriteFile(TestBlocknetConfFile(env.config.ConfigPath), []byte(TestBlocknetConf(env.config.Nodes)), 0644); err != nil {
		errors.Wrapf(err, "Failed to write blocknetdx.conf %s", TestBlocknetConfFile(env.config.ConfigPath))
	}

	// Stop all containers
	logrus.Info("Removing previous test containers...")
	if err := containers.StopAllContainers(ctx, env.docker, env.config.ContainerFilter, true); err != nil {
		logrus.Error(err)
	}

	// Start containers
	for _, c := range env.config.Nodes {
		if err := containers.CreateAndStart(ctx, env.docker, env.config.DefaultImage, c.Name, c.Ports); err != nil {
			return err
		}
		if c.DebuggerPort != "" {
			logrus.Infof("%s node running on %s, rpc on %s, gdb/lldb port on %s", c.Name, c.Port, c.RPCPort, c.DebuggerPort)
		} else {
			logrus.Infof("%s node running on %s, rpc on %s", c.Name, c.Port, c.RPCPort)
		}
	}

	// Start xwallet containers
	for _, w := range env.config.XWallets {
		// Ignore BYOW nodes (bring your own wallet)
		if w.BringOwn {
			continue
		}
		// Create node from xwallet
		xn := NodeForWallet(w, env.config.ContainerPrefix)
		env.xwalletNodes = append(env.xwalletNodes, xn)
		if err := containers.CreateContainer(ctx, env.docker, w.Container, xn.Name, xn.Ports); err != nil {
			return err
		}
		if xn.DebuggerPort != "" {
			logrus.Infof("%s node running on %s, rpc on %s, gdb/lldb port on %s", xn.Name, xn.Port, xn.RPCPort, xn.DebuggerPort)
		} else {
			logrus.Infof("%s node running on %s, rpc on %s", xn.Name, xn.Port, xn.RPCPort)
		}
	}

	// Import xwallet addresses if necessary
	for _, xn := range env.xwalletNodes {
		// Copy wallet conf
		xC := containers.FindContainer(env.docker, xn.Name)
		ctx2, cancel := context.WithTimeout(ctx, 10 * time.Second)
		defer cancel()
		rd, _, err := env.docker.CopyFromContainer(ctx2, xn.Name, fmt.Sprintf("/opt/blockchain/config/%s", xn.Conf))
		if err != nil {
			return errors.Wrapf(err, "Failed to copy %s from %s", xn.Conf, xn.Name)
		}
		defer rd.Close()

		// Read xwallet conf
		xwalletConf, err := readFromTar(xn.Conf, rd)
		if err != nil {
			return err
		}

		// Update xwallet conf
		newWalletConf := AddXWalletRPC(xwalletConf, xn.RPCUser, xn.RPCPass)
		tr := map[string][]byte{xn.Conf: newWalletConf}
		if buf, err := util.CreateTar(tr); err == nil {
			if err := env.docker.CopyToContainer(ctx, xC.ID, "/opt/blockchain/config/", buf, types.CopyToContainerOptions{}); err != nil {
				return errors.Wrapf(err, "Failed to write %s to %s", xn.Conf, xn.Name)
			}
		} else {
			return errors.Wrapf(err, "Failed to write %s to %s", xn.Conf, xn.Name)
		}
	}
	// Start xwallet containers
	for _, xn := range env.xwalletNodes {
		xC := containers.FindContainer(env.docker, xn.Name)
		if err := containers.StartContainer(ctx, env.docker, xC.ID); err != nil {
			return err
		}
	}
	logrus.Info("Waiting for xwallets to be ready...")
	if err := WaitForEnv(ctx, 120, env.xwalletNodes); err != nil {
		return err
	}
	for _, xn := range env.xwalletNodes {
		// Import address if required
		if xn.AddressKey == "" {
			continue
		}
		cmd := RPCCommand(xn.Name, xn.CLI, fmt.Sprintf("importprivkey %s coin", xn.AddressKey))
		if output, err := cmd.Output(); err != nil {
			logrus.Error(errors.Wrapf(err, "Problem importing xwallet address %s in %s", xn.Address, xn.Name))
		} else {
			logrus.Debug(string(output))
		}
	}

	logrus.Info("Waiting for nodes to be ready...")
	if err := WaitForEnv(ctx, 300, env.config.Nodes); err != nil {
		return err
	}

	// Setup blockchain
	if err := env.setupChain(ctx, env.docker); err != nil {
		return err
	}

	return nil
}

// Stop the environment, including performing necessary tear down.
func (env *TestEnv) Stop(ctx context.Context) error {
	if err := containers.StopAllContainers(ctx, env.docker, env.config.ContainerFilterFunc(""), false); err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

// setupChain will setup the DX environment, copy all configuration files, test RPC connectivity.
func (env *TestEnv) setupChain(ctx context.Context, docker *client.Client) error {
	// Nodes
	activator := env.config.Activator
	snodes := ServiceNodes(env.config.Nodes)
	activatorC := containers.FindContainer(docker, activator.Name)

	// Import wallet addresses
	for _, node := range env.config.Nodes {
		// skip nodes without wallets
		if node.AddressKey == "" {
			continue
		}
		cmd := RPCCommand(node.Name, node.CLI, fmt.Sprintf("importprivkey %s coin", node.AddressKey))
		if output, err := cmd.Output(); err != nil {
			return errors.Wrapf(err, "Failed to import wallet address %s to %s", node.Address, node.Name)
		} else {
			logrus.Debug(string(output))
		}
	}

	// First import test address into alias and then generate test coin
	for i := 0; i < 15; i++ {
		cmd := BlockRPCCommand(activator.Name, "setgenerate true 1")
		if output, err := cmd.Output(); err != nil || string(output) == "" {
			if err != nil {
				return errors.Wrap(err, "Failed to generate first 15 blocks")
			} else {
				return errors.New("Failed to generate first 15 blocks, empty output")
			}
		} else {
			logrus.Debug(string(output))
		}
		time.Sleep(time.Second)
	}

	// Import alias addresses
	cmd2 := BlockRPCCommands(activator.Name, []string{"importprivkey cRdLcWroNyJPJ1BH4Q24pamDQtE3JNdm7tGQoD6mm9brqpYuX1dC sn1", "importprivkey cMn9aiQGBYqeRzRuTFAModv459UQNxGsXkgPSRQ1W7XwGdGCp1JB sn2"})
	if output, err := cmd2.Output(); err != nil {
		return errors.Wrap(err, "Failed to import alias addresses")
	} else {
		logrus.Debug(string(output))
	}

	// Send coin to traders
	for _, node := range env.config.Nodes {
		if node.ID != Trader {
			continue
		}
		var cmds []string
		for i := 0; i < 20; i++ {
			cmds = append(cmds, fmt.Sprintf("sendtoaddress %s 500", node.Address))
		}
		cmd := BlockRPCCommands(activator.Name, cmds)
		if output, err := cmd.Output(); err != nil {
			return errors.Wrapf(err, "Failed to send coin from activator to %s", node.Name)
		} else {
			logrus.Debug(string(output))
		}
	}

	// Break up 10K inputs into 2.5k inputs to help with staking
	cmd4S := make([]string, 25)
	for i := 0; i < len(cmd4S); i++ {
		cmd4S[i] = fmt.Sprintf("sendtoaddress %s 2500", activator.Address)
	}
	cmd4 := BlockRPCCommands(activator.Name, cmd4S)
	if output, err := cmd4.Output(); err != nil {
		return errors.Wrap(err, "Failed to split coin")
	} else {
		logrus.Debug(string(output))
	}

	// Send 5k servicenode coin to each alias
	var snodeInputCmds []string
	for _, snode := range snodes {
		snodeInputCmds = append(snodeInputCmds, fmt.Sprintf("sendtoaddress %s 5000", snode.Address))
	}
	cmd5k := BlockRPCCommands(activator.Name, snodeInputCmds)
	if output, err := cmd5k.Output(); err != nil {
		return errors.Wrap(err, "Failed to send 5k servicenode coin")
	} else {
		logrus.Debug(string(output))
	}

	// Generate last PoW blocks
	for i := 0; i < 15; i++ {
		cmd5 := BlockRPCCommand(activator.Name, "setgenerate true 1")
		if output, err := cmd5.Output(); err != nil {
			return errors.Wrap(err, "Failed to generate blocks 15-30")
		} else {
			logrus.Debug(string(output))
		}
		time.Sleep(time.Second)
	}

	// Obtain servicenode keys
	var keys []string
	for _, snode := range snodes {
		cmdSnode := BlockRPCCommand(snode.Name, "servicenode genkey")
		if output, err := cmdSnode.Output(); err != nil {
			return errors.Wrapf(err, "Failed to call genkey on %s", snode.Name)
		} else {
			keys = append(keys, strings.TrimSpace(string(output)))
		}
	}

	// Setup activator servicenode.conf
	type OutputsResponse struct {
		TxID  string `json:"txhash"`
		TxPos int    `json:"outputidx"`
	}
	cmd7 := BlockRPCCommand(activator.Name, "servicenode outputs")
	output, err := cmd7.Output()
	if err != nil {
		return errors.Wrap(err, "Failed to parse servicenode outputs")
	}
	var outputs []OutputsResponse
	if err := json.Unmarshal(output, &outputs); err != nil {
		return errors.Wrap(err, "Failed to parse servicenode outputs")
	}

	// Create servicenode specific data provider
	var servicenodes []SNode
	for j, snode := range snodes {
		ssn := SNode{
			ID: snode.ID,
			Alias: snode.ShortName,
			IP: snode.IP(),
			Key: keys[j],
			CollateralID: outputs[j].TxID,
			CollateralPos: strconv.Itoa(outputs[j].TxPos),
		}
		servicenodes = append(servicenodes, ssn)
	}

	// Max wait time for all commands below
	ctx, cancel := context.WithTimeout(context.Background(), 600 * time.Second)
	defer cancel()

	// Generate activator servicenode.conf
	snConf := ServicenodeConf(servicenodes)
	// Copy activator servicenode.conf
	if servicenodeConf, err := util.CreateTar(map[string][]byte{"servicenode.conf": []byte(snConf)}); err == nil {
		if err := docker.CopyToContainer(ctx, activatorC.ID, "/opt/blockchain/dxregress/testnet4/", servicenodeConf, types.CopyToContainerOptions{}); err != nil {
			return errors.Wrap(err, "Failed to write servicenode.conf to activator")
		}
	} else {
		return errors.Wrap(err, "Failed to write servicenode.conf to activator")
	}

	// Update blocknetdx.conf on nodes (skip servicenodes)
	for _, node := range env.config.Nodes {
		if node.IsSnode {
			continue
		}
		nodeC := containers.FindContainer(docker, node.Name)
		blocknetConf := BlocknetdxConf(node, env.config.Nodes, true, "")
		if buf, err := util.CreateTar(map[string][]byte{"blocknetdx.conf": []byte(blocknetConf)}); err == nil {
			if err := docker.CopyToContainer(ctx, nodeC.ID, "/opt/blockchain/config/", buf, types.CopyToContainerOptions{}); err != nil {
				return errors.Wrapf(err, "Failed to write blocknetdx.conf to %s", node.Name)
			}
		} else {
			return errors.Wrapf(err, "Failed to write blocknetdx.conf to %s", node.Name)
		}
	}

	// Copy config files to servicenodes
	xbridgeConfSnode := XBridgeConf(env.config.XWallets)
	for _, ssn := range servicenodes {
		sn := NodeForID(ssn.ID, snodes)
		ssnC := containers.FindContainer(env.docker, sn.Name)

		// Update servicenodes blocknetdx.conf
		blocknetConfSn := BlocknetdxConf(sn, env.config.Nodes, true, ssn.Key)
		if bufSn, err := util.CreateTar(map[string][]byte{"blocknetdx.conf": []byte(blocknetConfSn)}); err == nil {
			if err := docker.CopyToContainer(ctx, ssnC.ID, "/opt/blockchain/config/", bufSn, types.CopyToContainerOptions{}); err != nil {
				return errors.Wrapf(err, "Failed to write blocknetdx.conf to %s", sn.ShortName)
			}
		} else {
			return errors.Wrapf(err, "Failed to write blocknetdx.conf to %s", sn.ShortName)
		}

		// Write servicenodes xbridge.conf
		if bufSn, err := util.CreateTar(map[string][]byte{"xbridge.conf": []byte(xbridgeConfSnode)}); err == nil {
			if err := docker.CopyToContainer(ctx, ssnC.ID, "/opt/blockchain/dxregress/", bufSn, types.CopyToContainerOptions{}); err != nil {
				return errors.Wrapf(err, "Failed to write xbridge.conf to %s", sn.ShortName)
			}
		} else {
			return errors.Wrapf(err, "Failed to copy xbridge.conf to %s", sn.ShortName)
		}
	}

	// Stop activator
	if err := containers.StopContainer(ctx, docker, activatorC.ID); err != nil {
		return err
	}

	// Restart all nodes except for wallet nodes
	if err := containers.RestartContainers(ctx, docker, env.config.ContainerFilterFunc("sn")); err != nil {
		return err
	}
	if err := containers.RestartContainers(ctx, docker, env.config.ContainerFilterFunc("act")); err != nil {
		return err
	}
	//if err := containers.RestartContainers(ctx, docker, env.config.ContainerFilterFunc("trader")); err != nil {
	//	return err
	//}

	// Wait for nodes to be ready
	logrus.Info("Waiting for nodes and wallets to be ready...")
	allContainers := append(env.config.Nodes, env.xwalletNodes...)
	if err := WaitForEnv(ctx, 500, allContainers); err != nil {
		return err
	}

	// Start servicenodes
	if err := StartServicenodesFrom(activator.Name); err != nil {
		return err
	}

	// Wait before restarting staker
	logrus.Info("Waiting to start staking on activator...")
	time.Sleep(10 * time.Second)

	// Restart the activator to trigger staking
	if err := containers.RestartContainers(ctx, docker, env.config.ContainerFilterFunc("act")); err != nil {
		return err
	}

	// Wait for activator to be ready
	logrus.Info("Waiting for activator to be ready...")
	if err := WaitForEnv(ctx, 45, []Node{activator}); err != nil {
		return err
	}

	// Call start servicenodes a second time to make sure they're started
	// Wait before re-running snode command
	time.Sleep(5 * time.Second)
	if err := StartServicenodesFrom(activator.Name); err != nil {
		return err
	}

	// TODO Check if wallets are reachable

	return nil
}

// NewTestEnvironment creates a new test environment instance.
func NewTestEnv(config *EnvConfig, docker *client.Client) *TestEnv {
	env := new(TestEnv)
	env.config = config
	env.docker = docker
	return env
}

// readFromTar reads the specified file from the tar.
func readFromTar(file string, tarFile io.Reader) ([]byte, error) {
	tr := tar.NewReader(tarFile)
	var fileBytes []byte
	// Iterate over tar entries
	for {
		header, err := tr.Next()
		if err == io.EOF {
			return fileBytes, errors.Wrapf(err, "Failed to find %s", file)
		} else if err != nil {
			return fileBytes, errors.Wrapf(err, "Failed to read %s", file)
		}
		info := header.FileInfo()
		// Check if file name match
		if info.Name() == file {
			var buf bytes.Buffer
			if _, err := io.Copy(&buf, tr); err != nil {
				return fileBytes, errors.Wrapf(err, "Failed to read %s", file)
			}
			// Read file bytes
			fileBytes := buf.Bytes()
			return fileBytes, nil
		}
	}
	return fileBytes, errors.New("Failed to find " + file)
}