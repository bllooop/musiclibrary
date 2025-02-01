package api

import (
	"github.com/bllooop/musiclibrary/internal/usecase"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	usecases *usecase.Usecase
}

func NewHandler(services *usecase.Usecase) *Handler {
	return &Handler{usecases: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	api := router.Group("/api")
	{
		songs := api.Group("/songs")
		{
			songs.POST("/", h.createSong)
			songs.GET("/", h.getSongs)
			songs.GET("/get-text", h.getSongById)
			songs.DELETE("/:id", h.deleteSong)
			songs.PUT("/:id", h.updateSong)
		}
	}
	return router
}
