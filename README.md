# JMComic Web

用 Go + Vue3 重构的 JMComic 漫画阅读 Web 端。

## 功能

- 🔍 搜索漫画
- 📖 在线阅读
- ⬇️ 下载管理
- ⭐ 收藏夹
- 📜 阅读历史
- 📂 分类浏览

## 技术栈

- **后端**: Go + Gin + GORM + SQLite
- **前端**: Vue3 + Vite + Element Plus
- **部署**: Docker + GitHub Actions CI

## 绿联云部署

```yaml
services:
  jmcomic-web:
    image: ghcr.io/mogvl/jm-web/jm-web-frontend:latest
    container_name: jmcomic-web
    ports:
      - "5000:80"
    environment:
      - TZ=Asia/Shanghai
    restart: unless-stopped

  jmcomic-server:
    image: ghcr.io/mogvl/jm-web/jm-web-server:latest
    container_name: jmcomic-server
    environment:
      - GIN_MODE=release
      - DB_PATH=/data/jmcomic.db
      - DOWNLOAD_DIR=/data/downloads
      - JM_BASE_URL=https://jmcomic.me
      - TZ=Asia/Shanghai
    volumes:
      - /volume1/JM:/data/downloads
      - /volume1/docker/JM:/data
    restart: unless-stopped
```

访问 `http://IP:5000`

## 本地开发

```bash
# 后端
cd server
GOPROXY=https://goproxy.cn,direct go mod tidy
go run cmd/main.go

# 前端
cd web
npm install
npm run dev
```

## API

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/search?q=&page= | 搜索 |
| GET | /api/comic/:id | 漫画详情 |
| GET | /api/comic/:id/chapters | 章节列表 |
| GET | /api/chapter/:id | 章节图片 |
| GET | /api/categories | 分类 |
| GET | /api/favorites | 收藏 |
| POST | /api/favorites | 添加收藏 |
| DELETE | /api/favorites/:id | 取消收藏 |
| GET | /api/history | 历史 |
| DELETE | /api/history/:id | 删除历史 |
| POST | /api/download | 创建下载 |
| GET | /api/downloads | 下载列表 |

## 环境变量

| 变量 | 默认值 | 说明 |
|------|--------|------|
| PORT | 8080 | 后端端口 |
| DB_PATH | /data/jmcomic.db | 数据库路径 |
| DOWNLOAD_DIR | /data/downloads | 下载目录 |
| JM_BASE_URL | https://jmcomic.me | API地址 |
