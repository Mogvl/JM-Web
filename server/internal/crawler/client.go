package crawler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"
)

type Client struct {
	baseURL    string
	auth       string
	httpClient *http.Client
}

func NewClient(baseURL, auth string) *Client {
	return &Client{
		baseURL:    baseURL,
		auth:       auth,
		httpClient: &http.Client{},
	}
}

func (c *Client) doRequest(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	if c.auth != "" {
		req.Header.Set("Authorization", c.auth)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func (c *Client) doRequestWithMethod(method, url string, body interface{}) ([]byte, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonBytes, _ := json.Marshal(body)
		bodyReader = bytes.NewBuffer(jsonBytes)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Content-Type", "application/json")
	if c.auth != "" {
		req.Header.Set("Authorization", c.auth)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func (c *Client) Search(query string, page int) (*SearchResult, error) {
	if page < 1 {
		page = 1
	}

	offset := (page - 1) * 20
	searchURL := fmt.Sprintf("%s/api/search?query=%s&offset=%d&limit=20", c.baseURL, url.QueryEscape(query), offset)

	body, err := c.doRequest(searchURL)
	if err != nil {
		log.Errorf("Search failed: %v", err)
		return nil, err
	}

	var resp APIResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	dataBytes, _ := json.Marshal(resp.Data)
	var searchData SearchData
	if err := json.Unmarshal(dataBytes, &searchData); err != nil {
		return nil, err
	}

	results := make([]ComicItem, len(searchData.List))
	for i, item := range searchData.List {
		author := ""
		if len(item.Author) > 0 {
			author = item.Author[0].Name
		}
		results[i] = ComicItem{
			ID:       item.PathWord,
			Title:    item.Name,
			Author:   author,
			CoverURL: item.Cover,
		}
	}

	totalPages := 1
	if searchData.Limit > 0 {
		totalPages = (searchData.Total + searchData.Limit - 1) / searchData.Limit
	}

	return &SearchResult{
		Items:      results,
		TotalPages: totalPages,
		Page:       page,
	}, nil
}

func (c *Client) GetComicDetail(comicID string) (*ComicDetail, error) {
	detailURL := fmt.Sprintf("%s/api/comic/%s", c.baseURL, comicID)

	body, err := c.doRequest(detailURL)
	if err != nil {
		log.Errorf("Get comic detail failed: %v", err)
		return nil, err
	}

	var resp APIResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	dataBytes, _ := json.Marshal(resp.Data)
	var detail ComicDetailData
	if err := json.Unmarshal(dataBytes, &detail); err != nil {
		return nil, err
	}

	author := ""
	if len(detail.Author) > 0 {
		author = detail.Author[0].Name
	}

	tags := make([]string, len(detail.Tags))
	for i, t := range detail.Tags {
		tags[i] = t.Name
	}

	return &ComicDetail{
		ID:          detail.PathWord,
		Title:       detail.Name,
		Author:      author,
		Description: detail.Description,
		CoverURL:    detail.Cover,
		Tags:        tags,
		Category:    detail.Category,
		Status:      detail.Status,
	}, nil
}

func (c *Client) GetChapters(comicID string) ([]ChapterItem, error) {
	chaptersURL := fmt.Sprintf("%s/api/comic/%s/chapters", c.baseURL, comicID)

	body, err := c.doRequest(chaptersURL)
	if err != nil {
		log.Errorf("Get chapters failed: %v", err)
		return nil, err
	}

	var resp APIResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	dataBytes, _ := json.Marshal(resp.Data)
	var chapters []ChapterData
	if err := json.Unmarshal(dataBytes, &chapters); err != nil {
		return nil, err
	}

	items := make([]ChapterItem, len(chapters))
	for i, ch := range chapters {
		items[i] = ChapterItem{
			ID:        ch.PathWord,
			Title:     ch.Name,
			SortOrder: i,
		}
	}

	return items, nil
}

func (c *Client) GetChapterImages(chapterID string) ([]string, error) {
	chapterURL := fmt.Sprintf("%s/api/chapter/%s", c.baseURL, chapterID)

	body, err := c.doRequest(chapterURL)
	if err != nil {
		log.Errorf("Get chapter images failed: %v", err)
		return nil, err
	}

	var resp APIResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	dataBytes, _ := json.Marshal(resp.Data)
	var imagesData ImagesData
	if err := json.Unmarshal(dataBytes, &imagesData); err != nil {
		return nil, err
	}

	return imagesData.Images, nil
}

func (c *Client) GetCategories() ([]CategoryItem, error) {
	categoriesURL := fmt.Sprintf("%s/api/categories", c.baseURL)

	body, err := c.doRequest(categoriesURL)
	if err != nil {
		log.Errorf("Get categories failed: %v", err)
		return nil, err
	}

	var resp APIResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	dataBytes, _ := json.Marshal(resp.Data)
	var categories []CategoryData
	if err := json.Unmarshal(dataBytes, &categories); err != nil {
		return nil, err
	}

	items := make([]CategoryItem, len(categories))
	for i, cat := range categories {
		items[i] = CategoryItem{
			ID:   cat.PathWord,
			Name: cat.Name,
		}
	}

	return items, nil
}

func (c *Client) GetRanking(rankType string, page int) (*SearchResult, error) {
	if page < 1 {
		page = 1
	}

	offset := (page - 1) * 20
	rankingURL := fmt.Sprintf("%s/api/ranking?type=%s&offset=%d&limit=20", c.baseURL, rankType, offset)

	body, err := c.doRequest(rankingURL)
	if err != nil {
		log.Errorf("Get ranking failed: %v", err)
		return nil, err
	}

	var resp APIResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	dataBytes, _ := json.Marshal(resp.Data)
	var searchData SearchData
	if err := json.Unmarshal(dataBytes, &searchData); err != nil {
		return nil, err
	}

	results := make([]ComicItem, len(searchData.List))
	for i, item := range searchData.List {
		author := ""
		if len(item.Author) > 0 {
			author = item.Author[0].Name
		}
		results[i] = ComicItem{
			ID:       item.PathWord,
			Title:    item.Name,
			Author:   author,
			CoverURL: item.Cover,
		}
	}

	totalPages := 1
	if searchData.Limit > 0 {
		totalPages = (searchData.Total + searchData.Limit - 1) / searchData.Limit
	}

	return &SearchResult{
		Items:      results,
		TotalPages: totalPages,
		Page:       page,
	}, nil
}

func (c *Client) GetComments(comicID string, page int) ([]CommentItem, error) {
	if page < 1 {
		page = 1
	}

	offset := (page - 1) * 20
	commentsURL := fmt.Sprintf("%s/api/comic/%s/comments?offset=%d&limit=20", c.baseURL, comicID, offset)

	body, err := c.doRequest(commentsURL)
	if err != nil {
		log.Errorf("Get comments failed: %v", err)
		return nil, err
	}

	var resp APIResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	dataBytes, _ := json.Marshal(resp.Data)
	var comments []CommentData
	if err := json.Unmarshal(dataBytes, &comments); err != nil {
		return nil, err
	}

	items := make([]CommentItem, len(comments))
	for i, comment := range comments {
		items[i] = CommentItem{
			ID:        comment.ID,
			Content:   comment.Content,
			Author:    comment.Author,
			Avatar:    comment.Avatar,
			LikeCount: comment.LikeCount,
			CreateTime: comment.CreateTime,
		}
	}

	return items, nil
}

func (c *Client) GetSubComments(commentID string, page int) ([]CommentItem, error) {
	if page < 1 {
		page = 1
	}

	offset := (page - 1) * 20
	commentsURL := fmt.Sprintf("%s/api/comment/%s/sub?offset=%d&limit=20", c.baseURL, commentID, offset)

	body, err := c.doRequest(commentsURL)
	if err != nil {
		log.Errorf("Get sub comments failed: %v", err)
		return nil, err
	}

	var resp APIResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	dataBytes, _ := json.Marshal(resp.Data)
	var comments []CommentData
	if err := json.Unmarshal(dataBytes, &comments); err != nil {
		return nil, err
	}

	items := make([]CommentItem, len(comments))
	for i, comment := range comments {
		items[i] = CommentItem{
			ID:        comment.ID,
			Content:   comment.Content,
			Author:    comment.Author,
			Avatar:    comment.Avatar,
			LikeCount: comment.LikeCount,
			CreateTime: comment.CreateTime,
		}
	}

	return items, nil
}

func (c *Client) Login(username, password string) (string, error) {
	loginURL := fmt.Sprintf("%s/api/login", c.baseURL)

	body, err := c.doRequestWithMethod("POST", loginURL, map[string]string{
		"username": username,
		"password": password,
	})
	if err != nil {
		log.Errorf("Login failed: %v", err)
		return "", err
	}

	var resp APIResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return "", err
	}

	dataBytes, _ := json.Marshal(resp.Data)
	var loginData LoginData
	if err := json.Unmarshal(dataBytes, &loginData); err != nil {
		return "", err
	}

	c.auth = loginData.Token
	return loginData.Token, nil
}

func (c *Client) GetUserInfo() (*UserInfo, error) {
	userInfoURL := fmt.Sprintf("%s/api/user/info", c.baseURL)

	body, err := c.doRequest(userInfoURL)
	if err != nil {
		log.Errorf("Get user info failed: %v", err)
		return nil, err
	}

	var resp APIResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	dataBytes, _ := json.Marshal(resp.Data)
	var userInfo UserInfo
	if err := json.Unmarshal(dataBytes, &userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

func (c *Client) Sign() error {
	signURL := fmt.Sprintf("%s/api/user/sign", c.baseURL)

	_, err := c.doRequestWithMethod("POST", signURL, nil)
	if err != nil {
		log.Errorf("Sign failed: %v", err)
		return err
	}

	return nil
}
