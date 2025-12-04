#!/bin/bash
mvn clean compile
mvn exec:java -Dexec.mainClass="org.apache.dubbo.samples.tri.streaming.StreamingClient"
