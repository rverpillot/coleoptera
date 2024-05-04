package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"runtime"
	"runtime/pprof"
	"strings"
	"time"

	_ "net/http/pprof"

	"github.com/rverpillot/coleoptera/pages"
	"github.com/rverpillot/ihui"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db          *gorm.DB
	contextRoot *string
	debug       *bool
	tmpDir      = path.Join(os.TempDir(), "coleoptera")
)

func start(session *ihui.Session) error {
	session.Set("db", db)
	session.Set("admin", *debug)
	session.Set("tmpDir", tmpDir)
	session.Set("debug", *debug)

	menu := pages.NewMenu()
	menu.Add("especes", "Espèces", pages.NewPageEspeces(menu))
	menu.Add("individus", "Individus", pages.NewPageIndividus(menu))
	menu.Add("plan", "Plan", pages.NewPagePlan(menu))
	if err := session.AddElement("menu", menu, nil); err != nil {
		return err
	}
	return menu.ShowItem(session, "individus")
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s database\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	contextRoot = flag.String("context", "/", "`context_root` (ex: /coleoptera)")
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

	level := logger.Error
	if *debug {
		level = logger.Info
		pages.SetDebugMode()
	}
	db, err = gorm.Open(sqlite.Open(baseName), &gorm.Config{
		Logger: logger.Default.LogMode(level),
	})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	// Create & purge tmp dir
	os.MkdirAll(tmpDir, 0755)
	go func() {
		for now := range time.Tick(time.Minute * 20) {
			if err := purgeTmp(tmpDir, now); err != nil {
				log.Print(err)
			}
		}
	}()

	if !strings.HasSuffix(*contextRoot, "/") {
		*contextRoot += "/"
	}
	staticsFS, _ := fs.Sub(pages.ResourcesFs, "statics")
	http.Handle(*contextRoot, http.StripPrefix(*contextRoot, http.FileServer(http.FS(staticsFS))))
	http.Handle(*contextRoot+"tmp/", http.StripPrefix(*contextRoot+"tmp", http.FileServer(http.Dir(tmpDir))))
	ihui.Handle(*contextRoot, start)

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

func purgeTmp(dirName string, now time.Time) error {
	tmpDir, err := os.Open(dirName)
	if err != nil {
		return err
	}
	defer tmpDir.Close()

	files, err := tmpDir.Readdir(0)
	if err != nil {
		return err
	}
	for _, file := range files {
		if now.Sub(file.ModTime()).Hours() > 48 {
			os.Remove(path.Join(dirName, file.Name()))
		}
	}
	return nil
}
