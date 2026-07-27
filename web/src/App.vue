<template>
  <el-container class="app">
    <el-header class="header">
      <div class="logo" @click="$router.push('/')">JMComic</div>
      <el-menu mode="horizontal" :router="true" :default-active="$route.path" class="nav">
        <el-menu-item index="/">首页</el-menu-item>
        <el-menu-item index="/favorites">收藏</el-menu-item>
        <el-menu-item index="/history">历史</el-menu-item>
        <el-menu-item index="/downloads">下载</el-menu-item>
      </el-menu>
      <el-input v-model="query" placeholder="搜索漫画..." @keyup.enter="search" class="search">
        <template #append>
          <el-button @click="search">搜索</el-button>
        </template>
      </el-input>
    </el-header>
    <el-main class="main">
      <router-view />
    </el-main>
  </el-container>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const query = ref('')

const search = () => {
  if (query.value.trim()) {
    router.push({ path: '/search', query: { q: query.value.trim() } })
  }
}
</script>

<style>
* { margin: 0; padding: 0; box-sizing: border-box; }
body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background: #f5f5f5; }
.app { min-height: 100vh; }
.header { display: flex; align-items: center; background: #fff; box-shadow: 0 2px 4px rgba(0,0,0,0.1); padding: 0 20px; }
.logo { font-size: 24px; font-weight: bold; color: #409eff; cursor: pointer; margin-right: 40px; }
.nav { flex: 1; }
.search { width: 300px; }
.main { padding: 20px; max-width: 1400px; margin: 0 auto; }
</style>
