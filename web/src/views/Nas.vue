<template>
  <div class="nas">
    <div class="header">
      <h2>NAS 同步</h2>
      <el-button type="primary" @click="showAdd = true">添加任务</el-button>
    </div>

    <el-table :data="tasks" style="width: 100%">
      <el-table-column prop="comic_title" label="漫画" width="200" />
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="statusType(row.status)">{{ statusText(row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="progress" label="进度" width="150">
        <template #default="{ row }">
          <el-progress :percentage="row.progress" />
        </template>
      </el-table-column>
      <el-table-column prop="path" label="NAS 路径" />
      <el-table-column label="操作" width="150">
        <template #default="{ row }">
          <el-button text @click="handleDelete(row.id)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="showAdd" title="添加 NAS 同步任务" width="500px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="漫画 ID">
          <el-input v-model="form.comic_id" placeholder="输入漫画 ID" />
        </el-form-item>
        <el-form-item label="NAS 路径">
          <el-input v-model="form.path" placeholder="/volume1/comics" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showAdd = false">取消</el-button>
        <el-button type="primary" @click="handleAdd">添加</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'

const tasks = ref([])
const showAdd = ref(false)
const form = ref({ comic_id: '', path: '/volume1/comics' })

const statusText = (s) => ({ pending: '等待中', syncing: '同步中', completed: '已完成', failed: '失败' }[s] || s)
const statusType = (s) => ({ pending: 'info', syncing: 'warning', completed: 'success', failed: 'danger' }[s] || 'info')

const handleAdd = () => {
  tasks.value.push({
    id: Date.now(),
    comic_title: `漫画 ${form.value.comic_id}`,
    status: 'pending',
    progress: 0,
    path: form.value.path
  })
  showAdd.value = false
  ElMessage.success('任务已添加')
}

const handleDelete = (id) => {
  tasks.value = tasks.value.filter(t => t.id !== id)
}
</script>

<style scoped>
.nas .header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
</style>
