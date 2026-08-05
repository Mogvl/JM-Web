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
      <el-autocomplete
        v-model="query"
        :fetch-suggestions="querySuggestions"
        placeholder="搜索漫画...支持 +、- 语法"
        size="large"
        clearable
        @keyup.enter="doSearch"
        @select="(item) => { query = item.value; doSearch() }"
      >
        <template #append><el-button @click="doSearch">搜索</el-button></template>
      </el-autocomplete>
    </div>

    <div v-if="searched" class="filter-bar">
      <span class="result-count">共 {{ items.length }} 条结果</span>
    </div>

    <div v-if="loading" class="loading-state">加载中...</div>
    <div v-else-if="searched && items.length === 0" class="empty-state">没有找到相关漫画</div>
    <div v-else-if="searched" class="comic-grid">
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
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { search } from '../api'

const route = useRoute()
const router = useRouter()
const query = ref('')
const items = ref([])
const loading = ref(false)
const searched = ref(false)
const page = ref(1)
const totalPages = ref(1)

// 搜索历史缓存（对齐原版 CacheWord）
const SEARCH_CACHE_KEY = 'searchCache'
const getSearchCache = () => JSON.parse(localStorage.getItem(SEARCH_CACHE_KEY) || '[]')
const addSearchCache = (word) => {
  const list = getSearchCache().filter(w => w !== word)
  list.unshift(word)
  localStorage.setItem(SEARCH_CACHE_KEY, JSON.stringify(list.slice(0, 10)))
}
const querySuggestions = (q, cb) => {
  const words = getSearchCache().filter(w => w.includes(q))
  cb(words.map(w => ({ value: w })))
}

const loadResults = async () => {
  if (!query.value.trim()) return
  searched.value = true
  loading.value = true
  try {
    const data = await search(query.value, page.value)
    items.value = data.items || []
    totalPages.value = data.total_pages || 1
  } catch (e) { items.value = [] }
  finally { loading.value = false }
}

// 把 q 和 p 都写入 URL，作为详情页返回时的恢复依据
const doSearch = async () => {
  if (!query.value.trim()) return
  page.value = 1
  addSearchCache(query.value.trim())
  router.replace({ query: { q: query.value.trim(), p: undefined } })
  await loadResults()
}

const loadPage = async (p = page.value) => {
  page.value = p
  router.replace({ query: { q: query.value.trim(), p: p > 1 ? String(p) : undefined } })
  await loadResults()
}

onMounted(() => {
  if (route.query.q) {
    query.value = route.query.q
    page.value = parseInt(route.query.p || '1')
    loadResults()
  }
})
</script>

<style scoped>
.search-bar { display: flex; align-items: center; gap: 8px; margin-bottom: 20px; max-width: 600px; }
.hint-icon { display: inline-flex; align-items: center; justify-content: center; width: 22px; height: 22px; border-radius: 50%; border: 1px solid var(--text-muted); color: var(--text-muted); font-size: 12px; cursor: help; flex-shrink: 0; }
.filter-bar { margin-bottom: 20px; }
.result-count { font-size: 13px; color: var(--text-muted); }
.pager { margin-top: 24px; display: flex; justify-content: center; }
</style>
