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

package org.apache.dubbo.samples.tls.provider;

import greet.Greet;
import greet.GreetService;
import greet.GreetServiceGrpc;
import io.grpc.stub.StreamObserver;

public class GreetServiceImpl extends GreetServiceGrpc.GreetServiceImplBase implements GreetService {
    
    /**
     * Dubbo RPC method - synchronous
     */
    @Override
    public Greet.GreetResponse greet(Greet.GreetRequest request) {
        return Greet.GreetResponse.newBuilder()
                .setGreeting(request.getName())
                .build();
    }
    
    /**
     * gRPC method - asynchronous (from GreetServiceImplBase)
     */
    @Override
    public void greet(Greet.GreetRequest request, StreamObserver<Greet.GreetResponse> responseObserver) {
        try {
            Greet.GreetResponse response = greet(request);  // Delegate to sync method
            responseObserver.onNext(response);
            responseObserver.onCompleted();
        } catch (Exception e) {
            // Send error to client to prevent hanging
            responseObserver.onError(e);
        }
    }
}
