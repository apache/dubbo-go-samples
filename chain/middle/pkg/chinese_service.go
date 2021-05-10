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
	"fmt"
)

import (
	"dubbo.apache.org/dubbo-go/v3/config"
)

type DogService struct {
	GetId   func() (int, error)
	GetName func() (string, error)
	Yell    func() (string, error)
}

func (d *DogService) Reference() string {
	return "DogService"
}

type TigerService struct {
	GetId   func() (int, error)
	GetName func() (string, error)
	Yell    func() (string, error)
}

func (t *TigerService) Reference() string {
	return "TigerService"
}

func init() {
	dog := new(DogService)
	config.SetConsumerService(dog)
	tiger := new(TigerService)
	config.SetConsumerService(tiger)

	config.SetProviderService(&ChineseService{
		dog:   dog,
		tiger: tiger,
	})
}

type ChineseService struct {
	dog   *DogService
	tiger *TigerService
}

func (c *ChineseService) Have() (string, error) {
	name, _ := c.dog.GetName()
	return "I'm Chinese and I have a " + name, nil
}

func (c *ChineseService) Hear() (string, error) {
	name, _ := c.tiger.GetName()
	yell, _ := c.tiger.Yell()
	return fmt.Sprintf("I'm Chinese and I heard a %s yells like %s", name, yell), nil
}

func (c *ChineseService) Reference() string {
	return "ChineseService"
}
