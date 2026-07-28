<template>
  <div class="category">
    <el-tabs v-model="activeTab" @tab-change="onTabChange">
      <el-tab-pane label="标题" name="titles">
        <div class="title-tags">
          <div v-for="(group, title) in titleGroups" :key="title" class="title-group">
            <h3>{{ title }}</h3>
            <el-tag v-for="tag in group" :key="tag" @click="searchTag(tag)" class="tag" effect="plain" size="large">{{ tag }}</el-tag>
          </div>
        </div>
      </el-tab-pane>
      <el-tab-pane v-for="cat in categories" :key="cat.id" :label="cat.name" :name="cat.id" />
    </el-tabs>

    <div class="controls">
      <el-select v-model="sort" @change="loadData" style="width: 120px">
        <el-option label="最新" value="mr" />
        <el-option label="总排行" value="mv" />
        <el-option label="月排行" value="mv_m" />
        <el-option label="周排行" value="mv_w" />
        <el-option label="日排行" value="mv_t" />
        <el-option label="最多图片" value="mp" />
        <el-option label="最多爱心" value="tf" />
      </el-select>
      <span class="page-info">页：{{ page }}/{{ totalPages }}</span>
      <el-input-number v-model="jumpPage" :min="1" :max="totalPages || 1" size="small" style="width: 80px" />
      <el-button @click="jump" size="small">跳转</el-button>
    </div>

    <div v-if="loading" class="loading">加载中...</div>
    <div v-else-if="items.length === 0 && activeTab !== 'titles'" class="empty">暂无数据</div>
    <div v-else class="grid">
      <div v-for="comic in items" :key="comic.id" class="comic-card" @click="$router.push(`/comic/${comic.id}`)">
        <img :src="comic.cover_url" :alt="comic.title" loading="lazy" />
        <div class="title">{{ comic.title }}</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getCategories } from '../api'
import { getCategoryFilter } from '../api'

const categories = ref([])
const titleGroups = ref({})
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
  } catch (e) {
    items.value = []
  } finally {
    loading.value = false
  }
}

const onTabChange = (tab) => {
  if (tab !== 'titles') {
    page.value = 1
    loadData()
  }
}

const searchTag = (tag) => {
  window.location.href = '/search?q=' + encodeURIComponent(tag)
}

const jump = () => {
  page.value = jumpPage.value
  loadData()
}

onMounted(async () => {
  try {
    const cats = await getCategories()
    categories.value = cats || []
  } catch (e) {
    categories.value = []
  }
})
</script>

<style scoped>
.category { padding: 0; }
.title-tags { padding: 10px 0; }
.title-group { margin-bottom: 20px; }
.title-group h3 { margin-bottom: 10px; font-size: 16px; }
.tag { cursor: pointer; margin: 4px; }
.tag:hover { color: #409eff; border-color: #409eff; }
.controls { display: flex; align-items: center; gap: 12px; margin: 15px 0; }
.page-info { font-size: 14px; color: #666; }
.loading, .empty { text-align: center; padding: 80px 0; color: #999; }
.grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(160px, 1fr)); gap: 16px; }
.comic-card { cursor: pointer; transition: all 0.3s; border-radius: 8px; overflow: hidden; background: #fff; box-shadow: 0 1px 4px rgba(0,0,0,0.08); }
.comic-card:hover { transform: translateY(-4px); box-shadow: 0 4px 12px rgba(0,0,0,0.15); }
.comic-card img { width: 100%; height: 220px; object-fit: cover; display: block; }
.comic-card .title { padding: 10px; font-size: 13px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
</style>
