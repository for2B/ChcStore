package chc

import (
	"ChenHC/chc/middleware"
	"ChenHC/chc/model"
	"ChenHC/chc/view"
	"ChenHC/internal/httpapi"
	"net/http"
	"github.com/gorilla/mux"
	"ChenHC/chc/controller"
	"ChenHC/chc/controller/project_controllet"
	"ChenHC/chc/model/project_model"
	"ChenHC/utils"
	"fmt"
)

type context struct {
	CHC *CHC
}

// 拦截非前端路由
var frontEndRoutes = []string{
	"/admin/mainshow",
	"/admin/adminmanage",
}

type httpServer struct {
	ctx               *context
	router            *mux.Router
	defaultMiddleWare httpapi.MiddlewareFunc
}

func newHTTPServer(ctx *context) *httpServer {

	router := mux.NewRouter()
	s := &httpServer{
		ctx:    ctx,  //CHC全部内容
		router: router, // mux.NewRouter()
		defaultMiddleWare: httpapi.CombinationMiddleware(middleware.Online,  //中间介函数
			middleware.DefaultDecode,
			middleware.Encode,
			middleware.Log(ctx.CHC.infrastructure.Logger)),//调用了Log函数 返回匿名的MiddlewareFunc函数，为了能够带参数
	}
	s.initRouter()
	initFrontEndMux(router)

	router.PathPrefix("/static").Handler(http.StripPrefix("/static", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		view.GzipFileServe(w, r, r.URL.Path, ctx.CHC.infrastructure.GetOpts().FrontDir)
	})))
	router.PathPrefix("/assets/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		view.GzipFileServe(w, r, r.URL.Path, ctx.CHC.infrastructure.GetOpts().FrontDir)
	})
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {  //处理所有非api接口，返回静态文件
		view.GzipServeFile(w, r, ctx.CHC.infrastructure.GetOpts().FrontDir+"/index.html")

	})
	return s
}

func (s *httpServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	req.Context()
	s.router.ServeHTTP(w, req)
}

func (s *httpServer) AllowOrigin(next httpapi.APIHandler) httpapi.APIHandler {
	return func(w http.ResponseWriter, r *http.Request) (response interface{}, err error) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Set("content-type", "application/json")
		return next(w, r)
	}
}

func initFrontEndMux(r *mux.Router) {
	for _, item := range frontEndRoutes {
		r.HandleFunc(item, view.LoadTemplate)
	}
}

func (s *httpServer) initRouter() {

	//guochuting add start-----------------


	//guochuting add end------------------


	//guozhenzhen add start* * * * * * * * * * *


	//guozhenzhen add end* * * * * * * * * * *

	//panjiawei add start+++++++++++++++++++


	//panjiawei add end+++++++++++++++++++++


	//chenhuiliang add start——  ——  ——  ——  ——  ——  ——


	//chenhuiliang add end——  ——  ——  ——  ——  ——  ——
	s.UploadFile(s.router) //上传文件
	s.DownFile(s.router) //下载文件接口
	//chencanxin add start~~~~~~~~~
	s.regLogController(s.router) //注册api和处理函数

	//chencanxin add end~~~~~~~~~~
}





//guochuting add start-----------------


//guochuting add end------------------


//guozhenzhen add start* * * * * * * * * * *


//guozhenzhen add end* * * * * * * * * * *

//panjiawei add start+++++++++++++++++++


//panjiawei add end+++++++++++++++++++++


//chenhuiliang add start——  ——  ——  ——  ——  ——  ——


//chenhuiliang add end——  ——  ——  ——  ——  ——  ——
func (s *httpServer) UploadFile(r *mux.Router){ //处理文件上传并保存
	c := &project_controllet.UpLoadFileController{
		UploadFileModel:project_model.GetUploadfilemodel(s.ctx.CHC.infrastructure,s.ctx.CHC.infrastructure.GetOpts().AllowOrigin),
	}
	r.Handle("/api/upload_file",httpapi.Decorate(c.UploadFile,s.AllowOrigin,middleware.Online,middleware.DefaultDecode,middleware.Log(s.ctx.CHC.infrastructure.Logger)))
}

func (s *httpServer) DownFile(r *mux.Router){
	path ,_:= utils.GetProDir()
	fmt.Println(path)
	http.Handle("/files/", http.StripPrefix("/files/",http.FileServer(http.Dir(path))))
}


//chencanxin add start~~~~~~~~~
func (s *httpServer) regLogController(r *mux.Router) { //配置api和对应的处理函数
	c := &controller.LogController{
		LogModel: model.GetLogModel(s.ctx.CHC.infrastructure),
	}
	r.Handle("/api/test", httpapi.Decorate(c.GetLogs, s.defaultMiddleWare))
}
//chencanxin add end~~~~~~~~~~

