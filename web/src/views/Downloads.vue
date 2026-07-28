<template>
  <div class="downloads-page">
    <div class="toolbar">
      <el-select v-model="filter" placeholder="筛选" style="width:120px" @change="applyFilter">
        <el-option label="全部" value="all" />
        <el-option label="未完成" value="pending" />
        <el-option label="已完成" value="completed" />
        <el-option label="失败" value="failed" />
      </el-select>
      <el-input v-model="searchText" placeholder="搜索标题..." clearable style="width:200px" @input="applyFilter" />
      <div class="spacer"></div>
      <el-button type="danger" plain :disabled="!list.length" @click="handleClear">清空全部</el-button>
    </div>

    <div v-if="loading" class="loading-state">加载中...</div>
    <div v-else-if="filteredList.length === 0" class="empty-state">没有下载任务</div>

    <div v-else class="table-wrap">
      <div class="table-head">
        <div class="col-id">ID</div>
        <div class="col-cover">封面</div>
        <div class="col-title">标题</div>
        <div class="col-progress">下载进度</div>
        <div class="col-pages">页数</div>
        <div class="col-status">状态</div>
        <div class="col-time">时间</div>
        <div class="col-actions">操作</div>
      </div>
      <div class="table-body">
        <div v-for="dl in filteredList" :key="dl.id" class="table-row">
          <div class="col-id">{{ dl.id }}</div>
          <div class="col-cover"><img v-if="dl.comic?.cover_url" :src="dl.comic.cover_url" class="cover" /></div>
          <div class="col-title" :title="dl.comic?.title">{{ dl.comic?.title || `漫画 ${dl.comic_id}` }}</div>
          <div class="col-progress">
            <el-progress :percentage="dl.progress" :status="progressStatus(dl.status)" :stroke-width="6" :show-text="false" />
            <span class="pct">{{ dl.progress }}%</span>
          </div>
          <div class="col-pages">{{ dl.downloaded || 0 }}/{{ dl.total_pages || 0 }}</div>
          <div class="col-status"><el-tag :type="statusType(dl.status)" size="small" effect="light">{{ statusText(dl.status) }}</el-tag></div>
          <div class="col-time">{{ formatDate(dl.created_at) }}</div>
          <div class="col-actions">
            <el-button v-if="dl.status === 'completed'" size="small" plain @click="openFolder(dl)">打开</el-button>
            <el-button v-if="dl.status === 'completed'" type="primary" size="small" plain @click="readDownload(dl)">阅读</el-button>
            <el-button v-if="dl.status === 'failed'" type="warning" size="small" plain @click="retry(dl)">重试</el-button>
            <el-button type="danger" size="small" text @click="handleDelete(dl.id)">删除</el-button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { getDownloads, deleteDownload, clearDownloads, createDownload } from '../api'
import { ElMessage, ElMessageBox } from 'element-plus'

const router = useRouter()
const list = ref([])
const loading = ref(true)
const filter = ref('all')
const searchText = ref('')
let timer = null

const statusText = (s) => ({ pending: '等待中', downloading: '下载中', completed: '已完成', failed: '失败' }[s] || s)
const statusType = (s) => ({ pending: 'info', downloading: 'warning', completed: 'success', failed: 'danger' }[s] || 'info')
const progressStatus = (s) => s === 'completed' ? 'success' : s === 'failed' ? 'exception' : undefined

const formatDate = (d) => d ? new Date(d).toLocaleString('zh-CN', { month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit' }) : ''

const filteredList = computed(() => {
  let result = list.value
  if (filter.value !== 'all') {
    if (filter.value === 'pending') result = result.filter(d => d.status === 'downloading' || d.status === 'pending')
    else result = result.filter(d => d.status === filter.value)
  }
  if (searchText.value.trim()) {
    const q = searchText.value.trim().toLowerCase()
    result = result.filter(d => (d.comic?.title || '').toLowerCase().includes(q))
  }
  return result
})

const applyFilter = () => {}

const loadList = async () => {
  try { list.value = await getDownloads() || [] } catch (e) { list.value = [] }
  finally { loading.value = false }
  if (list.value.some(d => d.status === 'downloading' || d.status === 'pending')) {
    if (!timer) timer = setInterval(loadList, 2000)
  } else if (timer) { clearInterval(timer); timer = null }
}

const handleDelete = async (id) => {
  await deleteDownload(id)
  list.value = list.value.filter(i => i.id !== id)
  ElMessage.success('已删除')
}

const handleClear = async () => {
  try {
    await ElMessageBox.confirm('确定清空所有下载记录吗？', '提示', { type: 'warning' })
    await clearDownloads()
    list.value = []
    ElMessage.success('已清空')
  } catch (e) {}
}

const readDownload = (dl) => router.push(`/comic/${dl.comic_id}`)
const openFolder = (dl) => { if (dl.file_path) ElMessage.info(`保存路径: ${dl.file_path}`) }
const retry = async (dl) => {
  await createDownload(dl.comic_id, { title: dl.comic?.title, author: dl.comic?.author, cover: dl.comic?.cover_url })
  ElMessage.success('已重新加入下载')
  loadList()
}

onMounted(loadList)
onUnmounted(() => { if (timer) clearInterval(timer) })
</script>

<style scoped>
.toolbar { display: flex; align-items: center; gap: 12px; margin-bottom: 20px; flex-wrap: wrap; }
.spacer { flex: 1; }
.table-wrap { background: var(--bg-surface); border: 1px solid var(--border); border-radius: var(--radius-md); overflow: hidden; }
.table-head, .table-row { display: grid; grid-template-columns: 50px 60px 1fr 160px 80px 90px 120px 180px; align-items: center; gap: 8px; padding: 10px 12px; }
.table-head { background: var(--bg-elevated); font-size: 12px; font-weight: 600; color: var(--text-muted); border-bottom: 1px solid var(--border); }
.table-row { border-bottom: 1px solid var(--border); font-size: 13px; color: var(--text-secondary); transition: var(--transition); }
.table-row:hover { background: var(--bg-hover); }
.table-row:last-child { border-bottom: none; }
.col-id { color: var(--text-muted); }
.cover { width: 40px; height: 56px; object-fit: cover; border-radius: 4px; background: var(--bg-elevated); }
.col-title { white-space: nowrap; overflow: hidden; text-overflow: ellipsis; color: var(--text-primary); font-weight: 500; }
.col-progress { display: flex; align-items: center; gap: 8px; }
.pct { font-size: 12px; color: var(--text-muted); min-width: 32px; }
.col-actions { display: flex; gap: 4px; flex-wrap: wrap; }

@media (max-width: 768px) {
  .table-head { display: none; }
  .table-row { grid-template-columns: 60px 1fr 80px; grid-template-areas: "cover title status" "cover actions actions"; padding: 8px; gap: 6px; }
  .col-id, .col-progress, .col-pages, .col-time { display: none; }
  .col-cover { grid-area: cover; }
  .col-title { grid-area: title; }
  .col-status { grid-area: status; justify-self: end; }
  .col-actions { grid-area: actions; }
  .cover { width: 48px; height: 64px; }
}
</style>
