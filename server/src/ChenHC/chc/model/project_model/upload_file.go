package project_model

import (
	"sync"
	"ChenHC/chc/infrastructure"
)

var onceGetUploadFile sync.Once

type UploadFileModel struct {
	*infrastructure.Infrastructure
	FilePath string
}

var uploadfilemodel *UploadFileModel

func GetUploadfilemodel (i *infrastructure.Infrastructure,filepath string) *UploadFileModel{
	onceGetUploadFile.Do(func() {
			uploadfilemodel = &UploadFileModel{
				Infrastructure:i,
				FilePath:filepath,
			}
	})
	return uploadfilemodel
}

func (u *UploadFileModel) Upload() (string,error){
	return "",nil
}