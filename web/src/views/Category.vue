<template>
  <div class="category">
    <div class="cat-buttons">
      <el-button :type="activeTab === 'titles' ? 'primary' : 'default'" @click="switchTab('titles')">标题</el-button>
      <el-button v-for="cat in categories" :key="cat.id" :type="activeTab === cat.id ? 'primary' : 'default'" @click="switchTab(cat.id)">{{ cat.name }}</el-button>
    </div>

    <!-- 标题页 -->
    <div v-if="activeTab === 'titles'" class="title-tags">
      <el-tag v-for="cat in categories" :key="cat.id" @click="switchCategory(cat.id)" class="tag" effect="plain" size="large">{{ cat.name }}</el-tag>
    </div>

    <!-- 分类内容页 -->
    <template v-else>
      <div class="controls">
        <el-select v-model="sort" @change="loadData" style="width:120px">
          <el-option label="最新" value="mr" />
          <el-option label="总排行" value="mv" />
          <el-option label="月排行" value="mv_m" />
          <el-option label="周排行" value="mv_w" />
          <el-option label="日排行" value="mv_t" />
          <el-option label="最多图片" value="mp" />
          <el-option label="最多爱心" value="tf" />
        </el-select>
        <span class="page-info">页：{{ page }}/{{ totalPages || 1 }}</span>
        <el-input-number v-model="jumpPage" :min="1" :max="totalPages||1" size="small" style="width:70px" />
        <el-button size="small" @click="jumpPageGo">跳转</el-button>
      </div>

      <div v-if="loading" class="loading">加载中...</div>
      <div v-else-if="items.length === 0" class="empty">暂无数据</div>
      <div v-else class="grid">
        <div v-for="comic in items" :key="comic.id" class="comic-card" @click="$router.push(`/comic/${comic.id}`)">
          <img :src="comic.cover_url" :alt="comic.title" loading="lazy" />
          <div class="title">{{ comic.title }}</div>
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
  } catch (e) {
    console.error('Load category failed:', e)
    items.value = []
  } finally {
    loading.value = false
  }
}

const switchTab = (tab) => {
  activeTab.value = tab
  if (tab === 'titles') return
  page.value = 1
  jumpPage.value = 1
  loadData()
}

const switchCategory = (id) => {
  activeTab.value = id
  page.value = 1
  jumpPage.value = 1
  loadData()
}

const jumpPageGo = () => {
  page.value = jumpPage.value
  loadData()
}

onMounted(async () => {
  try {
    const data = await getCategories()
    categories.value = data || []
  } catch (e) {
    console.error('Get categories failed:', e)
  }
})
</script>

<style scoped>
.cat-buttons { display: flex; flex-wrap: wrap; gap: 8px; margin-bottom: 20px; }
.title-tags { display: flex; flex-wrap: wrap; gap: 12px; }
.tag { cursor: pointer; font-size: 15px; padding: 8px 18px; }
.tag:hover { color: #409eff; border-color: #409eff; }
.controls { display: flex; align-items: center; gap: 12px; margin: 15px 0; flex-wrap: wrap; }
.page-info { font-size: 14px; color: #666; }
.loading, .empty { text-align: center; padding: 80px 0; color: #999; }
.grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(160px, 1fr)); gap: 16px; }
.comic-card { cursor: pointer; border-radius: 8px; overflow: hidden; background: #fff; box-shadow: 0 1px 4px rgba(0,0,0,0.08); transition: all 0.3s; }
.comic-card:hover { transform: translateY(-4px); box-shadow: 0 4px 12px rgba(0,0,0,0.15); }
.comic-card img { width: 100%; height: 220px; object-fit: cover; display: block; }
.comic-card .title { padding: 10px; font-size: 13px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
</style>
