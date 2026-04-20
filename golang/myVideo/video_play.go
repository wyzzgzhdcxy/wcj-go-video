package myVideo

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
	"github.com/wyzzgzhdcxy/wcj-go-common/core"
)

type CategorizeFilesReq struct {
	Dir        string `json:"dir"`
	Categories string `json:"categories"`
	Count      int    `json:"count"`
}

type VideoReq struct {
	Path string `json:"path"`
	Type string `json:"type"`
}

func getRealPath(path string) string {
	driveList := [3]string{"D:/", "E:/", "F:/"}
	for _, v := range driveList {
		if core.FileExist(v + path) {
			return v + path
		}
	}
	return ""
}

var preDeleteFile = ""
var preDeleteOneMinute = ""
var preAnticlockwise = ""
var preClockwise = ""

func anticlockwiseVideo(filepath string) {
	startTime := time.Now()
	cmdStr := "ffmpeg -i " + filepath + " -vf \"transpose=2\" D:\\11.mp4"
	log.Printf("逆时针旋转视频：" + cmdStr)
	core.ExecuteCmderLine(cmdStr)
	endTime := time.Now()
	duration := endTime.Sub(startTime)
	log.Printf("逆时针旋转视频,执行时间: %s\n", duration)
}

func clockwiseVideo(filepath string) {
	startTime := time.Now()
	cmdStr := "ffmpeg -i " + filepath + " -vf \"transpose=1\" " + filepath + ".clock.mp4"
	log.Printf("顺时针旋转视频：" + cmdStr)
	core.ExecuteCmderLine(cmdStr)
	endTime := time.Now()
	duration := endTime.Sub(startTime)
	log.Printf("顺时针旋转视频,执行时间: %s\n", duration)
}

// RotateVideoReq 视频旋转请求
type RotateVideoReq struct {
	FilePath  string `json:"filePath"`
	Angle     int    `json:"angle"`     // 旋转角度: 90, 180, 270
	Clockwise bool   `json:"clockwise"` // true: 顺时针, false: 逆时针
}

// RotateVideoRes 视频旋转结果
type RotateVideoRes struct {
	Success    bool   `json:"success"`
	OutputPath string `json:"outputPath"`
	Message    string `json:"message"`
	Cost       string `json:"cost"`
}

// RotateVideo 旋转视频文件
// angle: 旋转角度（90, 180, 270）
// clockwise: true=顺时针, false=逆时针
func RotateVideo(req RotateVideoReq, callback func([]byte)) RotateVideoRes {
	startTime := time.Now()
	callback([]byte("开始旋转视频: " + req.FilePath))

	// 根据角度和方向计算 ffmpeg transpose 参数
	// transpose=0: 逆时针90度
	// transpose=1: 顺时针90度
	// transpose=2: 逆时针90度并翻转
	// transpose=3: 顺时针90度并翻转
	var vfFilter string
	outputPath := req.FilePath

	// 去掉扩展名并添加后缀
	ext := ".mp4"
	basePath := req.FilePath
	if idx := strings.LastIndex(req.FilePath, "."); idx != -1 {
		ext = req.FilePath[idx:]
		basePath = req.FilePath[:idx]
	}

	direction := "clockwise"
	if !req.Clockwise {
		direction = "anticlockwise"
	}

	switch req.Angle {
	case 90:
		if req.Clockwise {
			vfFilter = "transpose=1"
		} else {
			vfFilter = "transpose=2"
		}
		outputPath = basePath + "_rotate90_" + direction + ext
	case 180:
		vfFilter = "transpose=1,transpose=1"
		outputPath = basePath + "_rotate180" + ext
	case 270:
		if req.Clockwise {
			vfFilter = "transpose=2"
		} else {
			vfFilter = "transpose=1"
		}
		outputPath = basePath + "_rotate270_" + direction + ext
	default:
		return RotateVideoRes{
			Success: false,
			Message: "不支持的旋转角度，请选择90、180或270度",
		}
	}

	cmdStr := fmt.Sprintf("ffmpeg -i \"%s\" -vf \"%s\" -c:a copy \"%s\"", req.FilePath, vfFilter, outputPath)
	log.Printf("执行视频旋转命令: %s", cmdStr)
	callback([]byte("执行命令: " + cmdStr))

	// 直接传参给 exec.Command，避免 cmd /c 嵌套引号解析失败
	execCmd := exec.Command("ffmpeg", "-i", req.FilePath, "-vf", vfFilter, "-c:a", "copy", outputPath)
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	if err := execCmd.Run(); err != nil {
		log.Printf("ffmpeg 执行失败: %v", err)
		return RotateVideoRes{
			Success: false,
			Message: "ffmpeg 执行失败: " + err.Error(),
		}
	}

	elapsed := time.Since(startTime)
	callback([]byte("旋转完成，耗时: " + elapsed.String()))

	return RotateVideoRes{
		Success:    true,
		OutputPath: outputPath,
		Message:    "视频旋转完成",
		Cost:       elapsed.String(),
	}
}

// ClipVideoReq 视频裁剪请求
type ClipVideoReq struct {
	FilePath  string `json:"filePath"`
	StartTime string `json:"startTime"` // 格式: HH:MM:SS 或 HH:MM:SS.ms
	EndTime   string `json:"endTime"`   // 格式: HH:MM:SS 或 HH:MM:SS.ms
}

// ClipVideoRes 视频裁剪结果
type ClipVideoRes struct {
	Success    bool   `json:"success"`
	OutputPath string `json:"outputPath"`
	Message    string `json:"message"`
	Cost       string `json:"cost"`
}

// ClipVideo 裁剪视频文件（指定开始和结束时间点）
func ClipVideo(req ClipVideoReq, callback func([]byte)) ClipVideoRes {
	startTime := time.Now()
	callback([]byte("开始裁剪视频：" + req.FilePath))

	if req.StartTime == "" {
		return ClipVideoRes{Success: false, Message: "请输入开始时间"}
	}

	// 如果没有输入结束时间，ffmpeg 会自动裁剪到视频末尾
	endTime := req.EndTime
	if endTime == "" {
		callback([]byte("未输入结束时间，将裁剪到视频末尾"))
	}

	// 生成输出路径
	ext := ".mp4"
	basePath := req.FilePath
	if idx := strings.LastIndex(req.FilePath, "."); idx != -1 {
		ext = req.FilePath[idx:]
		basePath = req.FilePath[:idx]
	}
	safeStart := strings.ReplaceAll(req.StartTime, ":", "-")
	var outputPath string
	if endTime == "" {
		// 只裁剪开始时间，输出文件名不带结束时间
		outputPath = fmt.Sprintf("%s_clip_%s%s", basePath, safeStart, ext)
		callback([]byte(fmt.Sprintf("裁剪范围：%s ~ 视频末尾", req.StartTime)))
	} else {
		safeEnd := strings.ReplaceAll(endTime, ":", "-")
		outputPath = fmt.Sprintf("%s_clip_%s_%s%s", basePath, safeStart, safeEnd, ext)
		callback([]byte(fmt.Sprintf("裁剪范围：%s ~ %s", req.StartTime, endTime)))
	}

	// 将 -ss 放在 -i 之前（input seeking），ffmpeg 直接从最近关键帧定位，速度快
	// -avoid_negative_ts 1：修正因从关键帧开始导致的时间戳偏移，避免开头卡帧
	// -c copy：无损复制流，不重新编码
	var execCmd *exec.Cmd
	if endTime == "" {
		// 只指定开始时间，裁剪到视频末尾
		execCmd = exec.Command("ffmpeg",
			"-ss", req.StartTime,
			"-i", req.FilePath,
			"-c", "copy",
			"-avoid_negative_ts", "1",
			outputPath)
	} else {
		// 指定开始和结束时间
		execCmd = exec.Command("ffmpeg",
			"-ss", req.StartTime,
			"-i", req.FilePath,
			"-to", endTime,
			"-c", "copy",
			"-avoid_negative_ts", "1",
			outputPath)
	}

	// Windows 下隐藏命令行窗口
	if runtime.GOOS == "windows" {
		execCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	}

	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	if err := execCmd.Run(); err != nil {
		log.Printf("ffmpeg 裁剪失败: %v", err)
		return ClipVideoRes{
			Success: false,
			Message: "ffmpeg 裁剪失败: " + err.Error(),
		}
	}

	elapsed := time.Since(startTime)
	callback([]byte("裁剪完成，耗时: " + elapsed.String()))

	return ClipVideoRes{
		Success:    true,
		OutputPath: outputPath,
		Message:    "视频裁剪完成",
		Cost:       elapsed.String(),
	}
}

// ClipAndReplaceReq 裁剪并替换原视频请求
type ClipAndReplaceReq struct {
	FilePath  string `json:"filePath"`
	StartTime string `json:"startTime"` // 格式: HH:MM:SS 或 HH:MM:SS.ms
	EndTime   string `json:"endTime"`   // 格式: HH:MM:SS 或 HH:MM:SS.ms，为空则裁剪到末尾
}

// ClipAndReplaceRes 裁剪并替换结果
type ClipAndReplaceRes struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Cost    string `json:"cost"`
}

// ClipAndReplaceVideo 裁剪视频并替换原文件
func ClipAndReplaceVideo(req ClipAndReplaceReq, callback func([]byte)) ClipAndReplaceRes {
	startTime := time.Now()
	callback([]byte("开始裁剪视频：" + req.FilePath))

	if req.FilePath == "" {
		return ClipAndReplaceRes{Success: false, Message: "请选择视频文件"}
	}

	if req.StartTime == "" {
		return ClipAndReplaceRes{Success: false, Message: "请输入开始时间"}
	}

	// 如果没有输入结束时间，ffmpeg 会自动裁剪到视频末尾
	endTime := req.EndTime
	if endTime == "" {
		callback([]byte("未输入结束时间，将裁剪到视频末尾"))
	}

	// 生成临时输出路径
	ext := ".mp4"
	basePath := req.FilePath
	if idx := strings.LastIndex(req.FilePath, "."); idx != -1 {
		ext = req.FilePath[idx:]
		basePath = req.FilePath[:idx]
	}
	tempPath := fmt.Sprintf("%s_clip_temp_%d%s", basePath, time.Now().UnixNano(), ext)

	// 将 -ss 放在 -i 之前（input seeking），ffmpeg 直接从最近关键帧定位，速度快
	// -avoid_negative_ts 1：修正因从关键帧开始导致的时间戳偏移，避免开头卡帧
	// -c copy：无损复制流，不重新编码
	var execCmd *exec.Cmd
	if endTime == "" {
		execCmd = exec.Command("ffmpeg",
			"-ss", req.StartTime,
			"-i", req.FilePath,
			"-c", "copy",
			"-avoid_negative_ts", "1",
			tempPath)
	} else {
		execCmd = exec.Command("ffmpeg",
			"-ss", req.StartTime,
			"-i", req.FilePath,
			"-to", endTime,
			"-c", "copy",
			"-avoid_negative_ts", "1",
			tempPath)
	}

	// Windows 下隐藏命令行窗口
	if runtime.GOOS == "windows" {
		execCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	}

	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	if err := execCmd.Run(); err != nil {
		log.Printf("ffmpeg 裁剪失败: %v", err)
		return ClipAndReplaceRes{
			Success: false,
			Message: "ffmpeg 裁剪失败: " + err.Error(),
		}
	}

	// 删除原文件
	if err := os.Remove(req.FilePath); err != nil {
		log.Printf("删除原文件失败: %v", err)
		os.Remove(tempPath) // 清理临时文件
		return ClipAndReplaceRes{
			Success: false,
			Message: "裁剪成功但删除原文件失败: " + err.Error(),
		}
	}

	// 将临时文件重命名为原文件路径
	if err := os.Rename(tempPath, req.FilePath); err != nil {
		log.Printf("重命名文件失败: %v", err)
		return ClipAndReplaceRes{
			Success: false,
			Message: "裁剪成功但重命名失败: " + err.Error() + "，原文件已删除",
		}
	}

	elapsed := time.Since(startTime)
	callback([]byte("视频裁剪完成，已替换原文件，耗时: " + elapsed.String()))

	return ClipAndReplaceRes{
		Success: true,
		Message: "视频裁剪完成，已替换原文件",
		Cost:    elapsed.String(),
	}
}

// RemoveVideoIntroReq 去除片头请求
type RemoveVideoIntroReq struct {
	FilePath      string  `json:"filePath"`
	IntroDuration float64 `json:"introDuration"` // 片头时长（秒）
}

// RemoveVideoIntroRes 去除片头结果
type RemoveVideoIntroRes struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Cost    string `json:"cost"`
}

// RemoveVideoIntro 去除视频片头（从指定时间裁剪到末尾，替换原文件）
func RemoveVideoIntro(req RemoveVideoIntroReq, callback func([]byte)) RemoveVideoIntroRes {
	startTime := time.Now()
	callback([]byte("开始去除片头：" + req.FilePath))

	if req.FilePath == "" {
		return RemoveVideoIntroRes{Success: false, Message: "请选择视频文件"}
	}

	if req.IntroDuration <= 0 {
		return RemoveVideoIntroRes{Success: false, Message: "片头时长必须大于0"}
	}

	// 生成临时输出路径
	ext := ".mp4"
	basePath := req.FilePath
	if idx := strings.LastIndex(req.FilePath, "."); idx != -1 {
		ext = req.FilePath[idx:]
		basePath = req.FilePath[:idx]
	}
	tempPath := fmt.Sprintf("%s_intro_temp_%d%s", basePath, time.Now().UnixNano(), ext)

	// 将 -ss 放在 -i 之前（input seeking），ffmpeg 直接从最近关键帧定位，速度快
	// -c copy：无损复制流，不重新编码
	execCmd := exec.Command("ffmpeg",
		"-ss", fmt.Sprintf("%.3f", req.IntroDuration),
		"-i", req.FilePath,
		"-c", "copy",
		"-avoid_negative_ts", "1",
		tempPath)

	// Windows 下隐藏命令行窗口
	if runtime.GOOS == "windows" {
		execCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	}

	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	if err := execCmd.Run(); err != nil {
		log.Printf("ffmpeg 去除片头失败: %v", err)
		return RemoveVideoIntroRes{
			Success: false,
			Message: "ffmpeg 去除片头失败: " + err.Error(),
		}
	}

	// 删除原文件
	if err := os.Remove(req.FilePath); err != nil {
		log.Printf("删除原文件失败: %v", err)
		os.Remove(tempPath) // 清理临时文件
		return RemoveVideoIntroRes{
			Success: false,
			Message: "去除片头成功但删除原文件失败: " + err.Error(),
		}
	}

	// 将临时文件重命名为原文件路径
	if err := os.Rename(tempPath, req.FilePath); err != nil {
		log.Printf("重命名文件失败: %v", err)
		return RemoveVideoIntroRes{
			Success: false,
			Message: "去除片头成功但重命名失败: " + err.Error() + "，原文件已删除",
		}
	}

	elapsed := time.Since(startTime)
	callback([]byte(fmt.Sprintf("去除片头完成，已替换原文件，耗时: %s，去除片头时长: %.1f秒", elapsed.String(), req.IntroDuration)))

	return RemoveVideoIntroRes{
		Success: true,
		Message: fmt.Sprintf("去除片头完成，已去除前 %.1f 秒", req.IntroDuration),
		Cost:    elapsed.String(),
	}
}

// BatchRemoveVideoIntroReq 批量去除片头请求
type BatchRemoveVideoIntroReq struct {
	FilePaths     []string `json:"filePaths"`
	IntroDuration float64  `json:"introDuration"` // 片头时长（秒）
}

// BatchRemoveVideoIntroRes 批量去除片头结果
type BatchRemoveVideoIntroRes struct {
	Success     bool     `json:"success"`
	Message     string   `json:"message"`
	Processed   int      `json:"processed"`   // 已处理数量
	Failed      int      `json:"failed"`      // 失败数量
	FailedFiles []string `json:"failedFiles"` // 失败文件列表
	TotalCost   string   `json:"totalCost"`   // 总耗时
}

// BatchRemoveVideoIntro 批量去除视频片头
func BatchRemoveVideoIntro(req BatchRemoveVideoIntroReq, callback func([]byte)) BatchRemoveVideoIntroRes {
	startTime := time.Now()
	callback([]byte(fmt.Sprintf("开始批量去除片头，共 %d 个视频，片头时长: %.1f 秒", len(req.FilePaths), req.IntroDuration)))

	if len(req.FilePaths) == 0 {
		return BatchRemoveVideoIntroRes{Success: false, Message: "请选择要处理的视频"}
	}

	if req.IntroDuration <= 0 {
		return BatchRemoveVideoIntroRes{Success: false, Message: "片头时长必须大于0"}
	}

	processed := 0
	failed := 0
	var failedFiles []string

	for i, filePath := range req.FilePaths {
		callback([]byte(fmt.Sprintf("正在处理 [%d/%d]: %s", i+1, len(req.FilePaths), filePath)))

		// 生成临时输出路径
		ext := ".mp4"
		basePath := filePath
		if idx := strings.LastIndex(filePath, "."); idx != -1 {
			ext = filePath[idx:]
			basePath = filePath[:idx]
		}
		tempPath := fmt.Sprintf("%s_intro_temp_%d%s", basePath, time.Now().UnixNano(), ext)

		execCmd := exec.Command("ffmpeg",
			"-ss", fmt.Sprintf("%.3f", req.IntroDuration),
			"-i", filePath,
			"-c", "copy",
			"-avoid_negative_ts", "1",
			tempPath)

		if runtime.GOOS == "windows" {
			execCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		}

		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr

		if err := execCmd.Run(); err != nil {
			log.Printf("去除片头失败 [%s]: %v", filePath, err)
			failed++
			failedFiles = append(failedFiles, filePath+" ("+err.Error()+")")
			continue
		}

		// 删除原文件
		if err := os.Remove(filePath); err != nil {
			log.Printf("删除原文件失败 [%s]: %v", filePath, err)
			os.Remove(tempPath)
			failed++
			failedFiles = append(failedFiles, filePath+" (删除原文件失败)")
			continue
		}

		// 重命名为原文件
		if err := os.Rename(tempPath, filePath); err != nil {
			log.Printf("重命名文件失败 [%s]: %v", filePath, err)
			failed++
			failedFiles = append(failedFiles, filePath+" (重命名失败)")
			continue
		}

		processed++
		callback([]byte(fmt.Sprintf("✅ 处理完成 [%d/%d]: %s", i+1, len(req.FilePaths), filePath)))
	}

	elapsed := time.Since(startTime)
	callback([]byte(fmt.Sprintf("批量处理完成，成功: %d，失败: %d，总耗时: %s", processed, failed, elapsed.String())))

	return BatchRemoveVideoIntroRes{
		Success:     failed == 0,
		Processed:   processed,
		Failed:      failed,
		FailedFiles: failedFiles,
		TotalCost:   elapsed.String(),
	}
}

// ExtractAudioReq 声音分离请求
type ExtractAudioReq struct {
	FilePath string `json:"filePath"`
	Format   string `json:"format"` // 输出格式: mp3, aac, wav, flac, 默认 mp3
}

// ExtractAudioRes 声音分离结果
type ExtractAudioRes struct {
	Success    bool   `json:"success"`
	OutputPath string `json:"outputPath"`
	Message    string `json:"message"`
	Cost       string `json:"cost"`
}

// ExtractAudio 从视频中分离音频
func ExtractAudio(req ExtractAudioReq, callback func([]byte)) ExtractAudioRes {
	startTime := time.Now()
	callback([]byte("开始提取音频: " + req.FilePath))

	if req.FilePath == "" {
		return ExtractAudioRes{Success: false, Message: "请选择视频文件"}
	}

	format := req.Format
	if format == "" {
		format = "mp3"
	}

	// 生成输出路径
	basePath := req.FilePath
	if idx := strings.LastIndex(req.FilePath, "."); idx != -1 {
		basePath = req.FilePath[:idx]
	}
	outputPath := fmt.Sprintf("%s_audio.%s", basePath, format)

	callback([]byte(fmt.Sprintf("输出格式: %s，输出路径: %s", format, outputPath)))

	// 使用 exec.Command 直接传参，避免路径含中文/空格解析失败
	// -vn: 不包含视频流
	// -acodec: 指定音频编解码器
	var execCmd *exec.Cmd
	switch format {
	case "wav":
		execCmd = exec.Command("ffmpeg", "-i", req.FilePath, "-vn", "-acodec", "pcm_s16le", "-y", outputPath)
	case "aac":
		execCmd = exec.Command("ffmpeg", "-i", req.FilePath, "-vn", "-acodec", "copy", "-y", outputPath)
	case "flac":
		execCmd = exec.Command("ffmpeg", "-i", req.FilePath, "-vn", "-acodec", "flac", "-y", outputPath)
	default: // mp3
		execCmd = exec.Command("ffmpeg", "-i", req.FilePath, "-vn", "-acodec", "libmp3lame", "-q:a", "2", "-y", outputPath)
	}

	// Windows 下隐藏命令行窗口
	if runtime.GOOS == "windows" {
		execCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	}

	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	if err := execCmd.Run(); err != nil {
		log.Printf("ffmpeg 音频提取失败: %v", err)
		return ExtractAudioRes{
			Success: false,
			Message: "ffmpeg 音频提取失败: " + err.Error(),
		}
	}

	elapsed := time.Since(startTime)
	callback([]byte("音频提取完成，耗时: " + elapsed.String()))

	return ExtractAudioRes{
		Success:    true,
		OutputPath: outputPath,
		Message:    "音频提取完成",
		Cost:       elapsed.String(),
	}
}

// VideoInfo 视频文件信息
type VideoInfo struct {
	FilePath   string  `json:"filePath"`   // 完整路径
	FileName   string  `json:"fileName"`   // 文件名（含扩展名）
	FileSizeMB float64 `json:"fileSizeMB"` // 文件大小 MB
	Resolution string  `json:"resolution"` // 分辨率：宽 x 高
	Width      int     `json:"width"`      // 宽度
	Height     int     `json:"height"`     // 高度
}

// ScanVideosReq 扫描视频请求
type ScanVideosReq struct {
	DirPath string `json:"dirPath"` // 要扫描的目录路径
}

// ScanVideosRes 扫描视频结果
type ScanVideosRes struct {
	Success        bool        `json:"success"`
	Message        string      `json:"message"`
	AllVideos      []VideoInfo `json:"allVideos"`      // 所有视频列表
	VerticalVideos []VideoInfo `json:"verticalVideos"` // 竖屏视频（高 > 宽）
}

// ScanVideos 扫描指定目录中的所有视频文件（不扫描子目录）
func ScanVideos(req ScanVideosReq, callback func([]byte)) ScanVideosRes {
	callback([]byte("开始扫描目录：" + req.DirPath))

	if req.DirPath == "" {
		return ScanVideosRes{
			Success: false,
			Message: "请选择要扫描的目录",
		}
	}

	if !core.FileExist(req.DirPath) {
		return ScanVideosRes{
			Success: false,
			Message: "目录不存在：" + req.DirPath,
		}
	}

	// 常见视频格式扩展名
	videoExts := map[string]bool{
		".mp4":  true,
		".avi":  true,
		".mkv":  true,
		".mov":  true,
		".wmv":  true,
		".flv":  true,
		".m4v":  true,
		".ts":   true,
		".rmvb": true,
		".webm": true,
	}

	var allVideos []VideoInfo
	var verticalVideos []VideoInfo

	// 只遍历当前目录，不递归子目录
	entries, err := os.ReadDir(req.DirPath)
	if err != nil {
		return ScanVideosRes{
			Success: false,
			Message: "无法读取目录：" + err.Error(),
		}
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue // 跳过子目录
		}

		fileName := entry.Name()
		ext := strings.ToLower(filepath.Ext(fileName))
		if !videoExts[ext] {
			continue // 只处理视频文件
		}

		filePath := filepath.Join(req.DirPath, fileName)

		// 获取文件信息
		info, err := entry.Info()
		if err != nil {
			log.Printf("无法获取文件信息：%s, 错误：%v", filePath, err)
			continue
		}

		// 获取文件大小（MB）
		fileSizeMB := float64(info.Size()) / 1024.0 / 1024.0

		// 使用 ffprobe 获取视频分辨率
		width, height, err := getVideoResolution(filePath)
		if err != nil {
			log.Printf("无法获取视频分辨率：%s, 错误：%v", filePath, err)
			continue
		}

		resolution := fmt.Sprintf("%dx%d", width, height)

		videoInfo := VideoInfo{
			FilePath:   filePath,
			FileName:   fileName,
			FileSizeMB: fileSizeMB,
			Resolution: resolution,
			Width:      width,
			Height:     height,
		}

		allVideos = append(allVideos, videoInfo)

		// 判断是否为竖屏视频（高 > 宽）
		isVertical := ""
		if height > width {
			verticalVideos = append(verticalVideos, videoInfo)
			isVertical = " [竖屏]"
		}

		// 每扫描一个文件，发送一条日志到前端
		logMsg := fmt.Sprintf("找到：%s | 大小：%.2f MB | 分辨率：%dx%d%s",
			fileName, fileSizeMB, width, height, isVertical)
		callback([]byte(logMsg))
	}

	if err != nil {
		return ScanVideosRes{
			Success: false,
			Message: "扫描目录失败：" + err.Error(),
		}
	}

	callback([]byte(fmt.Sprintf("扫描完成，共找到 %d 个视频文件，其中竖屏视频 %d 个", len(allVideos), len(verticalVideos))))

	return ScanVideosRes{
		Success:        true,
		Message:        fmt.Sprintf("找到 %d 个视频文件", len(allVideos)),
		AllVideos:      allVideos,
		VerticalVideos: verticalVideos,
	}
}

// getVideoResolution 使用 ffprobe 获取视频的宽度和高度
func getVideoResolution(filePath string) (int, int, error) {
	cmd := exec.Command("ffprobe",
		"-v", "quiet",
		"-print_format", "json",
		"-show_streams",
		filePath)

	// Windows 下隐藏命令行窗口
	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	}

	output, err := cmd.Output()
	if err != nil {
		return 0, 0, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		return 0, 0, err
	}

	streams, ok := result["streams"].([]interface{})
	if !ok || len(streams) == 0 {
		return 0, 0, fmt.Errorf("未找到视频流")
	}

	// 查找第一个视频流
	for _, s := range streams {
		stream, ok := s.(map[string]interface{})
		if !ok {
			continue
		}
		codecType, _ := stream["codec_type"].(string)
		if codecType == "video" {
			widthFloat, _ := stream["width"].(float64)
			heightFloat, _ := stream["height"].(float64)
			return int(widthFloat), int(heightFloat), nil
		}
	}

	return 0, 0, fmt.Errorf("未找到视频流信息")
}

// ExtractFramesReq 抽帧请求
type ExtractFramesReq struct {
	FilePath string `json:"filePath"` // 视频文件路径
	Count    int    `json:"count"`    // 抽取帧数量
}

// ExtractFramesRes 抽帧结果
type ExtractFramesRes struct {
	Success    bool     `json:"success"`
	Message    string   `json:"message"`
	OutputDir  string   `json:"outputDir"`  // 输出目录
	FramePaths []string `json:"framePaths"` // 所有帧的路径
	FrameCount int      `json:"frameCount"` // 实际抽取的帧数
	Cost       string   `json:"cost"`       // 耗时
}

// ExtractFrames 从视频中随机抽取 N 帧
func ExtractFrames(req ExtractFramesReq, callback func([]byte)) ExtractFramesRes {
	callback([]byte("开始抽取视频帧：" + req.FilePath))

	if req.FilePath == "" {
		return ExtractFramesRes{
			Success: false,
			Message: "请选择视频文件",
		}
	}

	if req.Count < 1 || req.Count > 100 {
		return ExtractFramesRes{
			Success: false,
			Message: "抽取数量必须在 1-100 之间",
		}
	}

	startTime := time.Now()

	// 获取视频时长
	duration, err := getVideoDuration(req.FilePath)
	if err != nil {
		return ExtractFramesRes{
			Success: false,
			Message: "获取视频时长失败：" + err.Error(),
		}
	}

	callback([]byte(fmt.Sprintf("视频时长：%.2f 秒", duration)))

	// 生成输出目录
	outputDir := req.FilePath + "_frames"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return ExtractFramesRes{
			Success: false,
			Message: "创建输出目录失败：" + err.Error(),
		}
	}

	// 生成随机时间点
	timestamps := generateRandomTimestamps(duration, req.Count)
	callback([]byte(fmt.Sprintf("生成 %d 个随机时间点", len(timestamps))))

	var framePaths []string

	// 抽取每一帧
	for i, ts := range timestamps {
		framePath := filepath.Join(outputDir, fmt.Sprintf("frame_%03d.jpg", i+1))

		// 使用 ffmpeg 抽取单帧
		cmd := exec.Command("ffmpeg",
			"-ss", fmt.Sprintf("%.3f", ts),
			"-i", req.FilePath,
			"-vf", "select=eq(pict_type\\,I)", // 视频过滤器：只选择帧类型为I的帧
			"-vsync", "0", // 防止帧数同步
			"-vframes", "1",
			"-q:v", "2",
			"-y",
			framePath)

		// Windows 下隐藏命令行窗口
		if runtime.GOOS == "windows" {
			cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		}

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			log.Printf("抽取第 %d 帧失败：%v", i+1, err)
			continue
		}

		framePaths = append(framePaths, framePath)
		callback([]byte(fmt.Sprintf("✅ 已抽取第 %d/%d 帧", i+1, len(timestamps))))
	}

	elapsed := time.Since(startTime)
	callback([]byte(fmt.Sprintf("抽帧完成，共提取 %d 张图片，耗时：%s", len(framePaths), elapsed.String())))

	return ExtractFramesRes{
		Success:    true,
		Message:    fmt.Sprintf("成功提取 %d 张图片", len(framePaths)),
		OutputDir:  outputDir,
		FramePaths: framePaths,
		FrameCount: len(framePaths),
		Cost:       elapsed.String(),
	}
}

// getVideoDuration 使用 ffprobe 获取视频时长（秒）
func getVideoDuration(filePath string) (float64, error) {
	cmd := exec.Command("ffprobe",
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		filePath)

	// Windows 下隐藏命令行窗口
	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	}

	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		return 0, err
	}

	format, ok := result["format"].(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("无法解析视频格式信息")
	}

	durationStr, _ := format["duration"].(string)
	if durationStr == "" {
		return 0, fmt.Errorf("无法获取视频时长")
	}

	return strconv.ParseFloat(durationStr, 64)
}

// generateRandomTimestamps 生成 N 个随机的时间点（秒）
func generateRandomTimestamps(duration float64, count int) []float64 {
	rand.Seed(time.Now().UnixNano())
	timestamps := make([]float64, count)

	// 避免在视频开头和结尾抽取
	safeDuration := duration * 0.9
	if safeDuration < 1 {
		safeDuration = duration
	}

	for i := 0; i < count; i++ {
		timestamps[i] = rand.Float64() * safeDuration
	}

	// 排序时间戳
	sort.Float64s(timestamps)
	return timestamps
}

func deleteFile(file string) {
	err := os.Remove(file)
	if err != nil {
		log.Printf("无法删除文件: %s", preDeleteFile)
	} else {
		log.Printf("文件已成功删除: %s", preDeleteFile)
	}
}

// ExtractVideoThumbnailReq 提取视频缩略图请求
type ExtractVideoThumbnailReq struct {
	FilePath  string  `json:"filePath"`  // 视频文件路径
	Timestamp float64 `json:"timestamp"` // 提取帧的时间点（秒），默认 1.0
}

// ExtractVideoThumbnailRes 提取视频缩略图结果
type ExtractVideoThumbnailRes struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	Thumbnail string `json:"thumbnail"` // Base64 编码的缩略图
	MimeType  string `json:"mimeType"`  // 图片 MIME 类型
}

// ExtractVideoThumbnail 从视频中提取一帧作为缩略图
func ExtractVideoThumbnail(req ExtractVideoThumbnailReq) ExtractVideoThumbnailRes {
	if req.FilePath == "" {
		return ExtractVideoThumbnailRes{
			Success: false,
			Message: "请选择视频文件",
		}
	}

	if !core.FileExist(req.FilePath) {
		return ExtractVideoThumbnailRes{
			Success: false,
			Message: "视频文件不存在",
		}
	}

	// 默认从第1秒提取
	timestamp := req.Timestamp
	if timestamp <= 0 {
		timestamp = 1.0
	}

	// 生成临时缩略图文件路径
	tempDir := core.GetTempDir()
	thumbnailPath := filepath.Join(tempDir, fmt.Sprintf("thumb_%d.jpg", time.Now().UnixNano()))

	// 使用 ffmpeg 提取单帧
	cmd := exec.Command("ffmpeg",
		"-ss", fmt.Sprintf("%.3f", timestamp),
		"-i", req.FilePath,
		"-vframes", "1",
		"-q:v", "2",
		"-y",
		thumbnailPath)

	// Windows 下隐藏命令行窗口
	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Printf("缩略图提取失败: %v", err)
		return ExtractVideoThumbnailRes{
			Success: false,
			Message: "缩略图提取失败: " + err.Error(),
		}
	}

	// 读取缩略图并转为 base64
	fileData, err := os.ReadFile(thumbnailPath)
	if err != nil {
		log.Printf("缩略图读取失败: %v", err)
		return ExtractVideoThumbnailRes{
			Success: false,
			Message: "缩略图读取失败: " + err.Error(),
		}
	}

	// 删除临时文件
	defer os.Remove(thumbnailPath)

	// base64 编码
	encoded := base64.StdEncoding.EncodeToString(fileData)

	return ExtractVideoThumbnailRes{
		Success:   true,
		Message:   "缩略图提取成功",
		Thumbnail: encoded,
		MimeType:  "image/jpeg",
	}
}

// BatchExtractAudioReq 批量提取音频请求
type BatchExtractAudioReq struct {
	DirPath     string `json:"dirPath"`     // 视频目录路径
	Format      string `json:"format"`      // 输出格式: mp3, aac, wav, flac, 默认 mp3
	ThreadCount int    `json:"threadCount"` // 并行线程数，默认 4
}

// BatchExtractAudioRes 批量提取音频结果
type BatchExtractAudioRes struct {
	Success      bool     `json:"success"`
	Message      string   `json:"message"`
	TotalCount   int      `json:"totalCount"`   // 总视频数
	SuccessCount int      `json:"successCount"` // 成功数
	FailedCount  int      `json:"failedCount"`  // 失败数
	FailedFiles  []string `json:"failedFiles"`  // 失败文件列表
	OutputDir    string   `json:"outputDir"`    // 输出目录
	TotalCost    string   `json:"totalCost"`    // 总耗时
}

// BatchExtractAudio 批量从视频中提取音频（多线程）
func BatchExtractAudio(req BatchExtractAudioReq, callback func([]byte)) BatchExtractAudioRes {
	startTime := time.Now()
	callback([]byte("开始批量提取音频..."))

	if req.DirPath == "" {
		return BatchExtractAudioRes{Success: false, Message: "请选择视频目录"}
	}

	if !core.FileExist(req.DirPath) {
		return BatchExtractAudioRes{Success: false, Message: "目录不存在：" + req.DirPath}
	}

	format := req.Format
	if format == "" {
		format = "mp3"
	}

	threadCount := req.ThreadCount
	if threadCount <= 0 {
		threadCount = 4
	}

	// 创建输出目录
	outputDir := filepath.Join(req.DirPath, "music")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return BatchExtractAudioRes{Success: false, Message: "创建输出目录失败：" + err.Error()}
	}

	callback([]byte(fmt.Sprintf("输出目录: %s，并行线程数: %d", outputDir, threadCount)))

	// 常见视频格式扩展名
	videoExts := map[string]bool{
		".mp4":  true,
		".avi":  true,
		".mkv":  true,
		".mov":  true,
		".wmv":  true,
		".flv":  true,
		".m4v":  true,
		".ts":   true,
		".rmvb": true,
		".webm": true,
	}

	// 扫描目录下所有视频文件
	entries, err := os.ReadDir(req.DirPath)
	if err != nil {
		return BatchExtractAudioRes{Success: false, Message: "无法读取目录：" + err.Error()}
	}

	var videoFiles []string
	var extCountMap = make(map[string]int)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		fileName := entry.Name()
		ext := strings.ToLower(filepath.Ext(fileName))
		if !videoExts[ext] {
			continue
		}
		videoFiles = append(videoFiles, filepath.Join(req.DirPath, fileName))
		extCountMap[ext]++
	}

	if len(videoFiles) == 0 {
		return BatchExtractAudioRes{Success: false, Message: "目录中没有找到视频文件"}
	}

	callback([]byte(fmt.Sprintf("共找到 %d 个视频文件", len(videoFiles))))

	// 统计各类型文件数
	var extStats []string
	for ext, count := range extCountMap {
		extStats = append(extStats, fmt.Sprintf("%s: %d个", ext, count))
	}
	callback([]byte("文件类型分布: " + strings.Join(extStats, ", ")))

	// 使用多线程处理
	var mu sync.Mutex
	var successCount, failedCount int
	var failedFiles []string
	var wg sync.WaitGroup

	// 控制并发数
	semaphore := make(chan struct{}, threadCount)

	for i, videoPath := range videoFiles {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(idx int, path string) {
			defer wg.Done()
			defer func() { <-semaphore }()

			callback([]byte(fmt.Sprintf("[%d/%d] 正在处理: %s", idx+1, len(videoFiles), filepath.Base(path))))

			// 生成输出路径
			baseName := filepath.Base(path)
			dotIdx := strings.LastIndex(baseName, ".")
			if dotIdx > 0 {
				baseName = baseName[:dotIdx]
			}
			outputPath := filepath.Join(outputDir, baseName+"."+format)

			// 构建 ffmpeg 命令
			var execCmd *exec.Cmd
			switch format {
			case "wav":
				execCmd = exec.Command("ffmpeg", "-i", path, "-vn", "-acodec", "pcm_s16le", "-y", outputPath)
			case "aac":
				execCmd = exec.Command("ffmpeg", "-i", path, "-vn", "-acodec", "copy", "-y", outputPath)
			case "flac":
				execCmd = exec.Command("ffmpeg", "-i", path, "-vn", "-acodec", "flac", "-y", outputPath)
			default: // mp3
				execCmd = exec.Command("ffmpeg", "-i", path, "-vn", "-acodec", "libmp3lame", "-q:a", "2", "-y", outputPath)
			}

			if runtime.GOOS == "windows" {
				execCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
			}

			execCmd.Stdout = os.Stdout
			execCmd.Stderr = os.Stderr

			if err := execCmd.Run(); err != nil {
				mu.Lock()
				failedCount++
				failedFiles = append(failedFiles, path+" (错误: "+err.Error()+")")
				mu.Unlock()
				callback([]byte(fmt.Sprintf("❌ [%d/%d] 失败: %s", idx+1, len(videoFiles), filepath.Base(path))))
			} else {
				mu.Lock()
				successCount++
				mu.Unlock()
				callback([]byte(fmt.Sprintf("✅ [%d/%d] 完成: %s", idx+1, len(videoFiles), filepath.Base(path))))
			}
		}(i, videoPath)
	}

	wg.Wait()

	elapsed := time.Since(startTime)
	totalCost := fmt.Sprintf("%.1f秒", elapsed.Seconds())

	callback([]byte(fmt.Sprintf("批量提取完成！成功: %d，失败: %d，总耗时: %s", successCount, failedCount, totalCost)))

	return BatchExtractAudioRes{
		Success:      failedCount == 0,
		Message:      fmt.Sprintf("批量提取完成，成功: %d，失败: %d", successCount, failedCount),
		TotalCount:   len(videoFiles),
		SuccessCount: successCount,
		FailedCount:  failedCount,
		FailedFiles:  failedFiles,
		OutputDir:    outputDir,
		TotalCost:    totalCost,
	}
}

// VideoExtInfo 视频扩展名统计信息
type VideoExtInfo struct {
	Ext   string `json:"ext"`   // 扩展名（如 .mp4）
	Count int    `json:"count"` // 数量
}

// ScanVideoDirReq 扫描视频目录请求
type ScanVideoDirReq struct {
	DirPath string `json:"dirPath"` // 要扫描的目录路径
}

// ScanVideoDirRes 扫描视频目录结果
type ScanVideoDirRes struct {
	Success    bool           `json:"success"`
	Message    string         `json:"message"`
	TotalCount int            `json:"totalCount"` // 视频总数
	ExtInfos   []VideoExtInfo `json:"extInfos"`   // 各类型统计
}

// ScanVideoDir 扫描视频目录，统计视频个数和各类型个数
func ScanVideoDir(req ScanVideoDirReq, callback func([]byte)) ScanVideoDirRes {
	callback([]byte("开始扫描目录：" + req.DirPath))

	if req.DirPath == "" {
		return ScanVideoDirRes{Success: false, Message: "请选择目录"}
	}

	if !core.FileExist(req.DirPath) {
		return ScanVideoDirRes{Success: false, Message: "目录不存在：" + req.DirPath}
	}

	// 常见视频格式扩展名
	videoExts := map[string]bool{
		".mp4":  true,
		".avi":  true,
		".mkv":  true,
		".mov":  true,
		".wmv":  true,
		".flv":  true,
		".m4v":  true,
		".ts":   true,
		".rmvb": true,
		".webm": true,
	}

	entries, err := os.ReadDir(req.DirPath)
	if err != nil {
		return ScanVideoDirRes{Success: false, Message: "无法读取目录：" + err.Error()}
	}

	extCountMap := make(map[string]int)
	totalCount := 0

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		fileName := entry.Name()
		ext := strings.ToLower(filepath.Ext(fileName))
		if !videoExts[ext] {
			continue
		}
		extCountMap[ext]++
		totalCount++
	}

	var extInfos []VideoExtInfo
	for ext, count := range extCountMap {
		extInfos = append(extInfos, VideoExtInfo{Ext: ext, Count: count})
	}

	// 按数量降序排序
	sort.Slice(extInfos, func(i, j int) bool {
		return extInfos[i].Count > extInfos[j].Count
	})

	callback([]byte(fmt.Sprintf("扫描完成，共找到 %d 个视频文件", totalCount)))

	return ScanVideoDirRes{
		Success:    true,
		Message:    fmt.Sprintf("共 %d 个视频", totalCount),
		TotalCount: totalCount,
		ExtInfos:   extInfos,
	}
}
