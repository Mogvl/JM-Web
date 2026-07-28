<template>
  <div class="category-page">
    <!-- 分类标签页 -->
    <el-tabs v-model="activeTab" @tab-change="onTabChange" class="cat-tabs">
      <el-tab-pane label="标题" name="titles">
        <div class="title-section">
          <div v-for="block in blocks" :key="block.title" class="title-group">
            <h4 class="group-title">{{ block.title }}</h4>
            <div class="tag-list">
              <el-tag v-for="tag in block.content" :key="tag" @click="searchTag(tag)" class="title-tag" effect="plain" size="large">{{ tag }}</el-tag>
            </div>
          </div>
          <div v-if="blocks.length === 0" class="empty-state">暂无数据</div>
        </div>
      </el-tab-pane>
      <el-tab-pane v-for="cat in categories" :key="cat.id" :label="cat.name" :name="cat.id" />
    </el-tabs>

    <!-- 分类内容 -->
    <template v-if="activeTab !== 'titles'">
      <!-- 子分类筛选 -->
      <div v-if="currentSubCats.length" class="subcat-bar">
        <button :class="['subcat-chip', !activeSubCat && 'active']" @click="selectSubCat('')">全部</button>
        <button v-for="sub in currentSubCats" :key="sub.slug" :class="['subcat-chip', activeSubCat === sub.slug && 'active']" @click="selectSubCat(sub.slug)">{{ sub.name }}</button>
      </div>

      <div v-if="loading" class="loading-state">加载中...</div>
      <div v-else-if="items.length === 0" class="empty-state">暂无数据</div>
      <div v-else class="comic-grid">
        <div v-for="comic in items" :key="comic.id" class="comic-card" @click="$router.push(`/comic/${comic.id}`)">
          <img :src="comic.cover_url" :alt="comic.title" loading="lazy" />
          <div class="card-title">{{ comic.title }}</div>
        </div>
      </div>

      <!-- 底部排序和分页 -->
      <div class="bottom-bar">
        <el-select v-model="sort" @change="loadData" class="sort-select">
          <el-option label="最新" value="mr" />
          <el-option label="总排行" value="mv" />
          <el-option label="月排行" value="mv_m" />
          <el-option label="周排行" value="mv_w" />
          <el-option label="日排行" value="mv_t" />
          <el-option label="最多图片" value="mp" />
          <el-option label="最多爱心" value="tf" />
        </el-select>
        <div class="spacer"></div>
        <span class="page-info">页：{{ page }}/{{ totalPages || 1 }}</span>
        <el-input-number v-model="jumpPage" :min="1" :max="totalPages||1" size="small" class="page-input" />
        <el-button size="small" @click="jump" class="jump-btn">跳转</el-button>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getCategories, getCategoryFilter } from '../api'

const router = useRouter()
const categories = ref([])
const blocks = ref([])
const items = ref([])
const loading = ref(false)
const activeTab = ref('titles')
const activeSubCat = ref('')
const sort = ref('mr')
const page = ref(1)
const totalPages = ref(1)
const jumpPage = ref(1)

const currentSubCats = computed(() => {
  const cat = categories.value.find(c => c.id === activeTab.value)
  return cat?.sub_categories || []
})

const loadData = async () => {
  if (activeTab.value === 'titles') return
  loading.value = true
  try {
    const cat = activeSubCat.value || activeTab.value
    const data = await getCategoryFilter(cat, sort.value, page.value)
    items.value = data.items || []
    totalPages.value = data.total_pages || 1
  } catch (e) { items.value = [] }
  finally { loading.value = false }
}

const onTabChange = (tab) => {
  if (tab === 'titles') return
  activeSubCat.value = ''
  page.value = 1
  jumpPage.value = 1
  loadData()
}

const selectSubCat = (slug) => {
  activeSubCat.value = slug
  page.value = 1
  jumpPage.value = 1
  loadData()
}

const searchTag = (tag) => {
  router.push({ path: '/search', query: { q: tag } })
}

const jump = () => {
  page.value = jumpPage.value
  loadData()
}

onMounted(async () => {
  try {
    const data = await getCategories()
    categories.value = data.categories || []
    blocks.value = data.blocks || []
  } catch (e) { categories.value = []; blocks.value = [] }
})
</script>

<style scoped>
.cat-tabs { margin-bottom: 0; }
.cat-tabs :deep(.el-tabs__content) { padding-top: 8px; }

.title-section { padding: 8px 0; }
.title-group { margin-bottom: 24px; }
.group-title { font-size: 14px; font-weight: 600; color: var(--text-secondary); margin-bottom: 12px; }
.tag-list { display: flex; flex-wrap: wrap; gap: 10px; }
.title-tag { cursor: pointer; font-size: 14px; padding: 8px 16px; border-radius: 20px; }
.title-tag:hover { color: var(--accent); border-color: var(--accent); }

.subcat-bar { display: flex; flex-wrap: wrap; gap: 8px; margin-bottom: 20px; padding: 12px 16px; background: var(--bg-surface); border: 1px solid var(--border); border-radius: var(--radius-md); }
.subcat-chip { padding: 4px 14px; border-radius: 16px; border: 1px solid var(--border); background: var(--bg-surface); color: var(--text-secondary); font-size: 13px; cursor: pointer; transition: var(--transition); }
.subcat-chip:hover { border-color: var(--border-light); color: var(--text-primary); }
.subcat-chip.active { background: var(--accent); border-color: var(--accent); color: #fff; font-weight: 500; }

.bottom-bar { display: flex; align-items: center; gap: 12px; margin-top: 24px; padding: 16px 0; border-top: 1px solid var(--border); }
.sort-select { width: 120px; }
.spacer { flex: 1; }
.page-info { font-size: 13px; color: var(--text-muted); white-space: nowrap; }
.page-input { width: 80px; }
.jump-btn { min-width: 60px; }
</style>
