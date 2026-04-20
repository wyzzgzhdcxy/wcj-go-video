package myVideo

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
	"github.com/wyzzgzhdcxy/wcj-go-common/core"
)

// MergeVideo 将符合格式的短视频合并
func MergeVideo(dir string, ignoreSort bool, callback func([]byte)) {
	files, err := os.ReadDir(dir) //读取目录下文件
	if err != nil {
		return
	}
	for _, file := range files {
		if file.IsDir() {
			startTime := time.Now() // 当前时间
			callback([]byte(file.Name() + "开始处理"))
			var nameList []string
			subFiles, _ := os.ReadDir(dir + "/" + file.Name()) //读取目录下文件
			for _, subFile := range subFiles {
				if strings.HasSuffix(subFile.Name(), ".mp4") {
					nameList = append(nameList, subFile.Name())
				}
			}
			log.Printf(file.Name() + ",nameList:" + core.ToJsonString(nameList))
			//按名称重低到高排序
			if !ignoreSort {
				sort.Slice(nameList, func(i, j int) bool {
					name1 := strings.Split(nameList[i], ".")[0]
					name2 := strings.Split(nameList[j], ".")[0]
					return core.StrToInt(name1) < core.StrToInt(name2)
				})
			}

			var builder strings.Builder
			for _, v := range nameList {
				builder.WriteString("file " + dir + "/" + file.Name() + "/" + v + "\n")
			}

			core.WriteStrToFile(dir+"/"+file.Name()+".txt", builder.String())

			//执行命令行,合并文件
			cmdStr := "ffmpeg -f concat -safe 0 -i " + file.Name() + ".txt -c copy " + file.Name() + ".mp4"
			log.Printf("执行命令:" + cmdStr)
			Execute(cmdStr, dir)
			elapsed := time.Now().Sub(startTime)
			callback([]byte(file.Name() + "处理结束,耗时:" + elapsed.String()))
		} else {
			fmt.Println("文件类型异常:" + file.Name())
		}
	}
}

func Execute(cmdStr string, dir string) int {
	cmd := exec.Command("cmd", "/c", cmdStr)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if core.StrIsNotEmpty(dir) {
		cmd.Dir = dir
	}
	//cmd.SysProcAttr = &syscall.SysProcAttr{
	//	HideWindow: true, // 这将创建一个新的会话组
	//}
	err := cmd.Start()
	pid := cmd.Process.Pid
	if err != nil {
		fmt.Printf("执行命令错误：%s\n", err)
	}
	// 等待命令执行完成
	err = cmd.Wait()
	if err != nil {
		fmt.Printf("执行命令错误：%s\n", err)
	}
	return pid
}
