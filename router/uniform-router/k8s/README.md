# Example of using Dubbbogo unified routing rules in the K8S cluster

## Premise
Docker and k8s clusters have been installed locally, and the clusters can be controlled through kubectl.

## Instructions
`sh build.sh`

The script will first create the dubbo-workplace namespace, construct the three services in turn, deploy the CRD resources and registry, package the image, and deploy the three services to the K8s cluster.

The client side will first read the routing rules from the file, and send all traffic to the server according to the rules, and no traffic will flow to server2.

Then try to modify the routing rules to remove the subset label restriction of dest_rule.yml:
```yaml
apiVersion: service.dubbo.apache.org/v1alpha1
kind: DestinationRule
metadata: {name: demo-route}
spec:
  host: demo
  subsets:
    -name: all
      labels: {ut: CENTER}
    -name: center
# labels: {ut: other} # Comment out the label requirements corresponding to center
    -name: other
```
Update routing rules:

```shell
kubectl apply -f ./go-client/conf/dest_rule.yml -n dubbo-workplace
```

After that, the request traffic can be viewed in both servers at the same time


## Delete namespace
```shell
kubectl delete ns dubbo-workplace
``` 