<template>
  <div class="back-nav">
    <button class="back-btn" @click="goBack">
      <el-icon><ArrowLeft /></el-icon><span>返回</span>
    </button>
    <span class="nav-title">{{ title }}</span>
    <slot name="right" />
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { ArrowLeft } from '@element-plus/icons-vue'

const props = defineProps({
  title: { type: String, default: '' },
  fallback: { type: String, default: '/search' }
})

const router = useRouter()
const goBack = () => {
  if (window.history.length > 1) router.back()
  else router.push(props.fallback)
}
</script>

<style scoped>
.back-nav { display: flex; align-items: center; gap: 16px; padding: 10px 24px; background: var(--bg-surface); border: 1px solid var(--border); border-radius: var(--radius-md); margin-bottom: 16px; }
.back-btn { display: inline-flex; align-items: center; gap: 4px; padding: 8px 16px; border: 1px solid var(--border); background: transparent; color: var(--text-primary); border-radius: var(--radius-sm); cursor: pointer; transition: var(--transition); font-size: 14px; font-weight: 500; font-family: inherit; flex-shrink: 0; }
.back-btn:hover { border-color: var(--accent); color: var(--accent); background: var(--accent-soft); }
.nav-title { font-size: 15px; font-weight: 600; color: var(--text-primary); flex: 1; min-width: 0; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }

@media (max-width: 768px) {
  .back-nav { padding: 8px 12px; gap: 8px; position: sticky; top: 0; z-index: 20; border-radius: var(--radius-sm); }
  .back-btn { padding: 7px 12px; }
}
</style>