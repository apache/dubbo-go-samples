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
package org.apache.dubbo.samples.http3;

import greet.DubboGreetServiceTriple;
import greet.GreetRequest;
import greet.GreetResponse;
import greet.GreetService;

import org.apache.dubbo.config.ServiceConfig;
import org.apache.dubbo.config.bootstrap.DubboBootstrap;

import java.util.concurrent.CompletableFuture;

public class Http3ServerApp {

    public static void main(String[] args) throws Exception {
        System.out.println("Starting Dubbo Triple HTTP/3 Server...");

        ServiceConfig<GreetService> serviceConfig = new ServiceConfig<>();
        serviceConfig.setInterface(GreetService.class);
        serviceConfig.setRef(new GreetServiceImpl());

        DubboBootstrap bootstrap = DubboBootstrap.getInstance();
        bootstrap.service(serviceConfig)
                .start();

        System.out.println("Server started on port 20000 - " + GreetService.SERVICE_NAME);
        bootstrap.await();
    }

    public static class GreetServiceImpl extends DubboGreetServiceTriple.GreetServiceImplBase {
        @Override
        public GreetResponse greet(GreetRequest request) {
            System.out.println("Received: " + request.getName());
            return GreetResponse.newBuilder()
                    .setGreeting("Hello " + request.getName() + ", from Java HTTP/3!")
                    .build();
        }
    }
}
