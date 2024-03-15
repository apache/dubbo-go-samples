package org.apache.dubbo.samples;

import org.springframework.context.support.ClassPathXmlApplicationContext;

import java.util.concurrent.CountDownLatch;

//TIP To <b>Run</b> code, press <shortcut actionId="Run"/> or
// click the <icon src="AllIcons.Actions.Execute"/> icon in the gutter.
public class Main {
    public static void main(String[] args) throws InterruptedException {
        System.setProperty("dubbo.application.register-mode", "interface");
        ClassPathXmlApplicationContext context = new ClassPathXmlApplicationContext(new String[]{"/spring/dubbo.server.xml"});
        context.start();
        System.out.println("dubbo service started");
        new CountDownLatch(1).await();
    }
}