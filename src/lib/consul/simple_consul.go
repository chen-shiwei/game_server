package consul

import (
	"fmt"
)

type ConsulConfig struct {
	Ip         string
	Port       int
	ServerName string
}

var (
	config  ConsulConfig
	address string
)

// get default consul config from file
func Init(cfg ConsulConfig) { // consul config initialize
	config = ConsulConfig{Ip: cfg.Ip, Port: cfg.Port,
		ServerName: cfg.ServerName}
	address = fmt.Sprintf("%s:%d", config.Ip, config.Port)
}

func Watch(service_name ...string) error {
	return DoWatch(address, service_name...)
}

func Discover(service_name string) error {
	return DoDiscover(address, service_name)
}

func RegisterService(service_ip string, service_port int, service_name string) bool {
	return ConsulRegistService(config.Ip, config.Port, service_ip, service_port, service_name, "/check")
}

//func RegisterGrpcDefault() bool {
//	return ConsulRegistGrpc(config)
//}

func RegisterGrpc(service_ip string, service_port int, service_name string) bool {
	return ConsulRegistGrpc5(config.Ip, config.Port, service_ip, service_port, service_name)
}
