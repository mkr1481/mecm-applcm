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
	"unsafe"

	log "github.com/sirupsen/logrus"
	"lcmbroker/pkg/handlers/pluginAdapter"
	"lcmbroker/util"
	"os"
)

type LcmController struct {
	beego.Controller
}

// Upload Config
func (c *LcmController) UploadConfig() {
	log.Info("Add configuration request received.")
	clientIp := c.Ctx.Request.Header.Get("X-Real-Ip")
	accessToken := c.Ctx.Request.Header.Get("access_token")
	err := util.ValidateAccessToken(accessToken)
	if err != nil {
		c.handleLoggingForError(clientIp, util.StatusUnauthorized, util.AuthorizationFailed)
		return
	}

	bKey := *(*[]byte)(unsafe.Pointer(&accessToken))
	hostIp, err := c.getHostIP(clientIp)
	if err != nil {
		util.ClearByteArray(bKey)
		return
	}

	file, header, err := c.GetFile("configFile")
	if err != nil {
		util.ClearByteArray(bKey)
		c.handleLoggingForError(clientIp, util.StatusInternalServerError, "Upload config file error")
		return
	}

	err = util.ValidateFileSize(header.Size, util.MaxConfigFile)
	if err != nil {
		util.ClearByteArray(bKey)
		c.handleLoggingForError(clientIp, util.StatusInternalServerError, "File size is larger than max size")
		return
	}

	pluginInfo := "helmplugin" + ":" + os.Getenv("HELM_PLUGIN_PORT")

	adapter := pluginAdapter.NewPluginAdapter(pluginInfo)
	_, err = adapter.UploadConfig(pluginInfo, file, hostIp, accessToken)
	util.ClearByteArray(bKey)
	if err != nil {
		c.handleLoggingForError(clientIp, util.StatusInternalServerError, "Upload configuration failed")
		return
	}

	c.ServeJSON()
}

// Remove Config
func (c *LcmController) RemoveConfig() {
	log.Info("Delete configuration request received.")
	clientIp := c.Ctx.Request.Header.Get("X-Real-Ip")
	accessToken := c.Ctx.Request.Header.Get("access_token")
	err := util.ValidateAccessToken(accessToken)
	if err != nil {
		c.handleLoggingForError(clientIp, util.StatusUnauthorized, util.AuthorizationFailed)
		return
	}
	bKey := *(*[]byte)(unsafe.Pointer(&accessToken))
	hostIp, err := c.getHostIP(clientIp)
	if err != nil {
		util.ClearByteArray(bKey)
		return
	}

	pluginInfo := "helmplugin" + ":" + os.Getenv("HELM_PLUGIN_PORT")

	adapter := pluginAdapter.NewPluginAdapter(pluginInfo)
	_, err = adapter.RemoveConfig(pluginInfo, hostIp, accessToken)
	util.ClearByteArray(bKey)
	if err != nil {
		c.handleLoggingForError(clientIp, util.StatusInternalServerError, "Remove configuration failed")
		return
	}
	c.ServeJSON()
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

func (c *LcmController) writeErrorResponse(errMsg string, code int) {
	log.Error(errMsg)
	c.writeResponse(errMsg, code)
}

func (c *LcmController) writeResponse(msg string, code int) {
	c.Data["json"] = msg
	c.Ctx.ResponseWriter.WriteHeader(code)
	c.ServeJSON()
}
// Get host IP
func (c *LcmController) getHostIP(clientIp string) (string, error) {
	hostIp := c.GetString("hostIp")
	err := util.ValidateIpv4Address(hostIp)
	if err != nil {
		c.handleLoggingForError(clientIp, util.BadRequest,"HostIp address is invalid")
		return "", err
	}
	return hostIp, nil
}
// Handled logging for error case
func (c *LcmController) handleLoggingForError(clientIp string, code int, errMsg string) {
	log.Info("Received message from ClientIP [" + clientIp + "] Operation [" + c.Ctx.Request.Method + "]" +
		" Resource [" + c.Ctx.Input.URL() + "]")
	c.writeErrorResponse(errMsg, code)
	log.Info("Response message for ClientIP [" + clientIp + "] Operation [" + c.Ctx.Request.Method + "]" +
		" Resource [" + c.Ctx.Input.URL() + "] Result [Failure: " + errMsg + ".]")
	return
}
