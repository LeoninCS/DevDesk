<!-- src/views/WorkPlan.vue -->
<template>
  <div class="workplan-page">
    <div class="wp-card">
      <h1 class="wp-title">工作计划 / TODO 看板</h1>
      <p class="wp-subtitle">用一个链接，管理你一天的工作安排</p>

      <div class="wp-steps">
        <h2>如何使用？</h2>
        <ol>
          <li>点击下面的「开始使用」按钮，系统会为你生成一个专属计划链接（hash）。</li>
          <li>浏览器会自动跳转到你的计划页面，并根据 hash 载入你的 TODO 列表。</li>
          <li>把这个链接保存到收藏夹 / 记事本，下次直接打开即可继续使用同一个计划。</li>
          <li>你可以分享这个链接给别人，一起看同一个 TODO 看板（注意：拿到链接的人都可以修改）。</li>
        </ol>
      </div>

      <div class="wp-actions">
        <button
          class="wp-start-btn"
          :disabled="loading"
          @click="handleStart"
        >
          <span v-if="!loading">🚀 开始使用</span>
          <span v-else>生成中...</span>
        </button>

        <p v-if="error" class="wp-error">
          {{ error }}
        </p>

        <p v-if="hash" class="wp-hash-tip">
          已生成计划：<code>{{ hash }}</code><br />
          系统已为你跳转到该计划页面，可将当前链接保存以便下次使用。
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { createWorkPlan } from '../api/workplan.ts'

const router = useRouter()

const loading = ref(false)
const error = ref('')
const hash = ref('')

const handleStart = async () => {
  if (loading.value) return
  loading.value = true
  error.value = ''
  try {
    const res = await createWorkPlan()
    hash.value = res.data.hash

    // 跳转到 TODO 视图页面
    router.push({
      name: 'WorkPlanView',
      params: { hash: hash.value },
    })
  } catch (e: any) {
    error.value = e?.response?.data?.error || e?.message || '生成失败，请稍后重试'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.workplan-page {
  max-width: 780px;
  margin: 0 auto;
  padding: 32px 16px;
}

.wp-card {
  background: var(--moment-card-bg, #ffffff);
  border-radius: 16px;
  box-shadow: 0 8px 30px rgba(15, 23, 42, 0.08);
  padding: 28px 24px 32px;
}

.wp-title {
  margin: 0;
  font-size: 26px;
  font-weight: 700;
}

.wp-subtitle {
  margin: 8px 0 20px;
  color: #6b7280;
  font-size: 14px;
}

.wp-steps h2 {
  margin: 0 0 8px;
  font-size: 18px;
}

.wp-steps ol {
  margin: 0 0 16px 18px;
  padding: 0;
  color: #4b5563;
  font-size: 14px;
}

.wp-steps li + li {
  margin-top: 4px;
}

.wp-actions {
  margin-top: 24px;
  text-align: center;
}

.wp-start-btn {
  min-width: 180px;
  padding: 10px 20px;
  border-radius: 999px;
  border: none;
  outline: none;
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
  background: #2563eb;
  color: #ffffff;
  transition: transform 0.05s ease, box-shadow 0.05s ease, opacity 0.2s;
  box-shadow: 0 10px 25px rgba(37, 99, 235, 0.35);
}

.wp-start-btn:disabled {
  opacity: 0.7;
  cursor: default;
  box-shadow: none;
}

.wp-start-btn:not(:disabled):hover {
  transform: translateY(-1px);
}

.wp-start-btn:not(:disabled):active {
  transform: translateY(0);
  box-shadow: 0 6px 16px rgba(37, 99, 235, 0.3);
}

.wp-error {
  margin-top: 12px;
  color: #dc2626;
  font-size: 13px;
}

.wp-hash-tip {
  margin-top: 12px;
  font-size: 13px;
  color: #4b5563;
}

.wp-hash-tip code {
  background: #f3f4f6;
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 12px;
}
</style>
