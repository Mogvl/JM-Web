package crawler

import (
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	appVersion    = "2.0.21"
	appSecret     = "18comicAPP"
	appDataSecret = "185Hcomic3PAPP7R"
)

type Client struct {
	baseURL    string
	auth       string
	httpClient *http.Client
	ts         int64
}

func NewClient(baseURL, auth string) *Client {
	return &Client{
		baseURL: baseURL,
		auth:    auth,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) SetAuth(token string) {
	c.auth = token
}

func (c *Client) getHeaders(method string) map[string]string {
	now := time.Now().Unix()
	c.ts = now
	param := fmt.Sprintf("%d%s", now, appSecret)
	token := fmt.Sprintf("%x", md5.Sum([]byte(param)))
	headers := map[string]string{
		"tokenparam":      fmt.Sprintf("%d,%s", now, appVersion),
		"token":           token,
		"accept-encoding": "gzip",
		"version":         appVersion,
	}
	if method == "POST" {
		headers["Content-Type"] = "application/x-www-form-urlencoded"
	}
	return headers
}

func (c *Client) readResponse(resp *http.Response) ([]byte, error) {
	var reader io.Reader
	if resp.Header.Get("Content-Encoding") == "gzip" {
		var err error
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
		defer reader.(*gzip.Reader).Close()
	} else {
		reader = resp.Body
	}
	return io.ReadAll(reader)
}

func (c *Client) decryptData(data string) (string, error) {
	dataB64, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	param := fmt.Sprintf("%d%s", c.ts, appDataSecret)
	keyBytes := md5.Sum([]byte(param))
	key := fmt.Sprintf("%x", keyBytes) // full 32 hex chars = 32 bytes
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	if len(dataB64)%blockSize != 0 {
		return "", fmt.Errorf("data not aligned to block size: %d %% %d", len(dataB64), blockSize)
	}
	decrypted := make([]byte, len(dataB64))
	for i := 0; i < len(dataB64); i += blockSize {
		block.Decrypt(decrypted[i:i+blockSize], dataB64[i:i+blockSize])
	}
	padding := int(decrypted[len(decrypted)-1])
	if padding > len(decrypted) || padding > blockSize {
		return string(decrypted), nil
	}
	decrypted = decrypted[:len(decrypted)-padding]
	return string(decrypted), nil
}

func (c *Client) doRequest(path string) ([]byte, error) {
	req, err := http.NewRequest("GET", c.baseURL+path, nil)
	if err != nil {
		return nil, err
	}
	return c.executeRequest(req, "GET")
}

func (c *Client) doPost(path string, data url.Values) ([]byte, error) {
	encoded := data.Encode()
	req, err := http.NewRequest("POST", c.baseURL+path, io.NopCloser(bytes.NewBufferString(encoded)))
	if err != nil {
		return nil, err
	}
	req.ContentLength = int64(len(encoded))
	return c.executeRequest(req, "POST")
}

func (c *Client) executeRequest(req *http.Request, method string) ([]byte, error) {
	headers := c.getHeaders(method)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	if c.auth != "" {
		req.Header.Set("Authorization", c.auth)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := c.readResponse(resp)
	if err != nil {
		return nil, err
	}

	// 检查是否需要解密
	var rawResp struct {
		Code     int             `json:"code"`
		Data     json.RawMessage `json:"data"`
		ErrorMsg string          `json:"errorMsg"`
	}
	if err := json.Unmarshal(body, &rawResp); err != nil {
		log.Debugf("Unmarshal failed: %v, body=%s", err, string(body[:min(100, len(body))]))
		return body, nil
	}
	if rawResp.Code == 200 && len(rawResp.Data) > 50 {
		var dataStr string
		if err := json.Unmarshal(rawResp.Data, &dataStr); err == nil && len(dataStr) > 50 {
			log.Debugf("Attempting decrypt, data len=%d", len(dataStr))
			decrypted, err := c.decryptData(dataStr)
			if err != nil {
				log.Warnf("decrypt failed: %v", err)
				return body, nil
			}
			log.Infof("Decrypted: %s", decrypted[:min(200, len(decrypted))])
			// 处理数组格式
			trimmed := strings.TrimSpace(decrypted)
			if strings.HasPrefix(trimmed, "[") {
				decrypted = `{"list":` + decrypted + `,"total":0,"limit":20,"offset":0}`
			}
			return []byte(fmt.Sprintf(`{"code":200,"data":%s,"errorMsg":""}`, decrypted)), nil
		}
	}
	return body, nil
}

func (c *Client) get(path string, params map[string]string) ([]byte, error) {
	fullPath := path
	if len(params) > 0 {
		values := url.Values{}
		for k, v := range params {
			values.Set(k, v)
		}
		fullPath = path + "/?" + values.Encode()
	}
	return c.doRequest(fullPath)
}

// ==================== 搜索 ====================

func (c *Client) Search(query string, page int) (*SearchResult, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * 20
	body, err := c.get("/api/search", map[string]string{
		"query":  query,
		"offset": fmt.Sprintf("%d", offset),
		"limit":  "20",
	})
	if err != nil {
		return nil, err
	}
	return c.parseAnyList(body, page)
}

func (c *Client) GetIndexInfo(page int) (*SearchResult, error) {
	body, err := c.get("/promote", map[string]string{"page": fmt.Sprintf("%d", page)})
	if err != nil {
		return nil, err
	}
	return c.parseAnyList(body, page+1)
}

func (c *Client) GetLatest(page int) (*SearchResult, error) {
	body, err := c.get("/latest", map[string]string{"page": fmt.Sprintf("%d", page)})
	if err != nil {
		return nil, err
	}
	return c.parseAnyList(body, page+1)
}

func (c *Client) GetComicDetail(comicID string) (*ComicDetail, error) {
	body, err := c.get("/album", map[string]string{
		"id":        comicID,
		"comicName": "",
	})
	if err != nil {
		return nil, err
	}

	log.Infof("Album raw: %s", string(body[:min(500, len(body))]))

	var resp struct {
		Code int             `json:"code"`
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	var dataArr []json.RawMessage
	if json.Unmarshal(resp.Data, &dataArr) == nil && len(dataArr) > 0 {
		return c.parseAlbumDetail(dataArr[0], comicID)
	}
	return c.parseAlbumDetail(resp.Data, comicID)
}

func (c *Client) parseAlbumDetail(raw json.RawMessage, comicID string) (*ComicDetail, error) {
	var album struct {
		ID          json.Number     `json:"id"`
		Name        string          `json:"name"`
		Author      json.RawMessage `json:"author"`
		Tags        json.RawMessage `json:"tags"`
		Description string          `json:"description"`
		Series      []struct {
			ID   json.Number `json:"id"`
			Sort int         `json:"sort"`
			Name string      `json:"name"`
		} `json:"series"`
		TotalPhotos int    `json:"total_photos"`
		TotalViews  string `json:"total_views"`
		Likes       string `json:"likes"`
		CommentTotal string `json:"comment_total"`
	}
	if err := json.Unmarshal(raw, &album); err != nil {
		return nil, err
	}

	author := parseAuthor(album.Author)
	tags := parseTags(album.Tags)

	detail := &ComicDetail{
		ID:          album.ID.String(),
		Title:       album.Name,
		Author:      author,
		Description: album.Description,
		CoverURL:    buildCoverURL(album.ID.String()),
		Tags:        tags,
	}

	// 解析章节
	for _, s := range album.Series {
		detail.Chapters = append(detail.Chapters, ChapterItem{
			ID:        s.ID.String(),
			Title:     s.Name,
			SortOrder: s.Sort,
		})
	}
	// 如果没有 series，创建一个默认章节
	if len(detail.Chapters) == 0 {
		detail.Chapters = append(detail.Chapters, ChapterItem{
			ID:        album.ID.String(),
			Title:     "第1章",
			SortOrder: 1,
		})
	}

	return detail, nil
}

func buildCoverURL(id string) string {
	return fmt.Sprintf("https://cdn-msp.jmapinodeudzn.net/media/albums/%s_3x4.jpg", id)
}

func (c *Client) GetChapters(comicID string) ([]ChapterItem, error) {
	// 章节信息从漫画详情中获取
	detail, err := c.GetComicDetail(comicID)
	if err != nil {
		return nil, err
	}
	return detail.Chapters, nil
}

func (c *Client) GetChapterImages(chapterID string) ([]string, error) {
	// 使用 /chapter_view_template 端点获取图片
	body, err := c.get("/chapter_view_template", map[string]string{
		"id":            chapterID,
		"mode":          "vertical",
		"page":          "0",
		"app_img_shunt": "NaN",
	})
	if err != nil {
		// 尝试备用端点
		body2, err2 := c.get("/chapter", map[string]string{
			"id":        chapterID,
			"comicName": "",
			"skip":      "",
		})
		if err2 != nil {
			return nil, err
		}
		body = body2
	}

	// 从 HTML 或 JSON 中提取图片 URL
	return c.parseChapterImages(body)
}

func (c *Client) parseChapterImages(body []byte) ([]string, error) {
	// 尝试 JSON 格式
	var resp struct {
		Code int             `json:"code"`
		Data json.RawMessage `json:"data"`
	}
	if json.Unmarshal(body, &resp) == nil {
		// 尝试数组格式
		var images []string
		if json.Unmarshal(resp.Data, &images) == nil {
			return images, nil
		}
		// 尝试对象格式
		var imgData struct {
			Images []string `json:"images"`
		}
		if json.Unmarshal(resp.Data, &imgData) == nil {
			return imgData.Images, nil
		}
	}

	return nil, fmt.Errorf("no images found")
}

func (c *Client) GetCategories() ([]CategoryItem, error) {
	body, err := c.doRequest("/categories")
	if err != nil {
		return nil, err
	}

	// 尝试解析: {"code":200,"data":{"categories":[...]},"errorMsg":""}
	var resp struct {
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	// data 是 {"categories": [...]}
	var dataObj struct {
		Categories json.RawMessage `json:"categories"`
	}
	if err := json.Unmarshal(resp.Data, &dataObj); err != nil {
		// 尝试直接解析 data 为数组
		var cats []CategoryData
		if json.Unmarshal(resp.Data, &cats) == nil {
			return toCategoryItems(cats), nil
		}
		return nil, err
	}

	var cats []CategoryData
	if err := json.Unmarshal(dataObj.Categories, &cats); err != nil {
		return nil, err
	}
	return toCategoryItems(cats), nil
}

func toCategoryItems(cats []CategoryData) []CategoryItem {
	items := make([]CategoryItem, len(cats))
	for i, cat := range cats {
		id := cat.PathWord
		if id == "" {
			id = cat.ID.String()
		}
		items[i] = CategoryItem{ID: id, Name: cat.Name}
	}
	return items
}

func (c *Client) GetRanking(rankType string, page int) (*SearchResult, error) {
	if page < 1 {
		page = 1
	}
	body, err := c.get("/api/ranking", map[string]string{
		"type":   rankType,
		"offset": fmt.Sprintf("%d", (page-1)*20),
		"limit":  "20",
	})
	if err != nil {
		return nil, err
	}
	return c.parseAnyList(body, page)
}

func (c *Client) GetCategoryFilter(category, sort string, page int) (*SearchResult, error) {
	params := map[string]string{
		"o": sort,
	}
	if page > 1 {
		params["page"] = fmt.Sprintf("%d", page)
	}
	if category != "" && category != "0" {
		params["category"] = category
	}

	body, err := c.get("/categories/filter", params)
	if err != nil {
		return nil, err
	}

	// 尝试解析
	var resp struct {
		Code int             `json:"code"`
		Data json.RawMessage `json:"data"`
	}
	if json.Unmarshal(body, &resp) != nil {
		return c.parseAnyList(body, page)
	}

	// data 可能是 {total, content}
	var catData struct {
		Total   int               `json:"total"`
		Content []json.RawMessage `json:"content"`
	}
	if json.Unmarshal(resp.Data, &catData) == nil {
		items := make([]ComicItem, 0)
		for _, item := range catData.Content {
			var raw RawComicItem
			if json.Unmarshal(item, &raw) != nil {
				continue
			}
			items = append(items, c.rawToItem(raw))
		}
		return &SearchResult{Items: items, TotalPages: (catData.Total + 19) / 20, Page: page}, nil
	}

	return c.parseAnyList(body, page)
}

func (c *Client) rawToItem(raw RawComicItem) ComicItem {
	author := parseAuthor(raw.Author)
	comicID := raw.PathWord
	if comicID == "" {
		comicID = raw.ID.String()
	}
	coverURL := raw.Cover
	if coverURL == "" {
		coverURL = raw.Image
	}
	if coverURL == "" {
		coverURL = buildCoverURL(comicID)
	}
	return ComicItem{ID: comicID, Title: raw.Name, Author: author, CoverURL: coverURL}
}

func (c *Client) GetOnlineFavorites(page int) (*SearchResult, error) {
	body, err := c.get("/favorite", map[string]string{
		"page":      fmt.Sprintf("%d", page),
		"folder_id": "0",
		"o":         "mr",
	})
	if err != nil {
		return nil, err
	}

	var resp struct {
		Code int             `json:"code"`
		Data json.RawMessage `json:"data"`
	}
	if json.Unmarshal(body, &resp) != nil || resp.Code != 200 {
		return c.parseAnyList(body, page)
	}

	// data 可能是 {list: [...]} 或 [...]
	return c.parseAnyList(body, page)
}

func (c *Client) GetOnlineHistory(page int) (*SearchResult, error) {
	body, err := c.get("/watch_list", map[string]string{
		"page": fmt.Sprintf("%d", page),
	})
	if err != nil {
		return nil, err
	}
	return c.parseAnyList(body, page)
}

func (c *Client) GetComments(comicID string, page int) ([]CommentItem, error) {
	if page < 1 {
		page = 1
	}
	body, err := c.get("/api/comic/"+comicID+"/comments", map[string]string{
		"offset": fmt.Sprintf("%d", (page-1)*20),
		"limit":  "20",
	})
	if err != nil {
		return nil, err
	}
	var resp struct {
		Code int             `json:"code"`
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	var comments []CommentData
	if err := json.Unmarshal(resp.Data, &comments); err != nil {
		return nil, err
	}
	items := make([]CommentItem, len(comments))
	for i, comment := range comments {
		items[i] = CommentItem{
			ID:         comment.ID,
			Content:    comment.Content,
			Author:     comment.Author,
			Avatar:     comment.Avatar,
			LikeCount:  comment.LikeCount,
			CreateTime: comment.CreateTime,
		}
	}
	return items, nil
}

func (c *Client) GetSubComments(commentID string, page int) ([]CommentItem, error) {
	body, err := c.get("/api/comment/"+commentID+"/sub", map[string]string{
		"offset": fmt.Sprintf("%d", (page-1)*20),
		"limit":  "20",
	})
	if err != nil {
		return nil, err
	}
	var resp struct {
		Code int             `json:"code"`
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	var comments []CommentData
	if err := json.Unmarshal(resp.Data, &comments); err != nil {
		return nil, err
	}
	items := make([]CommentItem, len(comments))
	for i, comment := range comments {
		items[i] = CommentItem{
			ID:         comment.ID,
			Content:    comment.Content,
			Author:     comment.Author,
			Avatar:     comment.Avatar,
			LikeCount:  comment.LikeCount,
			CreateTime: comment.CreateTime,
		}
	}
	return items, nil
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

func (c *Client) Login(username, password string) (*LoginUserData, error) {
	body, err := c.doPost("/login", url.Values{
		"username": {username},
		"password": {password},
	})
	if err != nil {
		return nil, err
	}

	var loginResp struct {
		Code     int             `json:"code"`
		Data     json.RawMessage `json:"data"`
		ErrorMsg string          `json:"errorMsg"`
	}
	if err := json.Unmarshal(body, &loginResp); err != nil {
		return nil, fmt.Errorf("parse failed: %s", string(body))
	}

	if loginResp.Code != 200 {
		return nil, fmt.Errorf("login failed: %s", loginResp.ErrorMsg)
	}

	var userData LoginUserData
	if err := json.Unmarshal(loginResp.Data, &userData); err != nil {
		log.Warnf("Parse user data failed: %v, data: %s...", err, string(loginResp.Data)[:min(300, len(loginResp.Data))])
		return nil, fmt.Errorf("parse user data failed")
	}

	log.Infof("Login user: %s, coin=%d, level=%d, token=%s", userData.Username, userData.Coin, userData.Level, userData.JWTToken[:min(20, len(userData.JWTToken))])

	if userData.JWTToken == "" {
		return nil, fmt.Errorf("no token in response")
	}

	c.auth = "Bearer " + userData.JWTToken
	return &userData, nil
}

func (c *Client) Register(username, email, password, passwordConfirm, gender string) error {
	_, err := c.doPost("/signup", url.Values{
		"username":         {username},
		"email":            {email},
		"password":         {password},
		"password_confirm": {passwordConfirm},
		"gender":           {gender},
		"age":              {"on"},
		"terms":            {"on"},
		"submit_signup":    {""},
	})
	return err
}

func (c *Client) GetUserInfo() (*UserInfo, error) {
	// 用户信息从登录响应获取，调用 /user 接口
	body, err := c.doRequest("/user")
	if err == nil {
		var resp struct {
			Code int             `json:"code"`
			Data json.RawMessage `json:"data"`
		}
		if json.Unmarshal(body, &resp) == nil && resp.Code == 200 {
			var userInfo UserInfo
			if json.Unmarshal(resp.Data, &userInfo) == nil {
				return &userInfo, nil
			}
		}
	}
	return nil, fmt.Errorf("not logged in")
}

func (c *Client) Sign() error {
	_, err := c.doRequest("/api/user/sign")
	return err
}

// ==================== 解析通用漫画列表 ====================

func parseAuthor(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	if raw[0] == '[' {
		var authors []Author
		if json.Unmarshal(raw, &authors) == nil && len(authors) > 0 {
			return authors[0].Name
		}
		var strArr []string
		if json.Unmarshal(raw, &strArr) == nil && len(strArr) > 0 {
			return strArr[0]
		}
	} else {
		var author string
		if json.Unmarshal(raw, &author) == nil {
			return author
		}
	}
	return ""
}

func parseTags(raw json.RawMessage) []string {
	if len(raw) == 0 {
		return nil
	}
	var strTags []string
	if json.Unmarshal(raw, &strTags) == nil {
		return strTags
	}
	var objTags []Tag
	if json.Unmarshal(raw, &objTags) == nil {
		tags := make([]string, len(objTags))
		for i, t := range objTags {
			tags[i] = t.Name
		}
		return tags
	}
	return nil
}

func (c *Client) parseAnyList(body []byte, page int) (*SearchResult, error) {
	// 尝试用标准格式解析
	result, err := c.parseComicList(body, page)
	if err == nil {
		return result, nil
	}
	// 尝试用数组格式解析
	var arrResp struct {
		Code int               `json:"code"`
		Data []json.RawMessage `json:"data"`
	}
	if json.Unmarshal(body, &arrResp) == nil {
		items := make([]ComicItem, 0)
		for _, item := range arrResp.Data {
			var raw RawComicItem
			if json.Unmarshal(item, &raw) != nil {
				continue
			}
			author := parseAuthor(raw.Author)
			comicID := raw.PathWord
			if comicID == "" {
				comicID = raw.ID.String()
			}
			coverURL := raw.Cover
			if coverURL == "" {
				coverURL = raw.Image
			}
			if coverURL == "" {
				coverURL = buildCoverURL(comicID)
			}
			items = append(items, ComicItem{
				ID:       comicID,
				Title:    raw.Name,
				Author:   author,
				CoverURL: coverURL,
			})
		}
		return &SearchResult{Items: items, TotalPages: 1, Page: page}, nil
	}
	return result, err
}

func (c *Client) parseComicList(body []byte, page int) (*SearchResult, error) {
	var resp struct {
		Code int             `json:"code"`
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	var searchData SearchData
	if err := json.Unmarshal(resp.Data, &searchData); err != nil {
		return nil, err
	}

	results := make([]ComicItem, len(searchData.List))
	for i, item := range searchData.List {
		author := parseAuthor(item.Author)
		// 兼容 path_word/id 和 cover/image
		comicID := item.PathWord
		if comicID == "" {
			comicID = item.ID.String()
		}
		coverURL := item.Cover
		if coverURL == "" {
			coverURL = item.Image
		}
		results[i] = ComicItem{
			ID:       comicID,
			Title:    item.Name,
			Author:   author,
			CoverURL: coverURL,
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
