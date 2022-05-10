package controller

import "web_app/models"

//专门用来放接口文档的model

type _ResponsePostList struct {
	Code    string                  `json:"code"`    //业务响应状态码
	Message string                  `json:"message"` //提示信息
	Data    []*models.ApiPostDetail `json:"data"`    //数据
}
