package api

import (
	"net/http"
	"strconv"

	"github.com/Mogvl/JM-Web/server/internal/model"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// ==================== 搜索和浏览 ====================

func (r *Router) Search(c *gin.Context) {
	query := c.Query("q")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "q is required"})
		return
	}

	result, err := r.client.Search(query, page)
	if err != nil {
		log.Errorf("Search failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Search failed"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (r *Router) GetComic(c *gin.Context) {
	id := c.Param("id")

	comic, err := r.client.GetComicDetail(id)
	if err != nil {
		log.Errorf("Get comic failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get comic"})
		return
	}

	r.db.SaveComic(&model.Comic{
		ID:          comic.ID,
		Title:       comic.Title,
		Author:      comic.Author,
		Description: comic.Description,
		CoverURL:    comic.CoverURL,
		Tags:        joinStrings(comic.Tags),
		Category:    comic.Category,
		Status:      comic.Status,
	})

	c.JSON(http.StatusOK, comic)
}

func (r *Router) GetChapters(c *gin.Context) {
	comicID := c.Param("id")

	chapters, err := r.client.GetChapters(comicID)
	if err != nil {
		log.Errorf("Get chapters failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chapters"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"chapters": chapters})
}

func (r *Router) GetChapterImages(c *gin.Context) {
	chapterID := c.Param("id")

	images, err := r.client.GetChapterImages(chapterID)
	if err != nil {
		log.Errorf("Get images failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get images"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"images": images})
}

func (r *Router) GetCategories(c *gin.Context) {
	categories, err := r.client.GetCategories()
	if err != nil {
		log.Errorf("Get categories failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get categories"})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// ==================== 收藏 ====================

func (r *Router) GetFavorites(c *gin.Context) {
	favorites, err := r.db.GetFavorites()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get favorites"})
		return
	}
	c.JSON(http.StatusOK, favorites)
}

func (r *Router) AddFavorite(c *gin.Context) {
	var req struct {
		ComicID string `json:"comic_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "comic_id required"})
		return
	}

	if err := r.db.AddFavorite(req.ComicID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add favorite"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Added"})
}

func (r *Router) RemoveFavorite(c *gin.Context) {
	comicID := c.Param("id")

	if err := r.db.RemoveFavorite(comicID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove favorite"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Removed"})
}

// ==================== 历史 ====================

func (r *Router) GetHistory(c *gin.Context) {
	history, err := r.db.GetHistory()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get history"})
		return
	}
	c.JSON(http.StatusOK, history)
}

func (r *Router) DeleteHistory(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := r.db.DeleteHistory(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

func (r *Router) ClearHistory(c *gin.Context) {
	if err := r.db.ClearHistory(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear history"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Cleared"})
}

// ==================== 下载 ====================

func (r *Router) CreateDownload(c *gin.Context) {
	var req struct {
		ComicID string `json:"comic_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "comic_id required"})
		return
	}

	dl, err := r.db.CreateDownload(req.ComicID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create download"})
		return
	}

	go r.processDownload(dl.ID, req.ComicID)

	c.JSON(http.StatusOK, dl)
}

func (r *Router) GetDownloads(c *gin.Context) {
	downloads, err := r.db.GetDownloads()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get downloads"})
		return
	}
	c.JSON(http.StatusOK, downloads)
}

func (r *Router) DeleteDownload(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := r.db.DeleteDownload(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete download"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

func (r *Router) ClearDownloads(c *gin.Context) {
	if err := r.db.ClearDownloads(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear downloads"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Cleared"})
}

// ==================== 排行榜 ====================

func (r *Router) GetIndex(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))

	result, err := r.client.GetIndexInfo(page)
	if err != nil {
		log.Errorf("Get index failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (r *Router) GetLatest(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))

	result, err := r.client.GetLatest(page)
	if err != nil {
		log.Errorf("Get latest failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (r *Router) GetRanking(c *gin.Context) {
	rankType := c.DefaultQuery("type", "daily")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	result, err := r.client.GetRanking(rankType, page)
	if err != nil {
		log.Errorf("Get ranking failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get ranking"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ==================== 评论 ====================

func (r *Router) GetComments(c *gin.Context) {
	comicID := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	comments, err := r.client.GetComments(comicID, page)
	if err != nil {
		log.Errorf("Get comments failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get comments"})
		return
	}

	c.JSON(http.StatusOK, comments)
}

func (r *Router) GetSubComments(c *gin.Context) {
	commentID := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	comments, err := r.client.GetSubComments(commentID, page)
	if err != nil {
		log.Errorf("Get sub comments failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get sub comments"})
		return
	}

	c.JSON(http.StatusOK, comments)
}

// ==================== 用户 ====================

func (r *Router) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username and password required"})
		return
	}

	userData, err := r.client.Login(req.Username, req.Password)
	if err != nil {
		log.Errorf("Login failed: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Login failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":    userData.JWTToken,
		"username": userData.Username,
		"coins":    userData.Coin,
		"level":    userData.Level,
		"level_name": userData.LevelName,
		"avatar":    userData.Photo,
		"favorites": userData.AlbumFavorites,
	})
}

func (r *Router) Register(c *gin.Context) {
	var req struct {
		Username        string `json:"username" binding:"required"`
		Email           string `json:"email" binding:"required"`
		Password        string `json:"password" binding:"required"`
		PasswordConfirm string `json:"password_confirm" binding:"required"`
		Gender          string `json:"gender"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid params"})
		return
	}

	if req.Gender == "" {
		req.Gender = "Male"
	}

	err := r.client.Register(req.Username, req.Email, req.Password, req.PasswordConfirm, req.Gender)
	if err != nil {
		log.Errorf("Register failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Register failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Register success"})
}

func (r *Router) GetUserInfo(c *gin.Context) {
	info, err := r.client.GetUserInfo()
	if err != nil {
		log.Errorf("Get user info failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}

	c.JSON(http.StatusOK, info)
}

func (r *Router) Sign(c *gin.Context) {
	err := r.client.Sign()
	if err != nil {
		log.Errorf("Sign failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Sign failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Signed"})
}

// ==================== 帮助 ====================

func (r *Router) GetHelp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version": "1.0.0",
		"faq": []map[string]string{
			{"q": "如何搜索漫画？", "a": "在顶部搜索框输入关键词，按回车即可搜索。"},
			{"q": "如何收藏漫画？", "a": "在漫画详情页点击收藏按钮。"},
			{"q": "如何下载漫画？", "a": "在漫画详情页点击下载按钮，下载任务会添加到下载管理中。"},
			{"q": "如何切换阅读模式？", "a": "在设置页面可以切换滚动/翻页模式。"},
		},
	})
}

func (r *Router) processDownload(downloadID int, comicID string) {
	r.db.UpdateDownload(downloadID, "downloading", 0, 0)

	comic, err := r.client.GetComicDetail(comicID)
	if err != nil {
		r.db.UpdateDownload(downloadID, "failed", 0, 0)
		return
	}

	chapters, err := r.client.GetChapters(comicID)
	if err != nil {
		r.db.UpdateDownload(downloadID, "failed", 0, 0)
		return
	}

	totalChapters := len(chapters)
	for i, chapter := range chapters {
		progress := (i * 100) / totalChapters
		r.db.UpdateDownload(downloadID, "downloading", progress, i)

		_, err := r.client.GetChapterImages(chapter.ID)
		if err != nil {
			log.Errorf("Download chapter %s failed: %v", chapter.ID, err)
			continue
		}

		log.Infof("Downloaded chapter %s for comic %s", chapter.ID, comic.Title)
	}

	r.db.UpdateDownload(downloadID, "completed", 100, totalChapters)
}

func joinStrings(strs []string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += ","
		}
		result += s
	}
	return result
}
