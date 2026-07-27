import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  { path: '/', component: () => import('../views/Home.vue'), meta: { title: '首页' } },
  { path: '/search', component: () => import('../views/Search.vue'), meta: { title: '搜索' } },
  { path: '/comic/:id', component: () => import('../views/Comic.vue'), meta: { title: '漫画详情' } },
  { path: '/read/:comicId/:chapterId', component: () => import('../views/Reader.vue'), meta: { title: '阅读' } },
  { path: '/favorites', component: () => import('../views/Favorites.vue'), meta: { title: '我的收藏' } },
  { path: '/history', component: () => import('../views/History.vue'), meta: { title: '阅读历史' } },
  { path: '/downloads', component: () => import('../views/Downloads.vue'), meta: { title: '下载管理' } },
  { path: '/category', component: () => import('../views/Category.vue'), meta: { title: '分类' } },
  { path: '/rank', component: () => import('../views/Rank.vue'), meta: { title: '排行榜' } },
  { path: '/settings', component: () => import('../views/Settings.vue'), meta: { title: '设置' } }
]

export default createRouter({
  history: createWebHistory(),
  routes
})
