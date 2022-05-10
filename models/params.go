package models

const (
	OrderTime  = "time"
	OrderScore = "score"
)

// 定义请求的参数结构体

// ParamSignUp 注册请求参数
type ParamSignUp struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"rePassword" binding:"required,eqfield=Password"`
}

// ParamLogin 登录请求参数
type ParamLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ParamVoteData struct {

	//UserID 从请求中获取当前的账户
	PostID    string `json:"post_id" binding:"required"`              //帖子id
	Direction int8   `json:"direction,string" binding:"oneof=1 0 -1"` //赞成票(1)还是反对票(-1)取消投票(0)
}

//获取帖子列表query string 参数
type ParamPostList struct {
	CommunityID int64  `json:"community_id" form:"community_id"`   // 社区ID　可以为空
	Page        int64  `json:"page" form:"page" example:"1"`       //页码
	Size        int64  `json:"size" form:"size" example:"10"`      //每页数据量
	Order       string `json:"order" form:"order" example:"score"` //排序依据
}

type ParamCommunityPostList struct {
	//Page        int64  `json:"page" form:"page"`
	//Size        int64  `json:"size" form:"size"`
	//Order       string `json:"order" form:"order"`
	ParamPostList
	CommunityID int64 `json:"community_id" form:"community_id"`
}
