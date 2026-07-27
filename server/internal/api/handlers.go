package api

import (
	"net/http"
	"strconv"

	"github.com/Mogvl/JM-Web/server/internal/model"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

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
