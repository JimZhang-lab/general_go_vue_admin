/*
 * @Author: JimZhang
 * @Date: 2025-07-24 22:11:35
 * @LastEditors: 很拉风的James
 * @LastEditTime: 2025-07-25 11:07:47
 * @FilePath: /server/api/service/captcha.go
 * @Description:
 *
 */
package service

import (
	"image/color"
	"server/common/utils"

	"github.com/mojocn/base64Captcha"
)

var store = utils.RedisStore{}

// 生成验证码

func CaptMake() (id, b64s string) {
	var driver base64Captcha.Driver
	var driverString base64Captcha.DriverString

	// 配置验证码信息
	captchaConfig := base64Captcha.DriverString{
		Height:          60,
		Width:           200,
		NoiseCount:      0,
		ShowLineOptions: 2 | 4,
		Length:          6,
		Source:          "1234567890qwertyuioplkjhgfdsazxcvbnm",
		BgColor: &color.RGBA{
			R: 3,
			G: 102,
			B: 214,
			A: 125,
		},
		Fonts: []string{"wqy-microhei.ttc"},
	}

	driverString = captchaConfig
	driver = driverString.ConvertFonts()
	captcha := base64Captcha.NewCaptcha(driver, store)
	lid, lb64s, _, _ := captcha.Generate()
	return lid, lb64s
}

func CaptVerify(id, capt string) bool {
	// var storeInstance utils.RedisStore
	if store.Verify(id, capt, false) {
		return true
	} else {
		return false
	}
}
