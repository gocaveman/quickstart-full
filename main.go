package main

import (
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gocaveman/caveman/filesystem/fsutil"
	"github.com/gocaveman/caveman/renderer"
	"github.com/gocaveman/caveman/webutil"
	"github.com/gocaveman/caveman/webutil/htmlmin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Change this to be the name of your project.
var PROGRAM_NAME = "quickstart-full"

func main() {

	pflag.StringP("workdir", "d", "", "Change directory to this value before resolving relative paths")
	pflag.StringP("http-listen", "l", ":8080", "IP:Port to listen on for HTTP")
	pflag.StringP("static-dir", "", "static", "Directory for static resources (CSS, JS, images, etc.)")
	pflag.StringP("views-dir", "", "views", "Directory for view template files (pages)")
	pflag.StringP("includes-dir", "", "includes", "Directory for include template files")
	pflag.BoolP("optimize", "o", false, "Enable optimizations such as CSS/JS file combining and minification")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	if wd := viper.GetString("workdir"); wd != "" {
		if err := os.Chdir(wd); err != nil {
			log.Fatalf("Error changing directory to %q: %v", wd, err)
		}
	}

	rend := renderer.New(
		http.Dir(fsutil.MustAbs(viper.GetString("views-dir"))),
		http.Dir(fsutil.MustAbs(viper.GetString("includes-dir"))),
	)

	var hl webutil.HandlerList

	hl = append(hl, webutil.NewContextCancelHandler())
	hl = append(hl, webutil.NewGzipHandler())

	hl = append(hl, webutil.NewCtxSetHandler("main.optimize", viper.GetBool("optimize")))
	if viper.GetBool("optimize") {
		hl = append(hl, htmlmin.NewHandler())
	}

	hl = append(hl, webutil.NewRequestContextHandler())

	hl = append(hl, renderer.NewHandler(rend))

	if staticDir := viper.GetString("static-dir"); staticDir != "" {
		hl = append(hl, webutil.NewStaticFileHandler(http.Dir(fsutil.MustAbs(staticDir))))
	}

	hl = append(hl, http.NotFoundHandler())

	var wg sync.WaitGroup

	webutil.StartHTTPServer(&http.Server{
		Addr:    viper.GetString("http-listen"),
		Handler: hl.WithCloseHandler(),
	}, &wg)

	wg.Wait()

}
