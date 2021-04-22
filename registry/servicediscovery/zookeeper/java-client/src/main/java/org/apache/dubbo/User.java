package org.apache.dubbo;

import java.io.Serializable;
import java.util.Date;

/**
 *@author Patrick Jiang
 *@Description
 *@date 2021/4/19 5:04 下午
 */
public class User implements Serializable {
    
    public User(String id, String name, int age) {
        this.id = id;
        this.name = name;
        this.age = age;
    }
    
    private String id;
    
    private String name;
    
    private int age;
    
    private Date time = new Date();
    
    public String getId() {
        return id;
    }
    
    public void setId(String id) {
        this.id = id;
    }
    
    public String getName() {
        return name;
    }
    
    public void setName(String name) {
        this.name = name;
    }
    
    public int getAge() {
        return age;
    }
    
    public void setAge(int age) {
        this.age = age;
    }
    
    public Date getTime() {
        return time;
    }
    
    public void setTime(Date time) {
        this.time = time;
    }
    
    @Override
    public String toString() {
        return "User{" +
                "id='" + id + '\'' +
                ", name='" + name + '\'' +
                ", age=" + age +
                ", time=" + time +
                '}';
    }
}
