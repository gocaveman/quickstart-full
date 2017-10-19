package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gocaveman/caveman/webutil"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Change this to be the name of your project.
var PROGRAM_NAME = "quickstart-full"

func main() {

	pflag.StringP("http-listen", "l", ":8080", "IP:Port to listen on for HTTP")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	var hl webutil.HandlerList

	hl = append(hl, webutil.NewContextCancelHandler())
	hl = append(hl, webutil.NewGzipHandler())
	hl = append(hl, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path == "/" {
			w.Header().Set("content-type", "text/html")
			fmt.Fprintf(w, "<html><body>Testing!</body></html>")
			return
		}

	}))
	hl = append(hl, http.NotFoundHandler())

	// TODO: I think we can package up the creation of the HTTP server
	// with good defaults, the log line and the listen into a helper
	// in webutil - would be cleaner in here - it should probably
	// start a new goroutine and return if error, so we can start
	// multiple etc.
	httpServer := http.Server{
		Addr:           viper.GetString("http-listen"),
		Handler:        hl.WithCloseHandler(),
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("Listening for HTTP %q", httpServer.Addr)
	log.Fatal(httpServer.ListenAndServe())

	select {} // block forever

}
