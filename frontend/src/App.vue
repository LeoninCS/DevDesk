<template>
  <div class="container" :data-theme="theme">
    <div
      v-if="isMobile && isMobileMenuOpen"
      class="sidebar-mask"
      @click="closeMobileMenu"
    ></div>

    <div class="sidebar" :class="sidebarClass">
      <div class="top-bar">
        <div class="title" @click="goHome">DevDesk</div>

        <a
          class="github"
          href="https://github.com/LeoninCS/Moment"
          target="_blank"
        >
          <svg viewBox="0 0 16 16" width="20" height="20" aria-hidden="true">
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

        <div class="sidebar-actions">
          <button class="toggle-btn" @click="toggleSidebar">
            {{
              isMobile
                ? isMobileMenuOpen
                  ? "收起"
                  : "展开"
                : sidebarCollapsed
                ? "展开"
                : "收起"
            }}
          </button>
        </div>
      </div>

      <input
        v-if="showSidebarContent"
        v-model="search"
        type="text"
        class="search"
        placeholder="搜索功能..."
      />

      <ul class="menu" v-if="showSidebarContent">
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

    <div class="content">
      <div class="header">
        <button
          v-if="isMobile"
          class="mobile-menu-btn"
          aria-label="展开菜单"
          @click="openMobileMenu"
        >
          ☰ 菜单
        </button>
        <button class="theme-btn" @click="toggleTheme">
          {{ theme === "light" ? "🌙 深色" : "☀️ 浅色" }}
        </button>
      </div>

      <router-view />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount } from "vue";
import { useRouter, useRoute } from "vue-router";

const router = useRouter();
const route = useRoute();

const features = [
  { name: "home", label: "关于本站", path: "/" },
  { name: "markdown", label: "Markdown 编辑", path: "/markdown" },
  { name: "httptest", label: "HTTP 接口测试", path: "/httptest" },
  { name: "workplan", label: "待办事项", path: "/workplan" },
  { name: "codeshare", label: "代码分享", path: "/codeshare" },
  { name: "htmlhost", label: "HTML 托管", path: "/html-host" },
];

const search = ref("");
const theme = ref("light");
const sidebarCollapsed = ref(false);
const isMobile = ref(false);
const isMobileMenuOpen = ref(false);

const filteredList = computed(() => {
  if (!search.value) return features;
  return features.filter((f) =>
    f.label.toLowerCase().includes(search.value.toLowerCase())
  );
});

const sidebarClass = computed(() => ({
  collapsed: sidebarCollapsed.value && !isMobile.value,
  mobile: isMobile.value,
  "mobile-open": isMobileMenuOpen.value,
}));

const showSidebarContent = computed(
  () =>
    (!sidebarCollapsed.value || isMobile.value) &&
    (!isMobile.value || isMobileMenuOpen.value)
);

function isActive(item) {
  return route.path === item.path;
}

function selectFeature(item) {
  router.push(item.path);
  closeMobileMenu();
}

function goHome() {
  router.push("/");
  closeMobileMenu();
}

function toggleTheme() {
  theme.value = theme.value === "light" ? "dark" : "light";
}

function toggleSidebar() {
  if (isMobile.value) {
    isMobileMenuOpen.value = !isMobileMenuOpen.value;
    return;
  }
  sidebarCollapsed.value = !sidebarCollapsed.value;
}

function openMobileMenu() {
  isMobileMenuOpen.value = true;
}

function closeMobileMenu() {
  if (isMobile.value) {
    isMobileMenuOpen.value = false;
  }
}

const handleResize = () => {
  isMobile.value = window.innerWidth <= 920;
  if (!isMobile.value) {
    isMobileMenuOpen.value = false;
  }
};

onMounted(() => {
  handleResize();
  window.addEventListener("resize", handleResize);
});

onBeforeUnmount(() => {
  window.removeEventListener("resize", handleResize);
});
</script>

<style>
.container {
  position: relative;
  display: flex;
  width: 100%;
  min-height: 100vh;
  background: var(--bg);
  color: var(--text);
}

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

.sidebar-mask {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.35);
  backdrop-filter: blur(1px);
  z-index: 8;
}

.sidebar {
  width: 260px;
  background: var(--sidebar);
  border-right: 1px solid var(--border);
  padding: 16px;
  display: flex;
  flex-direction: column;
  transition: width 0.25s, padding 0.25s, transform 0.25s ease;
  box-shadow: 2px 0 6px rgba(0, 0, 0, 0.05);
  z-index: 9;
}

.sidebar.mobile {
  position: fixed;
  inset: 0 auto 0 0;
  transform: translateX(-100%);
  max-width: 86vw;
  width: 280px;
  height: 100%;
  overflow-y: auto;
  border-right: none;
  box-shadow: 10px 0 30px rgba(0, 0, 0, 0.25);
}

.sidebar.mobile.mobile-open {
  transform: translateX(0);
}

.sidebar.collapsed:not(.mobile) {
  width: 64px;
  padding: 16px 8px;
}

.sidebar.collapsed:not(.mobile) .title,
.sidebar.collapsed:not(.mobile) .github,
.sidebar.collapsed:not(.mobile) .search,
.sidebar.collapsed:not(.mobile) .menu {
  display: none;
}

.sidebar.collapsed:not(.mobile) .top-bar {
  justify-content: center;
}

.top-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
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

.sidebar-actions {
  display: flex;
  align-items: center;
  gap: 6px;
}

.toggle-btn {
  padding: 6px 10px;
  cursor: pointer;
  border: 1px solid var(--border);
  background: var(--sidebar);
  border-radius: 6px;
  color: var(--text);
}

.toggle-btn.ghost {
  background: transparent;
}

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
  align-items: center;
  gap: 10px;
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

.mobile-menu-btn {
  display: none;
  padding: 8px 12px;
  border-radius: 6px;
  border: 1px solid var(--border);
  background: var(--sidebar);
  color: var(--text);
}

@media (max-width: 1100px) {
  .container {
    flex-direction: column;
  }

  .content {
    padding: 20px 18px;
  }
}

@media (max-width: 920px) {
  .content {
    padding: 16px 14px 22px;
  }

  .header {
    justify-content: space-between;
  }

  .mobile-menu-btn {
    display: inline-flex;
    align-items: center;
    gap: 6px;
  }
}
</style>
