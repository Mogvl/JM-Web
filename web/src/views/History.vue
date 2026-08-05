<template>
  <div class="history">
    <div class="header">
      <h2>阅读历史</h2>
      <el-button type="danger" @click="handleClear">清空历史</el-button>
    </div>
    <div v-if="loading" class="loading">加载中...</div>
    <div v-else-if="list.length === 0" class="empty">还没有记录</div>
    <div v-else class="list">
      <el-card v-for="item in list" :key="item.id" class="item">
        <div class="content">
          <img :src="item.comic?.cover_url" class="cover" />
          <div class="info">
            <h3 @click="$router.push(`/comic/${item.comic_id}`)">{{ item.comic?.title }}</h3>
            <p>{{ formatDate(item.last_read_at) }}</p>
          </div>
          <el-button type="danger" text @click="handleDelete(item.id)">删除</el-button>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getHistory, deleteHistory, clearHistory } from '../api'
import { ElMessage } from 'element-plus'

const list = ref([])
const loading = ref(true)

const formatDate = (d) => new Date(d).toLocaleString('zh-CN')

const handleDelete = async (id) => {
  await deleteHistory(id)
  list.value = list.value.filter(i => i.id !== id)
  ElMessage.success('已删除')
}

const handleClear = async () => {
  await clearHistory()
  list.value = []
  ElMessage.success('已清空')
}

onMounted(async () => {
  try {
    list.value = await getHistory()
  } catch (e) {
    list.value = []
  }
  loading.value = false
})
</script>

<style scoped>
.history .header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
.loading, .empty { text-align: center; padding: 100px 0; color: #999; }
.list { display: flex; flex-direction: column; gap: 16px; }
.item .content { display: flex; align-items: center; gap: 16px; }
.cover { width: 80px; height: 100px; object-fit: cover; border-radius: 4px; }
.info { flex: 1; }
.info h3 { cursor: pointer; margin-bottom: 8px; }
.info h3:hover { color: #409eff; }
.info p { color: #999; font-size: 14px; }
</style>
