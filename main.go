package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v2"
)

// Config 结构体用于解析配置文件
type Config struct {
	Path     string `yaml:"path"`
	Days     int    `yaml:"days"`
	Interval string `yaml:"interval"`
	LogFile  string `yaml:"log_file"`
}

// readConfig 从指定路径读取配置文件
func readConfig(configPath string) (*Config, error) {
	config := &Config{}

	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(file, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// deleteOldFiles 删除指定目录中早于指定天数的文件
func deleteOldFiles(root string, days int) (int, int64, error) {
	var filesDeleted int
	var totalSizeDeleted int64
	var dirsToDelete []string
	cutoff := time.Now().AddDate(0, 0, -days)

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			// 如果是目录，添加到待检查列表
			if path != root {
				dirsToDelete = append(dirsToDelete, path)
			}
			return nil
		}

		if info.ModTime().Before(cutoff) {
			size := info.Size()
			if err := os.Remove(path); err != nil {
				return err
			}
			filesDeleted++
			totalSizeDeleted += size
		}

		return nil
	})

	if err != nil {
		return filesDeleted, totalSizeDeleted, err
	}

	// 删除空目录
	for _, dir := range dirsToDelete {
		if isEmpty, _ := isDirEmpty(dir); isEmpty {
			os.Remove(dir)
		}
	}

	return filesDeleted, totalSizeDeleted, nil
}

// isDirEmpty 检查目录是否为空
func isDirEmpty(dir string) (bool, error) {
	f, err := os.Open(dir)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // 尝试读取至少一个项目
	if err == io.EOF {
		return true, nil
	}
	return false, err // 或许是权限错误或其他错误
}

// runAndLog 执行删除操作并记录日志
func runAndLog(config *Config, nextRun time.Time) {
	fmt.Println("运行检查：", time.Now().Format("2006-01-02 15:04:05"))
	filesDeleted, totalSizeDeleted, err := deleteOldFiles(config.Path, config.Days)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		appendToLog(config.LogFile, "Error: "+err.Error())
	} else {
		logMessage := fmt.Sprintf("检查完成。删除文件数：%d，总大小：%d字节。下次检查时间：%s",
			filesDeleted, totalSizeDeleted, nextRun.Format("2006-01-02 15:04:05"))
		fmt.Println(logMessage)
		appendToLog(config.LogFile, logMessage)
	}
}

// appendToLog 向日志文件追加信息
func appendToLog(filename, message string) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("无法打开日志文件:", err)
		return
	}
	defer f.Close()

	if _, err := f.WriteString(time.Now().Format("2006-01-02 15:04:05") + " - " + message + "\n"); err != nil {
		fmt.Println("无法写入日志文件:", err)
	}
}

// main 是程序的入口点
func main() {
	config, err := readConfig("config.yaml")
	if err != nil {
		fmt.Printf("Error reading config: %v\n", err)
		return
	}

	duration, err := time.ParseDuration(config.Interval)
	if err != nil {
		fmt.Printf("Error parsing interval: %v\n", err)
		return
	}

	fmt.Println("程序开始运行。正在监控：", config.Path)

	nextRun := time.Now().Add(duration)
	runAndLog(config, nextRun) // 首次运行

	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			nextRun = time.Now().Add(duration)
			runAndLog(config, nextRun)
		}
	}
}
