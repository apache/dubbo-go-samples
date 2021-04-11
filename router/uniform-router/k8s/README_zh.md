# K8S集群内使用 Dubbbogo 统一路由规则 示例

## 使用前提
本地已安装docker 和 k8s集群，可通过kubectl命令控制集群。

## 使用方法
`sh build.sh`

该脚本会先创建dubbo-workplace的命名空间，将三个服务依次构建，部署CRD资源和注册中心，打包镜像，并将三个服务部署至K8s集群。

client端会首先从文件中读入路由规则，根据规则将所有流量打入server，没有流量流至server2。

之后尝试修改路由规则，将dest_rule.yml 的子集标签限制去掉：
```yaml
apiVersion: service.dubbo.apache.org/v1alpha1
kind: DestinationRule
metadata: { name: demo-route }
spec:
  host: demo
  subsets:
    - name: all
      labels: { ut: CENTER }
    - name: center
#      labels: { ut: other } # 注释掉center对应的标签要求
    - name: other
```
更新路由规则:

```shell
kubectl apply -f ./go-client/conf/dest_rule.yml -n dubbo-workplace
```

之后可同时在两个server中查看到请求流量


## 删除命名空间
```shell
kubectl delete ns dubbo-workplace
```
