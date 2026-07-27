<template>
  <div class="local">
    <h2>本地阅读</h2>
    <div class="empty">
      <p>支持从本地文件夹导入漫画进行阅读</p>
      <el-button type="primary" @click="selectFolder">选择文件夹</el-button>
    </div>
  </div>
</template>

<script setup>
const selectFolder = () => {
  alert('功能开发中')
}
</script>

<style scoped>
.local h2 { margin-bottom: 20px; }
.empty { text-align: center; padding: 100px 0; color: #999; }
.empty p { margin-bottom: 20px; }
</style>
