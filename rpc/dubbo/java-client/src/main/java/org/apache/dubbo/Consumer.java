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

import java.text.SimpleDateFormat;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Date;
import java.util.HashMap;
import java.util.List;

import com.alibaba.dubbo.rpc.service.EchoService;
import org.apache.dubbo.common.constants.CommonConstants;
import org.apache.dubbo.config.ApplicationConfig;
import org.apache.dubbo.config.ReferenceConfig;
import org.apache.dubbo.config.RegistryConfig;
import org.springframework.context.support.ClassPathXmlApplicationContext;

public class Consumer {
    // Define a private variable (Required in Spring)
    private static UserProvider userProvider;
    private static UserProvider userProvider1;
    private static UserProvider userProvider2;

    public static void main(String[] args) throws Exception {
        ClassPathXmlApplicationContext context = new ClassPathXmlApplicationContext(new String[]{"META-INF/spring/dubbo.consumer.xml"});
        userProvider = (UserProvider)context.getBean("userProvider");
        userProvider1 = (UserProvider)context.getBean("userProvider1");
        userProvider2 = (UserProvider)context.getBean("userProvider2");

        start();
        startComplexConsumerService();
        // TODO when upgrade hessian version, remember to delete this comment
       startWrapperArrayClassService();
    }

    // Start the entry function for consumer (Specified in the configuration file)
    public static void start() throws Exception {
        System.out.println("\n\ntest");
        testGetUser();
        testGetUsers();
        System.out.println("\n\ntest1");
        testGetUser1();
        testGetUsers1();
        System.out.println("\n\ntest2");
        testGetUser2();
        testGetUsers2();
        Thread.sleep(2000);
    }

    private static void testGetUser() throws Exception {
        try {
            EchoService echoService = (EchoService)userProvider;
            Object status = echoService.$echo("OK");
            System.out.println("echo: "+status);
            User user1 = userProvider.GetUser("A003");
            System.out.println("[" + new SimpleDateFormat("HH:mm:ss").format(new Date()) + "] " +
                    " UserInfo, ID:" + user1.getId() + ", name:" + user1.getName() + ", sex:" + user1.getSex().toString()
                    + ", age:" + user1.getAge() + ", time:" + user1.getTime().toString());
            User user2 = userProvider.GetUser0("A003","Moorse");
            System.out.println("[" + new SimpleDateFormat("HH:mm:ss").format(new Date()) + "] " +
                    " UserInfo, ID:" + user2.getId() + ", name:" + user2.getName() + ", sex:" + user2.getSex().toString()
                    + ", age:" + user2.getAge() + ", time:" + user2.getTime().toString());
            User user3 = userProvider.getUser(1);
            System.out.println("[" + new SimpleDateFormat("HH:mm:ss").format(new Date()) + "] " +
                    " UserInfo, ID:" + user3.getId() + ", name:" + user3.getName() + ", sex:" + user3.getSex().toString()
                    + ", age:" + user3.getAge() + ", time:" + user3.getTime());
            User user4 = userProvider.getUser(1, "name");
            System.out.println("[" + new SimpleDateFormat("HH:mm:ss").format(new Date()) + "] " +
                    " UserInfo, ID:" + user4.getId() + ", name:" + user4.getName() + ", sex:" + user4.getSex().toString()
                    + ", age:" + user4.getAge() + ", time:" + user4.getTime());
            userProvider.GetUser3();
            System.out.println("GetUser3 succ");

            User user9 = userProvider.GetUser1("A003");
        } catch (Throwable e) {
            System.out.println("*************exception***********");
            e.printStackTrace();
        }
        try {
            userProvider.GetErr("A003");
        } catch (Throwable t) {
            System.out.println("*************exception***********");
            t.printStackTrace();
        }
    }

    private static void testGetUsers() throws Exception {
        try {
            List<String> userIDList = new ArrayList<String>();
            userIDList.add("A001");
            userIDList.add("A002");
            userIDList.add("A003");

            List<User> userList = userProvider.GetUsers(userIDList);

            for (int i = 0; i < userList.size(); i++) {
                User user = userList.get(i);
                System.out.println("[" + new SimpleDateFormat("HH:mm:ss").format(new Date()) + "] " +
                        " UserInfo, ID:" + user.getId() + ", name:" + user.getName() + ", sex:" + user.getSex().toString()
                        + ", age:" + user.getAge() + ", time:" + user.getTime().toString());
            }
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    private static  void testGetUser1() throws Exception {
        try {
            EchoService echoService = (EchoService)userProvider1;
            Object status = echoService.$echo("OK");
            System.out.println("echo: "+status);
            User user1 = userProvider1.GetUser("A003");
            System.out.println("[" + new SimpleDateFormat("HH:mm:ss").format(new Date()) + "] " +
                    " UserInfo, ID:" + user1.getId() + ", name:" + user1.getName() + ", sex:" + user1.getSex().toString()
                    + ", age:" + user1.getAge() + ", time:" + user1.getTime().toString());
            User user2 = userProvider1.GetUser0("A003","Moorse");
            System.out.println("[" + new SimpleDateFormat("HH:mm:ss").format(new Date()) + "] " +
                    " UserInfo, ID:" + user2.getId() + ", name:" + user2.getName() + ", sex:" + user2.getSex().toString()
                    + ", age:" + user2.getAge() + ", time:" + user2.getTime().toString());
            User user3 = userProvider1.getUser(1);
            System.out.println("[" + new SimpleDateFormat("HH:mm:ss").format(new Date()) + "] " +
                    " UserInfo, ID:" + user3.getId() + ", name:" + user3.getName() + ", sex:" + user3.getSex().toString()
                    + ", age:" + user3.getAge() + ", time:" + user3.getTime());
            User user4 = userProvider1.getUser(1, "name");
            System.out.println("[" + new SimpleDateFormat("HH:mm:ss").format(new Date()) + "] " +
                    " UserInfo, ID:" + user4.getId() + ", name:" + user4.getName() + ", sex:" + user4.getSex().toString()
                    + ", age:" + user4.getAge() + ", time:" + user4.getTime());
            userProvider1.GetUser3();
            System.out.println("GetUser3 succ");

            User user9 = userProvider1.GetUser1("A003");
        } catch (Throwable e) {
            System.out.println("*************exception***********");
            e.printStackTrace();
        }
        try {
            userProvider1.GetErr("A003");
        } catch (Throwable t) {
            System.out.println("*************exception***********");
            t.printStackTrace();
        }
    }

    private static  void testGetUsers1() throws Exception {
        try {
            List<String> userIDList = new ArrayList<String>();
            userIDList.add("A001");
            userIDList.add("A002");
            userIDList.add("A003");

            List<User> userList = userProvider1.GetUsers(userIDList);

            for (int i = 0; i < userList.size(); i++) {
                User user = userList.get(i);
                System.out.println("[" + new SimpleDateFormat("HH:mm:ss").format(new Date()) + "] " +
                        " UserInfo, ID:" + user.getId() + ", name:" + user.getName() + ", sex:" + user.getSex().toString()
                        + ", age:" + user.getAge() + ", time:" + user.getTime().toString());
            }
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    private static void testGetUser2() throws Exception {
        try {
            EchoService echoService = (EchoService)userProvider2;
            Object status = echoService.$echo("OK");
            System.out.println("echo: "+status);
            User user1 = userProvider2.GetUser("A003");
            System.out.println("[" + new SimpleDateFormat("HH:mm:ss").format(new Date()) + "] " +
                    " UserInfo, ID:" + user1.getId() + ", name:" + user1.getName() + ", sex:" + user1.getSex().toString()
                    + ", age:" + user1.getAge() + ", time:" + user1.getTime().toString());
            User user2 = userProvider2.GetUser0("A003","Moorse");
            System.out.println("[" + new SimpleDateFormat("HH:mm:ss").format(new Date()) + "] " +
                    " UserInfo, ID:" + user2.getId() + ", name:" + user2.getName() + ", sex:" + user2.getSex().toString()
                    + ", age:" + user2.getAge() + ", time:" + user2.getTime().toString());
            User user3 = userProvider2.getUser(1);
            System.out.println("[" + new SimpleDateFormat("HH:mm:ss").format(new Date()) + "] " +
                    " UserInfo, ID:" + user3.getId() + ", name:" + user3.getName() + ", sex:" + user3.getSex().toString()
                    + ", age:" + user3.getAge() + ", time:" + user3.getTime());
            User user4 = userProvider2.getUser(1, "name");
            System.out.println("[" + new SimpleDateFormat("HH:mm:ss").format(new Date()) + "] " +
                    " UserInfo, ID:" + user4.getId() + ", name:" + user4.getName() + ", sex:" + user4.getSex().toString()
                    + ", age:" + user4.getAge() + ", time:" + user4.getTime());
            userProvider2.GetUser3();
            System.out.println("GetUser3 succ");

            User user9 = userProvider2.GetUser1("A003");
        } catch (Throwable e) {
            System.out.println("*************exception***********");
            e.printStackTrace();
        }
        try {
            userProvider2.GetErr("A003");
        } catch (Throwable t) {
            System.out.println("*************exception***********");
            t.printStackTrace();
        }
    }

    private static void testGetUsers2() throws Exception {
        try {
            List<String> userIDList = new ArrayList<String>();
            userIDList.add("A001");
            userIDList.add("A002");
            userIDList.add("A003");

            List<User> userList = userProvider2.GetUsers(userIDList);

            for (int i = 0; i < userList.size(); i++) {
                User user = userList.get(i);
                System.out.println("[" + new SimpleDateFormat("HH:mm:ss").format(new Date()) + "] " +
                        " UserInfo, ID:" + user.getId() + ", name:" + user.getName() + ", sex:" + user.getSex().toString()
                        + ", age:" + user.getAge() + ", time:" + user.getTime().toString());
            }
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    public static void startComplexConsumerService() {
        ReferenceConfig<ComplexProvider> ref = new ReferenceConfig<>();
        ref.setInterface(ComplexProvider.class);
        ref.setCheck(false);
        ref.setProtocol(CommonConstants.DUBBO_PROTOCOL);
        ref.setLazy(true);
        ref.setTimeout(100000);
        ref.setApplication(new ApplicationConfig("demo-consumer"));

        ref.setRegistry(new RegistryConfig("zookeeper://127.0.0.1:2181"));
        final ComplexProvider complexProvider = ref.get();

//        complexProvider.invokeWithEmptyReq();
//        complexProvider.invokeWithSingleString("single string");
//        complexProvider.invokeWithMultiString("string1", "string2", "string3");
//        String [] strList = new String []{"first string"," second string"};
//        complexProvider.invokeWithStringList(strList );
//       String rsp = complexProvider.invokeWithEmptyReqStringRsp();
//       System.out.println("get rsp = "+  rsp);


       ComplexData cpxData = new ComplexData();
       cpxData.booleanData = true;
       cpxData.stringData = "test string";
       cpxData.byteData =  new byte[] {1, 12, 4, 3, 3,3};
       cpxData.int16Data =16;
       cpxData.intData = 32;
       cpxData.int64Data = 64;
       cpxData.arrayListData = new String[]{"array1", "array2"};
//       cpxData.arrayUserData = new User[]{new User(), new User(), new User()};
        cpxData.userDefinedData = new User();
        cpxData.userDefinedData.setAge(18);
        cpxData.userDefinedData.setId("iojfioj");
        cpxData.stringStringHashMap = new HashMap<>();
//        cpxData.stringStringHashMap.put("key1", "value");
//        cpxData.stringStringHashMap.put("key2", "value");
//        cpxData.stringUserDefinedPtrMapData = new HashMap<>();
//        cpxData.stringUserDefinedPtrMapData.put("key1", new User());
//        cpxData.stringUserDefinedPtrMapData.put("key2", new User());

//        ComplexData response = complexProvider.invokeWithComplexReqComplexRspPtr(cpxData);
//        System.out.println("get complex = "+  response);

        int rsp = complexProvider.InvokeWithMultiBasicData("str",new byte[]{1, 3, 4,6,7}, 32, true);
        System.out.println("get multi basic rsp = "+  rsp);
    }

    public static void startWrapperArrayClassService() {
        ReferenceConfig<WrapperArrayClassProvider> ref = new ReferenceConfig<>();
        ref.setInterface(WrapperArrayClassProvider.class);
        ref.setCheck(false);
        ref.setProtocol(CommonConstants.DUBBO_PROTOCOL);
        ref.setLazy(true);
        ref.setTimeout(100000);
        ref.setApplication(new ApplicationConfig("demo-consumer"));

        ref.setRegistry(new RegistryConfig("zookeeper://127.0.0.1:2181"));
        final WrapperArrayClassProvider wrapperArrayClassProvider = ref.get();
        System.out.println("---InvokeWithJavaByteArray:" + Arrays.toString(wrapperArrayClassProvider.InvokeWithJavaByteArray(new Byte[]{10, 100})));
        System.out.println("---InvokeWithJavaCharacterArray" + Arrays.toString(wrapperArrayClassProvider.InvokeWithJavaCharacterArray(new Character[]{'a', 'b', 'c'})));
        System.out.println("---InvokeWithJavaShortArray" + Arrays.toString(wrapperArrayClassProvider.InvokeWithJavaShortArray(new Short[]{1, 2, 3})));
        System.out.println("---InvokeWithJavaIntegerArray" + Arrays.toString(wrapperArrayClassProvider.InvokeWithJavaIntegerArray(new Integer[]{4, 5, 6})));
        System.out.println("---InvokeWithJavaLongArray" + Arrays.toString(wrapperArrayClassProvider.InvokeWithJavaLongArray(new Long[]{7L, 8L, 9000000000000L})));
        System.out.println("---InvokeWithJavaFloatArray" + Arrays.toString(wrapperArrayClassProvider.InvokeWithJavaFloatArray(new Float[]{1.2f, 2.3f, 3.0f})));
        System.out.println("---InvokeWithJavaDoubleArray" + Arrays.toString(wrapperArrayClassProvider.InvokeWithJavaDoubleArray(new Double[]{4.0, 5.1, 6.0})));
        System.out.println("---InvokeWithJavaBooleanArray" + Arrays.toString(wrapperArrayClassProvider.InvokeWithJavaBooleanArray(new Boolean[]{true, false, true})));
    }
}
