package vivo

import (
	"aliens/common/character"
	"aliens/log"
	"crypto/md5"
	"encoding/hex"
	"gok/module/passport/vivo/model"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"time"
)

var vivoConfig model.Config

func Init(config model.Config) {
	vivoConfig = config
}

func GetVivoAcessToken(code string) (string, error) {
	paramMap := make(map[string]string)
	paramMap["timestamp"] = character.Int64ToString(time.Now().UnixNano() / 1e6)
	paramMap["nonce"] = character.GetRandomString(32)
	paramMap["client_id"] = vivoConfig.AppID
	paramMap["grant_type"]= "authorization_code"
	paramMap["code"] = code

	paramLine := ConcatEncodeLine(paramMap)

	sign := MD5SaltHash(paramLine, vivoConfig.AppSecret)
	paramLine += "&sign=" + sign
	resp, err := http.Post("https://passport.vivo.com.cn/oauth/2.0/access_token", "application/x-www-form-urlencoded", strings.NewReader(paramLine))

	if err != nil || resp.StatusCode != http.StatusOK {
		log.Error("error: %v respCode:%v", err, resp.StatusCode)
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Error("vivo ioutil error %v", err)
	}
	log.Debug("vivo response: %v", string(body))
	return string(body), nil
}

func RefreshToken(refreshToken string) (string, error) {
	paramMap := make(map[string]string)
	paramMap["timestamp"] = character.Int64ToString(time.Now().UnixNano() / 1e6)
	paramMap["nonce"] = character.GetRandomString(32)
	paramMap["client_id"] = vivoConfig.AppID
	paramMap["refresh_token"] = refreshToken

	paramLine := ConcatEncodeLine(paramMap)
	sign := MD5SaltHash(paramLine, vivoConfig.AppSecret)
	paramLine += "&" + sign
	resp, err := http.Post("https://passport.vivo.cn/oauth/2.0/refresh_token","application/x-www-form-urlencoded", strings.NewReader(paramLine))
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Error("error: %v respCode:%v", err, resp.StatusCode)
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Error("vivo ioutil error %v", err)
	}
	log.Debug("vivo refresh response: %v", string(body))
	return string(body), nil
}

func ConcatEncodeLine(paramMap map[string]string) string {
	var keys []string
	for key := range paramMap {
		keys = append(keys, key)
	}
	var paramLine string
	sort.Sort(sort.StringSlice(keys))
	for index, key := range keys {
		paramLine += key + "=" + paramMap[key]
		if index == len(keys) - 1 {
			break
		}
		paramLine += "&"
	}
	return paramLine
}

func MD5SaltHash(param string, salt string) string {
	h := md5.New()
	h.Write([]byte(param))
	h.Write([]byte(salt))
	return hex.EncodeToString(h.Sum(nil))
}