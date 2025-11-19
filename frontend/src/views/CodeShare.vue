<template>
  <div class="codeshare">
    <h1>代码分享</h1>
    <p class="desc">在这里粘贴你的代码，生成一个可分享的片段。（已接入后端）</p>

    <form @submit.prevent="handleSubmit">
      <!-- 上面三列：发布者 / 语法 / 过期时间 -->
      <div class="row">
        <div class="field">
          <label>发布者</label>
          <input
            v-model="form.poster"
            placeholder="你的昵称"
            maxlength="30"
          />
          <small>最多 30 个字符</small>
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

      <!-- 内容区 -->
      <div class="field">
        <label>代码内容</label>
        <textarea
          v-model="form.content"
          rows="14"
          placeholder="在这里输入或粘贴代码……"
        />
      </div>

      <!-- 提交区 -->
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
  { value: "c_cpp", label: "C / C++" }
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
  content: ""
});

const loading = ref(false);
const shareUrl = ref("");

// 过期时间映射：把选择的项换成秒数
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
    const now = Math.floor(Date.now() / 1000); // 当前时间戳（秒）
    const delta = expirationSecondsMap[form.value.expiration] || (24 * 60 * 60);
    const destroyTime = now + delta; // 过期时间时间戳

    // 使用封装好的 API
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

    // 自动跳转到查看页
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
.codeshare {
  /* 宽度按百分比，左右留一点边距 */
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 4vw; /* 4% 视口宽度边距 */
  box-sizing: border-box;
}

h1 {
  margin: 0 0 8px;
  font-size: 24px;
}

.desc {
  margin: 0 0 24px;
  font-size: 14px;
  color: #6b7280;
}

/* 表单整体卡片，宽度跟随父容器百分比 */
form {
  width: 100%;
  background: rgba(255, 255, 255, 0.98);
  padding: 20px 24px;
  border-radius: 8px;
  border: 1px solid #e5e7eb;
}

/* 深色模式下的适配 */
:global([data-theme="dark"]) .codeshare form {
  background: #111827;
  border-color: #374151;
}

/* 顶部三列：桌面端三列布局 */
.row {
  display: grid;
  grid-template-columns: 1.2fr 1fr 1fr;
  gap: 16px;
  margin-bottom: 16px;
}

.field {
  display: flex;
  flex-direction: column;
  margin-bottom: 16px;
}

.field label {
  font-size: 14px;
  margin-bottom: 6px;
}

.field small {
  margin-top: 4px;
  font-size: 12px;
  color: #9ca3af;
}

input,
select,
textarea {
  padding: 8px 10px;
  border-radius: 4px;
  border: 1px solid #d1d5db;
  font: inherit;
  width: 100%;        /* 控件宽度占满父容器 */
  box-sizing: border-box;
}

/* 深色下输入框 */
:global([data-theme="dark"]) input,
:global([data-theme="dark"]) select,
:global([data-theme="dark"]) textarea {
  background: #111827;
  border-color: #4b5563;
  color: #e5e7eb;
}

/* 文本框高度按视口百分比，手机也能比较好用 */
textarea {
  resize: vertical;
  min-height: 30vh;   /* 至少占视口高度 30% */
}

/* 底部按钮区域 */
.actions {
  margin-top: 12px;
  display: flex;
  align-items: center;
  gap: 16px;
  flex-wrap: wrap;
}

.actions button {
  background: #16a34a;
  color: #fff;
  border-radius: 4px;
  padding: 8px 16px;
  border: none;
  cursor: pointer;
}

.actions button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.share-url {
  font-size: 14px;
}

.share-url a {
  color: #2563eb;
}

/* ===========================
   响应式：平板 / 手机适配
   =========================== */

@media (max-width: 768px) {
  .codeshare {
    padding: 0 3vw;
  }

  form {
    padding: 16px;
  }

  .row {
    grid-template-columns: 1fr;
  }

  h1 {
    font-size: 20px;
    text-align: center;
  }

  .desc {
    text-align: center;
  }

  textarea {
    min-height: 35vh;
  }

  .actions {
    flex-direction: column;
    align-items: stretch;
  }

  .actions button {
    width: 100%;
    text-align: center;
  }
}

@media (max-width: 480px) {
  form {
    padding: 12px;
  }

  .codeshare {
    padding: 0 2vw;
  }
}
</style>
