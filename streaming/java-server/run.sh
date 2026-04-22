#!/bin/bash
mvn -q clean compile
mvn -q exec:java -Dexec.mainClass="org.apache.dubbo.samples.tri.streaming.StreamingServer"
