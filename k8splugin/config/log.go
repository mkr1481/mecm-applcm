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
//KANAG: pls move this to log folder name to improve maintability and readability
// Beego log related configurations
package config

import (
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
	"os"
)

func init() {
//KANAG: Make it configuratble
	fileName := "/usr/app/log/k8splugin.log"
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
	if err == nil {
		err = file.Close()
		if err != nil {
			logrus.Error("Failed to close the log file")
			return
		}
		ioWriter := &lumberjack.Logger{
	//KANAG: Make them configurable
			Filename:   fileName,
			MaxSize:    20,   // megabytes
			MaxBackups: 50,   // max archived files
			MaxAge:     30,   // days
			Compress:   true, // compression
		}
		logrus.SetOutput(ioWriter)
	} else {
		logrus.Warn("Failed to log to file, using default stderr")
	}
	//KANAG: Make them configurable
	logrus.SetLevel(logrus.InfoLevel)
}
