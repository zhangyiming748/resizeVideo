package resizeVideo

import (
	"github.com/zhangyiming748/GetFileInfo"
	"github.com/zhangyiming748/log"
	"github.com/zhangyiming748/replace"
	"github.com/zhangyiming748/voiceAlert"
	"io/fs"
	"os"
	"os/exec"
	"strings"
)

const (
	UHD = "scale=-1:2160"
	QHD = "scale=-1:1440"
	FHD = "scale=-1:1080"
	HD  = "scale=-1:720"
)

func ResizeVideo(src, pattern, threads string, isDelete bool) {
	files := GetFileInfo.GetAllFileInfo(src, pattern)

	for _, file := range files {
		resize(file, threads, isDelete)
		voiceAlert.CustomizedOnMac(voiceAlert.Shanshan, "单个文件处理完成")

	}
	voiceAlert.CustomizedOnMac(voiceAlert.Shanshan, "全部文件处理完成")
}
func resize(in GetFileInfo.Info, threads string, isDelete bool) {
	defer func() {
		if err := recover(); err != nil {
			voiceAlert.CustomizedOnMac(voiceAlert.Shanshan, "文件处理发生错误")
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
	cmd := exec.Command("ffmpeg", "-threads", threads, "-i", in.FullPath, "-strict", "-2", "-vf", FHD, "-threads", threads, out)
	log.Debug.Printf("生成的命令是:%s\n", cmd)
	stdout, err := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout
	if err != nil {
		log.Debug.Panicf("cmd.StdoutPipe产生的错误:%v\n", err)
	}
	if err = cmd.Start(); err != nil {
		log.Debug.Panicf("cmd.Run产生的错误:%v\n", err)
	}
	for {
		tmp := make([]byte, 1024)
		_, err := stdout.Read(tmp)
		//写成输出日志
		t := string(tmp)
		t = replace.Replace(t)
		log.Info.Println(t)
		if err != nil {
			break
		}
	}
	if err = cmd.Wait(); err != nil {
		log.Debug.Panicf("命令执行中有错误产生:%v\n", err)
	}
	if isDelete {
		if err := os.Remove(in.FullPath); err != nil {
			log.Debug.Printf("删除源文件失败:%v\n", err)
		} else {
			log.Debug.Printf("删除源文件:%s\n", in.FullPath)
		}
	}

}
