<template>
  <el-container class="app">
    <!-- 左侧导航栏 -->
    <el-aside width="240px" class="sidebar">
      <div class="user-section">
        <el-avatar :size="80" src="/avatar.png" />
        <el-button type="primary" size="small" class="login-btn">登录</el-button>
        <div class="user-info">
          <div class="username">游客</div>
          <div class="coins">J Coins: 0</div>
          <div class="level">等级: 0</div>
        </div>
      </div>

      <el-menu :default-active="$route.path" router class="nav-menu">
        <el-menu-item-group>
          <template #title>用户</template>
          <el-menu-item index="/favorites">
            <el-icon><Star /></el-icon>
            <span>我的收藏</span>
          </el-menu-item>
          <el-menu-item index="/history">
            <el-icon><Clock /></el-icon>
            <span>阅读历史</span>
          </el-menu-item>
        </el-menu-item-group>

        <el-menu-item-group>
          <template #title>浏览</template>
          <el-menu-item index="/">
            <el-icon><HomeFilled /></el-icon>
            <span>首页</span>
          </el-menu-item>
          <el-menu-item index="/search">
            <el-icon><Search /></el-icon>
            <span>搜索</span>
          </el-menu-item>
          <el-menu-item index="/category">
            <el-icon><Grid /></el-icon>
            <span>分类</span>
          </el-menu-item>
          <el-menu-item index="/rank">
            <el-icon><Trophy /></el-icon>
            <span>排行榜</span>
          </el-menu-item>
        </el-menu-item-group>

        <el-menu-item-group>
          <template #title>工具</template>
          <el-menu-item index="/downloads">
            <el-icon><Download /></el-icon>
            <span>下载管理</span>
          </el-menu-item>
          <el-menu-item index="/settings">
            <el-icon><Setting /></el-icon>
            <span>设置</span>
          </el-menu-item>
        </el-menu-item-group>
      </el-menu>
    </el-aside>

    <!-- 右侧内容区 -->
    <el-container>
      <el-header class="content-header">
        <div class="breadcrumb">
          <el-breadcrumb separator="/">
            <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
            <el-breadcrumb-item v-if="$route.meta.title">{{ $route.meta.title }}</el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        <div class="search-box">
          <el-input v-model="searchQuery" placeholder="搜索漫画..." @keyup.enter="doSearch" clearable>
            <template #append>
              <el-button @click="doSearch">
                <el-icon><Search /></el-icon>
              </el-button>
            </template>
          </el-input>
        </div>
      </el-header>

      <el-main class="content-main">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { Star, Clock, HomeFilled, Search, Grid, Trophy, Download, Setting } from '@element-plus/icons-vue'

const router = useRouter()
const searchQuery = ref('')

const doSearch = () => {
  if (searchQuery.value.trim()) {
    router.push({ path: '/search', query: { q: searchQuery.value.trim() } })
  }
}
</script>

<style>
* { margin: 0; padding: 0; box-sizing: border-box; }
html, body, #app { height: 100%; }
body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background: #f5f5f5; }

.app { height: 100vh; }

.sidebar {
  background: #fff;
  border-right: 1px solid #e6e6e6;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
}

.user-section {
  padding: 20px;
  text-align: center;
  border-bottom: 1px solid #e6e6e6;
}

.login-btn { margin-top: 10px; }

.user-info { margin-top: 10px; font-size: 12px; color: #666; }
.user-info div { margin: 4px 0; }

.nav-menu {
  flex: 1;
  border-right: none;
}

.nav-menu .el-menu-item {
  height: 45px;
  line-height: 45px;
}

.content-header {
  background: #fff;
  border-bottom: 1px solid #e6e6e6;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
}

.search-box { width: 300px; }

.content-main {
  padding: 20px;
  overflow-y: auto;
}
</style>
