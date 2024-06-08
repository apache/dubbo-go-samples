
package org.apache.dubbo.sample;

import org.apache.dubbo.common.constants.CommonConstants;
import org.apache.dubbo.config.ApplicationConfig;
import org.apache.dubbo.config.ProtocolConfig;
import org.apache.dubbo.config.ServiceConfig;
import org.apache.dubbo.config.bootstrap.DubboBootstrap;

import java.io.IOException;

public class Provider {
    public static void main(String[] args) throws IOException {
        ServiceConfig<Greeter> service =  new ServiceConfig<>();
        service.setInterface(Greeter.class);
        service.setRef(new GreeterImpl());

        DubboBootstrap bootstrap = DubboBootstrap.getInstance();
        bootstrap.application(new ApplicationConfig("java-go-sample-server"))
                .protocol(new ProtocolConfig(CommonConstants.TRIPLE,50052))
                .service(service)
                .start();
        System.out.println("Dubbo triple java server started");
        System.in.read();
    }
}