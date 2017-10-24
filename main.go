package main

import (
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gocaveman/caveman/filesystem/fsutil"
	"github.com/gocaveman/caveman/renderer"
	"github.com/gocaveman/caveman/renderer/includeregistry"
	"github.com/gocaveman/caveman/renderer/viewregistry"
	"github.com/gocaveman/caveman/webutil"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	_ "github.com/gocaveman/caveman/demo/demotheme"
)

// EDITME: Change this to be the name of your project.
var PROGRAM_NAME = "quickstart-full" // FIXME: should we able to configure this via viper?

func main() {

	pflag.StringP("workdir", "d", "", "Change directory to this value before resolving relative paths")
	pflag.StringP("http-listen", "l", ":8080", "IP:Port to listen on for HTTP")
	pflag.StringP("static-dir", "", "static", "Directory for static resources (CSS, JS, images, etc.)")
	pflag.StringP("views-dir", "", "views", "Directory for view template files (pages)")
	pflag.StringP("includes-dir", "", "includes", "Directory for include template files")
	pflag.BoolP("optimize", "o", false, "Enable optimizations such as CSS/JS file combining and minification")
	pflag.BoolP("debug", "g", false, "Enable debug output (intended for development only)")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	if wd := viper.GetString("workdir"); wd != "" {
		if err := os.Chdir(wd); err != nil {
			log.Fatalf("Error changing directory to %q: %v", wd, err)
		}
	}

	// TODO: we should put in viper's config search stuff here, as well as the setup for environment
	// including the prefix and -_ replacement.

	includeSequence := includeregistry.Contents()

	// EDITME: You can modify items from the include registry here before using it.

	// includeSequence...

	// /EDITME

	includeFS := includeregistry.MakeFS(http.Dir(fsutil.MustAbs(viper.GetString("includes-dir"))), includeSequence, viper.GetBool("debug"))

	viewSequence := viewregistry.Contents()

	// EDITME: You can modify items from the view registry here before using it.

	// viewSequence...

	// /EDITME

	viewFS := viewregistry.MakeFS(http.Dir(fsutil.MustAbs(viper.GetString("views-dir"))), viewSequence, viper.GetBool("debug"))

	rend := renderer.New(
		viewFS,
		includeFS,
	)

	var hl webutil.HandlerList

	hl = append(hl, webutil.NewContextCancelHandler())
	hl = append(hl, webutil.NewGzipHandler())

	hl = append(hl, webutil.NewCtxSetHandler("main.optimize", viper.GetBool("optimize")))
	if viper.GetBool("optimize") {
		// disabling htmlmin handler for now, seems to be quite aggressive and removing things it should not
		// hl = append(hl, htmlmin.NewHandler())
	}

	hl = append(hl, webutil.NewRequestContextHandler())

	// EDITME: Custom handlers here, and you can edit the data from the handler registry as needed

	// ...

	// /EDITME

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
