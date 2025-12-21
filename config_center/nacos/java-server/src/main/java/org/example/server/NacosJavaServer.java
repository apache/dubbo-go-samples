package org.example.server;

import com.alibaba.nacos.api.NacosFactory;
import com.alibaba.nacos.api.config.ConfigService;
import greet.GreetService;
import greet.GreetRequest;
import greet.GreetResponse;
import org.apache.dubbo.config.ConfigCenterConfig;
import org.apache.dubbo.config.ServiceConfig;
import org.apache.dubbo.config.bootstrap.DubboBootstrap;
import java.util.Properties;


public class NacosJavaServer{
    // Nacos server address (used as config center)
    private static final String NACOS_ADDR = "127.0.0.1:8848";
    // Zookeeper address (used as registry)
    private static final String ZK_ADDR = "127.0.0.1:2181";
    // DataId used in Nacos config center
    private static final String DATA_ID = "dubbo-go-samples-configcenter-nacos-java-server";
    // Nacos group
    private static final String GROUP = "dubbo";
    // Configuration content pushed to Nacos config center
    // Dubbo will load these as application / registry / protocol configs
    private static final String CONFIG_CONTENT =
            "dubbo.application.name=dubbo-java-provider\n" +
                    "dubbo.registry.address=zookeeper://" + ZK_ADDR + "\n" +
                    "dubbo.protocol.name=tri\n" +
                    "dubbo.protocol.port=50051\n";
    public static void main(String[] args) throws Exception {
        // Publish Dubbo config to Nacos config center
        Properties properties = new Properties();
        properties.put("serverAddr", NACOS_ADDR);
        ConfigService configService = NacosFactory.createConfigService(properties);
        boolean success=configService.publishConfig(DATA_ID, GROUP, CONFIG_CONTENT);
        if (success) {
            System.out.println(" Succeed to publish config to Nacos, DATA_ID= "+DATA_ID);
        } else {
            System.err.println(" Failed to publish config to Nacos ");
        }

        // Tell Dubbo to use Nacos as config center
        ConfigCenterConfig configCenter = new ConfigCenterConfig();
        configCenter.setAddress("nacos://" + NACOS_ADDR);
        configCenter.setConfigFile(DATA_ID);
        configCenter.setGroup(GROUP);

        // Define service: only interface and implementation
        ServiceConfig<GreetServiceImpl> service = new ServiceConfig<>();
        service.setInterface(GreetService.class);
        service.setRef(new GreetServiceImpl());

        // Start DubboBootstrap with config center and service
        DubboBootstrap.getInstance()
                .application("dubbo-java-provider")
                .configCenter(configCenter)
                .service(service)
                .start()
                .await();
    }

    // Service implementation of GreetService
    static class GreetServiceImpl implements GreetService {
        @Override
        public GreetResponse greet(GreetRequest request) {
            System.out.println("Received request: " + request.getName());
            return GreetResponse.newBuilder()
                    .setGreeting("Hello, this is dubbo java server! I received: " + request.getName())
                    .build();
        }
    }
}