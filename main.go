/*
10.63
./DelSameFile -f /Volumes/CF2TSSD1/PicTmp2 &> CF2TSSD1_PicTmp2.log
./DelSameFile -f /Volumes/CF2TSSD1/PicTmp1 &> CF2TSSD1_PicTmp1.log
./DelSameFile -f /Volumes/CF2TSSD1/VideoJp &> CF2TSSD1_VideoJp.log
./DelSameFile -f /Volumes/CF2TSSD1/VideoTmp &> CF2TSSD1_VideoTmp.log
local
./DelSameFile -f /Users/yuankang/Movies/pic &> local_pic.log
./DelSameFile -f /Users/yuankang/Movies/222 &> local_222.log 691.5
10.97
./DelSameFile -f /Volumes/WD1TSSD/VideoTmp &> WD1TSSD_VideoTmp.log
./DelSameFile -f /Volumes/WD1TSSD/PicTmp &> WD1TSSD_PicTmp.log
20.76
./DelSameFile -f /Volumes/WD1T/VideoTmp1 &> WD1T_VideoTmp1.log
./DelSameFile -f /Volumes/WD1T/VideoTmp2 &> WD1T_VideoTmp2.log
11.75
./DelSameFile -f /Volumes/WD2T/VideoBig &> WD2T_VideoBig.log
./DelSameFile -f /Volumes/WD2T/VideoSmall &> WD2T_VideoSmall.log
./DelSameFile -f /Volumes/WD2T/VideoTmp &> WD2T_VideoTmp.log
*/
package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
)

func IsImageFile(filePath string) bool {
	fe := filepath.Ext(filePath)
	if fe == ".jpg" || fe == ".png" || fe == ".jpeg" || fe == ".bmp" || fe == ".gif" || fe == ".webp" {
		return true
	}
	log.Printf("%s is not a image file", filePath)
	return false
}

func FileMd5(fn string) (string, error) {
	file, err := os.Open(fn)
	if err != nil {
		log.Println("Error opening file:", err)
		return "", err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Println("Error getting file size:", err)
		return "", err
	}
	fileSize := fileInfo.Size()

	// 定义读取位置
	var start, middle, end int64
	bufferSize := int64(100)

	start = 0
	middle = fileSize/2 - bufferSize/2
	if middle < 0 {
		middle = 0
	}
	end = fileSize - bufferSize
	if end < 0 {
		end = 0
	}

	// 读取前100字节
	startBuffer := make([]byte, bufferSize)
	_, err = file.ReadAt(startBuffer, start)
	if err != nil && err != io.EOF {
		log.Println("Error reading start buffer:", err)
		return "", err
	}

	// 读取中间100字节
	middleBuffer := make([]byte, bufferSize)
	_, err = file.ReadAt(middleBuffer, middle)
	if err != nil && err != io.EOF {
		log.Println("Error reading middle buffer:", err)
		return "", err
	}

	// 读取后100字节
	endBuffer := make([]byte, bufferSize)
	_, err = file.ReadAt(endBuffer, end)
	if err != nil && err != io.EOF {
		log.Println("Error reading end buffer:", err)
		return "", err
	}

	// 连接所有字节
	data := append(startBuffer, middleBuffer...)
	data = append(data, endBuffer...)

	// 计算MD5值
	md5hash := md5.Sum(data)
	md5str := hex.EncodeToString(md5hash[:])
	return md5str, nil
}

// TraverseFolder 遍历指定文件夹并返回每个文件的全路径
func TraverseFolder(folder string) ([]string, error) {
	var files []string
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println(err)
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

//1 启动时加载sqlite数据库中的所有数据
//2 http下发任务, 指定检查路径
//3 获取路径下的所有文件的全路径
//4 循环计算每个文件的md5值
//5 如果md5值不存在, 就入库
//6 如果md5值存在, 就检查已存在md5值对应的文件是否存在
//7 如果文件存在 就删除新md5值对应的文件
//8 如果文件不存在, 就删除老md5的记录, 并入库新md5值的记录

//go run main.go -folder /path/to/your-folder
//./DelSameFile -folder ./pic/
func main() {
	//log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)

	// 定义命令行参数
	folder := flag.String("folder", ".", "folder to traverse")
	flag.Parse()

	files, err := TraverseFolder(*folder)
	if err != nil {
		log.Fatalf("Error traversing folder: %v", err)
	}

	for _, file := range files {
		md5hash, err := FileMd5(file)
		if err != nil {
			log.Println(err)
			continue
		}
		log.Printf("md5:%s, file:%s", md5hash, file)
	}
}
