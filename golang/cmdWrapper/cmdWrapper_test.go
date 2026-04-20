package cmdWrapper

import (
	"bytes"
	"os"
	"runtime"
	"testing"
)

func TestRunWithOutput(t *testing.T) {
	// 测试一个简单的命令
	output, err := RunWithOutput("echo", "hello world")
	if err != nil {
		t.Errorf("RunWithOutput failed: %v", err)
	}

	// 检查输出是否包含预期的内容
	// 注意：不同平台的echo输出可能不同
	if !bytes.Contains([]byte(output), []byte("hello world")) {
		t.Errorf("Expected output to contain 'hello world', got: %s", output)
	}
}

func TestRunWithDir(t *testing.T) {
	// 获取当前工作目录
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	// 创建一个临时目录测试
	tmpDir := t.TempDir()

	// 测试在指定目录中运行命令
	err = RunWithDir(tmpDir, "pwd")
	if err != nil {
		t.Errorf("RunWithDir failed: %v", err)
	}

	// 切换回原目录
	err = os.Chdir(cwd)
	if err != nil {
		t.Errorf("Failed to change back to original directory: %v", err)
	}
}

func TestRunWithOptions(t *testing.T) {
	// 测试使用自定义stdout
	var stdout bytes.Buffer
	err := RunWithOptions("echo", []string{"test output"}, WithStdout(&stdout))
	if err != nil {
		t.Errorf("RunWithOptions failed: %v", err)
	}

	output := stdout.String()
	if !bytes.Contains([]byte(output), []byte("test")) {
		t.Errorf("Expected stdout to contain 'test', got: %s", output)
	}
}

func TestStart(t *testing.T) {
	// 测试Start函数（启动但不等待）
	if runtime.GOOS == "windows" {
		// 在Windows上测试notepad或类似的简单命令
		err := Start("cmd", "/c", "echo", "test")
		if err != nil {
			t.Errorf("Start failed: %v", err)
		}
	} else {
		// 在Unix-like系统上测试
		err := Start("echo", "test")
		if err != nil {
			t.Errorf("Start failed: %v", err)
		}
	}
}
