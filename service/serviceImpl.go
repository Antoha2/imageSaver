package service

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	config "github.com/antoha2/images/config"
	"github.com/antoha2/images/helpers"
)

//

func (s *serviceImpl) ImageDownload(ctx context.Context, data []ServImagesData) (string, error) {

	imgHistory := helpers.NewHistory()
	ctx = context.WithValue(context.Background(), config.ImgMap, imgHistory)

	dt := time.Now()
	dirName1 := fmt.Sprintf("%s%s", config.ImgPath, dt.Format("01-02-2006 15-04-05.00"))
	err := os.Mkdir(dirName1, 0777)
	if err != nil {
		log.Println("os.Mkdir() err - ", err)
		return "", err
	}

	wg := new(sync.WaitGroup)
	wg.Add(len(data))

	log.Println("старт загрузки")
	errorChan := make(chan error)

	for index, imgData := range data {

		dirName2 := fmt.Sprintf("%s\\%d", dirName1, index+1)
		go s.RunDownload(ctx, imgData, dirName2, errorChan, wg)

	}

	for index := 0; index < len(data); index++ {

		errorRunDownload := <-errorChan
		if errorRunDownload != nil {
			log.Printf("%d ошибка загрузки - %s\n", index+1, errorRunDownload)
			continue
		}
	}
	wg.Wait()
	close(errorChan)

	response := imgHistory.Get()
	//log.Println(response)
	return response, nil
}

//скачиваник файлов
func (s *serviceImpl) RunDownload(ctx context.Context, imgData ServImagesData, dirName string, c chan error, wg *sync.WaitGroup) error {

	defer wg.Done()

	err := s.ExamImg(ctx, imgData.Urls) //imgFmt
	if err != nil {
		log.Println("s.ExamImg() err - ", err)
		c <- err
		return err
	}

	err = os.Mkdir(dirName, 0777)
	if err != nil {
		log.Println("os.Mkdir() err - ", err)
		c <- err
		return err
	}

	wgDwld := new(sync.WaitGroup)
	wgDwld.Add(imgData.Count)
	for index := 0; index < imgData.Count; index++ {

		imgName := fmt.Sprintf("%s\\(%d)", dirName, index+1)
		go s.PipelineDownload(ctx, &imgData, imgName, wgDwld) //
	}
	wgDwld.Wait()
	log.Printf("скачивание прошло успешно - %s \n", imgData.Urls)
	c <- err
	return nil
}

func (s *serviceImpl) PipelineDownload(ctx context.Context, imgData *ServImagesData, dirName string, wgDwld *sync.WaitGroup) error {

	defer wgDwld.Done()

	fileName := path.Base(imgData.Urls)
	res, err := http.Get(imgData.Urls)
	if err != nil {
		log.Println("http.Get() err - ", err)
		return err
	}
	// Задержка операции после defer, обычно используется для освобождения связанных переменных

	defer res.Body.Close()
	reader := bufio.NewReader(res.Body)
	file, err := os.Create(dirName + fileName)
	if err != nil {
		log.Println("os.Create() err - ", err)
		return err
	}

	// Получить объект записи файла

	writer := bufio.NewWriter(file)
	_, err = io.Copy(writer, reader)
	if err != nil {
		log.Println("io.Copy() err - ", err)
		return err
	}

	h := ctx.Value(config.ImgMap).(*helpers.History)
	h.Increment(imgData.Urls)
	return nil
}

//проверка соответсвия форматов изображения
func (s *serviceImpl) ExamImg(ctx context.Context, url string) error {

	for _, format := range config.ImgFormats {
		if contain := strings.Contains(strings.ToLower(url), format); contain { //сторока с конца
			return nil
		}
	}
	return errors.New("неверный формат файла - " + url)

}
