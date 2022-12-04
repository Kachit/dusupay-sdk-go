package dusupay

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
)

const stubSignature = `gYC3u1wUtk6UFpOVvCx+AyCnE3LXkS9Sg74fiRUQxRDDlllPu5vuRUrEbEqq/TEO90fYr76KGAWC6YSo
J9joYwk8RVftDQ1pNhROdfRkXL/yaQbrAvuT2gM2sO+HJhShCBLWbBcfXPOcjGcCedPCHSNFc5bq/Mk/
DszqlFEoH0dUN8hmqXQr673zyFivaKT76CpJTcmn5nvJi8r6IGoOJXb5uN8CdMTXbT6J08OmsbILPlfX
qe8PS/IlvXz11oy5xUaLXt+whhZL8rBrwQUsi9aNVf8Gd5m93D2ls1z03zDSOjSlb26Rvnvk97+XSM13
KuGbYjc3eJ6CUlQuIbTC1A==`

func BuildStubConfig() *Config {
	return &Config{
		Uri:       SandboxAPIUrl,
		PublicKey: "PublicKey",
		SecretKey: "SecretKey",
	}
}

func BuildStubHttpTransport() *Transport {
	return NewHttpTransport(BuildStubConfig(), &http.Client{})
}

func LoadStubResponseData(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

func BuildStubResponseFromString(statusCode int, json string) *http.Response {
	body := ioutil.NopCloser(strings.NewReader(json))
	return &http.Response{Body: body, StatusCode: statusCode}
}

func BuildStubResponseFromFile(statusCode int, path string) *http.Response {
	data, _ := LoadStubResponseData(path)
	body := ioutil.NopCloser(bytes.NewReader(data))
	return &http.Response{Body: body, StatusCode: statusCode}
}
