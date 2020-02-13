package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/siacentral/host-dashboard/daemon/persist"
	"github.com/siacentral/host-dashboard/daemon/sync"
	"github.com/siacentral/host-dashboard/daemon/web"
	"github.com/siacentral/host-dashboard/daemon/web/router"
)

var (
	dataPath   string
	listenAddr string
	siaAddr    string
	allowCORS  bool
	logFile    *os.File
)

func writeLine(format string, args ...interface{}) {
	os.Stdout.WriteString(fmt.Sprintf(format, args...) + "\n")
	log.Printf(format, args...)
}

func init() {
	var err error

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	flag.StringVar(&dataPath, "data-path", "data", "the data path to use")
	flag.StringVar(&listenAddr, "listen-addr", ":8884", "the address to listen on, defaults to :8884")
	flag.StringVar(&siaAddr, "sia-api-addr", "localhost:9980", "the url used to connect to Sia. Defaults to \"localhost:9980\"")
	flag.BoolVar(&allowCORS, "allow-cors", false, "enables cross-origin requests, this should only be enabled for development or specific use cases")
	flag.Parse()

	if err := os.MkdirAll(dataPath, 0770); err != nil && !os.IsExist(err) {
		log.Fatalf("error creating directory: %s", err)
	}

	if err := persist.InitializeDB(dataPath); err != nil {
		log.Fatalf("error initializing database: %s", err)
	}

	logFile, err = os.OpenFile(filepath.Join(dataPath, "log.log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatalf("error opening log: %s", err)
	}

	// log.SetOutput(logFile)
}

func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	if err != nil {
		log.Printf("Warn: unable to open browser: %s", err)
	}
}

func startAPI() {
	if err := web.Start(router.APIOptions{
		ListenAddress: listenAddr,
		CORS: router.CORSOptions{
			Enabled: allowCORS,
			Origins: []string{"*"},
			Methods: []string{"*"},
		},
		RateInterval: time.Second,
		RateLimit:    10,
	}); err != nil {
		writeLine("Error starting API: %s", err)
		os.Exit(1)
	}
}

func main() {
	writeLine("Starting Host Dashboard")

	if err := sync.Start(siaAddr); err != nil {
		log.Fatalf("error syncing data: %s", err)
	}

	go startAPI()

	log.Printf("Host Dashboard Ready at: http://localhost:%d", 8884)

	openbrowser("http://localhost:8884")

	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, os.Interrupt, os.Kill, syscall.SIGTERM)

	<-sigChan

	writeLine("Shutting down")

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)

	defer cancelFunc()

	if err := web.Shutdown(ctx); err != nil {
		log.Fatalln(err)
	}

	if err := persist.CloseDB(); err != nil {
		log.Fatalln(err)
	}

	if logFile != nil {
		if err := logFile.Close(); err != nil {
			writeLine("closing log: %s", err)
			os.Exit(1)
		}
	}
}
