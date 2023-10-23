/******************************************************************************
 * Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
 * PilotGo is licensed under the Mulan PSL v2.
 * You can use this software accodring to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: zhanghan
 * Date: 2021-11-18 13:03:16
 * LastEditTime: 2022-04-06 14:22:15
 * Description: provide configure yaml functions.
 ******************************************************************************/
package config

import (
	"fmt"
	"io"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

const (
	configType = "yaml"
)

func Load(file string, config interface{}) error {
	bytes, err := os.ReadFile(file)
	if err != nil {
		fmt.Printf("open %s failed! err = %s\n", file, err.Error())
		return err
	}

	err = yaml.Unmarshal(bytes, config)
	if err != nil {
		fmt.Printf("yaml Unmarshal %s failed!\n", string(bytes))
		return err
	}
	return nil
}

func Init(output io.Writer, configFile string) error {
	if output == nil {
		output = io.Discard
	}

	viper.SetConfigFile(configFile)
	viper.SetConfigType(configType) // or viper.SetConfigType("YAML")
	err := viper.ReadInConfig()     // Find and read the config file
	if err != nil {                 // Handle errors reading the config file
		return err
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Fprintf(output, "Config file changed %s \n", e.Name)
	})
	return nil
}

func MustInit(output io.Writer, conf string) { // MustInit if fail panic
	if err := Init(output, conf); err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
}
