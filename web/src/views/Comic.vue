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
        <span class="continue-reading" v-if="lastProgress.chapter_id">
          <el-button size="small" type="primary" plain @click="continueRead">继续上次阅读 →</el-button>
        </span>
      </div>
      <div v-if="chapters.length === 0" class="empty-state">暂无章节</div>
      <div v-else class="chapter-flow">
        <button v-for="ch in chapters" :key="ch.id" class="chapter-item" :class="{ 'current': ch.id === lastProgress.chapter_id }" @click="$router.push(`/read/${$route.params.id}/${ch.id}`)">
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

      <div v-if="canComment" class="comment-box">
        <el-input v-model="newComment" type="textarea" :rows="3" placeholder="发表你的评论..." resize="none" />
        <div class="comment-actions">
          <el-button type="primary" size="small" :loading="posting" @click="submitComment">发表评论</el-button>
        </div>
      </div>

      <div v-if="comments.length === 0" class="empty-state">暂无评论</div>
      <div v-else class="comment-list">
        <div v-for="comment in comments" :key="comment.id" class="comment-item">
          <el-avatar :size="36" :src="comment.avatar">{{ (comment.author || '?').charAt(0) }}</el-avatar>
          <div class="comment-body">
            <div class="comment-head">
              <span class="comment-author">{{ comment.author }}</span>
              <span class="comment-time">{{ comment.create_time }}</span>
              <span class="comment-like" v-if="comment.like_count"><el-icon><Pointer /></el-icon>{{ comment.like_count }}</span>
            </div>
            <div class="comment-text">{{ comment.content }}</div>
            <div class="comment-foot">
              <span class="reply-link" @click="toggleReply(comment)">{{ replyingTo === comment.id ? '取消回复' : '回复' }}</span>
              <span class="sub-link" v-if="comment.reply_count" @click="toggleSub(comment)">{{ subOpen === comment.id ? '收起' : `查看 ${comment.reply_count} 条回复` }}</span>
            </div>
            <div v-if="replyingTo === comment.id" class="reply-box">
              <el-input v-model="replyText" size="small" placeholder="回复..." @keyup.enter="submitReply(comment)" />
              <el-button type="primary" size="small" :loading="posting" @click="submitReply(comment)">发送</el-button>
            </div>
            <div v-if="subOpen === comment.id" class="sub-list">
              <div v-for="sub in subComments(comment.id)" :key="sub.id" class="sub-item">
                <el-avatar :size="26" :src="sub.avatar">{{ (sub.author || '?').charAt(0) }}</el-avatar>
                <div class="sub-body">
                  <span class="sub-author">{{ sub.author }}</span>
                  <span class="sub-text">{{ sub.content }}</span>
                </div>
              </div>
              <div v-if="!subLoaded(comment.id) && comment.reply_count" class="sub-load">
                <el-button size="small" text @click="loadSubs(comment)">加载回复</el-button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
  <div v-else class="loading-state">加载中...</div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getComic, getChapters, getComments, getFavorites, addFavorite, removeFavorite, createDownload, getSubComments, postComment, replyComment, getProgress } from '../api'
import { Star, Download, ChatDotRound, Pointer } from '@element-plus/icons-vue'
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
const newComment = ref('')
const replyText = ref('')
const posting = ref(false)
const replyingTo = ref('')
const subOpen = ref('')
const subMap = ref({})
const lastProgress = ref({})

const canComment = computed(() => localStorage.getItem('token') && localStorage.getItem('token') !== 'guest')

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
  try { lastProgress.value = await getProgress(id) || {} } catch (e) {}
})

const continueRead = () => {
  if (lastProgress.value.chapter_id) router.push(`/read/${route.params.id}/${lastProgress.value.chapter_id}`)
}

const toggleFavorite = async () => {
  const id = route.params.id
  if (isFav.value) { await removeFavorite(id); isFav.value = false; ElMessage.success('已取消收藏') }
  else { await addFavorite(id); isFav.value = true; ElMessage.success('已收藏') }
}

const startDownload = () => { showDownloadDialog.value = true }

const searchTag = (tag) => router.push({ path: '/search', query: { q: tag } })
const searchAuthor = (author) => { if (author) router.push({ path: '/search', query: { q: author } }) }
const scrollToComments = () => { document.getElementById('comments')?.scrollIntoView({ behavior: 'smooth' }) }

// ===== 评论交互 =====
const submitComment = async () => {
  if (!newComment.value.trim()) return
  posting.value = true
  try {
    await postComment(route.params.id, newComment.value.trim())
    ElMessage.success('评论已发表')
    newComment.value = ''
    comments.value = await getComments(route.params.id) || []
  } catch (e) { ElMessage.error('发表失败') }
  finally { posting.value = false }
}

const toggleReply = (comment) => { replyingTo.value = replyingTo.value === comment.id ? '' : comment.id; replyText.value = '' }
const submitReply = async (comment) => {
  if (!replyText.value.trim()) return
  posting.value = true
  try {
    await replyComment(route.params.id, comment.id, replyText.value.trim())
    ElMessage.success('回复成功')
    replyText.value = ''
    replyingTo.value = ''
  } catch (e) { ElMessage.error('回复失败') }
  finally { posting.value = false }
}

const toggleSub = (comment) => {
  if (subOpen.value === comment.id) { subOpen.value = ''; return }
  subOpen.value = comment.id
  if (!subMap.value[comment.id]) loadSubs(comment)
}
const subComments = (id) => subMap.value[id] || []
const subLoaded = (id) => !!subMap.value[id]
const loadSubs = async (comment) => {
  try { subMap.value[comment.id] = await getSubComments(comment.id) || [] }
  catch (e) { subMap.value[comment.id] = [] }
}
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
.tag { padding: 5px 14px; background: var(--bg-soft); color: var(--accent); border: 1px solid var(--border-light); border-radius: 999px; font-size: 12px; cursor: pointer; transition: var(--transition); }
.tag:hover { background: var(--accent-grad); color: #fff; border-color: transparent; }

.section { margin-bottom: 32px; }
.section-head { display: flex; align-items: center; gap: 8px; margin-bottom: 16px; }
.section-title { font-size: 18px; font-weight: 600; color: var(--text-primary); }
.count { font-size: 13px; color: var(--text-muted); }

.chapter-flow { display: flex; flex-wrap: wrap; gap: 8px; }
.chapter-item { padding: 8px 16px; border: 1px solid var(--border); background: var(--bg-surface); color: var(--text-secondary); border-radius: var(--radius-sm); font-size: 13px; cursor: pointer; transition: var(--transition); }
.chapter-item:hover { border-color: var(--accent); color: var(--accent); }
.chapter-item.current { border-color: var(--accent); color: var(--accent); background: var(--accent-soft); }
.continue-reading { margin-left: auto; }

.comment-list { display: flex; flex-direction: column; gap: 12px; }
.comment-item { display: flex; gap: 12px; padding: 14px; background: var(--bg-surface); border: 1px solid var(--border); border-radius: var(--radius-md); }
.comment-body { flex: 1; min-width: 0; }
.comment-head { display: flex; align-items: center; justify-content: space-between; margin-bottom: 6px; }
.comment-author { font-weight: 600; color: var(--text-primary); font-size: 13px; }
.comment-time { color: var(--text-muted); font-size: 11px; }
.comment-like { display: inline-flex; align-items: center; gap: 3px; color: var(--text-muted); font-size: 12px; }
.comment-text { font-size: 13px; color: var(--text-secondary); line-height: 1.6; word-break: break-word; }
.comment-foot { display: flex; gap: 16px; margin-top: 8px; }
.reply-link, .sub-link { font-size: 12px; color: var(--accent); cursor: pointer; }
.comment-box { background: var(--bg-surface); border: 1px solid var(--border); border-radius: var(--radius-md); padding: 12px; margin-bottom: 16px; }
.comment-actions { margin-top: 8px; text-align: right; }
.reply-box { display: flex; gap: 8px; margin-top: 8px; }
.sub-list { margin-top: 10px; padding-top: 8px; border-top: 1px dashed var(--border); display: flex; flex-direction: column; gap: 8px; }
.sub-item { display: flex; gap: 8px; align-items: flex-start; }
.sub-body { display: flex; flex-direction: column; font-size: 13px; min-width: 0; }
.sub-author { color: var(--accent); font-weight: 500; font-size: 12px; }
.sub-text { color: var(--text-secondary); line-height: 1.5; word-break: break-word; }
.sub-load { text-align: center; }

@media (max-width: 768px) {
  .detail-header { flex-direction: column; gap: 16px; align-items: center; }
  .cover-wrap { text-align: center; }
  .cover-actions { flex-direction: row; width: 100%; }
  .tool-btn { flex: 1; }
  .info { width: 100%; }
  .label { width: 52px; }
}
</style>
