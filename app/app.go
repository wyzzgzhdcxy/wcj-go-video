package app

import (
	"bufio"
	"bytes"
	"compress/flate"
	"compress/gzip"
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
	"wcj-go-video/golang/myVideo"
	"wcj-go-video/golang/cmdWrapper"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/wyzzgzhdcxy/wcj-go-common/core"
	"github.com/wyzzgzhdcxy/wcj-go-common/model"
	_ "modernc.org/sqlite"
)

// App struct
type App struct {
	ctx         context.Context
	SettingsDb  *sql.DB // 配置存储的 sqlite 数据库
	ProjectName string
	Assets      embed.FS
}

// InitSettingsDb 初始化配置数据库（sqlite）
func (a *App) InitSettingsDb() error {
	dbPath := core.GetTempDir() + "/data/video_settings.db"
	// 确保目录存在
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return fmt.Errorf("创建数据库目录失败: %v", err)
	}

	// 使用 database/sql + sqlite 驱动
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("打开数据库失败: %v", err)
	}

	// 设置 busy_timeout 避免数据库锁定错误（30秒）
	_, err = db.Exec("PRAGMA busy_timeout = 30000")
	if err != nil {
		db.Close()
		return fmt.Errorf("设置busy_timeout失败: %v", err)
	}

	// 启用 WAL 模式提升并发性能
	_, err = db.Exec("PRAGMA journal_mode = WAL")
	if err != nil {
		db.Close()
		return fmt.Errorf("设置WAL模式失败: %v", err)
	}

	// 创建配置表
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS settings (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			key TEXT NOT NULL UNIQUE,
			value TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		db.Close()
		return fmt.Errorf("创建表失败: %v", err)
	}

	// 创建索引
	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_settings_key ON settings(key)`)
	if err != nil {
		db.Close()
		return fmt.Errorf("创建索引失败: %v", err)
	}

	// 创建 emoji 图片表
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS emoji_images (
			id TEXT PRIMARY KEY,
			png_data BLOB NOT NULL,
			ico_data BLOB NOT NULL,
			emoji TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		db.Close()
		return fmt.Errorf("创建 emoji 图片表失败: %v", err)
	}

	a.SettingsDb = db
	log.Printf("配置数据库初始化成功: %s", dbPath)
	return nil
}

// NewApp creates a new App application struct
func NewApp(assets embed.FS) *App {
	return &App{
		Assets: assets,
	}
}

func (a *App) GetTempDir() string {
	return core.GetTempDir()
}

func (a *App) GetVideoFilesInDir(filePath string) []string {
	dir := filepath.Dir(filePath)
	var videoFiles []string
	var videoExts = []string{".mp4", ".avi", ".mov", ".mkv", ".flv", ".wmv", ".webm", ".mpg", ".mpeg"}
	_ = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		for _, ve := range videoExts {
			if ext == ve {
				videoFiles = append(videoFiles, path)
				break
			}
		}
		return nil
	})
	return videoFiles
}

// Startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	log.Printf("startup执行开始!%s", time.Now().Format("2006-01-02 15:04:05.000"))

	// 恢复窗口位置和大小（如果之前有保存）
	ws := a.GetWindowState()
	xSaved := a.getSetting("window.x_saved")
	if ws.Width > 0 && ws.Height > 0 {
		// 获取屏幕尺寸
		screens, err := runtime.ScreenGetAll(ctx)
		if err == nil && len(screens) > 0 {
			screenWidth := screens[0].Size.Width
			screenHeight := screens[0].Size.Height

			// 如果窗口尺寸超过屏幕尺寸-100px，则调整窗口尺寸
			maxWidth := screenWidth - 100
			maxHeight := screenHeight - 100
			newWidth := ws.Width
			newHeight := ws.Height
			if newWidth > maxWidth {
				newWidth = maxWidth
			}
			if newHeight > maxHeight {
				newHeight = maxHeight
			}
			runtime.WindowSetSize(ctx, newWidth, newHeight)
		}
	}
	// 如果窗口位置曾经被保存过，则恢复位置
	if xSaved == "1" {
		runtime.WindowSetPosition(ctx, ws.X, ws.Y)
	} else if ws.Width > 0 && ws.Height > 0 {
		runtime.WindowCenter(ctx)
	}

	log.Printf("startup执行结束! %s", time.Now().Format("2006-01-02 15:04:05.000"))
}

func (a *App) Shutdown(ctx context.Context) {
	if a.SettingsDb != nil {
		a.SettingsDb.Close()
	}
	log.Printf("shutdown执行结束!")
}

// saveSetting 保存配置到 sqlite
func (a *App) saveSetting(key, value string) {
	if a.SettingsDb == nil {
		return
	}
	_, err := a.SettingsDb.Exec(`
		INSERT INTO settings (key, value, updated_at)
		VALUES (?, ?, datetime('now'))
		ON CONFLICT(key) DO UPDATE SET value = excluded.value, updated_at = excluded.updated_at
	`, key, value)
	if err != nil {
		log.Printf("保存配置失败: %v", err)
	}
}

// getSetting 从 sqlite 获取配置
func (a *App) getSetting(key string) string {
	if a.SettingsDb == nil {
		return ""
	}
	var value string
	err := a.SettingsDb.QueryRow("SELECT value FROM settings WHERE key = ?", key).Scan(&value)
	if err != nil {
		return ""
	}
	return value
}

// WindowState 窗口状态
type WindowState struct {
	X      int
	Y      int
	Width  int
	Height int
}

// SaveWindowState 保存窗口状态到 sqlite
func (a *App) SaveWindowState(x, y, width, height int) {
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	if width < 600 {
		width = 600
	}
	if height < 400 {
		height = 400
	}
	a.saveSetting("window.x", fmt.Sprintf("%d", x))
	a.saveSetting("window.y", fmt.Sprintf("%d", y))
	a.saveSetting("window.width", fmt.Sprintf("%d", width))
	a.saveSetting("window.height", fmt.Sprintf("%d", height))
	a.saveSetting("window.x_saved", "1")
}

// GetWindowState 从 sqlite 获取窗口状态
func (a *App) GetWindowState() WindowState {
	ws := WindowState{}
	if a.SettingsDb == nil {
		return ws
	}
	xSaved := a.getSetting("window.x_saved")
	if xSaved != "1" {
		return ws
	}
	x := a.getSetting("window.x")
	if x == "" {
		return ws
	}
	ws.X, _ = strconv.Atoi(x)
	ws.Y, _ = strconv.Atoi(a.getSetting("window.y"))
	ws.Width, _ = strconv.Atoi(a.getSetting("window.width"))
	ws.Height, _ = strconv.Atoi(a.getSetting("window.height"))
	if ws.X < 0 {
		ws.X = 0
	}
	if ws.Y < 0 {
		ws.Y = 0
	}
	if ws.Width < 600 {
		ws.Width = 600
	}
	if ws.Height < 400 {
		ws.Height = 400
	}
	return ws
}

func (a *App) SendToFrontend(message []byte) {
	fmt.Println("send to frontend:", string(message))
	backMsg := model.BackMsg{
		Time: core.GetTime(),
		Msg:  string(message),
	}
	runtime.EventsEmit(a.ctx, "back_msg", core.ToJsonString(backMsg))
}

func (a *App) ReadAssetFile(path string) (string, error) {
	content, err := fs.ReadFile(a.Assets, "frontend/dist"+path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// OpenExplorer 打开文件资源管理器并定位到指定路径
func (a *App) OpenExplorer(path string) error {
	return cmdWrapper.OpenExplorer(path)
}

// SelectDirectory 打开目录选择对话框
func (a *App) SelectDirectory() (string, error) {
	selection, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择目录",
	})
	if err != nil {
		return "", err
	}
	return selection, nil
}

// SelectVideoFile 打开视频文件选择对话框
func (a *App) SelectVideoFile() (string, error) {
	selection, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择视频文件",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "视频文件 (*.mp4;*.avi;*.mkv;*.mov;*.wmv;*.flv;*.m4v)",
				Pattern:     "*.mp4;*.avi;*.mkv;*.mov;*.wmv;*.flv;*.m4v",
			},
		},
	})
	if err != nil {
		return "", err
	}
	return selection, nil
}

// RotateVideo 旋转视频文件
func (a *App) RotateVideo(req myVideo.RotateVideoReq) myVideo.RotateVideoRes {
	return myVideo.RotateVideo(req, a.SendToFrontend)
}

// ClipVideo 裁剪视频文件
func (a *App) ClipVideo(req myVideo.ClipVideoReq) myVideo.ClipVideoRes {
	return myVideo.ClipVideo(req, a.SendToFrontend)
}

// ExtractAudio 从视频中分离音频
func (a *App) ExtractAudio(req myVideo.ExtractAudioReq) myVideo.ExtractAudioRes {
	return myVideo.ExtractAudio(req, a.SendToFrontend)
}

// ScanVideos 扫描目录中的视频文件
func (a *App) ScanVideos(req myVideo.ScanVideosReq) myVideo.ScanVideosRes {
	return myVideo.ScanVideos(req, a.SendToFrontend)
}

// BatchExtractAudio 批量从视频中提取音频
func (a *App) BatchExtractAudio(req myVideo.BatchExtractAudioReq) myVideo.BatchExtractAudioRes {
	return myVideo.BatchExtractAudio(req, a.SendToFrontend)
}

// ScanVideoDir 扫描视频目录，统计视频个数和各类型个数
func (a *App) ScanVideoDir(req myVideo.ScanVideoDirReq) myVideo.ScanVideoDirRes {
	return myVideo.ScanVideoDir(req, a.SendToFrontend)
}

// ExtractFrames 从视频中抽取帧
func (a *App) ExtractFrames(req myVideo.ExtractFramesReq) myVideo.ExtractFramesRes {
	return myVideo.ExtractFrames(req, a.SendToFrontend)
}

// ExtractVideoThumbnail 从视频中提取缩略图
func (a *App) ExtractVideoThumbnail(req myVideo.ExtractVideoThumbnailReq) myVideo.ExtractVideoThumbnailRes {
	return myVideo.ExtractVideoThumbnail(req)
}

// ClipAndReplaceVideo 裁剪视频并替换原文件
func (a *App) ClipAndReplaceVideo(req myVideo.ClipAndReplaceReq) myVideo.ClipAndReplaceRes {
	return myVideo.ClipAndReplaceVideo(req, a.SendToFrontend)
}

// RemoveVideoIntro 去除视频片头
func (a *App) RemoveVideoIntro(req myVideo.RemoveVideoIntroReq) myVideo.RemoveVideoIntroRes {
	return myVideo.RemoveVideoIntro(req, a.SendToFrontend)
}

// BatchRemoveVideoIntro 批量去除视频片头
func (a *App) BatchRemoveVideoIntro(req myVideo.BatchRemoveVideoIntroReq) myVideo.BatchRemoveVideoIntroRes {
	return myVideo.BatchRemoveVideoIntro(req, a.SendToFrontend)
}

// ArchiveMobileVideos 将宽小于高的视频（竖屏/手机视频）移动到同目录下的 mobile 文件夹
func (a *App) ArchiveMobileVideos(req ArchiveMobileVideosReq) ArchiveMobileVideosRes {
	if req.DirPath == "" {
		return ArchiveMobileVideosRes{
			Success: false,
			Message: "请选择要归档的目录",
		}
	}

	// 先扫描该目录获取竖屏视频
	scanResult := myVideo.ScanVideos(myVideo.ScanVideosReq{DirPath: req.DirPath}, a.SendToFrontend)
	if !scanResult.Success {
		return ArchiveMobileVideosRes{
			Success: false,
			Message: "扫描目录失败：" + scanResult.Message,
		}
	}

	verticalVideos := scanResult.VerticalVideos
	if len(verticalVideos) == 0 {
		return ArchiveMobileVideosRes{
			Success:    true,
			Message:    "没有找到竖屏视频，无需归档",
			MovedCount: 0,
		}
	}

	// 创建 mobile 文件夹（与当前目录同级别）
	mobileDir := filepath.Join(req.DirPath, "mobile")
	if err := os.MkdirAll(mobileDir, 0755); err != nil {
		return ArchiveMobileVideosRes{
			Success: false,
			Message: "创建 mobile 文件夹失败：" + err.Error(),
		}
	}

	a.SendToFrontend([]byte(fmt.Sprintf("找到 %d 个竖屏视频，开始归档到 mobile 文件夹...", len(verticalVideos))))

	movedCount := 0
	var failedFiles []string

	for _, video := range verticalVideos {
		fileName := filepath.Base(video.FilePath)
		targetFile := filepath.Join(mobileDir, fileName)

		a.SendToFrontend([]byte(fmt.Sprintf("正在移动：%s -> %s", fileName, mobileDir)))

		// 检查目标文件是否已存在
		if core.FileExist(targetFile) {
			ext := filepath.Ext(targetFile)
			base := strings.TrimSuffix(targetFile, ext)
			timestamp := time.Now().Format("20060102_150405")
			targetFile = fmt.Sprintf("%s_%s%s", base, timestamp, ext)
		}

		// 移动文件
		if err := os.Rename(video.FilePath, targetFile); err != nil {
			failedFiles = append(failedFiles, fileName+" ("+err.Error()+")")
			a.SendToFrontend([]byte(fmt.Sprintf("❌ 移动失败：%s - %v", fileName, err)))
			continue
		}

		movedCount++
		a.SendToFrontend([]byte(fmt.Sprintf("✅ 移动成功：%s", fileName)))
	}

	msg := fmt.Sprintf("共 %d 个竖屏视频，成功归档 %d 个", len(verticalVideos), movedCount)
	if len(failedFiles) > 0 {
		msg += fmt.Sprintf("，失败 %d 个", len(failedFiles))
	}

	return ArchiveMobileVideosRes{
		Success:     true,
		Message:     msg,
		MovedCount:  movedCount,
		FailedFiles: failedFiles,
	}
}

// ArchiveMobileVideosReq 归档手机视频请求
type ArchiveMobileVideosReq struct {
	DirPath string `json:"dirPath"`
}

// ArchiveMobileVideosRes 归档手机视频结果
type ArchiveMobileVideosRes struct {
	Success     bool     `json:"success"`
	Message     string   `json:"message"`
	MovedCount  int      `json:"movedCount"`
	FailedFiles []string `json:"failedFiles"`
}

// ExtractM3u8LinksReq 提取m3u8链接请求
type ExtractM3u8LinksReq struct {
	Url string `json:"url"`
}

// ExtractM3u8LinksRes 提取m3u8链接结果
type ExtractM3u8LinksRes struct {
	Success bool       `json:"success"`
	Message string     `json:"message"`
	Links   []M3u8Info `json:"links"`
}

// M3u8Info m3u8视频信息
type M3u8Info struct {
	Url   string `json:"url"`
	Title string `json:"title"`
}

// ExtractM3u8Links 从网页中提取m3u8链接
func (a *App) ExtractM3u8Links(req ExtractM3u8LinksReq) ExtractM3u8LinksRes {
	if req.Url == "" {
		return ExtractM3u8LinksRes{
			Success: false,
			Message: "请输入URL",
		}
	}

	client := &http.Client{
		Transport: &http.Transport{
			DisableCompression: true,
		},
	}
	httpReq, err := http.NewRequest("GET", req.Url, nil)
	if err != nil {
		return ExtractM3u8LinksRes{
			Success: false,
			Message: "创建请求失败: " + err.Error(),
		}
	}

	httpReq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	httpReq.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	httpReq.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	httpReq.Header.Set("Accept-Encoding", "gzip, deflate")
	httpReq.Header.Set("Connection", "keep-alive")

	resp, err := client.Do(httpReq)
	if err != nil {
		return ExtractM3u8LinksRes{
			Success: false,
			Message: "获取网页失败: " + err.Error(),
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ExtractM3u8LinksRes{
			Success: false,
			Message: "读取网页内容失败: " + err.Error(),
		}
	}

	contentEncoding := resp.Header.Get("Content-Encoding")
	switch contentEncoding {
	case "gzip":
		reader, err := gzip.NewReader(bytes.NewReader(body))
		if err == nil {
			body, _ = io.ReadAll(reader)
			reader.Close()
		}
	case "deflate":
		reader := flate.NewReader(bytes.NewReader(body))
		body, _ = io.ReadAll(reader)
		reader.Close()
	}

	content := string(body)

	var links []M3u8Info
	seen := make(map[string]bool)

	m3u8Pattern := regexp.MustCompile(`(https?://[^\s"'\)]+\.m3u8[^\s"'\)]*)`)
	matches := m3u8Pattern.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if len(match) > 1 {
			url := strings.Trim(match[1], "\"'\r\n")
			if !seen[url] && strings.Contains(url, ".m3u8") {
				seen[url] = true
				title := filepath.Base(url)
				title = strings.TrimSuffix(title, ".m3u8")
				links = append(links, M3u8Info{Url: url, Title: title})
			}
		}
	}

	if len(links) == 0 {
		return ExtractM3u8LinksRes{
			Success: false,
			Message: "未找到m3u8链接",
		}
	}

	return ExtractM3u8LinksRes{
		Success: true,
		Message: fmt.Sprintf("找到 %d 个m3u8链接", len(links)),
		Links:   links,
	}
}

// MovieDownloadReq 电影下载请求
type MovieDownloadReq struct {
	MType      string `json:"mType"`
	MID        string `json:"mid"`
	FilterStr  string `json:"filterStr"`
}

// MovieDownload 下载电影
func (a *App) MovieDownload(req MovieDownloadReq) {
	srcPath, _ := os.Executable()
	srcDir := filepath.Dir(srcPath)
	cliPath := filepath.Join(srcDir, "movieDownload.exe")
	if !core.FileExist(cliPath) {
		a.SendToFrontend([]byte("错误: movieDownload.exe 不存在，请将CLI放在主程序同目录下"))
		return
	}
	args := []string{req.MType, req.MID}
	if req.FilterStr != "" {
		args = append(args, req.FilterStr)
	}
	cmd := exec.Command(cliPath, args...)
	cmd.Dir = srcDir
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		a.SendToFrontend([]byte("错误: 创建stdout管道失败 " + err.Error()))
		return
	}
	err = cmd.Start()
	if err != nil {
		a.SendToFrontend([]byte("错误: 启动CLI失败 " + err.Error()))
		return
	}
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		a.SendToFrontend([]byte(scanner.Text()))
	}
	cmd.Wait()
}

// MovieMergeReq 电影合并请求
type MovieMergeReq struct {
	Dir string `json:"dir"`
}

// MovieMerge 合并电影
func (a *App) MovieMerge(req MovieMergeReq) {
	myVideo.MergeVideo(req.Dir, false, a.SendToFrontend)
}

// MoveFilesReq 移动文件请求
type MoveFilesReq struct {
	SourceFiles []string `json:"sourceFiles"`
	TargetDir   string   `json:"targetDir"`
}

// MoveFilesRes 移动文件结果
type MoveFilesRes struct {
	Success     bool     `json:"success"`
	Message     string   `json:"message"`
	MovedCount  int      `json:"movedCount"`
	FailedFiles []string `json:"failedFiles"`
}

// MoveFiles 批量移动文件
func (a *App) MoveFiles(req MoveFilesReq) MoveFilesRes {
	if len(req.SourceFiles) == 0 {
		return MoveFilesRes{Success: false, Message: "请选择要移动的文件"}
	}
	if req.TargetDir == "" {
		return MoveFilesRes{Success: false, Message: "请选择目标目录"}
	}
	if err := os.MkdirAll(req.TargetDir, 0755); err != nil {
		return MoveFilesRes{Success: false, Message: "创建目标目录失败：" + err.Error()}
	}

	movedCount := 0
	var failedFiles []string

	for _, srcFile := range req.SourceFiles {
		fileName := filepath.Base(srcFile)
		targetFile := filepath.Join(req.TargetDir, fileName)

		a.SendToFrontend([]byte(fmt.Sprintf("正在移动：%s -> %s", fileName, targetFile)))

		if core.FileExist(targetFile) {
			ext := filepath.Ext(targetFile)
			base := strings.TrimSuffix(targetFile, ext)
			timestamp := time.Now().Format("20060102_150405")
			targetFile = fmt.Sprintf("%s_%s%s", base, timestamp, ext)
		}

		if err := os.Rename(srcFile, targetFile); err != nil {
			failedFiles = append(failedFiles, fileName+" ("+err.Error()+")")
			a.SendToFrontend([]byte(fmt.Sprintf("移动失败：%s - %v", fileName, err)))
			continue
		}

		movedCount++
		a.SendToFrontend([]byte(fmt.Sprintf("移动成功：%s", fileName)))
	}

	msg := fmt.Sprintf("共移动 %d 个文件，成功 %d 个", len(req.SourceFiles), movedCount)
	if len(failedFiles) > 0 {
		msg += fmt.Sprintf("，失败 %d 个", len(failedFiles))
	}

	return MoveFilesRes{
		Success:     true,
		Message:     msg,
		MovedCount:  movedCount,
		FailedFiles: failedFiles,
	}
}
