/*
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package org.apache.dubbo.sample;

import java.util.List;

public interface UserProvider {
    User[] GetUsers(String[] req) throws Exception;

    User GetErr(User req) throws Exception;

    User GetUser(User req) throws Exception;

    User GetUser0(String userId, String name) throws Exception;

    User GetUser2(int req) throws Exception;

    User GetUser3() throws Exception;

    Gender GetGender(Integer userId) throws Exception;

    User getUser(int req) throws Exception;
}
