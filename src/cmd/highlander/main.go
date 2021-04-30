package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sort"
	"time"
)

var cmd Cmd
var version = "0.1.1"

// like httputil.NewSingleHostReverseProxy but with BasicAuth and Host header rewrite
func NewSingleHostAuthReverseProxy(target *url.URL, user string, pass string) *httputil.ReverseProxy {
	targetQuery := target.RawQuery
	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.RawQuery = targetQuery + req.URL.RawQuery
		if _, ok := req.Header["User-Agent"]; !ok {
			req.Header.Set("User-Agent", "") // no default agent
		}

		// extensions to NewSingleHostReverseProxy:
		req.SetBasicAuth(user, pass)
		req.Host = target.Host
	}
	return &httputil.ReverseProxy{Director: director}
}

func main() {
	cmd = parseCmd()

	log.Println("Highander v", version, " is now starting. (use -h for help)")
	log.Println("binding address: ", cmd.bind)
	log.Println("remote address: ", cmd.remote)
	if cmd.remoteUser != "" {
		log.Println("remote user", cmd.remoteUser)
	}
	log.Println("health tick: ", cmd.healthTick.String())
	log.Println("health expiry: ", cmd.healthExpiry.String())
	log.Println("instance ID: ", cmd.instanceID)
	if cmd.clusterConnect != "" {
		log.Println("cluster connection string : ", cmd.clusterConnect, " (clustered)")
	} else {
		log.Println("cluster connection string : N/A (standalone)")
	}
	log.Println("cluster binding string: ", cmd.clusterBind)

	u, err := url.Parse(cmd.remote)
	if err != nil {
		panic(err)
	}

	// now connect to the cluster
	cluster, err := clusterJoin()
	if err != nil {
		panic(err)
	}

	// Ask for members of the cluster
	for _, member := range cluster.Members() {
		log.Printf("Cluster member: %s %s:%d\n", member.Name, member.Addr, member.Port)
	}

	var isLeader = false
	go func() {
		for {
			nodes := cluster.Members()
			sort.Slice(nodes, func(i, j int) bool {
				return nodes[i].Name < nodes[j].Name
			})
			if isLeader != (nodes[0] == cluster.LocalNode()) {
				isLeader = !isLeader
				log.Printf("isLeader changed to %t\n", isLeader)
			}
			time.Sleep(time.Second)
		}
	}()

	// then bind the service
	var proxy *httputil.ReverseProxy
	if cmd.remoteUser != "" {
		log.Printf("using authentication for user %s", cmd.remoteUser)
		proxy = NewSingleHostAuthReverseProxy(u, cmd.remoteUser, cmd.remotePassword)
	} else {
		proxy = httputil.NewSingleHostReverseProxy(u)
	}
	proxy.Transport = NewHighlanderProxy(
		cmd.healthTick,
		cmd.healthExpiry,
		func() bool {
			return isLeader
		})
	err = http.ListenAndServe(cmd.bind, proxy)
	if err != nil {
		panic(err)
	}

}
