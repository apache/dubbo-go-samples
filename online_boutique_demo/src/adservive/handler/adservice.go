package handler

import (
	"context"
	"math/rand"

	pb "github.com/apache/dubbo-go-demo/adservice/proto"
	"github.com/dubbogo/gost/log/logger"
)

const MAX_ADS_TO_SERVE = 2

var adsMap = createAdsMap()

type AdTripleService struct{}

func (s *AdTripleService) GetAds(ctx context.Context, req *pb.AdRequest) (res *pb.AdResponse, err error) {
	logger.Infof("received ad request (context_words= %v )", req.ContextKeys)
	allAds := make([]*pb.Ad, 0)
	res = new(pb.AdResponse)
	if len(req.ContextKeys) > 0 {
		for _, category := range req.ContextKeys {
			ads := getAdsByCategory(category)
			allAds = append(allAds, ads...)
		}
		if len(allAds) == 0 {
			allAds = getRandomAds()
		}
	} else {
		allAds = getRandomAds()
	}
	res.Ads = allAds
	return res, nil
}

func getAdsByCategory(category string) []*pb.Ad {
	return adsMap[category]
}

func getRandomAds() []*pb.Ad {
	ads := make([]*pb.Ad, 0, MAX_ADS_TO_SERVE)
	allAds := make([]*pb.Ad, 0, 7)
	for _, ads := range adsMap {
		allAds = append(allAds, ads...)
	}
	for i := 0; i < MAX_ADS_TO_SERVE; i++ {
		ads = append(ads, allAds[rand.Intn(len(allAds))])
	}
	return ads
}

func createAdsMap() map[string][]*pb.Ad {
	hairdryer := &pb.Ad{RedirectUrl: "/product/2ZYFJ3GM2N", Text: "Hairdryer for sale. 50% off."}
	tankTop := &pb.Ad{RedirectUrl: "/product/66VCHSJNUP", Text: "Tank top for sale. 20% off."}
	candleHolder := &pb.Ad{RedirectUrl: "/product/0PUK6V6EV0", Text: "Candle holder for sale. 30% off."}
	bambooGlassJar := &pb.Ad{RedirectUrl: "/product/9SIQT8TOJO", Text: "Bamboo glass jar for sale. 10% off."}
	watch := &pb.Ad{RedirectUrl: "/product/1YMWWN1N4O", Text: "Watch for sale. Buy one, get second kit for free"}
	mug := &pb.Ad{RedirectUrl: "/product/6E92ZMYYFZ", Text: "Mug for sale. Buy two, get third one for free"}
	loafers := &pb.Ad{RedirectUrl: "/product/L9ECAV7KIM", Text: "Loafers for sale. Buy one, get second one for free"}
	return map[string][]*pb.Ad{
		"clothing":    {tankTop},
		"accessories": {watch},
		"footwear":    {loafers},
		"hair":        {hairdryer},
		"decor":       {candleHolder},
		"kitchen":     {bambooGlassJar, mug},
	}
}
