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

package dbAdapter

import "github.com/astaxie/beego/orm"

// Database API's
type Database interface {
	InitDatabase() error
	InsertOrUpdateData(data interface{}, cols ...string) (err error)
	ReadData(data interface{}, cols ...string) (err error)
	DeleteData(data interface{}, cols ...string) (err error)
	QueryCount(tableName string) (int64, error)
	QueryCountForAppInfo(tableName, fieldName, fieldValue string) (int64, error)
	QueryTable(query string) orm.QuerySeter
	LoadRelated(md interface{}, name string) (int64, error)
}
