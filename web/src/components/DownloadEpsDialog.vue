<template>
  <el-dialog v-model="visible" title="章节下载" width="500px" @open="onOpen">
    <div class="eps-title">{{ comic?.title }}</div>
    <div class="eps-toolbar">
      <span class="label">章节</span>
      <div class="spacer"></div>
      <el-button size="small" @click="selectAll">全选</el-button>
      <el-button size="small" @click="invertSelect">反选</el-button>
    </div>
    <div class="eps-list">
      <el-checkbox v-for="ch in chapters" :key="ch.id" v-model="selected[ch.id]" class="eps-item">
        {{ ch.title || `第${ch.sort_order}话` }}
      </el-checkbox>
      <div v-if="chapters.length === 0" class="empty">暂无章节</div>
    </div>
    <template #footer>
      <el-button @click="visible = false">取消</el-button>
      <el-button type="primary" :loading="downloading" @click="startDownload">下载</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { getChapters, createDownload } from '../api'
import { ElMessage } from 'element-plus'

const props = defineProps({ modelValue: Boolean, comicId: String, comic: Object })
const emit = defineEmits(['update:modelValue'])
const router = useRouter()

const visible = ref(props.modelValue)
const chapters = ref([])
const selected = ref({})
const downloading = ref(false)

watch(() => props.modelValue, v => { visible.value = v })
watch(visible, v => emit('update:modelValue', v))

const onOpen = async () => {
  chapters.value = []
  selected.value = {}
  if (!props.comicId) return
  try {
    const data = await getChapters(props.comicId)
    chapters.value = data.chapters || []
    chapters.value.forEach(ch => { selected.value[ch.id] = true })
  } catch (e) {}
}

const selectAll = () => { chapters.value.forEach(ch => { selected.value[ch.id] = true }) }
const invertSelect = () => { chapters.value.forEach(ch => { selected.value[ch.id] = !selected.value[ch.id] }) }

const startDownload = async () => {
  const ids = chapters.value.filter(ch => selected.value[ch.id]).map(ch => ch.id)
  if (ids.length === 0) { ElMessage.warning('请选择章节'); return }
  downloading.value = true
  try {
    await createDownload(props.comicId, {
      title: props.comic?.title,
      author: props.comic?.author,
      cover: props.comic?.cover_url,
      chapters: ids,
    })
    ElMessage.success(`已添加 ${ids.length} 个章节到下载`)
    visible.value = false
  } catch (e) { ElMessage.error('下载失败') }
  finally { downloading.value = false }
}
</script>

<style scoped>
.eps-title { font-size: 15px; font-weight: 600; margin-bottom: 12px; color: var(--text-primary); }
.eps-toolbar { display: flex; align-items: center; gap: 8px; margin-bottom: 12px; }
.eps-toolbar .label { font-size: 13px; color: var(--text-muted); }
.spacer { flex: 1; }
.eps-list { max-height: 360px; overflow-y: auto; display: flex; flex-direction: column; gap: 6px; }
.eps-item { margin-right: 0; }
.empty { text-align: center; color: var(--text-muted); padding: 24px; }
</style>
