package main

import (
	Config "Server/Config"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {

	// ServerConfig 초기화
	var config Config.ServerConfig
	if err := config.ServerConfigLoad(); err != nil {
		return
	}

	// 1) 라우터 생성
	r := gin.New()
	// 처음엔 최소만: 필요 시 나중에 미들웨어 추가
	r.Use(gin.Recovery()) // 패닉 리커버리
	r.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	api := r.Group("/api/v1") // 버전 고정
	api.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "hello",
			"ts":      time.Now().UTC().Format(time.RFC3339),
		})
	})

	// 2) 서버 설정(타임아웃 포함)
	srv := &http.Server{
		Addr:              ":8080", // 환경변수로 바꾸고 싶으면 os.Getenv("PORT") 사용
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// 3) 서버 비동기 시작
	go func() {
		log.Printf("listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// 4) Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown error: %v", err)
	}
	log.Println("server shut down")
}
