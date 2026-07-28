<template>
  <div class="settings">
    <h2>设置</h2>
    <el-tabs v-model="activeTab">
      <el-tab-pane label="通用" name="general">
        <el-card class="card">
          <el-form label-width="120px">
            <el-form-item label="界面语言">
              <el-select v-model="lang" style="width: 200px">
                <el-option label="简体中文" value="zh-CN" />
                <el-option label="English" value="en" />
              </el-select>
            </el-form-item>
            <el-form-item label="主题">
              <el-radio-group v-model="theme">
                <el-radio label="light">浅色</el-radio>
                <el-radio label="dark">深色</el-radio>
                <el-radio label="auto">跟随系统</el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item label="API 地址">
              <el-input v-model="apiUrl" placeholder="https://www.cdnhjk.net" />
            </el-form-item>
            <el-form-item label="请求超时(秒)">
              <el-input-number v-model="timeout" :min="3" :max="60" />
            </el-form-item>
          </el-form>
        </el-card>
      </el-tab-pane>

      <el-tab-pane label="代理" name="proxy">
        <el-card class="card">
          <el-form label-width="120px">
            <el-form-item label="代理方式">
              <el-radio-group v-model="proxyType">
                <el-radio label="none">无代理</el-radio>
                <el-radio label="http">HTTP</el-radio>
                <el-radio label="socks5">SOCKS5</el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item label="代理地址">
              <el-input v-model="proxyHost" placeholder="127.0.0.1" />
            </el-form-item>
            <el-form-item label="代理端口">
              <el-input-number v-model="proxyPort" :min="1" :max="65535" />
            </el-form-item>
            <el-form-item label="DoH DNS">
              <el-input v-model="dohUrl" placeholder="https://dns.alidns.com/dns-query" />
            </el-form-item>
          </el-form>
        </el-card>
      </el-tab-pane>

      <el-tab-pane label="Waifu2x" name="waifu2x">
        <el-card class="card">
          <el-form label-width="120px">
            <el-form-item label="模型">
              <el-select v-model="waifu2xModel" style="width: 200px">
                <el-option label="CUNET" value="cunet" />
                <el-option label="UPSCALE" value="upscale" />
                <el-option label="NOISE" value="noise" />
              </el-select>
            </el-form-item>
            <el-form-item label="输出格式">
              <el-select v-model="outputFormat" style="width: 200px">
                <el-option label="PNG" value="png" />
                <el-option label="WEBP" value="webp" />
                <el-option label="JPG" value="jpg" />
              </el-select>
            </el-form-item>
            <el-form-item label="使用GPU">
              <el-switch v-model="useGpu" />
            </el-form-item>
          </el-form>
        </el-card>
      </el-tab-pane>

      <el-tab-pane label="下载与缓存" name="download">
        <el-card class="card">
          <el-form label-width="120px">
            <el-form-item label="下载目录">
              <el-input v-model="downloadDir" placeholder="/data/downloads" />
            </el-form-item>
            <el-form-item label="同时下载数">
              <el-input-number v-model="maxDownloads" :min="1" :max="10" />
            </el-form-item>
            <el-form-item label="图片质量">
              <el-select v-model="imageQuality" style="width: 200px">
                <el-option label="原图" value="original" />
                <el-option label="高清" value="high" />
                <el-option label="普通" value="normal" />
              </el-select>
            </el-form-item>
            <el-form-item label="缓存大小">
              <el-select v-model="cacheSize" style="width: 200px">
                <el-option label="256MB" value="256" />
                <el-option label="512MB" value="512" />
                <el-option label="1GB" value="1024" />
                <el-option label="2GB" value="2048" />
              </el-select>
            </el-form-item>
            <el-form-item label="自动缓存">
              <el-switch v-model="autoCache" />
            </el-form-item>
          </el-form>
        </el-card>
      </el-tab-pane>
    </el-tabs>

    <div class="actions">
      <el-button type="primary" @click="save">保存</el-button>
      <el-button @click="reset">重置</el-button>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { ElMessage } from 'element-plus'

const activeTab = ref('general')
const lang = ref(localStorage.getItem('lang') || 'zh-CN')
const theme = ref(localStorage.getItem('theme') || 'light')
const apiUrl = ref(localStorage.getItem('apiUrl') || 'https://www.cdnhjk.net')
const timeout = ref(parseInt(localStorage.getItem('timeout') || '10'))
const proxyType = ref(localStorage.getItem('proxyType') || 'none')
const proxyHost = ref(localStorage.getItem('proxyHost') || '')
const proxyPort = ref(parseInt(localStorage.getItem('proxyPort') || '1080'))
const dohUrl = ref(localStorage.getItem('dohUrl') || '')
const waifu2xModel = ref(localStorage.getItem('waifu2xModel') || 'cunet')
const outputFormat = ref(localStorage.getItem('outputFormat') || 'png')
const useGpu = ref(localStorage.getItem('useGpu') === 'true')
const downloadDir = ref(localStorage.getItem('downloadDir') || '/data/downloads')
const maxDownloads = ref(parseInt(localStorage.getItem('maxDownloads') || '3'))
const imageQuality = ref(localStorage.getItem('imageQuality') || 'original')
const cacheSize = ref(localStorage.getItem('cacheSize') || '512')
const autoCache = ref(localStorage.getItem('autoCache') !== 'false')

const save = () => {
  const items = { lang, theme, apiUrl, timeout, proxyType, proxyHost, proxyPort, dohUrl, waifu2xModel, outputFormat, useGpu, downloadDir, maxDownloads, imageQuality, cacheSize, autoCache }
  Object.entries(items).forEach(([k, v]) => localStorage.setItem(k, String(v.value)))
  ElMessage.success('设置已保存，部分设置需要刷新后生效')
}

const reset = () => {
  localStorage.clear()
  location.reload()
}
</script>

<style scoped>
.settings h2 { margin-bottom: 20px; }
.card { max-width: 700px; }
.actions { margin-top: 30px; display: flex; gap: 12px; }
</style>
