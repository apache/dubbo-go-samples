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
package org.apache.dubbo.samples.tri.streaming;

import org.apache.dubbo.common.stream.StreamObserver;
import org.apache.dubbo.samples.tri.streaming.api.DubboGreetServiceTriple;
import org.apache.dubbo.samples.tri.streaming.api.GreetRequest;
import org.apache.dubbo.samples.tri.streaming.api.GreetResponse;
import org.apache.dubbo.samples.tri.streaming.api.GreetStreamRequest;
import org.apache.dubbo.samples.tri.streaming.api.GreetStreamResponse;
import org.apache.dubbo.samples.tri.streaming.api.GreetClientStreamRequest;
import org.apache.dubbo.samples.tri.streaming.api.GreetClientStreamResponse;
import org.apache.dubbo.samples.tri.streaming.api.GreetServerStreamRequest;
import org.apache.dubbo.samples.tri.streaming.api.GreetServerStreamResponse;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.ArrayList;
import java.util.List;


public class GreeterImpl extends DubboGreetServiceTriple.GreetServiceImplBase {
    private static final Logger LOGGER = LoggerFactory.getLogger(GreeterImpl.class);

    // Unary: 1 request -> 1 response, returns "Hello {name}"
    @Override
    public GreetResponse greet(GreetRequest request) {
        LOGGER.info("Received unary request: {}", request.getName());
        
        GreetResponse response = GreetResponse.newBuilder()
                .setGreeting("Hello " + request.getName())
                .build();
        
        LOGGER.info("Sending unary response: {}", response.getGreeting());
        return response;
    }

    // BiStream: N requests -> N responses (1:1), returns "Echo from biStream: {name}"
    @Override
    public StreamObserver<GreetStreamRequest> greetStream(StreamObserver<GreetStreamResponse> responseObserver) {
        return new StreamObserver<GreetStreamRequest>() {
            @Override
            public void onNext(GreetStreamRequest request) {
                try {
                    LOGGER.info("Received biStream request: {}", request.getName());
                    
                    // Echo back the request with the same name
                    GreetStreamResponse response = GreetStreamResponse.newBuilder()
                            .setGreeting("Echo from biStream: " + request.getName())
                            .build();
                    
                    responseObserver.onNext(response);
                } catch (Exception e) {
                    LOGGER.error("Error processing biStream request", e);
                    responseObserver.onError(e);
                }
            }

            @Override
            public void onError(Throwable throwable) {
                LOGGER.error("BiStream error occurred", throwable);
            }

            @Override
            public void onCompleted() {
                LOGGER.info("BiStream completed");
                responseObserver.onCompleted();
            }
        };
    }

    // ClientStream: N requests -> 1 response, returns "Received {count} names: {name1}, {name2}, ..."
    @Override
    public StreamObserver<GreetClientStreamRequest> greetClientStream(
            StreamObserver<GreetClientStreamResponse> responseObserver) {
        return new StreamObserver<GreetClientStreamRequest>() {
            private final List<String> names = new ArrayList<>();
            
            @Override
            public void onNext(GreetClientStreamRequest request) {
                LOGGER.info("Received clientStream request: {}", request.getName());
                names.add(request.getName());
            }
            
            @Override
            public void onError(Throwable throwable) {
                LOGGER.error("ClientStream error occurred", throwable);
            }
            
            @Override
            public void onCompleted() {
                // Aggregate all received names and send a single response
                String greeting = "Received " + names.size() + " names: " + String.join(", ", names);
                GreetClientStreamResponse response = GreetClientStreamResponse.newBuilder()
                        .setGreeting(greeting)
                        .build();
                
                responseObserver.onNext(response);
                responseObserver.onCompleted();
                LOGGER.info("ClientStream completed with {} requests", names.size());
            }
        };
    }

    // ServerStream: 1 request -> 10 responses, returns "Response {i} from serverStream for {name}"
    @Override
    public void greetServerStream(GreetServerStreamRequest request, StreamObserver<GreetServerStreamResponse> responseObserver) {
        try {
            LOGGER.info("Received serverStream request: {}", request.getName());
            
            for (int i = 0; i < 10; i++) {
                GreetServerStreamResponse response = GreetServerStreamResponse.newBuilder()
                        .setGreeting("Response " + i + " from serverStream for " + request.getName())
                        .build();
                
                responseObserver.onNext(response);
                LOGGER.debug("Sent serverStream response {}", i);
            }
            
            responseObserver.onCompleted();
            LOGGER.info("ServerStream completed for request: {}", request.getName());
            
        } catch (Exception e) {
            LOGGER.error("Error in serverStream", e);
            responseObserver.onError(e);
        }
    }
}
