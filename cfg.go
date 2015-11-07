package tkits

import (
	"os"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/ini.v1"
	"gopkg.in/macaron.v1"
)

var (
	cfg     *ini.File
	stoken  *SimpleToken
	scrypto *Crypto
)

func init() {
	macaron.SetConfig(GetCfgFile())
	cfg = macaron.Config()

	secure, err := cfg.GetSection("secure")
	if err != nil {
		panic(err)
	}

	factor := secure.Key("factor").String()
	crc := secure.Key("crc").String()
	expire := secure.Key("token").MustFloat64(15.0)
	//println(factor, crc)

	scrypto, err = NewCrypto(factor, crc)
	if err != nil {
		panic(err)
	}

	stoken = NewSimpleToken(scrypto, expire)
}

func GetCfgFile() string {
	workPath, _ := os.Getwd()
	workPath, _ = filepath.Abs(workPath)
	// initialize default configurations
	appPath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	configPath := filepath.Join(appPath, "conf", "app.ini")
	if workPath != appPath {
		if FileExists(configPath) {
			os.Chdir(appPath)
		} else if strings.HasSuffix(workPath, "home") {
			configPath = filepath.Join(workPath, "conf", "app.ini")
		} else {
			configPath = filepath.Join(workPath, "../conf", "app.ini")
		}
	}
	log.Debug("config path=", configPath)
	return configPath
}

func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func GetSimpleToken() *SimpleToken {
	return stoken
}

func GetCrypto() *Crypto {
	return scrypto
}
