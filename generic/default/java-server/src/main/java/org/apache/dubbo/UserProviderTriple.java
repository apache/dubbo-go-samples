package org.apache.dubbo;

import java.util.ArrayList;
import java.util.List;
import java.util.Map;

/**
 * Created by liuyuecai on 2021/11/23.
 */
public interface UserProviderTriple {

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
