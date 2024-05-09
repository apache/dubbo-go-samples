package handler

import (
	"context"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"go-micro.dev/v4/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"

	pb "github.com/go-micro/demo/productcatalogservice/proto"
)

var reloadCatalog bool

type ProductCatalogService struct {
	sync.Mutex
	products []*pb.Product
}

func (s *ProductCatalogService) ListProducts(ctx context.Context, in *pb.Empty, out *pb.ListProductsResponse) error {
	out.Products = s.parseCatalog()
	return nil
}

func (s *ProductCatalogService) GetProduct(ctx context.Context, in *pb.GetProductRequest, out *pb.Product) error {
	var found *pb.Product
	products := s.parseCatalog()
	for _, p := range products {
		if in.Id == p.Id {
			found = p
		}
	}
	if found == nil {
		return status.Errorf(codes.NotFound, "no product with ID %s", in.Id)
	}
	out.Id = found.Id
	out.Name = found.Name
	out.Categories = found.Categories
	out.Description = found.Description
	out.Picture = found.Picture
	out.PriceUsd = found.PriceUsd
	return nil
}

func (s *ProductCatalogService) SearchProducts(ctx context.Context, in *pb.SearchProductsRequest, out *pb.SearchProductsResponse) error {
	var ps []*pb.Product
	products := s.parseCatalog()
	for _, p := range products {
		if strings.Contains(strings.ToLower(p.Name), strings.ToLower(in.Query)) ||
			strings.Contains(strings.ToLower(p.Description), strings.ToLower(in.Query)) {
			ps = append(ps, p)
		}
	}
	out.Results = ps
	return nil
}

func (s *ProductCatalogService) readCatalogFile() (*pb.ListProductsResponse, error) {
	s.Lock()
	defer s.Unlock()
	catalogJSON, err := ioutil.ReadFile("data/products.json")
	if err != nil {
		logger.Errorf("failed to open product catalog json file: %v", err)
		return nil, err
	}
	catalog := &pb.ListProductsResponse{}
	if err := protojson.Unmarshal(catalogJSON, catalog); err != nil {
		logger.Warnf("failed to parse the catalog JSON: %v", err)
		return nil, err
	}
	logger.Info("successfully parsed product catalog json")
	return catalog, nil
}

func (s *ProductCatalogService) parseCatalog() []*pb.Product {
	if reloadCatalog || len(s.products) == 0 {
		catalog, err := s.readCatalogFile()
		if err != nil {
			return []*pb.Product{}
		}
		s.products = catalog.Products
	}
	return s.products
}

func init() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGUSR1, syscall.SIGUSR2)
	go func() {
		for {
			sig := <-sigs
			logger.Infof("Received signal: %s", sig)
			if sig == syscall.SIGUSR1 {
				reloadCatalog = true
				logger.Infof("Enable catalog reloading")
			} else {
				reloadCatalog = false
				logger.Infof("Disable catalog reloading")
			}
		}
	}()
}
