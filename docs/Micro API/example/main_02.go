package main

import (
	"net/http"
	"strconv"

	// go-micro plugins
	_ "github.com/micro/go-plugins/registry/consul"
	_ "github.com/micro/go-plugins/registry/kubernetes"
	_ "github.com/micro/go-plugins/transport/grpc"
	_ "github.com/micro/go-plugins/transport/tcp"

	"github.com/micro/cli"
	"github.com/micro/micro/api"
	"github.com/micro/micro/cmd"

	// micro plugin
	"github.com/micro/micro/plugin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// micro cmd
func main() {
	cmd.Init()
}

// metrics
var (
	DefObjectives = map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001}
)

func init() {
	namespace := ""
	subsystem := ""
	disable := false
	api.Register(
		plugin.NewPlugin(
			plugin.WithName("metrics"),
			plugin.WithFlag(cli.BoolFlag{
				Name:  "metrics_disable",
				Usage: "disable metrics",
			}),
			plugin.WithInit(func(ctx *cli.Context) error {
				disable = ctx.Bool("metrics_disable")

				return nil
			}),
			plugin.WithCommand(),
			plugin.WithHandler(func(h http.Handler) http.Handler {
				if disable {

				}
				md := make(map[string]string)

				reqTotalCounter := prometheus.NewCounterVec(
					prometheus.CounterOpts{
						Namespace: namespace,
						Subsystem: subsystem,
						Name:      "request_total",
						Help:      "Total request count.",
					},
					[]string{"host", "status"},
				)

				reg := prometheus.NewRegistry()
				wrapreg := prometheus.WrapRegistererWith(md, reg)
				wrapreg.MustRegister(
					prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
					prometheus.NewGoCollector(),
					reqTotalCounter,
				)

				prometheus.DefaultGatherer = reg
				prometheus.DefaultRegisterer = wrapreg

				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// 拦截metrics path，默认"/metrics"
					if r.URL.Path == "/metrics" {
						promhttp.Handler().ServeHTTP(w, r)
						return
					}

					// 为了取到StatusCode对ResponseWriter进行装饰
					ww := WrapWriter{ResponseWriter: w}
					h.ServeHTTP(&ww, r)

					reqTotalCounter.WithLabelValues(r.Host, strconv.Itoa(ww.StatusCode)).Inc()
				})
			}),
		))

}

type WrapWriter struct {
	StatusCode  int
	Size        int64
	wroteHeader bool

	http.ResponseWriter
}

func (ww *WrapWriter) WriteHeader(statusCode int) {
	ww.wroteHeader = true
	ww.StatusCode = statusCode
	ww.ResponseWriter.WriteHeader(statusCode)
}

func (ww *WrapWriter) Write(b []byte) (n int, err error) {
	// 默认200
	if !ww.wroteHeader {
		ww.WriteHeader(http.StatusOK)
	}

	n, err = ww.ResponseWriter.Write(b)
	ww.Size += int64(n)
	return
}

// RequestSize returns the size of request object.
func RequestSize(r *http.Request) float64 {
	size := 0
	if r.URL != nil {
		size = len(r.URL.String())
	}

	size += len(r.Method)
	size += len(r.Proto)

	for name, values := range r.Header {
		size += len(name)
		for _, value := range values {
			size += len(value)
		}
	}
	size += len(r.Host)

	// r.Form and r.MultipartForm are assumed to be included in r.URL.
	if r.ContentLength != -1 {
		size += int(r.ContentLength)
	}
	return float64(size)
}
