/*
 * @Author: JimZhang
 * @Date: 2025-07-25 00:13:19
 * @LastEditors: 很拉风的James
 * @LastEditTime: 2025-07-25 00:13:30
 * @FilePath: /server/common/utils/encryption.go
 * @Description:
 *
 */
package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func EncryptionMd5(s string) string {
	ctx := md5.New()
	ctx.Write([]byte(s))
	return hex.EncodeToString(ctx.Sum(nil))
}
