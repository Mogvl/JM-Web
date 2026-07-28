<template>
  <div class="local-fav">
    <h2>本地收藏</h2>
    <div class="grid">
      <el-card v-if="items.length === 0" class="empty-card">
        <div class="empty">还没有本地收藏</div>
      </el-card>
      <el-card v-for="item in items" :key="item.id" @click="$router.push(`/comic/${item.id}`)" class="card" v-else>
        <img :src="item.cover" class="cover" />
        <div class="title">{{ item.title }}</div>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
const items = ref(JSON.parse(localStorage.getItem('localFav') || '[]'))
</script>

<style scoped>
.local-fav h2 { margin-bottom: 20px; }
.grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(180px, 1fr)); gap: 20px; }
.empty-card .empty { text-align: center; padding: 60px; color: #999; }
.card { cursor: pointer; }
.cover { width: 100%; height: 250px; object-fit: cover; }
.title { padding: 10px; font-size: 14px; }
</style>
