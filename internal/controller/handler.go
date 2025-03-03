package handler

import (
	"net/http"

	"url-shortener/internal/model"

	"github.com/gin-gonic/gin"
	_ "github.com/swaggo/files"
	swaggerFiles "github.com/swaggo/files"
	_ "github.com/swaggo/gin-swagger"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "url-shortener/docs"
)

const (
	extendUrl  = "/expand"
	shortenUrl = "/shorten"
)

// @Description Формат ответа об ошибке
type ErrorResponse struct {
	Message string `json:"message"`
}

type shortenerService interface {
	Shortening(string) (string, error)
	Expansion(string) (string, error)
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
// @Param shortUrl body model.ShortURL true "Короткая ссылка"
// @Success 200 {object} map[string]string "Расширенная длинная ссылка"
// @Failure 400 {object} ErrorResponse "Неверный ввод"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /expand [get]
func (h *Handler) Expansion(ctx *gin.Context) {
	var shortUrl model.ShortURL
	if err := ctx.ShouldBindJSON(&shortUrl); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		ctx.Abort()
		return
	}
	res, err := h.shortenerService.Expansion(shortUrl.URL)
	if err != nil {
		h.logger.Errorf("Ошибка при расширении: %v", err)
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, map[string]string{"long_url": res})
}

// @Summary Сократить длинную ссылку
// @Description Преобразует длинную ссылку в компактную форму.
// @Tags Сокращение URL
// @Accept json
// @Produce json
// @Param longUrl body model.LongURL true "Длинная ссылка"
// @Success 200 {object} map[string]string "Сокращённая ссылка"
// @Failure 400 {object} ErrorResponse "Неверный ввод"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /shorten [post]
func (h *Handler) Shortening(ctx *gin.Context) {
	var longUrl model.LongURL
	if err := ctx.ShouldBindJSON(&longUrl); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		ctx.Abort()
		return
	}

	res, err := h.shortenerService.Shortening(longUrl.URL)
	if err != nil {
		h.logger.Errorf("Ошибка при сокращении: %v", err)
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, map[string]string{"short_url": res})
}
