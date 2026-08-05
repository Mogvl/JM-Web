<template>
  <div class="search-page">
    <div class="search-bar">
      <el-tooltip placement="top" raw-content>
        <template #content>
          <div style="line-height:1.8;font-size:12px">
            <b>搜索语法:</b><br>包含搜索: 全彩 +人妻<br>排除搜索: 全彩 -人妻<br>普通搜索: 全彩 人妻
          </div>
        </template>
        <span class="hint-icon">?</span>
      </el-tooltip>
      <el-input v-model="query" placeholder="搜索漫画...支持 +、- 语法" @keyup.enter="doSearch" clearable size="large">
        <template #append><el-button @click="doSearch">搜索</el-button></template>
      </el-input>
    </div>

    <div v-if="searched" class="filter-bar">
      <span class="result-count">共 {{ items.length }} 条结果</span>
    </div>

    <div v-if="loading" class="loading-state">加载中...</div>
    <div v-else-if="searched && items.length === 0" class="empty-state">没有找到相关漫画</div>
    <div v-else class="comic-grid">
      <div v-for="comic in items" :key="comic.id" class="comic-card" @click="$router.push(`/comic/${comic.id}`)">
        <img :src="comic.cover_url" :alt="comic.title" loading="lazy" />
        <div class="card-title">{{ comic.title }}</div>
        <div v-if="comic.author" class="card-author">{{ comic.author }}</div>
      </div>
    </div>

    <div v-if="totalPages > 1" class="pager">
      <el-pagination v-model:current-page="page" :total="totalPages * 20" :page-size="20" layout="prev, pager, next" @current-change="loadPage" />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { search } from '../api'

const route = useRoute()
const router = useRouter()
const query = ref('')
const items = ref([])
const loading = ref(false)
const searched = ref(false)
const page = ref(parseInt(route.query.p || '1'))
const totalPages = ref(1)

const doSearch = async () => {
  if (!query.value.trim()) return
  searched.value = true
  page.value = 1
  loading.value = true
  try {
    const data = await search(query.value, page.value)
    items.value = data.items || []
    totalPages.value = data.total_pages || 1
  } catch (e) { items.value = [] }
  finally { loading.value = false }
}

const loadPage = async (p = page.value) => {
  page.value = p
  // 页码同步到 URL，返回时能恢复第几页
  router.replace({ query: { ...route.query, p: p > 1 ? String(p) : undefined } })
  loading.value = true
  try { const data = await search(query.value, p); items.value = data.items || [] }
  catch (e) { items.value = [] }
  finally { loading.value = false }
}

onMounted(() => {
  if (route.query.q) { query.value = route.query.q; searched.value = true; loadPage(parseInt(route.query.p || '1')) }
})

// 返回本页时若 URL 仍带页码，恢复对应页
watch(() => route.query.p, (p) => {
  const target = parseInt(p || '1')
  if (target !== page.value && route.query.q) { page.value = target; loadPage(target) }
})
</script>

<style scoped>
.search-bar { display: flex; align-items: center; gap: 8px; margin-bottom: 20px; max-width: 600px; }
.hint-icon { display: inline-flex; align-items: center; justify-content: center; width: 22px; height: 22px; border-radius: 50%; border: 1px solid var(--text-muted); color: var(--text-muted); font-size: 12px; cursor: help; flex-shrink: 0; }
.filter-bar { margin-bottom: 20px; }
.result-count { font-size: 13px; color: var(--text-muted); }
.pager { margin-top: 24px; display: flex; justify-content: center; }
</style>
