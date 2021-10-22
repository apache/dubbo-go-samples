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

package pkg

import (
	"math/rand"
)

import (
	"dubbo.apache.org/dubbo-go/v3/common/extension"
	"dubbo.apache.org/dubbo-go/v3/filter"

	"github.com/dubbogo/gost/log"
)

func init() {
	/*
	 * register your implementation and them using it like:
	 *
	 * UserProvider:
	 *   registry: hangzhouzk
	 *   protocol : dubbo
	 *   interface : com.ikurento.user.UserProvider
	 *   ... # other configuration
	 *   tps.limiter: method-service # the name of limiter
	 *   tps.limit.strategy: RandomLimitStrategy
	 */
	extension.SetTpsLimitStrategy("RandomLimitStrategy", &RandomTpsLimitStrategyCreator{})
}

/**
 * The RandomTpsLimitStrategy should not be singleton because different TpsLimiter will create many instances.
 * we won't want them affect each other.
 */
type RandomTpsLimitStrategy struct {
	rate     int
	interval int
}

/**
 * It implements the TpsLimitStrategy interface.
 * IsAllowable will return true if this invocation is not over limitation.
 */
func (r RandomTpsLimitStrategy) IsAllowable() bool {
	// this is a simple demo.
	gxlog.CInfo("Random IsAllowable!")
	randNum := rand.Int63n(2)
	return randNum == 0
}

type RandomTpsLimitStrategyCreator struct{}

/**
 * It implements the TpsLimitStrategyCreator interface.
 * TpsLimitStrategyCreator is the creator abstraction for TpsLimitStrategy.
 * Create will create an instance of TpsLimitStrategy.
 */
func (creator *RandomTpsLimitStrategyCreator) Create(rate int, interval int) filter.TpsLimitStrategy {
	return &RandomTpsLimitStrategy{
		rate:     rate,
		interval: interval,
	}
}
