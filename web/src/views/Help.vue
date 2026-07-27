<template>
  <div class="help">
    <h2>帮助</h2>
    <el-card class="help-card">
      <div class="version">版本: {{ version }}</div>
      <el-collapse v-model="activeName" accordion>
        <el-collapse-item v-for="(item, index) in faq" :key="index" :title="item.q" :name="index">
          <div>{{ item.a }}</div>
        </el-collapse-item>
      </el-collapse>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getHelp } from '../api'

const version = ref('1.0.0')
const faq = ref([])
const activeName = ref('0')

onMounted(async () => {
  try {
    const data = await getHelp()
    version.value = data.version || '1.0.0'
    faq.value = data.faq || []
  } catch (e) {
    faq.value = [
      { q: '如何搜索漫画？', a: '在顶部搜索框输入关键词，按回车即可搜索。' },
      { q: '如何收藏漫画？', a: '在漫画详情页点击收藏按钮。' },
      { q: '如何下载漫画？', a: '在漫画详情页点击下载按钮，下载任务会添加到下载管理中。' },
      { q: '如何切换阅读模式？', a: '在设置页面可以切换滚动/翻页模式。' }
    ]
  }
})
</script>

<style scoped>
.help h2 { margin-bottom: 20px; }
.help-card { max-width: 800px; }
.version { margin-bottom: 20px; color: #666; }
</style>
