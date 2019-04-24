package consul

import (
	"encoding/json"
	"fmt"
	"lib/util"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"

	"github.com/hashicorp/consul/api"
	"lib/sysload"
	"log"
	"strings"
)

type ServiceInfo struct {
	ServiceID string
	IP        string
	Port      int
	Load      int //负载
	Timestamp int //load updated ts
}

type ServiceList []ServiceInfo

type KVData struct {
	Load      int `json:"load"`
	Timestamp int `json:"ts"`
}

type ServicePool struct {
	Pos  int         // position
	List ServiceList // []ServiceInfo
}

type ServicesContainer struct {
	services         map[string]*ServicePool
	servicesUpdating map[string]*ServicePool
	mux              sync.RWMutex
	updateMux        sync.Mutex
}

type ConsulNode struct {
	Host        string
	Port        int
	ServiceName string
}

var (
	services_map = ServicesContainer{services: make(map[string]*ServicePool), servicesUpdating: make(map[string]*ServicePool)}
	consul_nodes = make(map[string]ConsulNode)
	discovering  = false
)

//排序
func (l ServiceList) Len() int           { return len(l) }
func (l ServiceList) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l ServiceList) Less(i, j int) bool { return l[i].Load < l[j].Load }

func ConsulRegistService(consul_host string, consul_port int, service_ip string, service_port int, service_name string, monitor_addr ...string) bool {
	var monitor string
	if service_ip == "0.0.0.0" || service_ip == "*" || service_ip == "" {
		service_ip = util.GetLocalIp()
	}
	if len(monitor_addr) > 0 {
		if strings.HasPrefix(monitor_addr[0], "http") {
			monitor = monitor_addr[0]
		} else {
			if strings.HasPrefix(monitor_addr[0], "/") {
				monitor_addr[0] = monitor_addr[0][1:]
			}
			monitor = fmt.Sprintf("http://%s:%d/%s", service_ip, service_port, monitor_addr[0])
		}
	} else {
		monitor = fmt.Sprintf("%s:%d", service_ip, service_port)
	}
	r := DoRegistService(fmt.Sprintf("%s:%d", consul_host, consul_port), monitor, service_name, service_ip, service_port)
	DiscoverServices(fmt.Sprintf("%s:%d", consul_host, consul_port), true, service_name)
	RegisteDiscoverProcess(ConsulNode{Host: consul_host, Port: consul_port, ServiceName: service_name})
	if !r {
		panic("ConsulRegistService error")
	}
	return r
}

//consul 注册grpc服务
func ConsulRegistGrpc5(consul_host string, consul_port int, service_ip string, service_port int, service_name string) bool {
	return ConsulRegistService(consul_host, consul_port, service_ip, service_port, service_name)
}

// 保留旧接口保持与之前的调用兼容
//func ConsulRegistGrpc(cfg ConsulConfig) bool {
//	return ConsulRegistService(cfg.Ip, cfg.Port), cfg.Grpc.ApiHost, util.Str2int(cfg.Grpc.ApiPort), cfg.Consul.Servername)
//}

func RegisteDiscoverProcess(node ConsulNode) bool {
	consulAddr := fmt.Sprintf("%s:%d", node.Host, node.Port)
	nodeKey := fmt.Sprintf("%s/%s", consulAddr, node.ServiceName)
	if _, ok := consul_nodes[nodeKey]; ok {
		return false
	}
	consul_nodes[nodeKey] = node
	go DoDiscover(consulAddr, node.ServiceName)
	return true
}

//注册服务HTTP
// consul_addr：consul地址
// monitor_addr：健康检查地址 http://127.0.0.1:1234/check
// service_name：服务名称
// ip：ip地址
// port：端口号
func DoRegistService(consul_addr string, monitor_addr string, service_name string, ip string, port int) bool {
	my_service_id := fmt.Sprintf("%s-%s:%d", service_name, ip, port)

	var tags []string
	service := &api.AgentServiceRegistration{
		ID:      my_service_id,
		Name:    service_name,
		Port:    port,
		Address: ip,
		Tags:    tags,
		Check: &api.AgentServiceCheck{
			HTTP:     monitor_addr,
			Interval: "30s",
			Timeout:  "1s",
		},
	}
	if !strings.HasPrefix(monitor_addr, "http") {
		service.Check = &api.AgentServiceCheck{
			TCP:      monitor_addr,
			Interval: "30s",
			Timeout:  "1`s",
		}
	}
	consulConf := api.DefaultConfig()
	consulConf.Address = consul_addr
	client, err := api.NewClient(consulConf)

	if err != nil {
		fmt.Println(err.Error() + "a1")
		return false
	}
	if err := client.Agent().ServiceRegister(service); err != nil {
		fmt.Println(err.Error() + "a2")
		return false
	}
	fmt.Printf("Registered service %q\n", service_name)
	go WaitToUnRegistService(client, my_service_id)
	return true
}

//监听系统信号，如果服务中断或者kill时触发；服务启动是启用携程执行
//consul_client:consul API Client
//my_service_id:service ID
func WaitToUnRegistService(consul_client *api.Client, my_service_id string) {
	quit := make(chan os.Signal, 1)
	//os.Interrupt 表示中断
	//os.Kill 杀死退出进程
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit

	if consul_client == nil {
		return
	}
	if err := consul_client.Agent().ServiceDeregister(my_service_id); err != nil {
		fmt.Println(err)
	}
	os.Exit(1)
}

//心跳检测check
//consul_addr 服务地址
//found_service 要查找的service name；target service;没有时为空
func DoDiscover(consul_addr string, found_service string, interval ...int) error {
	if watching {
		return ErrorWatchRunning
	}
	DiscoverServices(consul_addr, true, found_service)
	timeSet := 5
	if len(interval) > 0 && interval[0] > 0 {
		timeSet = interval[0]
	}
	t := time.NewTicker(time.Second * time.Duration(timeSet))
	discovering = true
	for {
		select {
		case <-t.C:
			DiscoverServices(consul_addr, true, found_service)
		}
	}
	return nil
}

//获取活跃的服务列表
// addr：consul服务地址
// healthyOnly:是否有心跳检测
// service_name：service_name,筛选条件，没有时为空
func DiscoverServices(addr string, healthyOnly bool, service_name string) (servics_map map[string]ServiceList) {
	consulConf := api.DefaultConfig()
	consulConf.Address = addr
	client, err := api.NewClient(consulConf)
	if err != nil {
		fmt.Println(err)
		return
	}

	services_map.servicesUpdating = make(map[string]*ServicePool)

	services, _, err := client.Catalog().Services(&api.QueryOptions{})
	if err != nil {
		fmt.Println(err)
		return
	}

	for name := range services {
		servicesData, _, err := client.Health().Service(name, "", healthyOnly, &api.QueryOptions{})
		CheckErr(err)
		UpdateServiceEntry(client, name, servicesData)
	}
	// assign servicesUpdating to services, start use new copy
	services_map.mux.Lock()
	defer services_map.mux.Unlock()
	services_map.services = services_map.servicesUpdating
	return
}

func UpdateServiceEntry(client *api.Client, name string, entries []*api.ServiceEntry) {
	services_map.updateMux.Lock()
	defer services_map.updateMux.Unlock()
	services_map.servicesUpdating[name] = &ServicePool{List: make([]ServiceInfo, 0)}
	for _, entry := range entries {
		UpdateHealthService(client, entry) // update entry
	}
}

func UpdateHealthService(client *api.Client, entry *api.ServiceEntry) error {
	for _, chk := range entry.Checks {
		if chk.ServiceName == "" {
			continue
		}
		if chk.Status == "critical" {
			continue
		}
		//fmt.Println("  health nodeid:", health.Node, " service_name:", health.ServiceName, " service_id:", health.ServiceID, " status:", health.Status, " ip:", entry.Service.Address, " port:", entry.Service.Port)
		var node ServiceInfo
		if FillServiceInfo(client, chk.ServiceName, chk.ServiceID, &node, entry) != nil {
			log.Println("fetch Service KV info failed:")
		}

		//fmt.Println("service node updated ip:", node.IP, " port:", node.Port, " serviceid:", node.ServiceID, " load:", node.Load, " ts:", node.Timestamp)
		if _, ok := services_map.servicesUpdating[chk.ServiceName]; !ok {
			services_map.servicesUpdating[chk.ServiceName] = &ServicePool{List: make([]ServiceInfo, 0)}
		}
		services_map.servicesUpdating[chk.ServiceName].List = append(services_map.servicesUpdating[chk.ServiceName].List, node)
	}
	return nil
}

func FillServiceInfo(client *api.Client, service_name string, service_id string, node *ServiceInfo, entry *api.ServiceEntry) error {
	node.IP = entry.Service.Address
	node.Port = entry.Service.Port
	node.ServiceID = service_id

	//get data from kv store
	s := GetKeyValue(client, service_name, node.IP, node.Port)
	if len(s) > 0 {
		var data KVData
		if err := json.Unmarshal([]byte(s), &data); err == nil {
			node.Load = data.Load
			node.Timestamp = data.Timestamp
		} else {
			return err
		}
	}
	return nil
}

//更新自己的负载信息到相应的key
func DoUpdateKeyValue(consul_client *api.Client, consul_addr string, service_name string, ip string, port int) {
	t := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-t.C:
			StoreKeyValue(consul_client, consul_addr, service_name, ip, port)
		}
	}
}

func StoreKeyValue(consul_client *api.Client, consul_addr string, service_name string, ip string, port int) {

	my_kv_key := service_name + "/" + ip + ":" + strconv.Itoa(port)

	var data KVData
	// data.Load = rand.Intn(100) //暂时是随机数，到时候更新成负载
	data.Load = sysload.GetSysLoad() // 暂时使用简单的负载评分算法
	data.Timestamp = int(time.Now().Unix())
	bys, _ := json.Marshal(&data)

	kv := &api.KVPair{
		Key:   my_kv_key,
		Flags: 0,
		Value: bys,
	}

	_, err := consul_client.KV().Put(kv, nil)
	CheckErr(err)
	fmt.Println(" store data key:", kv.Key, " value:", string(bys))
}

//获取负载信息
func GetKeyValue(consul_client *api.Client, service_name string, ip string, port int) string {
	key := service_name + "/" + ip + ":" + strconv.Itoa(port)

	kv, _, err := consul_client.KV().Get(key, nil)
	if kv == nil {
		return ""
	}
	CheckErr(err)

	return string(kv.Value)
}

func CheckErr(err error) {
	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}
}

func PoolingServiceInfo(service_name string) (service_info *ServiceInfo) {
	services_map.mux.RLock()
	defer services_map.mux.RUnlock()
	if pool, ok := services_map.services[service_name]; ok {
		if pool.Pos >= len(pool.List) {
			pool.Pos = 0
		}
		service_info = &pool.List[pool.Pos]
		pool.Pos++
	}
	return
}
