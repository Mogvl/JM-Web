# JMComic Web

基于 Go + Vue3 重构的 JMComic 漫画阅读 Web 端，完整复刻 [JMComic-qt](https://github.com/tonquer/JMComic-qt) 所有功能。

## 功能

### 用户
- 👤 登录 / 注册 / 游客访问
- ⭐ 收藏（支持文件夹分类）
- 📂 本地收藏
- 📜 阅读历史
- ☁️ 远程历史
- 💬 我的评论
- 📝 全部评论
- ✅ 每日签到

### 浏览
- 🏠 首页（最新漫画列表，支持分页）
- 🔍 搜索（关键词搜索）
- 🏷️ 分类浏览
- 🏆 排行榜（日/周/月榜）
- 📅 每周更新

### 漫画
- 📖 漫画详情（封面、作者、标签、章节、评论）
- 📄 阅读器（滚动/翻页模式）
- ⬇️ 下载管理
- 📦 批量下载
- 💬 评论查看

### 工具
- 📁 本地阅读（支持拖拽导入）
- 🖥️ NAS 同步
- ✨ Waifu2x 图片增强
- 🔬 批量超分辨率
- ⚙️ 设置（通用/代理/Waifu2x/下载与缓存）
- ❓ 帮助（FAQ）

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
    ports: ["5000:80"]
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
      - JM_BASE_URL=https://www.cdnhjk.net
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
cd server && GOPROXY=https://goproxy.cn,direct go mod tidy && go run cmd/main.go

# 前端
cd web && npm install && npm run dev
```

访问 http://localhost:5001

## 环境变量

| 变量 | 默认值 | 说明 |
|------|--------|------|
| PORT | 5000 | 后端端口 |
| DB_PATH | /data/jmcomic.db | 数据库路径 |
| DOWNLOAD_DIR | /data/downloads | 下载目录 |
| JM_BASE_URL | https://www.cdnhjk.net | API地址 |
