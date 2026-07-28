import axios from 'axios'

const api = axios.create({ baseURL: '/api', timeout: 30000 })

// 浏览
export const search = (q, page = 1) => api.get('/search', { params: { q, page } }).then(r => r.data)
export const getIndex = (page = 0) => api.get('/index', { params: { page } }).then(r => r.data)
export const getLatest = (page = 0) => api.get('/latest', { params: { page } }).then(r => r.data)
export const getComic = (id) => api.get(`/comic/${id}`).then(r => r.data)
export const getChapters = (id) => api.get(`/comic/${id}/chapters`).then(r => r.data)
export const getImages = (chapterId) => api.get(`/chapter/${chapterId}`).then(r => r.data)
export const getCategories = () => api.get('/categories').then(r => r.data)
export const getRanking = (type = 'daily', page = 1) => api.get('/ranking', { params: { type, page } }).then(r => r.data)
export const getCategoryFilter = (category, sort = 'mr', page = 1) => api.get('/categories/filter', { params: { category, sort, page } }).then(r => r.data)

// 评论
export const getComments = (comicId, page = 1) => api.get(`/comic/${comicId}/comments`, { params: { page } }).then(r => r.data)
export const getSubComments = (commentId, page = 1) => api.get(`/comment/${commentId}/sub`, { params: { page } }).then(r => r.data)

// 收藏
export const getFavorites = () => api.get('/favorites').then(r => r.data)
export const addFavorite = (comic_id) => api.post('/favorites', { comic_id })
export const removeFavorite = (id) => api.delete(`/favorites/${id}`)

// 历史
export const getHistory = () => api.get('/history').then(r => r.data)
export const deleteHistory = (id) => api.delete(`/history/${id}`)
export const clearHistory = () => api.delete('/history')

// 下载
export const createDownload = (comic_id) => api.post('/download', { comic_id })
export const getDownloads = () => api.get('/downloads').then(r => r.data)
export const deleteDownload = (id) => api.delete(`/download/${id}`)
export const clearDownloads = () => api.delete('/downloads')

// 用户
export const login = (username, password) => api.post('/login', { username, password }).then(r => r.data)

// 帮助
export const getHelp = () => api.get('/help').then(r => r.data)
