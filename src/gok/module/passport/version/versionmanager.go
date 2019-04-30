package version

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"aliens/log"
	"gok/module/passport/conf"
)

var VersionManager = &VersionInfo{}

type VersionInfo struct {
	VersionInfoMapping		map[string]string
}

func Init() {
	VersionManager.RefreshVersionInfo()
}

func (this *VersionInfo) RefreshVersionInfo() bool{
	httpResp,err := http.Get(conf.Server.VersionUrl)
	if err != nil {
		log.Error(err.Error())
		return  false
	}
	defer httpResp.Body.Close()
	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		log.Error(err.Error())
		return  false
	}
	this.VersionInfoMapping = make(map[string]string)
	versionMapping := this.VersionInfoMapping
	json.Unmarshal(body, &versionMapping)
	log.Info("client version mapping:%v",versionMapping)
	return true
}
