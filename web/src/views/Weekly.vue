<template>
  <div class="week-page">
    <!-- 顶部：周选择 + 更新提示 -->
    <div class="week-header">
      <el-select v-model="selectedWeek" @change="onWeekChange" placeholder="选择周期" class="week-select" :loading="initLoading">
        <el-option v-for="(cat, i) in weekCategories" :key="i" :label="cat.time + ' - ' + cat.title" :value="cat.id" />
      </el-select>
      <span class="update-tip">每周五 18:00更新</span>
    </div>

    <!-- 日漫/韩漫/其他 标签页 -->
    <el-tabs v-model="activeType" @tab-change="onTabChange" class="week-tabs">
      <el-tab-pane label="日漫" name="manga" />
      <el-tab-pane label="韩漫" name="hanman" />
      <el-tab-pane label="其他" name="another" />
    </el-tabs>

    <!-- 漫画列表 -->
    <div v-if="loading" class="loading-state">加载中...</div>
    <div v-else-if="items.length === 0" class="empty-state">暂无数据</div>
    <div v-else class="comic-grid">
      <div v-for="comic in items" :key="comic.id" class="comic-card" @click="$router.push(`/comic/${comic.id}`)">
        <img :src="comic.cover_url" :alt="comic.title" loading="lazy" />
        <div class="card-title">{{ comic.title }}</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getWeekCategories, getWeekFilter } from '../api'

const weekCategories = ref([])
const selectedWeek = ref('')
const activeType = ref('manga')
const items = ref([])
const loading = ref(false)
const initLoading = ref(true)
const ready = ref(false)

const loadData = async () => {
  if (!selectedWeek.value || !ready.value) return
  loading.value = true
  try {
    const data = await getWeekFilter(selectedWeek.value, activeType.value, 0)
    items.value = data.items || []
  } catch (e) { items.value = [] }
  finally { loading.value = false }
}

const onWeekChange = () => loadData()
const onTabChange = () => loadData()

onMounted(async () => {
  try {
    const cats = await getWeekCategories()
    weekCategories.value = cats || []
    if (weekCategories.value.length > 0) {
      selectedWeek.value = weekCategories.value[0].id
      ready.value = true
      await loadData()
    }
  } catch (e) { weekCategories.value = [] }
  finally { initLoading.value = false }
})
</script>

<style scoped>
.week-header { display: flex; align-items: center; gap: 16px; margin-bottom: 20px; }
.week-select { width: 360px; }
.update-tip { font-size: 13px; color: var(--text-muted); }
.week-tabs { margin-bottom: 20px; }
</style>
