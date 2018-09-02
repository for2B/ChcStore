package project_model

import (
	"sync"
	"ChenHC/chc/infrastructure"
)

var onceGetUploadFile sync.Once

type UploadFileModel struct {
	*infrastructure.Infrastructure
	AllowOrign string
}

var uploadfilemodel *UploadFileModel

func GetUploadfilemodel (i *infrastructure.Infrastructure,Alloworigin string) *UploadFileModel{
	onceGetUploadFile.Do(func() {
			uploadfilemodel = &UploadFileModel{
				Infrastructure:i,
				AllowOrign:Alloworigin,
			}
	})
	return uploadfilemodel
}

func (u *UploadFileModel) Upload() (string,error){
	return "",nil
}