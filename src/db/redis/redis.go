package driver

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
)

var (
	redisClient *redis.Pool
	host        = ""
	password    = ""
	db          = 1
)

func init() {
	redisClient = &redis.Pool{
		MaxIdle:     1,
		MaxActive:   10,
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) {
			r, err := redis.Dial(
				"tcp",
				host,
				redis.DialPassword(password), redis.DialDatabase(db),
			)
			if err != nil {
				log.Print("redis client connection error:", err.Error())
				os.Exit(0)
			}
			return r, err
		},
	}
}

// 执行redis
func RCall(comm string, args ...interface{}) (interface{}, error) {
	args_a := toArgs(args)
	rc := redisClient.Get()
	defer rc.Close()
	rslt, err := rc.Do(comm, args_a...)
	return rslt, err
}

// 执行redis 返回字符串
func RString(comm string, args ...interface{}) (string, error) {
	args_a := toArgs(args)
	str, err := redis.String(RCall(comm, args_a...))
	if err != nil && strings.LastIndexAny(err.Error(), "nil returned") != -1 {
		return "", nil
	}
	return str, err
}

// 执行redis返回数值
func RInt(comm string, args ...interface{}) (int, error) {
	args_a := toArgs(args)
	i, err := redis.Int(RCall(comm, args_a...))
	if err != nil && strings.LastIndexAny(err.Error(), "nil returned") != -1 {
		i = 0
		err = nil
	}
	return i, err
}

// 执行redis返回map
func RMap(comm string, args ...interface{}) (map[string]interface{}, error) {
	args_a := toArgs(args)
	values, err := RCall(comm, args_a...)
	if err != nil {
		return nil, err
	}

	data := values.([]interface{})
	l := len(data)
	rMap := make(map[string]interface{})
	for i := 0; i < l; i += 2 {
		rMap[fmt.Sprintf("%s", data[i].(interface{}))] = fmt.Sprintf("%s", data[i+1].(interface{}))
	}
	return rMap, nil
}

// 执行redis返回map
func RStruct(param interface{}, comm string, args ...interface{}) (interface{}, error) {
	args_a := toArgs(args)
	values, err := RCall(comm, args_a...)

	if err != nil {
		return nil, err
	}

	data := values.([]interface{})
	l := len(data)
	rMap := make(map[string]interface{})
	if strings.ToUpper(comm) == "HMGET" {
		for i := 0; i < l; i++ {
			key := fmt.Sprintf("%s", args_a[i+1])
			rMap[key] = data[i]
		}
	} else {
		for i := 0; i < l; i += 2 {
			key := fmt.Sprintf("%s", data[i])
			rMap[key] = data[i+1]
		}
	}
	err = mapToStruct(rMap, param)

	return rMap, err
}

// map 转换struct
func mapToStruct(src map[string]interface{}, dst interface{}) error {
	t := reflect.TypeOf(dst)
	v := reflect.ValueOf(dst)
	if t.Kind() == reflect.Ptr { // 是指针
		v = v.Elem()
		t = t.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		switch v.Field(i).Interface().(type) {
		case int:
			val := src[t.Field(i).Tag.Get("redis")]
			if val != nil {
				switch val.(type) {
				case int:
					v.Field(i).SetInt(int64(val.(int)))
				case string:
					vv, err := strconv.Atoi(val.(string))
					if err != nil {
						return err
					}
					v.Field(i).SetInt(int64(vv))
				}
			}
		case string:
			val := src[t.Field(i).Tag.Get("redis")]
			if val != nil {
				v.Field(i).SetString(fmt.Sprintf("%s", val))
			}
		}
	}
	return nil
}

// 转 args
func toArgs(param []interface{}) []interface{} {
	len := len(param)
	if len == 1 {
		switch param[0].(type) {
		case []interface{}:
			return param[0].([]interface{})
		case map[string]interface{}:
			return mapToArgs(param[0].(map[string]interface{}))
		case int, string:
			return param
		default:
			return structToArgs(param[0])
		}
	}
	if len == 2 {
		switch param[1].(type) {
		case []interface{}:
			tmp := []interface{}{param[0].(interface{})}
			return append(tmp, param[1].([]interface{})...)
		case map[string]interface{}:
			return append([]interface{}{param[0]}, mapToArgs(param[1].(map[string]interface{}))...)
		case int, string:
			return param
		default:
			return append([]interface{}{param[0]}, structToArgs(param[1])...)
		}
	}
	return param
}

// map转参数
func mapToArgs(param map[string]interface{}) []interface{} {
	var data = []interface{}{}
	for k, v := range param {
		data = append(data, strings.ToUpper(k))
		data = append(data, v)
	}
	return data
}

// struct 转参数
func structToArgs(param interface{}) []interface{} {
	t := reflect.TypeOf(param)
	v := reflect.ValueOf(param)
	if t.Kind() == reflect.Ptr { // 是指针
		v = v.Elem()
		t = t.Elem()
	}

	var data = []interface{}{}
	for i := 0; i < t.NumField(); i++ {
		data = append(data, strings.ToUpper(t.Field(i).Name))
		data = append(data, v.Field(i).Interface())
	}

	return data
}
