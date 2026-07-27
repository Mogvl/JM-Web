package model

import "time"

type Comic struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Description string    `json:"description"`
	CoverURL    string    `json:"cover_url"`
	Tags        string    `json:"tags"`
	Category    string    `json:"category"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Chapter struct {
	ID         string    `json:"id" gorm:"primaryKey"`
	ComicID    string    `json:"comic_id"`
	Title      string    `json:"title"`
	SortOrder  int       `json:"sort_order"`
	ImageCount int       `json:"image_count"`
	CreatedAt  time.Time `json:"created_at"`
}

type Favorite struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	ComicID   string    `json:"comic_id" gorm:"uniqueIndex"`
	Comic     Comic     `json:"comic" gorm:"foreignKey:ComicID"`
	CreatedAt time.Time `json:"created_at"`
}

type History struct {
	ID         int       `json:"id" gorm:"primaryKey;autoIncrement"`
	ComicID    string    `json:"comic_id"`
	ChapterID  string    `json:"chapter_id"`
	Page       int       `json:"page"`
	Comic      Comic     `json:"comic" gorm:"foreignKey:ComicID"`
	LastReadAt time.Time `json:"last_read_at"`
}

type Download struct {
	ID          int        `json:"id" gorm:"primaryKey;autoIncrement"`
	ComicID     string     `json:"comic_id"`
	Status      string     `json:"status" gorm:"default:pending"`
	Progress    int        `json:"progress" gorm:"default:0"`
	TotalPages  int        `json:"total_pages"`
	Downloaded  int        `json:"downloaded"`
	FilePath    string     `json:"file_path"`
	Comic       Comic      `json:"comic" gorm:"foreignKey:ComicID"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at"`
}
