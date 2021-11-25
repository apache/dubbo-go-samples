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

import java.util.ArrayList;
import java.util.List;
import java.util.Map;
//import org.apache.dubbo.rpc.filter.GenericFilter;

public interface UserProvider {

    boolean isLimit(Gender gender, String name);

    User GetUser1(String userId); // the first alpha is Upper case to compatible with golang.

    User GetUser2(String userId, String name);

    User GetUser3(int userCode);

    User GetUser4(int userCode, String name);

    User GetOneUser();

    List<User> GetUsers(List<String> userIdList);

    Map<String, User> GetUsersMap(List<String> userIdList);

    User queryUser(User user);

    List<User> queryUsers(ArrayList<User> userObjectList);

    Map<String, User> queryAll();

    User GetErr(String userId) throws Exception;

    int Calc(int a, int b);

    Response<Integer> Sum(int a, int b);
}
