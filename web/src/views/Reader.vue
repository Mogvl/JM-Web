<template>
  <div class="reader">
    <div class="toolbar">
      <el-button @click="$router.back()">返回</el-button>
      <el-button @click="prevChapter" :disabled="!hasPrev">上一章</el-button>
      <el-button @click="nextChapter" :disabled="!hasNext">下一章</el-button>
    </div>
    <div v-if="loading" class="loading">加载中...</div>
    <div v-else class="images">
      <img v-for="(img, i) in images" :key="i" :src="img" loading="lazy" />
    </div>
    <div class="toolbar">
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

const comicId = computed(() => route.params.comicId)
const chapterId = computed(() => route.params.chapterId)
const currentIndex = computed(() => chapters.value.findIndex(c => c.id === chapterId.value))
const hasPrev = computed(() => currentIndex.value > 0)
const hasNext = computed(() => currentIndex.value < chapters.value.length - 1)

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
.toolbar { display: flex; justify-content: center; gap: 20px; padding: 20px; background: #fff; margin-bottom: 20px; border-radius: 8px; }
.images { display: flex; flex-direction: column; align-items: center; gap: 10px; }
.images img { max-width: 100%; }
.loading { text-align: center; padding: 100px 0; color: #999; }
</style>
