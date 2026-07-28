# JMComic Web

基于 Go + Vue3 重构的 JMComic 漫画阅读 Web 端，完整复刻 [JMComic-qt](https://github.com/tonquer/JMComic-qt) 所有功能。

## 功能

### 用户
- 👤 登录 / 注册 / 游客访问（登录返回 J币、等级、头像）
- ⭐ 收藏（在线同步 + 本地存储）
- 📂 本地收藏
- 📜 阅读历史（在线同步）
- ☁️ 远程历史
- 💬 我的评论 / 全部评论
- ✅ 每日签到

### 浏览
- 🏠 首页（最新上传 + 推荐分类标签页）
- 🔍 搜索（支持 `+包含` / `-排除` 语法，排序筛选）
- 🏷️ 分类/排行（标题块 + 分类标签 + 子分类筛选 + 7种排序）
- 📅 每周必看（周选择下拉 + 日漫/韩漫/其他标签页）

### 漫画
- 📖 漫画详情（封面、作者、标签、观看数、爱心数、章节列表、评论）
- 📄 阅读器（上下章切换、图片流）
- ⬇️ 下载管理 / 批量下载

### 工具
- 📁 本地阅读
- 🖥️ NAS 同步
- ✨ Waifu2x 图片增强
- 🔬 批量超分辨率
- ⚙️ 设置
- ❓ 帮助

## 技术栈

- **后端**: Go + Gin + GORM + SQLite
- **前端**: Vue3 + Vite + Element Plus
- **部署**: Docker + GitHub Actions CI

## 设计特色

- **浅色主题** - 清新明亮，适合白天阅读
- **珊瑚红强调色** - `#FF4D6D`
- **移动端适配** - 侧边栏抽屉、自适应网格、触摸友好
- **统一设计系统** - CSS 变量管理颜色/圆角/阴影

## 端口配置

| 服务 | 端口 |
|------|------|
| 前端开发 | 5000 |
| 后端 API | 5001 |
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

## API 接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/search | 搜索 |
| GET | /api/latest | 最新上传 |
| GET | /api/index | 首页推荐 |
| GET | /api/comic/:id | 漫画详情 |
| GET | /api/comic/:id/chapters | 章节列表 |
| GET | /api/chapter/:id | 章节图片 |
| GET | /api/categories | 分类列表 |
| GET | /api/categories/filter | 分类筛选 |
| GET | /api/week | 每周必看列表 |
| GET | /api/week/filter | 每周必看漫画 |
| GET | /api/favorites | 收藏 |
| GET | /api/history | 历史 |
| POST | /api/login | 登录 |
| POST | /api/download | 创建下载 |

## 环境变量

| 变量 | 默认值 | 说明 |
|------|--------|------|
| PORT | 5000 | 后端端口 |
| DB_PATH | /data/jmcomic.db | 数据库路径 |
| DOWNLOAD_DIR | /data/downloads | 下载目录 |
| JM_BASE_URL | https://www.cdnhjk.net | API地址 |
