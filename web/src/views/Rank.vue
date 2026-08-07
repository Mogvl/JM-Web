<template>
  <div class="rank">
    <h2>排行榜</h2>
    <div class="tabs">
      <el-tabs v-model="activeTab" @tab-change="loadRanking">
        <el-tab-pane label="日榜" name="daily" />
        <el-tab-pane label="周榜" name="weekly" />
        <el-tab-pane label="月榜" name="monthly" />
      </el-tabs>
    </div>
    <div v-if="loading" class="loading">加载中...</div>
    <div v-else-if="items.length === 0" class="empty">暂无数据</div>
    <div v-else class="grid">
      <el-card v-for="(comic, index) in items" :key="comic.id" @click="$router.push(`/comic/${comic.id}`)" class="card">
        <div class="rank-num">{{ index + 1 }}</div>
        <img :src="comic.cover_url" :alt="comic.title" class="cover" />
        <div class="info">
          <div class="title">{{ comic.title }}</div>
          <div class="author">{{ comic.author }}</div>
        </div>
      </el-card>
    </div>
    <el-pagination v-if="totalPages > 1" v-model:current-page="page" :total="totalPages * 10" layout="prev, pager, next" @current-change="loadRanking" class="pagination" />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getRanking } from '../api'

const items = ref([])
const loading = ref(false)
const page = ref(1)
const totalPages = ref(1)
const activeTab = ref('daily')

const loadRanking = async () => {
  loading.value = true
  try {
    const data = await getRanking(activeTab.value, page.value)
    items.value = data.items || []
    totalPages.value = data.total_pages || 1
  } catch (e) { items.value = [] }
  finally { loading.value = false }
}

onMounted(loadRanking)
</script>

<style scoped>
.rank h2 { margin-bottom: 20px; }
.tabs { margin-bottom: 20px; }
.loading, .empty { text-align: center; padding: 100px 0; color: #999; }
.grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(180px, 1fr)); gap: 20px; }
.card { cursor: pointer; transition: all 0.3s; position: relative; }
.card:hover { transform: translateY(-5px); box-shadow: 0 4px 12px rgba(0,0,0,0.15); }
.rank-num { position: absolute; top: 10px; left: 10px; background: var(--accent-grad); color: #fff; width: 24px; height: 24px; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 12px; z-index: 1; }
.cover { width: 100%; height: 250px; object-fit: cover; }
.info { padding: 10px; }
.title { font-size: 14px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.author { font-size: 12px; color: #999; }
.pagination { margin-top: 30px; justify-content: center; display: flex; }
</style>
