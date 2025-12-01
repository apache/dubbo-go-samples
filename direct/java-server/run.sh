#!/bin/bash
mvn -q clean package
mvn -q exec:java -Dexec.mainClass=org.example.server.JavaGreetServer
