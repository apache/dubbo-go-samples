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
import org.apache.dubbo.config.ApplicationConfig;
import org.apache.dubbo.config.RegistryConfig;
import org.apache.dubbo.config.SslConfig;
import org.apache.dubbo.config.bootstrap.DubboBootstrap;

import javax.net.ssl.SSLException;

public class TlsTriProtoConsumer {
    public static void main(String[] args) throws SSLException {
        String caCertPath;
        if (args.length > 0) {
            caCertPath = args[0];
        } else {
            String root = System.getProperty("user.dir");
            caCertPath = root + "/../x509/server_ca_cert.pem";
        }
        java.io.File caFile = new java.io.File(caCertPath);
        if (!caFile.isAbsolute()) {
            caFile = caFile.getAbsoluteFile();
        }
        if (!caFile.exists()) {
            throw new IllegalArgumentException("CA cert not found: " + caFile);
        }
        caCertPath = caFile.getPath();

        String host = System.getProperty("tls.host", "127.0.0.1");
        String authority = System.getProperty("tls.authority", "dubbogo.test.example.com");
        int port = Integer.getInteger("tls.port", 20000);

        // Dubbo lifecycle and TLS config
        SslConfig ssl = new SslConfig();
        ssl.setClientTrustCertCollectionPath(caCertPath);

        DubboBootstrap bootstrap = DubboBootstrap.getInstance();
        bootstrap.application(new ApplicationConfig("java-tri-ssl-proto-consumer"))
                .registry(new RegistryConfig(RegistryConfig.NO_AVAILABLE))
                .ssl(ssl)
                .start();

        // Bridge client using gRPC stub under the hood
        GreetClientImpl client = new GreetClientImpl(host, port, authority, caCertPath);
        try {
            Greet.GreetResponse resp = client.greet(
                    Greet.GreetRequest.newBuilder().setName("hello world").build()
            );
            System.out.println("Greet response: " + resp.getGreeting());
        } finally {
            client.shutdown();
            bootstrap.destroy();
        }
    }
}
