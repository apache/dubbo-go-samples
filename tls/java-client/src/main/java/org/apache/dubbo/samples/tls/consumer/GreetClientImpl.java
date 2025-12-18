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

package org.apache.dubbo.samples.tls.consumer;

import greet.Greet;
import greet.GreetServiceGrpc;
import io.grpc.ManagedChannel;
import io.grpc.netty.shaded.io.grpc.netty.GrpcSslContexts;
import io.grpc.netty.shaded.io.grpc.netty.NettyChannelBuilder;
import io.grpc.netty.shaded.io.netty.handler.ssl.SslContext;

import javax.net.ssl.SSLException;
import java.io.File;

public class GreetClientImpl implements GreetClient {
    private final ManagedChannel channel;           // gRPC connection channel
    private final GreetServiceGrpc.GreetServiceBlockingStub stub;  // Synchronous RPC stub

    /**
     * Initialize TLS encrypted connection
     * @param targetHost Server hostname or IP address to connect to
     * @param port Server port number
     * @param authorityHost Hostname for TLS certificate verification (SNI)
     * @param caCertPath CA certificate path for server verification
     */
    public GreetClientImpl(String targetHost, int port, String authorityHost, String caCertPath) throws SSLException {
        // Configure SSL context and load CA certificate for server verification
        SslContext sslContext = GrpcSslContexts.forClient()
                .trustManager(new File(caCertPath))
                .build();

        // Create gRPC channel with TLS encryption enabled
        this.channel = NettyChannelBuilder.forAddress(targetHost, port)
                .sslContext(sslContext)
                .overrideAuthority(authorityHost)  // SNI: hostname used during certificate verification
                .build();
        
        // Create synchronous RPC service stub
        this.stub = GreetServiceGrpc.newBlockingStub(channel);
    }

    @Override
    public Greet.GreetResponse greet(Greet.GreetRequest request) {
        // Call remote service through TLS encrypted channel
        return stub.greet(request);
    }

    public void shutdown() {
        channel.shutdown();  // Close connection and release resources
    }
}
