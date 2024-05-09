package org.apache.dubbo.sample;

public class GreeterImpl extends DubboGreeterTriple.GreeterImplBase {
    @Override
    public HelloReply sayHello(HelloRequest request) {
        return HelloReply.newBuilder()
                .setMessage("Hello," + request.getName() + "!")
                .build();
    }
}