<template>
  <div class="home">
    <h2>JMComic Web</h2>
    <p>搜索你喜欢的漫画</p>
    <div class="links">
      <el-card @click="$router.push('/favorites')">⭐ 收藏</el-card>
      <el-card @click="$router.push('/history')">📜 历史</el-card>
      <el-card @click="$router.push('/downloads')">⬇️ 下载</el-card>
    </div>
  </div>
</template>

<style scoped>
.home { text-align: center; padding: 60px 0; }
.home h2 { font-size: 32px; margin-bottom: 10px; }
.home p { color: #666; margin-bottom: 40px; }
.links { display: flex; justify-content: center; gap: 30px; }
.links .el-card { width: 150px; cursor: pointer; text-align: center; font-size: 18px; }
.links .el-card:hover { transform: translateY(-5px); box-shadow: 0 4px 12px rgba(0,0,0,0.15); }
</style>
