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

import greet.GreetServiceGrpc;
import org.apache.dubbo.config.ApplicationConfig;
import org.apache.dubbo.config.ProtocolConfig;
import org.apache.dubbo.config.ServiceConfig;
import org.apache.dubbo.config.SslConfig;
import org.apache.dubbo.config.bootstrap.DubboBootstrap;

import java.util.concurrent.CountDownLatch;

public class TlsTriProvider {
    public static void main(String[] args) throws Exception {
        // Allow overriding cert paths via args: certChain, privateKey, [trustChain]
        SslConfig ssl = new SslConfig();
        String privateKeyPath;
        
        if (args.length > 0) {
            if (args.length < 2 || args.length > 3) {
                System.out.println("USAGE: TlsTriProvider certChainFilePath privateKeyFilePath [trustCertCollectionFilePath]");
                System.exit(1);
            }
            ssl.setServerKeyCertChainPath(args[0]);
            privateKeyPath = args[1];
            if (args.length == 3) {
                ssl.setServerTrustCertCollectionPath(args[2]);
            }
        } else {
            // default relative to this module: reuse go sample certs
            String root = System.getProperty("user.dir");
            // user.dir is typically .../java-tri-ssl/provider
            String base = root + "/../x509";
            ssl.setServerKeyCertChainPath(base + "/server2_cert.pem");
            privateKeyPath = base + "/server2_key.pem";
        }
        
        // Convert PKCS#1 to PKCS#8 if necessary with proper error handling
        try {
            privateKeyPath = Pkcs1ToPkcs8KeyConverter.loadAndConvertKey(privateKeyPath);
            ssl.setServerPrivateKeyPath(privateKeyPath);
        } catch (Exception e) {
            System.err.println("[TLS ERROR] Failed to load or convert private key: " + e.getMessage());
            System.err.println("Please ensure:");
            System.err.println("  1. The private key file exists and is readable");
            System.err.println("  2. The key is in valid PKCS#1 or PKCS#8 format");
            System.err.println("  3. OpenSSL is installed and available in PATH (required for PKCS#1 conversion)");
            throw new RuntimeException("TLS configuration failed", e);
        }

        ProtocolConfig protocol = new ProtocolConfig();
        protocol.setName("tri");
        protocol.setPort(20000);
        protocol.setSslEnabled(true);

        // Use the Dubbo interface from greet package (matches proto package.service)
        ServiceConfig<greet.GreetService> service = new ServiceConfig<>();
        service.setInterface(greet.GreetService.class);
        service.setRef(new GreetServiceImpl());
        service.setRegister(false); // no registry, direct export

        DubboBootstrap bootstrap = DubboBootstrap.getInstance();
        bootstrap.application(new ApplicationConfig("java-tri-ssl-provider"))
                .protocol(protocol)
                .ssl(ssl)
                .service(service)
                .start();

        // Add graceful shutdown hook
        Runtime.getRuntime().addShutdownHook(new Thread(() -> {
            System.out.println("[TLS] Shutting down Java triple TLS provider...");
            try {
                bootstrap.stop();
                System.out.println("[TLS] Dubbo provider stopped gracefully");
            } catch (Exception e) {
                System.err.println("[TLS] Error during shutdown: " + e.getMessage());
            }
        }));

        System.out.println("Java triple TLS provider started on tri://0.0.0.0:20000");
        System.out.println("Press Ctrl+C to stop the server");
        
        // Keep the application running and handle interruptions gracefully
        try {
            Thread.currentThread().join();
        } catch (InterruptedException e) {
            System.out.println("[TLS] Server interrupted, shutting down...");
            Thread.currentThread().interrupt();
        }
    }
}
