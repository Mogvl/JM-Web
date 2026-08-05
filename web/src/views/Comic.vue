<template>
  <div class="comic-page" v-if="comic">
    <BackNav :title="comic.title" fallback="/search" />
    <div class="detail-header">
      <div class="cover-wrap">
        <img :src="comic.cover_url" class="cover" />
        <div class="cover-actions">
          <button class="tool-btn" @click="toggleFavorite">
            <el-icon><Star /></el-icon>
            <span>{{ isFav ? '已收藏' : '收藏' }}</span>
          </button>
          <button class="tool-btn" @click="startDownload">
            <el-icon><Download /></el-icon>
            <span>下载</span>
          </button>
          <button class="tool-btn" @click="scrollToComments">
            <el-icon><ChatDotRound /></el-icon>
            <span>评论</span>
          </button>
        </div>
      </div>
      <div class="info">
        <h1 class="title">{{ comic.title }}</h1>
        <div class="info-row">
          <span class="label">id:</span>
          <span class="value">{{ comic.id }}</span>
        </div>
        <div class="info-row">
          <span class="label">作者：</span>
          <span class="author-tag" @click="searchAuthor(comic.author)">{{ comic.author || '未知' }}</span>
        </div>
        <div class="info-row">
          <span class="label">描述：</span>
          <p class="desc">{{ comic.description || '暂无描述' }}</p>
        </div>
        <div class="info-row">
          <span class="label">爱心数：</span>
          <span class="value">{{ comic.likes || 0 }}</span>
        </div>
        <div class="info-row">
          <span class="label">观看数：</span>
          <span class="value">{{ comic.total_views || 0 }}</span>
        </div>
        <div class="info-row">
          <span class="label">页数：</span>
          <span class="value">{{ comic.total_photos || 0 }}</span>
        </div>
        <div class="info-row" v-if="comic.tags && comic.tags.length">
          <span class="label">Tags：</span>
          <div class="tag-list">
            <button v-for="tag in comic.tags" :key="tag" class="tag" @click="searchTag(tag)">{{ tag }}</button>
          </div>
        </div>
      </div>
    </div>

    <div class="section">
      <div class="section-head">
        <h2 class="section-title">章节列表</h2>
        <span class="count">{{ chapters.length }} 话</span>
      </div>
      <div v-if="chapters.length === 0" class="empty-state">暂无章节</div>
      <div v-else class="chapter-flow">
        <button v-for="ch in chapters" :key="ch.id" class="chapter-item" @click="$router.push(`/read/${$route.params.id}/${ch.id}`)">
          {{ ch.title || `第${ch.sort_order}话` }}
        </button>
      </div>
    </div>

    <DownloadEpsDialog v-model="showDownloadDialog" :comic-id="route.params.id" :comic="comic" />

    <div class="section" id="comments">
      <div class="section-head">
        <h2 class="section-title">评论</h2>
        <span class="count">{{ comments.length }} 条</span>
      </div>
      <div v-if="comments.length === 0" class="empty-state">暂无评论</div>
      <div v-else class="comment-list">
        <div v-for="comment in comments" :key="comment.id" class="comment-item">
          <el-avatar :size="36" :src="comment.avatar">{{ (comment.author || '?').charAt(0) }}</el-avatar>
          <div class="comment-body">
            <div class="comment-head">
              <span class="comment-author">{{ comment.author }}</span>
              <span class="comment-time">{{ comment.create_time }}</span>
            </div>
            <div class="comment-text">{{ comment.content }}</div>
          </div>
        </div>
      </div>
    </div>
  </div>
  <div v-else class="loading-state">加载中...</div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getComic, getChapters, getComments, getFavorites, addFavorite, removeFavorite, createDownload } from '../api'
import { Star, Download, ChatDotRound } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import DownloadEpsDialog from '../components/DownloadEpsDialog.vue'
import BackNav from '../components/BackNav.vue'

const route = useRoute()
const router = useRouter()
const comic = ref(null)
const chapters = ref([])
const comments = ref([])
const isFav = ref(false)
const showDownloadDialog = ref(false)

onMounted(async () => {
  const id = route.params.id
  comic.value = await getComic(id)
  const data = await getChapters(id)
  chapters.value = data.chapters || []
  try { comments.value = await getComments(id) || [] } catch (e) { comments.value = [] }
  try {
    const favs = await getFavorites()
    isFav.value = favs.some(f => f.comic_id === id || f.id === id)
  } catch (e) {}
})

const toggleFavorite = async () => {
  const id = route.params.id
  if (isFav.value) { await removeFavorite(id); isFav.value = false; ElMessage.success('已取消收藏') }
  else { await addFavorite(id); isFav.value = true; ElMessage.success('已收藏') }
}

const startDownload = () => { showDownloadDialog.value = true }

const searchTag = (tag) => router.push({ path: '/search', query: { q: tag } })
const searchAuthor = (author) => { if (author) router.push({ path: '/search', query: { q: author } }) }
const scrollToComments = () => { document.getElementById('comments')?.scrollIntoView({ behavior: 'smooth' }) }
</script>

<style scoped>
.detail-header { display: flex; gap: 32px; margin-bottom: 40px; }
.cover-wrap { flex-shrink: 0; }
.cover { width: 240px; height: 320px; object-fit: cover; border-radius: var(--radius-lg); box-shadow: var(--shadow-hover); }
.cover-actions { display: flex; flex-direction: column; gap: 8px; margin-top: 12px; }
.tool-btn { display: flex; align-items: center; justify-content: center; gap: 8px; width: 100%; padding: 10px; border: 1px solid var(--border); background: var(--bg-surface); color: var(--text-secondary); border-radius: var(--radius-md); cursor: pointer; transition: var(--transition); font-size: 13px; }
.tool-btn:hover { border-color: var(--accent); color: var(--accent); background: var(--accent-soft); }
.info { flex: 1; min-width: 0; }
.title { font-size: 24px; font-weight: 700; color: var(--text-primary); margin-bottom: 16px; line-height: 1.3; word-wrap: break-word; }
.info-row { display: flex; gap: 8px; margin-bottom: 12px; align-items: flex-start; }
.label { flex-shrink: 0; width: 64px; color: var(--text-muted); font-size: 13px; text-align: right; }
.value { color: var(--text-secondary); font-size: 13px; }
.author-tag { color: var(--accent); cursor: pointer; font-size: 13px; }
.author-tag:hover { text-decoration: underline; }
.desc { flex: 1; color: var(--text-secondary); font-size: 13px; line-height: 1.6; max-height: 100px; overflow-y: auto; }
.tag-list { display: flex; flex-wrap: wrap; gap: 6px; flex: 1; }
.tag { padding: 3px 12px; background: #FBEEF3; color: #C45F7D; border: 1px solid #C45F7D; border-radius: 12px; font-size: 12px; cursor: pointer; transition: var(--transition); }
.tag:hover { background: #C45F7D; color: #fff; }

.section { margin-bottom: 32px; }
.section-head { display: flex; align-items: center; gap: 8px; margin-bottom: 16px; }
.section-title { font-size: 18px; font-weight: 600; color: var(--text-primary); }
.count { font-size: 13px; color: var(--text-muted); }

.chapter-flow { display: flex; flex-wrap: wrap; gap: 8px; }
.chapter-item { padding: 8px 16px; border: 1px solid var(--border); background: var(--bg-surface); color: var(--text-secondary); border-radius: var(--radius-sm); font-size: 13px; cursor: pointer; transition: var(--transition); }
.chapter-item:hover { border-color: var(--accent); color: var(--accent); }

.comment-list { display: flex; flex-direction: column; gap: 12px; }
.comment-item { display: flex; gap: 12px; padding: 14px; background: var(--bg-surface); border: 1px solid var(--border); border-radius: var(--radius-md); }
.comment-body { flex: 1; min-width: 0; }
.comment-head { display: flex; align-items: center; justify-content: space-between; margin-bottom: 6px; }
.comment-author { font-weight: 600; color: var(--text-primary); font-size: 13px; }
.comment-time { color: var(--text-muted); font-size: 11px; }
.comment-text { font-size: 13px; color: var(--text-secondary); line-height: 1.6; word-break: break-word; }

@media (max-width: 768px) {
  .detail-header { flex-direction: column; gap: 16px; align-items: center; }
  .cover-wrap { text-align: center; }
  .cover-actions { flex-direction: row; width: 100%; }
  .tool-btn { flex: 1; }
  .info { width: 100%; }
  .label { width: 52px; }
}
</style>
