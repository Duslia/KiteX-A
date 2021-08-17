package main

import (
	"context"
	"fmt"
	"github.com/Duslia997/KiteX-A/KiteX-A/kitex_gen/api"
	"github.com/Duslia997/KiteX-A/KiteX-A/kitex_gen/api/servicea"
	"github.com/KiteX-A/ServiceDiscovery/sd"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/connpool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net"
	"net/http"
	"runtime/debug"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var (
	serverAClient servicea.Client
	count         uint64
	errCount      uint64
)

const (
	Concurrent = 100
)


var (
	qpsCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kitex_rpc_qps",
			Help: "kitex qps",
		},
		[]string{"status"},
	)
)

func sendReq(workID int64, waitGroup *sync.WaitGroup) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Test failed %v, %s.\n", err, string(debug.Stack()))
		}
		waitGroup.Done()
	}()
	var cnt int64 = 0
	reqPrefix := "worker" + strconv.FormatInt(workID, 10) + "'s request"
	for {
		req := api.NewRequest()
		req.SetMessage(reqPrefix + strconv.FormatInt(cnt, 10))
		cnt++
		resp, err := serverAClient.ServiceA(context.Background(), req)
		if err != nil {
			log.Println("resp = ", resp, " err = ", err)
			atomic.AddUint64(&errCount, 1)
			qpsCounter.With(prometheus.Labels{"status": "err"}).Inc()
		} else {
			qpsCounter.With(prometheus.Labels{"status": "success"}).Inc()
		}

		atomic.AddUint64(&count, 1)
		time.Sleep(time.Millisecond * 50)
	}
}

func run() {
	var wg sync.WaitGroup
	wg.Add(Concurrent)
	for i := 0; i < Concurrent; i++ {
		go sendReq(int64(i), &wg)
	}
	wg.Wait()
}

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
	eps := sd.Lookup("kitex.service.a")
	if eps == nil || len(eps) == 0 {
		log.Println("empty service a list")
		eps = append(eps, sd.Endpoint{
			IP:   "kitex-service-a-default-prod",
			Port: "8888",
		})
	}
	options = append(options, client.WithHostPorts(net.JoinHostPort(eps[0].IP, eps[0].Port)))

	serverAClient, err = servicea.NewClient("servicea", options...)
	if err != nil {
		panic(err)
	}

	go func() {
		lastCount := count
		errLastCount := count
		for range time.Tick(time.Second) {
			curCount := atomic.LoadUint64(&count)
			log.Println("qps = ", curCount-lastCount)
			lastCount = curCount

			errCurCount := atomic.LoadUint64(&errCount)
			log.Println("err qps = ", errCurCount-errLastCount)
			errLastCount = errCurCount
		}
	}()
}

func initPrometheus() {
	prometheus.MustRegister(qpsCounter)

	// Expose the registered metrics via HTTP.
	http.Handle("/metrics", promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			// Opt into OpenMetrics to support exemplars.
			EnableOpenMetrics: true,
		},
	))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	go initPrometheus()

	fmt.Println("run")
	run()
	fmt.Println("run exit")
}
