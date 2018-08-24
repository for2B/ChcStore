package controller

import (
	"net/http"
	"ChenHC/chc/model"
)

type LogController struct {
	*model.LogModel
}

func (c *LogController) GetLogs(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	//	c.LogModel.Logger.Info("")
	return c.GetLog()
}
