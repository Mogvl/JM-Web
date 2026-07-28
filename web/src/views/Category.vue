<template>
  <div class="category-page">
    <div class="filter-bar">
      <button :class="['chip', activeTab === 'titles' && 'active']" @click="switchTab('titles')">标题</button>
      <button v-for="cat in categories" :key="cat.id" :class="['chip', activeTab === cat.id && 'active']" @click="switchTab(cat.id)">{{ cat.name }}</button>
    </div>

    <div v-if="activeTab === 'titles'" class="title-tags">
      <el-tag v-for="cat in categories" :key="cat.id" @click="switchTab(cat.id)" class="tag-chip" effect="plain" size="large">{{ cat.name }}</el-tag>
    </div>

    <template v-else>
      <div class="controls">
        <el-select v-model="sort" @change="loadData" style="width:120px" size="default">
          <el-option label="最新" value="mr" />
          <el-option label="总排行" value="mv" />
          <el-option label="月排行" value="mv_m" />
          <el-option label="周排行" value="mv_w" />
          <el-option label="日排行" value="mv_t" />
          <el-option label="最多图片" value="mp" />
          <el-option label="最多爱心" value="tf" />
        </el-select>
        <span class="page-info">页 {{ page }}/{{ totalPages || 1 }}</span>
        <el-input-number v-model="jumpPage" :min="1" :max="totalPages||1" size="small" style="width:70px" />
        <el-button size="small" @click="jump">跳转</el-button>
      </div>

      <div v-if="loading" class="loading-state">加载中...</div>
      <div v-else-if="items.length === 0" class="empty-state">暂无数据</div>
      <div v-else class="comic-grid">
        <div v-for="comic in items" :key="comic.id" class="comic-card" @click="$router.push(`/comic/${comic.id}`)">
          <img :src="comic.cover_url" :alt="comic.title" loading="lazy" />
          <div class="card-title">{{ comic.title }}</div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getCategories, getCategoryFilter } from '../api'

const categories = ref([])
const items = ref([])
const loading = ref(false)
const activeTab = ref('titles')
const sort = ref('mr')
const page = ref(1)
const totalPages = ref(1)
const jumpPage = ref(1)

const loadData = async () => {
  if (activeTab.value === 'titles') return
  loading.value = true
  try {
    const data = await getCategoryFilter(activeTab.value, sort.value, page.value)
    items.value = data.items || []
    totalPages.value = data.total_pages || 1
  } catch (e) { items.value = [] }
  finally { loading.value = false }
}

const switchTab = (tab) => {
  activeTab.value = tab
  if (tab === 'titles') return
  page.value = 1; jumpPage.value = 1; loadData()
}

const jump = () => { page.value = jumpPage.value; loadData() }

onMounted(async () => {
  try { categories.value = await getCategories() } catch (e) { categories.value = [] }
})
</script>

<style scoped>
.filter-bar { display: flex; flex-wrap: wrap; gap: 8px; margin-bottom: 24px; }
.chip { padding: 6px 16px; border-radius: 20px; border: 1px solid var(--border); background: var(--bg-surface); color: var(--text-secondary); font-size: 13px; cursor: pointer; transition: var(--transition); }
.chip:hover { border-color: var(--border-light); color: var(--text-primary); }
.chip.active { background: var(--accent); border-color: var(--accent); color: #fff; font-weight: 500; }
.title-tags { display: flex; flex-wrap: wrap; gap: 12px; }
.tag-chip { cursor: pointer; }
.controls { display: flex; align-items: center; gap: 12px; margin-bottom: 20px; }
.page-info { font-size: 13px; color: var(--text-muted); }
</style>
