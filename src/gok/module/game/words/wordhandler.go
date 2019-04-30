package words

import (
	"encoding/csv"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"gok/config"
)

var firstName []string = []string{}
var lastName []string = []string{}

//敏感词汇
var filters []string = []string{}

var regExp *regexp.Regexp

//func ValidateNickname(nickname string) {
//	if !regExp.MatchString(nickname) {
//		exception.GameException(exception.NICKNAME_INVALID_WORD)
//	}
//
//	for _, filter := range filters {
//		if strings.Contains(nickname, filter) {
//			//log.Debug("%v", filter)
//			exception.GameException(exception.NICKNAME_SENSITIVE_WORD)
//		}
//	}
//}

//随机姓名
func RandomName() string {
	//rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	return firstName[rand.Intn(len(firstName))] + lastName[rand.Intn(len(lastName))]
}

func init() {
	regexp, _ := regexp.Compile("^[\u4e00-\u9fa5a-zA-Z0-9_]+$")
	regExp = regexp

	rfile, _ := os.Open(config.GetConfigPath("conf/gok/game/name.csv"))
	r := csv.NewReader(rfile)
	for {
		strs, err := r.Read()
		if err != nil {
			break
		}
		firstName = append(firstName, strs[0])
		lastName = append(lastName, strs[1])
	}

	rSensitive, _ := os.Open(config.GetConfigPath("conf/gok/game/sensitive.csv"))
	r = csv.NewReader(rSensitive)
	for {
		strs, err := r.Read()
		if err != nil {
			break
		}
		for _, filter := range strs {
			filter = strings.TrimSpace(filter)
			if filter != "" {
				filters = append(filters, filter)
			}
		}
	}
}
