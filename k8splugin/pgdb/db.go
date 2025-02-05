/*
 * Copyright 2020 Huawei Technologies Co., Ltd.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package pgdb

// Database API's
type Database interface {
	//KANAG: change to bool instead and set return var name as err
	InitDatabase(dbSslMode string) error
	InsertOrUpdateData(data interface{}, cols ...string) (err error)
	//KANAG: its better to callout the 2nd param used for data id in below methods
	ReadData(data interface{}, cols ...string) (err error)
	DeleteData(data interface{}, cols ...string) (err error)
}
