package main

import (
	"flag"
	"os"
	"strconv"

	_ "net/http/pprof"

	"github.com/seiflotfy/skizze/config"
	"github.com/seiflotfy/skizze/server"
	"github.com/seiflotfy/skizze/utils"
)

var logger = utils.GetLogger()

func main() {
	var port uint
	flag.UintVar(&port, "p", 3596, "specifies the port for Counts to run on")
	flag.Parse()

	//TODO: Add arguments for dataDir and infoDir

	err := os.Setenv("COUNTS_PORT", strconv.Itoa(int(port)))
	utils.PanicOnError(err)

	logger.Info.Println("Starting counts...")
	conf := config.GetConfig()
	logger.Info.Println("Using data dir: ", conf.DataDir)
	server, err := server.New()
	utils.PanicOnError(err)
	server.Run()
}
