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
	date = stt.Date
	all := flightInformation()

	var rst []map[string]string
	for _, info := range all {
		if stt.Origin != "" && info["origin"] != stt.Origin {
			continue
		}
		if stt.Destination != "" && info["destination"] != stt.Destination {
			continue
		}
		// departure_time format: "date HH:MM", extract last 5 chars for time range filter
		dep_time := ""
		if t := info["departure_time"]; len(t) >= 5 {
			dep_time = t[len(t)-5:]
		}
		if stt.DepartureTimeStart != "" && dep_time < stt.DepartureTimeStart {
			continue
		}
		if stt.DepartureTimeEnd != "" && dep_time > stt.DepartureTimeEnd {
			continue
		}
		rst = append(rst, info)
	}

	if len(rst) == 0 {
		return "No relevant content was found", nil
	}

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
