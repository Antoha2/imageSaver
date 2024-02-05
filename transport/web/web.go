package web

import (
	"net/http"

	service "github.com/antoha2/images/service"
)

type Transport interface{}

type webImpl struct {
	service service.Service
	server  *http.Server
}

func NewWeb(service service.Service) *webImpl {
	return &webImpl{
		service: service,
	}
}

type WebImagesData struct {
	Urls  string `json:"urls"`
	Count int    `json:"count"`
}
