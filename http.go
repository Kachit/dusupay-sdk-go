package dusupay

import "net/http"

//Transport wrapper
type Transport struct {
	http *http.Client
}

//Response wrapper
type Response struct {
	raw *http.Response
}
