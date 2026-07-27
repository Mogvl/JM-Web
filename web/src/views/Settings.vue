<template>
  <div class="settings">
    <h2>设置</h2>
    <el-card class="settings-card">
      <el-form label-width="120px">
        <el-form-item label="API 地址">
          <el-input v-model="apiUrl" placeholder="https://jmcomic.me" />
        </el-form-item>
        <el-form-item label="代理设置">
          <el-input v-model="proxy" placeholder="不使用代理" />
        </el-form-item>
        <el-form-item label="图片质量">
          <el-select v-model="quality" style="width: 100%">
            <el-option label="原图" value="original" />
            <el-option label="高清" value="high" />
            <el-option label="普通" value="normal" />
          </el-select>
        </el-form-item>
        <el-form-item label="阅读模式">
          <el-radio-group v-model="readMode">
            <el-radio label="scroll">滚动</el-radio>
            <el-radio label="page">翻页</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="save">保存设置</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { ElMessage } from 'element-plus'

const apiUrl = ref(localStorage.getItem('apiUrl') || '')
const proxy = ref(localStorage.getItem('proxy') || '')
const quality = ref(localStorage.getItem('quality') || 'original')
const readMode = ref(localStorage.getItem('readMode') || 'scroll')

const save = () => {
  localStorage.setItem('apiUrl', apiUrl.value)
  localStorage.setItem('proxy', proxy.value)
  localStorage.setItem('quality', quality.value)
  localStorage.setItem('readMode', readMode.value)
  ElMessage.success('设置已保存')
}
</script>

<style scoped>
.settings h2 { margin-bottom: 20px; }
.settings-card { max-width: 600px; }
</style>
