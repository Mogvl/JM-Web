# JMComic Web

用 Go + Vue3 重构的 JMComic 漫画阅读 Web 端，参考 [JMComic-qt](https://github.com/tonquer/JMComic-qt)。

## 功能

| 功能 | 说明 |
|------|------|
| 🔍 搜索 | 关键词搜索漫画 |
| 📂 分类 | 按分类浏览 |
| 🏆 排行榜 | 日榜/周榜/月榜 |
| 📅 每周更新 | 每周更新列表 |
| 📖 漫画详情 | 封面、作者、标签、章节 |
| 📄 阅读器 | 滚动/翻页模式 |
| 💬 评论 | 查看/回复评论 |
| ⭐ 收藏 | 收藏漫画 |
| 📜 历史 | 阅读历史 |
| ⬇️ 下载 | 下载管理 |
| 👤 用户 | 登录/签到 |
| 📁 本地阅读 | 本地文件导入 |
| 🖥️ NAS 同步 | 同步到 NAS |
| ✨ 图片增强 | Waifu2x 超分辨率 |
| ⚙️ 设置 | API/代理/阅读模式 |
| ❓ 帮助 | FAQ |

## 技术栈

- **后端**: Go + Gin + GORM + SQLite
- **前端**: Vue3 + Vite + Element Plus
- **部署**: Docker + GitHub Actions CI

## 端口配置

| 服务 | 端口 |
|------|------|
| 后端 API | 5000 |
| 前端开发 | 5001 |
| Docker 部署 | 5000 |

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
# 后端（端口 5000）
cd server
GOPROXY=https://goproxy.cn,direct go mod tidy
go run cmd/main.go

# 前端（端口 5001）
cd web
npm install
npm run dev
```

访问 http://localhost:5001

## API 接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/search | 搜索 |
| GET | /api/comic/:id | 漫画详情 |
| GET | /api/comic/:id/chapters | 章节列表 |
| GET | /api/chapter/:id | 章节图片 |
| GET | /api/categories | 分类 |
| GET | /api/ranking | 排行榜 |
| GET | /api/comic/:id/comments | 评论 |
| GET | /api/favorites | 收藏 |
| POST | /api/favorites | 添加收藏 |
| DELETE | /api/favorites/:id | 取消收藏 |
| GET | /api/history | 历史 |
| DELETE | /api/history/:id | 删除历史 |
| POST | /api/download | 创建下载 |
| GET | /api/downloads | 下载列表 |
| POST | /api/login | 登录 |
| GET | /api/user/info | 用户信息 |
| POST | /api/user/sign | 签到 |
| GET | /api/help | 帮助 |

## 环境变量

| 变量 | 默认值 | 说明 |
|------|--------|------|
| PORT | 5000 | 后端端口 |
| DB_PATH | /data/jmcomic.db | 数据库路径 |
| DOWNLOAD_DIR | /data/downloads | 下载目录 |
| JM_BASE_URL | https://jmcomic.me | API地址 |
| JM_AUTH | | 认证 Token |
