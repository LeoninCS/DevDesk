<template>
  <div class="csv-page">
    <section class="csv-hero">
      <div>
        <p class="eyebrow">CodeShare</p>
        <h1>分享给队友的代码片段</h1>
        <p class="sub">
          打开链接即可查看代码。支持语法高亮、主题切换、复制代码和快速分享当前页面。
        </p>
        <div class="hero-actions">
          <button class="btn ghost" @click="copyLink">
            {{ linkCopied ? "已复制链接" : "复制当前链接" }}
          </button>
          <button class="btn ghost" @click="copyCode">
            {{ copied ? "已复制代码" : "复制代码" }}
          </button>
        </div>
      </div>
      <div class="hero-meta">
        <div class="meta-row">
          <span class="label">发布者</span>
          <span class="value">{{ author || "匿名" }}</span>
        </div>
        <div class="meta-row">
          <span class="label">语法</span>
          <span class="value">{{ language || "纯文本" }}</span>
        </div>
        <div class="meta-row">
          <span class="label">主题</span>
          <button class="btn inline" @click="toggleCodeTheme">
            {{ codeTheme === "dark" ? "切换为浅色" : "切换为深色" }}
          </button>
        </div>
      </div>
    </section>

    <section class="csv-card" :class="`csv-card--${codeTheme}`">
      <header class="card-head">
        <div>
          <h2>代码内容</h2>
          <p class="muted">已高亮显示，支持手动复制。</p>
        </div>
        <span v-if="language" class="lang-tag">{{ language }}</span>
      </header>

      <div v-if="loading" class="status">加载中...</div>
      <div v-else-if="error" class="status error">{{ error }}</div>

      <div v-else class="code-wrap">
        <pre class="pre-wrap">
          <code
            ref="codeEl"
            :class="['code-block', `language-${languageClass}`, `code-block--${codeTheme}`]"
          >{{ content }}</code>
        </pre>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, nextTick } from "vue";
import { useRoute } from "vue-router";
import hljs from "highlight.js";
import "highlight.js/styles/github-dark.css";

import { getCodeByHash } from "../api/codeshare.ts";

const route = useRoute();
const hash = route.params.hash as string;

const loading = ref(true);
const error = ref("");
const author = ref("");
const language = ref("");
const content = ref("");

const codeEl = ref<HTMLElement | null>(null);

const codeTheme = ref<"dark" | "light">("dark");
const copied = ref(false);
const linkCopied = ref(false);

const languageClass = computed(() => {
  if (!language.value) return "plaintext";
  if (language.value === "c_cpp") return "cpp";
  return language.value;
});

onMounted(async () => {
  try {
    const res = await getCodeByHash(hash);

    author.value = res.data.author || "";
    language.value = res.data.language || "";
    content.value = res.data.content || "";
  } catch (e) {
    console.error(e);
    error.value = "该代码不存在或已过期。";
  } finally {
    loading.value = false;
  }

  await nextTick();
  if (codeEl.value) {
    hljs.highlightElement(codeEl.value);
  }
});

function toggleCodeTheme() {
  codeTheme.value = codeTheme.value === "dark" ? "light" : "dark";
}

async function copyCode() {
  const text = content.value || "";
  if (!text) return;

  const fallbackCopy = () => {
    try {
      const textarea = document.createElement("textarea");
      textarea.value = text;
      textarea.style.position = "fixed";
      textarea.style.opacity = "0";
      textarea.style.left = "-9999px";
      document.body.appendChild(textarea);
      textarea.focus();
      textarea.select();
      const ok = document.execCommand("copy");
      document.body.removeChild(textarea);
      return ok;
    } catch (err) {
      console.error("fallback copy error", err);
      return false;
    }
  };

  try {
    if (
      typeof navigator !== "undefined" &&
      navigator.clipboard &&
      typeof navigator.clipboard.writeText === "function"
    ) {
      await navigator.clipboard.writeText(text);
    } else {
      const ok = fallbackCopy();
      if (!ok) throw new Error("fallback copy failed");
    }

    copied.value = true;
    setTimeout(() => {
      copied.value = false;
    }, 1500);
  } catch (e) {
    console.error(e);
    alert("复制失败，请手动选择代码复制");
  }
}

async function copyLink() {
  const url = typeof window !== "undefined" ? window.location.href : "";
  if (!url) return;

  const fallbackCopy = () => {
    try {
      const textarea = document.createElement("textarea");
      textarea.value = url;
      textarea.style.position = "fixed";
      textarea.style.opacity = "0";
      textarea.style.left = "-9999px";
      document.body.appendChild(textarea);
      textarea.focus();
      textarea.select();
      const ok = document.execCommand("copy");
      document.body.removeChild(textarea);
      return ok;
    } catch (err) {
      console.error("fallback copy link error", err);
      return false;
    }
  };

  try {
    if (
      typeof navigator !== "undefined" &&
      navigator.clipboard &&
      typeof navigator.clipboard.writeText === "function"
    ) {
      await navigator.clipboard.writeText(url);
    } else {
      const ok = fallbackCopy();
      if (!ok) throw new Error("fallback copy link failed");
    }

    linkCopied.value = true;
    setTimeout(() => {
      linkCopied.value = false;
    }, 1500);
  } catch (e) {
    console.error(e);
    alert("复制链接失败，请手动复制地址栏链接");
  }
}
</script>

<style scoped>
.csv-page {
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
  padding: 10px 6px 26px;
  display: flex;
  flex-direction: column;
  gap: 14px;
  box-sizing: border-box;
}

.csv-hero {
  background: linear-gradient(120deg, #0ea5e9 0%, #2563eb 50%, #111827 100%);
  color: #f8fafc;
  border-radius: 16px;
  padding: 20px 18px;
  display: grid;
  grid-template-columns: 1.4fr 0.8fr;
  gap: 16px;
  align-items: center;
}

.eyebrow {
  margin: 0;
  text-transform: uppercase;
  letter-spacing: 1px;
  font-size: 12px;
  opacity: 0.8;
}

.csv-hero h1 {
  margin: 6px 0 6px;
  font-size: 24px;
}

.sub {
  margin: 0 0 10px;
  line-height: 1.5;
  max-width: 720px;
}

.hero-actions {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.hero-meta {
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 14px;
  padding: 12px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.meta-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 8px;
}

.label {
  font-size: 12px;
  opacity: 0.8;
}

.value {
  font-weight: 700;
}

.csv-card {
  border-radius: 16px;
  padding: 16px 16px 18px;
  border: 1px solid #1f2937;
  overflow: hidden;
}

.csv-card--dark {
  background: #0f172a;
  color: #e5e7eb;
}

.csv-card--light {
  background: #f9fafb;
  color: #111827;
  border-color: #d1d5db;
}

.card-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 8px;
}

.card-head h2 {
  margin: 0;
  font-size: 20px;
}

.muted {
  margin: 2px 0 0;
  color: #9ca3af;
  font-size: 12px;
}

.lang-tag {
  padding: 6px 10px;
  border-radius: 999px;
  border: 1px solid currentColor;
  font-size: 12px;
}

.status {
  padding: 12px 0;
  font-size: 14px;
  color: #cbd5e1;
}

.status.error {
  color: #fca5a5;
}

.code-wrap {
  margin-top: 10px;
}

.pre-wrap {
  margin: 0;
  font-size: 13px;
  line-height: 1.5;
  max-height: 72vh;
  overflow: auto;
  background: transparent;
}

.code-block {
  display: block;
  white-space: pre;
}

.code-block--dark {
  background: #111827 !important;
  color: #f8fafc !important;
}

.code-block--light {
  background: #e5e7eb !important;
  color: #111827 !important;
}

.btn {
  padding: 8px 12px;
  border-radius: 10px;
  border: 1px solid rgba(255, 255, 255, 0.4);
  background: transparent;
  color: #f8fafc;
  cursor: pointer;
  font-weight: 600;
}

.btn.ghost {
  backdrop-filter: blur(6px);
}

.btn.inline {
  border-color: rgba(255, 255, 255, 0.5);
  color: #f8fafc;
}

:global([data-theme="light"]) .btn.inline {
  color: #111827;
  border-color: #d1d5db;
}

@media (max-width: 900px) {
  .csv-hero {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .csv-card {
    padding: 12px;
  }
}
</style>
