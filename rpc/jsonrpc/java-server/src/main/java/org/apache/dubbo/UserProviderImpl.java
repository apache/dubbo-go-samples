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

import java.util.HashMap;
import java.util.Map;

public class UserProviderImpl implements UserProvider {
    Map<String, User> userMap = new HashMap<String, User>();

    public UserProviderImpl() {
        userMap.put("A000", new User("0", "Alex Stocks", 31, Gender.MAN));
        userMap.put("A001", new User("A001", "ZhangSheng", 18, Gender.MAN));
        userMap.put("A002", new User("A002", "Lily", 20, Gender.WOMAN));
        userMap.put("A003", new User("A003", "Moorse", 30, Gender.MAN));
    }

    public User getUser(String userId) {
        User user = userMap.get(userId);
        if (user == null) {
            throw new IllegalArgumentException("invalid user id: " + userId);
        }
        return user;
    }

}
