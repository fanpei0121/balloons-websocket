package service

import (
	"balloons/util"
	"strconv"
	"time"
)

// 签名

// 获取签名
func GetSign(secretKey string) (string, error) {
	micro := "000" + strconv.Itoa(int(time.Now().Unix()))
	en, err := util.Encrypt([]byte(micro), []byte(secretKey))
	if err != nil {
		return "", err
	}
	sign := util.Base64Encode(en)
	return sign, nil
}

// 检测签名是否正确
func CheckSign(sign string, secretKey string) (string, error) {
	aa := util.Base64Decode(sign)
	str, err := util.Dncrypt(string(aa), []byte(secretKey))
	if err != nil {
		return "", err
	}
	return str, nil
}
