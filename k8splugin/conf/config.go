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

package conf

//KANAG: Move the ServerCOnfigurations under Configurations as there is only one field in this type
// Configurations exported
type Configurations struct {
	Server ServerConfigurations
}

// ServerConfigurations exported
type ServerConfigurations struct {
//KANAG: Strictly follow Naming notations 	for all fields
//KANAG: Should this be array to hold multiple ciphers
	Sslciphers    string
	Servername    string
	Sslnotenabled bool
	Certfilepath  string
	Keyfilepath   string
	Serverport    string
	Httpsaddr     string
	DbAdapter     string
//KANAG: is this to be bool 
	DbSslMode     string
}
