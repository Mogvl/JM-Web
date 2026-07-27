<template>
  <div class="weekly">
    <h2>每周更新</h2>
    <div v-if="loading" class="loading">加载中...</div>
    <div v-else class="empty">功能开发中...</div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
const loading = ref(false)
</script>

<style scoped>
.weekly h2 { margin-bottom: 20px; }
.loading, .empty { text-align: center; padding: 100px 0; color: #999; }
</style>
