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

package com.ikurento.user;

import java.text.SimpleDateFormat;
import java.util.Date;
import java.util.concurrent.CountDownLatch;

import io.grpc.stub.StreamObserver;


public class Consumer {
    // Define a private variable (Required in Spring)
    private UserProviderGrpc.IUserProvider userProvider;

    // Spring DI (Required in Spring)
    public void setUserProvider(UserProviderGrpc.IUserProvider u) {
        this.userProvider = u;
    }

    // Start the entry function for consumer (Specified in the configuration file)
    public void start() throws Exception {
        System.out.println("\n\ntest");
        testGetUser();
        System.out.println("\n\ntestList");
        testGetUserList();
        System.out.println("\n\ntestErr");
        testGetErr();
        System.out.println("\n\ntestStream");
        testGetUserByStream();
    }

    private void testGetUser() throws Exception {
        try {
            UserID userID = UserID.newBuilder().setID("A003").build();
            User user1 = userProvider.getUser(userID);
            System.out.println("[" + new SimpleDateFormat("HH:mm:ss").format(new Date()) + "] " +
                    " UserInfo, ID:" + user1.getID() + ", name:" + user1.getName() + ", sex:" + user1.getSex().toString()
                    + ", age:" + user1.getAge() + ", time:" + user1.getTime().toString());
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    private void testGetUserList() throws Exception {
        try {
            UserIDList userIDList = UserIDList.newBuilder().addID("A003").addID("A002").build();
            UserList userList = userProvider.getUserList(userIDList);
            for (User user : userList.getUserList()) {
                System.out.println("[" + new SimpleDateFormat("HH:mm:ss").format(new Date()) + "] " +
                        " UserInfo, ID:" + user.getID() + ", name:" + user.getName() + ", sex:" + user.getSex().toString()
                        + ", age:" + user.getAge() + ", time:" + user.getTime().toString());
            }
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    private void testGetErr() throws Exception {
        try {
            UserID userID = UserID.newBuilder().setID("A003").build();
            userProvider.getErr(userID);
        } catch (Throwable t) {
            System.out.println("*************exception***********");
            t.printStackTrace();
        }
    }

    private void testGetUserByStream() {
        try {
            UserID userID = UserID.newBuilder().setID("A001").build();
            userProvider.getUserByStream(userID, new StreamObserver<User>() {
                @Override
                public void onNext(User user) {
                    System.out.println("[" + new SimpleDateFormat("HH:mm:ss").format(new Date()) + "] " +
                            " UserInfo, ID:" + user.getID() + ", name:" + user.getName() + ", sex:" + user.getSex().toString()
                            + ", age:" + user.getAge() + ", time:" + user.getTime().toString());
                }

                @Override
                public void onError(Throwable t) {
                    t.printStackTrace();
                }

                @Override
                public void onCompleted() {
                    System.out.println("Done");
                }
            });

        } catch (Exception e) {
            System.out.println("*************exception***********");
            e.printStackTrace();
        }
    }
}
