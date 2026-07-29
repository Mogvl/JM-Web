package api

import (
	"encoding/json"
	"os/exec"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

func (r *Router) processDownload(downloadID int, comicID string, format string) {
	r.db.UpdateDownload(downloadID, "downloading", 0, 0)

	comic, err := r.client.GetComicDetail(comicID)
	if err != nil {
		r.db.UpdateDownload(downloadID, "failed", 0, 0)
		return
	}

	safeTitle := sanitizeFilename(comic.Title)
	if safeTitle == "" {
		safeTitle = comicID
	}
	comicDir := filepath.Join(r.config.Download.Dir, safeTitle+"_"+comicID)

	targetFmt := format
	if targetFmt == "" {
		targetFmt = "jpg"
	}
	if targetFmt == "jpg" {
		targetFmt = "jpg"
	}

	// 调用官方 Python jmcomic 库下载（自动处理图片重组/格式转换）
	scriptPath := filepath.Join("..", "download_jm.py")
	cmd := exec.Command("python3", scriptPath, comicID, comicDir, targetFmt)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorf("Python download failed: %v, output: %s", err, string(out))
		r.db.UpdateDownload(downloadID, "failed", 0, 0)
		return
	}

	var result struct {
		Success bool   `json:"success"`
		Count   int    `json:"count"`
		Dir     string `json:"dir"`
		Error   string `json:"error"`
	}
	if err := json.Unmarshal(out, &result); err != nil {
		log.Errorf("Parse python output failed: %v, output: %s", err, string(out))
		r.db.UpdateDownload(downloadID, "failed", 0, 0)
		return
	}

	if !result.Success {
		log.Errorf("Download failed: %s", result.Error)
		r.db.UpdateDownload(downloadID, "failed", 0, 0)
		return
	}

	r.db.SetDownloadTotal(downloadID, result.Count)
	r.db.UpdateDownload(downloadID, "completed", 100, result.Count)
	r.db.SetDownloadPath(downloadID, result.Dir)
	log.Infof("Download completed: %s -> %s (%d images)", comic.Title, result.Dir, result.Count)
}

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