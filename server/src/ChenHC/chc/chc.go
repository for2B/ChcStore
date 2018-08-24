package chc

import (
	"ChenHC/internal/httpapi"
	"ChenHC/chc/infrastructure"
	"ChenHC/chc/model"
	"ChenHC/chc/option"
	"ChenHC/chc/view"
	"ChenHC/utils"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

//remember to delete me
type CHC struct {
	httpListener net.Listener  //面向流的网络协议的公用的网络监听器接口

	sync.RWMutex  //读写锁互斥锁

	//isLoading 0->loading 1->load finish
	IsLoading int32  //??
	startTime time.Time  //服务启动的时间

	exitChan  chan struct{}  //退出管道?
	waitGroup utils.WaitGroupWrapper  //线程等待

	infrastructure *infrastructure.Infrastructure  //（基础设施）
}

func New(opts *option.Options) *CHC {
		p := &CHC{
		startTime:      time.Now(),
		exitChan:       make(chan struct{}),  //使用空结构，避免内存滥用。exitChan 仅做信号通知，没有实际价值
		infrastructure: infrastructure.NewInfrastructure(opts),
	}
	fmt.Println("初始化CHC完成")
	return p
}

func (p *CHC) Main() {
	opts := p.infrastructure.GetOpts()
	mm := model.GetModel(p.infrastructure)
	mm.InitAllTable()  //建表

	httpListener, err := net.Listen("tcp", ":"+opts.ServerPort) //监听对应端口
	if err != nil {
		p.infrastructure.Logger.Fatal(fmt.Sprintf("listen (%s) failed - %s", opts.ServerPort, err))
		os.Exit(1)
	}

	p.Lock()  //为什么这里要锁？
	p.httpListener = httpListener
	p.Unlock()
	view.Init(opts.FrontDir)  //view init
	server := newHTTPServer(&context{CHC: p})  //配置http api和对应处理函数，中间件配置
	p.waitGroup.Wrap(func() {
		httpapi.Serve(httpListener, server, p.infrastructure.Logger) //开启监听线程
	})

}

func (p *CHC) Exit() {
	if p.httpListener != nil {
		p.httpListener.Close()
	}

	p.WriteCurData()

	//TODO: ADD MORE

	close(p.exitChan)

	p.waitGroup.Wait()

	p.infrastructure.Logger.Sync()

	if p.infrastructure.DB != nil {
		p.infrastructure.DB.Close()
	}

}

func (p *CHC) LoadLastData() {
	//todo:load last data

}

func (p *CHC) WriteCurData() {
	//todo:write data before exit
}
