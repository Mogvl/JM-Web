<template>
  <div class="auth-page">
    <div class="auth-card">
      <div class="auth-header">
        <h1 class="auth-title">JMComic</h1>
        <p class="auth-subtitle">漫画阅读平台</p>
      </div>

      <el-tabs v-model="activeTab" class="auth-tabs">
        <el-tab-pane label="登录" name="login">
          <el-form :model="loginForm" :rules="loginRules" ref="loginRef" label-width="0" size="large">
            <el-form-item prop="username">
              <el-input v-model="loginForm.username" placeholder="用户名" prefix-icon="User" />
            </el-form-item>
            <el-form-item prop="password">
              <el-input v-model="loginForm.password" type="password" placeholder="密码" prefix-icon="Lock" show-password @keyup.enter="handleLogin" />
            </el-form-item>
            <div class="auth-options">
              <el-checkbox v-model="autoLogin">自动登录</el-checkbox>
              <el-checkbox v-model="savePassword">记住密码</el-checkbox>
            </div>
            <el-form-item>
              <el-button type="primary" size="large" style="width:100%" @click="handleLogin" :loading="loading" round>登录</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <el-tab-pane label="注册" name="register">
          <el-form :model="registerForm" :rules="registerRules" ref="registerRef" label-width="0" size="large">
            <el-form-item prop="username"><el-input v-model="registerForm.username" placeholder="用户名" prefix-icon="User" /></el-form-item>
            <el-form-item prop="email"><el-input v-model="registerForm.email" placeholder="邮箱" prefix-icon="Message" /></el-form-item>
            <el-form-item prop="password"><el-input v-model="registerForm.password" type="password" placeholder="密码" prefix-icon="Lock" show-password /></el-form-item>
            <el-form-item prop="password_confirm"><el-input v-model="registerForm.password_confirm" type="password" placeholder="确认密码" prefix-icon="Lock" show-password @keyup.enter="handleRegister" /></el-form-item>
            <el-form-item prop="gender">
              <el-radio-group v-model="registerForm.gender">
                <el-radio label="Male">男</el-radio>
                <el-radio label="Female">女</el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" size="large" style="width:100%" @click="handleRegister" :loading="loading" round>注册</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>
      </el-tabs>

      <div class="guest-link">
        <el-button text @click="handleGuest">游客访问</el-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { login, register } from '../api'
import { ElMessage } from 'element-plus'

const router = useRouter()
const activeTab = ref('login')
const loading = ref(false)
const loginRef = ref(null)
const registerRef = ref(null)
const autoLogin = ref(false)
const savePassword = ref(true)
const loginForm = ref({ username: '', password: '' })
const registerForm = ref({ username: '', email: '', password: '', password_confirm: '', gender: 'Male' })

const loginRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}
const registerRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  email: [{ required: true, message: '请输入邮箱', trigger: 'blur' }, { type: 'email', message: '邮箱格式不正确', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
  password_confirm: [{ required: true, message: '请确认密码', trigger: 'blur' }, { validator: (r, v, cb) => { v !== registerForm.value.password ? cb(new Error('两次密码不一致')) : cb() }, trigger: 'blur' }]
}

onMounted(() => {
  const u = localStorage.getItem('saved_username')
  const p = localStorage.getItem('saved_password')
  if (u) loginForm.value.username = u
  if (p) loginForm.value.password = p
})

const handleLogin = async () => {
  await loginRef.value.validate()
  loading.value = true
  try {
    const data = await login(loginForm.value.username, loginForm.value.password)
    localStorage.setItem('token', data.token)
    localStorage.setItem('username', data.username || loginForm.value.username)
    localStorage.setItem('coins', String(data.coins || 0))
    localStorage.setItem('level', String(data.level || 0))
    localStorage.setItem('level_name', data.level_name || '')
    localStorage.setItem('avatar', data.avatar || '')
    localStorage.setItem('favorites', String(data.favorites || 0))
    if (savePassword.value) { localStorage.setItem('saved_username', loginForm.value.username); localStorage.setItem('saved_password', loginForm.value.password) }
    else { localStorage.removeItem('saved_username'); localStorage.removeItem('saved_password') }
    ElMessage.success('登录成功')
    window.location.href = '/'
  } catch (e) {
    const msg = e?.response?.data?.error || e?.message || ''
    if (msg) ElMessage.error(msg)
    else ElMessage.error('登录失败，请检查用户名和密码')
  }
  finally { loading.value = false }
}

const handleRegister = async () => {
  await registerRef.value.validate()
  loading.value = true
  try {
    await register(registerForm.value.username, registerForm.value.email, registerForm.value.password, registerForm.value.password_confirm, registerForm.value.gender)
    ElMessage.success('注册成功，请登录')
    activeTab.value = 'login'
    loginForm.value.username = registerForm.value.username
    loginForm.value.password = ''
  } catch (e) { ElMessage.error('注册失败') }
  finally { loading.value = false }
}

const handleGuest = () => {
  localStorage.setItem('token', 'guest')
  localStorage.setItem('username', '游客')
  router.push('/')
}
</script>

<style scoped>
.auth-page { min-height: 100vh; display: flex; align-items: center; justify-content: center; background: linear-gradient(135deg, #F5F6F8 0%, #EBEDF0 100%); position: relative; overflow: hidden; }
.auth-page::before { content: ''; position: absolute; top: -50%; left: -50%; width: 200%; height: 200%; background: radial-gradient(circle at 30% 40%, rgba(255,77,109,0.1) 0%, transparent 50%), radial-gradient(circle at 70% 60%, rgba(99,102,241,0.08) 0%, transparent 50%); }
.auth-card { position: relative; width: 400px; padding: 40px; background: var(--bg-surface); border: 1px solid var(--border); border-radius: var(--radius-lg); box-shadow: 0 12px 40px rgba(0,0,0,0.08); }
.auth-header { text-align: center; margin-bottom: 28px; }
.auth-title { font-size: 32px; font-weight: 700; color: var(--accent); letter-spacing: -0.5px; }
.auth-subtitle { margin-top: 6px; color: var(--text-muted); font-size: 14px; }
.auth-options { display: flex; gap: 20px; margin-bottom: 20px; }
.guest-link { text-align: center; margin-top: 16px; }
</style>
