package consul

import (
	"errors"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/watch"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

var (
	watching = false
	cw       = &consulWatcher{
		watchers:   make(map[string]*watch.Plan),
		nodes:      make(map[string]string),
		consulAddr: "127.0.0.1:8500"}

	ErrorWatchRunning = errors.New("service watched")
)

type consulWatcher struct {
	nodeWatcher    *watch.Plan            // consul节点监控
	serviceWatcher *watch.Plan            // 可用service列表监控
	watchers       map[string]*watch.Plan // 指定service监控
	consulAddr     string                 // 当前有效的consul节点地址
	nodes          map[string]string      // 当前可用的consul节点列表

	nodeTicker   *time.Ticker // 更新节点信息计时
	client       *api.Client  // 当前打开的api.Client
	consulConfig *api.Config  // 当前consule api.Config

	sync.RWMutex
}

func DoWatch(consul_addr string, service_name ...string) error {
	var err error
	if consul_addr != "" {
		cw.consulAddr = consul_addr
	}
	consulConfig := api.DefaultConfig()
	consulConfig.Address = consul_addr
	cw.client, err = api.NewClient(consulConfig)
	if err != nil {
		return err
	}
	cw.consulConfig = consulConfig

	if cw.nodeWatcher == nil {
		cw.nodeWatcher, err = watch.Parse(map[string]interface{}{"type": "nodes"})
		if err != nil {
			return err
		}
		cw.nodeWatcher.Handler = nodesHandler
		go cw.nodeWatcher.Run(cw.consulAddr)
		RefreshConsulNode(10) //
		return nil
	}
	if len(service_name) == 0 {
		cw.serviceWatcher, err = watch.Parse(map[string]interface{}{"type": "services"})
		if err != nil {
			return err
		}
		cw.serviceWatcher.Handler = servicesHandler
		go cw.serviceWatcher.Run(cw.consulAddr)
		return nil
	}
	for _, sn := range service_name {
		if sn == "" {
			continue
		}
		doServiceWatcher(sn)
	}
	watching = true
	return nil
}

// 定时测试并寻找可用的consul Node地址
func RefreshConsulNode(interval int) {
	if interval == 0 {
		interval = 10
	}
	go func() {
		ticker := time.NewTicker(time.Duration(interval) * time.Second)
	REINTER:
		<-ticker.C
		if checkCurrentClient() {
			changeCurrentNode()
		}
		goto REINTER
	}()
}

// 测试查找可用的consul node
// 依赖nodeWatcher的执行
func checkCurrentClient() (change bool) {
	if cw.consulConfig == nil {
		cw.consulConfig = api.DefaultConfig()
		cw.consulConfig.Address = cw.consulAddr
	}
	if _, err := net.DialTimeout("tcp", cw.consulAddr, 1*time.Second); err == nil {
		return
	}
	delete(cw.nodes, cw.consulAddr)
	times := 0
	cw.RLock()
	defer cw.RUnlock()
	for address, _ := range cw.nodes {
		if times > len(cw.nodes) {
			break
		}
		cw.consulAddr = address
		cw.consulConfig.Address = address
		_, err := net.DialTimeout("tcp", address, 1*time.Second)
		if err == nil {
			cw.client, _ = api.NewClient(cw.consulConfig)
			change = true
			break
		}
		times++
	}
	return
}

func changeCurrentNode() {
	if cw.nodeWatcher != nil && !cw.nodeWatcher.IsStopped() {
		cw.nodeWatcher.Stop()
		cw.nodeWatcher.Run(cw.consulAddr)
	}
	if cw.serviceWatcher != nil && !cw.nodeWatcher.IsStopped() {
		cw.serviceWatcher.Stop()
		cw.serviceWatcher.Run(cw.consulAddr)
	}
	for k, _ := range cw.watchers {
		if cw.watchers[k] != nil && !cw.watchers[k].IsStopped() {
			cw.watchers[k].Stop()
			cw.watchers[k].Run(cw.consulAddr)
		}
	}
}

//
func doServiceWatcher(service_name string) error {
	if _, ok := cw.watchers[service_name]; ok {
		return ErrorWatchRunning
	}
	w, err := watch.Parse(map[string]interface{}{"type": "service", "service": service_name})
	if err != nil {
		return err
	}
	w.Handler = createHandler(service_name)
	go func() {
	REIN:
		err := w.Run(cw.consulAddr)
		if err != nil {
			log.Println(err.Error())
			goto REIN
		}
	}()
	cw.watchers[service_name] = w
	return nil
}

func nodesHandler(idx uint64, data interface{}) {
	nodes, ok := data.([]*api.Node)
	if !ok {
		return
	}
	addressValid := make(map[string]string)
	for _, node := range nodes {
		address := strings.Split(node.Address, ":")
		if address[0] == "" {
			address[0] = "127.0.0.1"
		}
		if len(address) < 2 {
			address = append(address, "8500")
		}
		addressValid[strings.Join(address[:2], ":")] = ""
	}
	cw.Lock()
	for k, _ := range cw.nodes {
		if _, ok := addressValid[k]; !ok {
			delete(cw.nodes, k)
		}
	}
	for k, v := range addressValid {
		cw.nodes[k] = v
	}
	cw.Unlock()
}

func servicesHandler(idx uint64, data interface{}) {
	entries, ok := data.(map[string][]string)
	if !ok {
		return
	}
	services := make(map[string]bool)
	cw.Lock()
	for service, _ := range entries {
		services[service] = true
		doServiceWatcher(service)
	}
	for k, w := range cw.watchers {
		if _, ok := services[k]; ok {
			continue
		}
		if w.IsStopped() {
			continue
		}
		w.Stop()
		delete(cw.watchers, k)
		delete(services_map.servicesUpdating, k)
	}
	cw.Unlock()
	services_map.mux.Lock()
	services_map.services = services_map.servicesUpdating
	services_map.mux.Unlock()
	return
}

func createHandler(name string) func(uint64, interface{}) {
	return func(idx uint64, data interface{}) {
		entries, ok := data.([]*api.ServiceEntry)
		if !ok {
			return
		}
		UpdateServiceEntry(cw.client, name, entries)
		services_map.mux.Lock()
		services_map.services = services_map.servicesUpdating
		services_map.mux.Unlock()
		return
	}
}
