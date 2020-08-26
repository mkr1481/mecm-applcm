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
package pluginAdapter

import (
	"context"
	"lcmbroker/util"
	"mime/multipart"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	chunkSize       = 1024
	rootCertificate = ""
)

// Plugin adapter which decides a specific client based on plugin info
type PluginAdapter struct {
	pluginInfo string
}

// Constructor of PluginAdapter
func NewPluginAdapter(pluginInfo string) *PluginAdapter {
	return &PluginAdapter{pluginInfo: pluginInfo}
}

// Instantiate application
func (c *PluginAdapter) Instantiate(pluginInfo string, host string, deployArtifact string,
	accessToken string, appInsId string) (error error, status string) {
	log.Info("Instantiation started")
	client, err := GetClient(pluginInfo)
	if err != nil {
		log.Errorf(util.FailedToCreateClient, err)
		return err, util.Failure
	}
	ctx, cancel := context.WithTimeout(context.Background(), util.Timeout*time.Second)
	defer cancel()

	status, err = client.Instantiate(ctx, deployArtifact, host, accessToken, appInsId)
	if err != nil {
		log.Error("failed to instantiate application")
		return err, util.Failure
	}
	log.Info("instantiation completed with status: ", status)
	return nil, status
}

// Query application
func (c *PluginAdapter) Query(pluginInfo string, host string) (status string, error error) {
	log.Info("Query started")
	client, err := GetClient(pluginInfo)
	if err != nil {
		log.Errorf(util.FailedToCreateClient, err)
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), util.Timeout*time.Second)
	defer cancel()

	status, err = client.Query(ctx, host)
	if err != nil {
		log.Errorf("failed to query information")
		return "", err
	}
	log.Info("query status: ", status)
	return status, nil
}

// Terminate application
func (c *PluginAdapter) Terminate(pluginInfo string, host string, accessToken string, appInsId string) (status string, error error) {
	log.Info("Terminate started")
	client, err := GetClient(pluginInfo)
	if err != nil {
		log.Errorf(util.FailedToCreateClient, err)
		return util.Failure, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), util.Timeout*time.Second)
	defer cancel()

	status, err = client.Terminate(ctx, host, accessToken, appInsId)
	if err != nil {
		log.Error("failed to terminate application")
		return util.Failure, err
	}

	log.Info("termination success with status: ", status)
	return status, nil
}

// Upload configuration
func (c *PluginAdapter) UploadConfig(pluginInfo string, file multipart.File, host string, accessToken string) (status string, error error) {
	log.Info("Upload config started")
	client, err := GetClient(pluginInfo)
	if err != nil {
		log.Errorf(util.FailedToCreateClient, err)
		return util.Failure, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), util.Timeout*time.Second)
	defer cancel()

	status, err = client.UploadConfig(ctx, file, host, accessToken)
	if err != nil {
		log.Error("failed to upload configuration")
		return util.Failure, err
	}

	log.Info("upload configuration is success with status: ", status)
	return status, nil
}

// Remove configuration
func (c *PluginAdapter) RemoveConfig(pluginInfo string, host string, accessToken string) (status string, error error) {
	log.Info("Remove config started")
	client, err := GetClient(pluginInfo)
	if err != nil {
		log.Errorf(util.FailedToCreateClient, err)
		return util.Failure, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), util.Timeout*time.Second)
	defer cancel()

	status, err = client.RemoveConfig(ctx, host, accessToken)
	if err != nil {
		log.Error("failed to remove configuration")
		return util.Failure, err
	}

	log.Info("remove configuration is success with status: ", status)
	return status, nil
}