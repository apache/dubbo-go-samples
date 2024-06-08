package org.apache.dubbo.sample;

import org.apache.dubbo.common.constants.CommonConstants;
import org.apache.dubbo.config.ApplicationConfig;
import org.apache.dubbo.config.ProtocolConfig;
import org.apache.dubbo.config.ReferenceConfig;
import org.apache.dubbo.config.RegistryConfig;
import org.apache.dubbo.config.bootstrap.DubboBootstrap;



public class Consumer {
    public static void main(String[] args) {
        DubboBootstrap bootstrap = DubboBootstrap.getInstance();
        ReferenceConfig<Greeter> ref = new ReferenceConfig<>();
        ref.setInterface(Greeter.class);
        ref.setTimeout(3000);
        ref.setUrl("tri://localhost:50052");
        bootstrap.application(new ApplicationConfig("dubbo-test"))
                .reference(ref)
                .start();

        Greeter greeter = ref.get();
        HelloRequest request = HelloRequest.newBuilder().setName("Demo").build();
        HelloReply reply = greeter.sayHello(request);
        System.out.println("Received reply:" + reply);
    }
}
