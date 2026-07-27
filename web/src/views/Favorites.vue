<template>
  <div class="favorites">
    <h2>我的收藏</h2>
    <div v-if="loading" class="loading">加载中...</div>
    <div v-else-if="list.length === 0" class="empty">还没有收藏</div>
    <div v-else class="grid">
      <el-card v-for="fav in list" :key="fav.id" @click="$router.push(`/comic/${fav.comic_id}`)" class="card">
        <img :src="fav.comic?.cover_url" class="cover" />
        <div class="title">{{ fav.comic?.title }}</div>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getFavorites } from '../api'

const list = ref([])
const loading = ref(true)

onMounted(async () => {
  list.value = await getFavorites()
  loading.value = false
})
</script>

<style scoped>
.favorites h2 { margin-bottom: 20px; }
.loading, .empty { text-align: center; padding: 100px 0; color: #999; }
.grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(180px, 1fr)); gap: 20px; }
.card { cursor: pointer; transition: all 0.3s; }
.card:hover { transform: translateY(-5px); }
.cover { width: 100%; height: 250px; object-fit: cover; }
.title { padding: 10px; font-size: 14px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
</style>
