package helpers

import "sync"

type DocumentData struct {
	DocumentData map[string]float32
	Mutex sync.Mutex
}

func (dd *DocumentData) PutTermFrequency(term string, freq float32){
	dd.Mutex.Lock()
	dd.DocumentData[term] = freq
	dd.Mutex.Unlock()
}

func (dd *DocumentData) GetTermFrequency(term string)float32{
	return dd.DocumentData[term]
}




