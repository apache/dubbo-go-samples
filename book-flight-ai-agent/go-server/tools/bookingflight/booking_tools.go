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
)

import (
	"github.com/apache/dubbo-go-samples/book-flight-ai-agent/go-server/tools"
)

var (
	date string
)

/*
SearchFlightTicketTool
*/
type SearchFlightTicketTool struct {
	tools.BaseTool
	Origin             string `json:"origin" validate:"required"`
	Destination        string `json:"destination" validate:"required"`
	Date               string `json:"date" validate:"required"`
	DepartureTimeStart string `json:"departure_time_start"`
	DepartureTimeEnd   string `json:"departure_time_end"`
}

// origin string, destination string, date string, departureTimeStart string, departureTimeEnd string
func (stt *SearchFlightTicketTool) Call(ctx context.Context, input string) (string, error) {
	err := json.Unmarshal([]byte(input), stt)
	if err != nil {
		return fmt.Sprintf("Error: %v", err), err
	}

	return stt.searchFlightTicket()
}

func (stt *SearchFlightTicketTool) searchFlightTicket() (string, error) {
	// Only the departure point is verified here, and other information is not verified
	if stt.Origin != "北京" {
		return "No relevant content was found", nil
	}

	date = stt.Date
	rst := flightInformation()
	rst_json, err := json.Marshal(rst)
	return string(rst_json), err
}

/*
PurchaseFlightTicketTool
*/
type PurchaseFlightTicketTool struct {
	tools.BaseTool
	FlightNumber string `json:"flight_number" validate:"required"`
}

func (ptt *PurchaseFlightTicketTool) Call(ctx context.Context, input string) (string, error) {
	err := json.Unmarshal([]byte(input), &ptt)
	if err != nil {
		return fmt.Sprintf("Error: %v", err), err
	}

	return ptt.purchaseFlightTicket()
}

func (ptt *PurchaseFlightTicketTool) purchaseFlightTicket() (string, error) {
	flightInfo := flightInformation()
	for _, info := range flightInfo {
		if ptt.FlightNumber == info["flight_number"] {
			info["message"] = "Successful purchase."
			rst_json, err := json.Marshal(info)
			return string(rst_json), err
		}
	}

	return fmt.Sprintf("The flight was not found: %v", ptt.FlightNumber), nil
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
