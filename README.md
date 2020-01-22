# MusicDownload

使用从音乐软件的缓存中还原出 mp3 文件，当前仅支持网易云音乐。

## 安装和使用

1. 如果 `go get github.com/bogem/id3v2` 安装不上，可以手动使用 golang 在 GitHub 上的镜像仓库安装，先确保 `$GOPATH/src/golang.org/x` 目录存在，之后再该目录下使用命令 `git clone https://github.com/golang/text.git`
2. `go build` 之后使用编译输出文件即可，或者直接 `go run main.go -s YOUR_PATH -o YOUR_PATH`

## 特性

针对于网易云音乐：

- 回填 cache 音乐里没有 ID3 Tag 信息（标题、艺术家、专辑）
  - [ ] 回填专辑封面图片
- 根据艺术家分文件夹存储
