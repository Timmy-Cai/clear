package main

import (
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"
)

var paths = make([]string, 0)

func main() {
	readConfig()
	for _, filePath := range paths {
		go checkFileTime(filePath)
	}

	select {}
}

func readConfig() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("配置文件读取失败 %s \n", err))
	}

	paths = viper.GetStringSlice("path")

	fmt.Println("监控目录")
	for _, filePath := range paths {
		fmt.Println(filePath)
	}
}

func checkFileTime(absPath string) {
	for {
		fmt.Println(" 检查目录", absPath)
		fileList, _ := ioutil.ReadDir(absPath)
		for i := range fileList {
			if fileList[i].ModTime().Add(time.Hour).Before(time.Now()) {
				checkPath := filepath.FromSlash(path.Join(absPath, fileList[i].Name()))
				fmt.Println("删除文件", checkPath)
				err := removeContents(checkPath)
				if err != nil {
					fmt.Println("发生错误:", err)
					panic(err)
				}
			}
		}
		time.Sleep(time.Duration(5) * time.Second)
	}
}

func removeContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		//删除文件
		err = os.RemoveAll(dir)
		return err
	}
	if len(names) == 0 {
		_ = os.RemoveAll(dir)
		return nil
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}
