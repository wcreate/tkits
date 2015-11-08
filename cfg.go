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

var (
	WebDomain   = ""
	WebName     = ""
	WebIntro    = ""
	WebListenIP = ""
	WebPort     = 8080
	EmailUser   = ""
	EmailHost   = ""
	EmailPasswd = ""
)

// initialize by the secure config
func init() {
	macaron.SetConfig(GetCfgFile())
	cfg = macaron.Config()

	getWebCfg()
	getSecureCfg()
	getEmailCfg()
}

func getWebCfg() {
	web, err := cfg.GetSection("web")
	if err != nil {
		panic(err)
	}

	WebListenIP = web.Key("ip").MustString("0.0.0.0")
	WebPort = web.Key("port").MustInt(8080)
	WebDomain = web.Key("domain").MustString("localhost")
	WebName = web.Key("name").String()
	WebIntro = web.Key("intro").String()
	log.Debug("web.domain=", WebDomain)
	log.Debug("web.name=", WebName)
	log.Debug("web.intro=", WebIntro)
}

func getEmailCfg() {
	email, err := cfg.GetSection("email")
	if err != nil {
		EmailUser = email.Key("user").String()
		EmailHost = email.Key("host").String()
		EmailPasswd = email.Key("password").String()
		EmailPasswd, _ = scrypto.DecryptStr(EmailPasswd)
	}
}

func getSecureCfg() {
	secure, err := cfg.GetSection("secure")
	if err != nil {
		panic(err)
	}

	factor := secure.Key("factor").String()
	crc := secure.Key("crc").String()
	expire := secure.Key("tokenexpire").MustInt64(15)

	scrypto, err = NewCrypto(factor, crc)
	if err != nil {
		panic(err)
	}

	stoken = NewSimpleToken(scrypto, expire)
}

// Get the config path
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

// retrive the SimpleToken using defalt config
func GetSimpleToken() *SimpleToken {
	return stoken
}

// retrive the Crypto using defalt config
func GetCrypto() *Crypto {
	return scrypto
}
