package httpparse

import (
	"context"
	"errors"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

func Parse(r *http.Request) (context.Context, error) {
	Form := make(map[string][]string)
	MultipartFormFile := make(map[string][]*multipart.FileHeader)
	MultipartFormValue := make(map[string][]string)
	ctx := r.Context()

	r.ParseMultipartForm(1 << 32)
	r.ParseForm()

	if len(r.Form) > 0 {
		for k, v := range r.Form {
			Form[k] = []string(v)
		}
	}

	if r.MultipartForm != nil && len(r.MultipartForm.File) > 0 {
		for k, v := range r.MultipartForm.File {
			MultipartFormFile[k] = []*multipart.FileHeader(v)
		}
	}

	if r.MultipartForm != nil && len(r.MultipartForm.Value) > 0 {
		for k, v := range r.MultipartForm.Value {
			MultipartFormValue[k] = []string(v)
		}
	}

	buf, _ := ioutil.ReadAll(r.Body)
	ctx = context.WithValue(ctx, "Body", string(buf))
	ctx = context.WithValue(ctx, "Form", Form)
	ctx = context.WithValue(ctx, "MultipartFormFile", MultipartFormFile)
	ctx = context.WithValue(ctx, "MultipartFormValue", MultipartFormValue)

	return ctx, nil
}

func GetBody(r *http.Request) (s string, err error) {
	body, ok := r.Context().Value("Body").(string)
	if !ok {
		return "", errors.New("get body fail : body is not string")
	}

	return body, nil
}

func GetForm(r *http.Request) (form map[string][]string, err error) {
	form, isOK := r.Context().Value("Form").(map[string][]string)
	if !isOK {
		return nil, errors.New("get form fail : form is not map[string][]string")
	}

	return form, nil
}

func GetMultipartFormFile(r *http.Request) (form map[string][]*multipart.FileHeader, err error) {
	form, ok := r.Context().Value("MultipartFormFile").(map[string][]*multipart.FileHeader)
	if !ok {
		return nil, errors.New("get multipart form file fail : multipart form file is not map[string][]*multipart.FileHeader")
	}

	return form, nil
}

func GetMultipartFormValue(r *http.Request) (form map[string][]string, err error) {
	form, ok := r.Context().Value("MultipartFormValue").(map[string][]string)
	if !ok {
		return nil, errors.New("get multipart form value fail : multipart form value is not map[string][]string")
	}

	return form, nil
}
