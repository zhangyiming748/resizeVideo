package resizeVideo

import (
	"github.com/zhangyiming748/GetAllFolder"
	"github.com/zhangyiming748/GetFileInfo"
	"github.com/zhangyiming748/log"
	"github.com/zhangyiming748/replace"
	"github.com/zhangyiming748/voiceAlert"
	"io/fs"
	"os"
	"os/exec"
	"strings"
)

func ResizeAllVideos(root, pattern, threads string) {
	ResizeVideos(root, pattern, threads)
	folders := GetAllFolder.ListFolders(root)
	for i, folder := range folders {
		log.Debug.Printf("正在处理第 %d/%d 个文件夹\n", i+1, len(folders))
		ResizeVideos(folder, pattern, threads)
	}
}
func ResizeVideos(src, pattern, threads string) {
	files := GetFileInfo.GetAllVideoFileInfo(src, pattern)
	for _, file := range files {
		//ResizeVideo(file,threads,true)
		//pretty.P(file)
		if file.Width <= 1920 || file.Height <= 1920 {
			log.Debug.Printf("跳过处理尺寸正常的视频:%v\n", file.FullName)
			continue
		} else {
			log.Debug.Printf("视频信息:%+v\n", file)
		}
		if file.Width > file.Height {
			log.Info.Printf("准备处理横屏视频:%+v\n", file)
			ResizeVideo(file, threads, "1920x1080")
		} else if file.Width < file.Height {
			log.Info.Printf("准备处理竖屏视频:%+v\n", file)
			ResizeVideo(file, threads, "1080x1920")
		} else {
			log.Info.Printf("准备处理正方形视频:%+v\n", file)
			ResizeVideo(file, threads, "1920x1920")
		}
	}
}
func ResizeVideo(in GetFileInfo.Info, threads string, p string) {
	defer func() {
		if err := recover(); err != nil {
			//pretty.P(err)
			log.Warn.Printf("出现错误的文件:%s\n", in.FullPath)
			voiceAlert.Customize("failed", voiceAlert.Samantha)
		}
	}()
	dst := strings.Trim(in.FullPath, in.FullName)   //原始目录
	dst = strings.Join([]string{dst, "resize"}, "") //二级目录
	fname := strings.Trim(in.FullName, in.ExtName)  //仅文件名
	mp4 := strings.Join([]string{fname, "mp4"}, ".")
	os.Mkdir(dst, fs.ModePerm)
	log.Info.Printf("开始处理文件:%v\n", in)
	out := strings.Join([]string{dst, mp4}, string(os.PathSeparator))
	log.Debug.Println("源文件:", in.FullPath)
	log.Debug.Println("输出文件:", out)
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
		log.Warn.Printf("出问题的视频:%+v\n", in)
	}
	log.Debug.Printf("生成的命令是:%s\n", cmd)
	stdout, err := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout
	if err != nil {
		log.Warn.Panicf("cmd.StdoutPipe产生的错误:%v\n", err)
	}
	if err = cmd.Start(); err != nil {
		log.Warn.Panicf("cmd.Run产生的错误:%v\n", err)
	}
	for {
		tmp := make([]byte, 1024)
		_, err := stdout.Read(tmp)
		//写成输出日志
		t := string(tmp)
		t = replace.Replace(t)
		log.TTY.Println(t)
		if err != nil {
			break
		}
	}
	if err = cmd.Wait(); err != nil {
		log.Warn.Panicf("命令执行中有错误产生:%v\n", err)
	}
	if err := os.Remove(in.FullPath); err != nil {
		log.Warn.Printf("删除源文件失败:%v\n", err)
	} else {
		log.Debug.Printf("删除源文件:%s\n", in.FullPath)
	}
}
