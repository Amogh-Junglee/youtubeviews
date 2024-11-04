package transport

import "net/http"

type Transport interface {
	Increment(writer http.ResponseWriter, request *http.Request)
	Get(writer http.ResponseWriter, request *http.Request)
	GetTopVideos(writer http.ResponseWriter, request *http.Request)
}
