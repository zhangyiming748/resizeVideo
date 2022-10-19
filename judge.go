package resizeVideo

import (
	"github.com/zhangyiming748/resizeVideo/log"
	"os"
	"runtime"
	"strconv"
)

func illegal(src, dst, threads string) bool {
	if src == dst {
		log.Debug.Println("输入输出目录相同")
		return true
	}
	if !exists(src) {
		log.Debug.Println("src目录不存在")
		return true
	}
	if !exists(dst) {
		log.Debug.Println("dst目录不存在")
		return true
	}
	if !isDir(src) {
		log.Debug.Println("src不是目录")
		return true
	}
	if !isDir(dst) {
		log.Debug.Println("dst不是目录")
		return true
	}
	if !allowThreads(threads) {
		log.Debug.Println("不允许的线程数")
		return true
	}
	return false
}

// 判断所给路径文件/文件夹是否存在
func exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// 判断所给路径是否为文件夹
func isDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断给定线程数是否合法
func allowThreads(threads string) bool {
	maxThreads := runtime.NumCPU()
	if t, err := strconv.Atoi(threads); err != nil {
		return false
	} else if t >= maxThreads {
		return false
	} else {
		return true
	}
}
