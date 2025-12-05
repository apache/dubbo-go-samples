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

import greet.GreetRequest;
import greet.GreetResponse;
import greet.GreetServiceGrpc;

import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;

public class JavaClientApp {

    public static void main(String[] args) {

        ManagedChannel channel = ManagedChannelBuilder
                .forAddress("127.0.0.1", 20000)
                .usePlaintext()
                .build();

        GreetServiceGrpc.GreetServiceBlockingStub client =
                GreetServiceGrpc.newBlockingStub(channel);

        // Direct call: 1 request -> 1 response
        GreetRequest req = GreetRequest.newBuilder()
                .setName("Java Client Dubbo")
                .build();

        GreetResponse resp = client.greet(req);

        System.out.println("Response: " + resp.getGreeting());

        channel.shutdown();
    }
}
