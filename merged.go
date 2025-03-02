//// Файл: ./cmd/main.go
////
package main

// @title URL Shortener API
// @version 1.0
// @description This is a sample API for a URL shortener with Swagger documentation.
// @host localhost:8080
// @BasePath /

import (
	"context"
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

	// 		init postgres
	pool, err := postgres.NewClient(context.Background(), cfg.DataBase)
	pstgrs := repository.NewDataBaseStorage(&pool)

	// 		init cache
	cache := repository.NewCacheStorage()

	storage := repository.NewStorage(pstgrs, cache)

	service := service.NewShortenerService(storage)

	// 	init router
	router := gin.Default()

	handler := handler.NewHandler(service, logger)
	handler.Register(router)
	start(router, cache, logger, cfg)

}

func start(router *gin.Engine, cache *repository.CacheStorage, logger *logging.Logger, cfg *config.Config) {
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

		go func() {
			defer cancel()
			cache.CashChecker(cfg.CacheTTL)
		}()
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

//// Файл: ./config/config.go
////
package config

import (
	"sync"
	"url-shortener/pkg/logging"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	Listen   Listen   `env:"LISTEN"`
	DataBase DataBase `env:"DATABASE"`
	CacheTTL int64    `env:"CACHE_TTL"`
}

type Listen struct {
	Type   string `env:"LISTEN_TYPE" env-default:"port"`
	BindIP string `env:"BIND_IP" env-default:"127.0.0.1"`
	Port   string `env:"PORT" env-default:"8080"`
}

type DataBase struct {
	Host     string `env:"DB_HOST" env-default:"postgres"`
	Port     string `env:"DB_PORT" env-default:"5432"`
	Username string `env:"DB_USERNAME" env-default:"postgres"`
	Password string `env:"DB_PASSWORD" env-default:"postgres"`
	DBName   string `env:"DB_NAME" env-default:"postgres"`
}

var instance *Config
var once sync.Once

func GetConfig(logFile, envFilePath string) *Config {
	once.Do(func() {
		logger, err := logging.GetLogger(logFile)
		if err != nil {
			panic(err)
		}
		logger.Info("read application configuration")
		instance = &Config{}
		if err := godotenv.Load(envFilePath); err != nil {
			logger.Fatal("Error loading .env file")
		}
		if err := env.Parse(instance); err != nil {
			logger.Fatalf("Failed to parse env variables: %v", err)
		}

	})
	return instance
}

//// Файл: ./internal/controller/handler.go
////
package handler

import (
	"context"
	"net/http"

	"url-shortener/internal/model"

	"github.com/gin-gonic/gin"

	_ "github.com/swaggo/gin-swagger"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/swaggo/files"
	swaggerFiles "github.com/swaggo/files"
	_ "url-shortener/docs"
)

const (
	extendUrl  = "/"
	shortenUrl = "/"
)

// ErrorResponse представляет структуру ответа об ошибке.
// @Description Формат ответа об ошибке
type ErrorResponse struct {
	Message string `json:"message"`
}

type shortenerService interface {
	Shortening(context.Context, string) (string, error)
	Expansion(context.Context, string) (string, error)
}

type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})

	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

type Handler struct {
	shortenerService
	logger Logger
}

func NewHandler(shortenerService shortenerService, logger Logger) *Handler {
	return &Handler{shortenerService: shortenerService, logger: logger}
}

func (h *Handler) Register(router *gin.Engine) {
	router.GET(extendUrl, h.Expansion)
	router.POST(shortenUrl, h.Shortening)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) // Добавляем Swagger UI
}

// @Summary Расширить короткую ссылку до её оригинальной формы
// @Description Преобразует короткую ссылку в исходную длинную ссылку.
// @Tags Расширение URL
// @Accept json
// @Produce json
// @Param shortUrl path string true "Короткая ссылка"
// @Success 200 {object} map[string]string "Расширенная длинная ссылка"
// @Failure 400 {object} ErrorResponse "Неверный ввод"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /:shortUrl [post]
func (h *Handler) Expansion(ctx *gin.Context) {
	var shortUrl model.ShortURL
	if err := ctx.ShouldBindJSON(&shortUrl); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		ctx.Abort()
		return
	}
	res, err := h.shortenerService.Expansion(ctx.Request.Context(), shortUrl.URL)
	if err != nil {
		h.logger.Errorf("Ошибка при расширении: %v", err)
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, map[string]string{"long url": res})
}

// @Summary Сократить длинную ссылку
// @Description Преобразует длинную ссылку в компактную форму.
// @Tags Сокращение URL
// @Accept json
// @Produce json
// @Param longUrl path string true "Длинная ссылка"
// @Success 200 {object} map[string]string "Сокращённая ссылка"
// @Failure 400 {object} ErrorResponse "Неверный ввод"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /:longUrl [get]
func (h *Handler) Shortening(ctx *gin.Context) {
	var longUrl model.LongURL
	if err := ctx.ShouldBindJSON(&longUrl); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		ctx.Abort()
		return
	}

	res, err := h.shortenerService.Shortening(ctx.Request.Context(), longUrl.URL)
	if err != nil {
		h.logger.Errorf("Ошибка при сокращении: %v", err)
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, map[string]string{"short url": res})
}

//// Файл: ./internal/repository/cache.go
////
package repository

import (
	"sync"
	"time"

	"url-shortener/pkg/storage"
)

type Node struct {
	data        string
	RequestTime time.Time
}

type CacheStorage struct {
	data map[string]Node
	sync.Mutex
}

func NewCacheStorage() *CacheStorage {
	return &CacheStorage{data: make(map[string]Node)}
}

func (c *CacheStorage) GetLongUrl(shortURL string) (string, error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	res, ok := c.data[shortURL]
	if !ok {
		return "", storage.ErrNotFound
	}
	c.data[shortURL] = Node{
		data:        res.data,
		RequestTime: time.Now(),
	}
	return res.data, nil
}

func (s *CacheStorage) Insert(shortURL, longURL string) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	if _, ok := s.data[shortURL]; ok {
		return storage.ErrAlreadyExists
	}
	s.data[shortURL] = Node{
		data:        longURL,
		RequestTime: time.Now(),
	}
	return nil
}

func (s *CacheStorage) CashChecker(lifeTime int64) {
	for {
		time.Sleep(time.Duration(lifeTime) * time.Millisecond)
		s.Mutex.Lock()
		for key := range s.data {
			if int64(time.Now().Sub(s.data[key].RequestTime).Minutes()) > lifeTime {
				delete(s.data, key)
			}
		}
		s.Unlock()
	}
}

//// Файл: ./internal/repository/repository.go
////
package repository

import (
	"context"
	"url-shortener/pkg/storage"
)

type Storage struct {
	db    *DataBaseStorage
	cache *CacheStorage
}

func NewStorage(db *DataBaseStorage, cache *CacheStorage) *Storage {
	return &Storage{db: db, cache: cache}
}

func (s *Storage) GetLongUrl(ctx context.Context, shortUrl string) (string, error) {
	res, err := s.cache.GetLongUrl(shortUrl)
	if err != nil {
		if err == storage.ErrNotFound {
			res, err = s.db.GetLongUrl(ctx, shortUrl)
			if err == nil {
				s.cache.Insert(shortUrl, res)
			}
		} else {
			res, err = s.db.GetLongUrl(ctx, shortUrl)
		}
	}
	return res, err

}

func (s *Storage) Insert(ctx context.Context, shortURL, longURL string) error {
	err := s.db.Insert(ctx, shortURL, longURL)
	if err != nil {
		return err
	}
	s.cache.Insert(shortURL, longURL)
	return nil
}

//// Файл: ./internal/repository/postgres.go
////
package repository

import (
	"context"
	// "errors"
	// "strconv"

	"url-shortener/pkg/storage"
	"url-shortener/pkg/storage/postgres"
)

type DataBaseStorage struct {
	pool *postgres.Pool
}

func NewDataBaseStorage(pool *postgres.Pool) *DataBaseStorage {
	return &DataBaseStorage{pool: pool}
}

func (s *DataBaseStorage) Insert(ctx context.Context, shortURL, longURL string) error {
	query := "INSERT INTO urls (short_url, long_url) VALUES ($1, $2) ON CONFLICT (short_url) DO NOTHING"
	_, err := s.pool.Query(ctx, query, shortURL, longURL)
	if postgres.IsDuplicateError(err) {
		return storage.ErrAlreadyExists
	}
	return err
}

func (s *DataBaseStorage) GetLongUrl(ctx context.Context, shortURL string) (string, error) {
	var longURL string
	err := s.pool.QueryRow(ctx, "SELECT long_url FROM urls WHERE short_url = $1", shortURL).Scan(&longURL)
	if err == postgres.ErrNotFound {
		return "", storage.ErrNotFound
	}
	return longURL, err
}

//// Файл: ./internal/service/user_service.go
////
package service

import (
	"context"
	"crypto/sha256"
	"errors"
	"url-shortener/pkg/storage"
)

const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz_"
const hashLength = 8 // Длина желаемого хэша
const maxIndex = len(alphabet) ^ 2 - 1

//go:generate mockgen -source=user_service.go -destination=mocks/mock.go
type Storage interface {
	GetLongUrl(ctx context.Context, shortUrl string) (string, error)
	Insert(ctx context.Context, shortUrl, longUrl string) error
}

type ShortenerService struct {
	Storage Storage
}

func NewShortenerService(Storage Storage) *ShortenerService {
	return &ShortenerService{Storage: Storage}
}

func (s ShortenerService) Shortening(ctx context.Context, longUrl string) (string, error) {
	id := 0
	hash := encodeHash(longUrl)
	for id < maxIndex {
		shortUrl := hash + IntToIndex63(id)
		longCheck, err := s.Storage.GetLongUrl(ctx, shortUrl)
		if err == storage.ErrNotFound {
			err = s.Storage.Insert(ctx, shortUrl, longUrl)
			return shortUrl, err
		} else if longCheck == longUrl {
			return shortUrl, err
		}
		id++
	}
	return "", errors.New("index out of range")
}

func (s ShortenerService) Expansion(ctx context.Context, shortUrl string) (string, error) {
	res, err := s.Storage.GetLongUrl(ctx, shortUrl)
	return res, err
}

// Функция для преобразования байтов в строку фиксированной длины
func encodeHash(input string) string {
	// Создаем хэш
	hasher := sha256.New()
	hasher.Write([]byte(input))
	hash := hasher.Sum(nil)

	result := ""
	for i := 0; i < hashLength; i++ {
		index := int(hash[i]) % len(alphabet) // Получаем индекс из хэша
		result += string(alphabet[index])     // Добавляем символ из алфавита
	}
	return result
}

// Функция для преобразования числа из 10-тичной системы счисления в 63-ную
func IntToIndex63(id int) string {
	res := ""
	res += string(alphabet[id/len(alphabet)])
	res += string(alphabet[id%len(alphabet)])
	return res

}

//// Файл: ./internal/model/model.go
////
package model

type LongURL struct {
	URL string `json:"long_url"`
}

type ShortURL struct {
	URL string `json:"short_url"`
}

//// Файл: ./pkg/logging/logging.go
////
package logging

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
)

type writerHook struct {
	Writer    []io.Writer
	LogLevels []logrus.Level
}

func (hook *writerHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}
	for _, w := range hook.Writer {
		w.Write([]byte(line))
	}
	return err
}

func (hook *writerHook) Levels() []logrus.Level {
	return hook.LogLevels
}

var e *logrus.Entry

var entryes = make(map[string]*logrus.Entry)

type Logger struct {
	*logrus.Entry
}

func GetLogger(file string) (*Logger, error) {
	e, ok := entryes[file]
	if !ok {
		return nil, errors.New(fmt.Sprintf("log file %s not found", file))
	}
	return &Logger{e}, nil
}

func (l *Logger) GetLoggerWithField(k string, v interface{}) *Logger {
	return &Logger{l.WithField(k, v)}
}

func init() {
	err := os.MkdirAll("logs", 0766)
	if err != nil {
		panic(err)
	}
}

func InitLogger(logFile string) {
	l := logrus.New()
	l.SetReportCaller(true)
	l.Formatter = &logrus.JSONFormatter{
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			filename := path.Base(frame.File)
			return fmt.Sprintf("%s()", frame.Function), fmt.Sprintf("%s:%d", filename, frame.Line)
		},
	}

	allFile, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		panic(err)
	}

	l.SetOutput(io.Discard)

	l.AddHook(&writerHook{
		Writer:    []io.Writer{allFile, os.Stdout},
		LogLevels: logrus.AllLevels,
	})

	l.SetLevel(logrus.TraceLevel)
	entryes[logFile] = logrus.NewEntry(l)
}

//// Файл: ./pkg/storage/inMemory/cache.go
////
package inmemory

import (
	"sync"
	"url-shortener/pkg/storage"
)

type MemoryStorage struct {
	data map[string]string
	sync.Mutex
}

func NewStorage(hasher func(string) string) *MemoryStorage {
	return &MemoryStorage{data: make(map[string]string)}
}

func (s *MemoryStorage) Insert(key, value string) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	_, ok := s.data[key]
	if !ok {
		return storage.ErrAlreadyExists
	}
	s.data[key] = value
	return nil
}

func (s *MemoryStorage) Get(key string) (string, error) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	res, ok := s.data[key]
	if !ok {
		return "", storage.ErrNotFound
	}
	return res, nil
}

//// Файл: ./pkg/storage/postgres/postgres.go
////
package postgres

import (
	"context"
	"fmt"
	"time"

	"errors"

	"url-shortener/config"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool struct {
	*pgxpool.Pool
}

var ErrNotFound = pgx.ErrNoRows

func NewClient(ctx context.Context, cfg config.DataBase) (Pool, error) {
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	p, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return Pool{}, err
	}
	return Pool{p}, nil
}

func IsDuplicateError(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return true
	}
	return false
}

//// Файл: ./pkg/storage/storage.go
////
package storage

import "errors"

var (
	ErrNotFound      = errors.New("url not found")
	ErrAlreadyExists = errors.New("url already exists")
)

//// Файл: ./docs/docs.go
////
// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/:longUrl": {
            "get": {
                "description": "Преобразует длинную ссылку в компактную форму.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Сокращение URL"
                ],
                "summary": "Сократить длинную ссылку",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Длинная ссылка",
                        "name": "longUrl",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Сокращённая ссылка",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "Неверный ввод",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/:shortUrl": {
            "post": {
                "description": "Преобразует короткую ссылку в исходную длинную ссылку.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Расширение URL"
                ],
                "summary": "Расширить короткую ссылку до её оригинальной формы",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Короткая ссылка",
                        "name": "shortUrl",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Расширенная длинная ссылка",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "Неверный ввод",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "handler.ErrorResponse": {
            "description": "Формат ответа об ошибке",
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "URL Shortener API",
	Description:      "This is a sample API for a URL shortener with Swagger documentation.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}

//// Файл: ./merged.go
////

