package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/bogem/id3v2"
)

type Song struct {
	title  string
	artist string
	album  string
}

var argIsHelp = flag.Bool("h", false, "show help")
var argOutputBasePath = flag.String("o", "./out", "decrypt mp3 file output path")
var argCachePath = flag.String("s", "", "music cache path")

func argParse() {
	flag.Parse()

	if *argIsHelp || *argCachePath == "" {
		flag.Usage()
		os.Exit(0)
	}
}

func main() {
	// TODO 更健壮（文件重复、输入路径无效...）
	// TODO 封面图片写入 tag
	// TODO 转换和获取名称使用 chan

	argParse()

	// cachePath := "./data"
	// outputBasePath := "./data/music"

	cachePath := *argCachePath
	outputBasePath := *argOutputBasePath

	cacheList := getCacheList(cachePath)
	for _, cacheName := range cacheList {
		// 网易云音乐：歌曲 id - 码率 - mp3 文件的 md5 摘要
		id := strings.Split(cacheName, "-")[0]
		song := getSongInfoByID(id)
		fmt.Println(song)

		srcFilePath := path.Join(cachePath, cacheName)
		dstFilePath := getOutputFilePath(outputBasePath, song)
		cacheTrans(srcFilePath, dstFilePath)
		fillTag(dstFilePath, song)
	}
}

func fillTag(dstFilePath string, song Song) {
	tag, err := id3v2.Open(dstFilePath, id3v2.Options{Parse: true})
	if err != nil {
		log.Fatal("Error while opening mp3 file: ", err)
	}
	defer tag.Close()

	tag.SetDefaultEncoding(id3v2.EncodingUTF8)

	tag.SetTitle(song.title)
	tag.SetArtist(song.artist)
	tag.SetAlbum(song.album)

	if err = tag.Save(); err != nil {
		log.Fatal("Error while saving a tag: ", err)
	}
}

func getOutputFilePath(outputBasePath string, song Song) string {
	filename := song.artist + " - " + song.title + ".mp3"
	filePath := path.Join(outputBasePath, song.artist)

	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		os.MkdirAll(filePath, os.ModePerm)
	}

	return path.Join(filePath, filename)
}

func getCacheList(cachePath string) []string {
	files, _ := ioutil.ReadDir(cachePath)
	fileList := make([]string, len(files))
	index := 0
	for _, f := range files {
		if !f.IsDir() && path.Ext(f.Name()) == ".uc" {
			fileList[index] = f.Name()
			index++
		}
	}
	return fileList[0:index]
}

func getSongInfoByID(id string) Song {
	resp, err := http.Get("https://api.imjad.cn/cloudmusic/?type=detail&id=" + id)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	ret := make(map[string]interface{})
	json.Unmarshal([]byte(data), &ret)
	song := ret["songs"].([]interface{})[0]
	// songs[0].name
	songName := song.(map[string]interface{})["name"]
	// songs[0].ar[0].name
	singerName := song.(map[string]interface{})["ar"].([]interface{})[0].(map[string]interface{})["name"]
	// songs[0].al.name
	albumName := song.(map[string]interface{})["al"].(map[string]interface{})["name"]

	return Song{title: songName.(string), artist: singerName.(string), album: albumName.(string)}
}

func cacheTrans(srcFilePath string, dstFilePath string) {
	f, err := os.Open(srcFilePath)
	if err != nil {
		panic(err.Error())
	}
	defer f.Close()

	buffer, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err.Error())
	}

	decrypt(buffer)

	err = ioutil.WriteFile(dstFilePath, buffer, 0644)
	if err != nil {
		panic(err.Error())
	}
}

func decrypt(buffer []byte) {
	// 网易云音乐，每个字节亦或 0xA3 即可
	for i := range buffer {
		buffer[i] = buffer[i] ^ 0xA3
	}
}
