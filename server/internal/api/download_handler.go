package api

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

func min(a, b int) int { if a < b { return a }; return b }

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
	scriptPath, _ := filepath.Abs("download_jm.py")
	absDir, _ := filepath.Abs(comicDir)
	os.MkdirAll(absDir, 0755)
	cmd := exec.Command("python", scriptPath, comicID, absDir, targetFmt)
	cmd.Dir = absDir // 设置工作目录为下载目录
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
	// 从输出中提取最后一行 JSON
	outputStr := string(out)
	// 找到第一个 { 开始的 JSON
	jsonStart := strings.LastIndex(outputStr, "{")
	jsonEnd := strings.LastIndex(outputStr, "}")
	if jsonStart >= 0 && jsonEnd > jsonStart {
		jsonStr := outputStr[jsonStart : jsonEnd+1]
		if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
			log.Errorf("Parse python output failed: %v, json: %s", err, jsonStr)
			r.db.UpdateDownload(downloadID, "failed", 0, 0)
			return
		}
	} else {
		log.Errorf("No JSON in python output: %s", outputStr[:min(200, len(outputStr))])
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