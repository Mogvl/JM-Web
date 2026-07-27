<template>
  <div class="home">
    <h2>欢迎使用 JMComic Web</h2>
    <p>搜索你喜欢的漫画，开始阅读吧！</p>
    <div class="quick-links">
      <el-card @click="$router.push('/favorites')">⭐ 收藏</el-card>
      <el-card @click="$router.push('/history')">📜 历史</el-card>
      <el-card @click="$router.push('/downloads')">⬇️ 下载</el-card>
      <el-card @click="$router.push('/category')">📂 分类</el-card>
    </div>
  </div>
</template>

<style scoped>
.home { text-align: center; padding: 60px 0; }
.home h2 { font-size: 32px; margin-bottom: 10px; }
.home p { color: #666; margin-bottom: 40px; }
.quick-links { display: flex; justify-content: center; gap: 30px; }
.quick-links .el-card { width: 150px; cursor: pointer; text-align: center; font-size: 18px; }
.quick-links .el-card:hover { transform: translateY(-5px); box-shadow: 0 4px 12px rgba(0,0,0,0.15); }
</style>
