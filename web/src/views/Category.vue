<template>
  <div class="category">
    <h2>分类浏览</h2>
    <div v-if="loading" class="loading">加载中...</div>
    <div v-else class="tags">
      <el-tag v-for="cat in categories" :key="cat.id" @click="goSearch(cat.name)" class="tag" effect="plain" size="large">
        {{ cat.name }}
      </el-tag>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getCategories } from '../api'

const router = useRouter()
const categories = ref([])
const loading = ref(true)

const goSearch = (name) => {
  router.push({ path: '/search', query: { q: name } })
}

onMounted(async () => {
  try {
    categories.value = await getCategories()
  } catch (e) {
    categories.value = []
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.category h2 { margin-bottom: 20px; }
.loading { text-align: center; padding: 100px 0; color: #999; }
.tags { display: flex; flex-wrap: wrap; gap: 12px; }
.tag { cursor: pointer; font-size: 16px; padding: 10px 20px; }
.tag:hover { color: #409eff; border-color: #409eff; }
</style>
