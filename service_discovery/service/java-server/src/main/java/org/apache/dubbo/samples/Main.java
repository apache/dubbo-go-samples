package org.apache.dubbo.samples;

import org.springframework.context.support.ClassPathXmlApplicationContext;

import java.io.IOException;
import java.util.concurrent.CountDownLatch;

//TIP To <b>Run</b> code, press <shortcut actionId="Run"/> or
// click the <icon src="AllIcons.Actions.Execute"/> icon in the gutter.
public class Main {
    public static void main(String[] args) throws IOException, InterruptedException {
        //TIP Press <shortcut actionId="ShowIntentionActions"/> with your caret at the highlighted text
        // to see how IntelliJ IDEA suggests fixing it.
        ClassPathXmlApplicationContext context = new ClassPathXmlApplicationContext(new String[]{"/spring/dubbo.server.xml"});
        context.start();
        System.out.println("dubbo service started");
        new CountDownLatch(1).await();
    }

//    public static void startComplexService() throws InterruptedException {
//        ServiceConfig<GreetService> service = new ServiceConfig<>();
//        service.setInterface(GreetService.class);
//        service.setRef(new GreetServiceImpl());
//        service.setProtocol(new ProtocolConfig(CommonConstants.DUBBO_PROTOCOL, 20010));
//        service.setApplication(new ApplicationConfig("user-info-server"));
//        service.setRegistry(new RegistryConfig("zookeeper://127.0.0.1:2181"));
//        service.export();
//        System.out.println("dubbo service started");
//    }
//
//    public static void startWrapperArrayClassProvider() throws InterruptedException {
//        ServiceConfig<WrapperArrayClassProvider> service = new ServiceConfig<>();
//        service.setInterface(WrapperArrayClassProvider.class);
//        service.setRef(new WrapperArrayClassProviderImpl());
//        service.setProtocol(new ProtocolConfig(CommonConstants.DUBBO_PROTOCOL, 20010));
//        service.setApplication(new ApplicationConfig("user-info-server"));
//        service.setRegistry(new RegistryConfig("zookeeper://127.0.0.1:2181"));
//        service.export();
//        System.out.println("dubbo service started");
//    }
}