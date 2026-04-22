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

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.*;

public class UserProviderImpl implements UserProvider {
    private static final Logger LOG = LoggerFactory.getLogger("userLogger");
    private Map<String, User> userMap = new HashMap<String, User>();

    UserProviderImpl() {
        userMap.put("001", new User("001", "other-zhangsan", 23, new Date(1998, 1, 2, 3, 4, 5), Gender.MAN));
        userMap.put("002", new User("002", "other-lisi", 25, new Date(1996, 1, 2, 3, 4, 5), Gender.MAN));
        userMap.put("003", new User("003", "other-lily", 28, new Date(1993, 1, 2, 3, 4, 5), Gender.WOMAN));
        userMap.put("004", new User("004", "other-lisa", 36, new Date(1985, 1, 2, 3, 4, 5), Gender.WOMAN));
    }

    public User GetUser1(String userId) {
        LOG.info("req:" + userId);
        User rsp = new User(userId, "Joe", 48, new Date(), Gender.MAN);
        LOG.info("rsp:" + rsp);
        return rsp;
    }

    public User GetUser2(String userId, String name) {
        LOG.info("req:" + userId + ", " + name);
        User rsp = new User(userId, name, 48, new Date(), Gender.MAN);
        LOG.info("rsp:" + rsp);
        return rsp;
    }

    public User GetUser3(int userCode) {
        LOG.info("req:" + userCode);
        User rsp = new User(String.valueOf(userCode), "Alex Stocks", 18, new Date(), Gender.MAN);
        LOG.info("rsp:" + rsp);
        return rsp;
    }

    public User GetUser4(int userCode, String name) {
        LOG.info("req:" + userCode + ", " + name);
        User rsp = new User(String.valueOf(userCode), name, 18, new Date(), Gender.MAN);
        LOG.info("rsp:" + rsp);
        return rsp;
    }

    public User GetOneUser() {
        return new User("1000", "xavierniu", 24, new Date(), Gender.MAN);
    }

    @Override
    public List<User> GetUsers(String[] userIdList) {
        LOG.info("req:" + Arrays.toString(userIdList));
        List<User> userList = new ArrayList<User>();
        for (String i : userIdList) {
            userList.add(userMap.get(i));
        }
        return userList;
    }

    @Override
    public Map<String, User> GetUsersMap(String[] userIdList) {
        LOG.info("req:" + Arrays.toString(userIdList));
        Map<String, User> users = new HashMap<String, User>();
        for (String i : userIdList) {
            users.put(i, userMap.get(i));
        }
        return users;
    }

    public User QueryUser(User user) {
        LOG.info("req1:" + user);
        User rsp = new User(user.getId(), user.getName(), user.getAge(), new Date(), user.getSex());
        LOG.info("rsp1:" + rsp);
        return rsp;
    }

    @Override
    public User[] QueryUsers(User[] users) {
        return users;
    }

    public Map<String, User> QueryAll() {
        Map<String, User> users = new HashMap<String, User>();
        users.put("001", new User("001", "Joe", 18, new Date(), Gender.MAN));
        users.put("002", new User("002", "Wen", 20, new Date(), Gender.MAN));
        return users;
    }
}