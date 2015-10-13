

CloudSQLIP
====

An app, that had to be deployed in a pod inside a kubernetes cluster, and 
monitors, nodes attached to it, adding new nodes to the gcloud sql allowed
hosts.

Running and customize:

In rc.yml, you have to add your cloudsql db server name. If you want to add
additional ips allowed to access cloudsql instance, add as extra.

Ensure your gke cluster is created with a service account with access to
cloudsql db.

Ensure that you had enabled cloudsql api in api console.


If you suspect is not woriing, and for debuging it, you can use:

kubectl get pods (got the actual pod name)
kubectl logs podname




