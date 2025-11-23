<!-- src/views/WorkPlanView.vue -->
<template>
  <div class="wpv-page">
    <div class="wpv-card">
      <header class="wpv-header">
        <div>
          <h1>我的工作计划</h1>
          <p class="wpv-tip">
            当前计划哈希：<code>{{ hash }}</code>
          </p>
        </div>

        <button class="wpv-back" @click="goIntro">
          ← 返回说明
        </button>
      </header>

      <!-- 分享提示：复制当前链接，让别人也能访问这个计划 -->
      <section v-if="!error" class="wpv-share">
        <span class="wpv-share-text">
          分享提示：把当前页面链接保存或发给同事，就可以访问同一份工作计划。
        </span>
        <button class="wpv-share-btn" @click="copyLink">
          {{ linkCopied ? '已复制链接' : '复制当前链接' }}
        </button>
      </section>

      <section class="wpv-info">
        <p>
          提示：这是一个基于 <strong>hash</strong> 的简单 TODO 看板。
          只要保留当前页面链接，未来都可以继续访问和修改这份计划。
        </p>
      </section>

      <section class="wpv-add">
        <input
          v-model.trim="newContent"
          class="wpv-input"
          type="text"
          placeholder="输入要做的事情，然后按回车或点击添加"
          @keyup.enter="handleAdd"
        />
        <button
          class="wpv-add-btn"
          :disabled="adding || !newContent"
          @click="handleAdd"
        >
          <span v-if="!adding">添加</span>
          <span v-else>添加中...</span>
        </button>
      </section>

      <p v-if="error" class="wpv-error">{{ error }}</p>

      <section class="wpv-list-section">
        <div v-if="loading" class="wpv-empty">
          正在载入 TODO 列表...
        </div>

        <div v-else-if="todos.length === 0" class="wpv-empty">
          还没有任何 TODO，先添加一条试试吧～
        </div>

        <ul v-else class="wpv-list">
          <li
            v-for="item in todos"
            :key="item.id"
            class="wpv-item"
            :class="{ done: item.done }"
          >
            <button
              class="wpv-check"
              :disabled="updatingId === item.id"
              @click="toggleDone(item)"
            >
              <span v-if="item.done">✓</span>
              <span v-else></span>
            </button>

            <div class="wpv-content" @dblclick="editItem(item)">
              {{ item.content }}
            </div>

            <div class="wpv-actions">
              <button
                class="wpv-action-btn"
                :disabled="updatingId === item.id"
                @click="editItem(item)"
              >
                编辑
              </button>
              <button
                class="wpv-action-btn danger"
                :disabled="updatingId === item.id"
                @click="deleteItem(item)"
              >
                删除
              </button>
            </div>
          </li>
        </ul>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import {
  fetchWorkPlanTodos,
  addWorkPlanTodo,
  toggleWorkPlanTodo,
  editWorkPlanTodo,
  deleteWorkPlanTodo,
  type TodoItem,
} from '../api/workplan.ts'

const route = useRoute()
const router = useRouter()

const hash = computed(() => route.params.hash as string)

const todos = ref<TodoItem[]>([])
const loading = ref(false)
const adding = ref(false)
const updatingId = ref<number | null>(null)
const newContent = ref('')
const error = ref('')

// 复制链接状态
const linkCopied = ref(false)

const fetchTodos = async () => {
  loading.value = true
  error.value = ''
  try {
    const res = await fetchWorkPlanTodos(hash.value)
    todos.value = res.data.todos || []
  } catch (e: any) {
    error.value =
      e?.response?.data?.error || e?.message || '加载 TODO 列表失败，请稍后重试'
  } finally {
    loading.value = false
  }
}

const handleAdd = async () => {
  if (!newContent.value || adding.value) return
  adding.value = true
  error.value = ''
  try {
    await addWorkPlanTodo(hash.value, newContent.value)
    newContent.value = ''
    await fetchTodos()
  } catch (e: any) {
    error.value =
      e?.response?.data?.error || e?.message || '添加失败，请稍后重试'
  } finally {
    adding.value = false
  }
}

const toggleDone = async (item: TodoItem) => {
  if (updatingId.value !== null) return
  updatingId.value = item.id
  error.value = ''
  try {
    await toggleWorkPlanTodo(hash.value, item.id)
    await fetchTodos()
  } catch (e: any) {
    error.value =
      e?.response?.data?.error || e?.message || '更新状态失败，请稍后重试'
  } finally {
    updatingId.value = null
  }
}

const editItem = async (item: TodoItem) => {
  const value = window.prompt('修改 TODO 内容：', item.content)
  if (value === null) return
  const trimmed = value.trim()
  if (!trimmed || trimmed === item.content) return

  if (updatingId.value !== null) return
  updatingId.value = item.id
  error.value = ''
  try {
    await editWorkPlanTodo(hash.value, item.id, trimmed)
    await fetchTodos()
  } catch (e: any) {
    error.value =
      e?.response?.data?.error || e?.message || '编辑失败，请稍后重试'
  } finally {
    updatingId.value = null
  }
}

const deleteItem = async (item: TodoItem) => {
  const ok = window.confirm(`确定要删除这条 TODO 吗？\n\n${item.content}`)
  if (!ok) return

  if (updatingId.value !== null) return
  updatingId.value = item.id
  error.value = ''
  try {
    await deleteWorkPlanTodo(hash.value, item.id)
    await fetchTodos()
  } catch (e: any) {
    error.value =
      e?.response?.data?.error || e?.message || '删除失败，请稍后重试'
  } finally {
    updatingId.value = null
  }
}

const goIntro = () => {
  router.push({ name: 'WorkPlan' })
}

// 复制当前页面链接
const copyLink = async () => {
  const url = typeof window !== 'undefined' ? window.location.href : ''
  if (!url) return

  const fallbackCopy = () => {
    try {
      const textarea = document.createElement('textarea')
      textarea.value = url
      textarea.style.position = 'fixed'
      textarea.style.opacity = '0'
      textarea.style.left = '-9999px'
      document.body.appendChild(textarea)
      textarea.focus()
      textarea.select()
      const ok = document.execCommand('copy')
      document.body.removeChild(textarea)
      return ok
    } catch (err) {
      console.error('fallback copy link error', err)
      return false
    }
  }

  try {
    if (
      typeof navigator !== 'undefined' &&
      navigator.clipboard &&
      typeof navigator.clipboard.writeText === 'function'
    ) {
      await navigator.clipboard.writeText(url)
    } else {
      const ok = fallbackCopy()
      if (!ok) throw new Error('fallback copy link failed')
    }

    linkCopied.value = true
    setTimeout(() => {
      linkCopied.value = false
    }, 1500)
  } catch (e) {
    console.error(e)
    alert('复制链接失败，请手动复制地址栏链接')
  }
}

onMounted(() => {
  if (!hash.value) {
    error.value = '缺少计划 hash'
    return
  }
  fetchTodos()
})
</script>

<style scoped>
.wpv-page {
  max-width: 900px;
  margin: 0 auto;
  padding: 24px 16px 40px;
}

.wpv-card {
  background: var(--moment-card-bg, #ffffff);
  border-radius: 16px;
  box-shadow: 0 8px 30px rgba(15, 23, 42, 0.08);
  padding: 22px 20px 26px;
}

.wpv-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 8px;
}

.wpv-header h1 {
  margin: 0;
  font-size: 22px;
}

.wpv-tip {
  margin: 4px 0 0;
  font-size: 12px;
  color: #6b7280;
}

.wpv-tip code {
  background: #f3f4f6;
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 11px;
}

.wpv-back {
  border: none;
  border-radius: 999px;
  padding: 6px 14px;
  font-size: 12px;
  cursor: pointer;
  background: #f3f4f6;
  color: #374151;
}

/* 分享提示条 */
.wpv-share {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  margin-bottom: 10px;
  padding: 8px 10px;
  border-radius: 8px;
  background: rgba(37, 99, 235, 0.06);
  border: 1px dashed rgba(37, 99, 235, 0.3);
  font-size: 13px;
  color: #374151;
}

.wpv-share-text {
  flex: 1;
}

.wpv-share-btn {
  flex-shrink: 0;
  padding: 4px 10px;
  font-size: 12px;
  border-radius: 999px;
  border: 1px solid #2563eb;
  background: #2563eb;
  color: #ffffff;
  cursor: pointer;
  transition: transform 0.05s ease, box-shadow 0.1s ease;
}

.wpv-share-btn:active {
  transform: scale(0.97);
}

.wpv-info {
  margin: 4px 0 18px;
  font-size: 13px;
  color: #4b5563;
}

.wpv-add {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 10px;
}

.wpv-input {
  flex: 1;
  padding: 8px 10px;
  border-radius: 10px;
  border: 1px solid #d1d5db;
  font-size: 14px;
  outline: none;
}

.wpv-input:focus {
  border-color: #2563eb;
}

.wpv-add-btn {
  padding: 7px 16px;
  border-radius: 999px;
  border: none;
  font-size: 13px;
  font-weight: 500;
  background: #2563eb;
  color: #ffffff;
  cursor: pointer;
}

.wpv-add-btn:disabled {
  opacity: 0.6;
  cursor: default;
}

.wpv-error {
  margin: 4px 0 10px;
  color: #dc2626;
  font-size: 12px;
}

.wpv-list-section {
  margin-top: 8px;
}

.wpv-empty {
  padding: 18px 10px;
  font-size: 13px;
  color: #6b7280;
  text-align: center;
}

.wpv-list {
  list-style: none;
  margin: 0;
  padding: 0;
}

.wpv-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 6px;
  border-bottom: 1px solid #e5e7eb;
}

.wpv-item:last-child {
  border-bottom: none;
}

.wpv-item.done .wpv-content {
  text-decoration: line-through;
  color: #9ca3af;
}

.wpv-check {
  width: 20px;
  height: 20px;
  border-radius: 999px;
  border: 1px solid #d1d5db;
  background: #ffffff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 13px;
  cursor: pointer;
}

.wpv-check span {
  line-height: 1;
}

.wpv-content {
  flex: 1;
  font-size: 14px;
  cursor: default;
  word-break: break-word;
}

.wpv-actions {
  display: flex;
  gap: 6px;
}

.wpv-action-btn {
  border: none;
  border-radius: 999px;
  padding: 4px 10px;
  font-size: 12px;
  cursor: pointer;
  background: #f3f4f6;
  color: #4b5563;
}

.wpv-action-btn.danger {
  background: #fee2e2;
  color: #b91c1c;
}

.wpv-action-btn:disabled {
  opacity: 0.6;
  cursor: default;
}

@media (max-width: 768px) {
  .wpv-share {
    flex-direction: column;
    align-items: flex-start;
  }

  .wpv-share-btn {
    align-self: flex-end;
  }
}
</style>
