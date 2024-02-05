package main

import (
	"os"
	"os/signal"
	"syscall"

	repository "github.com/antoha2/images/repository"
	service "github.com/antoha2/images/service"
	web "github.com/antoha2/images/transport/web"
)

func main() {

	Run()

}

func Run() {

	ImgRepository := repository.NewRepository()
	ImgService := service.NewService(ImgRepository)
	ImgTransport := web.NewWeb(ImgService)

	go ImgTransport.StartHTTP()

	// go ImgTransport.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit
	ImgTransport.Stop()
}
