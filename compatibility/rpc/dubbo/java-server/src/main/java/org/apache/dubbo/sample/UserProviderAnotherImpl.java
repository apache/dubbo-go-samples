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

import java.util.*;

public class UserProviderAnotherImpl extends UserProviderImpl {
    private Map<String, User> userMap = new HashMap<String, User>();

    public UserProviderAnotherImpl() {
        // userMap.put("001", new User("001", "other-zhangsan", 18, new Date(1998-1900, 1, 2, 3, 4, 5), Gender.MAN));
        userMap.put("001", new User("001", "other-zhangsan", 18, new Date(0x12345678), Gender.MAN));
        userMap.put("002", new User("002", "other-lisi", 20, new Date(1996-1900, 1, 2, 3, 4, 5), Gender.MAN));
        userMap.put("003", new User("003", "other-lily", 23, new Date(1993-1900, 1, 2, 3, 4, 5), Gender.WOMAN));
        userMap.put("004", new User("004", "other-lisa", 32, new Date(1985-1900, 1, 2, 3, 4, 5), Gender.WOMAN));
    }

    @Override
    public Map<String, User> getUserMap() {
        return userMap;
    }

}
