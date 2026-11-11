<style>
html, body {
  margin: 0;
  padding: 0;
  height: 100%;
  overflow: hidden;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "PingFang SC", "Microsoft YaHei", sans-serif;
}

#app-layout {
  display: flex;
  flex-direction: column;
  height: 100vh;
  width: 100vw;
}

.main-wrapper {
  display: flex;
  flex: 1;
  overflow: hidden;
}

.sidebar {
  width: 170px;
  background: #ffffff;
  border-right: 1px solid #e8eaed;
  display: flex;
  flex-direction: column;
  box-shadow: 2px 0 8px rgba(0, 0, 0, 0.05);
  z-index: 100;
  overflow-x: hidden;
}

.sidebar-nav {
  flex: 1;
  width: 100%;
  padding: 8px 0;
  overflow-y: auto;
  overflow-x: hidden;
}

.nav-item {
  width: 100%;
  padding: 8px 12px;
  display: flex;
  align-items: center;
  cursor: pointer;
  color: #606266;
  transition: background-color 0.2s;
  text-decoration: none;
}

.nav-item:hover {
  background-color: #f5f7fa;
  color: #409eff;
}

.nav-item.active {
  color: #409eff;
  background-color: #ecf5ff;
}

.nav-icon {
  font-size: 18px;
  margin-right: 8px;
}

.nav-text {
  font-size: 12px;
}

.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background: #f5f7fa;
}

.content-body {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
}
</style>

<template>
  <div id="app-layout">
    <div class="main-wrapper">
      <div class="sidebar">
        <div class="sidebar-nav">
          <div
            v-for="item in allTools"
            :key="item.link"
            class="nav-item"
            :class="{ active: isActive(item.link) }"
            @click="navigateTo(item.link)"
          >
            <div class="nav-icon">{{ item.emoji }}</div>
            <div class="nav-text">{{ item.text }}</div>
          </div>
        </div>
      </div>
      <div class="main-content">
        <div class="content-body">
          <router-view></router-view>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { computed, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { ReadAssetFile } from "../wailsjs/go/app/App.js";

export default {
  name: 'AppLayout',
  setup() {
    const route = useRoute();
    const router = useRouter();
    const allTools = ref([]);

    const normalizeLink = (link) => {
      if (!link) return '';
      return link.startsWith('/') ? link.substring(1) : link;
    };

    const normalizePath = (path) => {
      if (path === '/') return '/';
      return path.startsWith('/') ? path.substring(1) : path;
    };

    const currentRoute = computed(() => route.path);

    const isActive = (link) => {
      if (!link) return false;
      return normalizePath(route.path) === normalizeLink(link);
    };

    const loadTools = async () => {
      try {
        const menuData = await ReadAssetFile("/config/menu.json");
        const tools = JSON.parse(menuData);
        allTools.value = tools.map(tool => ({
          ...tool,
          link: normalizeLink(tool.link)
        }));
      } catch (error) {
        console.error('加载菜单失败:', error);
      }
    };

    const navigateTo = (link) => {
      const normalized = normalizeLink(link);
      router.push('/' + normalized);
    };

    loadTools();

    return {
      allTools,
      currentRoute,
      isActive,
      navigateTo,
    };
  }
};
</script>