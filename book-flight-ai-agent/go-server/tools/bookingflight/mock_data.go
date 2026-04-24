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
	"fmt"
)

func flightInformation() []map[string]string {
	return []map[string]string{
		// 北京 → 上海
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
			"departure_time": fmt.Sprintf("%v 10:30", date),
			"arrival_time":   fmt.Sprintf("%v 12:55", date),
			"price":          "1080.00",
			"seat_type":      "普通舱",
		},
		{
			"flight_number":  "CA1515",
			"origin":         "北京",
			"destination":    "上海",
			"departure_time": fmt.Sprintf("%v 15:45", date),
			"arrival_time":   fmt.Sprintf("%v 17:55", date),
			"price":          "980.00",
			"seat_type":      "普通舱",
		},
		{
			"flight_number":  "GS9012",
			"origin":         "北京",
			"destination":    "上海",
			"departure_time": fmt.Sprintf("%v 19:00", date),
			"arrival_time":   fmt.Sprintf("%v 21:15", date),
			"price":          "1250.00",
			"seat_type":      "头等舱",
		},
		{
			"flight_number":  "GS9013",
			"origin":         "北京",
			"destination":    "上海",
			"departure_time": fmt.Sprintf("%v 21:30", date),
			"arrival_time":   fmt.Sprintf("%v 23:45", date),
			"price":          "850.00",
			"seat_type":      "普通舱",
		},
		// 上海 → 北京
		{
			"flight_number":  "MU5101",
			"origin":         "上海",
			"destination":    "北京",
			"departure_time": fmt.Sprintf("%v 08:00", date),
			"arrival_time":   fmt.Sprintf("%v 10:20", date),
			"price":          "850.00",
			"seat_type":      "普通舱",
		},
		{
			"flight_number":  "CA1800",
			"origin":         "上海",
			"destination":    "北京",
			"departure_time": fmt.Sprintf("%v 12:30", date),
			"arrival_time":   fmt.Sprintf("%v 14:45", date),
			"price":          "920.00",
			"seat_type":      "普通舱",
		},
		{
			"flight_number":  "MU7788",
			"origin":         "上海",
			"destination":    "北京",
			"departure_time": fmt.Sprintf("%v 17:00", date),
			"arrival_time":   fmt.Sprintf("%v 19:20", date),
			"price":          "1100.00",
			"seat_type":      "头等舱",
		},
		{
			"flight_number":  "CZ6601",
			"origin":         "上海",
			"destination":    "北京",
			"departure_time": fmt.Sprintf("%v 20:15", date),
			"arrival_time":   fmt.Sprintf("%v 22:30", date),
			"price":          "980.00",
			"seat_type":      "普通舱",
		},
		// 北京 → 广州
		{
			"flight_number":  "CZ3001",
			"origin":         "北京",
			"destination":    "广州",
			"departure_time": fmt.Sprintf("%v 08:10", date),
			"arrival_time":   fmt.Sprintf("%v 11:20", date),
			"price":          "1350.00",
			"seat_type":      "普通舱",
		},
		{
			"flight_number":  "CZ3002",
			"origin":         "北京",
			"destination":    "广州",
			"departure_time": fmt.Sprintf("%v 13:00", date),
			"arrival_time":   fmt.Sprintf("%v 16:10", date),
			"price":          "1420.00",
			"seat_type":      "头等舱",
		},
		{
			"flight_number":  "CA1301",
			"origin":         "北京",
			"destination":    "广州",
			"departure_time": fmt.Sprintf("%v 18:45", date),
			"arrival_time":   fmt.Sprintf("%v 21:55", date),
			"price":          "1300.00",
			"seat_type":      "普通舱",
		},
		// 广州 → 北京
		{
			"flight_number":  "CZ3101",
			"origin":         "广州",
			"destination":    "北京",
			"departure_time": fmt.Sprintf("%v 07:50", date),
			"arrival_time":   fmt.Sprintf("%v 11:05", date),
			"price":          "1280.00",
			"seat_type":      "普通舱",
		},
		{
			"flight_number":  "CZ3102",
			"origin":         "广州",
			"destination":    "北京",
			"departure_time": fmt.Sprintf("%v 14:20", date),
			"arrival_time":   fmt.Sprintf("%v 17:35", date),
			"price":          "1380.00",
			"seat_type":      "头等舱",
		},
		{
			"flight_number":  "CA1302",
			"origin":         "广州",
			"destination":    "北京",
			"departure_time": fmt.Sprintf("%v 19:30", date),
			"arrival_time":   fmt.Sprintf("%v 22:40", date),
			"price":          "1250.00",
			"seat_type":      "普通舱",
		},
		// 上海 → 广州
		{
			"flight_number":  "MU2301",
			"origin":         "上海",
			"destination":    "广州",
			"departure_time": fmt.Sprintf("%v 09:00", date),
			"arrival_time":   fmt.Sprintf("%v 11:30", date),
			"price":          "780.00",
			"seat_type":      "普通舱",
		},
		{
			"flight_number":  "CZ8801",
			"origin":         "上海",
			"destination":    "广州",
			"departure_time": fmt.Sprintf("%v 16:00", date),
			"arrival_time":   fmt.Sprintf("%v 18:25", date),
			"price":          "860.00",
			"seat_type":      "普通舱",
		},
		{
			"flight_number":  "MU2302",
			"origin":         "上海",
			"destination":    "广州",
			"departure_time": fmt.Sprintf("%v 20:30", date),
			"arrival_time":   fmt.Sprintf("%v 22:55", date),
			"price":          "920.00",
			"seat_type":      "头等舱",
		},
		// 广州 → 上海
		{
			"flight_number":  "CZ8802",
			"origin":         "广州",
			"destination":    "上海",
			"departure_time": fmt.Sprintf("%v 08:30", date),
			"arrival_time":   fmt.Sprintf("%v 11:00", date),
			"price":          "800.00",
			"seat_type":      "普通舱",
		},
		{
			"flight_number":  "MU2401",
			"origin":         "广州",
			"destination":    "上海",
			"departure_time": fmt.Sprintf("%v 13:45", date),
			"arrival_time":   fmt.Sprintf("%v 16:10", date),
			"price":          "870.00",
			"seat_type":      "普通舱",
		},
		{
			"flight_number":  "CZ8803",
			"origin":         "广州",
			"destination":    "上海",
			"departure_time": fmt.Sprintf("%v 19:00", date),
			"arrival_time":   fmt.Sprintf("%v 21:20", date),
			"price":          "950.00",
			"seat_type":      "头等舱",
		},
		// 北京 → 成都
		{
			"flight_number":  "CA4101",
			"origin":         "北京",
			"destination":    "成都",
			"departure_time": fmt.Sprintf("%v 09:20", date),
			"arrival_time":   fmt.Sprintf("%v 12:05", date),
			"price":          "1180.00",
			"seat_type":      "普通舱",
		},
		{
			"flight_number":  "SC4201",
			"origin":         "北京",
			"destination":    "成都",
			"departure_time": fmt.Sprintf("%v 14:00", date),
			"arrival_time":   fmt.Sprintf("%v 16:50", date),
			"price":          "1260.00",
			"seat_type":      "头等舱",
		},
		{
			"flight_number":  "CA4102",
			"origin":         "北京",
			"destination":    "成都",
			"departure_time": fmt.Sprintf("%v 19:30", date),
			"arrival_time":   fmt.Sprintf("%v 22:15", date),
			"price":          "1150.00",
			"seat_type":      "普通舱",
		},
		// 成都 → 北京
		{
			"flight_number":  "CA4201",
			"origin":         "成都",
			"destination":    "北京",
			"departure_time": fmt.Sprintf("%v 08:00", date),
			"arrival_time":   fmt.Sprintf("%v 10:45", date),
			"price":          "1150.00",
			"seat_type":      "普通舱",
		},
		{
			"flight_number":  "SC4301",
			"origin":         "成都",
			"destination":    "北京",
			"departure_time": fmt.Sprintf("%v 13:30", date),
			"arrival_time":   fmt.Sprintf("%v 16:20", date),
			"price":          "1230.00",
			"seat_type":      "头等舱",
		},
		{
			"flight_number":  "CA4202",
			"origin":         "成都",
			"destination":    "北京",
			"departure_time": fmt.Sprintf("%v 18:50", date),
			"arrival_time":   fmt.Sprintf("%v 21:35", date),
			"price":          "1100.00",
			"seat_type":      "普通舱",
		},
	}
}
