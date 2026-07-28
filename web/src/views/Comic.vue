<template>
  <div class="comic-page" v-if="comic">
    <div class="detail-header">
      <div class="cover-wrap">
        <img :src="comic.cover_url" class="cover" />
      </div>
      <div class="info">
        <h1 class="title">{{ comic.title }}</h1>
        <div class="meta-row">
          <span class="meta-item" v-if="comic.author"><el-icon><User /></el-icon> {{ comic.author }}</span>
          <span class="meta-item" v-if="comic.total_views"><el-icon><View /></el-icon> {{ comic.total_views }}</span>
          <span class="meta-item" v-if="comic.likes"><el-icon><Star /></el-icon> {{ comic.likes }}</span>
          <span class="meta-item" v-if="comic.total_photos"><el-icon><Picture /></el-icon> {{ comic.total_photos }}页</span>
        </div>
        <div class="tags" v-if="comic.tags && comic.tags.length">
          <el-tag v-for="tag in comic.tags" :key="tag" size="small" effect="plain" round>{{ tag }}</el-tag>
        </div>
        <p class="desc" v-if="comic.description">{{ comic.description }}</p>
        <div class="actions">
          <el-button type="primary" round @click="toggleFavorite">{{ isFav ? '已收藏' : '收藏' }}</el-button>
          <el-button round @click="startDownload"><el-icon><Download /></el-icon> 下载</el-button>
          <el-button v-if="chapters.length" round type="success" @click="readFirst">开始阅读</el-button>
        </div>
      </div>
    </div>

    <div class="section" v-if="chapters.length">
      <h2 class="section-title">章节列表 <span class="count">{{ chapters.length }}</span></h2>
      <div class="chapter-list">
        <button v-for="ch in chapters" :key="ch.id" class="chapter-btn" @click="$router.push(`/read/${$route.params.id}/${ch.id}`)">
          {{ ch.title || `第${ch.sort_order}章` }}
        </button>
      </div>
    </div>

    <div class="section">
      <h2 class="section-title">评论 <span class="count">{{ comments.length }}</span></h2>
      <div v-if="comments.length === 0" class="empty-state">暂无评论</div>
      <div v-else class="comment-list">
        <div v-for="comment in comments" :key="comment.id" class="comment-item">
          <el-avatar :size="36" :src="comment.avatar">{{ (comment.author || '?').charAt(0) }}</el-avatar>
          <div class="comment-body">
            <div class="comment-head"><span class="comment-author">{{ comment.author }}</span></div>
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
import { User, View, Star, Picture, Download } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const route = useRoute()
const router = useRouter()
const comic = ref(null)
const chapters = ref([])
const comments = ref([])
const isFav = ref(false)

onMounted(async () => {
  const id = route.params.id
  comic.value = await getComic(id)
  const data = await getChapters(id)
  chapters.value = data.chapters || []
  try { comments.value = await getComments(id) || [] } catch (e) { comments.value = [] }
  try { const favs = await getFavorites(); isFav.value = favs.some(f => f.comic_id === id || f.id === id) } catch (e) {}
})

const toggleFavorite = async () => {
  const id = route.params.id
  if (isFav.value) { await removeFavorite(id); isFav.value = false; ElMessage.success('已取消收藏') }
  else { await addFavorite(id); isFav.value = true; ElMessage.success('已收藏') }
}

const startDownload = async () => { await createDownload(route.params.id); ElMessage.success('下载任务已创建'); router.push('/downloads') }
const readFirst = () => { if (chapters.value.length) router.push(`/read/${route.params.id}/${chapters.value[0].id}`) }
</script>

<style scoped>
.detail-header { display: flex; gap: 32px; margin-bottom: 40px; }
.cover-wrap { flex-shrink: 0; }
.cover { width: 240px; height: 320px; object-fit: cover; border-radius: var(--radius-lg); box-shadow: var(--shadow-hover); }
.info { flex: 1; min-width: 0; }
.title { font-size: 26px; font-weight: 700; color: var(--text-primary); margin-bottom: 16px; line-height: 1.3; }
.meta-row { display: flex; flex-wrap: wrap; gap: 20px; margin-bottom: 16px; }
.meta-item { display: inline-flex; align-items: center; gap: 6px; font-size: 14px; color: var(--text-secondary); }
.tags { display: flex; flex-wrap: wrap; gap: 8px; margin-bottom: 16px; }
.desc { font-size: 14px; color: var(--text-secondary); line-height: 1.7; margin-bottom: 24px; max-height: 120px; overflow-y: auto; }
.actions { display: flex; gap: 12px; flex-wrap: wrap; }

@media (max-width: 768px) {
  .detail-header { flex-direction: column; gap: 16px; align-items: center; text-align: center; }
  .cover { width: 180px; height: 240px; }
  .info { width: 100%; text-align: left; }
  .title { font-size: 20px; }
  .meta-row { justify-content: center; }
  .tags { justify-content: center; }
  .actions { justify-content: center; }
}
.section { margin-bottom: 32px; }
.section-title { font-size: 18px; font-weight: 600; color: var(--text-primary); margin-bottom: 16px; display: flex; align-items: center; gap: 8px; }
.count { font-size: 13px; color: var(--text-muted); font-weight: 400; }
.chapter-list { display: flex; flex-wrap: wrap; gap: 10px; }
.chapter-btn { padding: 8px 18px; border-radius: 20px; border: 1px solid var(--border); background: var(--bg-surface); color: var(--text-secondary); font-size: 13px; cursor: pointer; transition: var(--transition); }
.chapter-btn:hover { border-color: var(--accent); color: var(--accent); }
.comment-list { display: flex; flex-direction: column; gap: 16px; }
.comment-item { display: flex; gap: 12px; padding: 16px; background: var(--bg-surface); border-radius: var(--radius-md); border: 1px solid var(--border); }
.comment-body { flex: 1; }
.comment-author { font-weight: 600; color: var(--text-primary); font-size: 14px; }
.comment-text { margin-top: 6px; font-size: 14px; color: var(--text-secondary); line-height: 1.6; }
</style>
