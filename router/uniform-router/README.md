# Dubbo-go v3 Uniform router

## Introduction to routing rules
https://www.yuque.com/docs/share/c132d5db-0dcb-487f-8833-7c7732964bd4?#

"Draft Proposal for Uniform Routing Rules for Microservices V2"

## Introduction

Routing rules, in simple terms, are to send the traffic of a specific request to a specific service provider according to a specific condition. So as to realize the distribution of traffic.

In the definition of Dubbo3 unified routing rules, two resources in yaml format need to be provided: virtual service and destination rule. Its format is very similar to the routing rules defined by the service mesh.
-virtual service

Define host to establish contact with destination rule. \
Define service matching rules\
Define match matching rules\
After a specific request is matched, the target cluster is searched and verified. For empty cases, the fallback mechanism is used.

- destination rule

Define a specific cluster subset and the tags adapted to the subset. The tags are obtained from the URL exposed on the provider side and try to match.

## Ability provided
- Routing configuration for file reading

    [Example](./file/README.md)

- Routing configuration based on K8S dynamic update

    [Example](./k8s/README.md) 