<template>
  <div class="batch-sr">
    <h2>批量超分辨率</h2>

    <el-card class="card">
      <el-upload drag :auto-upload="false" :on-change="handleFiles" multiple accept="image/*">
        <el-icon class="el-icon--upload"><Upload /></el-icon>
        <div class="el-upload__text">拖拽图片到此处或 <em>点击选择</em></div>
      </el-upload>

      <div class="options">
        <el-form :inline="true">
          <el-form-item label="放大倍数">
            <el-select v-model="scale" style="width: 100px">
              <el-option label="2x" :value="2" />
              <el-option label="4x" :value="4" />
            </el-select>
          </el-form-item>
          <el-form-item label="降噪">
            <el-select v-model="noise" style="width: 100px">
              <el-option label="无" :value="0" />
              <el-option label="低" :value="1" />
              <el-option label="中" :value="2" />
              <el-option label="高" :value="3" />
            </el-select>
          </el-form-item>
          <el-form-item label="格式">
            <el-select v-model="format" style="width: 100px">
              <el-option label="PNG" value="png" />
              <el-option label="JPG" value="jpg" />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="startProcess" :loading="processing">{{ processing ? '处理中...' : '开始处理' }}</el-button>
          </el-form-item>
        </el-form>
      </div>
    </el-card>

    <div v-if="files.length > 0" class="files">
      <h3>文件列表 ({{ files.length }})</h3>
      <el-table :data="files" max-height="400">
        <el-table-column prop="name" label="文件名" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'done' ? 'success' : row.status === 'error' ? 'danger' : 'info'">{{ statusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="80">
          <template #default="{ row, index }">
            <el-button text @click="removeFile(index)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { Upload } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const files = ref([])
const processing = ref(false)
const scale = ref(2)
const noise = ref(1)
const format = ref('png')

const statusText = (s) => ({ pending: '等待', processing: '处理中', done: '完成', error: '失败' }[s] || s)

const handleFiles = (uploadFiles) => {
  uploadFiles.forEach(f => {
    if (!files.value.find(x => x.name === f.name)) {
      files.value.push({ name: f.name, raw: f.raw, status: 'pending' })
    }
  })
}

const removeFile = (index) => { files.value.splice(index, 1) }

const startProcess = async () => {
  if (files.value.length === 0) { ElMessage.warning('请先添加图片'); return }
  processing.value = true
  for (const file of files.value) {
    file.status = 'processing'
    await new Promise(r => setTimeout(r, 1000))
    file.status = 'done'
  }
  processing.value = false
  ElMessage.success('处理完成')
}
</script>

<style scoped>
.batch-sr h2 { margin-bottom: 20px; }
.card { margin-bottom: 30px; }
.options { margin-top: 20px; }
.files h3 { margin-bottom: 15px; }
</style>
