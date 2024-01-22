package main

import (
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/skeleton1231/gotal/internal/authzserver"
)

func main() {
	rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	if len(os.Getenv("GOMAXPROCS")) == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}
	authzserver.NewApp("gotal-authz-server").Run()
}
