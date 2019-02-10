/*
Copyright The Helm Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package helmpath

import (
	"github.com/casimir/xdg-go"
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
)

// New default helm home, with different paths for different OS:
//	- %APPDATA%\helm on Windows
//	- $HOME/Library/Preferences/helm on OSX
//  - $kXDG_CONFIG_HOME/helm (typically $HOME/.config/helm for Linux)
var defaultHelmHome = xdg.App{Name: "helm"}.ConfigPath("")

// Old default helm home, it's old good $HELM/.helm
var oldDefaultHelmHome = filepath.Join(homedir.HomeDir(), ".helm")

type DefaultConfigHomePath interface {
	xdgHomeExists() bool
	basicHomeExists() bool
}

type FSConfigHomePath struct{ DefaultConfigHomePath }

// Checks whether $XDG_CONFIG_HOME/helm exists
func (FSConfigHomePath) xdgHomeExists() bool {
	return DirExists(defaultHelmHome)
}

// Checks whether $HOME/.helm exists
func (FSConfigHomePath) basicHomeExists() bool {
	return DirExists(oldDefaultHelmHome)
}

var ConfigPath DefaultConfigHomePath = FSConfigHomePath{}

func DirExists(path string) bool {
	osStat, err := os.Stat(path)
	return err == nil && osStat.IsDir()
}

// Get configuration home dir.
func GetDefaultConfigHome() string {
	if ConfigPath.xdgHomeExists() || !ConfigPath.basicHomeExists() {
		return defaultHelmHome
	}
	return oldDefaultHelmHome
}
