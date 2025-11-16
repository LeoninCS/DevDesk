// src/utils/api.js
import axios from 'axios'

// 环境变量（如果没配置就默认 '/api'）
const baseURL = import.meta.env.VITE_API_BASE || '/'

export const api = axios.create({
  baseURL,
  headers: {
    'Content-Type': 'application/json'
  },
  timeout: 15000
})


// =====================================
// CodeShare API
// =====================================

// 上传代码 -> POST /api/codeshare/upload
export async function uploadCode({ author, language, content, ttl }) {
  const { data } = await api.post('/codeshare/upload', {
    author,
    language,
    content,
    ttl
  })
  return data
}

// 获取代码 -> GET /api/codeshare/code/:hash
export async function getCode(hash) {
  const { data } = await api.get(`/codeshare/code/${encodeURIComponent(hash)}`)
  return data
}

// =====================================
// TextCrypto API
// =====================================

// 获取随机密钥 -> GET /api/textcrypto/random-screat-key
export async function getRandomSecretKey() {
  const { data } = await api.get('/textcrypto/random-screat-key')

  if (typeof data === 'string') return data.trim()
  return (data.screat_key || data.key || '').trim()
}

// 加密 -> POST /api/textcrypto/encrypt-text
export async function encryptText({ text, screat_key }) {
  const { data } = await api.post('/textcrypto/encrypt-text', {
    text,
    screat_key
  })
  return data
}

// 解密 -> POST /api/textcrypto/decrypt-text
export async function decryptText({ encrypted_text, screat_key }) {
  const { data } = await api.post('/textcrypto/decrypt-text', {
    encrypted_text,
    screat_key
  })
  return data
}
