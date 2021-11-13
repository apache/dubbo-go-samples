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

package org.apache.dubbo;

import org.apache.dubbo.config.ApplicationConfig;
import org.apache.dubbo.config.ReferenceConfig;
import org.apache.dubbo.config.RegistryConfig;
import org.apache.dubbo.rpc.service.GenericService;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class ApiConsumer {
    private static final Logger logger = LoggerFactory.getLogger("userLogger"); // Output to com.dubbogo.user-server.log

    private static GenericService genericService;

    public static void main(String[] args) {
        initConfig();
        callGetUser();
    }

    private static void initConfig() {
        logger.info("\n\n\nstart to init config\n\n\n");
        ApplicationConfig applicationConfig = new ApplicationConfig();
        ReferenceConfig<GenericService> reference = new ReferenceConfig<>();
        applicationConfig.setName("user-info-server");
        reference.setApplication(applicationConfig);
        RegistryConfig registryConfig = new RegistryConfig();
        registryConfig.setAddress("zookeeper://127.0.0.1:2181");
        reference.setRegistry(registryConfig);
        reference.setGeneric(true);
        reference.setInterface("org.apache.dubbo.UserProvider");
        genericService = reference.get();
    }

    private static void callGetUser() {
        logger.info("\n\n\nCall GetUser");
        logger.info("Start to generic invoke");
        Object result;

        result = genericService.$invoke("GetUser", new String[]{"java.lang.String", "java.lang.String"}, new Object[]{"A003", "lily"});
        logger.info("\n\n\n" + "GetUser(String userId, String name) " + "res: " + result + "\n\n\n");
    }
}
