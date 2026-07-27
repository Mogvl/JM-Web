<template>
  <div class="comic" v-if="comic">
    <div class="header">
      <img :src="comic.cover_url" class="cover" />
      <div class="info">
        <h1>{{ comic.title }}</h1>
        <p v-if="comic.author">作者: {{ comic.author }}</p>
        <p v-if="comic.category">分类: {{ comic.category }}</p>
        <p v-if="comic.status">状态: {{ comic.status }}</p>
        <p v-if="comic.tags && comic.tags.length">标签: {{ comic.tags.join(', ') }}</p>
        <p v-if="comic.description" class="desc">{{ comic.description }}</p>
        <div class="actions">
          <el-button type="primary" @click="toggleFavorite">{{ isFav ? '取消收藏' : '收藏' }}</el-button>
          <el-button @click="startDownload">下载</el-button>
        </div>
      </div>
    </div>

    <h2>章节列表 ({{ chapters.length }})</h2>
    <div class="chapters">
      <el-button v-for="ch in chapters" :key="ch.id" @click="$router.push(`/read/${$route.params.id}/${ch.id}`)">
        {{ ch.title }}
      </el-button>
    </div>

    <h2>评论 ({{ comments.length }})</h2>
    <div class="comments">
      <div v-if="comments.length === 0" class="empty">暂无评论</div>
      <div v-else class="comment-list">
        <div v-for="comment in comments" :key="comment.id" class="comment-item">
          <el-avatar :size="40" :src="comment.avatar" />
          <div class="comment-content">
            <div class="comment-header">
              <span class="author">{{ comment.author }}</span>
              <span class="time">{{ comment.create_time }}</span>
            </div>
            <div class="comment-text">{{ comment.content }}</div>
            <div class="comment-actions">
              <el-button text size="small">👍 {{ comment.like_count }}</el-button>
              <el-button text size="small">回复</el-button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
  <div v-else class="loading">加载中...</div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getComic, getChapters, getComments, getFavorites, addFavorite, removeFavorite, createDownload } from '../api'
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

  try {
    const commentData = await getComments(id)
    comments.value = commentData || []
  } catch (e) {
    comments.value = []
  }

  const favs = await getFavorites()
  isFav.value = favs.some(f => f.comic_id === id)
})

const toggleFavorite = async () => {
  const id = route.params.id
  if (isFav.value) {
    await removeFavorite(id)
    isFav.value = false
    ElMessage.success('已取消收藏')
  } else {
    await addFavorite(id)
    isFav.value = true
    ElMessage.success('已收藏')
  }
}

const startDownload = async () => {
  await createDownload(route.params.id)
  ElMessage.success('下载任务已创建')
  router.push('/downloads')
}
</script>

<style scoped>
.header { display: flex; gap: 30px; margin-bottom: 40px; }
.cover { width: 250px; height: 350px; object-fit: cover; border-radius: 8px; }
.info h1 { font-size: 28px; margin-bottom: 15px; }
.info p { margin-bottom: 8px; color: #666; }
.desc { margin-top: 10px; line-height: 1.6; }
.actions { margin-top: 20px; display: flex; gap: 12px; }
.chapters { display: flex; flex-wrap: wrap; gap: 10px; margin-top: 20px; margin-bottom: 40px; }
.comments { margin-top: 20px; }
.empty { text-align: center; padding: 40px 0; color: #999; }
.comment-list { display: flex; flex-direction: column; gap: 16px; }
.comment-item { display: flex; gap: 12px; }
.comment-content { flex: 1; }
.comment-header { display: flex; justify-content: space-between; margin-bottom: 8px; }
.author { font-weight: bold; }
.time { color: #999; font-size: 12px; }
.comment-text { line-height: 1.6; }
.comment-actions { margin-top: 8px; }
.loading { text-align: center; padding: 100px 0; color: #999; }
</style>
