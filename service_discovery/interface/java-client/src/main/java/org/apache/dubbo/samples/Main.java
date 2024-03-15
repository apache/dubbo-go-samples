package org.apache.dubbo.samples;

import org.apache.dubbo.samples.proto.GreetRequest;
import org.apache.dubbo.samples.proto.GreetResponse;
import org.apache.dubbo.samples.proto.GreetService;
import org.springframework.context.support.ClassPathXmlApplicationContext;

import java.io.IOException;
import java.util.concurrent.CountDownLatch;

//TIP To <b>Run</b> code, press <shortcut actionId="Run"/> or
// click the <icon src="AllIcons.Actions.Execute"/> icon in the gutter.
public class Main {
    // Define a private variable (Required in Spring)
    private static GreetService greetService;
    public static void main(String[] args) throws InterruptedException, IOException {
        System.setProperty("dubbo.application.service-discovery.migration", "FORCE_INTERFACE");
        ClassPathXmlApplicationContext context = new ClassPathXmlApplicationContext(new String[]{"/spring/dubbo.consumer.xml"});
//      context.start();
        greetService = (GreetService)context.getBean("GreetService");
        GreetRequest req = GreetRequest.newBuilder().setName("Mamba").build();
        System.out.println("dubbo ref started");
        GreetResponse greet = greetService.greet(req);
        System.out.println("Greeting:" + greet.getGreeting() + "!!!!!!!!!!!!");
//        new CountDownLatch(1).await();
    }

}