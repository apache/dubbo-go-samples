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

import org.apache.dubbo.config.ApplicationConfig;
import org.apache.dubbo.config.ReferenceConfig;
import org.apache.dubbo.config.RegistryConfig;
import org.apache.dubbo.rpc.service.GenericService;

import java.util.HashMap;
import java.util.Map;

public class ApiTripleConsumer {
    private static final String ZOOKEEPER_ADDRESS = "zookeeper://127.0.0.1:2181";
    private static final String SERVICE_VERSION = "1.0.0";
    private static final String SERVICE_GROUP = "triple";

    public static void main(String[] args) {
        System.out.println("\n========== Java Triple Generic Call Test ==========");
        System.out.println("Registry: " + ZOOKEEPER_ADDRESS);
        System.out.println();

        ReferenceConfig<GenericService> reference = new ReferenceConfig<>();
        reference.setApplication(new ApplicationConfig("generic-java-client"));
        reference.setRegistry(new RegistryConfig(ZOOKEEPER_ADDRESS));
        reference.setInterface("org.apache.dubbo.samples.UserProvider");
        reference.setGroup(SERVICE_GROUP);
        reference.setVersion(SERVICE_VERSION);
        reference.setGeneric("true");
        reference.setProtocol("tri");
        reference.setCheck(false);

        GenericService genericService = reference.get();

        System.out.println("Connected to server via ZooKeeper, starting tests...\n");

        int passed = 0;
        int failed = 0;

        // ==================== 1. Basic type parameters ====================
        System.out.println("--- 1. Basic type parameters ---");

        // Test GetUser1(String)
        try {
            Object result = genericService.$invoke("GetUser1",
                    new String[]{"java.lang.String"},
                    new Object[]{"A003"});
            System.out.println("[PASS] GetUser1(String): " + result);
            passed++;
        } catch (Exception e) {
            System.out.println("[FAIL] GetUser1(String): " + e.getMessage());
            failed++;
        }

        // Test GetUser2(String, String)
        try {
            Object result = genericService.$invoke("GetUser2",
                    new String[]{"java.lang.String", "java.lang.String"},
                    new Object[]{"A003", "lily"});
            System.out.println("[PASS] GetUser2(String, String): " + result);
            passed++;
        } catch (Exception e) {
            System.out.println("[FAIL] GetUser2(String, String): " + e.getMessage());
            failed++;
        }

        // Test GetUser3(int)
        try {
            Object result = genericService.$invoke("GetUser3",
                    new String[]{"int"},
                    new Object[]{1});
            System.out.println("[PASS] GetUser3(int): " + result);
            passed++;
        } catch (Exception e) {
            System.out.println("[FAIL] GetUser3(int): " + e.getMessage());
            failed++;
        }

        // Test GetUser4(int, String)
        try {
            Object result = genericService.$invoke("GetUser4",
                    new String[]{"int", "java.lang.String"},
                    new Object[]{1, "zhangsan"});
            System.out.println("[PASS] GetUser4(int, String): " + result);
            passed++;
        } catch (Exception e) {
            System.out.println("[FAIL] GetUser4(int, String): " + e.getMessage());
            failed++;
        }

        // Test GetOneUser()
        try {
            Object result = genericService.$invoke("GetOneUser",
                    new String[]{},
                    new Object[]{});
            System.out.println("[PASS] GetOneUser(): " + result);
            passed++;
        } catch (Exception e) {
            System.out.println("[FAIL] GetOneUser(): " + e.getMessage());
            failed++;
        }

        // ==================== 2. Array/Collection types ====================
        System.out.println("\n--- 2. Array/Collection types ---");

        // Test GetUsers(String[]) -> User[]
        try {
            Object result = genericService.$invoke("GetUsers",
                    new String[]{"[Ljava.lang.String;"},
                    new Object[]{new String[]{"001", "002"}});
            System.out.println("[PASS] GetUsers(String[]): " + result);
            passed++;
        } catch (Exception e) {
            System.out.println("[FAIL] GetUsers(String[]): " + e.getMessage());
            failed++;
        }

        // Test GetUsersMap(String[]) -> Map<String, User>
        try {
            Object result = genericService.$invoke("GetUsersMap",
                    new String[]{"[Ljava.lang.String;"},
                    new Object[]{new String[]{"001", "002"}});
            System.out.println("[PASS] GetUsersMap(String[]): " + result);
            passed++;
        } catch (Exception e) {
            System.out.println("[FAIL] GetUsersMap(String[]): " + e.getMessage());
            failed++;
        }

        // Test QueryAll() -> Map<String, User>
        try {
            Object result = genericService.$invoke("QueryAll",
                    new String[]{},
                    new Object[]{});
            System.out.println("[PASS] QueryAll(): " + result);
            passed++;
        } catch (Exception e) {
            System.out.println("[FAIL] QueryAll(): " + e.getMessage());
            failed++;
        }

        // ==================== 3. Custom type parameters ====================
        System.out.println("\n--- 3. Custom type parameters ---");

        // Test QueryUser(User)
        try {
            Map<String, Object> userParam = new HashMap<>();
            userParam.put("class", "org.apache.dubbo.samples.User");
            userParam.put("id", "U001");
            userParam.put("name", "TestUser");
            userParam.put("age", 25);

            Object result = genericService.$invoke("QueryUser",
                    new String[]{"org.apache.dubbo.samples.User"},
                    new Object[]{userParam});
            System.out.println("[PASS] QueryUser(User): " + result);
            passed++;
        } catch (Exception e) {
            System.out.println("[FAIL] QueryUser(User): " + e.getMessage());
            e.printStackTrace();
            failed++;
        }

        // Test QueryUsers(User[])
        try {
            Map<String, Object> user1 = new HashMap<>();
            user1.put("class", "org.apache.dubbo.samples.User");
            user1.put("id", "U001");
            user1.put("name", "User1");
            user1.put("age", 20);

            Map<String, Object> user2 = new HashMap<>();
            user2.put("class", "org.apache.dubbo.samples.User");
            user2.put("id", "U002");
            user2.put("name", "User2");
            user2.put("age", 30);

            Object[] userArray = new Object[]{user1, user2};

            Object result = genericService.$invoke("QueryUsers",
                    new String[]{"[Lorg.apache.dubbo.samples.User;"},
                    new Object[]{userArray});
            System.out.println("[PASS] QueryUsers(User[]): " + result);
            passed++;
        } catch (Exception e) {
            System.out.println("[FAIL] QueryUsers(User[]): " + e.getMessage());
            e.printStackTrace();
            failed++;
        }

        // ==================== Test summary ====================
        System.out.println("\n========== Test Completed ==========");
        System.out.println("Passed: " + passed + ", Failed: " + failed + ", Total: " + (passed + failed));

        if (failed > 0) {
            System.out.println("\n[WARN] Some tests failed!");
        } else {
            System.out.println("\n[OK] All tests passed!");
        }

        System.exit(0);
    }
}
