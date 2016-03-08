package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang/groupcache"
)

var (
	gcGroup *groupcache.Group
)

func main() {
	hostname, _ := os.Hostname()
	groupcachePort := 8081
	groupcacheHost := fmt.Sprintf("http://%s:%d", hostname, groupcachePort)

	gcOptions := &groupcache.HTTPPoolOptions{
		Replicas: 2,
		BasePath: "/_groupcache/",
	}
	gcPool := groupcache.NewHTTPPoolOpts(groupcacheHost, gcOptions)
	gcGroup = groupcache.NewGroup("imageCache", 64<<20, groupcache.GetterFunc(imageCacheGetter))

	http.Handle(gcOptions.BasePath, gcPool)

	http.HandleFunc("/render", handleRender)
	http.ListenAndServe(":8080", nil)

	log.Println(gcPool)
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func imageCacheGetter(context groupcache.Context, key string, dest groupcache.Sink) error {
	log.Println("Calculating...", key)

	dest.SetString("Result")

	return nil
}

func handleRender(w http.ResponseWriter, r *http.Request) {
	var buf []byte

	gcGroup.Get(nil, strconv.Itoa(random(0, 10)), groupcache.AllocatingByteSliceSink(&buf))

	log.Println("Got", string(buf))
}
