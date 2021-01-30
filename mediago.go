package main

import (
	"errors"
	"flag"
	"mediago/parser"

	"mediago/engine"
	"mediago/utils"
)

func main() {
	var (
		err      error
		name     string
		videoDir string
		m3u8Url  string
	)

	flag.StringVar(&name, "name", "新影片", "string类型参数")
	flag.StringVar(&videoDir, "path", "", "string类型参数")
	flag.StringVar(&m3u8Url, "url", "", "string类型参数")
	flag.Parse()

	if !utils.IsUrl(m3u8Url) {
		panic(errors.New("参数 url 格式错误"))
	}

	videoDir = utils.NormalizePath(videoDir)

	e := &engine.Engine{}

	decoder := &engine.M3u8Decoder{}

	params := engine.DownloadParams{
		Name:     name,
		Local:    videoDir,
		Url:      m3u8Url,
		Decoder:  decoder,
		Download: parser.DownloadSegment,
	}

	if err = e.Run(params); err != nil {
		utils.Logger.Error(err)
		panic(err)
	}
}
