// Copyright @ 2020 OPS Inc.
//
// Author: Jinlong Yang
//

package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"s3/http/g"
)

var (
	cfgFile = flag.String("c", "../etc/dev.yaml", "yaml configuration file.")
	help    = flag.Bool("h", false, "show help info.")
)

func main() {
	flag.Parse()
	if *help {
		flag.PrintDefaults()
		return
	}
	g.ParseConfig(*cfgFile)
	g.InitLogger()

	ctx, cancel := context.WithCancel(context.Background())

	qs := make(chan os.Signal)
	signal.Notify(qs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGPIPE)

	r := gin.Default()
	urls(r)

	server := http.Server{
		Addr:           g.Config().Server.Address,
		Handler:        r,
		ReadTimeout:    0,    // 0表示没有限制
		WriteTimeout:   0,    // 0表示没有限制
		MaxHeaderBytes: 8192, // 8K
	}

	go func() {
		if err := server.ListenAndServe(); err == http.ErrServerClosed {
			log.Info("Listen and serve shutdown....")
		} else if err != nil {
			log.Errorf("Listen and serve boot failed: %s", err)
			cancel()
		}
	}()

	select {
	case sig := <-qs:
		if sig == syscall.SIGINT || sig == syscall.SIGTERM || sig == syscall.SIGQUIT {
			log.Info("Quit the server with Ctrl C.")
			server.Shutdown(ctx)
			cancel()
		} else if sig == syscall.SIGPIPE {
			log.Warn("Ignore broken pipe signal.")
		}
	}

	log.Info("S3 shutdown.....")
	time.Sleep(4 * time.Second)
}
