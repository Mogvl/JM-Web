import axios from 'axios'

const api = axios.create({ baseURL: '/api', timeout: 30000 })

export const search = (q, page = 1) => api.get('/search', { params: { q, page } }).then(r => r.data)
export const getComic = (id) => api.get(`/comic/${id}`).then(r => r.data)
export const getChapters = (id) => api.get(`/comic/${id}/chapters`).then(r => r.data)
export const getImages = (chapterId) => api.get(`/chapter/${chapterId}`).then(r => r.data)
export const getCategories = () => api.get('/categories').then(r => r.data)

export const getFavorites = () => api.get('/favorites').then(r => r.data)
export const addFavorite = (comic_id) => api.post('/favorites', { comic_id })
export const removeFavorite = (id) => api.delete(`/favorites/${id}`)

export const getHistory = () => api.get('/history').then(r => r.data)
export const deleteHistory = (id) => api.delete(`/history/${id}`)

export const createDownload = (comic_id) => api.post('/download', { comic_id })
export const getDownloads = () => api.get('/downloads').then(r => r.data)
