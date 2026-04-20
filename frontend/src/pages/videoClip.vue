<style scoped>
.clip-container {
  width: 100%;
  padding: 10px;
  box-sizing: border-box;
  overflow-x: hidden;
}

.form-row {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 12px;
  flex-wrap: wrap;
}

.label {
  font-size: 14px;
  color: #606266;
  min-width: 80px;
  flex-shrink: 0;
}

.main-layout {
  display: flex;
  gap: 16px;
  margin-top: 10px;
  width: 100%;
  box-sizing: border-box;
  overflow: hidden;
}

.left-panel {
  flex: 1 1 0;
  min-width: 0;
  overflow: hidden;
}

.right-panel {
  flex: 0 0 380px;
  min-width: 0;
  width: 380px;
  overflow: hidden;
}

.video-wrapper {
  width: 100%;
  background: #000;
  border-radius: 0 0 6px 6px;
  overflow: hidden;
  cursor: pointer;
  position: relative;
  box-sizing: border-box;
}

.video-wrapper video {
  width: 100%;
  max-width: 100%;
  display: block;
  max-height: 525px;
  object-fit: contain;
}

.time-display {
  font-size: 13px;
  color: #909399;
  margin-top: 4px;
  text-align: center;
}

.time-input-group {
  display: flex;
  align-items: center;
  gap: 8px;
}

.time-btn {
  white-space: nowrap;
}

.log-area {
  border: 1px solid #e4e7ed;
  border-radius: 4px;
  padding: 8px;
  background: #f8f9fa;
  height: 160px;
  overflow-y: auto;
  font-size: 12px;
  line-height: 1.8;
  color: #303133;
}

.result-card {
  margin-top: 12px;
  padding: 10px;
  border-radius: 6px;
  border: 1px solid #67c23a;
  background: #f0f9eb;
  font-size: 13px;
}

.result-card.error {
  border-color: #f56c6c;
  background: #fef0f0;
}

.drop-zone {
  display: flex;
  align-items: center;
  gap: 10px;
  flex: 1;
  border: 2px dashed #dcdfe6;
  border-radius: 6px;
  padding: 6px 10px;
  transition: border-color 0.2s, background 0.2s;
  cursor: pointer;
}

.drop-zone.dragging {
  border-color: #409eff;
  background: #ecf5ff;
}

.section-title {
  font-size: 13px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 8px;
  padding-bottom: 4px;
  border-bottom: 1px solid #ebeef5;
}
</style>

<template>
  <div class="clip-container">
    <div class="main-layout">
      <!-- 左侧：视频预览 + 抽帧预览 -->
      <div class="left-panel">
        <div style="display: flex; align-items: center; gap: 8px; background: #000; padding: 8px 12px; border-radius: 6px 6px 0 0;">
          <el-button circle size="small" @click="playPrevFile" title="上一个文件" :disabled="!hasPrevFile" style="color: #fff; background: rgba(255,255,255,0.2); border: 1px solid rgba(255,255,255,0.3);">◀</el-button>
          <el-button circle size="small" @click="playNextFile" title="下一个文件" :disabled="!hasNextFile" style="color: #fff; background: rgba(255,255,255,0.2); border: 1px solid rgba(255,255,255,0.3);">▶</el-button>
          <span style="flex: 1; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; color: #fff; font-size: 13px;">{{ filePath || '视频预览' }}</span>
          <el-button circle size="small" @click="selectVideoFile" title="选择文件" style="color: #fff; background: rgba(255,255,255,0.2); border: 1px solid rgba(255,255,255,0.3);">+</el-button>
        </div>
        <div class="video-wrapper" @dblclick="!videoSrc && selectVideoFile()">
          <video
              v-if="videoSrc"
              ref="videoPlayer"
              controls
              preload="metadata"
              @timeupdate="onTimeUpdate"
              @loadedmetadata="onMetadataLoaded"
          >
            <source :src="videoSrc" type="video/mp4">
            您的浏览器不支持 Video 标签。
          </video>
          <div v-else
               style="height: 240px; display: flex; align-items: center; justify-content: center; color: #909399; background: #1a1a1a; border-radius: 6px;">
            <div style="text-align: center;">
              <div style="font-size: 40px; margin-bottom: 8px;">🎬</div>
              <div style="font-size: 14px;">请先选择视频文件</div>
            </div>
          </div>
        </div>
        <div v-if="videoSrc" class="time-display">
          当前时间：{{ formatTime(currentTime) }} / 总时长：{{ formatTime(duration) }}
        </div>

        <!-- 预览提取的图片 -->
        <div v-if="extractedFrames.length > 0" style="margin-top: 16px;">
          <div class="section-title" style="font-size: 13px;">🖼️ 预览（{{ extractedFrames.length }}张）</div>
          <div style="display: grid; grid-template-columns: repeat(auto-fill, minmax(120px, 1fr)); gap: 8px; margin-top: 8px;">
            <div
                v-for="(frame, idx) in extractedFrames"
                :key="idx"
                style="position: relative; border-radius: 6px; overflow: hidden; border: 1px solid #e4e7ed; background: #000;"
            >
              <img :src="frame" style="width: 100%; display: block;" alt="视频帧"/>
              <div style="position: absolute; bottom: 0; left: 0; right: 0; background: rgba(0,0,0,0.6); color: #fff; font-size: 11px; padding: 2px 4px; text-align: center;">
                帧 #{{ idx + 1 }}
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 右侧：选项卡式功能区 -->
      <div class="right-panel" style="overflow-y: auto; padding-right: 4px;">
        <el-tabs v-model="activeTab" type="border-card">
          
          <!-- 选项卡 1：视频裁剪 -->
          <el-tab-pane label="✂️ 裁剪" name="clip">
            <div class="section-title" style="font-size: 13px; margin-bottom: 12px;">裁剪参数</div>

            <!-- 开始时间 -->
            <div class="form-row">
              <span class="label">开始时间：</span>
              <div class="time-input-group" style="flex:1">
                <el-input
                    v-model="startTime"
                    placeholder="00:00:00"
                    style="flex:1"
                    clearable
                />
                <el-button
                    class="time-btn"
                    size="small"
                    @click="setStartFromCurrent"
                    :disabled="!videoSrc"
                    title="使用当前播放位置"
                >📍当前位置
                </el-button>
              </div>
            </div>

            <!-- 结束时间 -->
            <div class="form-row">
              <span class="label">结束时间：</span>
              <div class="time-input-group" style="flex:1">
                <el-input
                    v-model="endTime"
                    placeholder="00:00:00"
                    style="flex:1"
                    clearable
                />
                <el-button
                    class="time-btn"
                    size="small"
                    @click="setEndFromCurrent"
                    :disabled="!videoSrc"
                    title="使用当前播放位置"
                >📍当前位置
                </el-button>
              </div>
            </div>

            <!-- 时间范围预览 -->
            <div v-if="startTime || endTime"
                 style="font-size: 12px; color: #909399; margin-bottom: 12px; padding: 6px 8px; background: #f5f7fa; border-radius: 4px;">
              裁剪片段：{{ startTime || '起始' }} → {{ endTime || '结尾' }}
              <span v-if="clipDurationText" style="margin-left: 6px; color: #409eff;">（约 {{ clipDurationText }}）</span>
            </div>

            <!-- 预览裁剪区段 -->
            <div class="form-row">
              <el-button
                  type="info"
                  size="small"
                  :disabled="!videoSrc || !startTime"
                  @click="previewClip"
              >▶ 预览片段
              </el-button>
              <el-button
                  size="small"
                  @click="clearTimes"
              >清空时间
              </el-button>
            </div>

            <el-divider style="margin: 10px 0"/>

            <!-- 执行裁剪 -->
            <div class="form-row">
              <el-button
                  type="primary"
                  :loading="clipping"
                  :disabled="!filePath || !startTime"
                  @click="clipVideo"
                  style="width: 100%"
              >
                开始裁剪
              </el-button>
            </div>

            <!-- 裁剪结果 -->
            <div v-if="result" class="result-card" :class="{ error: !result.success }" style="margin-top: 12px;">
              <div v-if="result.success">
                <strong style="color: #67c23a;">✅ 裁剪成功，已替换原视频</strong>
                <div style="margin-top: 6px; word-break: break-all;">
                  <div>耗时：{{ result.cost }}</div>
                </div>
                <div style="margin-top: 8px;">
                  <el-button size="small" type="success" @click="openOutputDir">打开目录</el-button>
                  <el-button size="small" @click="playOutput">播放视频</el-button>
                </div>
              </div>
              <div v-else>
                <strong style="color: #f56c6c;">❌ 裁剪失败</strong>
                <div style="margin-top: 6px; color: #606266;">{{ result.message }}</div>
              </div>
            </div>
          </el-tab-pane>

          <!-- 选项卡 2：视频旋转 -->
          <el-tab-pane label="🔄 旋转" name="rotate">
            <div class="section-title" style="font-size: 13px; margin-bottom: 12px;">旋转设置</div>
            
            <div class="form-row">
              <span class="label">旋转角度：</span>
              <el-radio-group v-model="rotateAngle" size="small">
                <el-radio-button label="90">90°</el-radio-button>
                <el-radio-button label="180">180°</el-radio-button>
                <el-radio-button label="270">270°</el-radio-button>
              </el-radio-group>
            </div>
            
            <div class="form-row">
              <span class="label">旋转方向：</span>
              <el-radio-group v-model="rotateDirection" size="small">
                <el-radio-button label="clockwise">顺时针</el-radio-button>
                <el-radio-button label="anticlockwise">逆时针</el-radio-button>
              </el-radio-group>
            </div>
            
            <div class="form-row">
              <el-button
                  type="warning"
                  :loading="rotating"
                  :disabled="!filePath"
                  @click="rotateVideo"
                  style="width: 100%"
              >
                开始旋转
              </el-button>
            </div>

            <!-- 旋转结果 -->
            <div v-if="rotateResult" class="result-card" :class="{ error: !rotateResult.success }" style="margin-top: 12px;">
              <div v-if="rotateResult.success">
                <strong style="color: #67c23a;">✅ 旋转成功</strong>
                <div style="margin-top: 6px; word-break: break-all;">
                  <div>输出：{{ rotateResult.outputPath }}</div>
                  <div>耗时：{{ rotateResult.cost }}</div>
                </div>
                <div style="margin-top: 8px;">
                  <el-button size="small" type="success" @click="openRotateDir">打开目录</el-button>
                  <el-button size="small" @click="playRotate">播放结果</el-button>
                </div>
              </div>
              <div v-else>
                <strong style="color: #f56c6c;">❌ 旋转失败</strong>
                <div style="margin-top: 6px; color: #606266;">{{ rotateResult.message }}</div>
              </div>
            </div>
          </el-tab-pane>

          <!-- 选项卡 3：视频抽帧 -->
          <el-tab-pane label="📸 抽帧" name="frames">
            <div class="section-title" style="font-size: 13px; margin-bottom: 12px;">抽帧设置</div>
            
            <div class="form-row">
              <span class="label">抽取数量：</span>
              <el-input-number
                  v-model="frameCount"
                  :min="1"
                  :max="100"
                  :step="1"
                  size="small"
                  style="width: 120px;"
              />
              <span style="font-size: 12px; color: #909399; margin-left: 8px;">（1-100 张）</span>
            </div>
            
            <div class="form-row">
              <el-button
                  type="success"
                  :loading="extractingFrames"
                  :disabled="!filePath"
                  @click="extractFrames"
                  style="width: 100%"
              >
                {{ extractingFrames ? '抽取中...' : '开始抽帧' }}
              </el-button>
            </div>

            <!-- 抽帧结果 -->
            <div v-if="framesResult" class="result-card" :class="{ error: !framesResult.success }" style="margin-top: 12px;">
              <div v-if="framesResult.success">
                <strong style="color: #67c23a;">✅ 抽帧成功</strong>
                <div style="margin-top: 6px;">
                  <div>共提取 {{ framesResult.frameCount }} 张图片，耗时：{{ framesResult.cost }}</div>
                  <div style="font-size: 12px; color: #909399; margin-top: 4px;">输出目录：{{ framesResult.outputDir }}</div>
                </div>
                <div style="margin-top: 8px;">
                  <el-button size="small" type="success" @click="openFramesDir">打开目录</el-button>
                </div>
              </div>
              <div v-else>
                <strong style="color: #f56c6c;">❌ 抽帧失败</strong>
                <div style="margin-top: 6px; color: #606266;">{{ framesResult.message }}</div>
              </div>
            </div>
          </el-tab-pane>

          <!-- 选项卡 4：声音分离 -->
          <el-tab-pane label="🎵 声音分离" name="audio">
            <div class="section-title" style="font-size: 13px; margin-bottom: 12px;">音频提取</div>
            
            <div class="form-row">
              <span class="label">输出格式：</span>
              <el-select v-model="audioFormat" style="width: 120px;">
                <el-option label="MP3" value="mp3"/>
                <el-option label="AAC" value="aac"/>
                <el-option label="WAV" value="wav"/>
                <el-option label="FLAC" value="flac"/>
              </el-select>
            </div>
            
            <div class="form-row">
              <el-button
                  type="warning"
                  :loading="extracting"
                  :disabled="!filePath"
                  @click="extractAudio"
                  style="width: 100%"
              >
                提取音频
              </el-button>
            </div>

            <!-- 音频提取结果 -->
            <div v-if="audioResult" class="result-card" :class="{ error: !audioResult.success }" style="margin-top: 12px;">
              <div v-if="audioResult.success">
                <strong style="color: #67c23a;">✅ 提取成功</strong>
                <div style="margin-top: 6px; word-break: break-all;">
                  <div>输出：{{ audioResult.outputPath }}</div>
                  <div>耗时：{{ audioResult.cost }}</div>
                </div>
                <div style="margin-top: 8px;">
                  <el-button size="small" type="success" @click="openAudioDir">打开目录</el-button>
                </div>
              </div>
              <div v-else>
                <strong style="color: #f56c6c;">❌ 提取失败</strong>
                <div style="margin-top: 6px; color: #606266;">{{ audioResult.message }}</div>
              </div>
            </div>
          </el-tab-pane>

        </el-tabs>

        <!-- 执行日志（固定在底部） -->
        <div style="margin-top: 16px;">
          <div class="section-title" style="font-size: 13px;">📝 执行日志</div>
          <div class="log-area" ref="logArea">
            <div v-if="logs.length === 0" style="color: #c0c4cc; text-align: center; padding-top: 40px;">暂无日志</div>
            <div v-for="(log, i) in logs" :key="i">{{ log }}</div>
          </div>
          <div style="margin-top: 6px; text-align: right;">
            <el-button size="small" text @click="clearLogs" :disabled="logs.length === 0">清空日志</el-button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import {ClipAndReplaceVideo, SelectVideoFile, OpenExplorer, ExtractAudio, RotateVideo, ExtractFrames, GetVideoFilesInDir} from "../wailsjs/go/app/App.js";
import {OnFileDrop, OnFileDropOff} from "../wailsjs/runtime/runtime.js";

// 本地文件服务端口
const LOCAL_FILE_PORT = 39870;

// 支持的视频扩展名
const VIDEO_EXTS = ['.mp4', '.avi', '.mkv', '.mov', '.wmv', '.flv', '.m4v', '.ts', '.rmvb'];

export default {
  name: 'videoClip',
  data() {
    return {
      filePath: '',
      dirFiles: [],
      currentIndex: -1,
      videoSrc: '',
      startTime: '',
      endTime: '',
      currentTime: 0,
      duration: 0,
      clipping: false,
      result: null,
      logs: [],
      isDragging: false,
      // 选项卡
      activeTab: 'clip',
      // 声音分离
      audioFormat: 'mp3',
      extracting: false,
      audioResult: null,
      // 视频旋转
      rotateAngle: '90',
      rotateDirection: 'clockwise',
      rotating: false,
      rotateResult: null,
      // 视频抽帧
      frameCount: 10,
      extractingFrames: false,
      framesResult: null,
      extractedFrames: [],
    };
  },
  computed: {
    clipDurationText() {
      if (!this.startTime || !this.endTime) return '';
      const s = this.parseTimeToSeconds(this.startTime);
      const e = this.parseTimeToSeconds(this.endTime);
      if (isNaN(s) || isNaN(e) || e <= s) return '';
      return this.formatTime(e - s);
    },
    hasPrevFile() {
      return this.currentIndex > 0;
    },
    hasNextFile() {
      return this.currentIndex >= 0 && this.currentIndex < this.dirFiles.length - 1;
    },
  },
  mounted() {
    window.runtime.EventsOn('back_msg', (message) => {
      const msg = JSON.parse(message);
      this.logs.push(msg.time + ' ' + msg.msg);
      this.$nextTick(() => {
        if (this.$refs.logArea) {
          this.$refs.logArea.scrollTop = this.$refs.logArea.scrollHeight;
        }
      });
    });

    // 监听文件拖放，获取真实本地路径
    // useDropTarget=false 表示全局监听，不限于 data-wails-drop-target 元素
    OnFileDrop((x, y, paths) => {
      if (!paths || paths.length === 0) return;
      const fp = paths[0];
      const ext = fp.substring(fp.lastIndexOf('.')).toLowerCase();
      if (!VIDEO_EXTS.includes(ext)) {
        this.$message.warning('请拖入视频文件（支持 mp4/avi/mkv/mov 等格式）');
        return;
      }
      this.filePath = fp;
      this.loadVideo(fp);
      this.result = null;
      this.isDragging = false;
    }, false);
  },
  beforeUnmount() {
    window.runtime.EventsOff('back_msg');
    OnFileDropOff();
  },
  methods: {
    async selectVideoFile() {
      try {
        const fp = await SelectVideoFile();
        if (fp) {
          this.filePath = fp;
          this.loadVideo(fp);
          this.result = null;
          // 加载同目录视频文件
          const files = await GetVideoFilesInDir(fp);
          this.dirFiles = files;
          this.currentIndex = files.indexOf(fp);
        }
      } catch (e) {
        console.error('选择文件失败:', e);
      }
    },

    playPrevFile() {
      if (this.currentIndex > 0) {
        this.currentIndex--;
        this.filePath = this.dirFiles[this.currentIndex];
        this.loadVideo(this.filePath);
      }
    },

    playNextFile() {
      if (this.currentIndex >= 0 && this.currentIndex < this.dirFiles.length - 1) {
        this.currentIndex++;
        this.filePath = this.dirFiles[this.currentIndex];
        this.loadVideo(this.filePath);
      }
    },

    onFilePathChange(val) {
      if (val) {
        this.loadVideo(val);
      } else {
        this.videoSrc = '';
      }
    },

    loadVideo(fp) {
      // 通过本地文件 HTTP 服务加载视频，支持 Range 请求（视频拖拽/快进必须）
      const encoded = encodeURIComponent(fp);
      this.videoSrc = `http://localhost:${LOCAL_FILE_PORT}/file?path=${encoded}`;
      this.currentTime = 0;
      this.duration = 0;
      // 重新加载 video 元素
      this.$nextTick(() => {
        if (this.$refs.videoPlayer) {
          this.$refs.videoPlayer.load();
        }
      });
    },

    onTimeUpdate(e) {
      this.currentTime = e.target.currentTime;
    },

    onMetadataLoaded(e) {
      this.duration = e.target.duration;
    },

    setStartFromCurrent() {
      this.startTime = this.formatTime(this.currentTime);
    },

    setEndFromCurrent() {
      this.endTime = this.formatTime(this.currentTime);
    },

    previewClip() {
      if (!this.$refs.videoPlayer || !this.startTime) return;
      const t = this.parseTimeToSeconds(this.startTime);
      if (!isNaN(t)) {
        this.$refs.videoPlayer.currentTime = t;
        this.$refs.videoPlayer.play();
        // 如果设置了结束时间，到达结束时间时暂停
        if (this.endTime) {
          const endSec = this.parseTimeToSeconds(this.endTime);
          const checkEnd = () => {
            if (this.$refs.videoPlayer && this.$refs.videoPlayer.currentTime >= endSec) {
              this.$refs.videoPlayer.pause();
              this.$refs.videoPlayer.removeEventListener('timeupdate', checkEnd);
            }
          };
          this.$refs.videoPlayer.addEventListener('timeupdate', checkEnd);
        }
      }
    },

    clearTimes() {
      this.startTime = '';
      this.endTime = '';
    },

    async clipVideo() {
      if (!this.filePath) {
        this.$message.warning('请先选择视频文件');
        return;
      }
      if (!this.startTime) {
        this.$message.warning('请输入开始时间');
        return;
      }

      try {
        await this.$confirm('裁剪后将替换原视频文件，是否继续？', '确认裁剪', {
          confirmButtonText: '确认',
          cancelButtonText: '取消',
          type: 'warning',
        });
      } catch {
        return;
      }

      this.clipping = true;
      this.result = null;
      this.logs = [];

      try {
        const res = await ClipAndReplaceVideo({
          filePath: this.filePath,
          startTime: this.startTime,
          endTime: this.endTime,
        });
        this.result = res;
        if (res.success) {
          this.$message.success('视频裁剪完成，已替换原文件！');
          // 重新加载裁剪后的视频
          this.loadVideo(this.filePath);
          this.startTime = '';
          this.endTime = '';
        } else {
          this.$message.error('裁剪失败：' + res.message);
        }
      } catch (e) {
        this.$message.error('裁剪失败: ' + e);
      } finally {
        this.clipping = false;
      }
    },

    openOutputDir() {
      if (this.filePath) {
        const path = this.filePath;
        const lastSep = Math.max(path.lastIndexOf('/'), path.lastIndexOf('\\'));
        const dir = lastSep !== -1 ? path.substring(0, lastSep) : path;
        OpenExplorer(dir);
      }
    },

    playOutput() {
      if (this.filePath) {
        this.loadVideo(this.filePath);
        this.startTime = '';
        this.endTime = '';
        this.result = null;
      }
    },

    async rotateVideo() {
      if (!this.filePath) {
        this.$message.warning('请先选择视频文件');
        return;
      }
      this.rotating = true;
      this.rotateResult = null;
      this.logs = [];

      try {
        const res = await RotateVideo({
          filePath: this.filePath,
          angle: parseInt(this.rotateAngle),
          clockwise: this.rotateDirection === 'clockwise',
        });
        this.rotateResult = res;
        if (res.success) {
          this.$message.success('视频旋转完成！');
        } else {
          this.$message.error('旋转失败：' + res.message);
        }
      } catch (e) {
        this.$message.error('旋转失败：' + e);
      } finally {
        this.rotating = false;
      }
    },

    openRotateDir() {
      if (this.rotateResult && this.rotateResult.outputPath) {
        const path = this.rotateResult.outputPath;
        const lastSep = Math.max(path.lastIndexOf('/'), path.lastIndexOf('\\'));
        const dir = lastSep !== -1 ? path.substring(0, lastSep) : path;
        OpenExplorer(dir);
      }
    },

    playRotate() {
      if (this.rotateResult && this.rotateResult.outputPath) {
        this.filePath = this.rotateResult.outputPath;
        this.loadVideo(this.rotateResult.outputPath);
        this.rotateResult = null;
      }
    },

    async extractFrames() {
      if (!this.filePath) {
        this.$message.warning('请先选择视频文件');
        return;
      }
      if (this.frameCount < 1 || this.frameCount > 100) {
        this.$message.warning('抽取数量必须在 1-100 之间');
        return;
      }

      this.extractingFrames = true;
      this.framesResult = null;
      this.extractedFrames = [];
      this.logs = [];

      try {
        const res = await ExtractFrames({
          filePath: this.filePath,
          count: this.frameCount,
        });
        this.framesResult = res;
        
        if (res.success) {
          this.$message.success(`成功提取 ${res.frameCount} 张图片！`);
          // 加载提取的图片用于预览
          this.extractedFrames = res.framePaths.map(path => 
            `http://localhost:${LOCAL_FILE_PORT}/file?path=${encodeURIComponent(path)}`
          );
        } else {
          this.$message.error('抽帧失败：' + res.message);
        }
      } catch (e) {
        this.$message.error('抽帧失败：' + e);
      } finally {
        this.extractingFrames = false;
      }
    },

    openFramesDir() {
      if (this.framesResult && this.framesResult.outputDir) {
        OpenExplorer(this.framesResult.outputDir);
      }
    },

    clearLogs() {
      this.logs = [];
    },

    async extractAudio() {
      if (!this.filePath) {
        this.$message.warning('请先选择视频文件');
        return;
      }
      this.extracting = true;
      this.audioResult = null;
      this.logs = [];

      try {
        const res = await ExtractAudio({
          filePath: this.filePath,
          format: this.audioFormat,
        });
        this.audioResult = res;
        if (res.success) {
          this.$message.success('音频提取完成！');
        } else {
          this.$message.error('提取失败：' + res.message);
        }
      } catch (e) {
        this.$message.error('提取失败: ' + e);
      } finally {
        this.extracting = false;
      }
    },

    openAudioDir() {
      if (this.audioResult && this.audioResult.outputPath) {
        const path = this.audioResult.outputPath;
        const lastSep = Math.max(path.lastIndexOf('/'), path.lastIndexOf('\\'));
        const dir = lastSep !== -1 ? path.substring(0, lastSep) : path;
        OpenExplorer(dir);
      }
    },

    // 将秒数格式化为 HH:MM:SS.ms
    formatTime(seconds) {
      if (isNaN(seconds) || seconds < 0) return '00:00:00';
      const h = Math.floor(seconds / 3600);
      const m = Math.floor((seconds % 3600) / 60);
      const s = Math.floor(seconds % 60);
      const ms = Math.round((seconds - Math.floor(seconds)) * 1000);
      const base = [
        String(h).padStart(2, '0'),
        String(m).padStart(2, '0'),
        String(s).padStart(2, '0'),
      ].join(':');
      return ms > 0 ? base + '.' + String(ms).padStart(3, '0') : base;
    },

    // 将 HH:MM:SS 或 HH:MM:SS.ms 转换为秒数
    parseTimeToSeconds(timeStr) {
      if (!timeStr) return NaN;
      const parts = timeStr.split(':');
      if (parts.length !== 3) return NaN;
      const [h, m, s] = parts.map(Number);
      if (isNaN(h) || isNaN(m) || isNaN(s)) return NaN;
      return h * 3600 + m * 60 + s;
    },
  },
};
</script>
