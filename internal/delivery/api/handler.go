package api

import (
	_ "github.com/bllooop/musiclibrary/docs"
	"github.com/bllooop/musiclibrary/internal/usecase"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	usecases *usecase.Usecase
}

func NewHandler(services *usecase.Usecase) *Handler {
	return &Handler{usecases: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
