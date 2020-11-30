package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"sync"
	"syscall"

	_ "github.com/lib/pq"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/reactserv/internal/constants"
	"github.com/qdm12/reactserv/internal/health"
	"github.com/qdm12/reactserv/internal/memfs"
	"github.com/qdm12/reactserv/internal/models"
	"github.com/qdm12/reactserv/internal/params"
	"github.com/qdm12/reactserv/internal/server"
	"github.com/qdm12/reactserv/internal/splash"
)

//nolint:gochecknoglobals
var (
	buildInfo models.BuildInformation
	version   = "unknown"
	commit    = "unknown"
	buildDate = "an unknown date"
)

func main() {
	buildInfo.Version = version
	buildInfo.Commit = commit
	buildInfo.BuildDate = buildDate
	ctx := context.Background()
	os.Exit(_main(ctx, os.Args))
}

func _main(ctx context.Context, _ []string) int {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	if health.IsClientMode(os.Args) {
		// Running the program in a separate instance through the Docker
		// built-in healthcheck, in an ephemeral fashion to query the
		// long running instance of the program about its status
		client := health.NewClient()
		if err := client.Query(ctx); err != nil {
			fmt.Println(err)
			return 1
		}
		return 0
	}
	paramsReader := params.NewReader()
	fmt.Println(splash.Splash(buildInfo))
	logger, err := createLogger(paramsReader)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	listeningPort, warning, err := paramsReader.GetListeningPort()
	if len(warning) > 0 {
		logger.Warn(warning)
	}
	if err != nil {
		logger.Error(err)
		return 1
	}
	rootURL, err := paramsReader.GetRootURL()
	if err != nil {
		logger.Error(err)
		return 1
	}
	rootDir, err := paramsReader.GetRootDir()
	if err != nil {
		logger.Error(err)
		return 1
	}

	oldToNew := map[*regexp.Regexp]string{
		constants.RegexStatic:       rootURL + "/static/",
		constants.RegexManifestJSON: rootURL + "/manifest.json",
	}
	memFSLogger := logger.WithPrefix("memory filesystem: ")
	memFS, err := memfs.New(rootDir, oldToNew, memFSLogger)
	if err != nil {
		logger.Error(err)
		return 1
	}

	wg := &sync.WaitGroup{}
	defer wg.Wait()

	go memFS.Watch(ctx, wg)

	serverLogger := logger.WithPrefix("http server: ")
	address := fmt.Sprintf("0.0.0.0:%d", listeningPort)
	server := server.New(address, rootURL, serverLogger, memFS, buildInfo)
	wg.Add(1)
	go server.Run(ctx, wg)

	const healthServerAddr = "127.0.0.1:9999"
	healthcheck := func() error { return nil }
	healthServer := health.NewServer(healthServerAddr, logger.WithPrefix("healthcheck server: "), healthcheck)
	wg.Add(1)
	go healthServer.Run(ctx, wg)

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals,
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	select {
	case signal := <-osSignals:
		message := fmt.Sprintf("Stopping program: caught OS signal %q", signal)
		logger.Warn(message)
		cancel()
		return 1
	case <-ctx.Done():
		message := fmt.Sprintf("Stopping program: %s", ctx.Err())
		logger.Warn(message)
		return 1
	}
}

func createLogger(paramsReader params.Reader) (logger logging.Logger, err error) {
	encoding, level, err := paramsReader.GetLoggerConfig()
	if err != nil {
		return nil, err
	}
	return logging.NewLogger(encoding, level)
}
