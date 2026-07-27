<template>
  <div class="comments">
    <h2>我的评论</h2>
    <div v-if="loading" class="loading">加载中...</div>
    <div v-else class="empty">请先登录查看评论</div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
const loading = ref(false)
</script>

<style scoped>
.comments h2 { margin-bottom: 20px; }
.loading, .empty { text-align: center; padding: 100px 0; color: #999; }
</style>
