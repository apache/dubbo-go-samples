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

package handler

import (
	"context"
	pb "github.com/apache/dubbo-go-samples/online_boutique_demo/productcatalogservice/proto"
	"github.com/dubbogo/gost/log/logger"
	"github.com/dubbogo/grpc-go/codes"
	"github.com/dubbogo/grpc-go/status"
	"google.golang.org/protobuf/encoding/protojson"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

var reloadCatalog bool

type ProductCatalogService struct {
	sync.Mutex
	products []*pb.Product
}

func (s *ProductCatalogService) ListProducts(ctx context.Context, in *pb.Empty) (*pb.ListProductsResponse, error) {
	out := &pb.ListProductsResponse{}
	out.Products = s.parseCatalog()
	return out, nil
}

func (s *ProductCatalogService) GetProduct(ctx context.Context, in *pb.GetProductRequest) (*pb.Product, error) {
	var found *pb.Product
	out := &pb.Product{}
	products := s.parseCatalog()
	for _, p := range products {
		if in.Id == p.Id {
			found = p
		}
	}
	if found == nil {
		return nil, status.Errorf(codes.NotFound, "Product not found with ID %s", in.Id)
	}
	out.Id = found.Id
	out.Name = found.Name
	out.Categories = found.Categories
	out.Description = found.Description
	out.Picture = found.Picture
	out.PriceUsd = found.PriceUsd
	return out, nil
}

func (s *ProductCatalogService) SearchProducts(ctx context.Context, in *pb.SearchProductsRequest) (*pb.SearchProductsResponse, error) {
	var ps []*pb.Product
	out := &pb.SearchProductsResponse{}
	products := s.parseCatalog()
	for _, p := range products {
		if strings.Contains(strings.ToLower(p.Name), strings.ToLower(in.Query)) ||
			strings.Contains(strings.ToLower(p.Description), strings.ToLower(in.Query)) {
			ps = append(ps, p)
		}
	}
	out.Results = ps
	return out, nil
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
	logger.Infof("successfully parsed product catalog json")
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
