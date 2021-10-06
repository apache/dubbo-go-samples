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

import java.util.Arrays;

public class ComplexProviderImpl implements ComplexProvider {
    @Override
    public int InvokeWithMultiBasicData(String str, byte[] data, int num, boolean boolValue) {
        System.out.println("InvokeWithMultiBasicData, str: " + str + ", data: " + Arrays.toString(data) + ", num: " + num + ", boolValue:" +
                " " + boolValue);
        return num;
    }

    @Override
    public void InvokeWithEmptyReq() {
        System.out.println("InvokeWithEmptyReq");
    }

    @Override
    public void InvokeWithSingleString(String req) {
        System.out.println("InvokeWithEmptyReq" + req);
    }

    @Override
    public void InvokeWithStringList(String[] req) {
        System.out.println("InvokeWithEmptyReq" + req);
    }

    @Override
    public void InvokeWithMultiString(String str1, String str2, String str3) {
        System.out.println("InvokeWithEmptyReq" + str1 + str2 + str3);
    }

    @Override
    public String InvokeWithEmptyReqStringRsp() {
        System.out.println("InvokeWithEmptyReq");
        return "invoke success";
    }

    @Override
    public ComplexData InvokeWithComplexReqComplexRspPtr(ComplexData complexData) {
        System.out.println("InvokeWithComplexReqComplexRspPtr = "+ complexData.GetString());
        return complexData;
    }
}
