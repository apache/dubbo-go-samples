package org.apache.dubbo.samples;

import org.springframework.context.support.ClassPathXmlApplicationContext;

import java.io.IOException;
import java.util.concurrent.CountDownLatch;

public class Main {
    public static void main(String[] args) throws IOException, InterruptedException {
        System.setProperty("dubbo.application.register-mode", "instance");
        ClassPathXmlApplicationContext context = new ClassPathXmlApplicationContext(new String[]{"/spring/dubbo.server.xml"});
        context.start();
        System.out.println("dubbo service started");
        new CountDownLatch(1).await();
    }
}