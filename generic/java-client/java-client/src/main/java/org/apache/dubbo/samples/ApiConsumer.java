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

package org.apache.dubbo.samples;

import org.apache.dubbo.config.ProtocolConfig;
import org.apache.dubbo.config.ReferenceConfig;
import org.apache.dubbo.config.RegistryConfig;
import org.apache.dubbo.config.bootstrap.DubboBootstrap;

public class ApiConsumer {
    public static void main(String[] args) {
        DubboBootstrap bootstrap = DubboBootstrap.getInstance();

        RegistryConfig registryConfig = new RegistryConfig("zookeeper://127.0.0.1:2181");
        ProtocolConfig protocolConfig = new ProtocolConfig("dubbo");

        bootstrap.application("generic-dubbo-client")
                .registry(registryConfig)
                .protocol(protocolConfig)
                .start();

        ReferenceConfig<UserProvider> reference = new ReferenceConfig<>();
        reference.setInterface(UserProvider.class);
        reference.setGroup("dubbo");
        reference.setVersion("1.0.0");
        reference.setGeneric("false");

        UserProvider userProvider = reference.get();

        System.out.println("\n=== Testing Dubbo Protocol ===");

        callGetUser(userProvider);
        callGetOneUser(userProvider);
        callGetUsers(userProvider);
        callGetUsersMap(userProvider);
    }

    private static void callGetUser(UserProvider userProvider) {
        System.out.println("GetUser1(String userId) res: " + userProvider.GetUser1("A003"));
        System.out.println("GetUser2(String userId, String name) res: " + userProvider.GetUser2("A003", "lily"));
        System.out.println("GetUser3(int userCode) res: " + userProvider.GetUser3(1));
        System.out.println("GetUser4(int userCode, String name) res: " + userProvider.GetUser4(1, "zhangsan"));
    }

    private static void callGetOneUser(UserProvider userProvider) {
        System.out.println("GetOneUser() res: " + userProvider.GetOneUser());
    }

    private static void callGetUsers(UserProvider userProvider) {
        System.out.println("Call GetUsers");
        String[] userIds = new String[]{"001", "002"};
        System.out.println("GetUsers res: " + userProvider.GetUsers(userIds));
    }

    private static void callGetUsersMap(UserProvider userProvider) {
        System.out.println("Call GetUsersMap");
        String[] userIds = new String[]{"001", "002"};
        System.out.println("GetUsersMap res: " + userProvider.GetUsersMap(userIds));
    }
}