package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"web_app/controller"
	"web_app/dao/mysql"
	"web_app/dao/redis"
	"web_app/logger"
	"web_app/pkg/snowflake"
	"web_app/router"
	"web_app/settings"

	"github.com/spf13/viper"

	"go.uber.org/zap"
)

/* swagger main 函数注释格式（写项目相关介绍信息）
// @title 这里写标题
// @version 1.0
// @description 这里写描述信息
// @termsOfService http://swagger.io/terms/

// @contact.name 这里写联系人信息
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host 这里写接口服务的host
// @BasePath 这里写base path
*/

// @title bluebell
// @version 1.0
// @description
// @tremsOfService http://swagger.io/terms/

// @contact.name doudoudoubi1220

// @contact.email 1715925630@nyist.edu.cn

// @license.name Apache 2.0
// @license.url http://www.apache.org/license/LICENSE-2.0html

// @host localhost:8081
// @BasePath /api/v1

//go web 开发较通用的脚手架模板
func main() {
	//1.加载配置
	if err := settings.Init(); err != nil {
		fmt.Println("init settings failed,err:", err)
		return
	}
	//2.初始化日志
	if err := logger.Init(settings.Conf.LogConfig, settings.Conf.Mode); err != nil {
		fmt.Println("init logger failed,err:", err)
		return
	}
	defer zap.L().Sync()
	zap.L().Debug("logger init success")
	//3.初始化MySQL连接
	if err := mysql.Init(settings.Conf.MySQLConfig); err != nil {
		fmt.Println("init settings failed,err:", err)
		return
	}
	defer mysql.Close()
	//4.初始化redis连接
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		fmt.Println("init settings failed,err:", err)
		return
	}
	defer redis.Close()
	//初始化gin框架内置的校验器私用的翻译器
	if err := controller.InitTrans("zh"); err != nil {
		fmt.Println("init validator trans failed,err:", err)
		return
	}
	//初始化雪花算法
	if err := snowflake.Init("2022-04-01", 1); err != nil {
		zap.L().Error("init snowflake failed", zap.Error(err))
		return
	}
	//5.注册路由
	r := router.Setup()
	//6.优雅关机
	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", viper.GetInt("port")),

		Handler: r,
	}

	go func() {
		// 开启一个goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("listen: %s\n", zap.Error(err))
		}
	}()

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	zap.L().Info("Shutdown Server ...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown: ", zap.Error(err))
	}

	zap.L().Info("Server exiting")
}
