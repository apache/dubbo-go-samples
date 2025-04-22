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

package bookingflight

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/apache/dubbo-go-samples/llm/book-flight/go-server/tools"
)

var (
	date string
)

/*
SearchFlightTicket
*/
type SearchFlightTicket struct {
	tools.BaseTool
}

type serachFlightTicketData struct {
	Origin             string `json:"origin" validate:"required"`
	Destination        string `json:"destination" validate:"required"`
	Date               string `json:"date" validate:"required"`
	DepartureTimeStart string `json:"departure_time_start"`
	DepartureTimeEnd   string `json:"departure_time_end"`
}

func NewSearchFlightTicket(name string, description string) SearchFlightTicket {
	return SearchFlightTicket{
		tools.NewBaseTool(
			name, description, tools.GetStructKeys(serachFlightTicketData{}), "", "", ""),
	}
}

// origin string, destination string, date string, departureTimeStart string, departureTimeEnd string
func (stt SearchFlightTicket) Call(ctx context.Context, input string) (string, error) {
	data := serachFlightTicketData{}
	err := json.Unmarshal([]byte(input), &data)
	if err != nil {
		return fmt.Sprintf("Error: %v", err), err
	}

	return stt.searchFlightTicket(data)
}

func (stt SearchFlightTicket) searchFlightTicket(data serachFlightTicketData) (string, error) {
	// 此处只做出发地校验，其他信息未进行校验
	if data.Origin != "北京" {
		return "未查询到相关内容", nil
	}

	date = data.Date
	rst := flightInformation()
	rst_json, err := json.Marshal(rst)
	return string(rst_json), err
}

/*
PurchaseFlightTicket
*/
type PurchaseFlightTicket struct {
	tools.BaseTool
}

type purchaseFlightTicketData struct {
	FlightNumber string `json:"flight_number" validate:"required"`
}

func NewPurchaseFlightTicket(name string, description string) PurchaseFlightTicket {
	return PurchaseFlightTicket{
		tools.NewBaseTool(
			name, description, tools.GetStructKeys(purchaseFlightTicketData{}), "", "", ""),
	}
}

func (ptt PurchaseFlightTicket) Call(ctx context.Context, input string) (string, error) {
	data := purchaseFlightTicketData{}
	err := json.Unmarshal([]byte(input), &data)
	if err != nil {
		return fmt.Sprintf("Error: %v", err), err
	}

	return ptt.purchaseFlightTicket(data)
}

func (ptt *PurchaseFlightTicket) purchaseFlightTicket(data purchaseFlightTicketData) (string, error) {
	flightInfo := flightInformation()
	for _, info := range flightInfo {
		if data.FlightNumber == info["flight_number"] {
			info["message"] = "购买成功"
			rst_json, err := json.Marshal(info)
			return string(rst_json), err
		}
	}

	return fmt.Sprintf("The flight was not found: %v", data.FlightNumber), nil
}

func flightInformation() []map[string]string {
	return []map[string]string{
		{
			"flight_number":  "MU5100",
			"origin":         "北京",
			"destination":    "上海",
			"departure_time": fmt.Sprintf("%v 07:00", date),
			"arrival_time":   fmt.Sprintf("%v 09:15", date),
			"price":          "900.00",
			"seat_type":      "头等舱",
		},
		{
			"flight_number":  "MU6865",
			"origin":         "北京",
			"destination":    "上海",
			"departure_time": fmt.Sprintf("%v 07:20", date),
			"arrival_time":   fmt.Sprintf("%v 09:25", date),
			"price":          "1160.00",
			"seat_type":      "头等舱",
		},
		{
			"flight_number":  "HM7601",
			"origin":         "北京",
			"destination":    "上海",
			"departure_time": fmt.Sprintf("%v 07:30", date),
			"arrival_time":   fmt.Sprintf("%v 09:55", date),
			"price":          "1080.00",
			"seat_type":      "普通舱",
		},
		{
			"flight_number":  "CA1515",
			"origin":         "北京",
			"destination":    "上海",
			"departure_time": fmt.Sprintf("%v 15:45", date),
			"arrival_time":   fmt.Sprintf("%v 17:55", date),
			"price":          "1080.00",
			"seat_type":      "普通舱",
		},
		{
			"flight_number":  "GS9012",
			"origin":         "北京",
			"destination":    "上海",
			"departure_time": fmt.Sprintf("%v 19:00", date),
			"arrival_time":   fmt.Sprintf("%v 23:00", date),
			"price":          "1250.00",
			"seat_type":      "头等舱",
		},
		{
			"flight_number":  "GS9013",
			"origin":         "北京",
			"destination":    "上海",
			"departure_time": fmt.Sprintf("%v 18:30", date),
			"arrival_time":   fmt.Sprintf("%v 22:00", date),
			"price":          "1200.00",
			"seat_type":      "头等舱",
		},
	}
}
