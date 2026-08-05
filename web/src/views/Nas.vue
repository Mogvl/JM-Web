<template>
  <div class="nas">
    <div class="header">
      <h2>NAS 上传</h2>
      <el-button type="primary" @click="showAdd = true">添加 NAS</el-button>
    </div>

    <!-- NAS 配置列表 -->
    <div class="nas-cards">
      <el-card v-for="nas in nasList" :key="nas.id" class="nas-card">
        <div class="nas-row">
          <div class="nas-name">{{ nas.title }}</div>
          <el-tag size="small">{{ typeText(nas.type) }}</el-tag>
          <div class="nas-addr">{{ nas.address }}{{ nas.port ? ':' + nas.port : '' }}</div>
        </div>
        <div class="nas-row">
          <el-button size="small" text @click="test(nas)">测试连接</el-button>
          <el-button size="small" text type="danger" @click="del(nas.id)">删除</el-button>
        </div>
      </el-card>
      <el-card v-if="nasList.length === 0" class="nas-card empty-card">
        <div class="empty">还没有 NAS 配置，点击右上角添加</div>
      </el-card>
    </div>

    <el-dialog v-model="showAdd" :title="editing ? '编辑 NAS' : '添加 NAS'" width="520px">
      <el-form :model="form" label-width="90px">
        <el-form-item label="名称">
          <el-input v-model="form.title" placeholder="例如：我的群晖" />
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="form.type" @change="onTypeChange">
            <el-option label="WebDAV" :value="0" />
            <el-option label="SMB" :value="1" />
            <el-option label="本地路径" :value="2" />
            <el-option label="SMBv3" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item label="地址">
          <el-input v-model="form.address" placeholder="http://192.168.1.100 或 192.168.1.100" />
        </el-form-item>
        <el-form-item label="端口" v-if="form.type !== 2">
          <el-input-number v-model="form.port" :min="0" :max="65535" />
        </el-form-item>
        <el-form-item label="用户名" v-if="form.type !== 2">
          <el-input v-model="form.user" />
        </el-form-item>
        <el-form-item label="密码" v-if="form.type !== 2">
          <el-input v-model="form.passwd" type="password" show-password />
        </el-form-item>
        <el-form-item label="保存路径" v-if="form.type === 2">
          <el-input v-model="form.path" placeholder="/volume1/comics" />
        </el-form-item>
        <el-form-item label="远程目录">
          <el-input v-model="form.dir" placeholder="远程目录，如 /comics" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showAdd = false">取消</el-button>
        <el-button @click="testCurrent">测试连接</el-button>
        <el-button type="primary" @click="save">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'

const nasList = ref([])
const showAdd = ref(false)
const editing = ref(null)
const form = ref({ title: '', type: 0, address: '', port: 0, user: '', passwd: '', path: '', dir: '' })

const typeText = (t) => ({ 0: 'WebDAV', 1: 'SMB', 2: '本地路径', 3: 'SMBv3' }[t] || '未知')

const onTypeChange = () => {
  if (form.value.type === 0) { form.value.address = form.value.address || 'http://192.168.1.28'; form.value.port = form.value.port || 5005 }
  else if (form.value.type === 1 || form.value.type === 3) { form.value.address = form.value.address || '192.168.1.28'; form.value.port = form.value.port || 0 }
}

const test = (nas) => {
  ElMessage.info(`测试连接 ${nas.address}...`)
  // 浏览器无法直接测 SMB；仅作提示
  setTimeout(() => ElMessage.success('连接测试已发起'), 600)
}

const testCurrent = () => {
  ElMessage.success(`将测试 ${form.value.address} 的连接`)
}

const del = async (id) => {
  try {
    await ElMessageBox.confirm('删除该 NAS 配置？', '提示', { type: 'warning' })
    nasList.value = nasList.value.filter(n => n.id !== id)
    ElMessage.success('已删除')
  } catch (e) {}
}

const save = () => {
  if (!form.value.title.trim()) { ElMessage.warning('请输入名称'); return }
  if (editing.value) {
    Object.assign(editing.value, { ...form.value })
  } else {
    nasList.value.push({ id: Date.now(), ...form.value })
  }
  showAdd.value = false
  editing.value = null
  ElMessage.success('已保存')
}

onMounted(() => {
  nasList.value = JSON.parse(localStorage.getItem('nasList') || '[]')
})
</script>

<style scoped>
.nas .header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
.nas-cards { display: grid; grid-template-columns: repeat(auto-fill, minmax(280px, 1fr)); gap: 16px; }
.nas-card { margin-bottom: 0; }
.nas-row { display: flex; align-items: center; gap: 10px; margin-bottom: 8px; }
.nas-name { font-weight: 600; }
.nas-addr { color: #999; font-size: 13px; flex: 1; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.empty-card .empty { text-align: center; color: #999; padding: 40px 0; }
</style>