package endpoints

import (
	"context"
	"log"

	service "github.com/antoha2/images/service"
	"github.com/go-kit/kit/endpoint"
)

func MakeDownloadEndpoints(s service.Service) *Endpoints {
	return &Endpoints{
		Download: MakeDownloadEndpoint(s),
	}
}

type Endpoints struct {
	Download endpoint.Endpoint
}

type DownloadRequest struct {
	ImgData []ReqImgData `json:"imgs"`
}

type ReqImgData struct {
	Url   string `json:"urls"`
	Count int    `json:"count"`
}

type DownloadResponse struct {
	Resp string `json:"response"`
	//ImgMap map[string]int `json:"imgMap"`
}

func MakeDownloadEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		log.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!! - ")
		webImgData := request.(DownloadRequest)
		servImagesData := make([]service.ServImagesData, len(webImgData.ImgData))

		for index, data := range webImgData.ImgData {
			servImagesData[index].Urls = data.Url
			servImagesData[index].Count = data.Count
		}

		response, err := s.ImageDownload(ctx, servImagesData)
		log.Println(response)
		if err != nil {
			return DownloadResponse{Resp: response}, err
		}

		return DownloadResponse{Resp: response}, nil
	}
}

/* func decodeDownloadResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var resp DownloadResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
} */

/* func encodeDownloadRequest(_ context.Context, req *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.ContentLength = int64(len(buf.Bytes()))
	req.Body = ioutil.NopCloser(&buf)
	return nil
}

func copyUrl(base *url.URL, path string) *url.URL {
	next := *base
	next.Path = path
	return &next
}
*/
