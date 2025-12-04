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
import org.apache.dubbo.config.ApplicationConfig;
import org.apache.dubbo.config.ReferenceConfig;
import org.apache.dubbo.config.bootstrap.DubboBootstrap;
import org.apache.dubbo.samples.tri.streaming.api.GreetService;
import org.apache.dubbo.samples.tri.streaming.api.GreetStreamRequest;
import org.apache.dubbo.samples.tri.streaming.api.GreetStreamResponse;
import org.apache.dubbo.samples.tri.streaming.api.GreetServerStreamRequest;
import org.apache.dubbo.samples.tri.streaming.api.GreetServerStreamResponse;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.concurrent.CountDownLatch;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicInteger;

public class StreamingClient {
    private static final Logger LOGGER = LoggerFactory.getLogger(StreamingClient.class);
    
    private static final String SERVER_URL = "tri://127.0.0.1:20000";
    private static final int TIMEOUT_SECONDS = 30;
    
    public static void main(String[] args) {
        System.out.println("\n" + "=".repeat(70));
        System.out.println("Starting Dubbo Streaming Client Tests");
        System.out.println("=".repeat(70));
        
        // Configure reference to the remote service
        ReferenceConfig<GreetService> reference = new ReferenceConfig<>();
        reference.setInterface(GreetService.class);
        reference.setUrl(SERVER_URL);
        reference.setTimeout(10000);
        
        // Start Dubbo Bootstrap
        DubboBootstrap bootstrap = DubboBootstrap.getInstance();
        bootstrap.application(new ApplicationConfig("streaming-client"))
                .reference(reference)
                .start();
        
        System.out.println("Connected to server: " + SERVER_URL);
        System.out.println("=".repeat(70) + "\n");
        
        // Get the service proxy
        GreetService greeter = reference.get();
        
        // Test results tracking
        boolean bidiStreamSuccess = false;
        boolean serverStreamSuccess = false;
        
        try {
            // Test bidirectional streaming
            System.out.println("\n" + "=".repeat(70));
            System.out.println("TEST 1: Bidirectional Streaming");
            System.out.println("=".repeat(70));
            bidiStreamSuccess = testBidiStream(greeter);
            
            // Test server streaming
            System.out.println("\n" + "=".repeat(70));
            System.out.println("TEST 2: Server Streaming");
            System.out.println("=".repeat(70));
            serverStreamSuccess = testServerStream(greeter);
            
        } catch (Exception e) {
            LOGGER.error("Error during testing", e);
        } finally {
            // Print test results summary
            printTestSummary(bidiStreamSuccess, serverStreamSuccess);
            
            // Shutdown
            bootstrap.stop();
            System.out.println("\n Client shutdown complete\n");
        }
    }
    
    private static boolean testBidiStream(GreetService greeter) {
        final CountDownLatch latch = new CountDownLatch(1);
        final AtomicInteger responseCount = new AtomicInteger(0);
        final AtomicInteger expectedResponses = new AtomicInteger(5);
        final boolean[] success = {false};
        
        try {
            // Create response observer
            StreamObserver<GreetStreamResponse> responseObserver = new StreamObserver<GreetStreamResponse>() {
                @Override
                public void onNext(GreetStreamResponse response) {
                    int count = responseCount.incrementAndGet();
                    System.out.println("  Received response #" + count + ": " + response.getGreeting());
                }
                
                @Override
                public void onError(Throwable throwable) {
                    System.err.println("  BiStream error: " + throwable.getMessage());
                    latch.countDown(); // Ensure latch is counted down on error
                }
                
                @Override
                public void onCompleted() {
                    System.out.println("\n   BiStream completed - Received " + responseCount.get() + " responses");
                    success[0] = (responseCount.get() == expectedResponses.get());
                    latch.countDown();
                }
            };
            
            // Get request observer
            StreamObserver<GreetStreamRequest> requestObserver = greeter.greetStream(responseObserver);
            
            // Send multiple requests
            for (int i = 0; i < expectedResponses.get(); i++) {
                GreetStreamRequest request = GreetStreamRequest.newBuilder()
                        .setName("Client-" + i)
                        .build();
                System.out.println("    Sending request #" + i + ": " + request.getName());
                requestObserver.onNext(request);
                
                // Small delay between requests
                Thread.sleep(100);
            }
            
            // Signal completion
            requestObserver.onCompleted();
            System.out.println("\n   All requests sent, waiting for responses...");
            
            // Wait for responses
            if (!latch.await(TIMEOUT_SECONDS, TimeUnit.SECONDS)) {
                System.err.println("   BiStream test timed out");
                return false;
            }
            
            return success[0];
            
        } catch (Exception e) {
            System.err.println("   BiStream test exception: " + e.getMessage());
            LOGGER.error("Error in testBidiStream", e);
            return false;
        }
    }
    
    private static boolean testServerStream(GreetService greeter) {
        final CountDownLatch latch = new CountDownLatch(1);
        final AtomicInteger responseCount = new AtomicInteger(0);
        final int expectedResponses = 10;
        final boolean[] success = {false};
        
        try {
            // Create request
            GreetServerStreamRequest request = GreetServerStreamRequest.newBuilder()
                    .setName("StreamingClient")
                    .build();
            
            System.out.println("    Sending request: " + request.getName());
            System.out.println("   Waiting for server stream responses...\n");
            
            // Create response observer
            StreamObserver<GreetServerStreamResponse> responseObserver = new StreamObserver<GreetServerStreamResponse>() {
                @Override
                public void onNext(GreetServerStreamResponse response) {
                    int count = responseCount.incrementAndGet();
                    System.out.println("  ⬅️  Received response #" + count + ": " + response.getGreeting());
                }
                
                @Override
                public void onError(Throwable throwable) {
                    System.err.println("   ServerStream error: " + throwable.getMessage());
                    latch.countDown(); // Ensure latch is counted down on error
                }
                
                @Override
                public void onCompleted() {
                    System.out.println("\n   ServerStream completed - Received " + responseCount.get() + " responses");
                    success[0] = (responseCount.get() == expectedResponses);
                    latch.countDown();
                }
            };
            
            // Call server streaming method
            greeter.greetServerStream(request, responseObserver);
            
            // Wait for responses
            if (!latch.await(TIMEOUT_SECONDS, TimeUnit.SECONDS)) {
                System.err.println("   ServerStream test timed out");
                return false;
            }
            
            return success[0];
            
        } catch (Exception e) {
            System.err.println("   ServerStream test exception: " + e.getMessage());
            LOGGER.error("Error in testServerStream", e);
            return false;
        }
    }
    
    /**
     * Print test results summary.
     * 
     * Requirements: 3.4
     * 
     * @param bidiStreamSuccess whether bidirectional streaming test passed
     * @param serverStreamSuccess whether server streaming test passed
     */
    private static void printTestSummary(boolean bidiStreamSuccess, boolean serverStreamSuccess) {
        System.out.println("\n" + "=".repeat(70));
        System.out.println(" TEST RESULTS SUMMARY");
        System.out.println("=".repeat(70));
        System.out.println("  Bidirectional Streaming: " + (bidiStreamSuccess ? " PASSED" : " FAILED"));
        System.out.println("  Server Streaming: " + (serverStreamSuccess ? " PASSED" : " FAILED"));
        System.out.println("-".repeat(70));
        if (bidiStreamSuccess && serverStreamSuccess) {
            System.out.println("   Overall: ALL TESTS PASSED!");
        } else {
            System.out.println("    Overall: SOME TESTS FAILED");
        }
        System.out.println("=".repeat(70));
    }
}
