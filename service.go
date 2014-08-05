package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/gorilla/pat"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	// in MB
	memory := r.URL.Query().Get(":memory")
	if memory == "" {
		memory = r.URL.Query().Get("memory")
	}
	MBToAlloc, _ := strconv.Atoi(memory)

	// Microseconds
	duration := r.URL.Query().Get(":duration")
	if duration == "" {
		duration = r.URL.Query().Get("duration")
	}
	allocDuration, _ := strconv.Atoi(duration)
	start := time.Now()
	fmt.Fprintf(w, "%dMB during %dms\n", MBToAlloc, allocDuration)

	alloc := make([]byte, MBToAlloc*1024*1024)
	for i := range alloc {
		alloc[i] = 'c'
	}

	time.Sleep(time.Duration(allocDuration) * time.Millisecond)
	alloc = nil
	debug.FreeOSMemory()

	end := time.Now()
	log.Printf("Memory: %vMB, Duration: %vms 200 OK (%v)", MBToAlloc, allocDuration, end.Sub(start).String())
}

func main() {
	r := pat.New()
	r.Get("/{memory}/{duration}", IndexHandler)
	r.Get("/", IndexHandler)
	http.Handle("/", r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Printf("Listen on 0.0.0.0:%s\n", port)
	http.ListenAndServe(":"+port, nil)
}
