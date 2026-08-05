<template>
  <div class="app-shell" v-if="!isLogin">
    <!-- 移动端遮罩 -->
    <div v-if="sidebarOpen" class="overlay" @click="sidebarOpen = false"></div>

    <!-- 左侧导航栏 -->
    <aside :class="['sidebar', sidebarOpen && 'open']">
      <div class="sidebar-close" @click="sidebarOpen = false"><el-icon><Close /></el-icon></div>
      <div class="user-section">
        <el-avatar :size="64" :src="userInfo.avatar || ''" class="user-avatar">
          {{ username.charAt(0).toUpperCase() }}
        </el-avatar>
        <div class="username">{{ username }}</div>
        <div class="user-stats" v-if="username !== '游客'">
          <span class="stat"><b>{{ userInfo.coins }}</b> J币</span>
          <span class="stat"><b>{{ userInfo.level_name || userInfo.level }}</b></span>
          <span class="stat"><b>{{ userInfo.favorites }}</b> 收藏</span>
        </div>
        <div class="user-actions">
          <el-button v-if="username !== '游客'" type="primary" size="small" round @click="handleSign">签到</el-button>
          <el-button v-if="username !== '游客'" type="primary" size="small" round plain @click="loadUserInfo">刷新信息</el-button>
          <el-button size="small" round plain @click="handleLogout">退出</el-button>
        </div>
      </div>

      <nav class="nav-menu" @click="onNavClick">
        <el-menu :default-active="$route.path" router>
          <el-menu-item-group>
            <template #title>用户</template>
            <el-menu-item index="/favorites"><el-icon><Star /></el-icon><span>收藏</span></el-menu-item>
            <el-menu-item index="/local-favorites"><el-icon><Folder /></el-icon><span>本地收藏</span></el-menu-item>
            <el-menu-item index="/history"><el-icon><Clock /></el-icon><span>历史</span></el-menu-item>
            <el-menu-item index="/remote-history"><el-icon><Upload /></el-icon><span>远程历史</span></el-menu-item>
            <el-menu-item index="/comments"><el-icon><ChatDotRound /></el-icon><span>我的评论</span></el-menu-item>
            <el-menu-item index="/all-comments"><el-icon><ChatLineSquare /></el-icon><span>全部评论</span></el-menu-item>
          </el-menu-item-group>
          <el-menu-item-group>
            <template #title>浏览</template>
            <el-menu-item index="/"><el-icon><HomeFilled /></el-icon><span>首页</span></el-menu-item>
            <el-menu-item index="/search"><el-icon><Search /></el-icon><span>搜索</span></el-menu-item>
            <el-menu-item index="/category"><el-icon><Grid /></el-icon><span>分类/排行</span></el-menu-item>
            <el-menu-item index="/weekly"><el-icon><Calendar /></el-icon><span>每周更新</span></el-menu-item>
          </el-menu-item-group>
          <el-menu-item-group>
            <template #title>工具</template>
            <el-menu-item index="/downloads"><el-icon><Download /></el-icon><span>下载管理</span></el-menu-item>
            <el-menu-item index="/download-all"><el-icon><List /></el-icon><span>批量下载</span></el-menu-item>
            <el-menu-item index="/local"><el-icon><FolderOpened /></el-icon><span>本地阅读</span></el-menu-item>
            <el-menu-item index="/nas"><el-icon><Connection /></el-icon><span>NAS</span></el-menu-item>
            <el-menu-item index="/waifu2x"><el-icon><MagicStick /></el-icon><span>图片增强</span></el-menu-item>
            <el-menu-item index="/batch-sr"><el-icon><FullScreen /></el-icon><span>批量超分</span></el-menu-item>
            <el-menu-item index="/settings"><el-icon><Setting /></el-icon><span>设置</span></el-menu-item>
            <el-menu-item index="/help"><el-icon><QuestionFilled /></el-icon><span>帮助</span></el-menu-item>
          </el-menu-item-group>
        </el-menu>
      </nav>
    </aside>

    <!-- 右侧内容区 -->
    <main class="main-area">
      <header class="top-bar">
        <button class="menu-toggle" @click="toggleSidebar">
          <el-icon><component :is="sidebarOpen ? Close : Menu" /></el-icon>
        </button>
        <el-breadcrumb separator="/" class="breadcrumb">
          <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
          <el-breadcrumb-item v-if="$route.meta.title">{{ $route.meta.title }}</el-breadcrumb-item>
        </el-breadcrumb>
        <div class="search-wrap">
          <el-input v-model="searchQuery" placeholder="搜索漫画..." @keyup.enter="doSearch" clearable>
            <template #append><el-button @click="doSearch"><el-icon><Search /></el-icon></el-button></template>
          </el-input>
        </div>
      </header>

      <div class="content-area">
        <router-view />
      </div>
    </main>
  </div>
  <router-view v-else />
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { Star, Clock, ChatDotRound, ChatLineSquare, HomeFilled, Search, Grid, Calendar, Download, FolderOpened, Setting, QuestionFilled, Connection, MagicStick, Upload, Folder, List, FullScreen, Menu, Close } from '@element-plus/icons-vue'
import { sign, getUserInfo } from './api'
import { ElMessage } from 'element-plus'

const router = useRouter()
const route = useRoute()
const searchQuery = ref('')
const sidebarOpen = ref(false)
const username = ref(localStorage.getItem('username') || '游客')
const userInfo = ref({
  coins: parseInt(localStorage.getItem('coins') || '0'),
  level: parseInt(localStorage.getItem('level') || '0'),
  level_name: localStorage.getItem('level_name') || '',
  avatar: localStorage.getItem('avatar') || '',
  favorites: parseInt(localStorage.getItem('favorites') || '0'),
})
const isLogin = computed(() => route.path === '/login')

// 刷新用户信息（对齐原版 SetUser/GetUserInfo）
const loadUserInfo = async () => {
  if (localStorage.getItem('token') === 'guest') return
  try {
    const info = await getUserInfo()
    if (info) {
      if (info.coin !== undefined) { userInfo.value.coins = info.coin; localStorage.setItem('coins', String(info.coin)) }
      if (info.level !== undefined) { userInfo.value.level = info.level; localStorage.setItem('level', String(info.level)) }
      if (info.photo) { userInfo.value.avatar = info.photo; localStorage.setItem('avatar', info.photo) }
      userInfo.value.level_name = info.level_name || userInfo.value.level_name
      localStorage.setItem('level_name', userInfo.value.level_name)
      if (info.album_favorites !== undefined) { userInfo.value.favorites = info.album_favorites; localStorage.setItem('favorites', String(info.album_favorites)) }
    }
  } catch (e) {}
}

const doSearch = () => {
  if (searchQuery.value.trim()) router.push({ path: '/search', query: { q: searchQuery.value.trim() } })
}

const toggleSidebar = () => { sidebarOpen.value = !sidebarOpen.value }
const onNavClick = () => { sidebarOpen.value = false }

const handleSign = async () => {
  try { await sign(); ElMessage.success('签到成功'); loadUserInfo() } catch (e) { ElMessage.error('签到失败') }
}

const handleLogout = () => { localStorage.clear(); router.push('/login') }

// 路由变化时关闭侧边栏
watch(() => route.path, () => { sidebarOpen.value = false })
</script>

<style scoped>
.user-section { padding: 24px 20px 20px; text-align: center; border-bottom: 1px solid var(--border); }
.user-avatar { background: var(--accent); color: #fff; font-weight: 600; }
.username { margin-top: 12px; font-size: 15px; font-weight: 600; color: var(--text-primary); }
.user-stats { margin-top: 8px; display: flex; justify-content: center; gap: 12px; font-size: 12px; color: var(--text-muted); }
.user-stats b { color: var(--text-secondary); }
.user-actions { margin-top: 14px; display: flex; justify-content: center; gap: 8px; }
.nav-menu { flex: 1; overflow-y: auto; padding-bottom: 16px; }
.sidebar-close { display: none; }
.top-bar { height: 56px; padding: 0 24px; display: flex; align-items: center; gap: 12px; background: var(--bg-surface); border-bottom: 1px solid var(--border); flex-shrink: 0; }
.menu-toggle { display: none; align-items: center; justify-content: center; width: 40px; height: 40px; border: none; background: transparent; color: var(--text-primary); cursor: pointer; border-radius: var(--radius-sm); font-size: 20px; }
.menu-toggle:hover { background: var(--bg-hover); }
.breadcrumb { flex-shrink: 0; }
.search-wrap { flex: 1; max-width: 320px; margin-left: auto; }
.content-area { flex: 1; overflow-y: auto; padding: 24px; }

@media (max-width: 768px) {
  .menu-toggle { display: inline-flex; }
  .sidebar-close { display: flex; align-items: center; justify-content: center; position: absolute; top: 12px; right: 12px; width: 36px; height: 36px; border-radius: 50%; background: var(--bg-hover); color: var(--text-secondary); cursor: pointer; z-index: 10; font-size: 18px; }
  .sidebar-close:hover { background: var(--border-light); color: var(--text-primary); }
  .user-section { padding-top: 32px; }
  .top-bar { padding: 0 12px; gap: 8px; }
  .breadcrumb { display: none; }
  .search-wrap { max-width: none; margin-left: 0; }
  .content-area { padding: 12px; }
}
</style>
