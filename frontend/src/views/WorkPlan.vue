<!-- src/views/WorkPlan.vue -->
<template>
  <div class="workplan-page">
    <section class="wp-hero">
      <div>
        <p class="eyebrow">TODO Board</p>
        <h1 class="wp-title">工作计划 / TODO 看板</h1>
        <p class="wp-subtitle">
          用一个哈希链接管理你的每日待办，与队友共享同一份计划，随时进入即可继续。
        </p>
        <ul class="wp-list">
          <li>⚡ 一键生成专属链接</li>
          <li>🧭 刷新后仍可继续使用同一看板</li>
          <li>🤝 分享链接即可多人协作</li>
        </ul>
        <div class="wp-actions">
          <button class="wp-start-btn" :disabled="loading" @click="handleStart">
            <span v-if="!loading">🚀 开始使用</span>
            <span v-else>生成中...</span>
          </button>
          <p v-if="error" class="wp-error">
            {{ error }}
          </p>
          <p v-if="hash" class="wp-hash-tip">
            已生成计划：<code>{{ hash }}</code><br />
            正在跳转到你的看板，可保存当前链接备用。
          </p>
        </div>
      </div>
      <div class="wp-card">
        <div class="wp-card-head">
          <span class="dot"></span>
          <span class="dot"></span>
          <span class="dot"></span>
        </div>
        <div class="wp-card-body">
          <div class="wp-card-col">
            <div class="wp-card-line" v-for="i in 5" :key="i"></div>
          </div>
          <div class="wp-card-col">
            <div class="wp-card-task">
              <span class="task-dot"></span>
              <div class="task-lines">
                <div class="task-line short"></div>
                <div class="task-line"></div>
              </div>
            </div>
            <div class="wp-card-task">
              <span class="task-dot done"></span>
              <div class="task-lines">
                <div class="task-line short"></div>
                <div class="task-line"></div>
              </div>
            </div>
            <div class="wp-card-task">
              <span class="task-dot"></span>
              <div class="task-lines">
                <div class="task-line short"></div>
                <div class="task-line"></div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";
import { createWorkPlan } from "../api/workplan.ts";

const router = useRouter();

const loading = ref(false);
const error = ref("");
const hash = ref("");

const handleStart = async () => {
  if (loading.value) return;
  loading.value = true;
  error.value = "";
  try {
    const res = await createWorkPlan();
    hash.value = res.data.hash;

    router.push({
      name: "WorkPlanView",
      params: { hash: hash.value },
    });
  } catch (e: any) {
    error.value =
      e?.response?.data?.error || e?.message || "生成失败，请稍后重试";
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
.workplan-page {
  max-width: 1100px;
  margin: 0 auto;
  padding: 24px 12px 32px;
}

.wp-hero {
  display: grid;
  grid-template-columns: 1.2fr 1fr;
  gap: 26px;
  align-items: center;
  background: linear-gradient(135deg, #6366f1 0%, #2563eb 45%, #0f172a 100%);
  border-radius: 18px;
  padding: 22px 22px 26px;
  color: #f8fafc;
  box-shadow: 0 16px 48px rgba(37, 99, 235, 0.38);
}

.eyebrow {
  margin: 0;
  font-size: 12px;
  text-transform: uppercase;
  letter-spacing: 1px;
  opacity: 0.8;
}

.wp-title {
  margin: 6px 0 6px;
  font-size: 28px;
  font-weight: 800;
}

.wp-subtitle {
  margin: 0 0 10px;
  font-size: 15px;
  color: #e0e7ff;
  line-height: 1.6;
}

.wp-list {
  margin: 0 0 14px;
  padding-left: 18px;
  color: #e2e8f0;
  font-size: 14px;
}

.wp-list li + li {
  margin-top: 4px;
}

.wp-actions {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.wp-start-btn {
  min-width: 180px;
  padding: 12px 20px;
  border-radius: 999px;
  border: none;
  outline: none;
  font-size: 15px;
  font-weight: 700;
  cursor: pointer;
  background: #fbbf24;
  color: #0f172a;
  transition: transform 0.05s ease, box-shadow 0.05s ease, opacity 0.2s;
  box-shadow: 0 12px 28px rgba(251, 191, 36, 0.4);
}

.wp-start-btn:disabled {
  opacity: 0.7;
  cursor: default;
  box-shadow: none;
}

.wp-error {
  margin: 0;
  color: #fecdd3;
  font-size: 13px;
}

.wp-hash-tip {
  margin: 0;
  font-size: 13px;
  color: #e2e8f0;
}

.wp-hash-tip code {
  background: rgba(255, 255, 255, 0.12);
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 12px;
  color: #fff;
}

.wp-card {
  border-radius: 16px;
  background: #0f172a;
  border: 1px solid rgba(255, 255, 255, 0.12);
  padding: 12px;
  box-shadow: 0 12px 30px rgba(0, 0, 0, 0.28);
}

.wp-card-head {
  height: 24px;
  display: flex;
  align-items: center;
  gap: 6px;
}

.dot {
  width: 10px;
  height: 10px;
  border-radius: 999px;
  background: #334155;
}

.wp-card-body {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
  margin-top: 8px;
}

.wp-card-col {
  background: rgba(255, 255, 255, 0.04);
  border-radius: 12px;
  padding: 10px;
}

.wp-card-line {
  height: 8px;
  border-radius: 999px;
  background: linear-gradient(90deg, #475569, #c7d2fe);
  opacity: 0.9;
  margin-bottom: 8px;
}

.wp-card-task {
  display: flex;
  gap: 8px;
  align-items: center;
  padding: 8px;
  border-radius: 10px;
  background: rgba(255, 255, 255, 0.04);
  margin-bottom: 8px;
}

.task-dot {
  width: 10px;
  height: 10px;
  border-radius: 999px;
  background: #22c55e;
}

.task-dot.done {
  background: #a855f7;
}

.task-lines {
  flex: 1;
}

.task-line {
  height: 6px;
  border-radius: 999px;
  background: #c7d2fe;
  margin-bottom: 4px;
}

.task-line.short {
  width: 70%;
}

@media (max-width: 900px) {
  .wp-hero {
    grid-template-columns: 1fr;
  }
}
</style>
