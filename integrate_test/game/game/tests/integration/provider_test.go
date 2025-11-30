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

package integration

import (
	"context"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

import (
	game "github.com/apache/dubbo-go-samples/game/proto/game"
)

func TestLogin(t *testing.T) {
	uid := t.Name()
	res := login(t, uid)
	data := res.Data.AsMap()
	assert.Equal(t, uid, data["to"])
	assert.EqualValues(t, 0, data["score"])
}

func TestScore(t *testing.T) {
	uid := t.Name()
	_ = login(t, uid)

	res, err := gameService.Score(context.TODO(), &game.ScoreRequest{Uid: uid, Score: "1"})
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, int32(0), res.Code)
	data := res.Data.AsMap()
	assert.Equal(t, uid, data["to"])
	assert.EqualValues(t, 1, data["score"])
}

func TestRank(t *testing.T) {
	uid := t.Name()
	_ = login(t, uid)

	// boost score so this user ranks first regardless of previous tests
	_, err := gameService.Score(context.TODO(), &game.ScoreRequest{Uid: uid, Score: "100"})
	assert.Nil(t, err)

	res, err := gameService.Rank(context.TODO(), &game.RankRequest{Uid: uid})
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, int32(0), res.Code)
	data := res.Data.AsMap()
	assert.Equal(t, uid, data["to"])
	assert.EqualValues(t, 1, data["rank"])
}

func login(t *testing.T, uid string) *game.Result {
	t.Helper()
	res, err := gameService.Login(context.TODO(), &game.LoginRequest{Uid: uid})
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, int32(0), res.Code)
	return res
}
