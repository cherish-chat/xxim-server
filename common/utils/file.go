package utils

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

func DownloadFile(url string, saveTo string) error {
	// Create the file
	out, err := os.Create(saveTo)
	if err != nil {
		return err
	}

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New("download error, status: " + resp.Status)
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func FileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	return !os.IsNotExist(err)
}

type xFile struct {
}

var File = &xFile{}

func (x *xFile) GetContentTypeByFilename(filename string) string {
	ext := x.GetExt(filename)
	switch ext {
	case ".png":
		return "image/png"
	case ".jpg":
		return "image/jpeg"
	case ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".mp4":
		return "video/mp4"
	case ".webm":
		return "video/webm"
	case ".ogg":
		return "video/ogg"
	case ".mov":
		return "video/quicktime"
	case ".avi":
		return "video/x-msvideo"
	case ".wmv":
		return "video/x-ms-wmv"
	case ".m3u8":
		return "application/x-mpegURL"
	case ".ts":
		return "video/MP2T"
	case ".mp3":
		return "audio/mpeg"
	case ".wav":
		return "audio/wav"
	case ".flac":
		return "audio/flac"
	case ".pdf":
		return "application/pdf"
	case ".doc":
		return "application/msword"
	case ".docx":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case ".xls":
		return "application/vnd.ms-excel"
	case ".xlsx":
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	case ".ppt":
		return "application/vnd.ms-powerpoint"
	case ".pptx":
		return "application/vnd.openxmlformats-officedocument.presentationml.presentation"
	case ".zip":
		return "application/zip"
	case ".rar":
		return "application/x-rar-compressed"
	case ".7z":
		return "application/x-7z-compressed"
	case ".tar":
		return "application/x-tar"
	case ".gz":
		return "application/gzip"
	case ".bz2":
		return "application/x-bzip2"
	case ".xz":
		return "application/x-xz"
	case ".exe":
		return "application/x-msdownload"
	case ".swf":
		return "application/x-shockwave-flash"
	case ".rtf":
		return "application/rtf"
	case ".eot":
		return "application/vnd.ms-fontobject"
	case ".otf":
		return "font/otf"
	case ".ttf":
		return "font/ttf"
	case ".woff":
		return "font/woff"
	case ".woff2":
		return "font/woff2"
	case ".svg":
		return "image/svg+xml"
	case ".svgz":
		return "image/svg+xml"
	case ".ico":
		return "image/x-icon"
	case ".json":
		return "application/json"
	case ".xml":
		return "application/xml"
	case ".txt":
		return "text/plain"
	case ".md":
		return "text/markdown"
	case ".csv":
		return "text/csv"
	case ".html":
		return "text/html"
	case ".htm":
		return "text/html"
	case ".js":
		return "text/javascript"
	case ".css":
		return "text/css"
	default:
		return "application/octet-stream"
	}
}

func (x *xFile) GetExt(filename string) string {
	return path.Ext(filename)
}

func (x *xFile) FilenameFromUrl(uri string) string {
	// example: https://example.com/abc.jpg?x1=1&x2=2
	// return abc.jpg
	// 先解码
	tmp, _ := url.Parse(uri)
	uri = tmp.Path
	// 去掉参数
	uri = strings.Split(uri, "?")[0]
	return path.Base(uri)
}
