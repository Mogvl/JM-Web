<template>
  <div class="category">
    <el-tabs v-model="activeTab" @tab-change="onTabChange" class="cat-tabs">
      <el-tab-pane label="标题" name="titles" />
      <el-tab-pane v-for="cat in categories" :key="cat.id" :label="cat.name" :name="cat.id" />
    </el-tabs>

    <div v-if="activeTab === 'titles'">
      <div class="title-tags">
        <el-tag v-for="cat in categories" :key="cat.id" @click="switchToCategory(cat.id)" class="tag" effect="plain" size="large">{{ cat.name }}</el-tag>
      </div>
    </div>

    <template v-else>
      <div class="controls">
        <el-select v-model="sort" @change="loadData(1)" style="width: 120px">
          <el-option label="最新" value="mr" />
          <el-option label="总排行" value="mv" />
          <el-option label="月排行" value="mv_m" />
          <el-option label="周排行" value="mv_w" />
          <el-option label="日排行" value="mv_t" />
          <el-option label="最多图片" value="mp" />
          <el-option label="最多爱心" value="tf" />
        </el-select>
        <span class="page-info">页：{{ page }}/{{ totalPages }}</span>
        <el-input-number v-model="jumpPage" :min="1" :max="totalPages||1" size="small" style="width:70px" />
        <el-button size="small" @click="jump">跳转</el-button>
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

const loadData = async (p) => {
  if (activeTab.value === 'titles') return
  loading.value = true
  page.value = p || page.value
  try {
    const data = await getCategoryFilter(activeTab.value, sort.value, page.value)
    items.value = data.items || []
    totalPages.value = data.total_pages || 1
  } catch (e) { items.value = [] }
  finally { loading.value = false }
}

const onTabChange = (tab) => {
  if (tab !== 'titles') loadData(1)
}

const switchToCategory = (id) => {
  activeTab.value = id
  page.value = 1
  loadData(1)
}

const jump = () => {
  if (jumpPage.value >= 1) loadData(jumpPage.value)
}

onMounted(async () => {
  try { categories.value = await getCategories() } catch (e) { categories.value = [] }
})
</script>

<style scoped>
.cat-tabs { margin-bottom: 10px; }
.title-tags { display: flex; flex-wrap: wrap; gap: 12px; padding: 10px 0; }
.tag { cursor: pointer; font-size: 15px; padding: 8px 18px; }
.tag:hover { color: #409eff; border-color: #409eff; }
.controls { display: flex; align-items: center; gap: 12px; margin: 15px 0; }
.page-info { font-size: 14px; color: #666; }
.loading, .empty { text-align: center; padding: 80px 0; color: #999; }
.grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(160px, 1fr)); gap: 16px; }
.comic-card { cursor: pointer; border-radius: 8px; overflow: hidden; background: #fff; box-shadow: 0 1px 4px rgba(0,0,0,0.08); transition: all 0.3s; }
.comic-card:hover { transform: translateY(-4px); box-shadow: 0 4px 12px rgba(0,0,0,0.15); }
.comic-card img { width: 100%; height: 220px; object-fit: cover; display: block; }
.comic-card .title { padding: 10px; font-size: 13px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
</style>
