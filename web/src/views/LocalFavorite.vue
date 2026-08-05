<template>
  <div class="local-fav">
    <h2>本地收藏</h2>
    <div class="grid">
      <el-card v-if="items.length === 0" class="empty-card">
        <div class="empty">还没有本地收藏</div>
      </el-card>
      <el-card v-for="item in items" :key="item.id" class="card" v-else>
        <img :src="item.cover" class="cover" @click="$router.push(`/comic/${item.id}`)" />
        <div class="title" @click="$router.push(`/comic/${item.id}`)">{{ item.title }}</div>
        <el-button type="danger" text size="small" @click="remove(item.id)">移除</el-button>
      </el-card>
    </div>
    <div v-if="items.length === 0 && canImport" class="import-section">
      <el-button type="primary" plain @click="importFromFav">从在线收藏导入到本地</el-button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { getFavorites } from '../api'
import { ElMessage } from 'element-plus'

const items = ref(JSON.parse(localStorage.getItem('localFav') || '[]'))
const inFlight = ref(false)
const canImport = computed(() => localStorage.getItem('token') && localStorage.getItem('token') !== 'guest')

const remove = (id) => {
  items.value = items.value.filter(i => i.id !== id)
  localStorage.setItem('localFav', JSON.stringify(items.value))
  ElMessage.success('已移除')
}

const importFromFav = async () => {
  if (inFlight.value) return
  inFlight.value = true
  try {
    const data = await getFavorites()
    const favs = data.list || data || []
    favs.forEach(f => {
      if (!items.value.find(x => x.id === f.comic_id)) {
        items.value.push({ id: f.comic_id, title: f.comic?.title, cover: f.comic?.cover_url })
      }
    })
    localStorage.setItem('localFav', JSON.stringify(items.value))
    ElMessage.success(`已导入 ${favs.length} 条收藏`)
  } catch (e) { ElMessage.error('导入失败') }
  finally { inFlight.value = false }
}

onMounted(() => {})
</script>

<style scoped>
.local-fav h2 { margin-bottom: 20px; }
.grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(180px, 1fr)); gap: 20px; }
.empty-card .empty { text-align: center; padding: 60px; color: #999; }
.card { cursor: pointer; }
.cover { width: 100%; height: 250px; object-fit: cover; }
.title { padding: 10px; font-size: 14px; }
</style>
