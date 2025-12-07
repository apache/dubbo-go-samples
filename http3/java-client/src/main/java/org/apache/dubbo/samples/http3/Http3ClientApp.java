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

import greet.GreetRequest;
import greet.GreetResponse;
import greet.GreetService;

import org.apache.dubbo.config.ReferenceConfig;
import org.apache.dubbo.config.bootstrap.DubboBootstrap;

public class Http3ClientApp {

    public static void main(String[] args) throws Exception {
        System.out.println("Starting Dubbo Triple HTTP/3 Client...");

        ReferenceConfig<GreetService> referenceConfig = new ReferenceConfig<>();
        referenceConfig.setInterface(GreetService.class);
        referenceConfig.setUrl("tri://127.0.0.1:20000");
        referenceConfig.setProtocol("tri");

        DubboBootstrap bootstrap = DubboBootstrap.getInstance();
        bootstrap.reference(referenceConfig)
                .start();

        GreetService greetService = referenceConfig.get();

        System.out.println("Connecting to 127.0.0.1:20000/" + GreetService.SERVICE_NAME);

        GreetRequest request = GreetRequest.newBuilder()
                .setName("Java Client")
                .build();

        try {
            GreetResponse response = greetService.greet(request);
            System.out.println("SUCCESS! Response: " + response.getGreeting());
        } catch (Exception e) {
            System.err.println("Failed: " + e.getMessage());
            e.printStackTrace();
        }

        bootstrap.stop();
        System.exit(0);
    }
}
