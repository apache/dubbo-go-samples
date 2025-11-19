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
import greet.GreetServiceGrpc;

import io.grpc.Server;
import io.grpc.ServerBuilder;
import io.grpc.stub.StreamObserver;

public class JavaServerApp {

    public static void main(String[] args) throws Exception {

        Server server = ServerBuilder.forPort(20000)
                .addService(new GreetServiceGrpc.GreetServiceImplBase() {

                    @Override
                    public void greet(GreetRequest request,
                                      StreamObserver<GreetResponse> respObserver) {

                        String name = request.getName();

                        GreetResponse resp = GreetResponse.newBuilder()
                                .setGreeting("Hello from Java Server, " + name)
                                .build();

                        respObserver.onNext(resp);
                        respObserver.onCompleted();
                    }
                })
                .build();

        server.start();
        System.out.println("Java Triple Server started on port 20000");
        server.awaitTermination();
    }
}
