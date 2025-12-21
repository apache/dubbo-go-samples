package org.example.client;

import greet.GreetService;
import greet.GreetRequest;
import greet.GreetResponse;
import org.apache.curator.framework.CuratorFramework;
import org.apache.curator.framework.CuratorFrameworkFactory;
import org.apache.curator.retry.ExponentialBackoffRetry;
import org.apache.dubbo.config.ConfigCenterConfig;
import org.apache.dubbo.config.ReferenceConfig;
import org.apache.dubbo.config.bootstrap.DubboBootstrap;
import java.nio.charset.StandardCharsets;

public class ZookeeperJavaClient {
    // Zookeeper address (used as both registry and config center)
    private static final String ZK_ADDR = "127.0.0.1:2181";
    // DataId used in Zookeeper config center
    private static final String DATA_ID = "dubbo-go-samples-configcenter-zookeeper-java-client";
    // Config center group
    private static final String GROUP = "dubbogo";
    // Configuration content pushed to Zookeeper config center
    // Dubbo will load these as application / registry / protocol configs
    private static final String CONFIG_CONTENT =
            "dubbo.application.name=dubbo-java-consumer\n" +
                    "dubbo.registry.address=zookeeper://" + ZK_ADDR + "\n" +
                    "dubbo.protocol.name=tri\n";

    public static void main(String[] args) throws Exception {
        // Publish Dubbo config to Zookeeper config center
        String configPath = "/dubbo/config/" + GROUP + "/" + DATA_ID;
        CuratorFramework client = CuratorFrameworkFactory.builder()
                .connectString(ZK_ADDR)
                .retryPolicy(new ExponentialBackoffRetry(1000, 3))
                .build();
        client.start();
        try {
            byte[] data = CONFIG_CONTENT.getBytes(StandardCharsets.UTF_8);
            if (client.checkExists().forPath(configPath) == null) {
                client.create()
                        .creatingParentsIfNeeded()
                        .forPath(configPath, data);
                System.out.println("Created client config in ZK, path=" + configPath);
            } else {
                client.setData().forPath(configPath, data);
                System.out.println("Updated client config in ZK, path=" + configPath);
            }
        } finally {
            client.close();
        }

        // Tell Dubbo to use Zookeeper as config center
        ConfigCenterConfig configCenter = new ConfigCenterConfig();
        configCenter.setAddress("zookeeper://" + ZK_ADDR);
        configCenter.setGroup(GROUP);
        configCenter.setConfigFile(DATA_ID);

        // Create reference config, only bind interface here
        ReferenceConfig<GreetService> reference = new ReferenceConfig<>();
        reference.setInterface(GreetService.class);

        // Start DubboBootstrap and load reference from config center
        DubboBootstrap bootstrap = DubboBootstrap.getInstance()
                .application("dubbo-java-consumer")
                .configCenter(configCenter)
                .reference(reference)
                .start();

        //Get remote service proxy
        GreetService greetService = reference.get();
        //Build request and call remote method
        GreetRequest request = GreetRequest.newBuilder()
                .setName("Hello, this is dubbo java client!")
                .build();

        GreetResponse response = greetService.greet(request);
        System.out.println("Server response: " + response.getGreeting());

        System.exit(0);
    }
}