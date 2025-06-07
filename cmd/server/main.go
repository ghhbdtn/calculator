package main

import (
	_ "calculator/docs"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"calculator/internal/app/service"
	"calculator/internal/app/strategy"
	grpcserver "calculator/internal/transport/grpc"
	httpserver "calculator/internal/transport/http"
)

func main() {
	evalStrategy := &strategy.EvaluationStrategy{}
	calculator := service.NewCalculator(evalStrategy)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Канал для ошибок серверов
	errChan := make(chan error, 2)

	// WaitGroup для ожидания завершения серверов
	var wg sync.WaitGroup

	// Запускаем HTTP сервер
	wg.Add(1)
	go func() {
		defer wg.Done()
		startHTTPServer(ctx, calculator, errChan)
	}()

	// Запускаем gRPC сервер
	wg.Add(1)
	go func() {
		defer wg.Done()
		startGRPCServer(ctx, calculator, errChan)
	}()

	// Обработка сигналов завершения
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		log.Println("Received shutdown signal")
		cancel()
	}()

	// Мониторинг ошибок
	go func() {
		for err := range errChan {
			if err != nil {
				log.Printf("Server error: %v", err)
				cancel()
				return
			}
		}
	}()

	// Ожидаем завершения работы серверов
	wg.Wait()
	log.Println("Shutdown completed")
}
func startHTTPServer(ctx context.Context, calculator *service.Calculator, errChan chan<- error) {
	// Создаем экземпляр нашего сервера
	httpServer := httpserver.NewServer(calculator)

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := httpServer.Stop(shutdownCtx); err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
		}
	}()

	// Запускаем сервер
	if err := httpServer.Start(":8080"); err != nil && err != http.ErrServerClosed {
		errChan <- err
	}
}

func startGRPCServer(ctx context.Context, calculator *service.Calculator, errChan chan<- error) {
	server := grpcserver.NewServer(calculator)

	go func() {
		<-ctx.Done()
		server.Stop()
	}()

	if err := server.Start(":50051"); err != nil {
		errChan <- err
	}
}
