<template>
  <div class="remote-history">
    <div class="header">
      <h2>远程历史</h2>
      <el-button v-if="list.length" size="small" @click="load">刷新</el-button>
    </div>
    <div v-if="loading" class="loading">加载中...</div>
    <div v-else-if="list.length === 0" class="empty">暂无远程历史记录</div>
    <div v-else class="list">
      <el-card v-for="item in list" :key="item.comic_id" class="item">
        <div class="content">
          <img :src="item.comic?.cover_url" class="cover" />
          <div class="info">
            <h3 @click="$router.push(`/comic/${item.comic_id}`)">{{ item.comic?.title }}</h3>
            <p>阅读于 {{ formatDate(item.last_read_at) }}</p>
          </div>
          <el-button type="primary" text @click="continueRead(item)">继续阅读</el-button>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getHistory } from '../api'

const router = useRouter()
const list = ref([])
const loading = ref(true)

const formatDate = (d) => d ? new Date(d).toLocaleString('zh-CN') : ''

const continueRead = (item) => {
  // 有章节ID则跳阅读器，否则详情页
  if (item.chapter_id) router.push(`/read/${item.comic_id}/${item.chapter_id}`)
  else router.push(`/comic/${item.comic_id}`)
}

const load = async () => {
  loading.value = true
  try { list.value = await getHistory() || [] } catch (e) { list.value = [] }
  finally { loading.value = false }
}

onMounted(load)
</script>

<style scoped>
.remote-history .header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
.loading, .empty { text-align: center; padding: 100px 0; color: #999; }
.list { display: flex; flex-direction: column; gap: 16px; }
.item .content { display: flex; align-items: center; gap: 16px; }
.cover { width: 80px; height: 100px; object-fit: cover; border-radius: 4px; }
.info { flex: 1; }
.info h3 { cursor: pointer; margin-bottom: 8px; }
.info h3:hover { color: var(--accent); }
.info p { color: #999; font-size: 14px; }
</style>
