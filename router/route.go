package router

import (
	"net/http"
	"time"
	"web_app/controller"
	_ "web_app/docs" // 千万不要忘了导入把你上一步生成的docs
	"web_app/logger"
	"web_app/middlewares"

	"github.com/gin-gonic/gin"
	gs "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"github.com/gin-contrib/pprof"
)

func Setup() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true), middlewares.RateLimitMiddleware(2*time.Second, 1))

	r.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")

	//注册业务路由
	v1.POST("/signup", controller.SignUpHandler)
	v1.POST("/login", controller.LoginHandler)
	v1.Use(middlewares.JWTAuthMiddleware()) //应用JWT认证中间件
	{
		v1.GET("/community", controller.CommunityHandler)
		v1.GET("/community/:id", controller.CommunityDetailHandler)

		v1.POST("/post", controller.CreatePostHandler)
		v1.GET("/post/:id", controller.GetPostHandler)
		v1.GET("/posts", controller.GetPostListHandler)

		v1.POST("/vote", controller.PostVoteController)

		//投票
		v1.GET("/post2", controller.GetPostListHandler2)
	}

	//r.GET("/ping", middlewares.JWTAuthMiddleware(), func(c *gin.Context) {
	//	////如果是登录的用户,判断请求头中是否有 有效的JWT ？
	//	//isLogin := true
	//	//c.Request.Header.Get("Authorization")
	//	//if isLogin {
	//	//	c.String(http.StatusOK, "pong")
	//	//} else {
	//	//	//否则就直接返回登陆
	//	//	c.String(http.StatusOK, "请登录")
	//	//}
	//	c.String(http.StatusOK, "pong")
	//})
	pprof.Register(r) //注册pprof相关路由
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})
	return r

}
