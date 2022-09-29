package io

import (
	"errors"
	"os"
)

//io辅助类

// CreateFolder 创建文件夹
// p:文件路径；ignoreExists:忽略已存在
func CreateFolder(p string, ignoreExists bool) error {
	//文件夹已存在，直接返回
	if FolderExists(p) == true && ignoreExists == false {
		return errors.New("folder exists")
	}
	//文件夹不存在，创建
	if FolderExists(p) == false {
		//0777表示：创建了一个普通文件，所有人拥有所有的读、写、执行权限
		err := os.MkdirAll(p, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

// FolderExists 判断filename文件夹是否存在
func FolderExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	if info == nil {
		return false
	}
	return info.IsDir()
}

// FileExists 判断filename文件是否存在
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	if info == nil {
		return false
	}
	return true
}
