package main

import (
	"flag"
	"github.com/jorgeluis594/go_indexer/indexer"
	"github.com/jorgeluis594/go_indexer/lib"
	"github.com/jorgeluis594/go_indexer/repository"
	"log"
	"os"
	"runtime/pprof"
)

func main() {
	f, _ := os.Create("cpu.pprof")
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	path := flag.String("path", "", "path to index")
	host := flag.String("host", "", "host of Zinc Search client")
	username := flag.String("username", "", "username of db")
	password := flag.String("password", "", "password of db")
	flag.Parse()

	clientHttp := repository.InitHttpClient(*host, *username, *password)
	repositoryDB := repository.InitRepository(clientHttp, "email_copy")
	directory, err := indexer.InitDirectory(*path)

	if err != nil {
		log.Fatal("Error reading directory: ", *path)
	}

	processor := lib.InitProcessor(directory.GetPaths(), repositoryDB)
	processor.Process()
}
