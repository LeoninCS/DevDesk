<!-- src/views/MarkdownEditorView.vue -->
<template>
  <div class="mde-page">
    <div class="mde-card">
      <!-- 顶部工具栏 -->
      <header class="mde-header">
        <div class="mde-header-left">
          <h1 class="mde-title">Markdown 编辑器</h1>
          <p class="mde-meta">
            文档哈希：<code>{{ hash }}</code>
          </p>
          <div class="mde-share-row">
            <span class="mde-meta-label">分享链接：</span>
            <input
              ref="shareInputRef"
              class="mde-share-input"
              :value="shareUrl"
              readonly
            />
            <button class="mde-copy-btn" @click="copyShareLink">复制</button>
          </div>
        </div>

        <div class="mde-header-right">
          <div class="mde-status-box">
            <span v-if="saving" class="mde-status saving">正在保存...</span>
            <span v-else-if="lastSavedAt" class="mde-status saved">
              已保存：{{ lastSavedAt }}
            </span>
            <span v-else class="mde-status idle">尚未保存</span>
          </div>
          <div class="mde-export-group">
            <button
              class="mde-btn"
              :disabled="exporting"
              @click="exportMarkdown"
            >
              导入Markdown
            </button>
            <button class="mde-btn" :disabled="exporting" @click="exportHtml">
              导出HTML
            </button>
            <button class="mde-btn" :disabled="exporting" @click="exportPdf">
              导出PDF
            </button>
          </div>
          <button class="mde-back" @click="goIntro">
            返回
          </button>
        </div>
      </header>

      <!-- 错误提示 -->
      <p v-if="error" class="mde-error">{{ error }}</p>

      <!-- 主编辑区：左右分栏，占满卡片剩余高度 -->
      <main class="mde-main">
        <section class="mde-pane">
          <div class="mde-pane-header">编辑区</div>
          <textarea
            v-model="content"
            class="mde-editor"
            placeholder="在这里输入 Markdown 内容..."
            @input="handleInput"
          ></textarea>
        </section>

        <section class="mde-pane">
          <div class="mde-pane-header">预览</div>
          <div ref="previewRef" class="mde-preview" v-html="previewHtml"></div>
        </section>
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { marked } from "marked";
import html2canvas from "html2canvas";
import { jsPDF } from "jspdf";

import {
  fetchMarkdownDoc,
  updateMarkdownDoc,
  type MarkdownDocResponse,
} from "../api/markdown.ts";

const route = useRoute();
const router = useRouter();
const apiBaseUrl =
  import.meta.env.VITE_API_BASE_URL ||
  (import.meta.env.DEV ? "http://localhost:8080/api" : "/api");

const hash = computed(() => route.params.hash as string);

const content = ref("");
const error = ref("");
const saving = ref(false);
const lastSavedAt = ref("");

const exporting = ref(false);

const previewRef = ref<HTMLElement | null>(null);
const shareInputRef = ref<HTMLInputElement | null>(null);

let sendTimer: number | null = null;
let es: EventSource | null = null;
let isRemoteUpdate = false;

// 预览 HTML
const previewHtml = computed(() => {
  return marked.parse(content.value || "");
});

// 分享链接：直接用当前地址
const shareUrl = computed(() => window.location.href);

const formatTime = (d: Date) => {
  const pad = (n: number) => (n < 10 ? "0" + n : "" + n);
  return (
    pad(d.getHours()) + ":" + pad(d.getMinutes()) + ":" + pad(d.getSeconds())
  );
};

const initDoc = async () => {
  if (!hash.value) {
    error.value = "缺少文档 hash";
    return;
  }
  error.value = "";

  try {
    const res = await fetchMarkdownDoc(hash.value);
    const data: MarkdownDocResponse = res.data;
    isRemoteUpdate = true;
    content.value = data.content || "";
    isRemoteUpdate = false;
  } catch (e: any) {
    error.value =
      e?.response?.data?.error || e?.message || "加载文档失败，请稍后重试";
  }
};

// 本地编辑 -> 防抖保存
const handleInput = () => {
  if (isRemoteUpdate) return;
  if (!hash.value) return;

  if (sendTimer) {
    window.clearTimeout(sendTimer);
  }

  sendTimer = window.setTimeout(async () => {
    saving.value = true;
    try {
      await updateMarkdownDoc(hash.value, content.value);
      lastSavedAt.value = formatTime(new Date());
    } catch (e: any) {
      error.value =
        e?.response?.data?.error || e?.message || "保存失败，请稍后重试";
    } finally {
      saving.value = false;
    }
  }, 400);
};

// 复制分享链接
const copyShareLink = async () => {
  try {
    if (navigator.clipboard) {
      await navigator.clipboard.writeText(shareUrl.value);
      alert("已复制链接，可以发给别人一起编辑");
    } else if (shareInputRef.value) {
      shareInputRef.value.select();
      document.execCommand("copy");
      alert("已复制链接，可以发给别人一起编辑");
    }
  } catch {
    alert("复制失败，请手动复制输入框内容");
  }
};

// 导出 .md
const exportMarkdown = () => {
  const blob = new Blob([content.value || ""], {
    type: "text/markdown;charset=utf-8",
  });
  const url = URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = `${hash.value || "document"}.md`;
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
  URL.revokeObjectURL(url);
};

// 导出 .html
const exportHtml = () => {
  const bodyContent = previewHtml.value;
  const fullHtml = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="UTF-8" />
<title>导出的 Markdown 文档</title>
<style>
body {
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue",
               Arial, "Noto Sans", "PingFang SC", "Microsoft YaHei", sans-serif;
  padding: 20px;
  line-height: 1.6;
}
h1 { font-size: 24px; margin: 16px 0 8px; }
h2 { font-size: 20px; margin: 14px 0 8px; }
h3 { font-size: 18px; margin: 12px 0 6px; }
p  { margin: 6px 0; }
code {
  background: #f2f2f2;
  padding: 2px 4px;
  border-radius: 3px;
  font-family: "JetBrains Mono", monospace;
  font-size: 13px;
}
pre code {
  display: block;
  padding: 10px;
  overflow-x: auto;
}
ul, ol { padding-left: 22px; margin: 6px 0; }
blockquote {
  border-left: 4px solid #ddd;
  margin: 6px 0;
  padding: 4px 10px;
  color: #666;
  background: #fafafa;
}
a { color: #2563eb; text-decoration: none; }
a:hover { text-decoration: underline; }
</style>
</head>
<body>
${bodyContent}
</body>
</html>`;

  const blob = new Blob([fullHtml], {
    type: "text/html;charset=utf-8",
  });
  const url = URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = `${hash.value || "document"}.html`;
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
  URL.revokeObjectURL(url);
};

// 导出 .pdf（基于预览区截图）
const exportPdf = async () => {
  if (!previewRef.value) {
    alert("预览区未就绪，稍后再试");
    return;
  }
  exporting.value = true;
  try {
    const canvas = await html2canvas(previewRef.value, { scale: 2 });
    const imgData = canvas.toDataURL("image/png");

    const pdf = new jsPDF("p", "pt", "a4");
    const pageWidth = pdf.internal.pageSize.getWidth();
    const pageHeight = pdf.internal.pageSize.getHeight();

    const imgWidth = pageWidth;
    const imgHeight = (canvas.height * imgWidth) / canvas.width;

    let heightLeft = imgHeight;
    let position = 0;

    pdf.addImage(imgData, "PNG", 0, position, imgWidth, imgHeight);
    heightLeft -= pageHeight;

    while (heightLeft > 0) {
      position = heightLeft - imgHeight;
      pdf.addPage();
      pdf.addImage(imgData, "PNG", 0, position, imgWidth, imgHeight);
      heightLeft -= pageHeight;
    }

    pdf.save(`${hash.value || "document"}.pdf`);
  } catch (e) {
    console.error(e);
    alert("导出 PDF 失败，请稍后重试");
  } finally {
    exporting.value = false;
  }
};

// SSE 协同
const initSSE = () => {
  if (!hash.value) return;
  es = new EventSource(`${apiBaseUrl}/markdown/stream/${hash.value}`);

  es.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data);
      const remoteContent = data.content || "";

      if (remoteContent === content.value) return;

      isRemoteUpdate = true;
      content.value = remoteContent;
      isRemoteUpdate = false;
    } catch (e) {
      console.error("parse sse data error", e);
    }
  };

  es.onerror = (err) => {
    console.error("SSE error", err);
    es?.close();
  };
};

const goIntro = () => {
  router.push({ name: "Markdown" });
};

onMounted(async () => {
  await initDoc();
  initSSE();
});

onBeforeUnmount(() => {
  if (es) {
    es.close();
    es = null;
  }
  if (sendTimer) {
    window.clearTimeout(sendTimer);
  }
});
</script>

<style scoped>
.mde-page {
  width: 100%;
  /* 原来有 max-width + margin: 0 auto，会被居中变窄，这里干掉 */
  margin: 0;
  /* 上下保留一点内边距，左右贴满 content */
  padding: 4px 0 16px;
}

.mde-card {
  background: var(--moment-card-bg, #ffffff);
  border-radius: 16px;
  box-shadow: 0 10px 30px rgba(15, 23, 42, 0.08);
  padding: 16px 18px 18px;
  display: flex;
  flex-direction: column;
  box-sizing: border-box;

  /* 关键：宽度铺满 content，高度尽量吃满视口 */
  width: 100%;
  height: calc(100vh - 120px); /* 这里 120px 预留给顶部导航+content padding，可以按需要微调 */
  max-height: none;
}

.mde-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
}

.mde-header-left {
  flex: 1;
  min-width: 0;
}

.mde-title {
  margin: 0;
  font-size: 20px;
  font-weight: 700;
}

.mde-meta {
  margin: 4px 0 0;
  font-size: 12px;
  color: #6b7280;
}

.mde-meta code {
  background: #e5e7eb;
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 11px;
}

.mde-share-row {
  margin-top: 6px;
  display: flex;
  align-items: center;
  gap: 6px;
  flex-wrap: wrap;
}

.mde-meta-label {
  font-size: 12px;
  color: #6b7280;
}

.mde-share-input {
  flex: 1;
  min-width: 220px;
  padding: 4px 8px;
  border-radius: 8px;
  border: 1px solid #d1d5db;
  font-size: 12px;
  background: #ffffff;
}

.mde-copy-btn {
  border-radius: 999px;
  border: none;
  padding: 4px 10px;
  font-size: 12px;
  cursor: pointer;
  background: #2563eb;
  color: #ffffff;
}

.mde-header-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.mde-status-box {
  min-width: 110px;
}

.mde-status {
  font-size: 12px;
}

.mde-status.saving {
  color: #f97316;
}

.mde-status.saved {
  color: #16a34a;
}

.mde-status.idle {
  color: #6b7280;
}

.mde-export-group {
  display: flex;
  gap: 4px;
}

.mde-btn {
  border-radius: 999px;
  border: none;
  padding: 4px 10px;
  font-size: 12px;
  cursor: pointer;
  background: #e5e7eb;
  color: #111827;
}

.mde-btn:disabled {
  opacity: 0.6;
  cursor: default;
}

.mde-back {
  border-radius: 999px;
  border: none;
  padding: 5px 12px;
  font-size: 12px;
  cursor: pointer;
  background: #f3f4f6;
  color: #374151;
}

.mde-error {
  margin: 8px 0 6px;
  font-size: 12px;
  color: #dc2626;
}

/* 主编辑区：占满卡片剩余空间 */
.mde-main {
  margin-top: 6px;
  flex: 1;
  display: flex;
  gap: 10px;
  min-height: 0; /* 这一句很关键，防止子元素撑破高度 */
}

.mde-pane {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.mde-pane-header {
  padding: 6px 10px;
  font-size: 12px;
  background: #f3f4f6;
  border-radius: 10px 10px 0 0;
  border-bottom: 1px solid #e5e7eb;
  color: #6b7280;
}

.mde-editor {
  flex: 1;
  border: 1px solid #e5e7eb;
  border-top: none;
  border-radius: 0 0 10px 10px;
  padding: 10px;
  font-family: "JetBrains Mono", "SF Mono", Menlo, Monaco, Consolas,
    "Courier New", monospace;
  font-size: 13px;
  outline: none;
  resize: none;
  background: #ffffff;
}

.mde-preview {
  flex: 1;
  border: 1px solid #e5e7eb;
  border-top: none;
  border-radius: 0 0 10px 10px;
  padding: 12px 14px;
  font-size: 14px;
  overflow: auto;
  background: #ffffff;
  line-height: 1.6;
}

/* Markdown 样式保持不变 */
.mde-preview h1 {
  font-size: 24px;
  margin: 16px 0 8px;
}
.mde-preview h2 {
  font-size: 20px;
  margin: 14px 0 8px;
}
.mde-preview h3 {
  font-size: 18px;
  margin: 12px 0 6px;
}
.mde-preview p {
  margin: 6px 0;
}
.mde-preview code {
  background: #f3f4f6;
  padding: 2px 4px;
  border-radius: 3px;
  font-family: "JetBrains Mono", monospace;
  font-size: 13px;
}
.mde-preview pre code {
  display: block;
  padding: 10px;
  overflow-x: auto;
}
.mde-preview ul,
.mde-preview ol {
  padding-left: 22px;
  margin: 6px 0;
}
.mde-preview blockquote {
  border-left: 4px solid #e5e7eb;
  margin: 6px 0;
  padding: 4px 10px;
  color: #6b7280;
  background: #f9fafb;
}
.mde-preview a {
  color: #2563eb;
  text-decoration: none;
}
.mde-preview a:hover {
  text-decoration: underline;
}

@media (max-width: 900px) {
  .mde-card {
    height: auto;
  }
  .mde-header {
    flex-direction: column;
    align-items: flex-start;
  }
  .mde-header-right {
    align-self: stretch;
    justify-content: flex-end;
  }
  .mde-main {
    flex-direction: column;
  }
}
</style>

