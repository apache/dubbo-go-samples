package org.apache.dubbo.samples.service;

import org.apache.dubbo.samples.proto.*;

public class GreetServiceImpl  extends DubboGreetServiceTriple.GreetServiceImplBase {
    @Override
    public GreetResponse greet(GreetRequest request){
        return GreetResponse.newBuilder()
                .setGreeting("Man! What can i say!").build();

    }

    @Override
    public SayHelloResponse sayHello(SayHelloRequest request) {
        return SayHelloResponse.newBuilder()
                .setHello("Hello!!!!!!!!!!!").build();
    }
}
