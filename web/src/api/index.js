import axios from 'axios'

const api = axios.create({ baseURL: '/api', timeout: 30000 })

// 自动添加 token
api.interceptors.request.use(config => {
  const token = localStorage.getItem('token')
  if (token && token !== 'guest') {
    config.headers.Authorization = 'Bearer ' + token
  }
  return config
})

// 浏览
export const search = (q, page = 1, sort = 'mr') => api.get('/search', { params: { q, page, sort } }).then(r => r.data)
export const getIndex = (page = 0) => api.get('/index', { params: { page } }).then(r => r.data)
export const getLatest = (page = 0) => api.get('/latest', { params: { page } }).then(r => r.data)
export const getComic = (id) => api.get(`/comic/${id}`).then(r => r.data)
export const getChapters = (id) => api.get(`/comic/${id}/chapters`).then(r => r.data)
export const getImages = (chapterId) => api.get(`/chapter/${chapterId}`).then(r => r.data)
export const getCategories = () => api.get('/categories').then(r => r.data)
export const getRanking = (type = 'daily', page = 1) => api.get('/ranking', { params: { type, page } }).then(r => r.data)
export const getCategoryFilter = (category, sort = 'mr', page = 1) => api.get('/categories/filter', { params: { category, sort, page } }).then(r => r.data)
export const getWeekCategories = () => api.get('/week').then(r => r.data)
export const getWeekFilter = (id, type, page = 0) => api.get('/week/filter', { params: { id, type, page } }).then(r => r.data)
// 阅读进度
export const getProgress = (comicId) => api.get(`/comic/${comicId}/progress`).then(r => r.data)
export const saveProgress = (comicId, chapterId, page) => api.post(`/comic/${comicId}/progress`, { chapter_id: chapterId, page })

// 评论
export const getComments = (comicId, page = 1) => api.get(`/comic/${comicId}/comments`, { params: { page } }).then(r => r.data)
export const getSubComments = (commentId, page = 1) => api.get(`/comment/${commentId}/sub`, { params: { page } }).then(r => r.data)
export const getMyComments = (page = 1) => api.get('/comments/mine', { params: { page } }).then(r => r.data)
export const getAllComments = (page = 1) => api.get('/comments/all', { params: { page } }).then(r => r.data)

// 收藏
export const getFavorites = () => api.get('/favorites').then(r => r.data)
export const addFavorite = (comic_id) => api.post('/favorites', { comic_id })
export const removeFavorite = (id) => api.delete(`/favorites/${id}`)
export const addFavoriteFolder = (name) => api.post('/favorite/folder', { name })
export const deleteFavoriteFolder = (fid) => api.delete(`/favorite/folder/${fid}`)
export const moveFavorite = (comic_id, folder_id) => api.post('/favorite/move', { comic_id, folder_id })

// 历史
export const getHistory = () => api.get('/history').then(r => r.data)
export const deleteHistory = (id) => api.delete(`/history/${id}`)
export const clearHistory = () => api.delete('/history')

// 下载
export const createDownload = (comic_id, info = {}) => api.post('/download', { comic_id, ...info })
export const getDownloads = () => api.get('/downloads').then(r => r.data)
export const deleteDownload = (id) => api.delete(`/download/${id}`)
export const clearDownloads = () => api.delete('/downloads')

// 用户
export const login = (username, password) => api.post('/login', { username, password }).then(r => r.data)
export const register = (username, email, password, password_confirm, gender) => api.post('/register', { username, email, password, password_confirm, gender }).then(r => r.data)
export const sign = () => api.post('/user/sign')
export const getUserInfo = () => api.get('/user/info').then(r => r.data)

// 帮助
export const getHelp = () => api.get('/help').then(r => r.data)
