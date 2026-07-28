<template>
  <div class="home">
    <el-tabs v-model="activeTab">
      <el-tab-pane label="最新上传" name="latest">
        <div v-if="loading" class="loading">加载中...</div>
        <div v-else class="grid">
          <div v-for="comic in items" :key="comic.id" class="comic-card" @click="$router.push(`/comic/${comic.id}`)">
            <img :src="comic.cover_url" :alt="comic.title" />
            <div class="title">{{ comic.title }}</div>
          </div>
        </div>
        <div class="pagination-bar">
          <el-pagination
            v-model:current-page="page"
            :page-size="20"
            :total="9999"
            layout="prev, pager, next"
            @current-change="load"
          />
        </div>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getLatest } from '../api'

const items = ref([])
const loading = ref(false)
const page = ref(1)
const activeTab = ref('latest')

const load = async () => {
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

onMounted(load)
</script>

<style scoped>
.home { padding: 0; }
.loading { text-align: center; padding: 100px 0; color: #999; }
.grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(160px, 1fr)); gap: 16px; }
.comic-card { cursor: pointer; transition: all 0.3s; border-radius: 8px; overflow: hidden; background: #fff; }
.comic-card:hover { transform: translateY(-5px); box-shadow: 0 4px 12px rgba(0,0,0,0.15); }
.comic-card img { width: 100%; height: 220px; object-fit: cover; display: block; }
.comic-card .title { padding: 10px; font-size: 13px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.pagination-bar { margin-top: 20px; display: flex; justify-content: center; }
</style>
