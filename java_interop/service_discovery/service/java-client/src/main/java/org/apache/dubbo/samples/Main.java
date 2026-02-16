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

import org.apache.dubbo.config.ReferenceConfig;
import org.apache.dubbo.config.RegistryConfig;
import org.apache.dubbo.config.bootstrap.DubboBootstrap;
import org.apache.dubbo.samples.proto.GreetRequest;
import org.apache.dubbo.samples.proto.GreetResponse;
import org.apache.dubbo.samples.proto.GreetService;

public class Main {
    public static void main(String[] args) {
        System.setProperty("dubbo.application.service-discovery.migration", "APPLICATION_FIRST");

        ReferenceConfig<GreetService> reference = new ReferenceConfig<>();
        reference.setInterface(GreetService.class);

        DubboBootstrap bootstrap = DubboBootstrap.getInstance()
                .application("greet-java-client")
                .registry(new RegistryConfig("nacos://127.0.0.1:8848"))
                .reference(reference)
                .start();

        GreetService greetService = reference.get();
        GreetRequest req = GreetRequest.newBuilder().setName("Mamba").build();
        System.out.println("dubbo ref started");
        GreetResponse greet = greetService.greet(req);
        System.out.println("Greeting:" + greet.getGreeting() + "!!!!!!!!!!!!");

        bootstrap.stop();
    }
}
