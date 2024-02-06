package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func main() {
	//log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
	dir := "pic/"
	//dir := "/Users/yuankang/Pictures/wallpaper"

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Println("无法读取目录:", err)
		return
	}

	//seen := make(map[string]bool)

	for _, file := range files {
		filePath := filepath.Join(dir, file.Name())

		if isImageFile(filePath) {
			//log.Println(filePath)

			bs, err := GetPartialData(filePath)
			if err != nil {
				log.Println("无法读取文件:", err)
				continue
			}
			//log.Printf("bs:%x", bs)

			md5 := calculateMD5(bs)
			log.Printf("fn:%s, md5:%x", file.Name(), md5)

			/*
				if !seen[hash] {
					seen[hash] = true
					log.Println("全新图片:", filePath)
				} else {
					log.Println("重复图片:", filePath)
				}
			*/
		}
	}
}

func isImageFile(filePath string) bool {
	fe := filepath.Ext(filePath)
	if fe == ".jpg" || fe == ".png" || fe == ".jpeg" || fe == ".bmp" || fe == ".gif" || fe == ".webp" {
		return true
	}
	log.Printf("%s is not a image file", filePath)
	return false
}

func calculateMD5(data []byte) string {
	hasher := md5.New()
	hasher.Write(data)

	hashInBytes := hasher.Sum(nil)
	md5String := hex.EncodeToString(hashInBytes)
	return md5String
}

func GetPartialData(filePath string) ([]byte, error) {
	f100b, err := getFirstNBytes(filePath, 10)
	if err != nil {
		log.Println("Error reading first 100 bytes:", err)
		return nil, err
	}
	//log.Printf("F100:%x", f100b)

	l100b, err := getLastNBytes(filePath, 10)
	if err != nil {
		log.Println("Error reading last 100 bytes:", err)
		return nil, err
	}
	//log.Printf("L100:%x", l100b)

	bs := append(f100b, l100b...)
	return bs, nil
}

//首部50字节+中间50字节+尾部50字节+文件长度
func GetPartialData0(filePath string) ([]byte, error) {
	f100b, err := getFirstNBytes(filePath, 10)
	if err != nil {
		log.Println("Error reading first 100 bytes:", err)
		return nil, err
	}
	//log.Printf("F100:%x", f100b)

	l100b, err := getLastNBytes(filePath, 10)
	if err != nil {
		log.Println("Error reading last 100 bytes:", err)
		return nil, err
	}
	//log.Printf("L100:%x", l100b)

	bs := append(f100b, l100b...)
	return bs, nil
}

func getFirstNBytes(filePath string, n int) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buffer := make([]byte, n)
	_, err = io.ReadFull(file, buffer)
	if err != nil {
		return nil, err
	}

	return buffer, nil
}

func getMiddleNBytes(filePath string, n int) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 获取文件信息
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	// 确保文件大小至少为 n 字节
	if fileInfo.Size() < int64(n) {
		return nil, fmt.Errorf("文件太小，无法获取中间 %d 字节", n)
	}

	// 计算中间位置
	middleOffset := fileInfo.Size()/2 - int64(n/2)

	// 移动文件指针到中间位置
	_, err = file.Seek(middleOffset, io.SeekStart)
	if err != nil {
		return nil, err
	}

	// 读取中间 n 字节的数据
	middleBytes := make([]byte, n)
	_, err = file.Read(middleBytes)
	if err != nil {
		return nil, err
	}

	return middleBytes, nil
}

func getLastNBytes(filePath string, n int) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	fileSize := fileInfo.Size()
	if fileSize < int64(n) {
		n = int(fileSize)
	}

	buffer := make([]byte, n)
	_, err = file.ReadAt(buffer, fileSize-int64(n))
	if err != nil {
		return nil, err
	}

	return buffer, nil
}

func getFileLength(filePath string) (int64, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	// 获取文件信息
	fileInfo, err := file.Stat()
	if err != nil {
		return 0, err
	}

	// 获取文件长度
	fileLength := fileInfo.Size()

	return fileLength, nil
}
