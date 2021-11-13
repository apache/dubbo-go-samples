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

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.*;

public class UserProviderImpl implements UserProvider {
    Map<String, User> userMap = new HashMap<>();

    public UserProviderImpl() {
        userMap.put("A001", new User("A001", "demo-zhangsan", 18));
        userMap.put("A002", new User("A002", "demo-lisi", 20));
        userMap.put("A003", new User("A003", "demo-lily", 23));
        userMap.put("A004", new User("A004", "demo-lisa", 32));
    }
    
    @Override
    public User GetUser(String userId) throws Exception {
        User user = userMap.get(userId);
        if (user == null) {
            throw new IllegalArgumentException("invalid user id: " + userId);
        }
        return user;
    }
}
