<template>
  <div class="favorites">
    <div class="header">
      <h2>我的收藏</h2>
      <el-button type="primary" @click="showFoldDialog = true">新建文件夹</el-button>
    </div>

    <div class="folders">
      <el-tag :type="activeFold === '' ? 'primary' : 'info'" @click="activeFold = ''; load()">全部</el-tag>
      <el-tag v-for="fold in folders" :key="fold.id" :type="activeFold === fold.id ? 'primary' : 'info'" @click="activeFold = fold.id; load()">{{ fold.name }}</el-tag>
    </div>

    <div v-if="loading" class="loading">加载中...</div>
    <div v-else-if="list.length === 0" class="empty">还没有收藏</div>
    <div v-else class="grid">
      <el-card v-for="fav in list" :key="fav.id" @click="$router.push(`/comic/${fav.comic_id}`)" class="card">
        <img :src="fav.comic?.cover_url" class="cover" />
        <div class="info">
          <div class="title">{{ fav.comic?.title }}</div>
          <el-button type="danger" text size="small" @click.stop="remove(fav.comic_id)">取消收藏</el-button>
        </div>
      </el-card>
    </div>

    <el-dialog v-model="showFoldDialog" title="新建文件夹" width="300px">
      <el-input v-model="newFoldName" placeholder="文件夹名称" />
      <template #footer>
        <el-button @click="showFoldDialog = false">取消</el-button>
        <el-button type="primary" @click="createFold">创建</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getFavorites, removeFavorite } from '../api'
import { ElMessage } from 'element-plus'

const list = ref([])
const folders = ref([])
const loading = ref(true)
const activeFold = ref('')
const showFoldDialog = ref(false)
const newFoldName = ref('')

const load = async () => {
  loading.value = true
  try {
    const data = await getFavorites()
    list.value = data || []
  } catch (e) {
    list.value = []
  } finally {
    loading.value = false
  }
}

const remove = async (id) => {
  await removeFavorite(id)
  list.value = list.value.filter(i => i.comic_id !== id)
  ElMessage.success('已取消收藏')
}

const createFold = () => {
  if (newFoldName.value.trim()) {
    folders.value.push({ id: Date.now().toString(), name: newFoldName.value })
    showFoldDialog.value = false
    newFoldName.value = ''
    ElMessage.success('文件夹已创建')
  }
}

onMounted(load)
</script>

<style scoped>
.favorites .header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
.folders { display: flex; gap: 10px; margin-bottom: 20px; flex-wrap: wrap; }
.folders .el-tag { cursor: pointer; }
.loading, .empty { text-align: center; padding: 100px 0; color: #999; }
.grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(180px, 1fr)); gap: 20px; }
.card { cursor: pointer; transition: all 0.3s; }
.card:hover { transform: translateY(-5px); }
.cover { width: 100%; height: 250px; object-fit: cover; }
.info { padding: 10px; }
.title { font-size: 14px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; margin-bottom: 8px; }
</style>
