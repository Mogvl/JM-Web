<template>
  <div class="local">
    <h2>本地阅读</h2>

    <div class="upload-area" @drop.prevent="handleDrop" @dragover.prevent @click="pickFiles">
      <el-icon :size="48"><FolderOpened /></el-icon>
      <p>拖拽漫画文件夹 / 压缩包到此处，或点击选择</p>
      <p class="hint">支持 ZIP / RAR / CBR / CBZ 及文件夹</p>
      <input ref="fileInput" type="file" multiple accept=".zip,.rar,.cbr,.cbz,application/zip,application/x-rar-compressed" style="display:none" @change="handlePick" />
    </div>

    <div class="import-actions" v-if="pending.length">
      <div class="pending-list">
        <div v-for="(p, i) in pending" :key="i" class="pending-item">
          <span>{{ p.name }}</span>
          <el-tag size="small" :type="p.status === 'done' ? 'success' : p.status === 'err' ? 'danger' : 'info'">{{ p.status === 'done' ? '已导入' : p.status === 'err' ? '失败' : '待导入' }}</el-tag>
        </div>
      </div>
      <el-button type="primary" size="small" @click="confirmImport">确认导入</el-button>
      <el-button size="small" @click="pending = []">清空</el-button>
    </div>

    <div v-if="localBooks.length > 0" class="books">
      <h3>本地漫画 <el-button size="small" text type="danger" @click="clearAll">清空</el-button></h3>
      <div class="grid">
        <el-card v-for="book in localBooks" :key="book.id" @click="readLocal(book)" class="card">
          <img :src="book.cover" class="cover" />
          <div class="title">{{ book.name }}</div>
        </el-card>
      </div>
    </div>

    <div v-if="localFavorites.length > 0" class="books">
      <h3>本地收藏</h3>
      <div class="grid">
        <el-card v-for="book in localFavorites" :key="book.id" @click="readLocal(book)" class="card">
          <img :src="book.cover" class="cover" />
          <div class="title">{{ book.name }}</div>
        </el-card>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { FolderOpened } from '@element-plus/icons-vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'

const router = useRouter()
const fileInput = ref(null)
const localBooks = ref(JSON.parse(localStorage.getItem('localBooks') || '[]'))
const localFavorites = ref(JSON.parse(localStorage.getItem('localFavorites') || '[]'))
const pending = ref([])

const pickFiles = () => fileInput.value?.click()

const handleDrop = (e) => {
  const files = Array.from(e.dataTransfer.files)
  addPending(files)
}

const handlePick = (e) => {
  addPending(Array.from(e.target.files))
  e.target.value = ''
}

const addPending = (files) => {
  files.forEach(file => {
    if (file.name.match(/\.(zip|rar|cbr|cbz)$/i) || file.name.endsWith('/')) {
      if (!pending.value.find(p => p.name === file.name)) {
        pending.value.push({ name: file.name.replace(/\.\w+$/, ''), file, status: 'pending' })
      }
    } else {
      ElMessage.warning(`跳过不支持的文件: ${file.name}`)
    }
  })
}

const confirmImport = () => {
  pending.value.forEach(p => {
    localBooks.value.push({ id: Date.now() + Math.random().toString(16).slice(2), name: p.name, cover: '', path: p.file?.path || '' })
    p.status = 'done'
  })
  localStorage.setItem('localBooks', JSON.stringify(localBooks.value))
  ElMessage.success(`已导入 ${pending.value.length} 个漫画`)
  pending.value = []
}

const clearAll = () => {
  localBooks.value = []
  localStorage.setItem('localBooks', '[]')
  ElMessage.success('已清空')
}

const readLocal = (book) => {
  router.push(`/read/local/${book.id}`)
}
</script>

<style scoped>
.local h2 { margin-bottom: 20px; }
.upload-area { border: 2px dashed var(--accent); border-radius: 12px; padding: 60px; text-align: center; cursor: pointer; margin-bottom: 20px; transition: background 0.2s; }
.upload-area:hover { background: rgba(64,158,255,0.05); }
.upload-area .el-icon { color: var(--accent); margin-bottom: 10px; }
.upload-area .hint { color: #999; font-size: 12px; margin-top: 10px; }
.import-actions { margin-bottom: 20px; display: flex; align-items: center; gap: 12px; flex-wrap: wrap; }
.pending-list { display: flex; flex-wrap: wrap; gap: 8px; flex: 1; min-width: 200px; }
.pending-item { display: flex; align-items: center; gap: 6px; padding: 4px 10px; background: var(--bg-surface); border: 1px solid var(--border); border-radius: var(--radius-sm); font-size: 12px; }
.books h3 { margin-bottom: 15px; }
.grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(150px, 1fr)); gap: 16px; margin-bottom: 30px; }
.card { cursor: pointer; }
.cover { width: 100%; height: 200px; object-fit: cover; }
.title { padding: 8px; font-size: 13px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
</style>
