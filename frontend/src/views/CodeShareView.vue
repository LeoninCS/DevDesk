<template>
  <div class="codeshare-view">
    <h1>代码片段</h1>

    <div v-if="loading" class="status">加载中...</div>
    <div v-else-if="error" class="status error">{{ error }}</div>

    <div
      v-else
      class="code-card"
      :class="`code-card--${codeTheme}`"
    >
      <div class="top-row">
        <div class="meta">
          <span v-if="author">发布者：{{ author }}</span>
          <span v-if="language" class="lang-tag">语法：{{ language }}</span>
        </div>

        <div class="actions">
          <button class="btn" @click="toggleCodeTheme">
            {{ codeTheme === 'dark' ? '切换为浅色' : '切换为深色' }}
          </button>
          <button class="btn" @click="copyCode">
            {{ copied ? '已复制' : '复制代码' }}
          </button>
        </div>
      </div>

      <pre class="pre-wrap">
        <code
          ref="codeEl"
          :class="['code-block', `language-${languageClass}`, `code-block--${codeTheme}`]"
        >{{ content }}</code>
      </pre>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, nextTick } from 'vue'
import { useRoute } from 'vue-router'
import hljs from 'highlight.js'
import 'highlight.js/styles/github-dark.css'

// ⭐ 用封装好的 API，而不是直接 axios
import { getCodeByHash } from '../api/codeshare.ts'

const route = useRoute()
const hash = route.params.hash as string

const loading = ref(true)
const error = ref('')
const author = ref('')
const language = ref('')
const content = ref('')

const codeEl = ref<HTMLElement | null>(null)

// 代码背景主题：dark / light
const codeTheme = ref<'dark' | 'light'>('dark')
// 复制状态
const copied = ref(false)

const languageClass = computed(() => {
  if (!language.value) return 'plaintext'
  if (language.value === 'c_cpp') return 'cpp'
  return language.value
})

onMounted(async () => {
  try {
    const res = await getCodeByHash(hash)

    author.value = res.data.author || ''
    language.value = res.data.language || ''
    content.value = res.data.content || ''
  } catch (e) {
    console.error(e)
    error.value = '该代码不存在或已过期。'
  } finally {
    loading.value = false
  }

  await nextTick()
  if (codeEl.value) {
    hljs.highlightElement(codeEl.value)
  }
})

function toggleCodeTheme() {
  codeTheme.value = codeTheme.value === 'dark' ? 'light' : 'dark'
}

async function copyCode() {
  const text = content.value || '';

  if (!text) return;

  // 兜底函数：用 textarea + execCommand 复制
  const fallbackCopy = () => {
    try {
      const textarea = document.createElement('textarea');
      textarea.value = text;
      textarea.style.position = 'fixed';
      textarea.style.opacity = '0';
      textarea.style.left = '-9999px';
      document.body.appendChild(textarea);
      textarea.focus();
      textarea.select();
      const ok = document.execCommand('copy');
      document.body.removeChild(textarea);
      return ok;
    } catch (err) {
      console.error('fallback copy error', err);
      return false;
    }
  };

  try {
    // SSR/非浏览器环境保护
    if (typeof navigator !== 'undefined' &&
        navigator.clipboard &&
        typeof navigator.clipboard.writeText === 'function') {
      await navigator.clipboard.writeText(text);
    } else {
      const ok = fallbackCopy();
      if (!ok) throw new Error('fallback copy failed');
    }

    // 复制成功
    copied.value = true;
    setTimeout(() => {
      copied.value = false;
    }, 1500);
  } catch (e) {
    console.error(e);
    alert('复制失败，请手动选择代码复制');
  }
}

</script>

<style scoped>
/* 你的原样式原封不动 */
.codeshare-view {
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px 4vw;
  box-sizing: border-box;
}

h1 {
  margin: 0 0 16px;
  font-size: 22px;
}

.status {
  padding: 12px 0;
  font-size: 14px;
  color: #6b7280;
}

.status.error {
  color: #ef4444;
}

.code-card {
  border-radius: 8px;
  padding: 16px 18px;
  border: 1px solid #1f2937;
  overflow: hidden;
}

/* 深色 / 浅色卡片背景 */
.code-card--dark {
  background: #21222c;
  color: #e5e7eb;
}

.code-card--light {
  background: #f9fafb;
  color: #111827;
  border-color: #d1d5db;
}

:global([data-theme='light']) .codeshare-view .code-card--dark {
  background: #111827;
  border-color: #1f2937;
}

/* 顶部：左信息 + 右侧按钮 */
.top-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  margin-bottom: 10px;
}

.meta {
  font-size: 13px;
  color: #9ca3af;
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  align-items: center;
}

.lang-tag {
  padding: 2px 8px;
  border-radius: 999px;
  border: 1px solid #4b5563;
  font-size: 12px;
}

/* 按钮区 */
.actions {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.btn {
  padding: 4px 10px;
  font-size: 12px;
  border-radius: 999px;
  border: 1px solid #4b5563;
  background: transparent;
  color: inherit;
  cursor: pointer;
  transition: background 0.2s, transform 0.05s;
}

.btn:hover {
  background: rgba(255, 255, 255, 0.08);
}

.btn:active {
  transform: scale(0.97);
}

/* 代码区域 */
.pre-wrap {
  margin: 0;
  font-size: 13px;
  line-height: 1.5;
  max-height: 70vh;
  overflow: auto;
  background: transparent;
}

/* code 本身根据主题切背景，!important 覆盖 highlight.js */
.code-block {
  display: block;
  white-space: pre;
}

.code-block--dark {
  background: #292a36 !important;
  color: #f5f6f8 !important;
}

.code-block--light {
  background: #eaeaeb !important;
  color: #111827 !important;
}

/* 小屏适配 */
@media (max-width: 768px) {
  .codeshare-view {
    padding: 16px 3vw;
  }

  .code-card {
    padding: 12px 12px;
  }

  h1 {
    font-size: 20px;
  }

  .top-row {
    flex-direction: column;
    align-items: flex-start;
  }

  .actions {
    align-self: stretch;
  }
}
</style>
