package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	wg = sync.WaitGroup{}
)

func main() {
	var path string

	if err := loadConfig(); nil != err {
		log.Println(err)
		fmt.Scanln(&path)
		return
	}

	log.Printf("请把需要导出的目录或文件拖进来，回车导出全部配置\n")
	fmt.Scanln(&path)

	if len(path) == 0 {
		path = "./"
	}

	f, err := os.Stat(path)
	if nil != err {
		log.Println(err)
		fmt.Scanln(&path)
		return
	}
	if f.IsDir() {
		filepath.Walk(path, walkPath)
	} else {
		wg.Add(1)
		go export(path)
	}

	wg.Wait()

	log.Println("按任意键退出")
	fmt.Scanln(&path)
}

func walkPath(fileName string, info fs.FileInfo, err error) error {
	log.Println(fileName)
	if nil == info {
		log.Println(err)
		return err
	}
	if info.IsDir() {
		return nil
	}
	ext := filepath.Ext(info.Name())
	if ".xlsx" == ext || "xls" == ext {
		wg.Add(1)
		go export(fileName)
	}
	return nil
}

func export(fName string) {
	start := time.Now()
	exportExcel(fName)
	if cost := time.Since(start).Seconds(); cost >= 0 {
		log.Printf("导出配置表[%s]耗时%v秒\n", fName, cost)
	}
	wg.Done()
}
