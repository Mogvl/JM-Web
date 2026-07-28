<template>
  <div class="reader-page">
    <div class="reader-bar">
      <el-button text @click="$router.back()"><el-icon><ArrowLeft /></el-icon> 返回</el-button>
      <span class="chapter-title">{{ currentTitle }}</span>
      <div class="nav-btns">
        <el-button size="small" :disabled="!hasPrev" @click="prevChapter">上一章</el-button>
        <el-button size="small" :disabled="!hasNext" @click="nextChapter">下一章</el-button>
      </div>
    </div>

    <div v-if="loading" class="loading-state">加载中...</div>
    <div v-else-if="images.length === 0" class="empty-state">暂无图片</div>
    <div v-else class="image-stream">
      <img v-for="(img, i) in images" :key="i" :src="img" loading="lazy" />
    </div>

    <div class="reader-bar bottom">
      <el-button size="small" :disabled="!hasPrev" @click="prevChapter">上一章</el-button>
      <el-button size="small" :disabled="!hasNext" @click="nextChapter">下一章</el-button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowLeft } from '@element-plus/icons-vue'
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
const hasNext = computed(() => currentIndex.value >= 0 && currentIndex.value < chapters.value.length - 1)
const currentTitle = computed(() => { const ch = chapters.value.find(c => c.id === chapterId.value); return ch ? (ch.title || `第${ch.sort_order}章`) : '' })

const loadImages = async () => {
  loading.value = true
  try { const data = await getImages(chapterId.value); images.value = data.images || [] }
  catch (e) { images.value = [] }
  finally { loading.value = false }
}

const prevChapter = () => { if (hasPrev.value) router.push(`/read/${comicId.value}/${chapters.value[currentIndex.value - 1].id}`) }
const nextChapter = () => { if (hasNext.value) router.push(`/read/${comicId.value}/${chapters.value[currentIndex.value + 1].id}`) }

onMounted(async () => {
  const data = await getChapters(comicId.value)
  chapters.value = data.chapters || []
  await loadImages()
})

watch(() => route.params.chapterId, loadImages)
</script>

<style scoped>
.reader-bar { display: flex; align-items: center; justify-content: space-between; padding: 12px 24px; background: var(--bg-surface); border: 1px solid var(--border); border-radius: var(--radius-md); margin-bottom: 16px; }
.reader-bar.bottom { margin-bottom: 0; margin-top: 16px; justify-content: center; }
.chapter-title { font-size: 15px; font-weight: 600; color: var(--text-primary); }
.nav-btns { display: flex; gap: 8px; }
.image-stream { display: flex; flex-direction: column; align-items: center; gap: 4px; background: var(--bg-surface); border-radius: var(--radius-md); padding: 16px; }
.image-stream img { max-width: 100%; border-radius: var(--radius-sm); }
</style>
