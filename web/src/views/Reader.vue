<template>
  <div class="reader-page">
    <BackNav :title="currentTitle" :fallback="`/comic/${comicId}`">
      <template #right>
        <div class="reader-tools">
          <select v-model="mode" class="mode-select" @change="onModeChange">
            <option value="scroll">滚动阅读</option>
            <option value="single">单页</option>
            <option value="double">双页</option>
          </select>
        </div>
        <div class="nav-btns">
          <button class="mini-btn" :disabled="!hasPrev" @click="prevChapter">上一章</button>
          <button class="mini-btn" :disabled="!hasNext" @click="nextChapter">下一章</button>
        </div>
      </template>
    </BackNav>

    <div v-if="loading" class="loading-state">加载中...</div>
    <div v-else-if="images.length === 0" class="empty-state">暂无图片</div>
    <div v-else class="reader-body" ref="readerBody" @scroll="onReaderScroll">
      <!-- 顶部进度提示 -->
      <div class="progress-hint" v-if="!hideProgress">{{ curPage }}/{{ images.length }}</div>
      <!-- 滚动 / 单页模式 -->
      <div v-if="mode !== 'double'" :class="['image-stream', mode === 'single' && 'single-mode']">
        <img v-for="(img, i) in proxiedImages" :key="i" :src="img" loading="lazy" :class="mode === 'single' && 'single-img'" />
      </div>
      <!-- 双页模式 -->
      <div v-else class="double-stream">
        <div v-for="n in doublePages" :key="n" class="double-page">
          <img v-if="normalized[n*2]" :src="normalized[n*2]" loading="lazy" />
          <div v-else class="blank"></div>
          <img v-if="normalized[n*2+1]" :src="normalized[n*2+1]" loading="lazy" />
          <div v-else class="blank"></div>
        </div>
      </div>
    </div>

    <div class="reader-bar bottom">
      <button class="mini-btn" :disabled="!hasPrev" @click="prevChapter">上一章</button>
      <button class="mini-btn" :disabled="!hasNext" @click="nextChapter">下一章</button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowLeft } from '@element-plus/icons-vue'
import { getImages, getChapters, saveProgress } from '../api'
import BackNav from '../components/BackNav.vue'

const route = useRoute()
const router = useRouter()
const images = ref([])
const chapters = ref([])
const loading = ref(false)
const mode = ref(localStorage.getItem('readerMode') || 'scroll')
const readerBody = ref(null)
const curPage = ref(1)
const hideProgress = ref(false)
let progressHideTimer = null

const comicId = computed(() => route.params.comicId)
const chapterId = computed(() => route.params.chapterId)
const currentIndex = computed(() => chapters.value.findIndex(c => c.id === chapterId.value))
const hasPrev = computed(() => currentIndex.value > 0)
const hasNext = computed(() => currentIndex.value >= 0 && currentIndex.value < chapters.value.length - 1)
const currentTitle = computed(() => { const ch = chapters.value.find(c => c.id === chapterId.value); return ch ? (ch.title || `第${ch.sort_order}章`) : '' })

// 图片通过后端代理，避免防盗链
const proxiedImages = computed(() => images.value.map(url => `/api/image?url=${encodeURIComponent(url)}`))

const onReaderScroll = () => {
  const el = readerBody.value
  if (!el) return
  const imgs = el.querySelectorAll('.image-stream img, .double-page img')
  const box = el.getBoundingClientRect()
  for (let i = 0; i < imgs.length; i++) {
    const r = imgs[i].getBoundingClientRect()
    if (r.top < box.top + el.clientHeight * 0.5 && r.bottom > box.top) { curPage.value = i + 1; break }
  }
  hideProgress.value = false
  clearTimeout(progressHideTimer)
  progressHideTimer = setTimeout(() => { hideProgress.value = true }, 1500)
}

// 双页模式：把图片序号重排为 0,2,4...1,3,5... 的数组，便于左页读单页
const normalized = computed(() => {
  const arr = proxiedImages.value
  if (arr.length <= 1) return arr
  const evens = arr.filter((_, i) => i % 2 === 0)
  const odds = arr.filter((_, i) => i % 2 === 1)
  return [...evens, ...odds]
})
const doublePages = computed(() => {
  const n = normalized.value.length
  const pages = []
  for (let i = 0; i < n; i += 2) pages.push(i / 2)
  return pages
})

const onModeChange = () => localStorage.setItem('readerMode', mode.value)

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

// 记录阅读进度（滚动到第几张图时上报，避免频繁请求）
let progressTimer = null
const reportProgress = () => {
  clearTimeout(progressTimer)
  progressTimer = setTimeout(async () => {
    const imgs = document.querySelectorAll('.image-stream img')
    let page = 0
    for (let i = 0; i < imgs.length; i++) {
      const r = imgs[i].getBoundingClientRect()
      if (r.top < window.innerHeight && r.bottom > 0) { page = i + 1; break }
    }
    try { await saveProgress(comicId.value, chapterId.value, page) } catch (e) {}
  }, 800)
}
onMounted(() => window.addEventListener('scroll', reportProgress, { passive: true }))
onUnmounted(() => { window.removeEventListener('scroll', reportProgress); clearTimeout(progressTimer) })
</script>

<style scoped>
.mini-btn { padding: 7px 14px; border: 1px solid var(--border); background: var(--bg-surface); color: var(--text-secondary); border-radius: var(--radius-sm); font-size: 13px; font-family: inherit; cursor: pointer; transition: var(--transition); }
.mini-btn:hover:not(:disabled) { border-color: var(--accent); color: var(--accent); background: var(--accent-soft); }
.mini-btn:disabled { opacity: 0.45; cursor: not-allowed; }
.nav-btns { display: flex; gap: 8px; flex-shrink: 0; }
.reader-tools { display: flex; align-items: center; flex-shrink: 0; }
.mode-select { padding: 6px 8px; border: 1px solid var(--border); border-radius: var(--radius-sm); background: var(--bg-surface); color: var(--text-secondary); font-size: 13px; font-family: inherit; cursor: pointer; }
.reader-bar.bottom { display: flex; justify-content: center; gap: 12px; margin: 16px auto 0; padding: 16px 24px; background: var(--bg-surface); border: 1px solid var(--border); border-radius: var(--radius-md); }
.reader-bar.bottom .mini-btn { border-radius: var(--radius-sm); }
.image-stream { display: flex; flex-direction: column; align-items: center; gap: 4px; background: var(--bg-surface); border-radius: var(--radius-md); padding: 16px; }
.image-stream img { max-width: 100%; border-radius: var(--radius-sm); }
.reader-body { position: relative; }
.progress-hint { position: sticky; top: 8px; z-index: 30; margin: 0 auto; width: fit-content; padding: 4px 12px; background: rgba(0,0,0,0.6); color: #fff; border-radius: 20px; font-size: 12px; pointer-events: none; transition: opacity 0.3s; }
.image-stream.single-mode { gap: 12px; }
.image-stream .single-img { max-width: 70%; }
.double-stream { display: flex; flex-direction: column; gap: 12px; align-items: center; }
.double-page { display: flex; gap: 4px; align-items: center; justify-content: center; background: var(--bg-surface); border-radius: var(--radius-md); padding: 16px; }
.double-page img { max-width: 49%; border-radius: var(--radius-sm); }
.double-page .blank { width: 49%; }

@media (max-width: 768px) {
  .reader-tools { display: none; }
  .image-stream .single-img { max-width: 100%; }
}
</style>
