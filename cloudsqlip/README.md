CloudSQLIP
====

An app, that had to be deployed in a pod inside a kubernetes cluster, and 
monitors, nodes attached to it, adding new nodes to the gcloud sql allowed
hosts.

Running and customize:
=====

**TL;DR**

Enable access to Cloud SQL in your Container Engine cluster as well as enable cloudsql api in api console. In rc.yml, you have to add your cloudsql db server name. If you want to add additional ips allowed to access cloudsql instance, add as extra.

**Longer version:**

1. Make sure that the "Cloud SQL Enabled"  permission is set for your Google Container Engine cluster (you can check this by going to `Google Developers Console` &gt; `Container Engine` &gt; `Container clusters` &gt; &lt;cluster&gt; and check "Permissions")
1. Enable API Access control: Goto `Google Developers Console` &gt; `API Manager` and search for "Google Cloud SQL API". Click on `Google Cloud SQL API` and enable it (there's no need to add credentials)
1. Fork the github repository and customize the [rc.yml](rc.yml) file. Here's an example:

	```yaml
	apiVersion: v1
	kind: ReplicationController
	metadata:
	  labels:
	    name: cloudsqlip
	    version: "0.3"
	  name: cloudsqlip
	spec:
	  replicas: 1
	  selector: 
	    name: cloudsqlip
	    version: v1
	  template:
	    metadata:
	      labels:
	        name: cloudsqlip
	        version: v1
	    spec:
	      containers:
	      - name: cloudsqlip
	        image: jordic/cloudsqlip:0.3
	        resources:
	          limits:
	            cpu: 10m
	            memory: 50Mi
	        command: 
	        - /main
	        - -db
	        - cloudsqldbserver
	        - -extra
	        - 188.166.20.115/32
	```

	Replace `cloudsqldbserver` with the name of your Google Cloud SQL instance name (for example if the instance id is `my-project:my-server`then the name should be `my-server`). 
	Also note the parameter called `extra`. This is optional and it's a way to tell `cloudsqlip` that it should not only allow the Kubernetes cluster to access the database but also some other network. For example you might want to put the subnet of your workplace here if you want to access the database from work. If you don't need this just leave the last two lines out.
1. Now all we need to do is to deploy the replication controller to Kubernetes:

	```bash
	$ kubectl create -f rc.yml
	```

	You can see that it's up and running correctly by listing all pods (`kubectl get pods`):

	```bash
	NAME                               READY     STATUS    RESTARTS   AGE
	cloudsqlip-ysps9                   1/1       Running   0          2s
	....
	```

	Watching the logs (`kubectl logs cloudsqlip-ysps9`)  should give you something like this:

	```bash
	A 2015-12-03 20:55:05.000 Updated [https://www.googleapis.com/sql/v1beta3/projects/my-project/instances/my-cluster].
	A 2015-12-03 20:55:05.000 Patching Cloud SQL instance... Patching Cloud SQL instance.../ Patching Cloud SQL instance...- Patching Cloud SQL instance...done.
	A 2015-12-03 20:55:05.000 {"project": "my-project", "instance": "my-cluster", "settings": {"ipConfiguration": {"authorizedNetworks": ["xxx.xxx.xx.xx/32", "yyy.yyy.yy.yyy/32", "yyy.yyy.yy.yyy/32"]}}}
	A 2015-12-03 20:55:05.000 2015/12/03 19:55:05 Gcloud stderr The following message will be used for the patch API method.
	A 2015-12-03 20:55:05.000 2015/12/03 19:55:05 Gcloud stdout
	A 2015-12-03 20:55:00.000 2015/12/03 19:55:00 Patching sql network: xxx.xxx.xx.xx/32,yyy.yyy.yy.yyy/32,yyy.yyy.yy.yyy/32,
	A 2015-12-03 20:55:00.000 2015/12/03 19:55:00 Node list changed
	A 2015-12-03 20:55:00.000 2015/12/03 19:55:00 Polling node list
	```

And that's it! The `cloudsqlip` pod will now poll the list of cluster nodes every 10 seconds (this is configurable by using `-poll`) and maintain the allowed hosts for you!

Note that you only need to deploy one instance of `cloudsqlip` *per database* in your cluster.

Aditional params:
====

* `-poll` 10s interval seconds for the polling, kube api 
* `-local` for running it in localhost, using the kubectl -proxy