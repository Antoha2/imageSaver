package web

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/antoha2/images/config"
	endpoints "github.com/antoha2/images/transport/web/endpoints"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func (wImpl *webImpl) StartHTTP() error {

	Options := []httptransport.ServerOption{
		//httptransport.ServerBefore(wImpl.UserIdentify),
	}

	DownloadHandler := httptransport.NewServer(
		endpoints.MakeDownloadEndpoint(wImpl.service),
		decodeMakeDownloadRequest,
		encodeResponse,
		Options...,
	)

	r := mux.NewRouter()
	r.Methods("POST").Path("/img_download").Handler(DownloadHandler)

	wImpl.server = &http.Server{Addr: config.HostAddr}                            //:8180
	log.Printf(" Запуск HTTP-сервера на http://127.0.0.1%s\n", wImpl.server.Addr) //:8180

	if err := http.ListenAndServe(wImpl.server.Addr, r); err != nil {
		log.Println(err)
	}

	return nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

func decodeMakeDownloadRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var request endpoints.DownloadRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func (wImpl *webImpl) Stop() {

	if err := wImpl.server.Shutdown(context.TODO()); err != nil {
		panic(err)
	}
}

/*
import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/antoha2/images/config"
	service "github.com/antoha2/images/service"
)

func (wImpl *webImpl) Start() error {

	wImpl.server = &http.Server{Addr: ":8180"}

	mux := http.NewServeMux()
	mux.HandleFunc("/img_download", wImpl.handlerImagesDownload)

	log.Printf("Запуск веб-сервера img_download на http://127.0.0.1:%s\n", wImpl.server.Addr) //:8180
	http.ListenAndServe(wImpl.server.Addr, mux)

	return nil
}



//декодеры JSON
func (wImpl *webImpl) Decoder(r *http.Request, data *[]WebImagesData) error {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, data)
	if err != nil {
		fmt.Println("can't unmarshal: ", err.Error())
		return err
	}
	return nil
}

//обработчик ImagesDownload
func (wImpl *webImpl) handlerImagesDownload(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		return
	}
	WebInput := new([]WebImagesData)
	err := wImpl.Decoder(r, WebInput)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	dataServ := make([]service.ServImagesData, len(*WebInput))
	for i, count := range *WebInput {

		data := service.ServImagesData{
			Urls:  count.Urls,
			Count: count.Count,
		}
		dataServ[i] = data
	}
	imgMapCtx := make(map[string]int, len(*WebInput))
	ctx := context.WithValue(context.Background(), config.ImgMap, imgMapCtx)
	err = wImpl.service.ImageDownload(ctx, dataServ)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	json, err := json.Marshal(imgMapCtx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(json)
}
*/
