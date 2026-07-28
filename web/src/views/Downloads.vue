<template>
  <div class="downloads-page">
    <div class="header">
      <h2>下载管理</h2>
      <el-button type="danger" plain @click="handleClear" :disabled="!list.length">清空下载</el-button>
    </div>

    <div v-if="loading" class="loading-state">加载中...</div>
    <div v-else-if="list.length === 0" class="empty-state">没有下载任务</div>
    <div v-else class="list">
      <el-card v-for="dl in list" :key="dl.id" class="item" shadow="never">
        <div class="content">
          <img :src="dl.comic?.cover_url" class="cover" />
          <div class="info">
            <h3 class="comic-title">{{ dl.comic?.title || `漫画 ${dl.comic_id}` }}</h3>
            <div class="progress-row">
              <el-progress :percentage="dl.progress" :status="progressStatus(dl.status)" :stroke-width="8" />
            </div>
            <p class="status">
              <el-tag :type="statusType(dl.status)" size="small" effect="light">{{ statusText(dl.status) }}</el-tag>
              <span class="page-count">{{ dl.downloaded || 0 }}/{{ dl.total_pages || 0 }} 页</span>
            </p>
          </div>
          <div class="actions">
            <el-button v-if="dl.status === 'completed'" type="primary" size="small" plain @click="readDownload(dl)">阅读</el-button>
            <el-button type="danger" size="small" text @click="handleDelete(dl.id)">删除</el-button>
          </div>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { getDownloads, deleteDownload, clearDownloads } from '../api'
import { ElMessage } from 'element-plus'

const router = useRouter()
const list = ref([])
const loading = ref(true)
let timer = null

const statusText = (s) => ({ pending: '等待中', downloading: '下载中', completed: '已完成', failed: '失败' }[s] || s)
const statusType = (s) => ({ pending: 'info', downloading: 'warning', completed: 'success', failed: 'danger' }[s] || 'info')
const progressStatus = (s) => s === 'completed' ? 'success' : s === 'failed' ? 'exception' : undefined

const loadList = async () => {
  try { list.value = await getDownloads() || [] } catch (e) { list.value = [] }
  finally { loading.value = false }
  // 下载中时自动刷新
  if (list.value.some(d => d.status === 'downloading' || d.status === 'pending')) {
    if (!timer) timer = setInterval(loadList, 2000)
  } else if (timer) {
    clearInterval(timer)
    timer = null
  }
}

const handleDelete = async (id) => {
  await deleteDownload(id)
  list.value = list.value.filter(i => i.id !== id)
  ElMessage.success('已删除')
}

const handleClear = async () => {
  await clearDownloads()
  list.value = []
  ElMessage.success('已清空')
}

const readDownload = (dl) => {
  router.push(`/comic/${dl.comic_id}`)
}

onMounted(loadList)
onUnmounted(() => { if (timer) clearInterval(timer) })
</script>

<style scoped>
.downloads-page .header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
.list { display: flex; flex-direction: column; gap: 12px; }
.item { border-radius: var(--radius-md); }
.item :deep(.el-card__body) { padding: 16px; }
.content { display: flex; align-items: center; gap: 16px; }
.cover { width: 64px; height: 88px; object-fit: cover; border-radius: var(--radius-sm); background: var(--bg-elevated); flex-shrink: 0; }
.info { flex: 1; min-width: 0; }
.comic-title { font-size: 14px; font-weight: 600; color: var(--text-primary); margin-bottom: 10px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.progress-row { margin-bottom: 8px; }
.status { display: flex; align-items: center; gap: 8px; font-size: 12px; color: var(--text-muted); }
.page-count { color: var(--text-muted); }
.actions { display: flex; flex-direction: column; gap: 8px; flex-shrink: 0; }

@media (max-width: 768px) {
  .content { gap: 10px; }
  .cover { width: 48px; height: 66px; }
  .actions { flex-direction: row; }
}
</style>
