package main

// @title URL Shortener API
// @version 1.0
// @description This is a sample API for a URL shortener with Swagger documentation.
// @host localhost:8080
// @BasePath /

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"url-shortener/config"

	_ "url-shortener/docs"
	"url-shortener/pkg/logging"
	"url-shortener/pkg/storage/postgres"

	"url-shortener/internal/controller"
	"url-shortener/internal/repository"
	"url-shortener/internal/service"
)

const (
	logFile            = "logs/server.log"
	serverStartTimeout = 10 * time.Second
)

func main() {
	storageFlag := flag.String("storage", "", "flag for specifying storage")
	flag.Parse()
	fmt.Println(*storageFlag)

	if *storageFlag != "cache" && *storageFlag != "postgres" {
		panic("the argument: %s, is not supported, specify argument cache or postgres")
	}

	// 	init logger
	logging.InitLogger(logFile)
	logger, err := logging.GetLogger(logFile)
	if err != nil {
		panic(err)
	}

	// 	init config
	projectRoot, err := os.Getwd()
	if err != nil {
		logger.Fatalf("Error getting working directory: %v", err)
	}
	envFilePath := filepath.Join(projectRoot, ".env")
	cfg := config.GetConfig(logFile, envFilePath)

	// 	init storage

	var storage service.Storage

	if *storageFlag == "cache" {
		storage = repository.NewCacheStorage()
	} else {
		pool, err := postgres.NewClient(context.Background(), cfg.DataBase)
		if err != nil {
			panic(err)
		}
		storage = repository.NewDataBaseStorage(&pool)
	}
	// 		init postgres
	// pool, err := postgres.NewClient(context.Background(), cfg.DataBase)
	// pstgrs := repository.NewDataBaseStorage(&pool)

	// // 		init cache
	// cache := repository.NewCacheStorage()
	//
	// storage := repository.NewStorage(pstgrs, cache)

	service := service.NewShortenerService(storage)
	// 	init router
	router := gin.Default()

	handler := handler.NewHandler(service, logger)
	handler.Register(router)
	start(router, storage, logger, cfg)

}

func start(router *gin.Engine, storage service.Storage, logger *logging.Logger, cfg *config.Config) {
	logger.Info("start application")
	var listener net.Listener
	var listenErr error

	if cfg.Listen.Type == "socket" {
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logger.Fatal(err)
		}
		logger.Info("create socket")
		socketPath := path.Join(appDir, "app.sock")
		logger.Debugf("socket path: %s", socketPath)

		logger.Info("listen unix socket")
		listener, listenErr = net.Listen("unix", socketPath)
		logger.Infof("server is listening on unix socket: %s", socketPath)

	} else {
		logger.Info("listen port")
		listener, listenErr = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port))
		logger.Infof("server is listening on %s:%s", cfg.Listen.BindIP, cfg.Listen.Port)
	}
	if listenErr != nil {
		logger.Fatal(listenErr)
	}

	func(ctx context.Context) {
		ctx, cancel := context.WithCancel(ctx)
		go func() {
			defer cancel()
			logger.Fatal(router.RunListener(listener))
		}()

		// go func() {
		// 	defer cancel()
		//
		// 	storage.CashChecker(cfg.CacheTTL)
		// }()
		notifyCtx, notify := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
		defer notify()

		go func() {
			defer cancel()
			<-notifyCtx.Done()
			closer := make(chan struct{})

			go func() {
				closer <- struct{}{}
			}()

			shutdownCtx, shutdown := context.WithTimeout(context.Background(), serverStartTimeout)
			defer shutdown()
			runtime.Gosched()

			select {
			case <-closer:
				logger.Info("shutting down gracefully")
			case <-shutdownCtx.Done():
				logger.Info("shutting down forcefully")
			}
		}()

		<-ctx.Done()
		cancel()

	}(context.Background())

}
