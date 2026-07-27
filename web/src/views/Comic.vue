<template>
  <div class="comic" v-if="comic">
    <div class="header">
      <img :src="comic.cover_url" class="cover" />
      <div class="info">
        <h1>{{ comic.title }}</h1>
        <p v-if="comic.author">作者: {{ comic.author }}</p>
        <p v-if="comic.tags && comic.tags.length">标签: {{ comic.tags.join(', ') }}</p>
        <p v-if="comic.description">{{ comic.description }}</p>
        <div class="actions">
          <el-button type="primary" @click="toggleFavorite">{{ isFav ? '取消收藏' : '收藏' }}</el-button>
          <el-button @click="startDownload">下载</el-button>
        </div>
      </div>
    </div>
    <h2>章节列表</h2>
    <div class="chapters">
      <el-button v-for="ch in chapters" :key="ch.id" @click="$router.push(`/read/${$route.params.id}/${ch.id}`)">{{ ch.title }}</el-button>
    </div>
  </div>
  <div v-else class="loading">加载中...</div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getComic, getChapters, getFavorites, addFavorite, removeFavorite, createDownload } from '../api'
import { ElMessage } from 'element-plus'

const route = useRoute()
const router = useRouter()
const comic = ref(null)
const chapters = ref([])
const isFav = ref(false)

onMounted(async () => {
  const id = route.params.id
  comic.value = await getComic(id)
  const data = await getChapters(id)
  chapters.value = data.chapters || []
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
.info p { margin-bottom: 10px; color: #666; }
.actions { margin-top: 20px; display: flex; gap: 12px; }
.chapters { display: flex; flex-wrap: wrap; gap: 10px; }
.loading { text-align: center; padding: 100px 0; color: #999; }
</style>
