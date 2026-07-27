package crawler

type APIResponse struct {
	Data    interface{} `json:"data"`
	Success bool        `json:"success"`
}

type SearchData struct {
	List   []RawComicItem `json:"list"`
	Total  int            `json:"total"`
	Limit  int            `json:"limit"`
	Offset int            `json:"offset"`
}

type RawComicItem struct {
	Name     string   `json:"name"`
	PathWord string   `json:"path_word"`
	Cover    string   `json:"cover"`
	Author   []Author `json:"author"`
}

type Author struct {
	Name     string `json:"name"`
	PathWord string `json:"path_word"`
}

type ComicDetailData struct {
	Name        string   `json:"name"`
	PathWord    string   `json:"path_word"`
	Cover       string   `json:"cover"`
	Description string   `json:"description"`
	Author      []Author `json:"author"`
	Tags        []Tag    `json:"tags"`
	Category    string   `json:"category"`
	Status      string   `json:"status"`
}

type Tag struct {
	Name string `json:"name"`
}

type ChapterData struct {
	Name     string `json:"name"`
	PathWord string `json:"path_word"`
}

type ImagesData struct {
	Images []string `json:"images"`
}

type CategoryData struct {
	Name     string `json:"name"`
	PathWord string `json:"path_word"`
}

type ComicItem struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	CoverURL string `json:"cover_url"`
}

type ComicDetail struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Author      string   `json:"author"`
	Description string   `json:"description"`
	CoverURL    string   `json:"cover_url"`
	Tags        []string `json:"tags"`
	Category    string   `json:"category"`
	Status      string   `json:"status"`
}

type ChapterItem struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	SortOrder int    `json:"sort_order"`
}

type CategoryItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type SearchResult struct {
	Items      []ComicItem `json:"items"`
	TotalPages int         `json:"total_pages"`
	Page       int         `json:"page"`
}
