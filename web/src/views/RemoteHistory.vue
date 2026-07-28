<template>
  <div class="remote-history">
    <h2>远程历史</h2>
    <div v-if="loading" class="loading">加载中...</div>
    <div v-else class="empty">需要登录才能查看远程历史</div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
const loading = ref(false)
</script>

<style scoped>
.remote-history h2 { margin-bottom: 20px; }
.loading, .empty { text-align: center; padding: 100px 0; color: #999; }
</style>
