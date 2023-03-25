package resizeVideo

import (
	"fmt"
	"github.com/zhangyiming748/GetAllFolder"
	"github.com/zhangyiming748/GetFileInfo"
	"github.com/zhangyiming748/replace"
	"github.com/zhangyiming748/voiceAlert"
	"golang.org/x/exp/slog"
	"io"
	"os"
	"os/exec"
	"strings"
)

func init() {
	logLevel := os.Getenv("LEVEL")
	//var level slog.Level
	var opt slog.HandlerOptions
	switch logLevel {
	case "Debug":
		opt = slog.HandlerOptions{ // 自定义option
			AddSource: true,
			Level:     slog.LevelDebug, // slog 默认日志级别是 info
		}
	case "Info":
		opt = slog.HandlerOptions{ // 自定义option
			AddSource: true,
			Level:     slog.LevelInfo, // slog 默认日志级别是 info
		}
	case "Warn":
		opt = slog.HandlerOptions{ // 自定义option
			AddSource: true,
			Level:     slog.LevelWarn, // slog 默认日志级别是 info
		}
	case "Err":
		opt = slog.HandlerOptions{ // 自定义option
			AddSource: true,
			Level:     slog.LevelError, // slog 默认日志级别是 info
		}
	default:
		slog.Warn("需要正确设置环境变量 Debug,Info,Warn or Err")
		slog.Info("默认使用Debug等级")
		opt = slog.HandlerOptions{ // 自定义option
			AddSource: true,
			Level:     slog.LevelDebug, // slog 默认日志级别是 info
		}

	}
	file := "GetFileInfo.log"
	logf, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		panic(err)
	}
	//defer logf.Close() //如果不关闭可能造成内存泄露
	logger := slog.New(opt.NewJSONHandler(io.MultiWriter(logf, os.Stdout)))
	slog.SetDefault(logger)
}
func ResizeAllVideos(root, pattern, threads string) {
	ResizeVideos(root, pattern, threads)
	folders := GetAllFolder.ListFolders(root)
	for i, folder := range folders {
		slog.Info(fmt.Sprintf("正在处理第 %d/%d 个文件夹", i+1, len(folders)))
		ResizeVideos(folder, pattern, threads)
	}
}
func ResizeVideos(src, pattern, threads string) {
	files := GetFileInfo.GetAllVideoFileInfo(src, pattern)
	for _, file := range files {
		if file.Width <= 1920 || file.Height <= 1920 {
			slog.Info("跳过", slog.Any("正常尺寸的视频", file.FullPath))
			continue
		}
		if file.Width > file.Height {
			slog.Info("横屏视频", slog.Any("视频信息", file))
			ResizeVideo(file, threads, "1920x1080")
		} else if file.Width < file.Height {
			slog.Info("竖屏视频", slog.Any("视频信息", file))
			ResizeVideo(file, threads, "1080x1920")
		} else {
			slog.Info("正方形视频", slog.Any("视频信息", file))
			ResizeVideo(file, threads, "1920x1920")
		}
	}
}
func ResizeVideo(in GetFileInfo.Info, threads string, p string) {
	defer func() {
		if err := recover(); err != nil {
			slog.Warn("错误", slog.Any("文件信息", in.FullPath))
			voiceAlert.Customize("failed", voiceAlert.Samantha)
		}
	}()
	dst := strings.Trim(in.FullPath, in.FullName)   //原始目录
	dst = strings.Join([]string{dst, "resize"}, "") //二级目录
	fname := strings.Trim(in.FullName, in.ExtName)  //仅文件名
	mp4 := strings.Join([]string{fname, "mp4"}, ".")
	os.Mkdir(dst, 0777)
	out := strings.Join([]string{dst, mp4}, string(os.PathSeparator))
	slog.Info("io", slog.Any("源文件:", in.FullPath), slog.Any("输出文件:", out))
	//ffmpeg -i 1.mp4 -strict -2 -vf scale=-1:1080 4.mp4
	// ffmpeg -threads 2 -i 4k_Saeko_Limo.mp4 -strict -2 -vf scale=-1:1080 -c:v libx265 -threads 2 1080.mp4
	var cmd *exec.Cmd
	switch p {
	case "1920x1080":
		cmd = exec.Command("ffmpeg", "-threads", threads, "-i", in.FullPath, "-strict", "-2", "-vf", "scale=-1:1080", "-threads", threads, out)
	case "1080x1920":
		cmd = exec.Command("ffmpeg", "-threads", threads, "-i", in.FullPath, "-strict", "-2", "-vf", "scale=-1:1920", "-threads", threads, out)
	case "1920x1920":
		cmd = exec.Command("ffmpeg", "-threads", threads, "-i", in.FullPath, "-strict", "-2", "-vf", "scale=1920:1920", "-threads", threads, out)
	default:
		slog.Warn("不正常的视频源", slog.Any("视频信息", in.FullPath))
	}
	slog.Info("ffmpeg", slog.Any("生成的命令", fmt.Sprintf("生成的命令是:%s\n", cmd)))
	stdout, err := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout
	if err != nil {
		slog.Warn("cmd.StdoutPipe产生错误", err)
		return
	}
	if err = cmd.Start(); err != nil {
		slog.Warn("cmd.Run产生的错误", err)
		return
	}
	for {
		tmp := make([]byte, 1024)
		_, err := stdout.Read(tmp)
		//写成输出日志
		t := string(tmp)
		t = replace.Replace(t)
		fmt.Println(t)
		if err != nil {
			break
		}
	}
	if err = cmd.Wait(); err != nil {
		slog.Warn("错误", slog.Any("命令执行中", err))
		return
	}
	if err := os.Remove(in.FullPath); err != nil {
		slog.Warn("删除失败", slog.Any("源文件", in.FullPath))
	} else {
		slog.Warn("删除成功", slog.Any("源文件", in.FullPath))
	}
}
