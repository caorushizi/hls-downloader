package download

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"
)

type Downloader struct {
	sync.WaitGroup

	Segments []string
	Chs      chan int  // 默认下载量
	Ans      chan bool // 每个进程的下载状态
}

func New(count int) Downloader {
	return Downloader{
		Chs: make(chan int, count),
		Ans: make(chan bool),
	}
}

// 初始化下载地址  根据项目确认使用配置文件的方式还是其他方式，此处使用爬虫处理没公开
func (downloader *Downloader) InitSegments(filename string) (err error) {
	var fileHandle *os.File
	if fileHandle, err = os.Open(filename); err != nil {
		return
	}

	defer fileHandle.Close()
	fileScanner := bufio.NewScanner(fileHandle)

	for fileScanner.Scan() {
		text := fileScanner.Text()
		switch {
		case strings.HasPrefix(text, "#EXT"):
		case strings.HasPrefix(text, "#"):
			continue
		default:
			downloader.Segments = append(downloader.Segments, text)
		}
	}
	return
}

// 实际的下载操作
func downloadHandle(url string, path string) {
	//需要根据下载内容作存储等处理
	err := Start(path, url)
	if err != nil {
		log.Print("下载文件时出错：", err)
	}
}

/**
每个线程的操作
url 下载地址
chs 默认下载量
ans 每个线程的下载状态
*/
func (downloader *Downloader) Work(segmentName string, urlString string, baseFolder string) {
	defer func() {
		<-downloader.Chs // 某个任务下载完成，让出
		downloader.Done()
	}()

	fmt.Println("开始下载" + segmentName)
	var (
		u1  *url.URL
		err error
	)

	if u1, err = url.Parse(urlString); err != nil {
		panic("invalid url")
	}

	u1.Path = path.Join(u1.Path, segmentName)
	newFilePath := fmt.Sprintf(path.Join(baseFolder, segmentName))
	fullUrl := u1.String()

	downloadHandle(fullUrl, newFilePath)

	fmt.Println("下载完成" + segmentName)
	downloader.Ans <- true // 告知下载完成
}
