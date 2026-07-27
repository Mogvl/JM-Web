<template>
  <div class="downloads">
    <h2>下载管理</h2>
    <div v-if="loading" class="loading">加载中...</div>
    <div v-else-if="list.length === 0" class="empty">没有下载任务</div>
    <div v-else class="list">
      <el-card v-for="dl in list" :key="dl.id" class="item">
        <div class="content">
          <img :src="dl.comic?.cover_url" class="cover" />
          <div class="info">
            <h3>{{ dl.comic?.title }}</h3>
            <el-progress :percentage="dl.progress" :status="dl.status === 'completed' ? 'success' : dl.status === 'failed' ? 'exception' : undefined" />
            <p class="status">{{ statusText(dl.status) }}</p>
          </div>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getDownloads } from '../api'

const list = ref([])
const loading = ref(true)

const statusText = (s) => ({ pending: '等待中', downloading: '下载中', completed: '已完成', failed: '失败' }[s] || s)

onMounted(async () => {
  list.value = await getDownloads()
  loading.value = false
})
</script>

<style scoped>
.downloads h2 { margin-bottom: 20px; }
.loading, .empty { text-align: center; padding: 100px 0; color: #999; }
.list { display: flex; flex-direction: column; gap: 16px; }
.item .content { display: flex; align-items: center; gap: 16px; }
.cover { width: 80px; height: 100px; object-fit: cover; border-radius: 4px; }
.info { flex: 1; }
.info h3 { margin-bottom: 12px; }
.status { color: #999; font-size: 14px; margin-top: 8px; }
</style>
