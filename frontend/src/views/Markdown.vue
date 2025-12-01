<!-- src/views/Markdown.vue -->
<template>
  <div class="md-root">
    <section class="md-hero">
      <div class="md-text">
        <p class="eyebrow">Markdown Lab</p>
        <h1 class="md-title">在线 Markdown 协同编辑</h1>
        <p class="md-subtitle">
          生成一个带哈希的文档链接，随时进入编辑。支持协作、实时保存，以及导出 HTML / PDF / MD。
        </p>

        <ul class="md-bullets">
          <li>🧑‍💻 左侧编辑，右侧实时预览</li>
          <li>🤝 分享链接即可多人协同</li>
          <li>📄 一键导出为 HTML / PDF / Markdown</li>
        </ul>

        <div class="md-actions">
          <button
            class="md-start-btn"
            :disabled="loading"
            @click="handleStart"
          >
            <span v-if="!loading">✏️ 创建新文档</span>
            <span v-else>生成中...</span>
          </button>

          <p v-if="error" class="md-error">
            {{ error }}
          </p>

          <p v-if="hash" class="md-hint">
            已生成文档：<code>{{ hash }}</code>，正在跳转到编辑页面。
          </p>
        </div>
      </div>

      <div class="md-hero-preview">
        <div class="md-preview-window">
          <div class="md-preview-header">
            <span></span><span></span><span></span>
          </div>
          <div class="md-preview-body">
            <div class="md-preview-left">
              <div class="md-line" v-for="i in 8" :key="i"></div>
            </div>
            <div class="md-preview-right">
              <div class="md-block-title"></div>
              <div class="md-block-line" v-for="i in 5" :key="i"></div>
            </div>
          </div>
        </div>
        <div class="md-preview-foot">
          <span class="foot-chip">实时保存</span>
          <span class="foot-chip">多人协作</span>
          <span class="foot-chip">导出 PDF/HTML/MD</span>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";
import { createMarkdownDoc } from "../api/markdown.ts";

const router = useRouter();

const loading = ref(false);
const error = ref("");
const hash = ref("");

const handleStart = async () => {
  if (loading.value) return;
  loading.value = true;
  error.value = "";

  try {
    const res = await createMarkdownDoc();
    hash.value = res.data.hash;

    router.push({
      name: "MarkdownEditorView",
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
.md-root {
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
  padding: 12px 8px 24px;
  box-sizing: border-box;
}

.md-hero {
  display: grid;
  grid-template-columns: 1.2fr 1fr;
  gap: 32px;
  align-items: center;
  background: linear-gradient(135deg, #0ea5e9 0%, #2563eb 45%, #0f172a 100%);
  border-radius: 18px;
  padding: 24px 24px 26px;
  color: #e5e7eb;
  box-shadow: 0 16px 50px rgba(15, 23, 42, 0.4);
}

.eyebrow {
  margin: 0;
  text-transform: uppercase;
  letter-spacing: 1px;
  font-size: 12px;
  opacity: 0.8;
}

.md-text {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.md-title {
  margin: 0;
  font-size: 30px;
  font-weight: 800;
}

.md-subtitle {
  margin: 0 0 8px;
  font-size: 15px;
  color: #cbd5e1;
  line-height: 1.6;
}

.md-bullets {
  margin: 0 0 12px;
  padding-left: 18px;
  font-size: 14px;
  color: #e2e8f0;
}

.md-bullets li + li {
  margin-top: 4px;
}

.md-actions {
  margin-top: 4px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.md-start-btn {
  min-width: 220px;
  padding: 11px 22px;
  border-radius: 999px;
  border: none;
  outline: none;
  font-size: 15px;
  font-weight: 700;
  cursor: pointer;
  background: #fbbf24;
  color: #0f172a;
  box-shadow: 0 14px 40px rgba(251, 191, 36, 0.35);
  transition: transform 0.08s ease, box-shadow 0.08s ease, opacity 0.2s;
}

.md-start-btn:disabled {
  opacity: 0.7;
  cursor: default;
  box-shadow: none;
}

.md-start-btn:not(:disabled):hover {
  transform: translateY(-1px);
}

.md-start-btn:not(:disabled):active {
  transform: translateY(1px);
  box-shadow: 0 8px 24px rgba(251, 191, 36, 0.45);
}

.md-error {
  margin: 0;
  color: #fecdd3;
  font-size: 13px;
}

.md-hint {
  margin: 0;
  font-size: 13px;
  color: #e2e8f0;
}

.md-hint code {
  background: rgba(255, 255, 255, 0.08);
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 12px;
  color: #fff;
}

.md-hero-preview {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
}

.md-preview-window {
  width: 100%;
  max-width: 440px;
  border-radius: 18px;
  background: #0f172a;
  box-shadow: 0 18px 40px rgba(15, 23, 42, 0.6);
  overflow: hidden;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.md-preview-header {
  height: 28px;
  padding: 8px 12px;
  display: flex;
  align-items: center;
  gap: 6px;
}

.md-preview-header span {
  width: 10px;
  height: 10px;
  border-radius: 999px;
  background: #475569;
}

.md-preview-body {
  display: flex;
  padding: 12px;
  gap: 10px;
  background: radial-gradient(circle at top left, #1f2937, #020617);
}

.md-preview-left {
  flex: 1;
  border-radius: 12px;
  background: rgba(15, 23, 42, 0.9);
  padding: 10px 8px;
}

.md-line {
  height: 6px;
  margin-bottom: 6px;
  border-radius: 999px;
  background: linear-gradient(90deg, #4b5563, #9ca3af);
  opacity: 0.85;
}

.md-preview-right {
  flex: 1;
  border-radius: 12px;
  background: #f9fafb;
  padding: 10px 10px 16px;
}

.md-block-title {
  width: 70%;
  height: 14px;
  border-radius: 999px;
  background: #e5e7eb;
  margin-bottom: 10px;
}

.md-block-line {
  width: 100%;
  height: 8px;
  border-radius: 999px;
  background: #e5e7eb;
  margin-bottom: 6px;
}

.md-preview-foot {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.foot-chip {
  padding: 6px 10px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.2);
  font-size: 12px;
}

@media (max-width: 900px) {
  .md-root {
    padding: 12px 4px 20px;
  }
  .md-hero {
    grid-template-columns: 1fr;
  }
  .md-hero-preview {
    order: -1;
    margin-bottom: 12px;
  }
}
</style>
