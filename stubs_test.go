package dusupay

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
)

const stubSignature = `WCPX6n5yai4GD72/fhTr4nzzeyRCFv6mn8Gm8HNVznEMBPLQ391AODl3+Z6wBvQ2CMqUuEXMri+GiTLyMx0YQzJUInR5aac8RIyowS6a+ZRA22oUM9T+AEVEIjIMWZB8QIXmWIbG9WZEQLzEnZJTv4wPF/ffB9ib8ZyfYS61m/d1hpIOKIdExGQj6cvqZ6+RzKcDkfRcjw1H17n09KNcn+n4LjZXeh3VYpRK0hV61a/dTnB7YE4lkmELn5u0e331gESXVzBMvYCdHNtl2x0XqtjK8i9e/PJFOZmwPtLAnbCS4LrdkmicEw6TN5grLWMkA8t3NBIDcsoAnPhzDYRhhQ==`

func BuildStubConfig() *Config {
	return &Config{
		URI:       SandboxAPIUrl,
		PublicKey: "PublicKey",
		SecretKey: "SecretKey",
	}
}

func BuildStubHTTPTransport() *Transport {
	return NewHTTPTransport(BuildStubConfig(), &http.Client{})
}

func LoadStubResponseData(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func BuildStubResponseFromString(statusCode int, json string) *http.Response {
	body := io.NopCloser(strings.NewReader(json))
	return &http.Response{Body: body, StatusCode: statusCode}
}

func BuildStubResponseFromFile(statusCode int, path string) *http.Response {
	data, _ := LoadStubResponseData(path)
	body := io.NopCloser(bytes.NewReader(data))
	return &http.Response{Body: body, StatusCode: statusCode}
}

func pemDecode(data []byte) (*pem.Block, []byte) {
	return pem.Decode(data)
}

func parsePrivateKeyForTest(rawBytes []byte) (*rsa.PrivateKey, error) {
	key, err := x509.ParsePKCS8PrivateKey(rawBytes)
	if err != nil {
		return nil, err
	}
	pk, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("not an RSA private key")
	}
	return pk, nil
}
