package helpers

import (
	"fmt"
	"sync"
)

type History struct {
	ImgHisMap map[string]int
	mu        sync.Mutex
}

func NewHistory() *History {
	return &History{ImgHisMap: make(map[string]int)}
}

func (h *History) Increment(url string) { // блокирует mux , наращивает count , разблокирует mux

	h.mu.Lock()
	h.ImgHisMap[url]++
	h.mu.Unlock()
}

func (h *History) Get() string { // формирование response

	var logImg string
	index := 1
	for key, value := range h.ImgHisMap {
		logImg = logImg + fmt.Sprintf("\n%d - %v - %d", index, key, value)
		index++
	}
	return logImg
}
