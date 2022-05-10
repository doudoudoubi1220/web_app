package controller

import (
	"strconv"
	"web_app/logic"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CommunityHandler(c *gin.Context) {
	//查询到所有的社区(community_id,community_name)以列表的形式返回
	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.GetCommunity() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy) //不轻易把服务端报错暴露给外面
		return
	}
	ResponseSuccess(c, data)
}

// CommunityDetailHandler 根据ID查询社区分类详情
func CommunityDetailHandler(c *gin.Context) {
	//获取社区id
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	data, err := logic.GetCommunityDetail(id)
	if err != nil {
		zap.L().Error("logic.GetCommunity() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy) //不轻易把服务端报错暴露给外面
		return
	}
	ResponseSuccess(c, data)
}
