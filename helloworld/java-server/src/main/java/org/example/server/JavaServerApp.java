package org.example.server;

import greet.GreetRequest;
import greet.GreetResponse;
import greet.GreetServiceGrpc;

import io.grpc.Server;
import io.grpc.ServerBuilder;
import io.grpc.stub.StreamObserver;

public class JavaServerApp {

    public static void main(String[] args) throws Exception {

        Server server = ServerBuilder.forPort(20000)
                .addService(new GreetServiceGrpc.GreetServiceImplBase() {

                    @Override
                    public void greet(GreetRequest request,
                                      StreamObserver<GreetResponse> respObserver) {

                        String name = request.getName();

                        GreetResponse resp = GreetResponse.newBuilder()
                                .setGreeting("Hello from Java Server, " + name)
                                .build();

                        respObserver.onNext(resp);
                        respObserver.onCompleted();
                    }
                })
                .build();

        server.start();
        System.out.println("Java Triple Server started on port 20000");
        server.awaitTermination();
    }
}
