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

// Image controller
package controllers

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"lcmcontroller/models"
	"lcmcontroller/util"
	"unsafe"
)

// Image Controller
type ImageController struct {
	BaseController
}

// @Title Create Image
// @Description creation of image
// @Param   tenantId        path 	string	true   "tenantId"
// @Param   appInstanceId   path 	string	true   "appInstanceId"
// @Param   access_token    header  string  true   "access token"
// @Param   vmId            body 	string	true   "vmId"
// @Success 200 ok
// @Failure 400 bad request
// @Failure 404 vm doesn't exist
// @router /tenants/:tenantId/app_instances/:appInstanceId/images [post]
func (c *ImageController) CreateImage() {
	log.Info("Image creation request received.")
	clientIp := c.Ctx.Input.IP()
	err := util.ValidateSrcAddress(clientIp)
	if err != nil {
		c.handleLoggingForError(clientIp, util.BadRequest, util.ClientIpaddressInvalid)
		return
	}
	c.displayReceivedMsg(clientIp)
	accessToken := c.Ctx.Request.Header.Get(util.AccessToken)
	bKey := *(*[]byte)(unsafe.Pointer(&accessToken))
	_, err = c.isPermitted(accessToken, clientIp)
	if err != nil {
		util.ClearByteArray(bKey)
		return
	}

	appInsId, err := c.getAppInstId(clientIp)
	if err != nil {
		util.ClearByteArray(bKey)
		return
	}

	appInfoRecord, err := c.getAppInfoRecord(appInsId, clientIp)
	if err != nil {
		util.ClearByteArray(bKey)
		return
	}

	vim, err := c.getVim(clientIp, appInfoRecord.HostIp)
	if err != nil {
		util.ClearByteArray(bKey)
		return
	}

	adapter, err := c.getPluginAdapter(appInfoRecord.DeployType, clientIp, vim)
	if err != nil {
		util.ClearByteArray(bKey)
		return
	}

	var request models.CreateVimRequest
	err = json.Unmarshal(c.Ctx.Input.RequestBody, &request)
	if err != nil {
		c.writeErrorResponse("failed to unmarshal request", util.BadRequest)
		util.ClearByteArray(bKey)
		return
	}

	response, err := adapter.CreateVmImage(appInfoRecord.HostIp, accessToken, appInfoRecord.AppInsId, request.VmId)
	util.ClearByteArray(bKey)
	if err != nil {
		// To check if any more error code needs to be returned.
		c.handleLoggingForError(clientIp, util.BadRequest, err.Error())
		return
	}

	_, err = c.Ctx.ResponseWriter.Write([]byte(response))
	if err != nil {
		c.handleLoggingForError(clientIp, util.StatusInternalServerError, util.FailedToWriteRes)
		return
	}

	c.handleLoggingForSuccess(clientIp, "VM Image creation is successful")
}


// @Title Delete Image
// @Description deletion of image
// @Param   tenantId        path 	string	true   "tenantId"
// @Param   appInstanceId   path 	string	true   "appInstanceId"
// @Param   imageId         path 	string	true   "imageId"
// @Param   access_token    header  string  true   "access token"
// @Success 204 no content
// @Failure 404 image doesn't exist
// @router /tenants/:tenantId/app_instances/:appInstanceId/images/:imageId [delete]
func (c *ImageController) DeleteImage() {
	log.Info("Image deletion request received.")
}

// @Title Query Image
// @Description query of image
// @Param   tenantId        path 	string	true   "tenantId"
// @Param   appInstanceId   path 	string	true   "appInstanceId"
// @Param   imageId         path 	string	true   "imageId"
// @Param   access_token    header  string  true   "access token"
// @Success 200 ok
// @Failure 404 image doesn't exist
// @router /tenants/:tenantId/app_instances/:appInstanceId/images/:imageId [get]
func (c *ImageController) GetImage() {
	log.Info("Query image request received.")
}

// @Title Download Image file
// @Description download a specific chunk of image file
// @Param   tenantId        path 	string	true   "tenantId"
// @Param   appInstanceId   path 	string	true   "appInstanceId"
// @Param   imageId         path 	string	true   "imageId"
// @Param   access_token    header  string  true   "access token"
// @Param   chunk_num       header  string  true   "chunk number"
// @Success 200 ok
// @Failure 404 image or chunk doesn't exist
// @router /tenants/:tenantId/app_instances/:appInstanceId/images/:imageId/file [get]
func (c *ImageController) GetImageFile() {
	log.Info("Download image file request received.")
}
