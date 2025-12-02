<template>
  <div class="host-page">
    <section class="hero">
      <div class="hero-left">
        <p class="eyebrow">HTML Hosting</p>
        <h1>托管一份 HTML，生成可访问链接</h1>
        <p class="sub">
          上传一份 HTML 文件，后端会帮你保存在静态目录下并生成链接，对方直接访问即可看到页面。适合快速演示、小型静态页分享。
        </p>
        <div class="hero-tags">
          <span class="pill">大小上限 {{ maxSizeLabel }}</span>
          <span class="pill">生成 /api/hosted 链接</span>
          <span class="pill">仅允许 .html / .htm</span>
        </div>
      </div>
      <div class="hero-card">
        <div class="meta">
          <p class="meta-label">当前上限</p>
          <p class="meta-value">{{ maxSizeLabel }}</p>
        </div>
        <div class="meta">
          <p class="meta-label">已选择</p>
          <p class="meta-value">
            {{ selectedFile ? selectedFile.name : "暂未选择" }}
          </p>
        </div>
      </div>
    </section>

    <section class="panel">
      <header class="panel-head">
        <div>
          <p class="eyebrow">上传文件</p>
          <h2>选择 HTML 并托管</h2>
          <p class="muted">
            支持拖拽或点击选择，超过大小限制会被拒收。链接可直接分享给他人访问。
          </p>
        </div>
        <div class="limit-pill">最大 {{ maxSizeLabel }}</div>
      </header>

      <div
        class="upload-area"
        :class="{ active: dragging }"
        @click="triggerSelect"
        @dragover.prevent="dragging = true"
        @dragleave.prevent="dragging = false"
        @drop.prevent="handleDrop"
      >
        <input
          ref="fileInput"
          type="file"
          class="file-input"
          accept=".html,.htm,text/html"
          @change="handleFileChange"
        />

        <template v-if="selectedFile">
          <div class="file-name">{{ selectedFile.name }}</div>
          <p class="file-meta">
            {{ formatBytes(selectedFile.size) }} ·
            {{ selectedFile.type || "text/html" }}
          </p>
          <p class="file-tip">点击或拖拽可以重新选择文件</p>
        </template>

        <template v-else>
          <p class="placeholder-title">拖拽 HTML 到这里，或点击选择</p>
          <p class="placeholder-sub">
            仅支持 .html / .htm · 最大 {{ maxSizeLabel }} · 将被托管到
            /api/hosted/
          </p>
        </template>
      </div>

      <div class="actions">
        <button
          class="primary"
          :disabled="uploading || !selectedFile"
          @click="handleUpload"
        >
          {{ uploading ? "正在托管..." : "上传并生成链接" }}
        </button>
        <button
          class="ghost"
          type="button"
          :disabled="uploading"
          @click="resetSelection"
        >
          清除选择
        </button>
        <span v-if="errorMessage" class="error">{{ errorMessage }}</span>
      </div>

      <div v-if="shareLink" class="result">
        <div class="result-row">
          <span class="label">访问链接</span>
          <a :href="shareLink" target="_blank" rel="noreferrer">
            {{ shareLink }}
          </a>
        </div>
        <div class="result-actions">
          <button class="ghost" type="button" @click="copyShareLink">
            复制链接
          </button>
          <a class="ghost" :href="shareLink" target="_blank" rel="noreferrer">
            新标签打开
          </a>
          <span v-if="copyState" class="hint">{{ copyState }}</span>
        </div>
      </div>

      <div v-if="previewText" class="preview">
        <div class="preview-head">
          <span>文件预览（前 800 个字符）</span>
          <span class="muted">仅供确认内容，实际以托管文件为准</span>
        </div>
        <pre>{{ previewText }}</pre>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from "vue";
import { uploadHtmlFile, type UploadHtmlResponse } from "../api/htmlHost";

const DEFAULT_LIMIT = 1024 * 1024;

const selectedFile = ref<File | null>(null);
const previewText = ref("");
const uploading = ref(false);
const errorMessage = ref("");
const shareLink = ref("");
const sizeLimitFromServer = ref<number | null>(null);
const copyState = ref("");
const dragging = ref(false);
const fileInput = ref<HTMLInputElement | null>(null);

const maxSize = computed(() => sizeLimitFromServer.value ?? DEFAULT_LIMIT);
const maxSizeLabel = computed(() => formatBytes(maxSize.value));

function formatBytes(bytes: number) {
  if (bytes >= 1024 * 1024) return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
  if (bytes >= 1024) return `${Math.ceil(bytes / 1024)} KB`;
  return `${bytes} B`;
}

function isHtmlFile(name: string) {
  const lower = name.toLowerCase();
  return lower.endsWith(".html") || lower.endsWith(".htm");
}

function triggerSelect() {
  fileInput.value?.click();
}

function resetSelection() {
  selectedFile.value = null;
  previewText.value = "";
  shareLink.value = "";
  copyState.value = "";
  errorMessage.value = "";
  if (fileInput.value) {
    fileInput.value.value = "";
  }
}

function setFile(file: File | null) {
  if (!file) {
    resetSelection();
    return;
  }

  if (!isHtmlFile(file.name)) {
    selectedFile.value = null;
    previewText.value = "";
    shareLink.value = "";
    copyState.value = "";
    errorMessage.value = "仅支持 .html / .htm 文件";
    return;
  }

  selectedFile.value = file;
  previewText.value = "";
  shareLink.value = "";
  copyState.value = "";
  errorMessage.value = "";

  const reader = new FileReader();
  reader.onload = () => {
    previewText.value = String(reader.result || "").slice(0, 800);
  };
  reader.readAsText(file);
}

function handleFileChange(event: Event) {
  const target = event.target as HTMLInputElement;
  const file = target.files?.[0] || null;
  setFile(file);
}

function handleDrop(event: DragEvent) {
  dragging.value = false;
  const file = event.dataTransfer?.files?.[0] || null;
  setFile(file);
}

function buildAbsoluteUrl(path?: string) {
  if (!path) return "";
  if (path.startsWith("http://") || path.startsWith("https://")) {
    return path;
  }
  const normalized = path.startsWith("/") ? path : `/${path}`;
  return `${window.location.origin}${normalized}`;
}

async function handleUpload() {
  if (!selectedFile.value) {
    errorMessage.value = "请选择要托管的 HTML 文件";
    return;
  }

  const limit = maxSize.value;
  if (selectedFile.value.size > limit) {
    errorMessage.value = `文件超过限制（最大 ${maxSizeLabel.value}）`;
    return;
  }

  uploading.value = true;
  errorMessage.value = "";
  shareLink.value = "";
  copyState.value = "";

  try {
    const res = await uploadHtmlFile(selectedFile.value);
    const data = res.data as UploadHtmlResponse;

    if (typeof data.limit_bytes === "number") {
      sizeLimitFromServer.value = data.limit_bytes;
    }

    const rawPath = data.url || (data.filename ? `/api/hosted/${data.filename}` : "");
    shareLink.value = data.full_url || buildAbsoluteUrl(rawPath);
  } catch (err: any) {
    const respLimit = err?.response?.data?.limit_bytes;
    if (typeof respLimit === "number") {
      sizeLimitFromServer.value = respLimit;
    }
    errorMessage.value =
      err?.response?.data?.error ||
      err?.message ||
      "上传失败，请检查后端是否运行";
  } finally {
    uploading.value = false;
  }
}

async function copyShareLink() {
  if (!shareLink.value) return;
  try {
    await navigator.clipboard.writeText(shareLink.value);
    copyState.value = "已复制";
  } catch {
    copyState.value = "复制失败，请手动复制";
  } finally {
    setTimeout(() => {
      copyState.value = "";
    }, 1600);
  }
}
</script>

<style scoped>
.host-page {
  max-width: 1100px;
  margin: 0 auto;
  padding: 12px 10px 26px;
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.hero {
  display: grid;
  grid-template-columns: 1.4fr 1fr;
  gap: 14px;
  align-items: center;
  background: linear-gradient(125deg, #0ea5e9 0%, #14b8a6 40%, #0f172a 100%);
  color: #ecfeff;
  padding: 20px 18px;
  border-radius: 16px;
  box-shadow: 0 16px 48px rgba(15, 23, 42, 0.35);
}

.hero-left h1 {
  margin: 6px 0 6px;
  font-size: 26px;
}

.sub {
  margin: 0 0 10px;
  line-height: 1.6;
  max-width: 720px;
}

.hero-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.pill {
  padding: 6px 10px;
  border-radius: 999px;
  border: 1px solid rgba(255, 255, 255, 0.3);
  background: rgba(255, 255, 255, 0.14);
  color: #ecfeff;
  font-size: 12px;
}

.hero-card {
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.24);
  border-radius: 14px;
  padding: 14px;
  display: grid;
  grid-template-columns: 1fr;
  gap: 12px;
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.2);
}

.meta {
  background: rgba(255, 255, 255, 0.08);
  border-radius: 10px;
  padding: 12px;
  border: 1px solid rgba(255, 255, 255, 0.14);
}

.meta-label {
  margin: 0;
  font-size: 12px;
  opacity: 0.82;
}

.meta-value {
  margin: 4px 0 0;
  font-size: 16px;
  font-weight: 700;
  color: #f8fafc;
  word-break: break-all;
}

.panel {
  background: var(--card-bg, #ffffff);
  border: 1px solid #e5e7eb;
  border-radius: 16px;
  padding: 18px;
  box-shadow: 0 12px 36px rgba(15, 23, 42, 0.08);
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.panel-head {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 10px;
}

.panel-head h2 {
  margin: 4px 0 6px;
  font-size: 22px;
}

.muted {
  margin: 0;
  color: #6b7280;
  font-size: 13px;
}

.limit-pill {
  padding: 8px 12px;
  border-radius: 999px;
  background: #ecfeff;
  color: #0369a1;
  border: 1px solid #bae6fd;
  align-self: flex-start;
  font-size: 13px;
}

.upload-area {
  border: 1.5px dashed #22d3ee;
  border-radius: 14px;
  padding: 22px 16px;
  background: linear-gradient(
    135deg,
    rgba(14, 165, 233, 0.06),
    rgba(20, 184, 166, 0.04)
  );
  text-align: center;
  cursor: pointer;
  transition: 0.18s ease;
}

.upload-area.active {
  border-color: #0ea5e9;
  background: linear-gradient(
    135deg,
    rgba(14, 165, 233, 0.12),
    rgba(20, 184, 166, 0.1)
  );
  box-shadow: 0 10px 26px rgba(14, 165, 233, 0.18);
}

.file-input {
  display: none;
}

.file-name {
  font-weight: 700;
  font-size: 16px;
  margin-bottom: 4px;
}

.file-meta,
.file-tip,
.placeholder-sub {
  margin: 4px 0;
  color: #475569;
  font-size: 13px;
}

.placeholder-title {
  margin: 0 0 4px;
  font-size: 16px;
  font-weight: 600;
  color: #0f172a;
}

.actions {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.primary,
.ghost {
  border-radius: 10px;
  padding: 10px 16px;
  font-weight: 700;
  cursor: pointer;
  border: none;
  transition: 0.18s ease;
}

.primary {
  background: linear-gradient(120deg, #0ea5e9, #14b8a6);
  color: #fff;
  box-shadow: 0 12px 28px rgba(14, 165, 233, 0.3);
}

.primary:disabled {
  opacity: 0.65;
  cursor: not-allowed;
  box-shadow: none;
}

.ghost {
  background: transparent;
  border: 1px solid #cbd5e1;
  color: #0f172a;
}

.ghost:hover {
  border-color: #0ea5e9;
  color: #0ea5e9;
}

.error {
  color: #b91c1c;
  font-size: 13px;
}

.result {
  border: 1px solid #bfdbfe;
  background: #eff6ff;
  border-radius: 12px;
  padding: 12px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.result-row {
  display: flex;
  gap: 8px;
  align-items: center;
  flex-wrap: wrap;
}

.result-row a {
  color: #1d4ed8;
  font-weight: 600;
  word-break: break-all;
  text-decoration: none;
}

.label {
  font-size: 12px;
  color: #0f172a;
  padding: 4px 8px;
  background: #dbeafe;
  border-radius: 8px;
}

.result-actions {
  display: flex;
  gap: 10px;
  align-items: center;
  flex-wrap: wrap;
}

.hint {
  font-size: 12px;
  color: #0f172a;
}

.preview {
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  background: #f8fafc;
  padding: 12px;
}

.preview-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.preview pre {
  margin: 0;
  white-space: pre-wrap;
  word-break: break-word;
  background: #0f172a;
  color: #e0f2fe;
  padding: 12px;
  border-radius: 10px;
  max-height: 320px;
  overflow: auto;
  font-family: "JetBrains Mono", Consolas, monospace;
  font-size: 13px;
}

.eyebrow {
  margin: 0;
  text-transform: uppercase;
  letter-spacing: 1px;
  font-size: 12px;
  color: #0ea5e9;
}

@media (max-width: 960px) {
  .hero {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 720px) {
  .host-page {
    padding: 10px 6px 20px;
  }

  .panel-head {
    flex-direction: column;
    align-items: flex-start;
  }

  .actions,
  .result-row,
  .result-actions,
  .hero-tags {
    flex-direction: column;
    align-items: flex-start;
  }

  .primary,
  .ghost {
    width: 100%;
    text-align: center;
  }
}

:global([data-theme="dark"]) .panel {
  background: #0b1222;
  border-color: #1f2a44;
}

:global([data-theme="dark"]) .muted,
:global([data-theme="dark"]) .file-meta,
:global([data-theme="dark"]) .file-tip,
:global([data-theme="dark"]) .placeholder-sub {
  color: #cbd5f5;
}

:global([data-theme="dark"]) .placeholder-title {
  color: #e2e8f0;
}

:global([data-theme="dark"]) .ghost {
  color: #e2e8f0;
  border-color: #334155;
}

:global([data-theme="dark"]) .result {
  background: rgba(59, 130, 246, 0.08);
  border-color: #1d4ed8;
}

:global([data-theme="dark"]) .label {
  background: rgba(59, 130, 246, 0.18);
  color: #e2e8f0;
}
</style>
