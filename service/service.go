package service

import (
	"context"
	"sync"

	repository "github.com/antoha2/images/repository"
)

type Service interface {
	ImageDownload(ctx context.Context, data []ServImagesData) (string, error)
	RunDownload(ctx context.Context, imgData ServImagesData, dirName string, c chan error, wg *sync.WaitGroup) error
	PipelineDownload(ctx context.Context, imgData *ServImagesData, dirName string, wg *sync.WaitGroup) error
}

type serviceImpl struct {
	rep repository.Repository
}

func NewService(rep repository.Repository) *serviceImpl {
	return &serviceImpl{}
}

type ServImagesData struct {
	Urls  string
	Count int
}

type History struct {
	ImgHisMap map[string]int
	mu        *sync.Mutex
}
