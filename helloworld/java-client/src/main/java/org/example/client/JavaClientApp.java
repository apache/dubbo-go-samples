package org.example.client;

import greet.GreetRequest;
import greet.GreetResponse;
import greet.GreetServiceGrpc;

import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;

public class JavaClientApp {

    public static void main(String[] args) {

        ManagedChannel channel = ManagedChannelBuilder
                .forAddress("127.0.0.1", 20000)
                .usePlaintext()
                .build();

        GreetServiceGrpc.GreetServiceBlockingStub client =
                GreetServiceGrpc.newBlockingStub(channel);

        GreetRequest req = GreetRequest.newBuilder()
                .setName("Java Client")
                .build();

        GreetResponse resp = client.greet(req);

        System.out.println("Response: " + resp.getGreeting());

        channel.shutdown();
    }
}
