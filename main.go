package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/zhangyiming748/goini"
	"golang.org/x/exp/slog"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	conf   *goini.Config
	logger *slog.Logger
)

const (
	configPath = "./settings.ini"
)

func setLevel(level string) {
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
	file := "uploadFile.log"
	logf, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		panic(err)
	}
	logger = slog.New(opt.NewJSONHandler(io.MultiWriter(logf, os.Stdout)))
}
func main() {
	conf = goini.SetConfig(configPath)
	level, _ := conf.GetValue("log", "level")
	setLevel(level)
	key, _ := conf.GetValue("person", "key")
	if len(os.Args) < 2 {
		logger.Warn("第二个参数为文件名")
	}
	file := os.Args[1]
	res := upload(file, key)
	var m Media
	json.Unmarshal(res, &m)
	mediaId, err := os.OpenFile("media.id", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0777)
	if err != nil {
		return
	}
	defer mediaId.Close()
	mediaId.WriteString(m.MediaId)
	logger.Info("返回", slog.String("上传文件获得的media_id", m.MediaId))
}
func getRes() {

}

type Media struct {
	Errcode   int    `json:"errcode"`
	Errmsg    string `json:"errmsg"`
	Type      string `json:"type"`
	MediaId   string `json:"media_id"`
	CreatedAt string `json:"created_at"`
}

func upload(media, key string) []byte {
	url := strings.Join([]string{"https://qyapi.weixin.qq.com/cgi-bin/webhook/upload_media?type=file&key", key}, "=")
	//url := "https://qyapi.weixin.qq.com/cgi-bin/webhook/upload_media?key=&type=file"
	logger.Info("拼接后的网址", slog.String("拼接的url", url))
	method := "POST"
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	file, errFile1 := os.Open(media)
	defer file.Close()
	part1,
		errFile1 := writer.CreateFormFile("media", filepath.Base("<media>"))
	_, errFile1 = io.Copy(part1, file)
	if errFile1 != nil {
		fmt.Println(errFile1)
		return nil
	}
	err := writer.Close()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return nil
	}
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://www.apifox.cn)")
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	fmt.Println(string(body))
	return body
}
