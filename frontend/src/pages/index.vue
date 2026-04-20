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
  width: 50px;
  background: #ffffff;
  border-right: 1px solid #e8eaed;
  display: flex;
  flex-direction: column;
  align-items: center;
  box-shadow: 2px 0 8px rgba(0, 0, 0, 0.05);
  z-index: 100;
}

.sidebar-logo {
  padding: 20px 0;
  width: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  border-bottom: 1px solid #e8eaed;
  position: relative;
}

.logo-icon {
  font-size: 24px;
  background: linear-gradient(135deg, #409eff, #6c8fff);
  width: 36px;
  height: 36px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.sidebar-nav {
  flex: 1;
  width: 100%;
  padding: 8px 4px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  overflow: hidden;
}

.nav-item {
  width: 40px;
  height: 40px;
  padding: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  position: relative;
  color: #606266;
  transition: background-color 0.2s;
  text-decoration: none;
  border-radius: 8px;
}

.nav-item:hover {
  background-color: #f5f7fa;
  color: #409eff;
}

.nav-item.active {
  color: #409eff;
  background-color: #ecf5ff;
}

.nav-item.active::before {
  content: '';
  position: absolute;
  left: 50%;
  top: 0;
  transform: translateX(-50%);
  width: 20px;
  height: 3px;
  background-color: #409eff;
  border-radius: 0 0 2px 2px;
}

.nav-icon {
  font-size: 18px;
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.nav-tooltip {
  position: absolute;
  left: calc(100% + 8px);
  top: 50%;
  transform: translateY(-50%);
  background: rgba(0, 0, 0, 0.8);
  color: white;
  padding: 6px 10px;
  border-radius: 4px;
  font-size: 11px;
  white-space: nowrap;
  z-index: 1000;
  opacity: 0;
  pointer-events: none;
  transition: opacity 0.2s;
}

.nav-item:hover .nav-tooltip {
  opacity: 1;
}

.sidebar-footer {
  width: 100%;
  padding: 12px 0;
  border-top: 1px solid #e8eaed;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
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
            <div class="nav-tooltip">{{ item.text }}</div>
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