<template>
  <div class="reader">
    <div class="toolbar">
      <el-button @click="$router.back()">返回</el-button>
      <span class="title">{{ currentChapterTitle }}</span>
      <el-button @click="prevChapter" :disabled="!hasPrev">上一章</el-button>
      <el-button @click="nextChapter" :disabled="!hasNext">下一章</el-button>
    </div>
    <div v-if="loading" class="loading">加载中...</div>
    <div v-else class="images" :class="readMode">
      <img v-for="(img, i) in images" :key="i" :src="img" loading="lazy" />
    </div>
    <div class="toolbar bottom">
      <el-button @click="prevChapter" :disabled="!hasPrev">上一章</el-button>
      <el-button @click="nextChapter" :disabled="!hasNext">下一章</el-button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getImages, getChapters } from '../api'

const route = useRoute()
const router = useRouter()
const images = ref([])
const chapters = ref([])
const loading = ref(false)
const readMode = ref(localStorage.getItem('readMode') || 'scroll')

const comicId = computed(() => route.params.comicId)
const chapterId = computed(() => route.params.chapterId)
const currentIndex = computed(() => chapters.value.findIndex(c => c.id === chapterId.value))
const hasPrev = computed(() => currentIndex.value > 0)
const hasNext = computed(() => currentIndex.value < chapters.value.length - 1)
const currentChapterTitle = computed(() => {
  const ch = chapters.value.find(c => c.id === chapterId.value)
  return ch ? ch.title : ''
})

const loadImages = async () => {
  loading.value = true
  try {
    const data = await getImages(chapterId.value)
    images.value = data.images || []
  } finally {
    loading.value = false
  }
}

const prevChapter = () => {
  if (hasPrev.value) router.push(`/read/${comicId.value}/${chapters.value[currentIndex.value - 1].id}`)
}
const nextChapter = () => {
  if (hasNext.value) router.push(`/read/${comicId.value}/${chapters.value[currentIndex.value + 1].id}`)
}

onMounted(async () => {
  const data = await getChapters(comicId.value)
  chapters.value = data.chapters || []
  await loadImages()
})

watch(() => route.params.chapterId, loadImages)
</script>

<style scoped>
.toolbar { display: flex; align-items: center; justify-content: center; gap: 20px; padding: 15px 20px; background: #fff; margin-bottom: 20px; border-radius: 8px; }
.toolbar .title { font-size: 16px; font-weight: bold; margin: 0 20px; }
.toolbar.bottom { margin-top: 20px; margin-bottom: 0; }
.images { display: flex; flex-direction: column; align-items: center; gap: 10px; }
.images img { max-width: 100%; }
.images.page img { max-width: 800px; }
.loading { text-align: center; padding: 100px 0; color: #999; }
</style>
