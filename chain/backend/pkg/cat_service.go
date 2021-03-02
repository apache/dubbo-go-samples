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
	"github.com/apache/dubbo-go/config"
)

func init() {
	config.SetProviderService(new(CatService))
}

type CatService struct {
}

func (c *CatService) GetId() (int, error) {
	return 1, nil
}

func (c *CatService) GetName() (string, error) {
	fmt.Println("I am a Cat!")
	return "Cat", nil
}

func (c *CatService) Yell() (string, error) {
	fmt.Println("Meow Meow!")
	return "Meow Meow!", nil
}

func (c *CatService) Reference() string {
	return "CatService"
}
