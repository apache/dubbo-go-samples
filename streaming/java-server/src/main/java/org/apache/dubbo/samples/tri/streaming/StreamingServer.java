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

import org.apache.dubbo.config.ApplicationConfig;
import org.apache.dubbo.config.ProtocolConfig;
import org.apache.dubbo.config.ServiceConfig;
import org.apache.dubbo.config.bootstrap.DubboBootstrap;
import org.apache.dubbo.samples.tri.streaming.api.GreetService;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class StreamingServer {
    private static final Logger LOGGER = LoggerFactory.getLogger(StreamingServer.class);
    
    private static final int SERVER_PORT = 20000;
    private static final String PROTOCOL_NAME = "tri";
    private static final String APPLICATION_NAME = "streaming-server";

    public static void main(String[] args) {
        try {
            LOGGER.info("Starting Dubbo Streaming Server...");
            
            // Create service configuration
            ServiceConfig<GreetService> service = new ServiceConfig<>();
            service.setInterface(GreetService.class);
            service.setRef(new GreeterImpl());
            
            // Configure and start Dubbo Bootstrap
            DubboBootstrap bootstrap = DubboBootstrap.getInstance();
            bootstrap.application(new ApplicationConfig(APPLICATION_NAME))
                    .protocol(new ProtocolConfig(PROTOCOL_NAME, SERVER_PORT))
                    .service(service)
                    .start();
            
            LOGGER.info("Dubbo Streaming Server started successfully on port {}", SERVER_PORT);
            LOGGER.info("Server is ready to accept connections...");
            
            // Keep the server running
            bootstrap.await();
            
        } catch (Exception e) {
            LOGGER.error("Failed to start Dubbo Streaming Server", e);
            System.exit(1);
        }
    }
}
