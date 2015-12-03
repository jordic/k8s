package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/template"

	aclient "k8s.io/kubernetes/pkg/api/unversioned"
	kclient "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/fields"
	"k8s.io/kubernetes/pkg/labels"
)

var (
	local = flag.Bool("local", false, "set to true if running on local machine")
	path  = flag.String("path", "/etc/nginx/conf.d", "Nginx conf output path")
)

var nginxTPL = `server {
    server_name {{.ServerName}};
    location / {
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header Host $http_host;
        proxy_set_header X-NginX-Proxy true;

        proxy_pass http://{{.Service}}:80;

    }
}`

func main() {
	flag.Parse()
	var cfg *kclient.Config
	var err error

	// If local flag is provided, will try to connect to the api, throught a
	// localhost proxy: you can get it up, with the following command:
	// kubectl proxy, and later launch the script: go run main.go
	if *local {
		cfg = &kclient.Config{
			Host: fmt.Sprintf("http://localhost:8001"),
		}
	} else {
		// This handles incluster config, when the script is runing inside
		// container, in a k8 cluster.
		cfg, err = kclient.InClusterConfig()
		if err != nil {
			log.Printf("failed to load incluster config %v", err)
			os.Exit(1)
		}
	}

	client, err := kclient.New(cfg)
	if err != nil {
		log.Printf("failed to create client %v", err)
	}
	// Query API for services, on the default namespace, matching
	// label proxy="true"
	services, err := client.Services("default").List(
		labels.SelectorFromSet(labels.Set{"proxy": "true"}),
		fields.Everything(), aclient.ListOptions{})

	// log.Printf("Services %v", services)
	// log.Printf("Name %s", services.Items[0].Name)
	// log.Printf("Proxy Name %s", services.Items[0].Labels["proxyName"])
	for _, k := range services.Items {
		// for every service found, write the desired nginx config file
		err = writeConfigFile(k.Name, k.Labels["proxyName"])
		if err != nil {
			log.Printf("can't write %s config file: %v\n", k.Name, err)
		}
	}
}

func writeConfigFile(service string, serverName string) error {

	log.Printf("writing service config file for %s, %s\n", service, serverName)
	fname := filepath.Join(*path, fmt.Sprintf("%s.conf", service))
	f, err := os.Create(fname)
	defer f.Close()
	if err != nil {
		return err
	}
	t, err := template.New("conf").Parse(nginxTPL)
	if err != nil {
		return err
	}
	w := bufio.NewWriter(f)
	err = t.Execute(w, map[string]interface{}{
		"ServerName": serverName,
		"Service":    service,
	})
	if err != nil {
		return err
	}
	w.Flush()
	return nil
}
