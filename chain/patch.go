package chain

import (
	"encoding/base64"

	"github.com/sirupsen/logrus"
)

// GenesisPatchV1 contains the patch for enabling the new test environment blockchain.
func GenesisPatchV1() string {
	patch := `ZGlmZiAtLWdpdCBhL3NyYy9jaGFpbnBhcmFtcy5jcHAgYi9zcmMvY2hhaW5wYXJhbXMuY3BwCmluZGV4IGY0NGJkMmU4Yi4uNzYyNDZhNzg0IDEwMDY0NAotLS0gYS9zcmMvY2hhaW5wYXJhbXMuY3BwCisrKyBiL3NyYy9jaGFpbnBhcmFtcy5jcHAKQEAgLTE5NCwzMiArMTk0LDMyIEBAIHB1YmxpYzoKICAgICAgICAgcGNoTWVzc2FnZVN0YXJ0WzFdID0gMHg3NjsKICAgICAgICAgcGNoTWVzc2FnZVN0YXJ0WzJdID0gMHg2NTsKICAgICAgICAgcGNoTWVzc2FnZVN0YXJ0WzNdID0gMHhiYTsKLSAgICAgICAgdkFsZXJ0UHViS2V5ID0gUGFyc2VIZXgoIjAwMDAxMGU4M2IyNzAzY2NmMzIyZjdkYmQ2MmRkNTg1NWFjN2MxMGJkMDU1ODE0Y2UxMjFiYTMyNjA3ZDU3M2I4ODEwYzAyYzA1ODJhZWQwNWI0ZGViOWM0Yjc3YjI2ZDkyNDI4YzYxMjU2Y2Q0Mjc3NGJhYmVhMGEwNzNiMmVkMGM5Iik7Ci0gICAgICAgIG5EZWZhdWx0UG9ydCA9IDQxNDc0OworLy8gICAgICAgIHZBbGVydFB1YktleSA9IFBhcnNlSGV4KCIwMDAwMTBlODNiMjcwM2NjZjMyMmY3ZGJkNjJkZDU4NTVhYzdjMTBiZDA1NTgxNGNlMTIxYmEzMjYwN2Q1NzNiODgxMGMwMmMwNTgyYWVkMDViNGRlYjljNGI3N2IyNmQ5MjQyOGM2MTI1NmNkNDI3NzRiYWJlYTBhMDczYjJlZDBjOSIpOworICAgICAgICBuRGVmYXVsdFBvcnQgPSA0MTQ3NjsKICAgICAgICAgbkVuZm9yY2VCbG9ja1VwZ3JhZGVNYWpvcml0eSA9IDUxOwogICAgICAgICBuUmVqZWN0QmxvY2tPdXRkYXRlZE1ham9yaXR5ID0gNzU7CiAgICAgICAgIG5Ub0NoZWNrQmxvY2tVcGdyYWRlTWFqb3JpdHkgPSAxMDA7CiAgICAgICAgIG5NaW5lclRocmVhZHMgPSAwOwotICAgICAgICBuVGFyZ2V0VGltZXNwYW4gPSAxICogNjA7IC8vIEJsb2NrbmV0RFg6IDEgZGF5Ci0gICAgICAgIG5UYXJnZXRTcGFjaW5nID0gMSAqIDYwOyAgLy8gQmxvY2tuZXREWDogMSBtaW51dGUKLSAgICAgICAgbkxhc3RQT1dCbG9jayA9IDIwMDA7Ci0gICAgICAgIG5NYXR1cml0eSA9IDE1OwotICAgICAgICBuTW9kaWZpZXJVcGRhdGVCbG9jayA9IDUxMTk3OyAvL2FwcHJveCBNb24sIDE3IEFwciAyMDE3IDA0OjAwOjAwIEdNVAorICAgICAgICBuVGFyZ2V0VGltZXNwYW4gPSA2MDsgLy8gQmxvY2tuZXREWDogMSBtaW51dGUKKyAgICAgICAgblRhcmdldFNwYWNpbmcgPSAzMDsgIC8vIEJsb2NrbmV0RFg6IDMwIHNlY29uZHMKKyAgICAgICAgbkxhc3RQT1dCbG9jayA9IDUwOworICAgICAgICBuTWF0dXJpdHkgPSAxMDsKKyAgICAgICAgbk1vZGlmaWVyVXBkYXRlQmxvY2sgPSAyXjMyOwogCiAgICAgICAgIC8vISBNb2RpZnkgdGhlIHRlc3RuZXQgZ2VuZXNpcyBibG9jayBzbyB0aGUgdGltZXN0YW1wIGlzIHZhbGlkIGZvciBhIGxhdGVyIHN0YXJ0LgotICAgICAgICBnZW5lc2lzLm5UaW1lID0gMTUwMzU3MTAwMDsKLSAgICAgICAgZ2VuZXNpcy5uTm9uY2UgPSAyMTU4OTYyOworICAgICAgICBiblByb29mT2ZXb3JrTGltaXQgPSB+dWludDI1NigpID4+IDE7CisgICAgICAgIGdlbmVzaXMublRpbWUgPSAxNTA5MTYzMDY3OworICAgICAgICBnZW5lc2lzLm5Ob25jZSA9IDEyMzQ1Njc7CiAKICAgICAgICAgaGFzaEdlbmVzaXNCbG9jayA9IGdlbmVzaXMuR2V0SGFzaCgpOwotICAgICAgICBhc3NlcnQoaGFzaEdlbmVzaXNCbG9jayA9PSB1aW50MjU2KCIweDAwMDAwZjkwYWMyNjA4NTllNDUxNTM1NjcxOWQ5NGM5ZmI4Y2FkYjFhM2RkYTE4NmE2NGFjNDFjZTRjM2M3YTciKSk7CisgICAgICAgIGFzc2VydChoYXNoR2VuZXNpc0Jsb2NrID09IHVpbnQyNTYoIjY5YWYwMzliZjZhOWVkMzdmNTcxNzkyNTkwZGVkZGMwZWYyZDQ4YmU3MDg0OTljODk5ODBiMWNiMzljNDQxYjEiKSk7CiAKICAgICAgICAgdkZpeGVkU2VlZHMuY2xlYXIoKTsKICAgICAgICAgdlNlZWRzLmNsZWFyKCk7Ci0JCi0JCi0gICAgICAgIHZTZWVkcy5wdXNoX2JhY2soQ0ROU1NlZWREYXRhKCIxNzguNjIuOTAuMjEzIiwgIjE3OC42Mi45MC4yMTMiKSk7IC8vIHNlZWQgbm9kZQotICAgICAgICB2U2VlZHMucHVzaF9iYWNrKENETlNTZWVkRGF0YSgiMTM4LjE5Ny43My4yMTQiLCAiMTM4LjE5Ny43My4yMTQiKSk7IC8vIHNlZWQgbm9kZQotICAgICAgICB2U2VlZHMucHVzaF9iYWNrKENETlNTZWVkRGF0YSgiMzQuMjM1LjQ5LjI0OCIsICIzNC4yMzUuNDkuMjQ4IikpOyAvLyBzZWVkIG5vZGUKKworLy8gICAgICAgIHZTZWVkcy5wdXNoX2JhY2soQ0ROU1NlZWREYXRhKCIxNzguNjIuOTAuMjEzIiwgIjE3OC42Mi45MC4yMTMiKSk7IC8vIHNlZWQgbm9kZQorLy8gICAgICAgIHZTZWVkcy5wdXNoX2JhY2soQ0ROU1NlZWREYXRhKCIxMzguMTk3LjczLjIxNCIsICIxMzguMTk3LjczLjIxNCIpKTsgLy8gc2VlZCBub2RlCisvLyAgICAgICAgdlNlZWRzLnB1c2hfYmFjayhDRE5TU2VlZERhdGEoIjM0LjIzNS40OS4yNDgiLCAiMzQuMjM1LjQ5LjI0OCIpKTsgLy8gc2VlZCBub2RlCiAKICAgICAgICAgYmFzZTU4UHJlZml4ZXNbUFVCS0VZX0FERFJFU1NdID0gc3RkOjp2ZWN0b3I8dW5zaWduZWQgY2hhcj4oMSwgMTM5KTsgLy8gVGVzdG5ldCBibG9ja25ldGR4IGFkZHJlc3NlcyBzdGFydCB3aXRoICd4JyBvciAneScKICAgICAgICAgYmFzZTU4UHJlZml4ZXNbU0NSSVBUX0FERFJFU1NdID0gc3RkOjp2ZWN0b3I8dW5zaWduZWQgY2hhcj4oMSwgMTkpOyAgLy8gVGVzdG5ldCBibG9ja25ldGR4IHNjcmlwdCBhZGRyZXNzZXMgc3RhcnQgd2l0aCAnOCcgb3IgJzknCkBAIC0yMzEsMTggKzIzMSwyMCBAQCBwdWJsaWM6CiAgICAgICAgIC8vIFRlc3RuZXQgYmxvY2tuZXRkeCBCSVA0NCBjb2luIHR5cGUgaXMgJzEnIChBbGwgY29pbidzIHRlc3RuZXQgZGVmYXVsdCkKICAgICAgICAgYmFzZTU4UHJlZml4ZXNbRVhUX0NPSU5fVFlQRV0gPSBib29zdDo6YXNzaWduOjpsaXN0X29mKDB4ODApKDB4MDApKDB4MDApKDB4MDEpLmNvbnZlcnRfdG9fY29udGFpbmVyPHN0ZDo6dmVjdG9yPHVuc2lnbmVkIGNoYXI+ID4oKTsKIAotICAgICAgICBjb252ZXJ0U2VlZDYodkZpeGVkU2VlZHMsIHBuU2VlZDZfdGVzdCwgQVJSQVlMRU4ocG5TZWVkNl90ZXN0KSk7CisvLyAgICAgICAgY29udmVydFNlZWQ2KHZGaXhlZFNlZWRzLCBwblNlZWQ2X3Rlc3QsIEFSUkFZTEVOKHBuU2VlZDZfdGVzdCkpOwogCiAgICAgICAgIGZSZXF1aXJlUlBDUGFzc3dvcmQgPSB0cnVlOwotICAgICAgICBmTWluaW5nUmVxdWlyZXNQZWVycyA9IHRydWU7Ci0gICAgICAgIGZBbGxvd01pbkRpZmZpY3VsdHlCbG9ja3MgPSB0cnVlOworICAgICAgICBmTWluaW5nUmVxdWlyZXNQZWVycyA9IGZhbHNlOworICAgICAgICBmQWxsb3dNaW5EaWZmaWN1bHR5QmxvY2tzID0gZmFsc2U7CiAgICAgICAgIGZEZWZhdWx0Q29uc2lzdGVuY3lDaGVja3MgPSBmYWxzZTsKLSAgICAgICAgZlJlcXVpcmVTdGFuZGFyZCA9IGZhbHNlOwotICAgICAgICBmTWluZUJsb2Nrc09uRGVtYW5kID0gZmFsc2U7Ci0gICAgICAgIGZUZXN0bmV0VG9CZURlcHJlY2F0ZWRGaWVsZFJQQyA9IHRydWU7CisgICAgICAgIGZSZXF1aXJlU3RhbmRhcmQgPSB0cnVlOworICAgICAgICBmTWluZUJsb2Nrc09uRGVtYW5kID0gdHJ1ZTsKKyAgICAgICAgZlNraXBQcm9vZk9mV29ya0NoZWNrID0gdHJ1ZTsKKyAgICAgICAgZlRlc3RuZXRUb0JlRGVwcmVjYXRlZEZpZWxkUlBDID0gZmFsc2U7CisgICAgICAgIGZIZWFkZXJzRmlyc3RTeW5jaW5nQWN0aXZlID0gZmFsc2U7CiAKICAgICAgICAgblBvb2xNYXhUcmFuc2FjdGlvbnMgPSAyOwotICAgICAgICBzdHJTcG9ya0tleSA9ICIwNDU2NWY0MjliOGM2OGRhYmRjZGYwOTYwOGJlMDViNGExMGZjNzA0ZjdkZTE4NjZhYWVlZDI4YTcyOWVjNWI4YzQxOGQ5MGY5NTEwYmExMTVjMGYzNWYzNTNiMWFlYTk4M2ZlOTkzOTdkMWMyMDY4NWQ2YWIzZWQwZDBiN2JhM2VhIjsKKyAgICAgICAgc3RyU3BvcmtLZXkgPSAiMDQ2YWU3ZGU3YzhmNmJiMmU2YmNiYzFkYmRiNzFjYzQzNDFlYzRiNGVjOWYxZDdhMTZiNmJjYzYxYzQ2MGQ1NDIxYWZlMWIwNDg3YWNjYzQ2ZjJlZTAwNDQ5Mjc3ZjliMjE2MWU1ZTUxZWJmYjA4YTNjNTM2MWI4MTFlYTIyZDNiOSI7CiAgICAgICAgIHN0ck9iZnVzY2F0aW9uUG9vbER1bW15QWRkcmVzcyA9ICJ5NTdjcWZHUmtla1J5RFJOZUppTHRZVkVidmhYck5ibW94IjsKICAgICAgICAgblN0YXJ0U2VydmljZW5vZGVQYXltZW50cyA9IDE0MjA4Mzc1NTg7IC8vRnJpLCAwOSBKYW4gMjAxNSAyMTowNTo1OCBHTVQKICAgICB9CmRpZmYgLS1naXQgYS9zcmMvaW5pdC5jcHAgYi9zcmMvaW5pdC5jcHAKaW5kZXggMjViOWNlNzNlLi4wMGJlMDkwNzIgMTAwNjQ0Ci0tLSBhL3NyYy9pbml0LmNwcAorKysgYi9zcmMvaW5pdC5jcHAKQEAgLTM1MSw3ICszNTEsNyBAQCBzdGQ6OnN0cmluZyBIZWxwTWVzc2FnZShIZWxwTWVzc2FnZU1vZGUgbW9kZSkKICAgICBzdHJVc2FnZSArPSBIZWxwTWVzc2FnZU9wdCgiLW9ubHluZXQ9PG5ldD4iLCBfKCJPbmx5IGNvbm5lY3QgdG8gbm9kZXMgaW4gbmV0d29yayA8bmV0PiAoaXB2NCwgaXB2NiBvciBvbmlvbikiKSk7CiAgICAgc3RyVXNhZ2UgKz0gSGVscE1lc3NhZ2VPcHQoIi1wZXJtaXRiYXJlbXVsdGlzaWciLCBzdHJwcmludGYoXygiUmVsYXkgbm9uLVAyU0ggbXVsdGlzaWcgKGRlZmF1bHQ6ICV1KSIpLCAxKSk7CiAgICAgc3RyVXNhZ2UgKz0gSGVscE1lc3NhZ2VPcHQoIi1wZWVyYmxvb21maWx0ZXJzIiwgc3RycHJpbnRmKF8oIlN1cHBvcnQgZmlsdGVyaW5nIG9mIGJsb2NrcyBhbmQgdHJhbnNhY3Rpb24gd2l0aCBibG9vbSBmaWx0ZXJzIChkZWZhdWx0OiAldSkiKSwgZmFsc2UpKTsKLSAgICBzdHJVc2FnZSArPSBIZWxwTWVzc2FnZU9wdCgiLXBvcnQ9PHBvcnQ+Iiwgc3RycHJpbnRmKF8oIkxpc3RlbiBmb3IgY29ubmVjdGlvbnMgb24gPHBvcnQ+IChkZWZhdWx0OiAldSBvciB0ZXN0bmV0OiAldSkiKSwgNDE0MTIsIDQxNDc0KSk7CisgICAgc3RyVXNhZ2UgKz0gSGVscE1lc3NhZ2VPcHQoIi1wb3J0PTxwb3J0PiIsIHN0cnByaW50ZihfKCJMaXN0ZW4gZm9yIGNvbm5lY3Rpb25zIG9uIDxwb3J0PiAoZGVmYXVsdDogJXUgb3IgdGVzdG5ldDogJXUpIiksIDQxNDEyLCA0MTQ3NikpOwogICAgIHN0clVzYWdlICs9IEhlbHBNZXNzYWdlT3B0KCItcHJveHk9PGlwOnBvcnQ+IiwgXygiQ29ubmVjdCB0aHJvdWdoIFNPQ0tTNSBwcm94eSIpKTsKICAgICBzdHJVc2FnZSArPSBIZWxwTWVzc2FnZU9wdCgiLXNlZWRub2RlPTxpcD4iLCBfKCJDb25uZWN0IHRvIGEgbm9kZSB0byByZXRyaWV2ZSBwZWVyIGFkZHJlc3NlcywgYW5kIGRpc2Nvbm5lY3QiKSk7CiAgICAgc3RyVXNhZ2UgKz0gSGVscE1lc3NhZ2VPcHQoIi10aW1lb3V0PTxuPiIsIHN0cnByaW50ZihfKCJTcGVjaWZ5IGNvbm5lY3Rpb24gdGltZW91dCBpbiBtaWxsaXNlY29uZHMgKG1pbmltdW06IDEsIGRlZmF1bHQ6ICVkKSIpLCBERUZBVUxUX0NPTk5FQ1RfVElNRU9VVCkpOwpkaWZmIC0tZ2l0IGEvc3JjL2tlcm5lbC5jcHAgYi9zcmMva2VybmVsLmNwcAppbmRleCBhMzc2YTg1NmIuLmU4NzdkMjQ5ZSAxMDA2NDQKLS0tIGEvc3JjL2tlcm5lbC5jcHAKKysrIGIvc3JjL2tlcm5lbC5jcHAKQEAgLTI1Niw2ICsyNTYsOCBAQCBib29sIEdldEtlcm5lbFN0YWtlTW9kaWZpZXIodWludDI1NiBoYXNoQmxvY2tGcm9tLCB1aW50NjRfdCYgblN0YWtlTW9kaWZpZXIsIGludAogICAgIC8vIGxvb3AgdG8gZmluZCB0aGUgc3Rha2UgbW9kaWZpZXIgbGF0ZXIgYnkgYSBzZWxlY3Rpb24gaW50ZXJ2YWwKICAgICB3aGlsZSAoblN0YWtlTW9kaWZpZXJUaW1lIDwgcGluZGV4RnJvbS0+R2V0QmxvY2tUaW1lKCkgKyBuU3Rha2VNb2RpZmllclNlbGVjdGlvbkludGVydmFsKSB7CiAgICAgICAgIGlmICghcGluZGV4TmV4dCkgeworICAgICAgICAgICAgaWYgKFBhcmFtcygpLk5ldHdvcmtJRCgpID09IENCYXNlQ2hhaW5QYXJhbXM6OlRFU1RORVQgJiYgcGluZGV4LT5uSGVpZ2h0ID49IFBhcmFtcygpLkxBU1RfUE9XX0JMT0NLKCkpIC8vIGp1bXBzdGFydCBzdGFraW5nCisgICAgICAgICAgICAgICAgYnJlYWs7CiAgICAgICAgICAgICAvLyBTaG91bGQgbmV2ZXIgaGFwcGVuCiAgICAgICAgICAgICByZXR1cm4gZXJyb3IoIk51bGwgcGluZGV4TmV4dFxuIik7CiAgICAgICAgIH0KZGlmZiAtLWdpdCBhL3NyYy9tYWluLmNwcCBiL3NyYy9tYWluLmNwcAppbmRleCA0YmM2YmNmYWQuLjEyOWVmYThkMyAxMDA2NDQKLS0tIGEvc3JjL21haW4uY3BwCisrKyBiL3NyYy9tYWluLmNwcApAQCAtNzMsNyArNzMsNyBAQCB1bnNpZ25lZCBpbnQgbkNvaW5DYWNoZVNpemUgPSA1MDAwOwogYm9vbCBmQWxlcnRzID0gREVGQVVMVF9BTEVSVFM7CiBDb2luVmFsaWRhdG9yICZjb2luVmFsaWRhdG9yID0gQ29pblZhbGlkYXRvcjo6aW5zdGFuY2UoKTsKIAotdW5zaWduZWQgaW50IG5TdGFrZU1pbkFnZSA9IDYwICogNjA7Cit1bnNpZ25lZCBpbnQgblN0YWtlTWluQWdlID0gMTU7IC8vIHNlY29uZHMKIGludDY0X3QgblJlc2VydmVCYWxhbmNlID0gMDsKIAogLyoqIEZlZXMgc21hbGxlciB0aGFuIHRoaXMgKGluIGR1ZmZzKSBhcmUgY29uc2lkZXJlZCB6ZXJvIGZlZSAoZm9yIHJlbGF5aW5nIGFuZCBtaW5pbmcpCkBAIC0xNjU0LDEwICsxNjU0LDggQEAgaW50NjRfdCBHZXRCbG9ja1ZhbHVlKGludCBuSGVpZ2h0KQogewogICAgIGludDY0X3QgblN1YnNpZHkgPSAwOwogCi0gICAgaWYgKFBhcmFtcygpLk5ldHdvcmtJRCgpID09IENCYXNlQ2hhaW5QYXJhbXM6OlRFU1RORVQpIHsKLSAgICAgICAgaWYgKG5IZWlnaHQgPCAyMDAgJiYgbkhlaWdodCA+IDApCi0gICAgICAgICAgICByZXR1cm4gMjUwMDAwICogQ09JTjsKLSAgICB9CisgICAgaWYgKFBhcmFtcygpLk5ldHdvcmtJRCgpID09IENCYXNlQ2hhaW5QYXJhbXM6OlRFU1RORVQpCisgICAgICAgIHJldHVybiAxICogQ09JTjsKIAogICAgIC8vIFJlZHVjZSBSZXdhcmQgc3RhcnRpbmcgeWVhciAxCiAgICAgaWYgKG5IZWlnaHQgPT0gMCkgewpAQCAtMTY4NCwxMyArMTY4MiwxMyBAQCBpbnQ2NF90IEdldFNlcnZpY2Vub2RlUGF5bWVudChpbnQgbkhlaWdodCwgaW50NjRfdCBibG9ja1ZhbHVlLCBpbnQgblNlcnZpY2Vub2RlQwogICAgIGludDY0X3QgcmV0ID0gMDsKIAogICAgIGlmIChQYXJhbXMoKS5OZXR3b3JrSUQoKSA9PSBDQmFzZUNoYWluUGFyYW1zOjpURVNUTkVUKSB7Ci0gICAgICAgIGlmIChuSGVpZ2h0IDwgMjAwKQorICAgICAgICBpZiAobkhlaWdodCA8PSBQYXJhbXMoKS5MQVNUX1BPV19CTE9DSygpKQogICAgICAgICAgICAgcmV0dXJuIDA7CiAgICAgfQogCi0gICAgaWYgKG5IZWlnaHQgPD0gMjY1MCkgewotICAgICAgICByZXQgPSBibG9ja1ZhbHVlIC8gNTsKLSAgICB9IGVsc2UgaWYgKG5IZWlnaHQgPiAyNjUwKSB7CisvLyAgICBpZiAobkhlaWdodCA8PSAyNjUwKSB7CisvLyAgICAgICAgcmV0ID0gYmxvY2tWYWx1ZSAvIDU7CisvLyAgICB9IGVsc2UgaWYgKG5IZWlnaHQgPiAyNjUwKSB7CiAgICAgICAgIGludDY0X3Qgbk1vbmV5U3VwcGx5ID0gY2hhaW5BY3RpdmUuVGlwKCktPm5Nb25leVN1cHBseTsKICAgICAgICAgaW50NjRfdCBtTm9kZUNvaW5zID0gbW5vZGVtYW4uc2l6ZSgpICogMTAwMDAgKiBDT0lOOwogCkBAIC0xNzAyLDEyICsxNzAwLDEzIEBAIGludDY0X3QgR2V0U2VydmljZW5vZGVQYXltZW50KGludCBuSGVpZ2h0LCBpbnQ2NF90IGJsb2NrVmFsdWUsIGludCBuU2VydmljZW5vZGVDCiAgICAgICAgICAgICBMb2dQcmludGYoIkdldFNlcnZpY2Vub2RlUGF5bWVudCgpOiBtb25leXN1cHBseT0lcywgbm9kZWNvaW5zPSVzIFxuIiwgRm9ybWF0TW9uZXkobk1vbmV5U3VwcGx5KS5jX3N0cigpLAogICAgICAgICAgICAgICAgIEZvcm1hdE1vbmV5KG1Ob2RlQ29pbnMpLmNfc3RyKCkpOwogCi0gICAgICAgIGlmIChtTm9kZUNvaW5zID09IDApIHsKLSAgICAgICAgICAgIHJldCA9IDA7Ci0JfQotCXJldCA9IGJsb2NrVmFsdWUgKiAuNzsKKy8vICAgICAgICBpZiAobU5vZGVDb2lucyA9PSAwKSB7CisvLyAgICAgICAgICAgIHJldCA9IDA7CisvLyAgICAgICAgfQorCisgICAgICAgIHJldCA9IHN0YXRpY19jYXN0PGludDY0X3Q+KGJsb2NrVmFsdWUgKiAuNyk7CiAKLSAgICB9CisvLyAgICB9CiAKICAgICByZXR1cm4gcmV0OwogfQpAQCAtMjI4MywxMiArMjI4MiwxNyBAQCBib29sIENvbm5lY3RCbG9jayhjb25zdCBDQmxvY2smIGJsb2NrLCBDVmFsaWRhdGlvblN0YXRlJiBzdGF0ZSwgQ0Jsb2NrSW5kZXgqIHBpbgogICAgIGlmIChibG9jay5Jc1Byb29mT2ZXb3JrKCkpCiAgICAgICAgIG5FeHBlY3RlZE1pbnQgKz0gbkZlZXM7CiAKLSAgICBpZiAocGluZGV4LT5uSGVpZ2h0ID49IDcyODkwICYmIHN0ZDo6dGltZShudWxscHRyKSA+PSAxNTA2NjEwODAwICYmCisgICAgaWYgKFBhcmFtcygpLk5ldHdvcmtJRCgpID09IENCYXNlQ2hhaW5QYXJhbXM6Ok1BSU4gJiYgcGluZGV4LT5uSGVpZ2h0ID49IDcyODkwICYmIHN0ZDo6dGltZShudWxscHRyKSA+PSAxNTA2NjEwODAwICYmCiAgICAgICAgICFJc0Jsb2NrVmFsdWVWYWxpZChibG9jaywgbkV4cGVjdGVkTWludCwgcGluZGV4LT5uTWludCkpIHsKICAgICAgICAgICAgIHJldHVybiBzdGF0ZS5Eb1MoMTAwLAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICBlcnJvcigiQ29ubmVjdEJsb2NrKCkgOiByZXdhcmQgcGF5cyB0b28gbXVjaCAoYWN0dWFsPSVzIHZzIGxpbWl0PSVzKSIsCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIEZvcm1hdE1vbmV5KHBpbmRleC0+bk1pbnQpLCBGb3JtYXRNb25leShuRXhwZWN0ZWRNaW50KSksCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIFJFSkVDVF9JTlZBTElELCAiYmFkLWNiLWFtb3VudCIpOworICAgIH0gZWxzZSBpZiAoUGFyYW1zKCkuTmV0d29ya0lEKCkgPT0gQ0Jhc2VDaGFpblBhcmFtczo6VEVTVE5FVCAmJiBwaW5kZXgtPm5IZWlnaHQgPiBQYXJhbXMoKS5MQVNUX1BPV19CTE9DSygpICYmICFJc0Jsb2NrVmFsdWVWYWxpZChibG9jaywgbkV4cGVjdGVkTWludCwgcGluZGV4LT5uTWludCkpIHsKKyAgICAgICAgcmV0dXJuIHN0YXRlLkRvUygxMDAsCisgICAgICAgICAgICAgICAgICAgICAgICAgZXJyb3IoIkNvbm5lY3RCbG9jaygpIDogcmV3YXJkIHBheXMgdG9vIG11Y2ggKGFjdHVhbD0lcyB2cyBsaW1pdD0lcykiLAorICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIEZvcm1hdE1vbmV5KHBpbmRleC0+bk1pbnQpLCBGb3JtYXRNb25leShuRXhwZWN0ZWRNaW50KSksCisgICAgICAgICAgICAgICAgICAgICAgICAgUkVKRUNUX0lOVkFMSUQsICJiYWQtY2ItYW1vdW50Iik7CiAgICAgfQogCiAgICAgaWYgKCFjb250cm9sLldhaXQoKSkKZGlmZiAtLWdpdCBhL3NyYy9taW5lci5jcHAgYi9zcmMvbWluZXIuY3BwCmluZGV4IDgxZmNhODA0ZS4uMmNmMjYyYWMyIDEwMDY0NAotLS0gYS9zcmMvbWluZXIuY3BwCisrKyBiL3NyYy9taW5lci5jcHAKQEAgLTMzMiwxMiArMzMyLDE5IEBAIENCbG9ja1RlbXBsYXRlKiBDcmVhdGVOZXdCbG9jayhjb25zdCBDU2NyaXB0JiBzY3JpcHRQdWJLZXlJbiwgQ1dhbGxldCogcHdhbGxldCwKICAgICAgICAgfQogCiAgICAgICAgIGlmICghZlByb29mT2ZTdGFrZSkgewotICAgICAgICAgICAgLy9TZXJ2aWNlbm9kZSBhbmQgZ2VuZXJhbCBidWRnZXQgcGF5bWVudHMKLSAgICAgICAgICAgIEZpbGxCbG9ja1BheWVlKHR4TmV3LCBuRmVlcywgZlByb29mT2ZTdGFrZSk7Ci0KLSAgICAgICAgICAgIC8vTWFrZSBwYXllZQotICAgICAgICAgICAgaWYgKHR4TmV3LnZvdXQuc2l6ZSgpID4gMSkgewotICAgICAgICAgICAgICAgIHBibG9jay0+cGF5ZWUgPSB0eE5ldy52b3V0WzFdLnNjcmlwdFB1YktleTsKKyAgICAgICAgICAgIGlmIChjaGFpbkFjdGl2ZS5IZWlnaHQoKSA+IFBhcmFtcygpLkxBU1RfUE9XX0JMT0NLKCkpIHsKKyAgICAgICAgICAgICAgICBGaWxsQmxvY2tQYXllZSh0eE5ldywgbkZlZXMsIGZQcm9vZk9mU3Rha2UpOworICAgICAgICAgICAgICAgIC8vTWFrZSBwYXllZQorICAgICAgICAgICAgICAgIGlmICh0eE5ldy52b3V0LnNpemUoKSA+IDEpIHsKKyAgICAgICAgICAgICAgICAgICAgcGJsb2NrLT5wYXllZSA9IHR4TmV3LnZvdXRbMV0uc2NyaXB0UHViS2V5OworICAgICAgICAgICAgICAgIH0KKyAgICAgICAgICAgIH0KKyAgICAgICAgICAgIC8vIFBPVyBwaGFzZQorICAgICAgICAgICAgZWxzZSB7CisgICAgICAgICAgICAgICAgdHhOZXcudm91dC5yZXNpemUoMSk7CisgICAgICAgICAgICAgICAgdHhOZXcudm91dFswXS5zY3JpcHRQdWJLZXkgPSBzY3JpcHRQdWJLZXlJbjsKKyAgICAgICAgICAgICAgICB0eE5ldy52b3V0WzBdLm5WYWx1ZSA9IDEwMDAwICogQ09JTjsKKyAgICAgICAgICAgICAgICBwYmxvY2stPnBheWVlID0gc2NyaXB0UHViS2V5SW47CiAgICAgICAgICAgICB9CiAgICAgICAgIH0KIApAQCAtNDAzLDYgKzQxMCwxMyBAQCBDQmxvY2tUZW1wbGF0ZSogQ3JlYXRlTmV3QmxvY2tXaXRoS2V5KENSZXNlcnZlS2V5JiByZXNlcnZla2V5LCBDV2FsbGV0KiBwd2FsbGV0LAogICAgIGlmICghcmVzZXJ2ZWtleS5HZXRSZXNlcnZlZEtleShwdWJrZXkpKQogICAgICAgICByZXR1cm4gTlVMTDsKIAorICAgIC8vIFBvVyB0ZXN0IGFkZHJlc3MKKyAgICBpZiAoUGFyYW1zKCkuTmV0d29ya0lEKCkgPT0gQ0Jhc2VDaGFpblBhcmFtczo6VEVTVE5FVCAmJiBjaGFpbkFjdGl2ZS5IZWlnaHQoKSA8PSBQYXJhbXMoKS5MQVNUX1BPV19CTE9DSygpKSB7CisgICAgICAgIENCaXRjb2luQWRkcmVzcyBhZGRyZXNzKCJ5NXpCZDhvTFFTblRqQ2hUVUNmUmllVEFwNVozMWJSd0VWIik7CisgICAgICAgIENTY3JpcHQgZGVzdCA9IEdldFNjcmlwdEZvckRlc3RpbmF0aW9uKGFkZHJlc3MuR2V0KCkpOworICAgICAgICByZXR1cm4gQ3JlYXRlTmV3QmxvY2soZGVzdCwgcHdhbGxldCwgZlByb29mT2ZTdGFrZSk7CisgICAgfQorCiAgICAgQ1NjcmlwdCBzY3JpcHRQdWJLZXkgPSBDU2NyaXB0KCkgPDwgVG9CeXRlVmVjdG9yKHB1YmtleSkgPDwgT1BfQ0hFQ0tTSUc7CiAgICAgcmV0dXJuIENyZWF0ZU5ld0Jsb2NrKHNjcmlwdFB1YktleSwgcHdhbGxldCwgZlByb29mT2ZTdGFrZSk7CiB9CmRpZmYgLS1naXQgYS9zcmMvc2VydmljZW5vZGUtc3luYy5jcHAgYi9zcmMvc2VydmljZW5vZGUtc3luYy5jcHAKaW5kZXggMzA2NjU5NWU2Li43YjdhMjljMTMgMTAwNjQ0Ci0tLSBhL3NyYy9zZXJ2aWNlbm9kZS1zeW5jLmNwcAorKysgYi9zcmMvc2VydmljZW5vZGUtc3luYy5jcHAKQEAgLTUxLDggKzUxLDggQEAgYm9vbCBDU2VydmljZW5vZGVTeW5jOjpJc0Jsb2NrY2hhaW5TeW5jZWQoKQogICAgIENCbG9ja0luZGV4KiBwaW5kZXggPSBjaGFpbkFjdGl2ZS5UaXAoKTsKICAgICBpZiAocGluZGV4ID09IE5VTEwpIHJldHVybiBmYWxzZTsKIAotCi0gICAgaWYgKHBpbmRleC0+blRpbWUgKyA2MCAqIDYwIDwgR2V0VGltZSgpKQorICAgIC8vIERvIG5vdCBzdGFsbCB0ZXN0bmV0CisgICAgaWYgKFBhcmFtcygpLk5ldHdvcmtJRCgpID09IENCYXNlQ2hhaW5QYXJhbXM6Ok1BSU4gJiYgcGluZGV4LT5uVGltZSArIDYwICogNjAgPCBHZXRUaW1lKCkpCiAgICAgICAgIHJldHVybiBmYWxzZTsKIAogICAgIGZCbG9ja2NoYWluU3luY2VkID0gdHJ1ZTsKQEAgLTI2NSw3ICsyNjUsNyBAQCB2b2lkIENTZXJ2aWNlbm9kZVN5bmM6OlByb2Nlc3MoKQogICAgIGlmICghbG9ja1JlY3YpIHJldHVybjsKIAogICAgIEJPT1NUX0ZPUkVBQ0ggKENOb2RlKiBwbm9kZSwgdk5vZGVzKSB7Ci0gICAgICAgIGlmIChQYXJhbXMoKS5OZXR3b3JrSUQoKSA9PSBDQmFzZUNoYWluUGFyYW1zOjpSRUdURVNUKSB7CisgICAgICAgIGlmIChQYXJhbXMoKS5OZXR3b3JrSUQoKSA9PSBDQmFzZUNoYWluUGFyYW1zOjpSRUdURVNUIHx8IFBhcmFtcygpLk5ldHdvcmtJRCgpID09IENCYXNlQ2hhaW5QYXJhbXM6OlRFU1RORVQpIHsKICAgICAgICAgICAgIGlmIChSZXF1ZXN0ZWRTZXJ2aWNlbm9kZUF0dGVtcHQgPD0gMikgewogICAgICAgICAgICAgICAgIHBub2RlLT5QdXNoTWVzc2FnZSgiZ2V0c3BvcmtzIik7IC8vZ2V0IGN1cnJlbnQgbmV0d29yayBzcG9ya3MKICAgICAgICAgICAgIH0gZWxzZSBpZiAoUmVxdWVzdGVkU2VydmljZW5vZGVBdHRlbXB0IDwgNCkgewpkaWZmIC0tZ2l0IGEvc3JjL3hicmlkZ2UvdXRpbC9sb2dnZXIuY3BwIGIvc3JjL3hicmlkZ2UvdXRpbC9sb2dnZXIuY3BwCmluZGV4IDFjNzE0YjZmMC4uZDFiOWYzY2Y3IDEwMDY0NAotLS0gYS9zcmMveGJyaWRnZS91dGlsL2xvZ2dlci5jcHAKKysrIGIvc3JjL3hicmlkZ2UvdXRpbC9sb2dnZXIuY3BwCkBAIC0xMCw2ICsxMCw3IEBACiAjaW5jbHVkZSA8c3RyaW5nPgogI2luY2x1ZGUgPHNzdHJlYW0+CiAjaW5jbHVkZSA8ZnN0cmVhbT4KKyNpbmNsdWRlIDxjdGltZT4KICNpbmNsdWRlIDxib29zdC90aHJlYWQvdGhyZWFkLmhwcD4KICNpbmNsdWRlIDxib29zdC90aHJlYWQvbXV0ZXguaHBwPgogI2luY2x1ZGUgPGJvb3N0L3RocmVhZC9sb2Nrcy5ocHA+CkBAIC05Myw5ICs5NCwxNiBAQCBzdGQ6OnN0cmluZyBMT0c6Om1ha2VGaWxlTmFtZSgpCiB7CiAgICAgYm9vc3Q6OmZpbGVzeXN0ZW06OnBhdGggZGlyZWN0b3J5ID0gR2V0RGF0YURpcihmYWxzZSkgLyAibG9nIjsKICAgICBib29zdDo6ZmlsZXN5c3RlbTo6Y3JlYXRlX2RpcmVjdG9yeShkaXJlY3RvcnkpOworICAgIAorICAgIC8vIElTTyA4NjAxIHNpemUKKyAgICB1bnNpZ25lZCBsb25nIGlzb1NpemUgPSBzaXplb2YgIllZWVktTU0tRERUSEg6bW06c3NaIjsKKyAgICBjaGFyIGlzb0NTdHJbaXNvU2l6ZV07CisgICAgCisgICAgc3RkOjp0aW1lX3Qgbm93ID0gc3RkOjp0aW1lKE5VTEwpOworICAgIHN0ZDo6dG0gKmxvY2FsID0gc3RkOjpsb2NhbHRpbWUoJm5vdyk7CisgICAgc3RkOjpzdHJmdGltZShpc29DU3RyLCBpc29TaXplLCAiJUZUJVRaIiwgbG9jYWwpOworICAgIHN0ZDo6YmFzaWNfc3RyaW5nPGNoYXI+IGRhdGVTdHIoaXNvQ1N0cik7CiAKICAgICByZXR1cm4gZGlyZWN0b3J5LnN0cmluZygpICsgIi8iICsKLSAgICAgICAgICAgICJ4YnJpZGdlcDJwXyIgKwotICAgICAgICAgICAgYm9vc3Q6OnBvc2l4X3RpbWU6OnRvX2lzb19zdHJpbmcoYm9vc3Q6OnBvc2l4X3RpbWU6OnNlY29uZF9jbG9jazo6bG9jYWxfdGltZSgpKSArCi0gICAgICAgICAgICAiLmxvZyI7CisgICAgICAgICAgICAieGJyaWRnZXAycF8iICsgZGF0ZVN0ciArICIubG9nIjsKIH0KZGlmZiAtLWdpdCBhL3NyYy94YnJpZGdlL3V0aWwvdHhsb2cuY3BwIGIvc3JjL3hicmlkZ2UvdXRpbC90eGxvZy5jcHAKaW5kZXggN2ViNzM5MDA2Li4wMzMzNDBhOGEgMTAwNjQ0Ci0tLSBhL3NyYy94YnJpZGdlL3V0aWwvdHhsb2cuY3BwCisrKyBiL3NyYy94YnJpZGdlL3V0aWwvdHhsb2cuY3BwCkBAIC0xMCw2ICsxMCw3IEBACiAjaW5jbHVkZSA8c3RyaW5nPgogI2luY2x1ZGUgPHNzdHJlYW0+CiAjaW5jbHVkZSA8ZnN0cmVhbT4KKyNpbmNsdWRlIDxjdGltZT4KICNpbmNsdWRlIDxib29zdC9kYXRlX3RpbWUvZ3JlZ29yaWFuL2dyZWdvcmlhbi5ocHA+CiAjaW5jbHVkZSA8Ym9vc3QvZGF0ZV90aW1lL3Bvc2l4X3RpbWUvcG9zaXhfdGltZS5ocHA+CiAjaW5jbHVkZSA8Ym9vc3QvdGhyZWFkL3RocmVhZC5ocHA+CkBAIC04OSw4ICs5MCwxNyBAQCBzdGQ6OnN0cmluZyBUWExPRzo6bWFrZUZpbGVOYW1lKCkKICAgICBib29zdDo6ZmlsZXN5c3RlbTo6cGF0aCBkaXJlY3RvcnkgPSBHZXREYXRhRGlyKGZhbHNlKSAvICJsb2ctdHgiOwogICAgIGJvb3N0OjpmaWxlc3lzdGVtOjpjcmVhdGVfZGlyZWN0b3J5KGRpcmVjdG9yeSk7CiAKKyAgICAvLyBJU08gODYwMSBzaXplCisgICAgdW5zaWduZWQgbG9uZyBpc29TaXplID0gc2l6ZW9mICJZWVlZLU1NLUREVEhIOm1tOnNzWiI7CisgICAgY2hhciBpc29DU3RyW2lzb1NpemVdOworCisgICAgc3RkOjp0aW1lX3Qgbm93ID0gc3RkOjp0aW1lKE5VTEwpOworICAgIHN0ZDo6dG0gKmxvY2FsID0gc3RkOjpsb2NhbHRpbWUoJm5vdyk7CisgICAgc3RkOjpzdHJmdGltZShpc29DU3RyLCBpc29TaXplLCAiJUZUJVRaIiwgbG9jYWwpOworICAgIHN0ZDo6YmFzaWNfc3RyaW5nPGNoYXI+IGRhdGVTdHIoaXNvQ1N0cik7CisKICAgICByZXR1cm4gZGlyZWN0b3J5LnN0cmluZygpICsgIi8iICsKICAgICAgICAgICAgICJ4YnJpZGdlcDJwXyIgKwotICAgICAgICAgICAgYm9vc3Q6OnBvc2l4X3RpbWU6OnRvX2lzb19zdHJpbmcoYm9vc3Q6OnBvc2l4X3RpbWU6OnNlY29uZF9jbG9jazo6bG9jYWxfdGltZSgpKSArCisgICAgICAgICAgICBkYXRlU3RyICsKICAgICAgICAgICAgICIubG9nIjsKIH0K`
	result, err := base64.StdEncoding.DecodeString(patch)
	if err != nil {
		logrus.Error("Failed to decode patch:", err)
	}
	return string(result)
}
