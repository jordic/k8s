
### NGiNx Service Proxy for Kubernetes

A docker image, and a Kubernetes resource controller and service, for building
a nginx proxy that proxies labeled services. 

The docker image is public at jordic/nginxsp

For using it inside a kubernetes cluster, first of all, you need to create the
service with something like:

```
kubectl -f create service.yaml
```
In the current yaml file, the service exposes a fixed NodePort, that can be
accessed on any public ip address of the cluster. (You can put outside a LB)

Later you need to label, each desired service with these labels:

proxy="true"
proxyName="subdomain.domain.com"

And create the resource controller:

```
kubectl -f create rc.yaml
```

On container start, the entry point, will launch the main.go, this will query
the API, for services matching label: proxy="true" and build the nginx conf
files, on the directory /etc/nginx/conf.d 

As an example. for a service named frontend-test it will create the file:
/etc/nginx/conf.d/frontend-test.conf 
```
server {
    server_name subdomain.domain.com;
    location / {
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header Host $http_host;
        proxy_set_header X-NginX-Proxy true;
        proxy_pass http://frontend-test:80;

    }
}
```

When done, will start the nginx proxy, with the available services.

Requirements: You will need a cluster with the kube-dns service.


#### Motivation

The main motivation for building this image, is a simple use case to proxy
some services that matches a label. I know there are plenty of other options,
and also the Ingress controller exists, but, for now, I prefer to play, and
learn the kubernetes API with simple and easy use cases.

Also, as said, this is a simple example of using the Kubernetes API, for querying 
available services by label from a go program.

You can just build it with:

```
go get -u github.com/jordic/k8s/nginx_proxy
```


#### Improvements

- Handle ssl termination on nginx-proxy. Actually, we have the SSL
    termination, at LB level.
- Think on a better way of handling service reloads.... Perhaps with
    rolling-update.
- Rebuild the service, using the watch API, and dinamicaly proxying to
    kubernetes services.



