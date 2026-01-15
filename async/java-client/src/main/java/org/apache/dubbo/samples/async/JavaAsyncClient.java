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

import org.apache.dubbo.config.ApplicationConfig;
import org.apache.dubbo.config.ReferenceConfig;
import org.apache.dubbo.config.bootstrap.DubboBootstrap;
import org.apache.dubbo.samples.async.proto.*;

import java.util.concurrent.CompletableFuture;
import java.util.concurrent.CountDownLatch;

public class JavaAsyncClient {

    public static void main(String[] args) throws Exception {
        ReferenceConfig<UserProvider> userProviderRef = new ReferenceConfig<>();
        userProviderRef.setInterface(UserProvider.class);
        userProviderRef.setUrl("tri://127.0.0.1:20000");

        ReferenceConfig<UserProviderV2> userProviderV2Ref = new ReferenceConfig<>();
        userProviderV2Ref.setInterface(UserProviderV2.class);
        userProviderV2Ref.setUrl("tri://127.0.0.1:20000");

        DubboBootstrap bootstrap = DubboBootstrap.getInstance();
        bootstrap.application(new ApplicationConfig("java-async-client"))
                .reference(userProviderRef)
                .reference(userProviderV2Ref)
                .start();

        UserProvider userProvider = userProviderRef.get();
        UserProviderV2 userProviderV2 = userProviderV2Ref.get();

        CountDownLatch latch1 = new CountDownLatch(1);
        GetUserRequest getUserReq = GetUserRequest.newBuilder().setId("003").build();
        CompletableFuture<GetUserResponse> future1 = userProvider.getUserAsync(getUserReq);

        System.out.println("async request sent, doing other work...");

        future1.whenComplete((response, throwable) -> {
            if (throwable != null) {
                System.err.println("error: " + throwable.getMessage());
            } else {
                User user = response.getUser();
                System.out.println("received user: " + user.getName() + ", age: " + user.getAge());
            }
            latch1.countDown();
        });

        CountDownLatch latch2 = new CountDownLatch(1);
        SayHelloRequest sayHelloReq = SayHelloRequest.newBuilder().setUserId("002").build();
        CompletableFuture<SayHelloResponse> future2 = userProviderV2.sayHelloAsync(sayHelloReq);

        future2.whenComplete((response, throwable) -> {
            if (throwable == null) {
                System.out.println("sayHello completed");
            }
            latch2.countDown();
        });

        latch1.await();
        latch2.await();
        bootstrap.stop();
    }
}
