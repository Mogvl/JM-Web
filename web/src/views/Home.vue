<template>
  <div class="home">
    <el-tabs v-model="activeTab" @tab-change="onTabChange">
      <el-tab-pane label="最新上传" name="latest">
        <div v-if="loading" class="loading"><el-icon class="is-loading"><Loading /></el-icon> 加载中...</div>
        <div v-else-if="items.length === 0" class="empty">暂无数据</div>
        <div v-else class="grid">
          <div v-for="comic in items" :key="comic.id" class="comic-card" @click="$router.push(`/comic/${comic.id}`)">
            <img :src="comic.cover_url" :alt="comic.title" loading="lazy" />
            <div class="title">{{ comic.title }}</div>
          </div>
        </div>
        <div class="pagination-bar">
          <el-pagination v-model:current-page="page" :total="9999" :page-size="20" layout="prev, pager, next" @current-change="loadLatest" />
        </div>
      </el-tab-pane>
      <el-tab-pane v-for="cat in categories" :key="cat.id" :label="cat.name" :name="cat.id">
        <div v-if="loading" class="loading">加载中...</div>
        <div v-else-if="!catItems[cat.id] || catItems[cat.id].length === 0" class="empty">暂无数据</div>
        <div v-else class="grid">
          <div v-for="comic in catItems[cat.id]" :key="comic.id" class="comic-card" @click="$router.push(`/comic/${comic.id}`)">
            <img :src="comic.cover_url" :alt="comic.title" loading="lazy" />
            <div class="title">{{ comic.title }}</div>
          </div>
        </div>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getLatest, getIndex, getCategories } from '../api'
import { Loading } from '@element-plus/icons-vue'

const items = ref([])
const categories = ref([])
const catItems = ref({})
const loading = ref(false)
const page = ref(1)
const activeTab = ref('latest')

const loadLatest = async () => {
  loading.value = true
  try {
    const data = await getLatest(page.value - 1)
    items.value = data.items || []
  } catch (e) {
    items.value = []
  } finally {
    loading.value = false
  }
}

const loadCategory = async (catId) => {
  if (catItems.value[catId]) return
  loading.value = true
  try {
    const data = await getIndex(0)
    catItems.value[catId] = data.items || []
  } catch (e) {
    catItems.value[catId] = []
  } finally {
    loading.value = false
  }
}

const onTabChange = (tab) => {
  if (tab !== 'latest') {
    loadCategory(tab)
  }
}

onMounted(async () => {
  await loadLatest()
  try {
    const data = await getCategories()
    categories.value = data || []
  } catch (e) {
    categories.value = []
  }
})
</script>

<style scoped>
.home { padding: 0; }
.loading, .empty { text-align: center; padding: 80px 0; color: #999; }
.grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(160px, 1fr)); gap: 16px; }
.comic-card { cursor: pointer; transition: all 0.3s; border-radius: 8px; overflow: hidden; background: #fff; box-shadow: 0 1px 4px rgba(0,0,0,0.08); }
.comic-card:hover { transform: translateY(-4px); box-shadow: 0 4px 12px rgba(0,0,0,0.15); }
.comic-card img { width: 100%; height: 220px; object-fit: cover; display: block; }
.comic-card .title { padding: 10px; font-size: 13px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.pagination-bar { margin-top: 20px; display: flex; justify-content: center; }
</style>
