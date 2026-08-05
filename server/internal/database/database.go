package database

import (
	"os"
	"path/filepath"
	"time"

	"github.com/Mogvl/JM-Web/server/internal/model"
	"github.com/glebarez/sqlite"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

func New(dbPath string) (*DB, error) {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Errorf("Failed to connect to database: %v", err)
		return nil, err
	}

	if err := db.AutoMigrate(
		&model.Comic{},
		&model.Chapter{},
		&model.Favorite{},
		&model.History{},
		&model.Download{},
	); err != nil {
		log.Errorf("Failed to migrate database: %v", err)
		return nil, err
	}

	log.Info("Database initialized")
	return &DB{db}, nil
}

// ==================== Comic ====================

func (db *DB) SaveComic(comic *model.Comic) error {
	return db.Save(comic).Error
}

func (db *DB) GetComic(id string) (*model.Comic, error) {
	var comic model.Comic
	err := db.First(&comic, "id = ?", id).Error
	return &comic, err
}

// ==================== Chapter ====================

func (db *DB) SaveChapter(chapter *model.Chapter) error {
	return db.Save(chapter).Error
}

func (db *DB) GetChapters(comicID string) ([]model.Chapter, error) {
	var chapters []model.Chapter
	err := db.Where("comic_id = ?", comicID).Order("sort_order").Find(&chapters).Error
	return chapters, err
}

// ==================== Favorite ====================

func (db *DB) AddFavorite(comicID string) error {
	return db.Create(&model.Favorite{ComicID: comicID}).Error
}

func (db *DB) RemoveFavorite(comicID string) error {
	return db.Where("comic_id = ?", comicID).Delete(&model.Favorite{}).Error
}

func (db *DB) IsFavorite(comicID string) bool {
	var count int64
	db.Model(&model.Favorite{}).Where("comic_id = ?", comicID).Count(&count)
	return count > 0
}

func (db *DB) GetFavorites() ([]model.Favorite, error) {
	var favorites []model.Favorite
	err := db.Preload("Comic").Order("created_at DESC").Find(&favorites).Error
	return favorites, err
}

// ==================== History ====================

func (db *DB) AddHistory(comicID, chapterID string, page int) error {
	history := &model.History{
		ComicID:   comicID,
		ChapterID: chapterID,
		Page:      page,
	}
	return db.Save(history).Error
}

// UpsertHistory 记录阅读进度（同漫画更新章节/页码/时间）
func (db *DB) UpsertHistory(comicID, chapterID string, page int) error {
	var h model.History
	err := db.Where("comic_id = ?", comicID).First(&h).Error
	if err == nil {
		h.ChapterID = chapterID
		h.Page = page
		h.LastReadAt = time.Now()
		return db.Save(&h).Error
	}
	if err == gorm.ErrRecordNotFound {
		return db.Create(&model.History{
			ComicID:    comicID,
			ChapterID:  chapterID,
			Page:       page,
			LastReadAt: time.Now(),
		}).Error
	}
	return err
}

// GetHistoryByComic 取某漫画最近的阅读进度
func (db *DB) GetHistoryByComic(comicID string) (*model.History, error) {
	var h model.History
	err := db.Preload("Comic").Where("comic_id = ?", comicID).Order("last_read_at DESC").First(&h).Error
	if err != nil {
		return nil, err
	}
	return &h, nil
}

func (db *DB) GetHistory() ([]model.History, error) {
	var history []model.History
	err := db.Preload("Comic").Order("last_read_at DESC").Find(&history).Error
	return history, err
}

func (db *DB) DeleteHistory(id int) error {
	return db.Delete(&model.History{}, id).Error
}

func (db *DB) ClearHistory() error {
	return db.Where("1 = 1").Delete(&model.History{}).Error
}

// ==================== Download ====================

func (db *DB) CreateDownload(comicID string) (*model.Download, error) {
	dl := &model.Download{
		ComicID: comicID,
		Status:  "pending",
	}
	err := db.Create(dl).Error
	return dl, err
}

func (db *DB) UpdateDownload(id int, status string, progress, downloaded int) error {
	return db.Model(&model.Download{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":     status,
		"progress":   progress,
		"downloaded": downloaded,
	}).Error
}

func (db *DB) GetDownloads() ([]model.Download, error) {
	var downloads []model.Download
	err := db.Preload("Comic").Order("created_at DESC").Find(&downloads).Error
	return downloads, err
}

func (db *DB) GetDownload(id int) (*model.Download, error) {
	var dl model.Download
	err := db.Preload("Comic").First(&dl, id).Error
	return &dl, err
}

func (db *DB) DeleteDownload(id int) error {
	return db.Delete(&model.Download{}, id).Error
}

func (db *DB) ClearDownloads() error {
	return db.Where("1 = 1").Delete(&model.Download{}).Error
}

func (db *DB) SetDownloadPath(id int, path string) error {
	return db.Model(&model.Download{}).Where("id = ?", id).Update("file_path", path).Error
}

func (db *DB) SetDownloadTotal(id int, total int) error {
	return db.Model(&model.Download{}).Where("id = ?", id).Update("total_pages", total).Error
}
