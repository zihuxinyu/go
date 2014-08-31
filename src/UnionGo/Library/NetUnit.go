package Library

import (
	"strings"
	"os"
	// "path/filepath"
	"net/http"
	"io/ioutil"
)
// net的util
// toPath 文件保存的目录
// 默认是/tmp
// 返回文件的完整目录
func WriteUrl(url string, toPath string) (path string, ok bool) {
	if url == "" {
		return;
	}
	content, err := GetContent(url)
	if err != nil {
		return;
	}
	// a.html?a=a11&xxx
	url = trimQueryParams(url)
	_, ext := SplitFilename(url)
	if toPath == "" {
		toPath = "/tmp"
	}
	// dir := filepath.Dir(toPath)
	newFilename := NewGuid() + ext
	fullPath := toPath + "/" + newFilename
	/*
	if err := os.MkdirAll(dir, 0777); err != nil {
	return
	}
	*/
	// 写到文件中
	file, err := os.Create(fullPath)
	defer file.Close()
	if err != nil {
		return
	}
	file.Write(content)
	path = fullPath
	ok = true
	return
}
// 得到内容
func GetContent(url string) (content []byte, err error) {
	var resp *http.Response
	resp, err = http.Get(url)
	if(resp != nil && resp.Body != nil) {
		defer resp.Body.Close()
	} else {
	}
	if resp == nil || resp.Body == nil || err != nil || resp.StatusCode != http.StatusOK {
		return
	}
	var buf []byte
	buf, err = ioutil.ReadAll(resp.Body)
	if(err != nil) {
		return
	}
	content = buf;
	err = nil
	return
}
// 将url ?, #后面的字符串去掉
func trimQueryParams(url string) string {
	pos := strings.Index(url, "?");
	if pos != -1 {
		url = Substr(url, 0, pos);
	}
	pos = strings.Index(url, "#");
	if pos != -1 {
		url = Substr(url, 0, pos);
	}
	pos = strings.Index(url, "!");
	if pos != -1 {
		url = Substr(url, 0, pos);
	}
	return url;
}
