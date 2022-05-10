// @Title controller
// @Description 帖子相关
// @Author doudoudoubi1220 2022-04-07
// @Update doudoudoubi1220 2022-05-07
package controller

import (
	"strconv"
	"web_app/logic"
	"web_app/models"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// CreatePostHandler 创建帖子
// @Summary 创建帖子接口
// @Description 根据用户输入的数据创建一个帖子
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer JWT"
// @Param object body models.Post false "帖子参数"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostList
// @Router /post [post]
func CreatePostHandler(c *gin.Context) {

	//1.获取参数及参数的校验
	p := new(models.Post)
	if err := c.ShouldBindJSON(p); err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	//从 c 取到当前发请求的用户ID
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.AuthorID = userID
	//2.创建帖子
	if err := logic.CreatePost(p); err != nil {
		zap.L().Error("logic.CreatePost(p) failed", zap.Error(err))
		return
	}
	//3.返回响应
	ResponseSuccess(c, CodeSuccess)

}

// GetPostHandler 获取帖子详情的接口
// @Summary 获取帖子详情的接口
// @Description 根据传入的postid查询帖子的详细信息
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer JWT"
// @Param id path string true "帖子ID"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostList
// @Router /post/{id} [get]
func GetPostHandler(c *gin.Context) {
	//1.获取参数(从URL中获取帖子的id)
	pidStr := c.Param("id")
	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		zap.L().Error("GetPostDetail with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	//2.根据id取出帖子数据
	data, err := logic.GetPostByID(pid)
	if err != nil {
		zap.L().Error("logic.GetPostByID(pid) failed,", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c, data)
}

// GetPostListHandler 获取帖子列表的接口
// @Summary 获取帖子列表的接口
// @Description 获取所有帖子列表，根据传递的参数进行分页，按照发部顺序进行排序
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer JWT"
// @Param page query int false "页码"
// @Param size query int false "每页数量"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostList
// @Router /posts [get]
func GetPostListHandler(c *gin.Context) {
	//获取分页参数
	page, size := GetPageInfo(c)

	//获取数据
	data, err := logic.GetPostList(page, size)
	if err != nil {
		zap.L().Error("logic.GetPostList() failed,", zap.Error(err))
		return
	}
	//返回响应
	ResponseSuccess(c, data)
}

// GetPostListHandler2 升级版帖子列表接口
// @Summary 升级版帖子列表接口
// @Description 可按社区按时间或分数排序查询帖子列表接口
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object body models.ParamPostList false "查询参数"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostList
// @Router /post2 [get]
func GetPostListHandler2(c *gin.Context) {

	//GET请求参数(query string)： /api/v1/post2?page=1&size=10&order=time
	//初始化结构体时指定初始参数
	p := &models.ParamPostList{
		Page:  1,
		Size:  10,
		Order: models.OrderTime,
	}
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetPostListHandler2 with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	//获取数据
	data, err := logic.GetPostListNew(p)
	if err != nil {
		zap.L().Error("logic.GetPostList() failed,", zap.Error(err))
		return
	}
	//返回响应
	ResponseSuccess(c, data)
}

//func GetCommunityPostListHandler(c *gin.Context) {
//
//	//GET请求参数(query string)： /api/v1/post2?page=1&size=10&order=time
//	//初始化结构体时指定初始参数
//	p := &models.ParamCommunityPostList{
//		ParamPostList: models.ParamPostList{
//			Page:  1,
//			Size:  10,
//			Order: models.OrderTime,
//		},
//	}
//	if err := c.ShouldBindQuery(p); err != nil {
//		zap.L().Error("GetCommunityPostListHandler with invalid params", zap.Error(err))
//		ResponseError(c, CodeInvalidParam)
//		return
//	}
//
//	//获取数据
//	data, err := logic.GetCommunityPostListHandler(p)
//	if err != nil {
//		zap.L().Error("logic.GetPostList() failed,", zap.Error(err))
//		return
//	}
//	//返回响应
//	ResponseSuccess(c, data)
//}
