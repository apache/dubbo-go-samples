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

package org.apache.dubbo.sample;

import org.apache.dubbo.common.constants.CommonConstants;
import org.apache.dubbo.config.ApplicationConfig;
import org.apache.dubbo.config.ProtocolConfig;
import org.apache.dubbo.config.RegistryConfig;
import org.apache.dubbo.config.ServiceConfig;
import org.springframework.context.support.ClassPathXmlApplicationContext;

public class Provider {

    /**
     * To get ipv6 address to work, add
     * System.setProperty("java.net.preferIPv6Addresses", "true");
     * before running your application.
     */
    public static void main(String[] args) throws Exception {
        ClassPathXmlApplicationContext context = new ClassPathXmlApplicationContext(new String[]{"META-INF/spring/dubbo.provider.xml"});
        context.start();
        startComplexService();
        startWrapperArrayClassProvider();
        System.in.read(); // press any key to exit
    }

    public static void startComplexService() throws InterruptedException {
        ServiceConfig<ComplexProvider> service = new ServiceConfig<>();
        service.setInterface(ComplexProvider.class);
        service.setRef(new ComplexProviderImpl());
        service.setProtocol(new ProtocolConfig(CommonConstants.DUBBO_PROTOCOL, 20010));
        service.setApplication(new ApplicationConfig("user-info-server"));
        service.setRegistry(new RegistryConfig("zookeeper://127.0.0.1:2181"));
        service.export();
        System.out.println("dubbo service started");
    }

    public static void startWrapperArrayClassProvider() throws InterruptedException {
        ServiceConfig<WrapperArrayClassProvider> service = new ServiceConfig<>();
        service.setInterface(WrapperArrayClassProvider.class);
        service.setRef(new WrapperArrayClassProviderImpl());
        service.setProtocol(new ProtocolConfig(CommonConstants.DUBBO_PROTOCOL, 20010));
        service.setApplication(new ApplicationConfig("user-info-server"));
        service.setRegistry(new RegistryConfig("zookeeper://127.0.0.1:2181"));
        service.export();
        System.out.println("dubbo service started");
    }
}
