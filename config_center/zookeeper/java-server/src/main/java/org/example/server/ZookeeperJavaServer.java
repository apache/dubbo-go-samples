/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package org.example.server;

import greet.GreetRequest;
import greet.GreetResponse;
import greet.GreetService;
import org.apache.curator.framework.CuratorFramework;
import org.apache.curator.framework.CuratorFrameworkFactory;
import org.apache.curator.retry.ExponentialBackoffRetry;
import org.apache.dubbo.config.ConfigCenterConfig;
import org.apache.dubbo.config.ServiceConfig;
import org.apache.dubbo.config.bootstrap.DubboBootstrap;

import java.nio.charset.StandardCharsets;

public class ZookeeperJavaServer {
    // Zookeeper address (used as both registry and config center)
    private static final String ZK_ADDR = "127.0.0.1:2181";
    // DataId used in Zookeeper config center
    private static final String DATA_ID = "dubbo-go-samples-configcenter-zookeeper-java-server";
    // Config center group
    private static final String GROUP = "dubbogo";
    // Configuration content pushed to Zookeeper config center
    // Dubbo will load these as application / registry / protocol configs
    private static final String CONFIG_CONTENT =
            "dubbo.application.name=dubbo-java-provider\n" +
                    "dubbo.registry.address=zookeeper://" + ZK_ADDR + "\n" +
                    "dubbo.protocol.name=tri\n" +
                    "dubbo.protocol.port=50051\n";
    public static void main(String[] args) throws Exception {
        // Publish Dubbo config to Zookeeper config center
        String configPath = "/dubbo/config/" + GROUP + "/" + DATA_ID;
        CuratorFramework client = CuratorFrameworkFactory
                .builder()
                .connectString(ZK_ADDR)
                .retryPolicy(new ExponentialBackoffRetry(1000, 3))
                .build();
        client.start();

        try {
            if (client.checkExists().forPath(configPath) == null) {
                client.create()
                        .creatingParentsIfNeeded()
                        .forPath(configPath, CONFIG_CONTENT.getBytes(StandardCharsets.UTF_8));
                System.out.println("Config created in ZK, path=" + configPath);
            } else {
                client.setData()
                        .forPath(configPath, CONFIG_CONTENT.getBytes(StandardCharsets.UTF_8));
                System.out.println("Config updated in ZK, path=" + configPath);
            }
        } finally {
            client.close();
        }

        // Tell Dubbo to use Zookeeper as config center
        ConfigCenterConfig configCenter = new ConfigCenterConfig();
        configCenter.setAddress("zookeeper://" + ZK_ADDR);
        configCenter.setGroup(GROUP);
        configCenter.setConfigFile(DATA_ID);

        // Define service: only interface and implementation
        ServiceConfig<GreetService> service = new ServiceConfig<>();
        service.setInterface(GreetService.class);
        service.setRef(new GreetServiceImpl());

        // Start DubboBootstrap and load reference from config center
        DubboBootstrap.getInstance()
                .application("dubbo-java-provider")
                .configCenter(configCenter)
                .service(service)
                .start()
                .await();
    }

    // Service implementation of GreetService
    static class GreetServiceImpl implements GreetService {
        @Override
        public GreetResponse greet(GreetRequest request) {
            System.out.println("Received request: " + request.getName());
            return GreetResponse.newBuilder()
                    .setGreeting("Hello, this is dubbo java server! I received: " + request.getName())
                    .build();
        }
    }
}