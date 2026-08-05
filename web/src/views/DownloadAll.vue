<template>
  <div class="download-all">
    <div class="header">
      <h2>批量下载</h2>
      <span class="hint">支持多行，每行一个漫画ID（纯数字或 JMxxxx）</span>
    </div>
    <el-card class="card">
      <el-form label-width="100px">
        <el-form-item label="漫画ID">
          <el-input v-model="comicIds" type="textarea" :rows="6" placeholder="每行一个漫画ID，例如：\n421926\n422xxx\n..." />
        </el-form-item>
        <el-form-item label="格式">
          <el-radio-group v-model="format">
            <el-radio-button label="jpg">JPG</el-radio-button>
            <el-radio-button label="webp">WEBP</el-radio-button>
            <el-radio-button label="png">PNG</el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="running" @click="startBatch">开始批量下载</el-button>
          <el-button v-if="running" @click="cancel = true">停止</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <div v-if="results.length" class="results">
      <h3>任务结果</h3>
      <el-table :data="results" max-height="400">
        <el-table-column prop="id" label="漫画ID" width="120" />
        <el-table-column label="状态" width="120">
          <template #default="{ row }">
            <el-tag :type="row.status === 'ok' ? 'success' : 'danger'" size="small">{{ row.status === 'ok' ? '已添加' : '失败' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="msg" label="信息" />
      </el-table>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { createDownload } from '../api'
import { ElMessage } from 'element-plus'

const comicIds = ref('')
const format = ref('jpg')
const running = ref(false)
const cancel = ref(false)
const results = ref([])

const startBatch = async () => {
  const ids = comicIds.value.split('\n').map(s => s.trim().replace(/^jm/gi, '')).filter(Boolean)
  if (ids.length === 0) { ElMessage.warning('请输入漫画ID'); return }
  running.value = true
  cancel.value = false
  results.value = []
  for (const id of ids) {
    if (cancel.value) break
    try {
      await createDownload(id, { format: format.value })
      results.value.push({ id, status: 'ok', msg: '已添加到下载管理' })
    } catch (e) {
      results.value.push({ id, status: 'err', msg: '添加失败' })
    }
  }
  running.value = false
  ElMessage[results.value.some(r => r.status === 'err') ? 'warning' : 'success'](`完成：成功 ${results.value.filter(r => r.status === 'ok').length} / ${results.value.length}`)
}
</script>

<style scoped>
.download-all .header { display: flex; align-items: baseline; gap: 12px; margin-bottom: 20px; }
.download-all .hint { color: #999; font-size: 13px; }
.card { max-width: 520px; }
.results { max-width: 520px; margin-top: 24px; }
.results h3 { margin-bottom: 12px; }
</style>
