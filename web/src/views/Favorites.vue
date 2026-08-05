<template>
  <div class="favorites">
    <div class="header">
      <h2>我的收藏</h2>
      <el-button type="primary" @click="showFoldDialog = true">新建文件夹</el-button>
    </div>

    <div class="folders">
      <el-tag :type="activeFold === '' ? 'primary' : 'info'" @click="activeFold = ''; load()">全部</el-tag>
      <el-tag v-for="fold in folders" :key="fold.id" closable :type="activeFold === fold.id ? 'primary' : 'info'" @click="activeFold = fold.id; load()" @close="delFolder(fold.id)">{{ fold.name }}</el-tag>
    </div>

    <div v-if="loading" class="loading">加载中...</div>
    <div v-else-if="list.length === 0" class="empty">还没有收藏</div>
    <div v-else class="grid">
      <el-card v-for="fav in list" :key="fav.comic_id" @click="$router.push(`/comic/${fav.comic_id}`)" class="card">
        <img :src="fav.comic?.cover_url" class="cover" />
        <div class="info">
          <div class="title">{{ fav.comic?.title }}</div>
          <div class="ops">
            <el-select v-model="fav.folder_id" size="small" placeholder="移动文件夹" @change="move(fav, $event)" @click.stop>
              <el-option label="全部" value="" />
              <el-option v-for="f in folders" :key="f.id" :label="f.name" :value="f.id" />
            </el-select>
            <el-button type="danger" text size="small" @click.stop="remove(fav.comic_id)">取消收藏</el-button>
          </div>
        </div>
      </el-card>
    </div>

    <el-dialog v-model="showFoldDialog" title="新建文件夹" width="300px">
      <el-input v-model="newFoldName" placeholder="文件夹名称" @keyup.enter="createFold" />
      <template #footer>
        <el-button @click="showFoldDialog = false">取消</el-button>
        <el-button type="primary" @click="createFold">创建</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getFavorites, removeFavorite, addFavoriteFolder, deleteFavoriteFolder, moveFavorite } from '../api'
import { ElMessage, ElMessageBox } from 'element-plus'

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
    list.value = (data.list || data || []).map(f => ({ ...f, folder_id: '' }))
    if (data.folders) {
      folders.value = data.folders.map(f => ({ id: f.FID, name: f.name }))
    }
  } catch (e) {
    list.value = []
  } finally {
    loading.value = false
  }
}

const remove = async (id) => {
  try {
    await removeFavorite(id)
    list.value = list.value.filter(i => i.comic_id !== id)
    ElMessage.success('已取消收藏')
  } catch (e) { ElMessage.error('操作失败') }
}

const createFold = async () => {
  if (!newFoldName.value.trim()) return
  try {
    await addFavoriteFolder(newFoldName.value.trim())
    folders.value.push({ id: 'f' + Date.now(), name: newFoldName.value.trim() })
    showFoldDialog.value = false
    newFoldName.value = ''
    ElMessage.success('文件夹已创建')
  } catch (e) { ElMessage.error('创建失败') }
}

const delFolder = async (fid) => {
  try {
    await ElMessageBox.confirm('删除文件夹？', '提示', { type: 'warning' })
    await deleteFavoriteFolder(fid)
    folders.value = folders.value.filter(f => f.id !== fid)
    if (activeFold.value === fid) activeFold.value = ''
    ElMessage.success('已删除')
  } catch (e) {}
}

const move = async (fav, fid) => {
  try {
    await moveFavorite(fav.comic_id, fid)
    ElMessage.success('已移动')
  } catch (e) {
    fav.folder_id = ''
    ElMessage.error('移动失败')
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
.ops { display: flex; align-items: center; gap: 6px; }
.ops .el-select { width: 110px; }
</style>