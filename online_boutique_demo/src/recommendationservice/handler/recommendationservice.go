package handler

import (
	"context"
	"github.com/dubbogo/gost/log/logger"
	"math/rand"

	pb "github.com/apache/dubbo-go-samples/online_boutique_demo/recommendationservice/proto"
)

type RecommendationService struct {
	ProductCatalogService pb.ProductCatalogService
}

func (s *RecommendationService) ListRecommendations(ctx context.Context, in *pb.ListRecommendationsRequest) (*pb.ListRecommendationsResponse, error) {
	maxResponsesCount := 5
	var out *pb.ListRecommendationsResponse
	// # fetch list of products from product catalog stub
	catalog, err := s.ProductCatalogService.ListProducts(ctx, &pb.Empty{})
	if err != nil {
		return nil, nil
	}
	filteredProductsIDs := make([]string, 0, len(catalog.Products))
	for _, p := range catalog.Products {
		if contains(p.Id, in.ProductIds) {
			continue
		}
		filteredProductsIDs = append(filteredProductsIDs, p.Id)
	}
	productIDs := sample(filteredProductsIDs, maxResponsesCount)
	logger.Infof("[Recv ListRecommendations] product_ids=%v", productIDs)
	out.ProductIds = productIDs
	return out, nil
}

func contains(target string, source []string) bool {
	for _, s := range source {
		if target == s {
			return true
		}
	}
	return false
}

func sample(source []string, c int) []string {
	n := len(source)
	if n <= c {
		return source
	}
	indices := make([]int, n)
	for i := 0; i < n; i++ {
		indices[i] = i
	}
	for i := n - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		indices[i], indices[j] = indices[j], indices[i]
	}
	result := make([]string, 0, c)
	for i := 0; i < c; i++ {
		result = append(result, source[indices[i]])
	}
	return result
}
