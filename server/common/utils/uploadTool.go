/*
 * @Author: JimZhang
 * @Date: 2023-07-22 18:21:10
 * @LastEditors: 很拉风的James
 * @LastEditTime: 2025-07-25 20:33:25
 * @FilePath: /server/common/utils/uploadTool.go
 * @Description:
 *
 */

package utils

import "os"

// CreateDir 创建目录
func CreateDir(filePath string) error {
	if !IsExist(filePath) {
		err := os.MkdirAll(filePath, os.ModePerm)
		return err
	}
	return nil
}

// IsExist 判断是否存在
func IsExist(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		return os.IsExist(err)
	}
	return true
}
