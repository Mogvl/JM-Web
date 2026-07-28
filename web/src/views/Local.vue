<template>
  <div class="local">
    <h2>本地阅读</h2>

    <div class="upload-area" @drop.prevent="handleDrop" @dragover.prevent>
      <el-icon :size="48"><FolderOpened /></el-icon>
      <p>拖拽漫画文件夹到此处</p>
      <p class="hint">支持 ZIP / RAR / 文件夹格式</p>
    </div>

    <div v-if="localBooks.length > 0" class="books">
      <h3>本地漫画</h3>
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

const router = useRouter()
const localBooks = ref(JSON.parse(localStorage.getItem('localBooks') || '[]'))
const localFavorites = ref(JSON.parse(localStorage.getItem('localFavorites') || '[]'))

const handleDrop = (e) => {
  const files = Array.from(e.dataTransfer.files)
  files.forEach(file => {
    if (file.name.match(/\.(zip|rar|cbr|cbz)$/i) || file.name.endsWith('/')) {
      localBooks.value.push({
        id: Date.now().toString(),
        name: file.name.replace(/\.\w+$/, ''),
        cover: '',
        path: file.path
      })
    }
  })
  localStorage.setItem('localBooks', JSON.stringify(localBooks.value))
}

const readLocal = (book) => {
  router.push(`/read/local/${book.id}`)
}
</script>

<style scoped>
.local h2 { margin-bottom: 20px; }
.upload-area { border: 2px dashed #409eff; border-radius: 12px; padding: 60px; text-align: center; cursor: pointer; margin-bottom: 30px; }
.upload-area .el-icon { color: #409eff; margin-bottom: 10px; }
.upload-area .hint { color: #999; font-size: 12px; margin-top: 10px; }
.books h3 { margin-bottom: 15px; }
.grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(150px, 1fr)); gap: 16px; margin-bottom: 30px; }
.card { cursor: pointer; }
.cover { width: 100%; height: 200px; object-fit: cover; }
.title { padding: 8px; font-size: 13px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
</style>
