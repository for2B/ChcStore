package model

import (
	"ChenHC/chc/infrastructure"
	"fmt"
	"log"
	"sync"
)

type Model struct {
	*infrastructure.Infrastructure
}

var onceModel sync.Once  //只执行一次
var model *Model

func GetModel(i *infrastructure.Infrastructure) *Model {
	onceModel.Do(func() {
		model = &Model{
			Infrastructure: i,
		}
	})
	return model
}

func (m *Model) InitAllTable() {

	// 创建syslog表
	_, err := m.DB.Exec(`create table if not exists syslog(
		id SERIAL NOT NULL,
		log jsonb NOT NULL,
		PRIMARY KEY ("id")
	);`)
	if err != nil {
		log.Panicln("create table ChenHC.syslog failed: " + err.Error())
	}

	fmt.Println("init db success!")
}
