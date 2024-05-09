package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	mgrpc "github.com/go-micro/plugins/v4/client/grpc"
	_ "github.com/go-micro/plugins/v4/registry/kubernetes"
	mhttp "github.com/go-micro/plugins/v4/server/http"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"

	"github.com/go-micro/demo/frontend/config"
	pb "github.com/go-micro/demo/frontend/proto"
)

const (
	name    = "frontend"
	version = "1.0.0"

	defaultCurrency = "USD"
	cookieMaxAge    = 60 * 60 * 48

	cookiePrefix    = "shop_"
	cookieSessionID = cookiePrefix + "session-id"
	cookieCurrency  = cookiePrefix + "currency"
)

var (
	whitelistedCurrencies = map[string]bool{
		"USD": true,
		"EUR": true,
		"CAD": true,
		"JPY": true,
		"GBP": true,
		"TRY": true,
	}
)

type ctxKeySessionID struct{}

type frontendServer struct {
	adService             pb.AdService
	cartService           pb.CartService
	checkoutService       pb.CheckoutService
	currencyService       pb.CurrencyService
	productCatalogService pb.ProductCatalogService
	recommendationService pb.RecommendationService
	shippingService       pb.ShippingService
}

func main() {
	// Load conigurations
	if err := config.Load(); err != nil {
		logger.Fatal(err)
	}

	// Create service
	srv := micro.NewService(
		micro.Server(mhttp.NewServer()),
		micro.Client(mgrpc.NewClient()),
	)
	opts := []micro.Option{
		micro.Name(name),
		micro.Version(version),
		micro.Address(config.Address()),
	}
	if cfg := config.Tracing(); cfg.Enable {
		tp, err := newTracerProvider(name, srv.Server().Options().Id, cfg.Jaeger.URL)
		if err != nil {
			logger.Fatal(err)
		}
		defer func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()
			if err := tp.Shutdown(ctx); err != nil {
				logger.Fatal(err)
			}
		}()
		otel.SetTracerProvider(tp)
		otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	}
	srv.Init(opts...)

	log := logrus.New()
	log.Level = logrus.DebugLevel
	log.Formatter = &logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyMsg:   "message",
		},
		TimestampFormat: time.RFC3339Nano,
	}
	log.Out = os.Stdout

	cfg, client := config.Get(), srv.Client()
	svc := &frontendServer{
		adService:             pb.NewAdService(cfg.AdService, client),
		cartService:           pb.NewCartService(cfg.CartService, client),
		checkoutService:       pb.NewCheckoutService(cfg.CheckoutService, client),
		currencyService:       pb.NewCurrencyService(cfg.CurrencyService, client),
		productCatalogService: pb.NewProductCatalogService(cfg.ProductCatalogService, client),
		recommendationService: pb.NewRecommendationService(cfg.RecommendationService, client),
		shippingService:       pb.NewShippingService(cfg.ShippingService, client),
	}

	r := mux.NewRouter()
	r.HandleFunc("/", svc.homeHandler).Methods(http.MethodGet, http.MethodHead)
	r.HandleFunc("/product/{id}", svc.productHandler).Methods(http.MethodGet, http.MethodHead)
	r.HandleFunc("/cart", svc.viewCartHandler).Methods(http.MethodGet, http.MethodHead)
	r.HandleFunc("/cart", svc.addToCartHandler).Methods(http.MethodPost)
	r.HandleFunc("/cart/empty", svc.emptyCartHandler).Methods(http.MethodPost)
	r.HandleFunc("/setCurrency", svc.setCurrencyHandler).Methods(http.MethodPost)
	r.HandleFunc("/logout", svc.logoutHandler).Methods(http.MethodGet)
	r.HandleFunc("/cart/checkout", svc.placeOrderHandler).Methods(http.MethodPost)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	r.HandleFunc("/robots.txt", func(w http.ResponseWriter, _ *http.Request) { fmt.Fprint(w, "User-agent: *\nDisallow: /") })
	r.HandleFunc("/_healthz", func(w http.ResponseWriter, _ *http.Request) { fmt.Fprint(w, "ok") })

	var handler http.Handler = r
	handler = &logHandler{log: log, next: handler} // add logging
	handler = ensureSessionID(handler)             // add session ID
	// handler = tracing(handler)                     // add opentelemetry instrumentation
	r.Use(otelmux.Middleware(name))
	r.Use(tracingContextWrapper)
	if err := micro.RegisterHandler(srv.Server(), handler); err != nil {
		logger.Fatal(err)
	}

	logger.Infof("starting server on %s", config.Address())
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
