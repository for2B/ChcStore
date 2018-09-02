package project_controllet

import (
	"net/http"
	"ChenHC/chc/model/project_model"
	"time"
	"crypto/md5"
	"io"
	"strconv"
	"fmt"
	"os"
	"encoding/json"
)

type Link struct {
	Link string `json:"link"`
}

type UpLoadFileController struct {
	*project_model.UploadFileModel
}

func (c *UpLoadFileController) UploadFile(w http.ResponseWriter,r *http.Request) (interface{},error){

	if r.Method == "POST"{
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h,strconv.FormatInt(crutime,10))
		token:=fmt.Sprintf("%x",h.Sum(nil))
		fmt.Println(token)
		//r.ParseMultipartForm(32<<20)
		file,handler,err:=r.FormFile("filename")
		if err!=nil{
			fmt.Println("get r.FormFile failed",err)
			return nil,err
		}
		defer file.Close()
		fmt.Fprint(w,"%v",handler.Header)
			f,err := os.OpenFile("bin/project_item_src/"+token+handler.Filename,os.O_WRONLY|os.O_APPEND|os.O_CREATE,0666)
			if err!=nil{
				fmt.Println("OpenFile failed",err)
				return nil,err
			}
			defer f.Close()
			_,err = io.Copy(f,file)
			if err!=nil{
				fmt.Println("copy file failed",err)
			return nil,err
		}
		lk := Link{Link:"http://chl.ish2b.cn/ospf/"}
		buf, err := json.Marshal(lk)
		if err!=nil{
			fmt.Println("json Marshal faild",err)
			return nil,err
		}
		fmt.Println(string(buf))
		fmt.Fprint(w, string(buf))
	}
	return c.Upload()
}

