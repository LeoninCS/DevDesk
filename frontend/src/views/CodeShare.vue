<template>
  <div class="cs-page">
    <section class="cs-hero">
      <div>
        <p class="eyebrow">CodeShare</p>
        <h1>生成可分享的代码片段</h1>
        <p class="sub">
          贴上代码，选择语法高亮和过期时间，一键获取可分享的链接。已接入后端，分享出去即可在线查看。
        </p>
        <div class="chips">
          <span class="chip">语法高亮</span>
          <span class="chip">临时/短链</span>
          <span class="chip">便捷复制</span>
        </div>
      </div>
      <div class="hero-card">
        <div class="hero-stat">
          <p class="hero-label">默认过期</p>
          <p class="hero-value">24h</p>
        </div>
        <div class="hero-divider"></div>
        <div class="hero-stat">
          <p class="hero-label">支持</p>
          <p class="hero-value">多语言</p>
        </div>
      </div>
    </section>

    <section class="cs-card">
      <header class="cs-card-head">
        <div>
          <p class="eyebrow">创建分享</p>
          <h2>填写信息并生成链接</h2>
          <p class="muted">支持作者昵称、语法选择和自定义过期时间。</p>
        </div>
        <div v-if="shareUrl" class="share-chip">
          <span>已生成链接</span>
          <a :href="shareUrl" target="_blank">{{ shareUrl }}</a>
        </div>
      </header>

      <form @submit.prevent="handleSubmit">
        <div class="row">
          <div class="field">
            <label>发布者</label>
            <input
              v-model="form.poster"
              placeholder="你的昵称"
              maxlength="30"
            />
            <small>可为空，最长 30 个字符</small>
          </div>

          <div class="field">
            <label>语法高亮</label>
            <select v-model="form.syntax">
              <option
                v-for="opt in syntaxes"
                :key="opt.value"
                :value="opt.value"
              >
                {{ opt.label }}
              </option>
            </select>
          </div>

          <div class="field">
            <label>过期时间</label>
            <select v-model="form.expiration">
              <option
                v-for="opt in expirations"
                :key="opt.value"
                :value="opt.value"
              >
                {{ opt.label }}
              </option>
            </select>
            <small>大致时间，实际以后端配置为准</small>
          </div>
        </div>

        <div class="field">
          <label>代码内容</label>
          <textarea
            v-model="form.content"
            rows="14"
            placeholder="在这里输入或粘贴代码..."
          />
        </div>

        <div class="actions">
          <button type="submit" :disabled="!form.content.trim() || loading">
            {{ loading ? "提交中..." : "生成分享链接" }}
          </button>

          <span v-if="shareUrl" class="share-url">
            分享链接：
            <a :href="shareUrl" target="_blank">{{ shareUrl }}</a>
          </span>
        </div>
      </form>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";
import { uploadCode } from "../api/codeshare.ts";

const router = useRouter();

const syntaxes = [
  { value: "plaintext", label: "纯文本" },
  { value: "javascript", label: "JavaScript" },
  { value: "typescript", label: "TypeScript" },
  { value: "vue", label: "Vue" },
  { value: "python", label: "Python" },
  { value: "go", label: "Go" },
  { value: "c_cpp", label: "C / C++" },
];

const expirations = [
  { value: "10min", label: "10 分钟" },
  { value: "1hour", label: "1 小时" },
  { value: "1day", label: "1 天" },
  { value: "1week", label: "1 周" },
  { value: "1month", label: "1 个月" },
];

const form = ref({
  poster: "",
  syntax: "plaintext",
  expiration: "1day",
  content: "",
});

const loading = ref(false);
const shareUrl = ref("");

const expirationSecondsMap: Record<string, number> = {
  "10min": 10 * 60,
  "1hour": 60 * 60,
  "1day": 24 * 60 * 60,
  "1week": 7 * 24 * 60 * 60,
  "1month": 30 * 24 * 60 * 60,
};

async function handleSubmit() {
  if (!form.value.content.trim()) return;

  loading.value = true;
  shareUrl.value = "";

  try {
    const now = Math.floor(Date.now() / 1000);
    const delta = expirationSecondsMap[form.value.expiration] || 24 * 60 * 60;
    const destroyTime = now + delta;

    const res = await uploadCode({
      author: form.value.poster,
      language: form.value.syntax,
      content: form.value.content,
      destroy_time: destroyTime,
    });

    const hash = res.data.hash;
    const urlFromServer = res.data.url;
    const finalUrl = urlFromServer || `${window.location.origin}/c/${hash}`;

    shareUrl.value = finalUrl;
    router.push(`/code/${hash}`);
  } catch (err) {
    console.error(err);
    alert("提交失败，请检查后端是否已启动 /api/codeshare/upload");
  } finally {
    loading.value = false;
  }
}
</script>

<style scoped>
.cs-page {
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  gap: 16px;
  padding: 8px 6px 24px;
  box-sizing: border-box;
}

.cs-hero {
  background: linear-gradient(120deg, #0ea5e9 0%, #2563eb 40%, #111827 100%);
  color: #f8fafc;
  border-radius: 16px;
  padding: 20px 18px;
  display: grid;
  grid-template-columns: 1.5fr 1fr;
  gap: 14px;
  align-items: center;
}

.eyebrow {
  margin: 0;
  text-transform: uppercase;
  letter-spacing: 1px;
  font-size: 12px;
  opacity: 0.8;
}

.cs-hero h1 {
  margin: 6px 0 4px;
  font-size: 26px;
}

.sub {
  margin: 0 0 10px;
  max-width: 760px;
  line-height: 1.5;
}

.chips {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.chip {
  padding: 6px 10px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.12);
  border: 1px solid rgba(255, 255, 255, 0.3);
  font-size: 12px;
  color: #e2e8f0;
}

.hero-card {
  display: flex;
  align-items: center;
  justify-content: space-evenly;
  background: rgba(255, 255, 255, 0.08);
  border-radius: 14px;
  padding: 14px 10px;
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.hero-stat {
  text-align: center;
}

.hero-label {
  margin: 0;
  font-size: 12px;
  opacity: 0.8;
}

.hero-value {
  margin: 4px 0 0;
  font-size: 20px;
  font-weight: 800;
}

.hero-divider {
  width: 1px;
  height: 36px;
  background: rgba(255, 255, 255, 0.25);
}

.cs-card {
  background: var(--moment-card-bg, #ffffff);
  border-radius: 16px;
  border: 1px solid #e5e7eb;
  box-shadow: 0 12px 38px rgba(15, 23, 42, 0.08);
  padding: 18px 18px 22px;
  box-sizing: border-box;
}

.cs-card-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 10px;
  margin-bottom: 10px;
}

.cs-card-head h2 {
  margin: 4px 0 4px;
  font-size: 22px;
}

.muted {
  margin: 0;
  color: #6b7280;
  font-size: 13px;
}

.share-chip {
  background: #ecfdf3;
  border: 1px solid #bbf7d0;
  color: #15803d;
  padding: 10px 12px;
  border-radius: 12px;
  font-size: 12px;
  max-width: 320px;
}

.share-chip a {
  display: block;
  color: inherit;
  text-decoration: none;
  font-weight: 600;
  margin-top: 4px;
  word-break: break-all;
}

form {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.row {
  display: grid;
  grid-template-columns: 1.2fr 1fr 1fr;
  gap: 16px;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.field label {
  font-size: 14px;
  font-weight: 600;
  color: #111827;
}

.field small {
  font-size: 12px;
  color: #9ca3af;
}

input,
select,
textarea {
  padding: 10px 12px;
  border-radius: 10px;
  border: 1px solid #d1d5db;
  font: inherit;
  width: 100%;
  box-sizing: border-box;
  background: #ffffff;
}

textarea {
  resize: vertical;
  min-height: 32vh;
  font-family: "JetBrains Mono", Consolas, monospace;
}

.actions {
  margin-top: 6px;
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.actions button {
  background: linear-gradient(120deg, #2563eb, #4f46e5);
  color: #fff;
  border-radius: 12px;
  padding: 10px 18px;
  border: none;
  cursor: pointer;
  font-weight: 700;
  box-shadow: 0 10px 24px rgba(79, 70, 229, 0.35);
}

.actions button:disabled {
  opacity: 0.65;
  cursor: not-allowed;
  box-shadow: none;
}

.share-url {
  font-size: 13px;
  color: #374151;
}

.share-url a {
  color: #2563eb;
  word-break: break-all;
}

:global([data-theme="dark"]) .cs-card {
  background: #0b1222;
  border-color: #1f2a44;
}

:global([data-theme="dark"]) input,
:global([data-theme="dark"]) select,
:global([data-theme="dark"]) textarea {
  background: #0f172a;
  border-color: #334155;
  color: #e5e7eb;
}

@media (max-width: 960px) {
  .cs-hero {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .row {
    grid-template-columns: 1fr;
  }

  .cs-card-head {
    flex-direction: column;
    align-items: flex-start;
  }

  .hero-card {
    justify-content: center;
  }
}
</style>
