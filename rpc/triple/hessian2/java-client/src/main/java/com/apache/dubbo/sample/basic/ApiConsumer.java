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

package com.apache.dubbo.sample.basic;

import org.apache.dubbo.common.constants.CommonConstants;
import org.apache.dubbo.config.ApplicationConfig;
import org.apache.dubbo.config.ReferenceConfig;
import org.apache.dubbo.config.RegistryConfig;

import java.io.IOException;
import java.util.HashMap;

public class ApiConsumer {
    public static void main(String[] args) throws InterruptedException, IOException {
        ReferenceConfig<ComplexProvider> ref = new ReferenceConfig<>();
        ref.setInterface(ComplexProvider.class);
        ref.setCheck(false);
        ref.setProtocol(CommonConstants.TRIPLE);
        ref.setLazy(true);
        ref.setTimeout(100000);
        ref.setApplication(new ApplicationConfig("demo-consumer"));

        ref.setRegistry(new RegistryConfig("zookeeper://127.0.0.1:2181"));
        final ComplexProvider complexProvider = ref.get();

        complexProvider.invokeWithEmptyReq();
        complexProvider.invokeWithSingleString("single string");
        complexProvider.invokeWithMultiString("string1", "string2", "string3");
        String[] strList = new String[]{"first string", " second string"};
        complexProvider.invokeWithStringList(strList);
        String rsp = complexProvider.invokeWithEmptyReqStringRsp();
        System.out.println("get rsp = " + rsp);

        ComplexData cpxData = new ComplexData();
        cpxData.booleanData = true;
        cpxData.stringData = "test string";
        cpxData.byteData = new byte[]{1, 12, 4, 3, 3, 3};
        cpxData.int16Data = 16;
        cpxData.intData = 32;
        cpxData.int64Data = 64;
        cpxData.arrayListData = new String[]{"array1", "array2"};
        cpxData.arrayUserData = new User[]{new User(), new User(), new User()};
        cpxData.userDefinedData = new User();
        cpxData.userDefinedData.age = 18;
        cpxData.userDefinedData.id = "iojfioj";
        cpxData.stringUserDefinedPtrMapData = new HashMap<>();
        cpxData.stringUserDefinedPtrMapData.put("key1", new User());
        cpxData.stringUserDefinedPtrMapData.put("key2", new User());

        ComplexData response = complexProvider.invokeWithComplexReqComplexRspPtr(cpxData);
        System.out.println("get complex = " + response);

        int rsp1 = complexProvider.invokeWithMultiBasicData("str", new byte[]{1, 3, 4, 44, 7}, 32, true);
        System.out.println("get multi basic rsp = " + rsp1);
    }
}
