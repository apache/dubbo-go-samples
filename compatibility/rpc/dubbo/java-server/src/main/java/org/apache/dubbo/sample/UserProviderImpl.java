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

package org.apache.dubbo.sample;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class UserProviderImpl implements UserProvider {
    private static final Logger LOG = LoggerFactory.getLogger("UserLogger"); // Output to com.dubbogo.user-server.log
    Map<String, User> userMap = new HashMap<String, User>();

    public UserProviderImpl() {
        userMap.put("001", new User("001", "demo-zhangsan", 18));
        userMap.put("002", new User("002", "demo-lisi", 20));
        userMap.put("003", new User("003", "demo-lily", 23));
        userMap.put("004", new User("004", "demo-lisa", 32));
    }

    @Override
    public Map<String, User> getUserMap() {
        return userMap;
    }


    @Override
    public User GetUser(User user) {
        return new User(user.getId(), "zhangsan", 18);
    }

    @Override
    public User GetErr(User req) throws Exception {
        throw new Exception("exception");
    }

    @Override
    public User GetUser0(String userId, String name) throws Exception {
        return new User(userId, name, 18);
    }

    @Override
    public User GetUser2(int req) throws Exception {
        return new User(String.valueOf(req));
    }

    @Override
    public User[] GetUsers(String[] userIdList) {
        User[] userList = new User[userIdList.length];
        LOG.warn("@userIdList size:" + userIdList.length);

        for (int i = 0; i < userIdList.length; i++) {
            String id = userIdList[i];
            LOG.info("GetUsers(@uid:" + id + ")");
            if (getUserMap().containsKey(id)) {
                userList[i] = getUserMap().get(id);
                LOG.info("id:" + id + ", com.dubbogo.user:" + getUserMap().get(id));
            }
        }
        
        for (String id : userIdList) {
            
        }
        return userList;
    }

    @Override
    public User GetUser3() {
        LOG.info("this is GetUser3 of impl");
        return getUserMap().get("004");
    }

    @Override
    public User getUser(int userCode) {
        return new User(String.valueOf(userCode), "userCode get", 48);
    }

    @Override
    public Gender GetGender(int gender) {
        return 1 == gender ? Gender.WOMAN : Gender.MAN;
    }
}
