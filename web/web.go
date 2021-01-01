package web

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"

	"time"

	"github.com/aaaasmile/live-streamer/conf"
	"github.com/aaaasmile/live-streamer/util"
	"github.com/aaaasmile/live-streamer/web/live"
)

func RunService(configfile string) error {

	conf.ReadConfig(configfile)
	log.Println("Configuration is read")
	if err := live.InitFromConfig(conf.Current.DebugVerbose, conf.Current.DBPath); err != nil {
		return err
	}
	util.UseRelativeRoot = conf.Current.UseRelativeRoot

	serverurl := conf.Current.ServiceURL
	finalServURL := fmt.Sprintf("http://%s%s", strings.Replace(serverurl, "0.0.0.0", "localhost", 1), conf.Current.RootURLPattern)
	finalServURL = strings.Replace(finalServURL, "127.0.0.1", "localhost", 1)
	log.Println("Server started with URL ", serverurl)
	log.Println("Try this url: ", finalServURL)

	http.Handle(conf.Current.RootURLPattern+"static/", http.StripPrefix(conf.Current.RootURLPattern+"static", http.FileServer(http.Dir(util.GetFullPath("static")))))
	http.HandleFunc(conf.Current.RootURLPattern, live.APiHandler)
	http.HandleFunc("/websocket", live.WsHandler)

	srv := &http.Server{
		Addr: serverurl,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      nil,
	}

	chShutdown := make(chan struct{}, 1)
	go func(chs chan struct{}) {
		if err := srv.ListenAndServe(); err != nil {
			log.Println("Server is not listening anymore: ", err)
			chs <- struct{}{}
		}
	}(chShutdown)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt) //We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	log.Println("Enter in server loop")
loop:
	for {
		select {
		case <-sig:
			log.Println("stop because interrupt")
			break loop
		case <-chShutdown:
			log.Println("stop because service shutdown on listening")
			log.Fatal("Force with an error to restart")
			break loop
		}
	}

	live.HandlerShutdown()

	wait := 3 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	err := srv.Shutdown(ctx)
	if err != nil {
		log.Println("Shutdown error: ", err)
	}

	log.Println("Bye, service")
	return nil
}
