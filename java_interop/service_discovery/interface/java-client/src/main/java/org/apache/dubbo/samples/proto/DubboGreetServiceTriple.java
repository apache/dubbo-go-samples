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

import com.google.protobuf.Message;
import org.apache.dubbo.common.URL;
import org.apache.dubbo.common.stream.StreamObserver;
import org.apache.dubbo.rpc.*;
import org.apache.dubbo.rpc.model.MethodDescriptor;
import org.apache.dubbo.rpc.model.ServiceDescriptor;
import org.apache.dubbo.rpc.model.StubMethodDescriptor;
import org.apache.dubbo.rpc.model.StubServiceDescriptor;
import org.apache.dubbo.rpc.stub.*;

import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.CompletableFuture;
import java.util.function.BiConsumer;

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
    GreetRequest.class, GreetResponse.class, MethodDescriptor.RpcType.UNARY,
    obj -> ((Message) obj).toByteArray(), obj -> ((Message) obj).toByteArray(), GreetRequest::parseFrom,
    GreetResponse::parseFrom);

    private static final StubMethodDescriptor greetAsyncMethod = new StubMethodDescriptor("Greet",
    GreetRequest.class, CompletableFuture.class, MethodDescriptor.RpcType.UNARY,
    obj -> ((Message) obj).toByteArray(), obj -> ((Message) obj).toByteArray(), GreetRequest::parseFrom,
    GreetResponse::parseFrom);

    private static final StubMethodDescriptor greetProxyAsyncMethod = new StubMethodDescriptor("GreetAsync",
    GreetRequest.class, GreetResponse.class, MethodDescriptor.RpcType.UNARY,
    obj -> ((Message) obj).toByteArray(), obj -> ((Message) obj).toByteArray(), GreetRequest::parseFrom,
    GreetResponse::parseFrom);
    private static final StubMethodDescriptor sayHelloMethod = new StubMethodDescriptor("SayHello",
    SayHelloRequest.class, SayHelloResponse.class, MethodDescriptor.RpcType.UNARY,
    obj -> ((Message) obj).toByteArray(), obj -> ((Message) obj).toByteArray(), SayHelloRequest::parseFrom,
    SayHelloResponse::parseFrom);

    private static final StubMethodDescriptor sayHelloAsyncMethod = new StubMethodDescriptor("SayHello",
    SayHelloRequest.class, CompletableFuture.class, MethodDescriptor.RpcType.UNARY,
    obj -> ((Message) obj).toByteArray(), obj -> ((Message) obj).toByteArray(), SayHelloRequest::parseFrom,
    SayHelloResponse::parseFrom);

    private static final StubMethodDescriptor sayHelloProxyAsyncMethod = new StubMethodDescriptor("SayHelloAsync",
    SayHelloRequest.class, SayHelloResponse.class, MethodDescriptor.RpcType.UNARY,
    obj -> ((Message) obj).toByteArray(), obj -> ((Message) obj).toByteArray(), SayHelloRequest::parseFrom,
    SayHelloResponse::parseFrom);




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
        public GreetResponse greet(GreetRequest request){
            return StubInvocationUtil.unaryCall(invoker, greetMethod, request);
        }

        public CompletableFuture<GreetResponse> greetAsync(GreetRequest request){
            return StubInvocationUtil.unaryCall(invoker, greetAsyncMethod, request);
        }

        public void greet(GreetRequest request, StreamObserver<GreetResponse> responseObserver){
            StubInvocationUtil.unaryCall(invoker, greetMethod , request, responseObserver);
        }
        @Override
        public SayHelloResponse sayHello(SayHelloRequest request){
            return StubInvocationUtil.unaryCall(invoker, sayHelloMethod, request);
        }

        public CompletableFuture<SayHelloResponse> sayHelloAsync(SayHelloRequest request){
            return StubInvocationUtil.unaryCall(invoker, sayHelloAsyncMethod, request);
        }

        public void sayHello(SayHelloRequest request, StreamObserver<SayHelloResponse> responseObserver){
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
        public CompletableFuture<GreetResponse> greetAsync(GreetRequest request){
                return CompletableFuture.completedFuture(greet(request));
        }
        @Override
        public CompletableFuture<SayHelloResponse> sayHelloAsync(SayHelloRequest request){
                return CompletableFuture.completedFuture(sayHello(request));
        }

        /**
        * This server stream type unary method is <b>only</b> used for generated stub to support async unary method.
        * It will not be called if you are NOT using Dubbo3 generated triple stub and <b>DO NOT</b> implement this method.
        */
        public void greet(GreetRequest request, StreamObserver<GreetResponse> responseObserver){
            greetAsync(request).whenComplete((r, t) -> {
                if (t != null) {
                    responseObserver.onError(t);
                } else {
                    responseObserver.onNext(r);
                    responseObserver.onCompleted();
                }
            });
        }
        public void sayHello(SayHelloRequest request, StreamObserver<SayHelloResponse> responseObserver){
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

            pathResolver.addNativeStub( "/" + SERVICE_NAME + "/Greet" );
            pathResolver.addNativeStub( "/" + SERVICE_NAME + "/GreetAsync" );
            pathResolver.addNativeStub( "/" + SERVICE_NAME + "/SayHello" );
            pathResolver.addNativeStub( "/" + SERVICE_NAME + "/SayHelloAsync" );

            BiConsumer<GreetRequest, StreamObserver<GreetResponse>> greetFunc = this::greet;
            handlers.put(greetMethod.getMethodName(), new UnaryStubMethodHandler<>(greetFunc));
            BiConsumer<GreetRequest, StreamObserver<GreetResponse>> greetAsyncFunc = syncToAsync(this::greet);
            handlers.put(greetProxyAsyncMethod.getMethodName(), new UnaryStubMethodHandler<>(greetAsyncFunc));
            BiConsumer<SayHelloRequest, StreamObserver<SayHelloResponse>> sayHelloFunc = this::sayHello;
            handlers.put(sayHelloMethod.getMethodName(), new UnaryStubMethodHandler<>(sayHelloFunc));
            BiConsumer<SayHelloRequest, StreamObserver<SayHelloResponse>> sayHelloAsyncFunc = syncToAsync(this::sayHello);
            handlers.put(sayHelloProxyAsyncMethod.getMethodName(), new UnaryStubMethodHandler<>(sayHelloAsyncFunc));




            return new StubInvoker<>(this, url, GreetService.class, handlers);
        }


        @Override
        public GreetResponse greet(GreetRequest request){
            throw unimplementedMethodException(greetMethod);
        }

        @Override
        public SayHelloResponse sayHello(SayHelloRequest request){
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
