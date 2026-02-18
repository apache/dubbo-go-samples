/*
 *  Licensed to the Apache Software Foundation (ASF) under one or more
 *  contributor license agreements.  See the NOTICE file distributed with
 *  this work for additional information regarding copyright ownership.
 *  The ASF licenses this file to You under the Apache License, Version 2.0
 *  (the "License"); you may not use this file except in compliance with
 *  the License.  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package org.apache.dubbo.samples;

import org.apache.dubbo.config.ProtocolConfig;
import org.apache.dubbo.config.RegistryConfig;
import org.apache.dubbo.config.ServiceConfig;
import org.apache.dubbo.config.bootstrap.DubboBootstrap;

import java.util.concurrent.CountDownLatch;

public class ApiProvider {
    private static final int DUBBO_PORT = 20000;
    private static final int TRIPLE_PORT = 50052;
    private static final String SERVICE_VERSION = "1.0.0";

    public static void main(String[] args) throws InterruptedException {
        // Dubbo protocol service (group=dubbo)
        ServiceConfig<UserProvider> dubboService = new ServiceConfig<>();
        dubboService.setInterface(UserProvider.class);
        dubboService.setRef(new UserProviderImpl());
        dubboService.setGroup("dubbo");
        dubboService.setVersion(SERVICE_VERSION);
        dubboService.setProtocol(new ProtocolConfig("dubbo", DUBBO_PORT));

        // Triple protocol service (group=triple)
        ServiceConfig<UserProvider> tripleService = new ServiceConfig<>();
        tripleService.setInterface(UserProvider.class);
        tripleService.setRef(new UserProviderImpl());
        tripleService.setGroup("triple");
        tripleService.setVersion(SERVICE_VERSION);
        tripleService.setProtocol(new ProtocolConfig("tri", TRIPLE_PORT));

        // Use direct export without registry
        RegistryConfig registryConfig = new RegistryConfig();
        registryConfig.setAddress("N/A");

        DubboBootstrap bootstrap = DubboBootstrap.getInstance();
        bootstrap.application("generic-java-server")
                .registry(registryConfig)
                .service(dubboService)
                .service(tripleService)
                .start();

        System.out.println("Generic Java server started:");
        System.out.println("  - Dubbo protocol on port " + DUBBO_PORT + " (group=dubbo)");
        System.out.println("  - Triple protocol on port " + TRIPLE_PORT + " (group=triple)");

        new CountDownLatch(1).await();
    }
}
