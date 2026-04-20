<style scoped>
.video-browse-container {
  width: 100%;
  padding: 20px;
  box-sizing: border-box;
  background: #f5f7fa;
  min-height: calc(100vh - 40px);
}

.page-title {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 20px;
  display: flex;
  align-items: center;
  gap: 8px;
}

/* 路径选择区 */
.path-section {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
  margin-bottom: 20px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
}

.path-input-group {
  display: flex;
  gap: 10px;
  align-items: center;
}

.path-input {
  flex: 1;
}

/* 工具栏 */
.toolbar {
  display: flex;
  gap: 12px;
  align-items: center;
  flex-wrap: wrap;
}

.toolbar-left {
  display: flex;
  gap: 12px;
  align-items: center;
}

.toolbar-right {
  display: flex;
  gap: 12px;
  align-items: center;
  margin-left: auto;
}

.toolbar-btn {
  padding: 6px 16px;
  border-radius: 4px;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  gap: 4px;
}

.toolbar-btn.primary {
  background: #409eff;
  border: none;
  color: #fff;
}

.toolbar-btn.primary:hover {
  background: #66b1ff;
}

.toolbar-btn.danger {
  background: #f56c6c;
  border: none;
  color: #fff;
}

.toolbar-btn.danger:hover {
  background: #f78989;
}

.toolbar-btn:disabled {
  background: #c0c4cc;
  cursor: not-allowed;
}

/* 视频统计 */
.stats-bar {
  display: flex;
  gap: 16px;
  margin-top: 12px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: #606266;
}

.stat-value {
  font-weight: 600;
  color: #409eff;
}

/* 视频网格 */
.video-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
  gap: 16px;
}

.video-card {
  background: #fff;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
  cursor: pointer;
  transition: transform 0.2s, box-shadow 0.2s;
  position: relative;
}

.video-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(64, 158, 255, 0.15);
}

.video-card.selected {
  box-shadow: 0 0 0 3px #409eff;
}

.video-card .checkbox {
  position: absolute;
  top: 8px;
  left: 8px;
  z-index: 10;
  width: 20px;
  height: 20px;
  border-radius: 4px;
  background: rgba(255, 255, 255, 0.9);
  border: 2px solid #409eff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  color: #409eff;
  opacity: 0;
  transition: opacity 0.2s;
}

.video-card:hover .checkbox,
.video-card.selected .checkbox {
  opacity: 1;
}

.video-card.selected .checkbox {
  background: #409eff;
  color: #fff;
}

.thumbnail-wrapper {
  position: relative;
  width: 100%;
  padding-top: 56.25%;
  background: #000;
  overflow: hidden;
}

.thumbnail-wrapper img {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.thumbnail-placeholder {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);
  color: #fff;
}

.thumbnail-placeholder .play-icon {
  font-size: 48px;
  margin-bottom: 8px;
  opacity: 0.8;
}

.video-info {
  padding: 12px;
}

.video-name {
  font-size: 13px;
  color: #303133;
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-bottom: 8px;
}

.video-meta {
  display: flex;
  justify-content: space-between;
  font-size: 11px;
  color: #909399;
}

.video-meta .resolution {
  color: #67c23a;
  font-weight: 500;
}

/* 加载状态 */
.loading-mask {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 14px;
}

/* 空状态 */
.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: #909399;
}

.empty-icon {
  font-size: 64px;
  margin-bottom: 16px;
}

.empty-text {
  font-size: 16px;
  margin-bottom: 8px;
}

.empty-hint {
  font-size: 13px;
  color: #c0c4cc;
}

/* 视频播放对话框 */
.video-dialog {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.85);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  z-index: 2000;
}

.video-dialog-header {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  padding: 16px 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: linear-gradient(to bottom, rgba(0,0,0,0.7), transparent);
  color: #fff;
}

.video-dialog-title {
  font-size: 16px;
  font-weight: 500;
  max-width: 70%;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.close-btn {
  background: rgba(255, 255, 255, 0.2);
  border: none;
  color: #fff;
  width: 36px;
  height: 36px;
  border-radius: 50%;
  font-size: 20px;
  cursor: pointer;
  transition: background 0.2s;
}

.close-btn:hover {
  background: rgba(255, 255, 255, 0.3);
}

.video-player {
  max-width: 90%;
  max-height: 80vh;
  border-radius: 8px;
  overflow: hidden;
}

.video-dialog-footer {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 16px 24px;
  background: linear-gradient(to top, rgba(0,0,0,0.7), transparent);
  color: #fff;
  display: flex;
  justify-content: space-around;
  flex-wrap: wrap;
  gap: 12px;
}

.video-detail-item {
  text-align: center;
}

.video-detail-label {
  font-size: 11px;
  opacity: 0.8;
  margin-bottom: 4px;
}

.video-detail-value {
  font-size: 14px;
  font-weight: 500;
}

.video-action-btn {
  background: rgba(255, 255, 255, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.3);
  color: #fff;
  padding: 6px 12px;
  border-radius: 4px;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s;
}

.video-action-btn:hover {
  background: rgba(255, 255, 255, 0.3);
}

.video-action-btn.danger {
  background: rgba(245, 108, 108, 0.8);
  border-color: #f56c6c;
}

.video-action-btn.danger:hover {
  background: #f56c6c;
}

/* 扫描进度 */
.scan-progress {
  margin-top: 12px;
}

.progress-bar {
  height: 4px;
  background: #e4e7ed;
  border-radius: 2px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #409eff, #67c23a);
  transition: width 0.3s;
}

.progress-text {
  font-size: 12px;
  color: #909399;
  margin-top: 6px;
}

/* 裁剪对话框 */
.clip-dialog {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 3000;
}

.clip-dialog-content {
  background: #fff;
  border-radius: 8px;
  padding: 24px;
  width: 400px;
  max-width: 90%;
  box-shadow: 0 12px 48px rgba(0, 0, 0, 0.3);
}

.clip-dialog-title {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 16px;
}

.clip-form-row {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
}

.clip-form-label {
  font-size: 13px;
  color: #606266;
  min-width: 70px;
}

.clip-form-hint {
  font-size: 11px;
  color: #909399;
  margin-left: 78px;
  margin-top: -8px;
  margin-bottom: 12px;
}

.clip-dialog-footer {
  display: flex;
  gap: 10px;
  justify-content: flex-end;
  margin-top: 20px;
}

/* 批量处理对话框 */
.batch-dialog {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 3000;
}

.batch-dialog-content {
  background: #fff;
  border-radius: 8px;
  padding: 24px;
  width: 420px;
  max-width: 90%;
  box-shadow: 0 12px 48px rgba(0, 0, 0, 0.3);
}

.batch-dialog-title {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 16px;
}

.batch-info {
  background: #f5f7fa;
  border-radius: 6px;
  padding: 12px;
  margin-bottom: 16px;
  font-size: 13px;
  color: #606266;
}

.batch-form-row {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
}

.batch-form-label {
  font-size: 13px;
  color: #606266;
  min-width: 80px;
}

.batch-dialog-footer {
  display: flex;
  gap: 10px;
  justify-content: flex-end;
  margin-top: 20px;
}

/* 处理日志 */
.process-log {
  background: #f8f9fa;
  border: 1px solid #e4e7ed;
  border-radius: 6px;
  padding: 12px;
  max-height: 200px;
  overflow-y: auto;
  font-size: 12px;
  line-height: 1.8;
  color: #303133;
  font-family: 'Consolas', 'Monaco', monospace;
  margin-top: 16px;
}

.process-log-item {
  word-break: break-all;
}

.process-log-item.success {
  color: #67c23a;
}

.process-log-item.fail {
  color: #f56c6c;
}
</style>

<template>
  <div class="video-browse-container">
    <div class="page-title">🎬 视频浏览</div>

    <!-- 路径选择 -->
    <div class="path-section">
      <div class="path-input-group">
        <el-input
            v-model="dirPath"
            placeholder="请选择要浏览的文件夹路径"
            class="path-input"
            clearable
        />
        <el-button type="primary" @click="selectDirectory">选择目录</el-button>
        <el-button
            type="success"
            :disabled="!dirPath || scanning"
            @click="startScan"
        >{{ scanning ? '扫描中...' : '开始扫描' }}</el-button>
      </div>

      <!-- 工具栏 -->
      <div v-if="videos.length > 0" class="toolbar" style="margin-top: 12px;">
        <div class="toolbar-left">
          <el-checkbox
              v-model="selectAll"
              :indeterminate="isIndeterminate"
              @change="toggleSelectAll"
          >全选</el-checkbox>
          <span style="font-size: 13px; color: #606266;">
            已选择 {{ selectedVideos.length }} / {{ videos.length }}
          </span>
        </div>
        <div class="toolbar-right">
          <el-input-number
              v-model="batchIntroDuration"
              :min="1"
              :max="3600"
              :step="1"
              size="small"
              style="width: 120px;"
              placeholder="片头时长"
          />
          <span style="font-size: 12px; color: #909399;">秒</span>
          <button
              class="toolbar-btn danger"
              :disabled="selectedVideos.length === 0 || batchProcessing"
              @click="openBatchRemoveIntroDialog"
          >
            🗑️ 批量去除片头
          </button>
          <button
              class="toolbar-btn primary"
              :disabled="selectedVideos.length === 0 || moveProcessing"
              @click="openMoveDialog"
          >
            📦 移动视频
          </button>
          <button
              class="toolbar-btn primary"
              :disabled="videos.length === 0 || archiveProcessing"
              @click="archiveMobileVideos"
          >
            📱 归档手机视频
          </button>
        </div>
      </div>

      <!-- 统计信息 -->
      <div v-if="videos.length > 0" class="stats-bar">
        <div class="stat-item">
          <span>总视频数:</span>
          <span class="stat-value">{{ videos.length }}</span>
        </div>
        <div class="stat-item">
          <span>已加载缩略图:</span>
          <span class="stat-value">{{ loadedThumbnails }} / {{ videos.length }}</span>
        </div>
      </div>

      <!-- 扫描进度 -->
      <div v-if="scanning || generatingThumbnails" class="scan-progress">
        <div class="progress-bar">
          <div class="progress-fill" :style="{ width: scanProgress + '%' }"></div>
        </div>
        <div class="progress-text">{{ scanStatus }}</div>
      </div>
    </div>

    <!-- 视频网格 -->
    <div v-if="videos.length > 0" class="video-grid">
      <div
          v-for="(video, idx) in videos"
          :key="idx"
          class="video-card"
          :class="{ selected: isSelected(video) }"
          @click="toggleSelect(video)"
      >
        <div class="checkbox">
          <span v-if="isSelected(video)">✓</span>
        </div>
        <div class="thumbnail-wrapper" @click.stop="playVideo(video)">
          <img
              v-if="video.thumbnail"
              :src="video.thumbnail"
              :alt="video.fileName"
          />
          <div v-else class="thumbnail-placeholder">
            <div class="play-icon">▶</div>
          </div>
          <div v-if="generatingThumbnails && !video.thumbnail" class="loading-mask">
            生成中...
          </div>
        </div>
        <div class="video-info">
          <div class="video-name" :title="video.fileName">{{ video.fileName }}</div>
          <div class="video-meta">
            <span class="resolution">{{ video.resolution || '未知' }}</span>
            <span>{{ video.fileSizeMB.toFixed(2) }} MB</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 空状态 -->
    <div v-else-if="!scanning && dirPath" class="empty-state">
      <div class="empty-icon">📂</div>
      <div class="empty-text">该目录下没有找到视频文件</div>
      <div class="empty-hint">请选择包含视频文件的文件夹</div>
    </div>

    <div v-else-if="!dirPath" class="empty-state">
      <div class="empty-icon">🎬</div>
      <div class="empty-text">请选择一个文件夹</div>
      <div class="empty-hint">将显示文件夹中的所有视频缩略图</div>
    </div>

    <!-- 视频播放对话框 -->
    <div v-if="playingVideo" class="video-dialog" @click.self="closeVideo">
      <div class="video-dialog-header">
        <div class="video-dialog-title">{{ playingVideo.fileName }}</div>
        <button class="close-btn" @click="closeVideo">×</button>
      </div>
      <video
          class="video-player"
          controls
          autoplay
          :src="playingVideo.videoUrl"
      ></video>
      <div class="video-dialog-footer">
        <div class="video-detail-item">
          <div class="video-detail-label">分辨率</div>
          <div class="video-detail-value">{{ playingVideo.resolution }}</div>
        </div>
        <div class="video-detail-item">
          <div class="video-detail-label">文件大小</div>
          <div class="video-detail-value">{{ playingVideo.fileSizeMB.toFixed(2) }} MB</div>
        </div>
        <div class="video-detail-item">
          <div class="video-detail-label">文件路径</div>
          <div class="video-detail-value" style="font-size: 11px; max-width: 200px; overflow: hidden; text-overflow: ellipsis;">
            {{ playingVideo.filePath }}
          </div>
        </div>
        <div style="display: flex; gap: 8px; align-items: center;">
          <button class="video-action-btn danger" @click.stop="openClipDialog">
            ✂️ 裁剪
          </button>
        </div>
      </div>
    </div>

    <!-- 裁剪对话框 -->
    <div v-if="showClipDialog" class="clip-dialog" @click.self="closeClipDialog">
      <div class="clip-dialog-content">
        <div class="clip-dialog-title">✂️ 裁剪视频（将替换原文件）</div>

        <div class="clip-form-row">
          <span class="clip-form-label">开始时间：</span>
          <el-input
              v-model="clipStartTime"
              placeholder="00:00:00"
              style="flex: 1"
              clearable
          />
        </div>
        <div class="clip-form-hint">格式：HH:MM:SS，例如 00:01:30</div>

        <div class="clip-form-row">
          <span class="clip-form-label">结束时间：</span>
          <el-input
              v-model="clipEndTime"
              placeholder="留空则裁剪到末尾"
              style="flex: 1"
              clearable
          />
        </div>
        <div class="clip-form-hint">格式：HH:MM:SS</div>

        <div style="font-size: 12px; color: #909399; margin-top: 8px;">
          裁剪范围：{{ clipStartTime || '起始' }} → {{ clipEndTime || '末尾' }}
        </div>

        <div class="clip-dialog-footer">
          <el-button @click="closeClipDialog">取消</el-button>
          <el-button
              type="danger"
              :loading="clipping"
              :disabled="!clipStartTime"
              @click="executeClip"
          >
            {{ clipping ? '裁剪中...' : '确认裁剪' }}
          </el-button>
        </div>
      </div>
    </div>

    <!-- 批量去除片头对话框 -->
    <div v-if="showBatchDialog" class="batch-dialog" @click.self="closeBatchDialog">
      <div class="batch-dialog-content">
        <div class="batch-dialog-title">🗑️ 批量去除片头</div>

        <div class="batch-info">
          即将处理 <strong>{{ selectedVideos.length }}</strong> 个视频文件
        </div>

        <div class="batch-form-row">
          <span class="batch-form-label">片头时长：</span>
          <el-input-number
              v-model="batchIntroDuration"
              :min="1"
              :max="3600"
              :step="1"
              style="flex: 1;"
          />
          <span>秒</span>
        </div>

        <div style="font-size: 12px; color: #909399; margin-top: 8px;">
          将从每个视频的前 {{ batchIntroDuration }} 秒处开始裁剪，替换原文件
        </div>

        <!-- 处理日志 -->
        <div v-if="batchLogs.length > 0" class="process-log">
          <div
              v-for="(log, idx) in batchLogs"
              :key="idx"
              class="process-log-item"
              :class="{ success: log.includes('✅'), fail: log.includes('❌') || log.includes('失败') }"
          >{{ log }}</div>
        </div>

        <div class="batch-dialog-footer">
          <el-button @click="closeBatchDialog" :disabled="batchProcessing">取消</el-button>
          <el-button
              type="danger"
              :loading="batchProcessing"
              :disabled="!batchIntroDuration || batchIntroDuration <= 0"
              @click="executeBatchRemoveIntro"
          >
            {{ batchProcessing ? '处理中...' : '开始处理' }}
          </el-button>
        </div>
      </div>
    </div>

    <!-- 批量移动视频对话框 -->
    <div v-if="showMoveDialog" class="batch-dialog" @click.self="closeMoveDialog">
      <div class="batch-dialog-content">
        <div class="batch-dialog-title">📦 批量移动视频</div>

        <div class="batch-info">
          即将移动 <strong>{{ selectedVideos.length }}</strong> 个视频文件到目标目录
        </div>

        <div class="batch-form-row">
          <span class="batch-form-label">目标目录：</span>
          <el-input
              v-model="moveTargetDir"
              placeholder="选择目标目录"
              style="flex: 1;"
              clearable
          />
          <el-button @click="selectMoveTargetDir">选择</el-button>
        </div>

        <div style="font-size: 12px; color: #909399; margin-top: 8px;">
          默认目标：上一级目录（父文件夹）
        </div>

        <!-- 处理日志 -->
        <div v-if="moveLogs.length > 0" class="process-log">
          <div
              v-for="(log, idx) in moveLogs"
              :key="idx"
              class="process-log-item"
              :class="{ success: log.includes('✅'), fail: log.includes('❌') || log.includes('失败') }"
          >{{ log }}</div>
        </div>

        <div class="batch-dialog-footer">
          <el-button @click="closeMoveDialog" :disabled="moveProcessing">取消</el-button>
          <el-button
              type="primary"
              :loading="moveProcessing"
              :disabled="!moveTargetDir"
              @click="executeMoveVideos"
          >
            {{ moveProcessing ? '移动中...' : '开始移动' }}
          </el-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import {SelectDirectory, ScanVideos, ExtractVideoThumbnail, ClipAndReplaceVideo, BatchRemoveVideoIntro, MoveFiles, ArchiveMobileVideos} from "../wailsjs/go/app/App.js";
import {EventsOn, EventsOff} from "../wailsjs/runtime/runtime.js";

export default {
  name: 'VideoBrowse',
  data() {
    return {
      dirPath: '',
      videos: [],
      scanning: false,
      scanProgress: 0,
      scanStatus: '',
      playingVideo: null,
      localServerPort: 39870,
      generatingThumbnails: false,
      loadedThumbnails: 0,
      // 裁剪相关
      showClipDialog: false,
      clipStartTime: '',
      clipEndTime: '',
      clipping: false,
      // 批量选择
      selectedVideos: [],
      // 批量去除片头
      showBatchDialog: false,
      batchIntroDuration: 30,
      batchProcessing: false,
      batchLogs: [],
      // 批量移动视频
      showMoveDialog: false,
      moveTargetDir: '',
      moveProcessing: false,
      moveLogs: [],
      // 归档手机视频
      archiveProcessing: false,
    };
  },
  computed: {
    selectAll() {
      return this.videos.length > 0 && this.selectedVideos.length === this.videos.length;
    },
    isIndeterminate() {
      return this.selectedVideos.length > 0 && this.selectedVideos.length < this.videos.length;
    },
  },
  methods: {
    async selectDirectory() {
      try {
        const path = await SelectDirectory();
        if (path) {
          this.dirPath = path;
          this.videos = [];
          this.playingVideo = null;
          this.selectedVideos = [];
        }
      } catch (e) {
        console.error('选择目录失败:', e);
        this.$message.error('选择目录失败：' + e);
      }
    },

    async startScan() {
      if (!this.dirPath) {
        this.$message.warning('请选择目录');
        return;
      }

      this.scanning = true;
      this.videos = [];
      this.scanProgress = 0;
      this.scanStatus = '正在扫描视频文件...';
      this.loadedThumbnails = 0;
      this.selectedVideos = [];

      ScanVideos({dirPath: this.dirPath})
          .then(res => {
            if (res.success) {
              this.videos = res.allVideos.map(v => ({
                ...v,
                thumbnail: null,
                videoUrl: null,
              }));
              this.scanProgress = 100;
              this.scanStatus = `扫描完成，共找到 ${this.videos.length} 个视频`;

              if (this.videos.length > 0) {
                this.generateThumbnails();
              }
            } else {
              this.$message.error('扫描失败：' + res.message);
              this.scanStatus = '扫描失败';
            }
          })
          .catch(e => {
            this.$message.error('扫描失败：' + e);
            this.scanStatus = '扫描失败';
          })
          .finally(() => {
            this.scanning = false;
          });
    },

    async generateThumbnails() {
      if (this.videos.length === 0 || this.generatingThumbnails) return;

      this.generatingThumbnails = true;
      this.loadedThumbnails = 0;

      for (let i = 0; i < this.videos.length; i++) {
        const video = this.videos[i];
        this.scanStatus = `正在生成缩略图 ${i + 1}/${this.videos.length}`;
        this.scanProgress = Math.round(((i + 1) / this.videos.length) * 100);

        try {
          const result = await ExtractVideoThumbnail({filePath: video.filePath});
          if (result.success && result.thumbnail) {
            video.thumbnail = `data:${result.mimeType};base64,${result.thumbnail}`;
            this.loadedThumbnails++;
          }
        } catch (e) {
          console.warn('缩略图生成失败:', video.fileName, e);
        }
      }

      this.generatingThumbnails = false;
      this.scanStatus = '缩略图生成完成';
    },

    isSelected(video) {
      return this.selectedVideos.some(v => v.filePath === video.filePath);
    },

    toggleSelect(video) {
      const idx = this.selectedVideos.findIndex(v => v.filePath === video.filePath);
      if (idx === -1) {
        this.selectedVideos.push(video);
      } else {
        this.selectedVideos.splice(idx, 1);
      }
    },

    toggleSelectAll(val) {
      if (val) {
        this.selectedVideos = [...this.videos];
      } else {
        this.selectedVideos = [];
      }
    },

    playVideo(video) {
      if (this.localServerPort > 0) {
        const encodedPath = encodeURIComponent(video.filePath);
        video.videoUrl = `http://localhost:${this.localServerPort}/file?path=${encodedPath}`;
      } else {
        video.videoUrl = `file://${video.filePath.replace(/\\/g, '/')}`;
      }
      this.playingVideo = video;
    },

    closeVideo() {
      this.playingVideo = null;
    },

    openClipDialog() {
      this.clipStartTime = '';
      this.clipEndTime = '';
      this.showClipDialog = true;
    },

    closeClipDialog() {
      if (!this.clipping) {
        this.showClipDialog = false;
      }
    },

    async executeClip() {
      if (!this.playingVideo) return;
      if (!this.clipStartTime) {
        this.$message.warning('请输入开始时间');
        return;
      }

      try {
        await this.$confirm('裁剪后将替换原视频，是否继续？', '确认裁剪', {
          confirmButtonText: '确认裁剪',
          cancelButtonText: '取消',
          type: 'warning',
        });
      } catch {
        return;
      }

      this.clipping = true;

      try {
        const res = await ClipAndReplaceVideo({
          filePath: this.playingVideo.filePath,
          startTime: this.clipStartTime,
          endTime: this.clipEndTime,
        });

        if (res.success) {
          this.$message.success('视频裁剪完成，已替换原文件！');
          this.showClipDialog = false;
          this.closeVideo();
          this.startScan();
        } else {
          this.$message.error('裁剪失败：' + res.message);
        }
      } catch (e) {
        this.$message.error('裁剪失败：' + e);
      } finally {
        this.clipping = false;
      }
    },

    openBatchRemoveIntroDialog() {
      if (this.selectedVideos.length === 0) {
        this.$message.warning('请先选择要处理的视频');
        return;
      }
      this.batchLogs = [];
      this.showBatchDialog = true;
    },

    closeBatchDialog() {
      if (!this.batchProcessing) {
        this.showBatchDialog = false;
      }
    },

    // 打开移动视频对话框
    openMoveDialog() {
      if (this.selectedVideos.length === 0) {
        this.$message.warning('请先选择要移动的视频');
        return;
      }
      // 默认目标目录为当前目录的上一级
      if (this.dirPath) {
        const parentPath = this.dirPath.replace(/\\$/, '').split('\\').slice(0, -1).join('\\');
        this.moveTargetDir = parentPath;
      }
      this.moveLogs = [];
      this.showMoveDialog = true;
    },

    closeMoveDialog() {
      if (!this.moveProcessing) {
        this.showMoveDialog = false;
      }
    },

    async selectMoveTargetDir() {
      try {
        const path = await SelectDirectory();
        if (path) {
          this.moveTargetDir = path;
        }
      } catch (e) {
        console.error('选择目录失败:', e);
        this.$message.error('选择目录失败：' + e);
      }
    },

    async executeMoveVideos() {
      if (this.selectedVideos.length === 0) return;
      if (!this.moveTargetDir) {
        this.$message.warning('请选择目标目录');
        return;
      }

      this.moveProcessing = true;
      this.moveLogs = [];

      // 监听后端日志
      const logHandler = (message) => {
        try {
          const msg = JSON.parse(message);
          this.moveLogs.push(msg.time + ' ' + msg.msg);
          this.$nextTick(() => {
            const logEl = document.querySelector('.process-log');
            if (logEl) {
              logEl.scrollTop = logEl.scrollHeight;
            }
          });
        } catch (e) {
          this.moveLogs.push(message);
        }
      };
      window.runtime.EventsOn('back_msg', logHandler);

      try {
        const filePaths = this.selectedVideos.map(v => v.filePath);
        const res = await MoveFiles({
          sourceFiles: filePaths,
          targetDir: this.moveTargetDir,
        });

        if (res.success && res.movedCount > 0) {
          this.$message.success(`移动完成！成功: ${res.movedCount}${res.failedFiles && res.failedFiles.length > 0 ? '，失败: ' + res.failedFiles.length : ''}`);
          this.showMoveDialog = false;
          this.selectedVideos = [];
          this.startScan();
        } else if (res.success && res.movedCount === 0) {
          this.$message.warning('没有文件被移动');
        } else {
          this.$message.error(res.message || '移动失败');
        }
      } catch (e) {
        this.$message.error('移动失败：' + e);
      } finally {
        this.moveProcessing = false;
        window.runtime.EventsOff('back_msg', logHandler);
      }
    },

    async executeBatchRemoveIntro() {
      if (this.selectedVideos.length === 0) return;
      if (!this.batchIntroDuration || this.batchIntroDuration <= 0) {
        this.$message.warning('请输入有效的片头时长');
        return;
      }

      this.batchProcessing = true;
      this.batchLogs = [];

      // 监听后端日志
      const logHandler = (message) => {
        try {
          const msg = JSON.parse(message);
          this.batchLogs.push(msg.time + ' ' + msg.msg);
          this.$nextTick(() => {
            const logEl = document.querySelector('.process-log');
            if (logEl) {
              logEl.scrollTop = logEl.scrollHeight;
            }
          });
        } catch (e) {
          this.batchLogs.push(message);
        }
      };
      window.runtime.EventsOn('back_msg', logHandler);

      try {
        const filePaths = this.selectedVideos.map(v => v.filePath);
        const res = await BatchRemoveVideoIntro({
          filePaths: filePaths,
          introDuration: this.batchIntroDuration,
        });

        if (res.processed > 0) {
          this.$message.success(`批量处理完成！成功: ${res.processed}，失败: ${res.failed}`);
          this.showBatchDialog = false;
          this.selectedVideos = [];
          this.startScan();
        } else {
          this.$message.error('批量处理失败');
        }
      } catch (e) {
        this.$message.error('批量处理失败：' + e);
      } finally {
        this.batchProcessing = false;
        window.runtime.EventsOff('back_msg', logHandler);
      }
    },

    async archiveMobileVideos() {
      if (this.videos.length === 0) {
        this.$message.warning('没有可归档的视频');
        return;
      }

      try {
        await this.$confirm('将把所有竖屏视频（宽<高）移动到同目录下的 mobile 文件夹，是否继续？', '确认归档', {
          confirmButtonText: '确认',
          cancelButtonText: '取消',
          type: 'warning',
        });
      } catch {
        return;
      }

      this.archiveProcessing = true;

      // 监听后端日志
      const logHandler = (message) => {
        try {
          const msg = JSON.parse(message);
          console.log(msg.time + ' ' + msg.msg);
        } catch (e) {
          console.log(message);
        }
      };
      window.runtime.EventsOn('back_msg', logHandler);

      try {
        const res = await ArchiveMobileVideos({dirPath: this.dirPath});

        if (res.success) {
          this.$message.success(res.message);
          this.startScan();
        } else {
          this.$message.error('归档失败：' + res.message);
        }
      } catch (e) {
        this.$message.error('归档失败：' + e);
      } finally {
        this.archiveProcessing = false;
        window.runtime.EventsOff('back_msg', logHandler);
      }
    },
  },
  beforeUnmount() {
    this.closeVideo();
  },
};
</script>
