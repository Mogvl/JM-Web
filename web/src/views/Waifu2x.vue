<template>
  <div class="waifu2x">
    <h2>Waifu2x 图片增强</h2>

    <el-card class="upload-card">
      <el-upload
        class="upload-area"
        drag
        :auto-upload="false"
        :on-change="handleFileChange"
        accept="image/*"
      >
        <el-icon class="el-icon--upload"><Upload /></el-icon>
        <div class="el-upload__text">拖拽图片到此处或 <em>点击选择</em></div>
      </el-upload>

      <div class="options">
        <el-form :inline="true">
          <el-form-item label="放大倍数">
            <el-select v-model="scale" style="width: 120px">
              <el-option label="2x" :value="2" />
              <el-option label="4x" :value="4" />
            </el-select>
          </el-form-item>
          <el-form-item label="降噪等级">
            <el-select v-model="noiseLevel" style="width: 120px">
              <el-option label="无" :value="0" />
              <el-option label="低" :value="1" />
              <el-option label="中" :value="2" />
              <el-option label="高" :value="3" />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="handleProcess" :loading="processing" :disabled="!originalImage">
              开始处理
            </el-button>
          </el-form-item>
        </el-form>
      </div>
    </el-card>

    <div class="preview" v-if="originalImage || enhancedImage">
      <el-row :gutter="20">
        <el-col :span="12">
          <h3>原图</h3>
          <div class="image-box">
            <img v-if="originalImage" :src="originalImage" />
          </div>
        </el-col>
        <el-col :span="12">
          <h3>增强后</h3>
          <div class="image-box">
            <img v-if="enhancedImage" :src="enhancedImage" />
            <div v-else class="placeholder">{{ processing ? '处理中...' : '等待处理' }}</div>
          </div>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { Upload } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const originalImage = ref('')
const enhancedImage = ref('')
const processing = ref(false)
const scale = ref(2)
const noiseLevel = ref(1)

const handleFileChange = (file) => {
  const reader = new FileReader()
  reader.onload = (e) => {
    originalImage.value = e.target.result
    enhancedImage.value = ''
  }
  reader.readAsDataURL(file.raw)
}

const handleProcess = async () => {
  processing.value = true
  // 模拟处理
  setTimeout(() => {
    enhancedImage.value = originalImage.value
    processing.value = false
    ElMessage.success('处理完成')
  }, 2000)
}
</script>

<style scoped>
.waifu2x h2 { margin-bottom: 20px; }
.upload-card { margin-bottom: 30px; }
.upload-area { width: 100%; }
.options { margin-top: 20px; }
.preview h3 { margin-bottom: 10px; text-align: center; }
.image-box { border: 1px solid #ddd; border-radius: 8px; overflow: hidden; min-height: 300px; display: flex; align-items: center; justify-content: center; }
.image-box img { max-width: 100%; max-height: 500px; }
.placeholder { color: #999; }
</style>
