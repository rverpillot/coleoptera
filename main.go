package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"runtime"
	"runtime/pprof"
	"rverpi/coleoptera/pages"

	"rverpi/ihui.v2"

	_ "net/http/pprof"

	"github.com/facebookgo/inject"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	db          *gorm.DB
	contextRoot *string
	debug       *bool
)

func start(session *ihui.Session) {
	session.Set("db", db)
	session.Set("admin", *debug)

	menu := pages.NewMenu("individus")
	inject.Populate(menu, pages.NewPageIndividus(menu), pages.NewPageEspeces(menu), pages.NewPagePlan(menu))
	session.ShowPage("Coleoptera", menu.PageIndividus)
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s database\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	contextRoot = flag.String("context", "/", "`context_root` (ex: /app)")
	var address = flag.String("listen", "127.0.0.1:8080", "`address` (ip:port)")
	var secure = flag.Bool("secure", false, "SSL")
	var certFile = flag.String("certFile", "cert.pem", "`certFile`")
	var keyFile = flag.String("keyFile", "key.pem", "`keyFile`")
	var cpuprofile = flag.String("cpuprofile", "", "write cpu profile `file`")
	var memprofile = flag.String("memprofile", "", "write memory profile to `file`")
	debug = flag.Bool("debug", false, "debug mode")
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(0)
	}

	baseName := flag.Arg(0)

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
		f.Close()
	}

	_, err := os.Stat(baseName)
	if err != nil {
		log.Fatalln(err)
	}

	db, err = gorm.Open("sqlite3", baseName)
	if err != nil {
		log.Fatal("failed to connect database")
	}
	defer db.Close()
	db.LogMode(*debug)

	paths := []string{"js", "css", "images"}
	for _, p := range paths {
		path := path.Join(*contextRoot, p) + "/"
		log.Println(path)
		http.Handle(path, http.StripPrefix(*contextRoot, http.FileServer(pages.ResourcesBox.HTTPBox())))
	}

	http.Handle(os.path.Join(*contextRoot, "/app"), ihui.NewHTTPHandler(start))

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func(c chan os.Signal) {
		<-c
		pprof.StopCPUProfile()
		log.Println("Exit.")
		os.Exit(0)
	}(c)

	if *secure {
		log.Printf("Listen to https://%s\n", *address)
		log.Fatal(http.ListenAndServeTLS(*address, *certFile, *keyFile, nil))
	} else {
		log.Printf("Listen to http://%s\n", *address)
		log.Fatal(http.ListenAndServe(*address, nil))
	}

}
