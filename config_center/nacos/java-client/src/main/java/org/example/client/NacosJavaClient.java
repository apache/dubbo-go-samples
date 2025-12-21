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

package org.example.client;

import com.alibaba.nacos.api.NacosFactory;
import com.alibaba.nacos.api.config.ConfigService;
import greet.GreetService;
import greet.GreetRequest;
import greet.GreetResponse;
import org.apache.dubbo.config.ConfigCenterConfig;
import org.apache.dubbo.config.ReferenceConfig;
import org.apache.dubbo.config.bootstrap.DubboBootstrap;
import java.util.Properties;

public class NacosJavaClient {
    // Nacos server address (used as config center)
    private static final String NACOS_ADDR = "127.0.0.1:8848";
    // Zookeeper address (used as registry)
    private static final String ZK_ADDR = "127.0.0.1:2181";
    // DataId used in Nacos config center
    private static final String DATA_ID = "dubbo-go-samples-configcenter-nacos-java-client";
    // Nacos group
    private static final String GROUP = "dubbo";
    // Configuration content pushed to Nacos config center
    // Dubbo will load these as consumer / registry / protocol configs
    private static final String CONFIG_CONTENT =
            "dubbo.application.name=dubbo-java-consumer\n" +
                    "dubbo.registry.address=zookeeper://" + ZK_ADDR + "\n" +
                    "dubbo.protocol.name=tri\n" +
                    "dubbo.consumer.timeout=10000\n";

    public static void main(String[] args) throws Exception {
        //Publish Dubbo config to Nacos config center
        Properties properties = new Properties();
        properties.put("serverAddr", NACOS_ADDR);
        ConfigService configService = NacosFactory.createConfigService(properties);
        boolean success = configService.publishConfig(DATA_ID, GROUP, CONFIG_CONTENT);
        if (success) {
            System.out.println(" Succeed to publish config to Nacos, DATA_ID= " + DATA_ID);
        } else {
            System.err.println(" Failed to publish config to Nacos ");
        }


        ConfigCenterConfig configCenter = new ConfigCenterConfig();
        configCenter.setAddress("nacos://" + NACOS_ADDR);
        configCenter.setConfigFile(DATA_ID);
        configCenter.setGroup(GROUP);

        // Create reference config, only bind interface here
        ReferenceConfig<GreetService> reference = new ReferenceConfig<>();
        reference.setInterface(GreetService.class);

        // Create reference config, only bind interface here
        DubboBootstrap bootstrap = DubboBootstrap.getInstance()
                .application("dubbo-java-consumer")
                .configCenter(configCenter)
                .reference(reference)
                .start();

        // Get remote service proxy
        GreetService greetService = reference.get();
        // Build request and call remote method
        GreetRequest request = GreetRequest.newBuilder()
                .setName("Hello, this is dubbo java client!")
                .build();

        GreetResponse response = greetService.greet(request);
        System.out.println("server response: " + response.getGreeting());

        System.exit(0);
    }
}