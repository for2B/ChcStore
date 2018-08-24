package httpapi

import (
	"fmt"
	"go.uber.org/zap"
	"net"
	"net/http"
	"strings"
)
//开启http服务线程
func Serve(listener net.Listener, handler http.Handler, logger *zap.Logger) {
	server := &http.Server{Handler: handler} //Server类型定义了运行HTTP服务端的参数。Handler: 调用的处理器，

	err := server.Serve(listener)//Serve会接收监听器收到的每一个连接，并为每一个连接创建一个新的服务go程。该go程会读取请求，然后调用srv.Handler回复请求。

	if err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
		logger.Error(fmt.Sprintf("http.Serve() - %s", err))
	}
	logger.Info(fmt.Sprintf("%s: closing", "http"),
		zap.String("port", listener.Addr().String()))
}
