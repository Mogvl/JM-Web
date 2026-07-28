<template>
  <div class="home">
    <div class="cat-buttons">
      <el-button :type="activeTab === 'latest' ? 'primary' : 'default'" @click="switchTab('latest')">最新上传</el-button>
      <el-button v-for="(cats, name) in promoteMap" :key="name" :type="activeTab === name ? 'primary' : 'default'" @click="switchPromote(name)">{{ name }}</el-button>
    </div>

    <div v-if="loading" class="loading"><el-icon class="is-loading"><Loading /></el-icon> 加载中...</div>
    <div v-else-if="items.length === 0" class="empty">暂无数据</div>
    <div v-else class="grid">
      <div v-for="comic in items" :key="comic.id" class="comic-card" @click="$router.push(`/comic/${comic.id}`)">
        <img :src="comic.cover_url" :alt="comic.title" loading="lazy" />
        <div class="title">{{ comic.title }}</div>
      </div>
    </div>

    <div v-if="activeTab === 'latest'" class="pagination-bar">
      <el-pagination v-model:current-page="page" :total="9999" :page-size="20" layout="prev, pager, next" @current-change="loadLatest" />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getLatest, getIndex } from '../api'
import { Loading } from '@element-plus/icons-vue'

const items = ref([])
const loading = ref(false)
const activeTab = ref('latest')
const page = ref(1)
const promoteMap = ref({})

const loadLatest = async () => {
  loading.value = true
  try {
    const data = await getLatest(page.value - 1)
    items.value = data.items || []
  } catch (e) { items.value = [] }
  finally { loading.value = false }
}

const switchTab = (tab) => {
  activeTab.value = tab
  if (tab === 'latest') {
    loadLatest()
  }
}

const switchPromote = (name) => {
  activeTab.value = name
  items.value = promoteMap.value[name] || []
}

onMounted(async () => {
  await loadLatest()
  // 加载推广分类
  try {
    const data = await getIndex(0)
    if (data.items && data.items.length > 0) {
      promoteMap.value = { '推荐': data.items }
    }
  } catch (e) {}
})
</script>

<style scoped>
.home { padding: 0; }
.cat-buttons { display: flex; flex-wrap: wrap; gap: 8px; margin-bottom: 20px; }
.loading, .empty { text-align: center; padding: 80px 0; color: #999; }
.grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(160px, 1fr)); gap: 16px; }
.comic-card { cursor: pointer; border-radius: 8px; overflow: hidden; background: #fff; box-shadow: 0 1px 4px rgba(0,0,0,0.08); transition: all 0.3s; }
.comic-card:hover { transform: translateY(-4px); box-shadow: 0 4px 12px rgba(0,0,0,0.15); }
.comic-card img { width: 100%; height: 220px; object-fit: cover; display: block; }
.comic-card .title { padding: 10px; font-size: 13px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.pagination-bar { margin-top: 20px; display: flex; justify-content: center; }
</style>
