package rpc

import (
	"log"
	"net"
	"time"

	"github.com/Duslia997/KiteX-A/KiteX-B/kitex_gen/api/serviceb"
	"github.com/Duslia997/KiteX-A/pkg/zlog"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/connpool"
	"github.com/KiteX-A/ServiceDiscovery/sd"
)

var ServerBClient serviceb.Client

func init() {
	var err error
	options := []client.Option{}
	options = append(options, client.WithLongConnection(connpool.IdleConfig{
		MaxIdlePerAddress: 100,
		MaxIdleGlobal:     1000,
		MaxIdleTimeout:    60 * time.Second,
	}))
	options = append(options, client.WithRPCTimeout(time.Second*5))
	options = append(options, client.WithConnectTimeout(time.Millisecond*50))
	options = append(options, client.WithMiddlewareBuilder(zlog.LogMiddleware))
	options = append(options, client.WithLogger(zlog.New()))
	eps := sd.Lookup("kitex.service.b")
	if eps == nil || len(eps) == 0 {
		log.Println("empty service a list")
		eps = append(eps, sd.Endpoint{
			IP:   "kitex-service-b-default-prod",
			Port: "8888",
		})
	}
	options = append(options, client.WithHostPorts(net.JoinHostPort(eps[0].IP, eps[0].Port)))


	ServerBClient, err = serviceb.NewClient("serviceb", options...)
	if err != nil {
		panic(err)
	}
}
