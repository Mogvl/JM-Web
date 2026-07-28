package crawler

import (
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
	if strings.HasPrefix(token, "Bearer ") {
		c.auth = token
	} else {
		c.auth = "Bearer " + token
	}
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
	key := fmt.Sprintf("%x", keyBytes)
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	if len(dataB64)%blockSize != 0 {
		return "", fmt.Errorf("data not aligned to block size")
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

func (c *Client) doPostForm(path string, data url.Values) ([]byte, error) {
	encoded := data.Encode()
	req, err := http.NewRequest("POST", c.baseURL+path, strings.NewReader(encoded))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
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

	var rawResp struct {
		Code     int             `json:"code"`
		Data     json.RawMessage `json:"data"`
		ErrorMsg string          `json:"errorMsg"`
	}
	if err := json.Unmarshal(body, &rawResp); err != nil {
		return body, nil
	}
	if rawResp.Code == 200 && len(rawResp.Data) > 50 {
		var dataStr string
		if err := json.Unmarshal(rawResp.Data, &dataStr); err == nil && len(dataStr) > 50 {
			decrypted, err := c.decryptData(dataStr)
			if err != nil {
				fmt.Printf("DECRYPT FAIL: %v, data=%s...\n", err, dataStr[:60])
				return body, nil
			}
			fmt.Printf("DECRYPT OK: %s...\n", decrypted[:min(100, len(decrypted))])
			return []byte(fmt.Sprintf(`{"code":200,"data":%s,"errorMsg":""}`, decrypted)), nil
		}
	}
	return body, nil
}

func (c *Client) get(path string, params map[string]string) ([]byte, error) {
	fullPath := c.baseURL + path
	if len(params) > 0 {
		values := url.Values{}
		for k, v := range params {
			values.Set(k, v)
		}
		fullPath = fullPath + "/?" + values.Encode()
	}
	req, err := http.NewRequest("GET", fullPath, nil)
	if err != nil {
		return nil, err
	}
	return c.executeRequest(req, "GET")
}

func (c *Client) Search(query string, page int, sort string) (*SearchResult, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * 20
	params := map[string]string{
		"query":  query,
		"offset": fmt.Sprintf("%d", offset),
		"limit":  "20",
	}
	if sort != "" {
		params["o"] = sort
	}
	body, err := c.get("/api/search", params)
	if err != nil {
		return nil, err
	}

	result, err := c.parseAnyList(body, page)
	if err != nil {
		// 尝试直接解析为最新格式
		var fallback struct {
			Code int `json:"code"`
			Data json.RawMessage `json:"data"`
		}
		if json.Unmarshal(body, &fallback) == nil && fallback.Code == 200 {
			// data 可能是空数组或空对象
			return &SearchResult{Items: []ComicItem{}, Page: page}, nil
		}
		return nil, err
	}
	return result, nil
}

func (c *Client) GetLatest(page int) (*SearchResult, error) {
	body, err := c.get("/latest", map[string]string{"page": fmt.Sprintf("%d", page)})
	if err != nil {
		return nil, err
	}
	return c.parseAnyList(body, page+1)
}

func (c *Client) GetComicDetail(comicID string) (*AlbumDetail, error) {
	body, err := c.get("/album", map[string]string{"id": comicID, "comicName": ""})
	if err != nil {
		return nil, err
	}

	var resp struct {
		Code int             `json:"code"`
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(body, &resp); err != nil || resp.Code != 200 {
		return nil, fmt.Errorf("failed to get album")
	}

	// 有可能是数组
	var arr []json.RawMessage
	if json.Unmarshal(resp.Data, &arr) == nil && len(arr) > 0 {
		return parseAlbum(arr[0])
	}
	return parseAlbum(resp.Data)
}

func parseAlbum(raw json.RawMessage) (*AlbumDetail, error) {
	var rawAlbum struct {
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
		TotalPhotos  int    `json:"total_photos"`
		TotalViews   string `json:"total_views"`
		Likes        string `json:"likes"`
		CommentTotal string `json:"comment_total"`
	}
	if err := json.Unmarshal(raw, &rawAlbum); err != nil {
		return nil, err
	}

	detail := &AlbumDetail{
		ID:           rawAlbum.ID.String(),
		Title:        rawAlbum.Name,
		Author:       parseAuthor(rawAlbum.Author),
		Tags:         parseTags(rawAlbum.Tags),
		Description:  rawAlbum.Description,
		CoverURL:     buildCoverURL(rawAlbum.ID.String()),
		TotalPhotos:  rawAlbum.TotalPhotos,
		TotalViews:   rawAlbum.TotalViews,
		Likes:        rawAlbum.Likes,
		CommentTotal: rawAlbum.CommentTotal,
	}

	if len(rawAlbum.Series) > 0 {
		for _, s := range rawAlbum.Series {
			detail.Chapters = append(detail.Chapters, ChapterItem{
				ID: s.ID.String(), Title: s.Name, SortOrder: s.Sort,
			})
		}
	} else {
		detail.Chapters = append(detail.Chapters, ChapterItem{
			ID: rawAlbum.ID.String(), Title: "第1章", SortOrder: 1,
		})
	}
	return detail, nil
}

func buildCoverURL(id string) string {
	return fmt.Sprintf("https://cdn-msp.jmapinodeudzn.net/media/albums/%s_3x4.jpg", id)
}

func (c *Client) GetChapters(comicID string) ([]ChapterItem, error) {
	detail, err := c.GetComicDetail(comicID)
	if err != nil {
		return nil, err
	}
	return detail.Chapters, nil
}

func (c *Client) GetChapterImages(chapterID string) ([]string, error) {
	body, err := c.get("/chapter_view_template", map[string]string{
		"id": chapterID, "mode": "vertical", "page": "0", "app_img_shunt": "NaN",
	})
	if err != nil {
		body, err = c.get("/chapter", map[string]string{"id": chapterID, "comicName": "", "skip": ""})
		if err != nil {
			return nil, err
		}
	}

	var resp struct {
		Code int             `json:"code"`
		Data json.RawMessage `json:"data"`
	}
	if json.Unmarshal(body, &resp) != nil {
		return nil, fmt.Errorf("no images found")
	}

	var arr []string
	if json.Unmarshal(resp.Data, &arr) == nil {
		return arr, nil
	}
	var obj struct {
		Images []string `json:"images"`
	}
	if json.Unmarshal(resp.Data, &obj) == nil {
		return obj.Images, nil
	}
	return nil, fmt.Errorf("no images found")
}

func (c *Client) GetCategories() ([]CategoryItem, error) {
	body, err := c.get("/categories", nil)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Data json.RawMessage `json:"data"`
	}
	if json.Unmarshal(body, &resp) != nil {
		return nil, err
	}

	var dataObj struct {
		Categories json.RawMessage `json:"categories"`
	}
	if json.Unmarshal(resp.Data, &dataObj) != nil {
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

func (c *Client) GetCategoryFilter(category, sort string, page int) (*SearchResult, error) {
	params := map[string]string{"o": sort}
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

	var resp struct {
		Data json.RawMessage `json:"data"`
	}
	if json.Unmarshal(body, &resp) != nil {
		return c.parseAnyList(body, page)
	}

	var catData struct {
		Total   int               `json:"total"`
		Content []json.RawMessage `json:"content"`
	}
	if json.Unmarshal(resp.Data, &catData) == nil {
		items := make([]ComicItem, len(catData.Content))
		for i, item := range catData.Content {
			var raw RawComicItem
			json.Unmarshal(item, &raw)
			items[i] = rawToComicItem(raw)
		}
		return &SearchResult{Items: items, TotalPages: (catData.Total + 19) / 20, Page: page}, nil
	}

	return c.parseAnyList(body, page)
}

func (c *Client) GetComments(comicID string, page int) ([]CommentItem, error) {
	body, err := c.get("/api/comic/"+comicID+"/comments", map[string]string{
		"offset": fmt.Sprintf("%d", (page-1)*20), "limit": "20",
	})
	if err != nil {
		return nil, err
	}
	// 尝试直接从 /forum 获取
	body2, err2 := c.get("/forum", map[string]string{"mode": "undefined", "aid": comicID, "page": fmt.Sprintf("%d", page)})
	if err2 == nil {
		body = body2
	}

	var resp struct {
		Data json.RawMessage `json:"data"`
	}
	if json.Unmarshal(body, &resp) != nil {
		return nil, err
	}
	var comments []CommentData
	if json.Unmarshal(resp.Data, &comments) != nil {
		return nil, err
	}
	items := make([]CommentItem, len(comments))
	for i, c := range comments {
		items[i] = CommentItem{
			ID: c.ID, Content: c.Content, Author: c.Author,
			Avatar: c.Avatar, LikeCount: c.LikeCount, CreateTime: c.CreateTime,
		}
	}
	return items, nil
}

func (c *Client) Login(username, password string) (*LoginUserData, error) {
	body, err := c.doPostForm("/login", url.Values{"username": {username}, "password": {password}})
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
		return nil, fmt.Errorf("parse user data failed")
	}
	if userData.JWTToken == "" {
		return nil, fmt.Errorf("no token in response")
	}
	c.auth = "Bearer " + userData.JWTToken
	return &userData, nil
}

func (c *Client) GetOnlineFavorites(page int) (*SearchResult, error) {
	body, err := c.get("/favorite", map[string]string{
		"page": fmt.Sprintf("%d", page), "folder_id": "0", "o": "mr",
	})
	if err != nil {
		return nil, err
	}
	return c.parseAnyList(body, page)
}

func (c *Client) GetOnlineHistory(page int) (*SearchResult, error) {
	body, err := c.get("/watch_list", map[string]string{"page": fmt.Sprintf("%d", page)})
	if err != nil {
		return nil, err
	}
	return c.parseAnyList(body, page)
}

// ==================== Stubs ====================

func (c *Client) GetIndexInfo(page int) (map[string][]ComicItem, error) {
	body, err := c.get("/promote", map[string]string{"page": fmt.Sprintf("%d", page)})
	if err != nil {
		return nil, err
	}

	var resp struct {
		Code int             `json:"code"`
		Data json.RawMessage `json:"data"`
	}
	if json.Unmarshal(body, &resp) != nil || resp.Code != 200 {
		return nil, fmt.Errorf("promote failed")
	}

	// promote返回数组: [{title: "分类名", content: [{id,name,author,image}, ...]}]
	var items []struct {
		Title   string `json:"title"`
		Content []struct {
			ID     json.Number    `json:"id"`
			Name   string         `json:"name"`
			Author json.RawMessage `json:"author"`
			Image  string         `json:"image"`
		} `json:"content"`
	}
	if err := json.Unmarshal(resp.Data, &items); err != nil {
		return nil, err
	}

	result := make(map[string][]ComicItem)
	for _, item := range items {
		if len(item.Content) == 0 {
			continue
		}
		comics := make([]ComicItem, len(item.Content))
		for i, c := range item.Content {
			comics[i] = ComicItem{
				ID:       c.ID.String(),
				Title:    c.Name,
				Author:   parseAuthor(c.Author),
				CoverURL: buildCoverURL(c.ID.String()),
			}
		}
		result[item.Title] = comics
	}
	return result, nil
}

func (c *Client) GetRanking(rankType string, page int) (*SearchResult, error) {
	return c.GetCategoryFilter("0", mapSort(rankType), page)
}

func mapSort(t string) string {
	m := map[string]string{"daily": "mv_t", "weekly": "mv_w", "monthly": "mv_m"}
	if s, ok := m[t]; ok {
		return s
	}
	return "mr"
}

func (c *Client) GetSubComments(commentID string, page int) ([]CommentItem, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *Client) Register(username, email, password, passwordConfirm, gender string) error {
	_, err := c.doPostForm("/signup", url.Values{
		"username": {username}, "email": {email}, "password": {password},
		"password_confirm": {passwordConfirm}, "gender": {gender},
		"age": {"on"}, "terms": {"on"}, "submit_signup": {""},
	})
	return err
}

func (c *Client) GetUserInfo() (*LoginUserData, error) {
	return nil, fmt.Errorf("not logged in")
}

func (c *Client) Sign() error {
	return nil
}

// ==================== Helper Functions ====================

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

func rawToComicItem(raw RawComicItem) ComicItem {
	author := parseAuthor(raw.Author)
	id := raw.PathWord
	if id == "" {
		id = raw.ID.String()
	}
	cover := raw.Cover
	if cover == "" {
		cover = raw.Image
	}
	if cover == "" {
		cover = buildCoverURL(id)
	}
	return ComicItem{ID: id, Title: raw.Name, Author: author, CoverURL: cover}
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

func (c *Client) parseAnyList(body []byte, page int) (*SearchResult, error) {
	// Try standard format
	resp, err := c.parseComicList(body, page)
	if err == nil {
		return resp, nil
	}
	// Try array format (data as array)
	var arrResp struct {
		Code int               `json:"code"`
		Data []json.RawMessage `json:"data"`
	}
	if json.Unmarshal(body, &arrResp) == nil && len(arrResp.Data) > 0 {
		return parseRawItems(arrResp.Data, page), nil
	}
	// Try bare array in body
	var bareArr []json.RawMessage
	if json.Unmarshal(body, &bareArr) == nil && len(bareArr) > 0 {
		return parseRawItems(bareArr, page), nil
	}
	return resp, err
}

func parseRawItems(rawItems []json.RawMessage, page int) *SearchResult {
	items := make([]ComicItem, 0)
	for _, item := range rawItems {
		var raw RawComicItem
		if json.Unmarshal(item, &raw) != nil {
			continue
		}
		items = append(items, rawToComicItem(raw))
	}
	return &SearchResult{Items: items, TotalPages: 1, Page: page}
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
		results[i] = rawToComicItem(item)
	}
	totalPages := 1
	if searchData.Limit > 0 {
		totalPages = (searchData.Total + searchData.Limit - 1) / searchData.Limit
	}
	return &SearchResult{Items: results, TotalPages: totalPages, Page: page}, nil
}
