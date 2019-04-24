package consul

import (
	"errors"
	"fmt"
	"google.golang.org/grpc"
)

// Grpc
func DoGrpcRequest(service_name string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	serverInfo := PoolingServiceInfo(service_name)
	if serverInfo == nil {
		return nil, errors.New("no valid server")
	}
	address := fmt.Sprintf("%s:%d", serverInfo.IP, serverInfo.Port)
	return grpc.Dial(address, opts...)
}
