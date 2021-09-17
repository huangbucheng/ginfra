package utils

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
)

//IsDir 判断路径是否目录
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

//Exists 判断路径是否存在
func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

// CopyFile copies a single file from src to dst
func CopyFile(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}

// CopyDir copies a whole directory recursively
func CopyDir(src string, dst string) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = CopyDir(srcfp, dstfp); err != nil {
				return err
			}
		} else {
			if err = CopyFile(srcfp, dstfp); err != nil {
				return err
			}
		}
	}
	return nil
}

//WriteFile write content to file
func WriteFile(filename, content string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(content)
	if err != nil {
		return err
	}

	err = f.Sync()
	if err != nil {
		return err
	}

	return nil
}

//ReadNLines 读取文件的前N行
func ReadNLines(path string, n int) ([]string, error) {
	var lines []string
	var count int
	inFile, err := os.Open(path)
	if err != nil {
		return lines, err
	}
	defer inFile.Close()

	scanner := bufio.NewScanner(inFile)
	for scanner.Scan() {
		lines = append(lines, scanner.Text()) // the line
		count += 1
		if count >= n {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return lines, err
	}
	return lines, nil
}

//ReadNLinesFromString 读取字符串内容的前N行
func ReadNLinesFromString(content string, n int) ([]string, error) {
	var lines []string
	var count int
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		lines = append(lines, scanner.Text()) // the line
		count += 1
		if count >= n {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return lines, err
	}
	return lines, nil
}

//ReadFile 读取文件内容，支持COS、HTTP、HTTPS、本地文件
func ReadFile(logfile string) (string, error) {
	// try cos
	if strings.HasPrefix(logfile, CosBucketUrl) {
		if len(CosSecretKey) > 0 {
			body, err := GetFileFromCos(logfile)
			if err != nil {
				return "", err
			}
			return string(body), nil
		}
	}

	// http url
	if strings.HasPrefix(logfile, "http://") || strings.HasPrefix(logfile, "https://") {
		client := http.Client{
			CheckRedirect: func(r *http.Request, via []*http.Request) error {
				r.URL.Opaque = r.URL.Path
				return nil
			},
		}
		// Put content on file
		resp, err := client.Get(logfile)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		return string(body), nil
	}

	// local file
	if Exists(logfile) {
		body, _ := ioutil.ReadFile(logfile)
		return string(body), nil
	}

	return "", fmt.Errorf("cannot handle file:%s", logfile)
}

//AppendToFile write to end of file
func AppendToFile(filename, content string) error {
	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	if _, err := f.Write([]byte(content + "\n")); err != nil {
		return err
	}
	if err := f.Sync(); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return nil
}
