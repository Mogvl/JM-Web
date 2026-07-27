<template>
  <div class="search">
    <h2>搜索: {{ $route.query.q }}</h2>
    <div v-if="loading" class="loading">加载中...</div>
    <div v-else-if="items.length === 0" class="empty">没有找到漫画</div>
    <div v-else class="grid">
      <el-card v-for="comic in items" :key="comic.id" @click="$router.push(`/comic/${comic.id}`)" class="card">
        <img :src="comic.cover_url" :alt="comic.title" class="cover" />
        <div class="info">
          <div class="title">{{ comic.title }}</div>
          <div class="author">{{ comic.author }}</div>
        </div>
      </el-card>
    </div>
    <el-pagination v-if="totalPages > 1" v-model:current-page="page" :total="totalPages * 10" layout="prev, pager, next" @current-change="doSearch" class="pagination" />
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { search } from '../api'

const route = useRoute()
const items = ref([])
const loading = ref(false)
const page = ref(1)
const totalPages = ref(1)

const doSearch = async () => {
  loading.value = true
  try {
    const data = await search(route.query.q, page.value)
    items.value = data.items || []
    totalPages.value = data.total_pages || 1
  } catch (e) {
    items.value = []
  } finally {
    loading.value = false
  }
}

onMounted(doSearch)
watch(() => route.query.q, doSearch)
</script>

<style scoped>
.search h2 { margin-bottom: 20px; }
.loading, .empty { text-align: center; padding: 100px 0; color: #999; }
.grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(180px, 1fr)); gap: 20px; }
.card { cursor: pointer; transition: all 0.3s; }
.card:hover { transform: translateY(-5px); box-shadow: 0 4px 12px rgba(0,0,0,0.15); }
.cover { width: 100%; height: 250px; object-fit: cover; }
.info { padding: 10px; }
.title { font-size: 14px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.author { font-size: 12px; color: #999; }
.pagination { margin-top: 30px; justify-content: center; display: flex; }
</style>
