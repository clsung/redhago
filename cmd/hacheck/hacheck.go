package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/clsung/redhago"
	"github.com/golang/glog"
	"github.com/jessevdk/go-flags"
)

var version = "v0.0.1"

type cmdOpts struct {
	OptVersion    bool   `short:"v" long:"version" description:"print the version and exit"`
	OptConfigfile string `long:"config" description:"config file" optional:"yes"`
}

func main() {
	var err error
	var exitCode int

	defer func() { os.Exit(exitCode) }()
	done := make(chan bool)

	opts := &cmdOpts{}
	p := flags.NewParser(opts, flags.Default)
	p.Usage = "[OPTIONS] REDIS1[,REDIS2...]"

	args, err := p.Parse()

	if opts.OptVersion {
		fmt.Fprintf(os.Stderr, "hacheck: %s\n", version)
		return
	}

	if err != nil || len(args) == 0 {
		p.WriteHelp(os.Stderr)
		exitCode = 1
		return
	}

	config := redhago.Config{}
	if opts.OptConfigfile != "" {
		file, err := os.Open(opts.OptConfigfile)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("Using default config")
			} else {
				fmt.Printf("error: %v\n", err)
				exitCode = 1
				return
			}
		} else {
			decoder := json.NewDecoder(file)
			err = decoder.Decode(&config)
			if err != nil {
				fmt.Printf("error: %v\n", err)
				exitCode = 1
				return
			}
		}
	}

	// TODO:
	for i, rs := range config.RedisServer {
		pool := redhago.NewPool(rs.Address, rs.Password, rs.MaxIdleConn, rs.MaxIdleTimeout)
		glog.Infof("%v", pool)
		go func(number int, pool *redhago.Pool) {
			glog.Infof("%d: idleTimeout: %v", number, pool)
			times := 1
			for {
				for j := 0; j < 10; j++ {
					glog.Infof("%d: %d Calling RedisSetWaitGet %d", number, times, j)
					redhago.RedisSetWaitGet(
						pool,
						fmt.Sprintf("pool%d_idle_conn_%d_timeout_%d", number, rs.MaxIdleConn, rs.MaxIdleTimeout),
						"value",
						5, // wait 5 sec
						0,
					)
					glog.Infof("%d: %d Called RedisSetWaitGet %d", number, times, j)
				}
				glog.Infof("%d: Sleep 300 sec", number)
				time.Sleep(300 * time.Second)
				times++
			}
		}(i, pool)
	}

	<-done
	fmt.Println("Exit with 0")
}
