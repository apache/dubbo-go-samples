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

package org.apache.dubbo.samples.async;

import org.apache.dubbo.common.constants.CommonConstants;
import org.apache.dubbo.config.ApplicationConfig;
import org.apache.dubbo.config.ProtocolConfig;
import org.apache.dubbo.config.ServiceConfig;
import org.apache.dubbo.config.bootstrap.DubboBootstrap;
import org.apache.dubbo.samples.async.proto.*;

import java.io.IOException;
import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.CompletableFuture;

public class JavaAsyncServer {

    public static void main(String[] args) throws IOException {
        ServiceConfig<UserProvider> userProviderService = new ServiceConfig<>();
        userProviderService.setInterface(UserProvider.class);
        userProviderService.setRef(new UserProviderImpl());

        ServiceConfig<UserProviderV2> userProviderV2Service = new ServiceConfig<>();
        userProviderV2Service.setInterface(UserProviderV2.class);
        userProviderV2Service.setRef(new UserProviderV2Impl());

        DubboBootstrap bootstrap = DubboBootstrap.getInstance();
        bootstrap.application(new ApplicationConfig("java-async-server"))
                .protocol(new ProtocolConfig(CommonConstants.TRIPLE, 50051))
                .service(userProviderService)
                .service(userProviderV2Service)
                .start();

        System.out.println("Dubbo triple java async server started on port 50051");
        System.in.read();
    }

    /**
     * UserProvider async implementation
     */
    static class UserProviderImpl extends DubboUserProviderTriple.UserProviderImplBase {
        private static final Map<String, User> userMap = new HashMap<>();

        static {
            userMap.put("000", User.newBuilder()
                    .setId("000")
                    .setName("Alex Stocks")
                    .setAge(31)
                    .setTime(System.currentTimeMillis() / 1000)
                    .setSex(Gender.MAN)
                    .build());

            userMap.put("001", User.newBuilder()
                    .setId("001")
                    .setName("ZhangSheng")
                    .setAge(18)
                    .setTime(System.currentTimeMillis() / 1000)
                    .setSex(Gender.MAN)
                    .build());

            userMap.put("002", User.newBuilder()
                    .setId("002")
                    .setName("Lily")
                    .setAge(20)
                    .setTime(System.currentTimeMillis() / 1000)
                    .setSex(Gender.WOMAN)
                    .build());

            userMap.put("003", User.newBuilder()
                    .setId("003")
                    .setName("Moorse")
                    .setAge(30)
                    .setTime(System.currentTimeMillis() / 1000)
                    .setSex(Gender.WOMAN)
                    .build());
        }

        @Override
        public CompletableFuture<GetUserResponse> getUserAsync(GetUserRequest request) {
            return CompletableFuture.supplyAsync(() -> {
                System.out.println("Received GetUser request, id: " + request.getId());
                User user = userMap.get(request.getId());
                if (user == null) {
                    return GetUserResponse.getDefaultInstance();
                }
                System.out.println("Returning user: " + user.getName() + ", age: " + user.getAge());
                return GetUserResponse.newBuilder()
                        .setUser(user)
                        .build();
            });
        }
    }

    /**
     * UserProviderV2 async implementation
     */
    static class UserProviderV2Impl extends DubboUserProviderV2Triple.UserProviderV2ImplBase {
        @Override
        public CompletableFuture<SayHelloResponse> sayHelloAsync(SayHelloRequest request) {
            return CompletableFuture.supplyAsync(() -> {
                System.out.println("Received SayHello request, userId: " + request.getUserId());
                return SayHelloResponse.getDefaultInstance();
            });
        }
    }
}
