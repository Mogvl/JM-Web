import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  { path: '/', component: () => import('../views/Home.vue') },
  { path: '/search', component: () => import('../views/Search.vue') },
  { path: '/comic/:id', component: () => import('../views/Comic.vue') },
  { path: '/read/:comicId/:chapterId', component: () => import('../views/Reader.vue') },
  { path: '/favorites', component: () => import('../views/Favorites.vue') },
  { path: '/history', component: () => import('../views/History.vue') },
  { path: '/downloads', component: () => import('../views/Downloads.vue') }
]

export default createRouter({
  history: createWebHistory(),
  routes
})
