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

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.lang.Integer;

public class ApiConsumer {
    private static final Logger logger = LoggerFactory.getLogger("userLogger"); // Output to com.dubbogo.user-server.log

    private static GenericService genericService;

     public static void main(String[] args) {
        initConfig();

         callGetUser();
        callGetOneUser();
        callGetUsers();
        callGetUsersMap();
        callQueryUser();
        callQueryUsers();
         callQueryAll();
     }

     private static void initConfig(){
        logger.info("\n\n\nstart to init config\n\n\n");
        ApplicationConfig applicationConfig = new ApplicationConfig();
        ReferenceConfig<GenericService> reference = new ReferenceConfig<GenericService>();
        applicationConfig.setName("user-info-server");
        reference.setApplication(applicationConfig);
        RegistryConfig registryConfig = new RegistryConfig();
        registryConfig.setAddress("zookeeper://127.0.0.1:2181");
        reference.setRegistry(registryConfig);
        reference.setGeneric(true);
        reference.setInterface("org.apache.dubbo.UserProvider");
        genericService = reference.get();
     }

     private static void callGetUser(){
         logger.info("\n\n\nCall GetUser");
         logger.info("Start to generic invoke");
         Object result;

        result = genericService.$invoke("GetUser1", new String[]{"java.lang.String"} , new Object[]{"A003"});
         logger.info("\n\n\n" + "GetUser1(String userId) " + "res: " + result + "\n\n\n");

         result = genericService.$invoke("GetUser2", new String[]{"java.lang.String", "java.lang.String"} , new Object[]{"A003", "lily"});
         logger.info("\n\n\n" + "GetUser2(String userId, String name) " + "res: " + result + "\n\n\n");

         result = genericService.$invoke("GetUser3", new String[]{"int"} , new Object[]{1});
         logger.info("\n\n\n" + "GetUser3(int userCode) " + "res: " + result + "\n\n\n");

         result = genericService.$invoke("GetUser4", new String[]{"int", "java.lang.String"} , new Object[]{1, "zhangsan"});
         logger.info("\n\n\n" + "GetUser4(int userCode, String name) " + "res: " + result + "\n\n\n");
     }

    private static void callGetOneUser(){
        logger.info("\n\n\nCall GetOneUser");
        logger.info("Start to generic invoke");

        Object result = genericService.$invoke("GetOneUser", new String[]{} , new Object[]{});
        logger.info("\n\n\n" + "GetOneUser() " + "res: " + result + "\n\n\n");
    }

    private static void callGetUsers(){
        logger.info("\n\n\nCall GetUsers");
        logger.info("Start to generic invoke");

        List<String> userIdList = new ArrayList<String>();
        userIdList.add("001");
        userIdList.add("002");
        userIdList.add("003");
        userIdList.add("004");

        Object result = genericService.$invoke("GetUsers", new String[]{"java.util.List"}, new Object[]{userIdList});
        logger.info("\n\n\n" + "GetUsers(List<String> userIdList) " + "res: " + result + "\n\n\n");
    }

    private static void callGetUsersMap(){
        logger.info("\n\n\nCall GetUsersMap");
        logger.info("Start to generic invoke");
        Object result;

        List<String> userIdList = new ArrayList<String>();
        userIdList.add("001");
        userIdList.add("002");
        userIdList.add("003");
        userIdList.add("004");

        result = genericService.$invoke("GetUsersMap", new String[]{"java.util.List"}, new Object[]{userIdList});
        logger.info("\n\n\n" + "GetUserMap(List<String> userIdList) " + "res: " + result + "\n\n\n");
    }

     private static void callQueryUser(){
         logger.info("\n\n\nCall QueryUser");
         logger.info("Start to generic invoke");

        User user = new User();
        user.setName("Patrick");
        user.setId("id");
        user.setAge(10);
        Object result = genericService.$invoke("queryUser", new String[]{"org.apache.dubbo.User"} , new Object[]{user});
         logger.info("\n\n\n" + "queryUser(User user) " + "res: " + result + "\n\n\n");
     }

     private static void callQueryUsers(){
         logger.info("\n\n\nCall QueryUsers");
         logger.info("Start to generic invoke");

        ArrayList<Map> userArr = new ArrayList<>();
        Map<Object, Object> userMap1 = new HashMap<>();
        userMap1.put("id", "A001");
        userMap1.put("name", "Patrick");
        userMap1.put("age", 10);
        Map<Object, Object> userMap2 = new HashMap<>();
        userMap2.put("id", "A002");
        userMap2.put("name", "xavier-niu");
        userMap2.put("age", 24);
        userArr.add(userMap1);
        userArr.add(userMap2);

        Object result = genericService.$invoke("queryUsers", new String[]{"java.util.ArrayList"} , new Object[]{userArr});
         logger.info("\n\n\n" + "queryUsers(ArrayList<User> userList) " + "res: " + result + "\n\n\n");
     }

    private static void callQueryAll(){
        logger.info("\n\n\nCall queryAll");
        logger.info("Start to generic invoke");

        Object result = genericService.$invoke("queryAll", new String[]{} , new Object[]{});
        logger.info("\n\n\n" + "queryAll() " + "res: " + result + "\n\n\n");
    }
}
