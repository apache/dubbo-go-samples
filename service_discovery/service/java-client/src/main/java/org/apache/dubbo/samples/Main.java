package org.apache.dubbo.samples;

import org.apache.dubbo.common.constants.CommonConstants;
import org.apache.dubbo.config.ApplicationConfig;
import org.apache.dubbo.config.ReferenceConfig;
import org.apache.dubbo.config.RegistryConfig;
import org.apache.dubbo.samples.proto.GreetRequest;
import org.apache.dubbo.samples.proto.GreetResponse;
import org.apache.dubbo.samples.proto.GreetService;
import org.springframework.context.support.ClassPathXmlApplicationContext;

import java.io.IOException;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.TimeUnit;

//TIP To <b>Run</b> code, press <shortcut actionId="Run"/> or
// click the <icon src="AllIcons.Actions.Execute"/> icon in the gutter.
public class Main {
    // Define a private variable (Required in Spring)
//    private static GreetService greetService;
    public static void main(String[] args) throws InterruptedException, IOException {
//            ClassPathXmlApplicationContext context = new ClassPathXmlApplicationContext(new String[]{"/spring/dubbo.consumer.xml"});
//            context.start();
//            greetService = (GreetService)context.getBean("greetService");
//            GreetRequest req = GreetRequest.getDefaultInstance();
//            GreetResponse greet = greetService.greet(req);
//            System.out.println(greet.getGreeting());
//            new CountDownLatch(1).await();
        ReferenceConfig<GreetService> ref = new ReferenceConfig<>();
        ref.setInterface(GreetService.class);
        ref.setCheck(false);
        ref.setProtocol(CommonConstants.TRIPLE);
        ref.setLazy(true);
        ref.setTimeout(100000);
        ref.setApplication(new ApplicationConfig("demo-consumer"));
        ref.setRegistry(new RegistryConfig("nacos://127.0.0.1:8848"));
        final GreetService greetService = ref.get();

        System.out.println("dubbo ref started");
        GreetRequest req = GreetRequest.newBuilder().setName("laurence").build();
        try {
            final GreetResponse reply =greetService.greet(req);
            TimeUnit.SECONDS.sleep(1);
            System.out.println("Reply:" + reply);
        } catch (Throwable t) {
            t.printStackTrace();
        }
        System.in.read();
    }

}