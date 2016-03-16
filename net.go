package main

import (
	"fmt"
	"log"
	"net"

	"github.com/miekg/dns"
)

// getLocalIP returns the non loopback local IP of the host
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}

	return ""
}

// queryDNSServer executes a given query on DNS Server
func queryDNSServer(query, server string, queryType uint16) ([]dns.RR, error) {
	c := dns.Client{}
	m := dns.Msg{}
	if string(query[len(query)-1]) != "." {
		query = query + "."
	}
	m.SetQuestion(query, queryType)

	r, _, err := c.Exchange(&m, server)
	if err != nil {
		return nil, err
	}

	return r.Answer, nil
}

// getSRVHostListFromDNSServer executes SRV request for given query on DNS Server
// and then tries to lookup the ip addresses of each entry.
func getSRVHostListFromDNSServer(query, server string) []string {
	var result []string
	answer, err := queryDNSServer(query, server, dns.TypeSRV)
	if err != nil {
		log.Printf("Could not get SRV Records from DNS Server (%s)\n", err)
		return result
	}

	for _, ans := range answer {
		srvRecord := ans.(*dns.SRV)
		ipAnswer, err := queryDNSServer(srvRecord.Target, server, dns.TypeA)
		if err != nil {
			log.Printf("Could not get A Record for %s from DNS Server (%s)\n", srvRecord.Target, err)
			continue
		}
		aRecord := ipAnswer[0].(*dns.A)
		result = append(result, fmt.Sprintf("http://%v:%d", aRecord.A.String(), srvRecord.Port))
	}

	return result
}
