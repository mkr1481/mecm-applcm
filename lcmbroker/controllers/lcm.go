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

// token controller
package controllers

import (
	log "github.com/sirupsen/logrus"
	"github.com/astaxie/beego"
)

type LcmController struct {
	beego.Controller
}

func (c *LcmController) UploadConfig() {
	log.Info("Add configuration request received.")
}

func (c *LcmController) RemoveConfig() {
	log.Info("Delete configuration request received.")
}

func (c *LcmController) Instantiate() {
	log.Info("Application instantiation request received.")
}

func (c *LcmController) Terminate() {
	log.Info("Application termination request received.")
}

func (c *LcmController) Query() {
	log.Info("Application query request received.")
}

func (c *LcmController) QueryKPI() {
	log.Info("Query KPI request received.")
}

func (c *LcmController) QueryMepCapabilities() {
	log.Info("Query mep capabilities request received.")
}
