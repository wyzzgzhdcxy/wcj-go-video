package main

import (
	"embed"
	"log"
	"time"
	"wcj-go-video/app"

	"github.com/wyzzgzhdcxy/wcj-go-common/utils"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	utils.InitLog(true)

	application := app.NewApp(assets)
	log.Printf("%s", "log init finish! "+time.Now().Format("2006-01-02 15:04:05.000"))

	// 初始化配置数据库（sqlite）
	log.Printf("开始初始化配置数据库...")
	if err := application.InitSettingsDb(); err != nil {
		log.Printf("初始化配置数据库失败: %v", err)
	}
	log.Printf("配置数据库初始化完成")

	// 默认窗口尺寸
	defaultWidth := 1200
	defaultHeight := 800

	// 优先从 SQLite 加载窗口状态，否则使用默认尺寸
	ws := application.GetWindowState()
	width := ws.Width
	height := ws.Height
	if width == 0 || height == 0 {
		width = defaultWidth
		height = defaultHeight
	}

	// Create application with options
	err := wails.Run(&options.App{
		Title:         "视频工具箱",
		Width:         width,
		DisableResize: false,
		Height:        height,
		Frameless:     false,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		DragAndDrop: &options.DragAndDrop{
			EnableFileDrop:     true,
			DisableWebViewDrop: true,
		},
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 1},
		OnStartup:        application.Startup,
		OnShutdown:       application.Shutdown,
		Bind: []interface{}{
			application,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
	log.Printf("%s", "main start finish! "+time.Now().Format("2006-01-02 15:04:05.000"))
}
