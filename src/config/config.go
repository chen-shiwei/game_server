package config

import (
	"os"
	"strconv"
	"strings"
)

var (
	configServer     string
	configServerProd string = "cfg.h.jskj.cn"
	configServerDev  string = "config.localhost"

	configItem = make(map[string]interface{})
	loaded     = false
)

func init() {
	if strings.ToUpper(os.Getenv("SERVER_ENV")) == "DEV" {
		configServer = configServerDev
	} else {
		configServer = configServerProd
	}
	loadConfig(false)
}

func loadConfig(force bool) error {
	if force || !loaded {
		// just for development enviroment
		configItem["port"] = 10086
		configItem["logFile"] = "./server.log"
		configItem["userServer"] = "localhost:8899"
	}
	return nil
}

func Load(force bool) error {
	return loadConfig(force)
}

func GetInt(name string) int {
	val := getItem(name, "int")
	return val.(int)
}

func GetBool(name string) bool {
	val := getItem(name, "bool")
	return val.(bool)
}

func GetString(name string) string {
	val := getItem(name, "string")
	return val.(string)
}

func getItem(name string, typ string) (value interface{}) {
	if _, ok := configItem[name]; !ok {
		return
	} else {
		value = configItem[name]
	}
	realType := typeAsString(value)
	if typ == realType {
		return
	}
	if typ == "int" && realType == "string" {
		value, _ = strconv.Atoi(value.(string))
	} else if typ == "string" && realType == "int" {
		value = strconv.Itoa(value.(int))
	} else if typ == "bool" && realType == "int" {
		if value.(int) == 0 {
			value = false
		} else {
			value = true
		}
	} else if typ == "bool" && realType == "string" {
		if len(value.(string)) > 0 {
			value = true
		} else {
			value = false
		}
	}
	return
}

func typeAsString(v interface{}) string {
	switch v.(type) {
	case int, uint, int64, uint64, int32, uint32:
		return "int"
	case string:
		return "string"
	case bool:
		return "bool"
	default:
		return "unknown"
	}
}
