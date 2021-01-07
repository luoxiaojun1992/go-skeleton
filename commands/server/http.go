package server

import (
	"context"
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/gookit/config/v2"
	"github.com/luoxiaojun1992/go-skeleton/commands"
	httpServerConsts "github.com/luoxiaojun1992/go-skeleton/consts/server/http"
	"github.com/luoxiaojun1992/go-skeleton/router"
	"github.com/luoxiaojun1992/go-skeleton/services/helper"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type HttpServer struct {
	commands.BaseCommand
}

func (hs *HttpServer) startServer() {
	switch config.String("app.runMode", httpServerConsts.RUN_MODE_RELEASE) {
	case httpServerConsts.RUN_MODE_TEST:
		gin.SetMode(gin.TestMode)
	case httpServerConsts.RUN_MODE_DEBUG:
		gin.SetMode(gin.DebugMode)
	case httpServerConsts.RUN_MODE_RELEASE:
		fallthrough
	default:
		gin.SetMode(gin.ReleaseMode)
	}

	srv := &http.Server{
		Addr:    ":" + config.String("app.server.port", httpServerConsts.DEFAULT_PORT),
		Handler: router.Register(),
	}

	go func() {
		if errListen := srv.ListenAndServe(); errListen != nil && errListen != http.ErrServerClosed {
			helper.CheckErrThenPanic("failed to listen", errListen)
		}
	}()

	log.Println("Server started")

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if errShutdown := srv.Shutdown(ctx); errShutdown != nil {
		if errShutdown == context.DeadlineExceeded {
			log.Panicln("Server Shutdown: timeout of 3 seconds.")
		} else {
			helper.CheckErrThenPanic("failed to shutdown server", errShutdown)
		}
	}
	select {
	case <-ctx.Done():
	}
	log.Println("Server exited")
}

func (hs *HttpServer) Handle() {
	hs.startServer()
}

func (hs *HttpServer) ParseOptions(flag *flag.FlagSet) {
	//
}
