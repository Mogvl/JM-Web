<template>
  <div class="home">
    <div class="filter-bar">
      <button :class="['chip', activeTab === 'latest' && 'active']" @click="switchTab('latest')">最新上传</button>
      <button v-for="(items, name) in promoteMap" :key="name" :class="['chip', activeTab === name && 'active']" @click="switchTab(name)">{{ name }}</button>
    </div>

    <div v-if="loading" class="loading-state"><el-icon class="is-loading"><Loading /></el-icon> 加载中...</div>
    <div v-else-if="items.length === 0" class="empty-state">暂无数据</div>
    <div v-else class="comic-grid">
      <div v-for="comic in items" :key="comic.id" class="comic-card" @click="openComic(comic)">
        <img :src="comic.cover_url" :alt="comic.title" loading="lazy" />
        <div class="card-title">{{ comic.title }}</div>
        <div v-if="comic.author" class="card-author">{{ comic.author }}</div>
      </div>
    </div>

    <div v-if="activeTab === 'latest'" class="pager">
      <el-pagination v-model:current-page="page" :total="9999" :page-size="20" layout="prev, pager, next" @current-change="loadLatest" />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getLatest, getIndex } from '../api'
import { Loading } from '@element-plus/icons-vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const items = ref([])
const loading = ref(false)
const activeTab = ref('latest')
const page = ref(1)
const promoteMap = ref({})

const openComic = (comic) => {
  router.push(`/comic/${comic.id}`)
}

const loadLatest = async () => {
  loading.value = true
  try { const data = await getLatest(page.value - 1); items.value = data.items || [] }
  catch (e) { items.value = [] }
  finally { loading.value = false }
}

const switchTab = (name) => {
  activeTab.value = name
  if (name === 'latest') loadLatest()
  else items.value = promoteMap.value[name] || []
}

onMounted(async () => {
  await loadLatest()
  try {
    const data = await getIndex(0)
    if (data && typeof data === 'object') {
      const map = {}
      for (const [key, val] of Object.entries(data)) if (Array.isArray(val) && val.length > 0) map[key] = val
      if (Object.keys(map).length > 0) promoteMap.value = map
    }
  } catch (e) {}
})
</script>

<style scoped>
.filter-bar { display: flex; flex-wrap: wrap; gap: 8px; margin-bottom: 24px; }
.chip { padding: 6px 16px; border-radius: 20px; border: 1px solid var(--border); background: var(--bg-surface); color: var(--text-secondary); font-size: 13px; cursor: pointer; transition: var(--transition); }
.chip:hover { border-color: var(--border-light); color: var(--text-primary); }
.chip.active { background: var(--accent); border-color: var(--accent); color: #fff; font-weight: 500; }
.pager { margin-top: 24px; display: flex; justify-content: center; }
</style>
