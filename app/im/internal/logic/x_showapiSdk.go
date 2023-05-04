package logic

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"github.com/satori/go.uuid"
	"io"
	"mime"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type NormalReq struct {
	url        string
	textMap    url.Values
	uploadMap  map[string]interface{}
	timeout    time.Duration
	headMap    url.Values
	charset    string
	headString string
	bodyString string
	resHeadMap map[string]string
}

// ShowApiRequest 用于请求官网
func ShowApiRequest(reqUrl string, appid int, sign string) *NormalReq {
	values := make(url.Values)
	values.Set("showapi_appid", strconv.Itoa(appid))
	values.Set("showapi_sign", sign)
	return &NormalReq{reqUrl, values, nil, 10 * time.Second, values, "utf-8", "", "", make(map[string]string)}
}

func (request *NormalReq) AddTextPara(key, value string) {
	request.textMap.Set(key, value)
}

func (request *NormalReq) Url() string {
	reqUrl := request.url
	return reqUrl
}

func (request *NormalReq) Post() (string, error) {
	var dial = &net.Dialer{
		Timeout:       30 * time.Second,
		KeepAlive:     30 * time.Second,
		FallbackDelay: -1,
	}
	var ctx = func(ctx context.Context, network, addr string) (net.Conn, error) {
		network = "tcp4" //仅使用ipv4
		return dial.DialContext(ctx, network, addr)
	}
	tr := &http.Transport{
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           ctx,
		MaxIdleConns:          20,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	client := &http.Client{
		Timeout:   request.timeout,
		Transport: tr,
	}
	var resp *http.Response
	if request != nil && request.uploadMap != nil && len(request.uploadMap) > 0 {
		var fileName, filePath string
		for k, v := range request.uploadMap {
			fileName = k
			typeOfUpload := reflect.TypeOf(v)
			if typeOfUpload.Kind() == reflect.String {
				filePath = v.(string)
			} else if typeOfUpload.Kind() == reflect.Ptr {
				file := v.(*os.File)
				if file == nil {
					return "", errors.New("请确认上传的文件是否正确")
				}
				filePath = file.Name()
			} else if typeOfUpload.Kind() == reflect.Slice {
				fileBytes := v.([]byte)
				fileType := http.DetectContentType(fileBytes) //文件后缀
				u := uuid.NewV4()
				fileName := strings.Replace(u.String(), "-", "", -1)
				fileEndings, err := mime.ExtensionsByType(fileType)
				if err != nil {
					return "", err
				}
				newPath := filepath.Join("/NFS/hby/upload", fileName+fileEndings[0]) //生成完整的文件名
				newFile, err := os.Create(newPath)
				if err != nil {
					return "", err
				}
				if _, err := newFile.Write(fileBytes); err != nil {
					return "", err
				}
				defer os.Remove(newFile.Name())
				filePath = newFile.Name()
				newFile.Close()
			} else {
				return "", errors.New("上传格式为string或File或[]Byte")
			}
		}
		//只支持单文件上传
		file, err := os.Open(filePath)
		if err != nil {
			return "", err
		}
		defer file.Close()
		postbody := &bytes.Buffer{}
		writer := multipart.NewWriter(postbody)
		part, err := writer.CreateFormFile(fileName, filepath.Base(filePath))
		if err != nil {
			return "", err
		}
		_, err = io.Copy(part, file)
		for k, v := range request.textMap {
			writer.WriteField(k, v[0])
		}
		err = writer.Close()
		if err != nil {
			return "request err", err
		}
		req, err := http.NewRequest("POST", request.Url(), postbody)
		req.Header.Add("Content-Type", writer.FormDataContentType())
		for k, v := range request.headMap {
			if k != "Content-Type" {
				req.Header.Add(k, v[0])
			}
		}
		if request.headString != "" {
			strlist := strings.Split(request.headString, "\r\n")
			if len(strlist) <= 1 {
				strlist = strings.Split(request.headString, "\n")
			}
			for i := 0; i < len(strlist); i++ {
				setheadlist := strings.Split(strlist[i], ":")
				req.Header.Add(strings.TrimSpace(setheadlist[0]), strings.TrimSpace(setheadlist[1]))
			}
		}
		resp, err = client.Do(req)
		if err != nil {
			return "request err", err
		}

	} else {
		bodystr := strings.TrimSpace(request.textMap.Encode())
		if request.bodyString != "" {
			bodystr = bodystr + strings.TrimSpace(request.bodyString)
		}
		req, err := http.NewRequest("POST", strings.TrimSpace(request.url), strings.NewReader(bodystr))
		if request.headString != "" {
			strlist := strings.Split(request.headString, "\r\n")
			if len(strlist) <= 1 {
				strlist = strings.Split(request.headString, "\n")
			}
			for i := 0; i < len(strlist); i++ {
				setheadlist := strings.Split(strlist[i], ":")
				req.Header.Add(strings.TrimSpace(setheadlist[0]), strings.TrimSpace(setheadlist[1]))
			}
		}
		for k, v := range request.headMap {
			if k == "Content-Type" && !strings.Contains(v[0], "charset") {
				req.Header.Add(k, v[0]+";charset="+request.charset)
			} else {
				req.Header.Add(k, v[0])
			}
		}
		resp, err = client.Do(req)
		if err != nil {
			return "request err", err
		}
	}
	for k, v := range resp.Header {
		request.resHeadMap[k] = v[0]
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return "response err", err
	}
	return string(body), nil
}
