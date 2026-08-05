<template>
  <div class="comments">
    <div class="header">
      <h2>我的评论</h2>
      <el-button v-if="list.length" size="small" @click="load">刷新</el-button>
    </div>
    <div v-if="loading" class="loading">加载中...</div>
    <div v-else-if="list.length === 0" class="empty">还没有发表过评论</div>
    <div v-else class="list">
      <el-card v-for="c in list" :key="c.id" class="item">
        <div class="content">
          <el-avatar :size="36" :src="c.avatar">{{ (c.author || '?').charAt(0) }}</el-avatar>
          <div class="info">
            <div class="head"><span class="author">{{ c.author }}</span><span class="time">{{ c.create_time }}</span></div>
            <div class="text">{{ c.content }}</div>
          </div>
        </div>
      </el-card>
    </div>
    <div v-if="more" class="load-more"><el-button size="small" @click="load(page + 1)">加载更多</el-button></div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getMyComments } from '../api'

const list = ref([])
const loading = ref(true)
const page = ref(1)
const more = ref(false)

const load = async (p = 1) => {
  loading.value = true
  try {
    const data = await getMyComments(p)
    if (p === 1) list.value = data || []
    else list.value = list.value.concat(data || [])
    page.value = p
    more.value = (data || []).length >= 20
  } catch (e) { list.value = [] }
  finally { loading.value = false }
}

onMounted(() => load())
</script>

<style scoped>
.comments .header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
.loading, .empty { text-align: center; padding: 100px 0; color: #999; }
.list { display: flex; flex-direction: column; gap: 12px; }
.item .content { display: flex; gap: 12px; }
.info { flex: 1; min-width: 0; }
.head { display: flex; justify-content: space-between; align-items: center; margin-bottom: 6px; }
.author { font-weight: 600; }
.time { color: #999; font-size: 12px; }
.text { color: #666; line-height: 1.6; word-break: break-word; }
.load-more { text-align: center; margin-top: 20px; }
</style>