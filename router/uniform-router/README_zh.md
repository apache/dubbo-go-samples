# Dubbo-go v3 统一路由规则

## 路由规则介绍
https://www.yuque.com/docs/share/c132d5db-0dcb-487f-8833-7c7732964bd4?# 

《微服务统一路由规则方案草案V2》

## 简介

路由规则，简单来说就是根据**特定的条件**，将**特定的请求**流量发送到**特定的服务提供者**。从而实现流量的分配。

在Dubbo3统一路由规则的定义中，需要提供两个yaml格式的资源：virtual service 和 destination rule。其格式和service mesh定义的路由规则非常相似。
- virtual service

定义host，用于和destination rule建立联系。\
定义 service 匹配规则\
定义 match 匹配规则\
匹配到特定请求后，进行目标集群的查找和验证，对于为空情况，使用fallback机制。

- destination rule

定义特定集群子集，以及子集所适配的标签，标签从provider 端暴露的 url 中获取，并尝试匹配。

## 提供能力
- 文件读入的路由配置
  
    [例子](./file/README_zh.md)
  
- 基于 K8S 动态更新的路由配置
  
    [例子](./k8s/README_zh.md)
