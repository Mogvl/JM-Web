package crawler

import "encoding/json"

type RawComicItem struct {
	Name     string          `json:"name"`
	PathWord string          `json:"path_word"`
	ID       json.Number     `json:"id"`
	Cover    string          `json:"cover"`
	Image    string          `json:"image"`
	Author   json.RawMessage `json:"author"`
}

type ComicItem struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	CoverURL string `json:"cover_url"`
}

type ChapterItem struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	SortOrder int    `json:"sort_order"`
}

type SearchResult struct {
	Items      []ComicItem `json:"items"`
	TotalPages int         `json:"total_pages"`
	Page       int         `json:"page"`
}

type SearchData struct {
	List   []RawComicItem `json:"list"`
	Total  int            `json:"total"`
	Limit  int            `json:"limit"`
	Offset int            `json:"offset"`
}

type Author struct {
	Name     string `json:"name"`
	PathWord string `json:"path_word"`
}

type Tag struct {
	Name string `json:"name"`
}

type CategoryData struct {
	ID            json.Number   `json:"id"`
	Name          string        `json:"name"`
	PathWord      string        `json:"slug"`
	Type          string        `json:"type"`
	TotalAlbums   json.Number   `json:"total_albums"`
	SubCategories []SubCategory `json:"sub_categories"`
}

type SubCategory struct {
	CID  string `json:"CID"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type CategoryItem struct {
	ID             string         `json:"id"`
	Name           string         `json:"name"`
	Type           string         `json:"type"`
	TotalAlbums    string         `json:"total_albums"`
	SubCategories  []SubCategory  `json:"sub_categories"`
}

type CommentData struct {
	ID         string `json:"id"`
	Content    string `json:"content"`
	Author     string `json:"author"`
	Avatar     string `json:"avatar"`
	LikeCount  int    `json:"like_count"`
	CreateTime string `json:"create_time"`
}

type CommentItem struct {
	ID         string `json:"id"`
	Content    string `json:"content"`
	Author     string `json:"author"`
	Avatar     string `json:"avatar"`
	LikeCount  int    `json:"like_count"`
	CreateTime string `json:"create_time"`
}

type AlbumDetail struct {
	ID           string        `json:"id"`
	Title        string        `json:"title"`
	Author       string        `json:"author"`
	Tags         []string      `json:"tags"`
	Description  string        `json:"description"`
	CoverURL     string        `json:"cover_url"`
	TotalPhotos  int           `json:"total_photos"`
	TotalViews   string        `json:"total_views"`
	Likes        string        `json:"likes"`
	CommentTotal string        `json:"comment_total"`
	Chapters     []ChapterItem `json:"chapters"`
}

type LoginUserData struct {
	UID           string `json:"uid"`
	Username      string `json:"username"`
	Email         string `json:"email"`
	Photo         string `json:"photo"`
	Gender        string `json:"gender"`
	Coin          int    `json:"coin"`
	Level         int    `json:"level"`
	LevelName     string `json:"level_name"`
	Exp           string `json:"exp"`
	AlbumFavorites int   `json:"album_favorites"`
	JWTToken      string `json:"jwttoken"`
}

type CategoryBlock struct {
	Title   string   `json:"title"`
	Content []string `json:"content"`
}
