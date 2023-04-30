package resizeVideo

import (
	"fmt"
	"github.com/zhangyiming748/GetAllFolder"
	"github.com/zhangyiming748/GetFileInfo"
	"github.com/zhangyiming748/voiceAlert"
	"golang.org/x/exp/slog"
	"io"
	"os"
	"os/exec"
	"strings"
)

var mylog *slog.Logger

func setLog(level string) {
	var opt slog.HandlerOptions
	switch level {
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
	file := "resizeVideo.log"
	logf, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		panic(err)
	}
	//defer logf.Close() //如果不关闭可能造成内存泄露
	mylog = slog.New(opt.NewJSONHandler(io.MultiWriter(logf, os.Stdout)))
}
func init() {
	l := os.Getenv("LEVEL")
	setLog(l)
}
func ResizeAllVideos(root, pattern, threads string) {
	ResizeVideos(root, pattern, threads)
	folders := GetAllFolder.List(root)
	for i, folder := range folders {
		mylog.Info(fmt.Sprintf("正在处理第 %d/%d 个文件夹", i+1, len(folders)))
		ResizeVideos(folder, pattern, threads)
	}
}
func ResizeVideos(src, pattern, threads string) {
	files := GetFileInfo.GetAllVideoFileInfo(src, pattern)
	for _, file := range files {
		if file.Width <= 1920 || file.Height <= 1920 {
			mylog.Info("跳过", slog.String("正常尺寸的视频", file.FullPath))
			continue
		}
		if file.Width > file.Height {
			mylog.Info("横屏视频", slog.Any("视频信息", file))
			ResizeVideo(file, threads, "1920x1080")
		} else if file.Width < file.Height {
			mylog.Info("竖屏视频", slog.Any("视频信息", file))
			ResizeVideo(file, threads, "1080x1920")
		} else {
			mylog.Info("正方形视频", slog.Any("视频信息", file))
			ResizeVideo(file, threads, "1920x1920")
		}
	}
}
func ResizeVideo(in GetFileInfo.Info, threads string, p string) {
	defer func() {
		if err := recover(); err != nil {
			mylog.Warn("错误", slog.String("文件信息", in.FullPath))
			voiceAlert.Customize("failed", voiceAlert.Samantha)
		}
	}()
	dst := strings.Trim(in.FullPath, in.FullName)   //原始目录
	dst = strings.Join([]string{dst, "resize"}, "") //二级目录
	fname := strings.Trim(in.FullName, in.ExtName)  //仅文件名
	mp4 := strings.Join([]string{fname, "mp4"}, ".")
	os.Mkdir(dst, 0777)
	mylog.Debug("新建文件夹", slog.String("全名", dst))
	out := strings.Join([]string{dst, mp4}, string(os.PathSeparator))
	mylog.Info("io", slog.String("源文件:", in.FullPath), slog.String("输出文件:", out))
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
		mylog.Warn("不正常的视频源", slog.Any("视频信息", in.FullPath))
	}
	mylog.Info("ffmpeg", slog.String("生成的命令", fmt.Sprintf("生成的命令是:%s\n", cmd)))
	stdout, err := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout
	if err != nil {
		mylog.Warn("cmd.StdoutPipe", slog.Any("错误", err))
		return
	}
	if err = cmd.Start(); err != nil {
		mylog.Warn("cmd.Run", slog.Any("错误", err))
		return
	}
	for {
		tmp := make([]byte, 1024)
		_, err = stdout.Read(tmp)
		//写成输出日志
		t := string(tmp)
		t = strings.Replace(t, "\u0000", "", -1)
		fmt.Println(t)
		if err != nil {
			break
		}
	}
	if err = cmd.Wait(); err != nil {
		mylog.Warn("命令执行中", slog.Any("错误", err))
		return
	}
	if err = os.Remove(in.FullPath); err != nil {
		mylog.Warn("删除失败", slog.String("源文件", in.FullPath), slog.Any("错误文本", err))
	} else {
		mylog.Warn("删除成功", slog.String("源文件", in.FullPath))
	}
}

func FixAll4x3s(root, pattern, threads string) {
	folders := GetAllFolder.List(root)
	folders = append(folders, root)
	for i, folder := range folders {
		mylog.Info(fmt.Sprintf("正在处理第 %d/%d 个文件夹", i+1, len(folders)))
		Fix4x3s(folder, pattern, threads)
	}
}

func Fix4x3s(src, pattern, threads string) {
	files := GetFileInfo.GetAllVideoFileInfo(src, pattern)
	for _, file := range files {
		//if file.Width <= 1920 || file.Height <= 1920 {
		//	mylog.Info("跳过", slog.String("正常尺寸的视频", file.FullPath))
		//	continue
		//}
		if file.Width > file.Height {
			mylog.Info("横屏视频", slog.Any("视频信息", file))
			Fix4x3(file, threads)
		} else {
			mylog.Warn("视频不正确", slog.Any("视频信息", file))
		}
	}
}

func Fix4x3(in GetFileInfo.Info, threads string) {
	defer func() {
		if err := recover(); err != nil {
			mylog.Warn("错误", slog.String("文件信息", in.FullPath))
			voiceAlert.Customize("failed", voiceAlert.Samantha)
		}
	}()
	dst := strings.Trim(in.FullPath, in.FullName)       //原始目录
	dst = strings.Join([]string{dst, "resolution"}, "") //二级目录
	fname := strings.Trim(in.FullName, in.ExtName)      //仅文件名
	mp4 := strings.Join([]string{fname, "mp4"}, ".")
	os.Mkdir(dst, 0777)
	mylog.Debug("新建文件夹", slog.String("全名", dst))
	out := strings.Join([]string{dst, mp4}, string(os.PathSeparator))
	mylog.Info("io", slog.String("源文件:", in.FullPath), slog.String("输出文件:", out))
	//ffmpeg -i 1.mp4 -strict -2 -vf scale=-1:1080 4.mp4
	// ffmpeg -threads 2 -i 4k_Saeko_Limo.mp4 -strict -2 -vf scale=-1:1080 -c:v libx265 -threads 2 1080.mp4
	var cmd *exec.Cmd = exec.Command("ffmpeg", "-i", in.FullPath, "-aspect", "4:3", "-c:v", "libx265", "-tag:v", "hvc1", "-threads", threads, out)
	mylog.Info("ffmpeg", slog.String("生成的命令", fmt.Sprintf("生成的命令是:%s\n", cmd)))
	defer func() {
		mylog.Warn(fmt.Sprintf("本次错误生成的命令:%v", cmd))
	}()
	stdout, err := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout
	if err != nil {
		mylog.Warn("cmd.StdoutPipe", slog.Any("错误", err))
		return
	}
	if err = cmd.Start(); err != nil {
		mylog.Warn("cmd.Run", slog.Any("错误", err))
		return
	}
	for {
		tmp := make([]byte, 1024)
		_, err = stdout.Read(tmp)
		t := string(tmp)
		t = strings.Replace(t, "\u0000", "", -1)
		fmt.Println(t)
		if err != nil {
			break
		}
	}
	if err = cmd.Wait(); err != nil {
		mylog.Warn("命令执行中", slog.Any("错误", err))
		return
	}
	if err = os.Remove(in.FullPath); err != nil {
		mylog.Warn("删除失败", slog.String("源文件", in.FullPath), slog.Any("错误文本", err))
	} else {
		mylog.Warn("删除成功", slog.String("源文件", in.FullPath))
	}
}
