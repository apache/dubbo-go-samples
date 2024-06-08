/*
* Licensed to the Apache Software Foundation (ASF) under one or more
* contributor license agreements.  See the NOTICE file distributed with
* this work for additional information regarding copyright ownership.
* The ASF licenses this file to You under the Apache License, Version 2.0
* (the "License"); you may not use this file except in compliance with
* the License.  You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*/

package org.apache.dubbo.samples.proto;

import org.apache.dubbo.common.stream.StreamObserver;
import org.apache.dubbo.common.URL;
import org.apache.dubbo.rpc.Invoker;
import org.apache.dubbo.rpc.PathResolver;
import org.apache.dubbo.rpc.RpcException;
import org.apache.dubbo.rpc.ServerService;
import org.apache.dubbo.rpc.TriRpcStatus;
import org.apache.dubbo.rpc.model.MethodDescriptor;
import org.apache.dubbo.rpc.model.ServiceDescriptor;
import org.apache.dubbo.rpc.model.StubMethodDescriptor;
import org.apache.dubbo.rpc.model.StubServiceDescriptor;
import org.apache.dubbo.rpc.stub.BiStreamMethodHandler;
import org.apache.dubbo.rpc.stub.ServerStreamMethodHandler;
import org.apache.dubbo.rpc.stub.StubInvocationUtil;
import org.apache.dubbo.rpc.stub.StubInvoker;
import org.apache.dubbo.rpc.stub.StubMethodHandler;
import org.apache.dubbo.rpc.stub.StubSuppliers;
import org.apache.dubbo.rpc.stub.UnaryStubMethodHandler;

import com.google.protobuf.Message;

import java.util.HashMap;
import java.util.Map;
import java.util.function.BiConsumer;
import java.util.concurrent.CompletableFuture;

public final class DubboGreetServiceTriple {

    public static final String SERVICE_NAME = GreetService.SERVICE_NAME;

    private static final StubServiceDescriptor serviceDescriptor = new StubServiceDescriptor(SERVICE_NAME,GreetService.class);

    static {
        org.apache.dubbo.rpc.protocol.tri.service.SchemaDescriptorRegistry.addSchemaDescriptor(SERVICE_NAME,GreetProto.getDescriptor());
        StubSuppliers.addSupplier(SERVICE_NAME, DubboGreetServiceTriple::newStub);
        StubSuppliers.addSupplier(GreetService.JAVA_SERVICE_NAME,  DubboGreetServiceTriple::newStub);
        StubSuppliers.addDescriptor(SERVICE_NAME, serviceDescriptor);
        StubSuppliers.addDescriptor(GreetService.JAVA_SERVICE_NAME, serviceDescriptor);
    }

    @SuppressWarnings("all")
    public static GreetService newStub(Invoker<?> invoker) {
        return new GreetServiceStub((Invoker<GreetService>)invoker);
    }

    private static final StubMethodDescriptor greetMethod = new StubMethodDescriptor("Greet",
    org.apache.dubbo.samples.proto.GreetRequest.class, org.apache.dubbo.samples.proto.GreetResponse.class, MethodDescriptor.RpcType.UNARY,
    obj -> ((Message) obj).toByteArray(), obj -> ((Message) obj).toByteArray(), org.apache.dubbo.samples.proto.GreetRequest::parseFrom,
    org.apache.dubbo.samples.proto.GreetResponse::parseFrom);

    private static final StubMethodDescriptor greetAsyncMethod = new StubMethodDescriptor("Greet",
    org.apache.dubbo.samples.proto.GreetRequest.class, java.util.concurrent.CompletableFuture.class, MethodDescriptor.RpcType.UNARY,
    obj -> ((Message) obj).toByteArray(), obj -> ((Message) obj).toByteArray(), org.apache.dubbo.samples.proto.GreetRequest::parseFrom,
    org.apache.dubbo.samples.proto.GreetResponse::parseFrom);

    private static final StubMethodDescriptor greetProxyAsyncMethod = new StubMethodDescriptor("GreetAsync",
    org.apache.dubbo.samples.proto.GreetRequest.class, org.apache.dubbo.samples.proto.GreetResponse.class, MethodDescriptor.RpcType.UNARY,
    obj -> ((Message) obj).toByteArray(), obj -> ((Message) obj).toByteArray(), org.apache.dubbo.samples.proto.GreetRequest::parseFrom,
    org.apache.dubbo.samples.proto.GreetResponse::parseFrom);
    private static final StubMethodDescriptor sayHelloMethod = new StubMethodDescriptor("SayHello",
    org.apache.dubbo.samples.proto.SayHelloRequest.class, org.apache.dubbo.samples.proto.SayHelloResponse.class, MethodDescriptor.RpcType.UNARY,
    obj -> ((Message) obj).toByteArray(), obj -> ((Message) obj).toByteArray(), org.apache.dubbo.samples.proto.SayHelloRequest::parseFrom,
    org.apache.dubbo.samples.proto.SayHelloResponse::parseFrom);

    private static final StubMethodDescriptor sayHelloAsyncMethod = new StubMethodDescriptor("SayHello",
    org.apache.dubbo.samples.proto.SayHelloRequest.class, java.util.concurrent.CompletableFuture.class, MethodDescriptor.RpcType.UNARY,
    obj -> ((Message) obj).toByteArray(), obj -> ((Message) obj).toByteArray(), org.apache.dubbo.samples.proto.SayHelloRequest::parseFrom,
    org.apache.dubbo.samples.proto.SayHelloResponse::parseFrom);

    private static final StubMethodDescriptor sayHelloProxyAsyncMethod = new StubMethodDescriptor("SayHelloAsync",
    org.apache.dubbo.samples.proto.SayHelloRequest.class, org.apache.dubbo.samples.proto.SayHelloResponse.class, MethodDescriptor.RpcType.UNARY,
    obj -> ((Message) obj).toByteArray(), obj -> ((Message) obj).toByteArray(), org.apache.dubbo.samples.proto.SayHelloRequest::parseFrom,
    org.apache.dubbo.samples.proto.SayHelloResponse::parseFrom);




    static{
        serviceDescriptor.addMethod(greetMethod);
        serviceDescriptor.addMethod(greetProxyAsyncMethod);
        serviceDescriptor.addMethod(sayHelloMethod);
        serviceDescriptor.addMethod(sayHelloProxyAsyncMethod);
    }

    public static class GreetServiceStub implements GreetService{
        private final Invoker<GreetService> invoker;

        public GreetServiceStub(Invoker<GreetService> invoker) {
            this.invoker = invoker;
        }

        @Override
        public org.apache.dubbo.samples.proto.GreetResponse greet(org.apache.dubbo.samples.proto.GreetRequest request){
            return StubInvocationUtil.unaryCall(invoker, greetMethod, request);
        }

        public CompletableFuture<org.apache.dubbo.samples.proto.GreetResponse> greetAsync(org.apache.dubbo.samples.proto.GreetRequest request){
            return StubInvocationUtil.unaryCall(invoker, greetAsyncMethod, request);
        }

        public void greet(org.apache.dubbo.samples.proto.GreetRequest request, StreamObserver<org.apache.dubbo.samples.proto.GreetResponse> responseObserver){
            StubInvocationUtil.unaryCall(invoker, greetMethod , request, responseObserver);
        }
        @Override
        public org.apache.dubbo.samples.proto.SayHelloResponse sayHello(org.apache.dubbo.samples.proto.SayHelloRequest request){
            return StubInvocationUtil.unaryCall(invoker, sayHelloMethod, request);
        }

        public CompletableFuture<org.apache.dubbo.samples.proto.SayHelloResponse> sayHelloAsync(org.apache.dubbo.samples.proto.SayHelloRequest request){
            return StubInvocationUtil.unaryCall(invoker, sayHelloAsyncMethod, request);
        }

        public void sayHello(org.apache.dubbo.samples.proto.SayHelloRequest request, StreamObserver<org.apache.dubbo.samples.proto.SayHelloResponse> responseObserver){
            StubInvocationUtil.unaryCall(invoker, sayHelloMethod , request, responseObserver);
        }



    }

    public static abstract class GreetServiceImplBase implements GreetService, ServerService<GreetService> {

        private <T, R> BiConsumer<T, StreamObserver<R>> syncToAsync(java.util.function.Function<T, R> syncFun) {
            return new BiConsumer<T, StreamObserver<R>>() {
                @Override
                public void accept(T t, StreamObserver<R> observer) {
                    try {
                        R ret = syncFun.apply(t);
                        observer.onNext(ret);
                        observer.onCompleted();
                    } catch (Throwable e) {
                        observer.onError(e);
                    }
                }
            };
        }

        @Override
        public CompletableFuture<org.apache.dubbo.samples.proto.GreetResponse> greetAsync(org.apache.dubbo.samples.proto.GreetRequest request){
                return CompletableFuture.completedFuture(greet(request));
        }
        @Override
        public CompletableFuture<org.apache.dubbo.samples.proto.SayHelloResponse> sayHelloAsync(org.apache.dubbo.samples.proto.SayHelloRequest request){
                return CompletableFuture.completedFuture(sayHello(request));
        }

        /**
        * This server stream type unary method is <b>only</b> used for generated stub to support async unary method.
        * It will not be called if you are NOT using Dubbo3 generated triple stub and <b>DO NOT</b> implement this method.
        */
        public void greet(org.apache.dubbo.samples.proto.GreetRequest request, StreamObserver<org.apache.dubbo.samples.proto.GreetResponse> responseObserver){
            greetAsync(request).whenComplete((r, t) -> {
                if (t != null) {
                    responseObserver.onError(t);
                } else {
                    responseObserver.onNext(r);
                    responseObserver.onCompleted();
                }
            });
        }
        public void sayHello(org.apache.dubbo.samples.proto.SayHelloRequest request, StreamObserver<org.apache.dubbo.samples.proto.SayHelloResponse> responseObserver){
            sayHelloAsync(request).whenComplete((r, t) -> {
                if (t != null) {
                    responseObserver.onError(t);
                } else {
                    responseObserver.onNext(r);
                    responseObserver.onCompleted();
                }
            });
        }

        @Override
        public final Invoker<GreetService> getInvoker(URL url) {
            PathResolver pathResolver = url.getOrDefaultFrameworkModel()
            .getExtensionLoader(PathResolver.class)
            .getDefaultExtension();
            Map<String,StubMethodHandler<?, ?>> handlers = new HashMap<>();

            pathResolver.addNativeStub( "/" + SERVICE_NAME + "/Greet");
            pathResolver.addNativeStub( "/" + SERVICE_NAME + "/GreetAsync");
            // for compatibility
            pathResolver.addNativeStub( "/" + JAVA_SERVICE_NAME + "/Greet");
            pathResolver.addNativeStub( "/" + JAVA_SERVICE_NAME + "/GreetAsync");

            pathResolver.addNativeStub( "/" + SERVICE_NAME + "/SayHello");
            pathResolver.addNativeStub( "/" + SERVICE_NAME + "/SayHelloAsync");
            // for compatibility
            pathResolver.addNativeStub( "/" + JAVA_SERVICE_NAME + "/SayHello");
            pathResolver.addNativeStub( "/" + JAVA_SERVICE_NAME + "/SayHelloAsync");


            BiConsumer<org.apache.dubbo.samples.proto.GreetRequest, StreamObserver<org.apache.dubbo.samples.proto.GreetResponse>> greetFunc = this::greet;
            handlers.put(greetMethod.getMethodName(), new UnaryStubMethodHandler<>(greetFunc));
            BiConsumer<org.apache.dubbo.samples.proto.GreetRequest, StreamObserver<org.apache.dubbo.samples.proto.GreetResponse>> greetAsyncFunc = syncToAsync(this::greet);
            handlers.put(greetProxyAsyncMethod.getMethodName(), new UnaryStubMethodHandler<>(greetAsyncFunc));
            BiConsumer<org.apache.dubbo.samples.proto.SayHelloRequest, StreamObserver<org.apache.dubbo.samples.proto.SayHelloResponse>> sayHelloFunc = this::sayHello;
            handlers.put(sayHelloMethod.getMethodName(), new UnaryStubMethodHandler<>(sayHelloFunc));
            BiConsumer<org.apache.dubbo.samples.proto.SayHelloRequest, StreamObserver<org.apache.dubbo.samples.proto.SayHelloResponse>> sayHelloAsyncFunc = syncToAsync(this::sayHello);
            handlers.put(sayHelloProxyAsyncMethod.getMethodName(), new UnaryStubMethodHandler<>(sayHelloAsyncFunc));




            return new StubInvoker<>(this, url, GreetService.class, handlers);
        }


        @Override
        public org.apache.dubbo.samples.proto.GreetResponse greet(org.apache.dubbo.samples.proto.GreetRequest request){
            throw unimplementedMethodException(greetMethod);
        }

        @Override
        public org.apache.dubbo.samples.proto.SayHelloResponse sayHello(org.apache.dubbo.samples.proto.SayHelloRequest request){
            throw unimplementedMethodException(sayHelloMethod);
        }





        @Override
        public final ServiceDescriptor getServiceDescriptor() {
            return serviceDescriptor;
        }
        private RpcException unimplementedMethodException(StubMethodDescriptor methodDescriptor) {
            return TriRpcStatus.UNIMPLEMENTED.withDescription(String.format("Method %s is unimplemented",
                "/" + serviceDescriptor.getInterfaceName() + "/" + methodDescriptor.getMethodName())).asException();
        }
    }

}
