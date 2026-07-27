<template>
  <el-container class="app">
    <!-- 左侧导航栏 -->
    <el-aside width="240px" class="sidebar">
      <div class="user-section">
        <el-avatar :size="80" :src="userInfo.avatar || '/avatar.png'" />
        <div class="user-info">
          <div class="username">{{ username }}</div>
          <div class="coins">J Coins: {{ userInfo.coins || 0 }}</div>
          <div class="level">等级: {{ userInfo.level || 0 }}</div>
          <div class="favorites">收藏数: {{ userInfo.favorites || 0 }}</div>
        </div>
        <div class="user-actions">
          <el-button v-if="username !== '游客'" type="success" size="small" @click="handleSign">签到</el-button>
          <el-button type="danger" size="small" @click="handleLogout">退出</el-button>
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
          <el-menu-item index="/comments">
            <el-icon><ChatDotRound /></el-icon>
            <span>我的评论</span>
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
          <el-menu-item index="/weekly">
            <el-icon><Calendar /></el-icon>
            <span>每周更新</span>
          </el-menu-item>
        </el-menu-item-group>

        <el-menu-item-group>
          <template #title>工具</template>
          <el-menu-item index="/downloads">
            <el-icon><Download /></el-icon>
            <span>下载管理</span>
          </el-menu-item>
          <el-menu-item index="/local">
            <el-icon><FolderOpened /></el-icon>
            <span>本地阅读</span>
          </el-menu-item>
          <el-menu-item index="/nas">
            <el-icon><Connection /></el-icon>
            <span>NAS 同步</span>
          </el-menu-item>
          <el-menu-item index="/waifu2x">
            <el-icon><MagicStick /></el-icon>
            <span>图片增强</span>
          </el-menu-item>
          <el-menu-item index="/settings">
            <el-icon><Setting /></el-icon>
            <span>设置</span>
          </el-menu-item>
          <el-menu-item index="/help">
            <el-icon><QuestionFilled /></el-icon>
            <span>帮助</span>
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
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Star, Clock, ChatDotRound, HomeFilled, Search, Grid, Trophy, Calendar, Download, FolderOpened, Setting, QuestionFilled, Connection, MagicStick } from '@element-plus/icons-vue'
import { getUserInfo, sign } from './api'
import { ElMessage } from 'element-plus'

const router = useRouter()
const searchQuery = ref('')
const username = ref(localStorage.getItem('username') || '游客')
const userInfo = ref({})

const doSearch = () => {
  if (searchQuery.value.trim()) {
    router.push({ path: '/search', query: { q: searchQuery.value.trim() } })
  }
}

const handleSign = async () => {
  try {
    await sign()
    ElMessage.success('签到成功')
    loadUserInfo()
  } catch (e) {
    ElMessage.error('签到失败')
  }
}

const handleLogout = () => {
  localStorage.removeItem('token')
  localStorage.removeItem('username')
  router.push('/login')
}

const loadUserInfo = async () => {
  try {
    const data = await getUserInfo()
    userInfo.value = data
  } catch (e) {
    // 未登录或游客
  }
}

onMounted(() => {
  if (username.value !== '游客') {
    loadUserInfo()
  }
})
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

.user-info { margin-top: 10px; font-size: 12px; color: #666; }
.user-info div { margin: 4px 0; }

.user-actions { margin-top: 10px; display: flex; justify-content: center; gap: 8px; }

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
