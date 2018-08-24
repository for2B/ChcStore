package testtool

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"ChenHC/internal/httpapi"
	"ChenHC/chc/constant"
	"ChenHC/chc/infrastructure"
	"ChenHC/chc/model"
	"ChenHC/chc/option"
	"path/filepath"
	"strings"
	"testing"
)

func TestEnv() (i *infrastructure.Infrastructure) {
	opts := option.NewOpts("../../../../bin")
	i = infrastructure.NewInfrastructure(opts)
	return
}

type FB struct {
	DistWriter http.ResponseWriter `json:"-"`
	FbCode     int                 `json:"code"`
	FbMsg      string              `json:"msg,omitempty"`
	FbData     interface{}         `json:"data,omitempty"`
}

// Case 表示一个测试样例
type Case struct {
	Path     string
	Method   string            //http请求方法
	RawQuery string            //Get方法携带参数
	Data     string            //POST方法携带数据
	HTTPFunc http.HandlerFunc  //测试函数
	Target   string            //测试目的
	WantCode int               //期望返回值
	UserInfo *constant.Session //用户信息
	FeedBack *FB
	FBdata   interface{}
}

// 模拟请求
func (c *Case) MockRequest(t *testing.T) error {
	req := httptest.NewRequest(c.Method, c.Path, strings.NewReader(c.Data))
	w := httptest.NewRecorder()
	if c.RawQuery != "" {
		req.URL.RawQuery = c.RawQuery

	}
	ctx := req.Context()
	ctx = context.WithValue(ctx, "session", c.UserInfo)

	req = req.WithContext(ctx)
	c.HTTPFunc(w, req)
	//判断http status
	if status := w.Code; status != http.StatusOK {
		t.Errorf("http status is not expected: %v ", status)
		return nil
	}

	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read result body failed: %v", err)
	}
	defer resp.Body.Close()

	//skip check code
	if c.WantCode == 0 {
		return nil
	}

	if c.FeedBack == nil {
		var fb FB
		if err = json.Unmarshal(body, &fb); err != nil {
			return fmt.Errorf("unmarshal feedback data failed: %v", err)
		}

		if fb.FbCode != c.WantCode {
			return fmt.Errorf("want %d but got %d", c.WantCode, fb.FbCode)
		}
	} else {
		if err = json.Unmarshal(body, &c.FeedBack); err != nil {
			return fmt.Errorf("unmarshal feedback data failed: %v", err)
		}
		if c.FeedBack.FbCode != c.WantCode {
			return fmt.Errorf("want %d but got %d", c.WantCode, c.FeedBack.FbCode)
		}
	}

	body, err = json.Marshal(c.FeedBack.FbData)
	if err != nil {
		t.Errorf("json.Marshal(getAnswerCase.FeedBack.FbData) failed: %v", err)
	}
	data := model.Data{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		t.Errorf("json.Unmarshal(body, &data) failed:%v", err)
	}
	c.FBdata = data
	return nil
}

func NewfileUploadRequest(uri string /*params map[string]string,*/, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	/*	for key, val := range params {
		_ = writer.WriteField(key, val)
	}*/

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req := httptest.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req, nil
}

func AllowOrigin(next httpapi.APIHandler) httpapi.APIHandler {
	return func(w http.ResponseWriter, r *http.Request) (response interface{}, err error) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Set("content-type", "application/json")
		return next(w, r)
	}
}
