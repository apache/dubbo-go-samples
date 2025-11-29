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

public interface UserProvider {
    User GetUser1(String userId);

    User GetUser2(String userId, String name);

    User GetUser3(int userCode);

    User GetUser4(int userCode, String name);

    User GetOneUser();

    List<User> GetUsers(String[] userIdList);

    Map<String, User> GetUsersMap(String[] userIdList);

    User QueryUser(User user);

    List<User> QueryUsers(List<User> userObjectList);

    Map<String, User> QueryAll();
}