package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/felixge/fgprof"

	"github.com/chhetripradeep/chtop/cmd"
)

func main() {
	http.DefaultServeMux.Handle("/debug/fgprof", fgprof.Handler())
	go func() {
		fmt.Fprintln(os.Stderr, http.ListenAndServe(":6060", nil))
	}()
	cmd.Execute()
}
