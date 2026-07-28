<template>
  <div class="search">
    <div class="search-bar">
      <el-tooltip placement="top" raw-content>
        <template #content>
          <div style="line-height:1.8;font-size:12px">
            <b>搜索语法:</b><br>
            包含搜索: 全彩 +人妻<br>
            排除搜索: 全彩 -人妻<br>
            普通搜索: 全彩 人妻
          </div>
        </template>
        <span class="search-hint">?</span>
      </el-tooltip>
      <el-input v-model="query" placeholder="搜索漫画...支持 +、- 语法" @keyup.enter="doSearch" clearable class="search-input">
        <template #append>
          <el-button @click="doSearch">搜索</el-button>
        </template>
      </el-input>
    </div>

    <div v-if="searched" class="controls">
      <el-select v-model="sort" @change="doSearch" style="width:110px">
        <el-option label="最新" value="mr" />
        <el-option label="最多点击" value="mv" />
        <el-option label="最多图片" value="mp" />
        <el-option label="最多爱心" value="tf" />
      </el-select>
      <span class="page-info">页：{{ page }}/{{ totalPages || 1 }}</span>
      <el-input-number v-model="jumpPage" :min="1" :max="totalPages||1" size="small" style="width:70px" />
      <el-button size="small" @click="jump">跳转</el-button>
    </div>

    <div v-if="loading" class="loading">加载中...</div>
    <div v-else-if="searched && items.length === 0" class="empty">没有找到漫画</div>
    <div v-else class="grid">
      <div v-for="comic in items" :key="comic.id" class="comic-card" @click="$router.push(`/comic/${comic.id}`)">
        <img :src="comic.cover_url" :alt="comic.title" loading="lazy" />
        <div class="title">{{ comic.title }}</div>
        <div v-if="comic.author" class="author">{{ comic.author }}</div>
      </div>
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
const sort = ref('mr')
const page = ref(1)
const totalPages = ref(1)
const jumpPage = ref(1)

const doSearch = async () => {
  if (!query.value.trim()) return
  searched.value = true
  page.value = 1
  jumpPage.value = 1
  loading.value = true
  try {
    const data = await search(query.value, page.value, sort.value)
    items.value = data.items || []
    totalPages.value = data.total_pages || 1
  } catch (e) { items.value = [] }
  finally { loading.value = false }
}

const jump = () => {
  page.value = jumpPage.value
  doSearch()
}

onMounted(() => {
  if (route.query.q) {
    query.value = route.query.q
    doSearch()
  }
})
</script>

<style scoped>
.search-bar { display: flex; align-items: center; gap: 8px; margin-bottom: 15px; }
.search-input { flex: 1; max-width: 500px; }
.search-hint { display: inline-flex; align-items: center; justify-content: center; width: 20px; height: 20px; border-radius: 50%; border: 1px solid #409eff; color: #409eff; font-size: 12px; cursor: help; }
.controls { display: flex; align-items: center; gap: 12px; margin-bottom: 20px; }
.page-info { font-size: 14px; color: #666; }
.loading, .empty { text-align: center; padding: 80px 0; color: #999; }
.grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(160px, 1fr)); gap: 16px; }
.comic-card { cursor: pointer; border-radius: 8px; overflow: hidden; background: #fff; box-shadow: 0 1px 4px rgba(0,0,0,0.08); transition: all 0.3s; }
.comic-card:hover { transform: translateY(-4px); box-shadow: 0 4px 12px rgba(0,0,0,0.15); }
.comic-card img { width: 100%; height: 220px; object-fit: cover; display: block; }
.comic-card .title { padding: 10px 10px 2px; font-size: 13px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.comic-card .author { padding: 0 10px 10px; font-size: 12px; color: #999; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
</style>
