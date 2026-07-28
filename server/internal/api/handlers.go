package api

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Mogvl/JM-Web/server/internal/crawler"
	"github.com/Mogvl/JM-Web/server/internal/model"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// ==================== 搜索和浏览 ====================

func (r *Router) Search(c *gin.Context) {
	query := c.Query("q")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	sort := c.DefaultQuery("sort", "")

	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "q is required"})
		return
	}

	result, err := r.client.Search(query, page, sort)
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
	categories, blocks, err := r.client.GetCategories()
	if err != nil {
		log.Errorf("Get categories failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get categories"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"categories": categories,
		"blocks":     blocks,
	})
}

func (r *Router) GetCategoryFilter(c *gin.Context) {
	category := c.DefaultQuery("category", "0")
	sort := c.DefaultQuery("sort", "mr")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	result, err := r.client.GetCategoryFilter(category, sort, page)
	if err != nil {
		log.Errorf("Get category filter failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ==================== 收藏 ====================

func (r *Router) GetFavorites(c *gin.Context) {
	// 先从本地获取
	favorites, err := r.db.GetFavorites()
	if err == nil && len(favorites) > 0 {
		c.JSON(http.StatusOK, favorites)
		return
	}
	// 再从在线获取
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	online, err := r.client.GetOnlineFavorites(page)
	if err != nil {
		c.JSON(http.StatusOK, []interface{}{})
		return
	}
	c.JSON(http.StatusOK, online.Items)
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
		Title   string `json:"title"`
		Author  string `json:"author"`
		Cover   string `json:"cover"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "comic_id required"})
		return
	}

	// 保存漫画信息到数据库
	r.db.SaveComic(&model.Comic{
		ID:       req.ComicID,
		Title:    req.Title,
		Author:   req.Author,
		CoverURL: req.Cover,
	})

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

	promoteMap, err := r.client.GetIndexInfo(page)
	if err != nil {
		result, err2 := r.client.GetLatest(page)
		if err2 != nil {
			c.JSON(http.StatusOK, gin.H{"latest": []interface{}{}})
			return
		}
		c.JSON(http.StatusOK, gin.H{"latest": result.Items})
		return
	}

	c.JSON(http.StatusOK, promoteMap)
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

// ==================== 每周必看 ====================

func (r *Router) GetWeekCategories(c *gin.Context) {
	cats, err := r.client.GetWeekCategories()
	if err != nil {
		log.Errorf("Get week categories failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed"})
		return
	}
	c.JSON(http.StatusOK, cats)
}

func (r *Router) GetWeekFilter(c *gin.Context) {
	id := c.Query("id")
	catType := c.DefaultQuery("type", "manga")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))

	result, err := r.client.GetWeekFilter(id, catType, page)
	if err != nil {
		log.Errorf("Get week filter failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed"})
		return
	}
	c.JSON(http.StatusOK, result)
}

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

// ==================== 图片代理 ====================

func (r *Router) ProxyImage(c *gin.Context) {
	imgURL := c.Query("url")
	if imgURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url required"})
		return
	}

	req, err := http.NewRequest("GET", imgURL, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid url"})
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Referer", r.config.JMComic.BaseURL)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "fetch failed"})
		return
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/jpeg"
	}
	c.Header("Content-Type", contentType)
	c.Header("Cache-Control", "public, max-age=86400")
	c.Status(resp.StatusCode)
	io.Copy(c.Writer, resp.Body)
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

	// 先获取所有章节的图片列表，计算总数
	type chapterImages struct {
		chapter crawler.ChapterItem
		images  []string
	}
	allChapters := make([]chapterImages, 0, len(chapters))
	totalImgs := 0
	for _, ch := range chapters {
		imgs, err := r.client.GetChapterImages(ch.ID)
		if err != nil {
			log.Warnf("Get chapter %s images failed: %v", ch.ID, err)
			continue
		}
		allChapters = append(allChapters, chapterImages{chapter: ch, images: imgs})
		totalImgs += len(imgs)
	}
	r.db.SetDownloadTotal(downloadID, totalImgs)

	// 创建下载目录
	safeTitle := sanitizeFilename(comic.Title)
	if safeTitle == "" {
		safeTitle = comicID
	}
	comicDir := filepath.Join(r.config.Download.Dir, safeTitle+"_"+comicID)
	if err := os.MkdirAll(comicDir, 0755); err != nil {
		log.Errorf("Create download dir failed: %v", err)
		r.db.UpdateDownload(downloadID, "failed", 0, 0)
		return
	}

	downloadedImgs := 0
	for i, ci := range allChapters {
		progress := (downloadedImgs * 100) / maxInt(1, totalImgs)
		r.db.UpdateDownload(downloadID, "downloading", progress, downloadedImgs)

		chapterDir := filepath.Join(comicDir, fmt.Sprintf("%02d_%s", i+1, sanitizeFilename(ci.chapter.Title)))
		os.MkdirAll(chapterDir, 0755)

		for j, imgURL := range ci.images {
			if imgURL == "" {
				continue
			}
			ext := ".jpg"
			if strings.HasSuffix(imgURL, ".png") {
				ext = ".png"
			} else if strings.HasSuffix(imgURL, ".webp") {
				ext = ".webp"
			} else if strings.HasSuffix(imgURL, ".gif") {
				ext = ".gif"
			}
			imgPath := filepath.Join(chapterDir, fmt.Sprintf("%03d%s", j+1, ext))
			if err := r.client.DownloadImage(imgURL, imgPath); err != nil {
				log.Warnf("Download image %s failed: %v", imgURL, err)
				continue
			}
			downloadedImgs++
		}
		log.Infof("Downloaded chapter %s (%d images) for %s", ci.chapter.ID, len(ci.images), comic.Title)
	}

	r.db.UpdateDownload(downloadID, "completed", 100, downloadedImgs)
	r.db.SetDownloadPath(downloadID, comicDir)
	log.Infof("Download completed: %s -> %s (%d/%d images)", comic.Title, comicDir, downloadedImgs, totalImgs)
}

func maxInt(a, b int) int { if a > b { return a }; return b }

func sanitizeFilename(name string) string {
	name = strings.TrimSpace(name)
	for _, ch := range `<>:"/\|?*` {
		name = strings.ReplaceAll(name, string(ch), "_")
	}
	return name
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
