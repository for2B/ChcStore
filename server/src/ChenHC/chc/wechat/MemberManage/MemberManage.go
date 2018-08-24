package MemberManage

import (
	"ChenHC/internal/httpapi"
	"ChenHC/chc/constant"
	"ChenHC/chc/wechat/context"
	"ChenHC/chc/wechat/util"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

const (
	// 对应【获取部门成员详情】
	getMemberDetailURL = "https://qyapi.weixin.qq.com/cgi-bin/user/list?access_token=%s&department_id=%s&fetch_child=%s"
	// 对应【获取部门列表】
	getDepartmentNameURL = "https://qyapi.weixin.qq.com/cgi-bin/department/list?access_token=%s&id=%s"
)

//MemberManage 成员管理
type MemberManage struct {
	*context.Context
}

//NewMemberManage 实例化授权信息
func NewMemberManage(context *context.Context, w http.ResponseWriter, r *http.Request) *MemberManage {
	memManage := new(MemberManage)
	memManage.Context = context
	memManage.Context.Request = r
	memManage.Context.Writer = w
	return memManage
}

// ResAccessToken 获取用户授权access_token的返回结果
type ResMemberDetail struct {
	util.CommonError
	Userlist []User `json:"userlist"`
}

// 获取部门列表返回结果
type ResDepartmentList struct {
	util.CommonError
	Department []Department `json:"department"`
}

type User struct {
	Userid       string `json:"userid"`
	Name         string `json:"name"`
	Departmentid []int  `json:"department"`
	Mobile       string `json:"mobile"`
}

type Department struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Parentid int    `json:"parentid"`
}

// getMemberDetail
func (MemberManage *MemberManage) GetMemberDetail(accessToken string, departmentID int) (result ResMemberDetail, err error) {

	urlStr := fmt.Sprintf(getMemberDetailURL, accessToken, strconv.Itoa(departmentID), "1")

	var response []byte
	response, err = util.HTTPGet(urlStr)
	if err != nil {
		return
	}
	err = json.Unmarshal(response, &result)
	if err != nil {
		return
	}
	if result.ErrCode != 0 {
		err = fmt.Errorf("getMemberDetail error : errcode=%v , errmsg=%v", result.ErrCode, result.ErrMsg)
		return
	}
	return
}

// get department list 若departmentID参数为-1则获取企业所有部门
func (MemberManage *MemberManage) GetDepartmentList(accessToken string, departmentID int) (result ResDepartmentList, err error) {
	urlStr := ""
	if departmentID == -1 {
		urlStr = fmt.Sprintf(getDepartmentNameURL, accessToken)
	} else {
		urlStr = fmt.Sprintf(getDepartmentNameURL, accessToken, strconv.Itoa(departmentID))
	}
	var response []byte
	response, err = util.HTTPGet(urlStr)
	if err != nil {
		return
	}
	err = json.Unmarshal(response, &result)
	if err != nil {
		return
	}

	if result.ErrCode != 0 {
		err = fmt.Errorf("GetDepartmentList error : errcode=%v , errmsg=%v", result.ErrCode, result.ErrMsg)
		return
	}

	if len(result.Department) == 0 {
		return result, httpapi.NewErr(constant.GLOBAL_SYS_ERR, "the departments list from wechat server is empty!", err)
	}
	return
}
