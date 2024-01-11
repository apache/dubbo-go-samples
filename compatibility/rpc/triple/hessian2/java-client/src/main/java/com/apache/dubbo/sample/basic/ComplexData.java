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

import java.io.Serializable;
import java.util.Arrays;
import java.util.HashMap;
import java.util.StringJoiner;

public class ComplexData implements Serializable {
    boolean booleanData;
    String stringData;

    short int16Data;
    int intData;
    long int64Data;

    User userDefinedData;
    byte [] byteData;
    HashMap<String, String> stringStringHashMap;
    HashMap<String, User> stringUserDefinedPtrMapData;
    String[] arrayListData;
    User[] arrayUserData;

    @Override
    public String toString() {
        return new StringJoiner(", ", ComplexData.class.getSimpleName() + "[", "]")
                .add("booleanData=" + booleanData)
                .add("stringData='" + stringData + "'")
                .add("int16Data=" + int16Data)
                .add("intData=" + intData)
                .add("int64Data=" + int64Data)
                .add("userDefinedData=" + userDefinedData)
                .add("byteData=" + Arrays.toString(byteData))
                .add("stringStringHashMap=" + stringStringHashMap)
                .add("stringUserDefinedPtrMapData=" + stringUserDefinedPtrMapData)
                .add("arrayListData=" + Arrays.toString(arrayListData))
                .add("arrayUserData=" + Arrays.toString(arrayUserData))
                .toString();
    }
}
