#!/bin/bash

cd "$(dirname "$0")"
mvn -q compile exec:java -Dexec.mainClass="org.apache.dubbo.samples.h3.H3ClientApp"
