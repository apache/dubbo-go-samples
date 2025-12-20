#!/bin/bash

cd "$(dirname "$0")"
mvn compile exec:java -Dexec.mainClass="org.apache.dubbo.samples.h3.H3ClientApp"
