<template>
  <div class="container" :data-theme="theme">
    <!-- 左侧边栏 -->
    <div class="sidebar" :class="{ collapsed }">
      <!-- 顶部标题区 -->
      <div class="top-bar">
        <div class="title" @click="goHome">Moment</div>

        <a
          class="github"
          href="https://github.com/LeoninCS/Moment"
          target="_blank"
        >
          <svg viewBox="0 0 16 16" width="20" height="20">
            <path
              fill="currentColor"
              d="M8 0C3.58 0 0 3.58 0 8a8 8 0 005.47 7.59c.4.07.55-.17.55-.38 
              0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52
              -.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 
              2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 
              0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 
              0 0 .67-.21 2.2.82a7.6 7.6 0 012 0c1.53-1.03 2.2-.82 2.2-.82 
              .44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 
              0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 
              0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8 8 0 008 0z"
            />
          </svg>
        </a>

        <button class="toggle-btn" @click="collapsed = !collapsed">
          {{ collapsed ? "→" : "←" }}
        </button>
      </div>

      <!-- 搜索栏 -->
      <input
        v-if="!collapsed"
        v-model="search"
        type="text"
        class="search"
        placeholder="搜索功能..."
      />

      <!-- 菜单列表 -->
      <ul class="menu" v-if="!collapsed">
        <li
          v-for="item in filteredList"
          :key="item.name"
          :class="{ active: isActive(item) }"
          @click="selectFeature(item)"
        >
          {{ item.label }}
        </li>
      </ul>
    </div>

    <!-- 右侧内容区 -->
    <div class="content">
      <div class="header">
        <button class="theme-btn" @click="toggleTheme">
          {{ theme === "light" ? "🌙 深色" : "☀️ 浅色" }}
        </button>
      </div>

      <!-- 这里改成路由出口 -->
      <router-view />
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from "vue";
import { useRouter, useRoute } from "vue-router";

const router = useRouter();
const route = useRoute();

// 功能列表：用 path 来驱动路由
const features = [
  { name: "home", label: "关于本站", path: "/" },
  { name: "codeshare", label: "代码分享", path: "/codeshare" },
  { name: "workplan", label: "待办事项", path: "/workplan" },
  { name: "markdown", label: "Markdown 编辑器", path: "/markdown" },
];

const search = ref("");
const collapsed = ref(false);
const theme = ref("light"); // 默认主题

// 搜索过滤
const filteredList = computed(() => {
  if (!search.value) return features;
  return features.filter((f) =>
    f.label.toLowerCase().includes(search.value.toLowerCase())
  );
});

// 当前菜单是否激活
function isActive(item) {
  // 简单按 path 匹配；如果希望 /code/:hash 也高亮“代码分享”，可以在这里加逻辑
  return route.path === item.path;
}

// 点击功能项 —— 跳转路由
function selectFeature(item) {
  router.push(item.path);
}

// 回到首页
function goHome() {
  router.push("/");
}

// 切换主题
function toggleTheme() {
  theme.value = theme.value === "light" ? "dark" : "light";
}
</script>

<!-- 注意：这里不再加 scoped -->
<style>
/* =========================================
   布局基础
========================================= */
.container {
  display: flex;
  width: 100%;
  min-height: 100vh;      /* 新的：至少占满一屏 */
  /* height: 100%;        删掉 */
  /* overflow: hidden;    删掉，不然纵向滚动全没了 */
  background: var(--bg);
  color: var(--text);
}


/* =========================================
   主题系统
========================================= */
:root {
  --bg: #f5f6fa;
  --sidebar: #dad7d7;
  --border: #e5e7eb;
  --text: #374151;
  --hover: #eef2ff;
  --active: #3b82f6;
  --content-bg: #f0ebeb;
}

[data-theme="dark"] {
  --bg: #111827;
  --sidebar: #1f2933;
  --border: #374151;
  --text: #e5e7eb;
  --hover: #374151;
  --active: #2563eb;
  --content-bg: #090d1d;
}

/* =========================================
   左侧 Sidebar
========================================= */
.sidebar {
  width: 260px;
  background: var(--sidebar);
  border-right: 1px solid var(--border);
  padding: 16px;
  display: flex;
  flex-direction: column;
  transition: width 0.25s, padding 0.25s;
  box-shadow: 2px 0 6px rgba(0,0,0,0.05);
}

/* 折叠状态 */
.sidebar.collapsed {
  width: 64px;
  padding: 16px 8px;
}

/* 折叠时隐藏多余元素，只保留按钮 */
.sidebar.collapsed .title,
.sidebar.collapsed .github,
.sidebar.collapsed .search,
.sidebar.collapsed .menu {
  display: none;
}

.sidebar.collapsed .top-bar {
  justify-content: center;
}

/* 顶部区域 */
.top-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.title {
  font-size: 20px;
  font-weight: bold;
  cursor: pointer;
  color: var(--text);
}

.github {
  color: var(--text);
  margin-left: 10px;
}

.toggle-btn {
  margin-left: 10px;
  padding: 4px 6px;
  cursor: pointer;
  border: 1px solid var(--border);
  background: var(--sidebar);
  border-radius: 4px;
  color: var(--text);
}

/* 搜索框 */
.search {
  padding: 10px;
  border-radius: 8px;
  border: 1px solid var(--border);
  margin: 14px 0;
  outline: none;
  color: var(--text);
  background: var(--bg);
}

.search:focus {
  border-color: var(--active);
}

/* 菜单 */
.menu {
  list-style: none;
  padding: 0;
  margin: 0;
  overflow-y: auto;
  flex-grow: 1;
}

.menu li {
  padding: 12px 14px;
  margin-bottom: 6px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 15px;
  color: var(--text);
  transition: 0.2s;
}

.menu li:hover {
  background: var(--hover);
}

.menu li.active {
  background: var(--active);
  color: white;
}

/* =========================================
   右侧内容
========================================= */
.content {
  flex-grow: 1;
  padding: 26px;
  background: var(--content-bg);
  color: var(--text);
  overflow-y: auto;
  display: flex;
  flex-direction: column;
}

.header {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 16px;
}

.theme-btn {
  padding: 8px 14px;
  border-radius: 6px;
  border: 1px solid var(--border);
  cursor: pointer;
  background: var(--sidebar);
  color: var(--text);
  transition: 0.2s;
}

.theme-btn:hover {
  background: var(--hover);
}
</style>
