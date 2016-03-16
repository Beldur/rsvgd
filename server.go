package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang/groupcache"
)

type server struct {
	gcPool  *groupcache.HTTPPool
	gcGroup *groupcache.Group
	gcPort  int
	port    int
}

func newServer(bindAddress string, port int) *server {
	groupcachePort := port + 1
	groupcacheHost := fmt.Sprintf("http://%s:%d", bindAddress, groupcachePort)
	gcOptions := &groupcache.HTTPPoolOptions{
		Replicas: 2,
		BasePath: "/_groupcache/",
	}

	gcPool := groupcache.NewHTTPPoolOpts(groupcacheHost, gcOptions)
	gcGroup := groupcache.NewGroup("cache", 64<<20, groupcache.GetterFunc(cacheGetter))

	return &server{
		gcPool:  gcPool,
		port:    port,
		gcPort:  groupcachePort,
		gcGroup: gcGroup,
	}
}

func (s *server) start() {
	go http.ListenAndServe(fmt.Sprintf(":%d", s.gcPort), s.gcPool)
	http.HandleFunc("/render", s.handleRender)
	http.HandleFunc("/info", s.handleInfo)
	http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil)
}

// SetPeers updates the servers peer list
func (s *server) SetPeers(peerList []string) {
	s.gcPool.Set(peerList...)
}

// handleRender takes the incoming SVG file, renders it and returns the
// resulting image to the client.
func (s *server) handleRender(w http.ResponseWriter, r *http.Request) {
	var buf []byte

	s.gcGroup.Get(nil, strconv.Itoa(random(0, 1000)), groupcache.AllocatingByteSliceSink(&buf))

	log.Println("Got", string(buf))
}

// handleInfo sends service information to client
func (s *server) handleInfo(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(struct {
		Version         string
		Uptime          string
		CachePeerLoads  int64
		CachePeerErrors int64
		ServerRequests  int64
		Gets            int64
	}{
		ServiceVersion,
		time.Since(ServiceStartTime).String(),
		int64(s.gcGroup.Stats.PeerLoads),
		int64(s.gcGroup.Stats.PeerErrors),
		int64(s.gcGroup.Stats.ServerRequests),
		int64(s.gcGroup.Stats.Gets),
	})
}
