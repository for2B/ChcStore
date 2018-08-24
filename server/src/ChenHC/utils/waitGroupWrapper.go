package utils

import (
	"sync"
)
//WaitGroupWrapper == sync.WaitGroup
type WaitGroupWrapper struct {
	sync.WaitGroup
}
//Wrap 开启一个goroutine
func (w *WaitGroupWrapper) Wrap(cb func()) {
	w.Add(1)
	go func() {
		cb()
		w.Done()
	}()
}
