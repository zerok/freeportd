package main

import (
	"fmt"
	"net/http"
	"time"

	"gitlab.com/zerok/freeportd"

	"github.com/bluele/gcache"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

type jobResult struct {
	port int
	err  error
}

type job struct {
	responseChannel chan jobResult
}

func main() {
	var httpAddr string
	var gracePeriod time.Duration
	var cacheSize int

	pflag.StringVar(&httpAddr, "http-addr", "localhost:8888", "Address to listen on for HTTP requests")
	pflag.DurationVar(&gracePeriod, "grace-period", time.Minute*5, "Graceperiod for how long ports are kept in the internal store")
	pflag.IntVar(&cacheSize, "cache-size", 1024, "Maximum number of ports in the cache")
	pflag.Parse()

	cache := gcache.New(cacheSize).Simple().Build()
	jobs := make(chan job, 10)

	// We have a single method that is responsible for trying to retrieve
	// and available port
	go func() {
	jobloop:
		for j := range jobs {
			for a := 0; a < cacheSize; a++ {
				port, err := freeportd.GetTCPPort()
				if err != nil {
					log.WithError(err).Error("Failed to fetch port")
					j.responseChannel <- jobResult{err: err}
					close(j.responseChannel)
					continue jobloop
				} else {
					// Check that we don't yet have this port in the cache
					_, err := cache.Get(port)
					if err == gcache.KeyNotFoundError {
						cache.SetWithExpire(port, port, gracePeriod)
						j.responseChannel <- jobResult{port: port}
						close(j.responseChannel)
						continue jobloop
					}
				}
			}
			e := fmt.Errorf("No port found")
			log.WithError(e).Error("Failed to fetch port")
			j.responseChannel <- jobResult{err: e}
			close(j.responseChannel)
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		resp := make(chan jobResult, 1)
		jobs <- job{responseChannel: resp}
		result := <-resp
		if result.err != nil {
			http.Error(w, "Failed to acquire port", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "%d", result.port)
	})

	log.Infof("Starting server on %s", httpAddr)
	http.ListenAndServe(httpAddr, nil)
}
