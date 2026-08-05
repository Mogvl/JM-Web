<template>
  <div class="reader-page">
    <BackNav :title="currentTitle" :fallback="`/comic/${comicId}`">
      <template #right>
        <div class="nav-btns">
          <button class="mini-btn" :disabled="!hasPrev" @click="prevChapter">上一章</button>
          <button class="mini-btn" :disabled="!hasNext" @click="nextChapter">下一章</button>
        </div>
      </template>
    </BackNav>

    <div v-if="loading" class="loading-state">加载中...</div>
    <div v-else-if="images.length === 0" class="empty-state">暂无图片</div>
    <div v-else class="image-stream">
      <img v-for="(img, i) in proxiedImages" :key="i" :src="img" loading="lazy" />
    </div>

    <div class="reader-bar bottom">
      <button class="mini-btn" :disabled="!hasPrev" @click="prevChapter">上一章</button>
      <button class="mini-btn" :disabled="!hasNext" @click="nextChapter">下一章</button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowLeft } from '@element-plus/icons-vue'
import { getImages, getChapters } from '../api'
import BackNav from '../components/BackNav.vue'

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

// 图片通过后端代理，避免防盗链
const proxiedImages = computed(() => images.value.map(url => `/api/image?url=${encodeURIComponent(url)}`))

const loadImages = async () => {
  loading.value = true
  try { const data = await getImages(chapterId.value); images.value = data.images || [] }
  catch (e) { images.value = [] }
  finally { loading.value = false }
}

const prevChapter = () => { if (hasPrev.value) router.push(`/read/${comicId.value}/${chapters.value[currentIndex.value - 1].id}`) }
const nextChapter = () => { if (hasNext.value) router.push(`/read/${comicId.value}/${chapters.value[currentIndex.value + 1].id}`) }

onMounted(async () => {
  try {
    const data = await getChapters(comicId.value)
    chapters.value = data.chapters || []
  } catch (e) { chapters.value = [] }
  await loadImages()
})

watch(() => route.params.chapterId, loadImages)
</script>

<style scoped>
.mini-btn { padding: 7px 14px; border: 1px solid var(--border); background: var(--bg-surface); color: var(--text-secondary); border-radius: var(--radius-sm); font-size: 13px; font-family: inherit; cursor: pointer; transition: var(--transition); }
.mini-btn:hover:not(:disabled) { border-color: var(--accent); color: var(--accent); background: var(--accent-soft); }
.mini-btn:disabled { opacity: 0.45; cursor: not-allowed; }
.nav-btns { display: flex; gap: 8px; flex-shrink: 0; }
.reader-bar.bottom { display: flex; justify-content: center; gap: 12px; margin: 16px auto 0; padding: 16px 24px; background: var(--bg-surface); border: 1px solid var(--border); border-radius: var(--radius-md); }
.reader-bar.bottom .mini-btn { border-radius: var(--radius-sm); }
.image-stream { display: flex; flex-direction: column; align-items: center; gap: 4px; background: var(--bg-surface); border-radius: var(--radius-md); padding: 16px; }
.image-stream img { max-width: 100%; border-radius: var(--radius-sm); }
</style>
