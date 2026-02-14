/*
 *  Licensed to the Apache Software Foundation (ASF) under one or more
 *  contributor license agreements.  See the NOTICE file distributed with
 *  this work for additional information regarding copyright ownership.
 *  The ASF licenses this file to You under the Apache License, Version 2.0
 *  (the "License"); you may not use this file except in compliance with
 *  the License.  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package org.apache.dubbo.samples;

import org.apache.dubbo.config.ProtocolConfig;
import org.apache.dubbo.config.RegistryConfig;
import org.apache.dubbo.config.ServiceConfig;
import org.apache.dubbo.config.bootstrap.DubboBootstrap;
import org.apache.dubbo.common.constants.CommonConstants;
import org.apache.dubbo.samples.proto.GreetService;
import org.apache.dubbo.samples.service.GreetServiceImpl;

public class Main {
    public static void main(String[] args) {
        System.setProperty("dubbo.application.register-mode", "instance");

        ServiceConfig<GreetService> service = new ServiceConfig<>();
        service.setInterface(GreetService.class);
        service.setRef(new GreetServiceImpl());

        DubboBootstrap.getInstance()
                .application("greet-java-server")
                .registry(new RegistryConfig("nacos://127.0.0.1:8848"))
                .protocol(new ProtocolConfig(CommonConstants.TRIPLE, 20055))
                .service(service)
                .start()
                .await();
    }
}
