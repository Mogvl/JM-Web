<template>
  <div class="download-all">
    <h2>批量下载</h2>
    <el-card class="card">
      <el-form label-width="100px">
        <el-form-item label="漫画ID">
          <el-input v-model="comicIds" type="textarea" :rows="5" placeholder="每行一个漫画ID" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="startBatch">开始批量下载</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { createDownload } from '../api'
import { ElMessage } from 'element-plus'

const comicIds = ref('')

const startBatch = async () => {
  const ids = comicIds.value.split('\n').map(s => s.trim()).filter(Boolean)
  if (ids.length === 0) { ElMessage.warning('请输入漫画ID'); return }
  for (const id of ids) {
    try {
      await createDownload(id)
      ElMessage.success(`漫画 ${id} 已添加下载`)
    } catch (e) {
      ElMessage.error(`漫画 ${id} 添加失败`)
    }
  }
}
</script>

<style scoped>
.download-all h2 { margin-bottom: 20px; }
.card { max-width: 500px; }
</style>
