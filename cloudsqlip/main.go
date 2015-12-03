package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/jordic/k8s/cloudsqlip/Godeps/_workspace/src/k8s.io/kubernetes/pkg/api"
	aclient "github.com/jordic/k8s/cloudsqlip/Godeps/_workspace/src/k8s.io/kubernetes/pkg/api/unversioned"
	kclient "github.com/jordic/k8s/cloudsqlip/Godeps/_workspace/src/k8s.io/kubernetes/pkg/client/unversioned"
)

var (
	local        = flag.Bool("local", false, "set to true if running on local machine not within cluster")
	localPort    = flag.Int("localport", 8001, "port that kubectl proxy is running on (local must be true)")
	pollInterval = flag.Duration("poll", 10, "Interval in seconds for polling the api to discover new nodes")
	dbName       = flag.String("db", "", "Db instance to patch ips..")
	extraIP      = flag.String("extra", "", "Extra ip in the form.. 11.11.11.11/32")

	nodeList []string
	client   *kclient.Client
)

func main() {

	flag.Parse()

	if *dbName == "" {
		log.Println("Must provide a dbinstance name")
		os.Exit(1)
	}

	var (
		cfg *kclient.Config
		err error
	)
	if *local {
		cfg = &kclient.Config{Host: fmt.Sprintf("http://localhost:%d", *localPort)}
	} else {
		cfg, err = kclient.InClusterConfig()
		if err != nil {
			log.Printf("failed to load config: %v", err)
			os.Exit(1)
		}
	}

	client, err = kclient.New(cfg)
	nodeList = make([]string, 0)
	PollNodes()
	ticker := time.NewTicker(*pollInterval * time.Second)
	for {
		select {
		case <-ticker.C:
			PollNodes()
		}
	}

}

// PollNodes polls for new nodes added to cluster
func PollNodes() {

	log.Println("Polling node list")
	nodes, err := client.Nodes().List(aclient.ListOptions{})

	if err != nil {
		log.Printf("failed to get node list %v", err)
	}
	var tnodes []string
	for _, n := range nodes.Items {
		ip := getExternalIP(n.Status.Addresses)
		tnodes = append(tnodes, ip)
	}

	if changed(nodeList, tnodes) {
		log.Println("Node list changed")
		nodeList = tnodes
		updateDb()
	}

}

func updateDb() {
	networks := ""
	for _, k := range nodeList {
		networks += fmt.Sprintf("%s/32,", k)
	}
	networks += *extraIP
	log.Printf("Patching sql network: %s\n", networks)
	cmd := exec.Command("gcloud", "sql", "instances", "patch", *dbName, "--authorized-networks", networks)
	var out bytes.Buffer
	var eout bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &eout
	err := cmd.Run()
	if err != nil {
		log.Printf("Error executing command %v", err)
	}
	log.Printf("Gcloud stdout %s\n", out.String())
	log.Printf("Gcloud stderr %s/n", eout.String())

}

func changed(l, p []string) bool {
	if len(l) != len(p) {
		return true
	}
	for _, k := range l {
		found := false
		for _, j := range p {
			if k == j {
				found = true
			}
		}
		if found == false {
			return true
		}
	}
	return false
}

func getExternalIP(addr []api.NodeAddress) string {
	for _, el := range addr {
		if el.Type == "ExternalIP" {
			return el.Address
		}
	}
	return ""
}
