package main

import (
	"dubbo.apache.org/dubbo-go/v3"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/registry"
	"github.com/apache/dubbo-go-samples/online_boutique_demo/checkoutservice/handler"
	pb "github.com/apache/dubbo-go-samples/online_boutique_demo/checkoutservice/proto"
	"github.com/dubbogo/gost/log/logger"
)

var (
	name    = "checkoutservice"
	version = "1.0.0"
)

func main() {

	ins, err := dubbo.NewInstance(
		dubbo.WithName(name),
		dubbo.WithRegistry(
			registry.WithZookeeper(),
			registry.WithAddress("127.0.0.1:2181"),
		),
		dubbo.WithProtocol(
			protocol.WithTriple(),
			protocol.WithPort(20002),
		),
	)
	if err != nil {
		panic(err)
	}
	//server
	srv, err := ins.NewServer()
	if err != nil {
		panic(err)
	}
	client, err := ins.NewClient()
	cartService, err := pb.NewCartService(client)
	currencyService, err := pb.NewCurrencyService(client)
	emailService, err := pb.NewEmailService(client)
	paymentService, err := pb.NewPaymentService(client)
	productCatalogService, err := pb.NewProductCatalogService(client)
	shippingService, err := pb.NewShippingService(client)
	if err := pb.RegisterCheckoutServiceHandler(srv, &handler.CheckoutService{
		CartService:           cartService,
		CurrencyService:       currencyService,
		EmailService:          emailService,
		PaymentService:        paymentService,
		ProductCatalogService: productCatalogService,
		ShippingService:       shippingService,
	}); err != nil {
		panic(err)
	}

	if err := srv.Serve(); err != nil {
		logger.Error(err)
	}
}
