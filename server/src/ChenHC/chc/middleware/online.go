package middleware

import (
	"net/http"
	"ChenHC/internal/httpapi"
	"sync"
)

func Online(next httpapi.APIHandler) httpapi.APIHandler {
	recordOnlineCount := 0
	return func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
		//todo: do what you want to do
		lock := new(sync.Mutex)
		lock.Lock()
		recordOnlineCount++
		// fmt.Println("users online:" + strconv.Itoa(recordOnlineCount))
		lock.Unlock()
		response, err := next(w, r)
		lock.Lock()
		recordOnlineCount--
		// fmt.Println("users online:" + strconv.Itoa(recordOnlineCount))
		lock.Unlock()
		return response, err
	}
}
