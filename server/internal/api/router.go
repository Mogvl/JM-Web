package api

import (
	"github.com/Mogvl/JM-Web/server/config"
	"github.com/Mogvl/JM-Web/server/internal/crawler"
	"github.com/Mogvl/JM-Web/server/internal/database"
	"github.com/gin-gonic/gin"
)

type Router struct {
	engine *gin.Engine
	client *crawler.Client
	db     *database.DB
	config *config.Config
}

func NewRouter(cfg *config.Config, db *database.DB) *Router {
	engine := gin.Default()

	engine.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	client := crawler.NewClient(cfg.JMComic.BaseURL, cfg.JMComic.Auth)

	r := &Router{
		engine: engine,
		client: client,
		db:     db,
		config: cfg,
	}

	r.setupRoutes()
	return r
}

func (r *Router) setupRoutes() {
	api := r.engine.Group("/api")
	{
		api.GET("/search", r.Search)
		api.GET("/comic/:id", r.GetComic)
		api.GET("/comic/:id/chapters", r.GetChapters)
		api.GET("/chapter/:id", r.GetChapterImages)
		api.GET("/categories", r.GetCategories)

		api.GET("/favorites", r.GetFavorites)
		api.POST("/favorites", r.AddFavorite)
		api.DELETE("/favorites/:id", r.RemoveFavorite)

		api.GET("/history", r.GetHistory)
		api.DELETE("/history/:id", r.DeleteHistory)

		api.POST("/download", r.CreateDownload)
		api.GET("/downloads", r.GetDownloads)
	}
}

func (r *Router) Run(addr string) error {
	return r.engine.Run(addr)
}
