package resizeVideo

import (
	"github.com/zhangyiming748/resizeVideo/log"
	"github.com/zhangyiming748/voiceAlert"
	"os/exec"

	"os"
	"strings"
)

const (
	success  = iota + 1 // 单次转码成功
	failed              // 转码失败,程序退出
	complete            // 转码进程完成
)

func ResizeVideo(src, dst, pattern, threads string, isDelete bool) {
	if illegal(src, dst, threads) {
		os.Exit(1)
	}
	files := getFiles(src, pattern)
	log.Info.Println("文件目录", files)
	total := len(files)
	for index, file := range files {
		resize(src, dst, file, threads, index, total, isDelete)
	}
	voiceAlert.Voice(complete)
}
func resize(src, dst, file, threads string, index, total int, isDelete bool) {
	defer func() {
		if err := recover(); err != nil {
			voiceAlert.Voice(failed)
		} else {
			voiceAlert.Voice(success)
		}
	}()
	in := strings.Join([]string{src, file}, "/")
	log.Info.Printf("开始处理文件:%v\n", in)
	out := strings.Join([]string{dst, file}, "/")
	log.Debug.Println("源文件目录:", src)
	log.Debug.Println("输出文件目录:", dst)
	log.Debug.Println("开始处理文件:", in)
	log.Debug.Println("输出文件:", out)
	//ffmpeg -i 1.mp4 -strict -2 -vf scale=-1:1080 4.mp4
	cmd := exec.Command("ffmpeg", "-threads", threads, "-i", in, "-strict", "-2", "-vf", "scale=-1:1080", "-threads", threads, out)
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
		log.Info.Printf("正在处理第 %d/%d 个文件: %s\n", index+1, total, file)
		t := string(tmp)
		t = strings.Replace(t, "\u0000", "", -1)
		log.Info.Println(t)
		if err != nil {
			break
		}
	}
	if err = cmd.Wait(); err != nil {
		log.Debug.Panicf("命令执行中有错误产生:%v\n", err)
	}
	log.Debug.Printf("完成当前文件的处理:源文件是%s\t目标文件是%s\n", in, file)
	if isDelete {
		if err := os.RemoveAll(in); err != nil {
			log.Debug.Printf("删除源文件失败:%v\n", err)
		} else {
			log.Debug.Printf("删除源文件:%s\n", in)
		}
	}

}

func getFiles(dir, pattern string) []string {
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Debug.Printf("读取文件目录产生的错误:%v\n", err)
	}
	var aim []string
	types := strings.Split(pattern, ";") //"wmv;rm"
	for _, f := range files {
		if l := strings.Split(f.Name(), ".")[0]; len(l) != 0 {
			for _, v := range types {
				if strings.HasSuffix(f.Name(), v) {
					log.Debug.Printf("有效的目标文件:%v\n", f.Name())
					aim = append(aim, f.Name())
				}
			}
		}
	}
	return aim
}
