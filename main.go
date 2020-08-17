package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"webapp/dao/mysql"
	"webapp/dao/redis"
	"webapp/logger"
	"webapp/routers"
	"webapp/settings"

	"go.uber.org/zap"
)

// Go Web开发较通用的脚手架模板

func init() {
	// 1. 加载配置--视频中使用的是viper但是我用的是goini这个库
	settings.Setup()

	// 2. 初始化日志
	if err := logger.Init(); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
	}
	fmt.Println(zap.L())
	zap.L().Debug("logger init success")
	// 3. 初始化MySQL连接
	if err := mysql.InitDB(); err != nil {
		fmt.Printf("init mysql failed, err:%v\n", err)
	}
	// 4. 初始化Redis连接
	if err := redis.Init(); err != nil {
		fmt.Printf("init redist failed, err:%v\n", err)
	}

}
func main() {

	// 5. 注册路由
	r := routers.SetUp()

	defer mysql.Close()
	defer redis.Close()
	defer zap.L().Sync()
	// 启动服务（优雅关机）
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", settings.AppSetting.Port),
		Handler: r,
	}

	go func() {
		// 开启一个goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("listen error", zap.Error(err))
		}
	}()

	// 等待中断信号来优雅地关闭服务
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // 阻塞在此， 当接收到上述两种信号时才会往下进行
	zap.L().Info("ShutDown Server ...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（处理未处理完的请求再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server shutdown:", zap.Error(err))
	}

	zap.L().Info("server exiting")
}
