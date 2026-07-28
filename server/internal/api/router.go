package api

import (
	"strings"

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
	client := crawler.NewClient(cfg.JMComic.BaseURL, cfg.JMComic.Auth)

	r := &Router{
		client: client,
		db:     db,
		config: cfg,
	}

	engine := gin.Default()

	engine.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		auth := c.GetHeader("Authorization")
		if auth != "" && auth != "Bearer guest" && auth != "Bearer " && strings.HasPrefix(auth, "Bearer ") && len(auth) > 20 {
			r.client.SetAuth(auth)
		}
		c.Next()
	})

	r.engine = engine
	r.setupRoutes()
	return r
}

func (r *Router) setupRoutes() {
	api := r.engine.Group("/api")
	{
		// 浏览
		api.GET("/search", r.Search)
		api.GET("/index", r.GetIndex)
		api.GET("/latest", r.GetLatest)
		api.GET("/ranking", r.GetRanking)
		api.GET("/comic/:id", r.GetComic)
		api.GET("/comic/:id/chapters", r.GetChapters)
		api.GET("/chapter/:id", r.GetChapterImages)
		api.GET("/categories", r.GetCategories)
		api.GET("/categories/filter", r.GetCategoryFilter)

		// 评论
		api.GET("/comic/:id/comments", r.GetComments)
		api.GET("/comment/:id/sub", r.GetSubComments)

		// 收藏
		api.GET("/favorites", r.GetFavorites)
		api.POST("/favorites", r.AddFavorite)
		api.DELETE("/favorites/:id", r.RemoveFavorite)

		// 历史
		api.GET("/history", r.GetHistory)
		api.DELETE("/history/:id", r.DeleteHistory)
		api.DELETE("/history", r.ClearHistory)

		// 下载
		api.POST("/download", r.CreateDownload)
		api.GET("/downloads", r.GetDownloads)
		api.DELETE("/download/:id", r.DeleteDownload)
		api.DELETE("/downloads", r.ClearDownloads)

		// 用户
		api.POST("/login", r.Login)
		api.POST("/register", r.Register)
		api.GET("/user/info", r.GetUserInfo)
		api.POST("/user/sign", r.Sign)

		// 帮助
		api.GET("/help", r.GetHelp)
	}
}

func (r *Router) Run(addr string) error {
	return r.engine.Run(addr)
}
