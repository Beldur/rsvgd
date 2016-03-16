package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/golang/groupcache"
)

var (
	// ServiceStartTime holds the time when service was started
	ServiceStartTime time.Time
	// ServiceVersion holds build information
	ServiceVersion    string
	dnsPeerLookupName = flag.String("dnsPeerLookupName", "", "DNS SRV Lookup name to find peers.")
	dnsPeerServer     = flag.String("dnsPeerServer", "", "DNS Server to find peers.")
	bindAddress       = flag.String("bind", getLocalIP(), "Bind server to given ip address.")
	port              = flag.Int("port", 8080, "Port to listen on.")
)

func init() {
	rand.Seed(time.Now().Unix())
	ServiceStartTime = time.Now()
}

func main() {
	flag.Parse()

	server := newServer(*bindAddress, *port)

	if *dnsPeerLookupName != "" && *dnsPeerServer != "" {
		startDNSRefreshTicker(server, *dnsPeerLookupName, *dnsPeerServer)
	}

	server.start()
}

// startDNSRefreshTicker queries srv records from DNS server and updates the
// server's groupcache peer list.
func startDNSRefreshTicker(server *server, dnsPeerLookupName, dnsPeerServer string) {
	dnsRefreshChan := time.NewTicker(5 * time.Second)
	go func() {
		for {
			select {
			case <-dnsRefreshChan.C:
				peerList := getSRVHostListFromDNSServer(dnsPeerLookupName, dnsPeerServer)
				fmt.Println(peerList)
				server.SetPeers(peerList)
			}
		}
	}()
}

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

func cacheGetter(context groupcache.Context, key string, dest groupcache.Sink) error {
	log.Println("Calculating...", key)

	dest.SetString("Result")

	return nil
}
