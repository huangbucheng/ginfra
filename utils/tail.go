package utils

import (
	"bytes"
	"os"
)

const (
	defaultBufSize = 4096
)

//Tail 获取文件的最后N行
func Tail(filename string, n int) (lines []string, err error) {
	var f *os.File
	fInfo, err := os.Stat(filename)
	if err != nil {
		return lines, err
	}
	size := fInfo.Size()

	f, err = os.Open(filename)
	if err != nil {
		return lines, err
	}
	defer f.Close()

	// b - read bytes
	b := make([]byte, defaultBufSize)
	// bBuf - buffer for incompeted line content between two read
	bBuf := bytes.NewBuffer([]byte{})
	startpos := size
	toReadLines := n
	flag := true
	for flag {
		if startpos < defaultBufSize {
			startpos = 0
		} else {
			startpos -= defaultBufSize
		}
		_, err = f.Seek(startpos, os.SEEK_SET)
		if err != nil {
			return lines, err
		}
		rsize, err := f.Read(b)
		if err != nil {
			return lines, err
		}
		if rsize <= 0 {
			return lines, nil
		}
		j := rsize
		for i := rsize - 1; i >= 0; i-- {
			if b[i] == '\n' {
				bLine := bytes.NewBuffer([]byte{})
				bLine.Write(b[i+1 : j])
				j = i
				if bBuf.Len() > 0 {
					bLine.Write(bBuf.Bytes())
					bBuf.Reset()
				}
				if bLine.Len() > 0 && toReadLines > 0 { //skip last "\n"
					lines = append(lines, bLine.String())
					toReadLines--
				}
				if toReadLines <= 0 {
					flag = false
					break
				}
			}
		}
		// incompeted line content
		if flag && j > 0 {
			if startpos == 0 {
				bLine := bytes.NewBuffer([]byte{})
				bLine.Write(b[:j])
				if bBuf.Len() > 0 {
					bLine.Write(bBuf.Bytes())
					bBuf.Reset()
				}
				lines = append(lines, bLine.String())
				flag = false
			} else {
				bb := make([]byte, bBuf.Len())
				copy(bb, bBuf.Bytes())
				bBuf.Reset()
				bBuf.Write(b[:j])
				bBuf.Write(bb)
			}
		}
	}
	return lines, nil
}
