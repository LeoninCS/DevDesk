package Server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	httpServer *http.Server
	Port       string
}

func NewServer(addr string, handler http.Handler) (*Server, error) {
	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}
	return &Server{
		httpServer: srv,
		Port:       addr,
	}, nil
}

func (s *Server) Start() error {
	// 在goroutine中启动服务器以便能监听中断信号
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务器启动失败: %v", err)
		}
	}()

	log.Printf("服务器正在监听地址 %s", s.httpServer.Addr)

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("正在关闭服务器...")

	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Fatal("服务器关闭错误:", err)
	}

	log.Println("服务器已关闭")
	return nil
}
