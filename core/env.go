/*
 * Copyright 2021-2023 JetBrains s.r.o.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package core

import (
	"fmt"
	cienvironment "github.com/cucumber/ci-environment/go"
	log "github.com/sirupsen/logrus"
	"os"
	"runtime"
	"strings"
)

const (
	qodanaEnv              = "QODANA_ENV"
	QodanaToken            = "QODANA_TOKEN"
	QodanaLicenseOnlyToken = "QODANA_LICENSE_ONLY_TOKEN"
	qodanaJobUrl           = "QODANA_JOB_URL"
	qodanaRemoteUrl        = "QODANA_REMOTE_URL"
	qodanaBranch           = "QODANA_BRANCH"
	qodanaRevision         = "QODANA_REVISION"
	qodanaCliContainerName = "QODANA_CLI_CONTAINER_NAME"
	qodanaCliContainerKeep = "QODANA_CLI_CONTAINER_KEEP"
	qodanaCliUsePodman     = "QODANA_CLI_USE_PODMAN"
	qodanaDockerEnv        = "QODANA_DOCKER"
	QodanaConfEnv          = "QODANA_CONF"
	QodanaToolEnv          = "QODANA_TOOL"
	QodanaDistEnv          = "QODANA_DIST"
	QodanaEndpoint         = "ENDPOINT"
	qodanaCorettoSdk       = "QODANA_CORETTO_SDK"
	androidSdkRoot         = "ANDROID_SDK_ROOT"
	QodanaLicenseEndpoint  = "LICENSE_ENDPOINT"
	QodanaLicense          = "QODANA_LICENSE"
	QodanaTreatAsRelease   = "QODANA_TREAT_AS_RELEASE"
	qodanaClearKeyring     = "QODANA_CLEAR_KEYRING"
)

func ExtractQodanaEnvironment() {
	ci := cienvironment.DetectCIEnvironment()
	qEnv := "qodana"
	if ci != nil {
		qEnv = strings.ReplaceAll(strings.ToLower(ci.Name), " ", "-")
		setEnv(qodanaJobUrl, validateJobUrl(ci.URL, qEnv))
		if ci.Git != nil {
			setEnv(qodanaRemoteUrl, validateRemoteUrl(ci.Git.Remote))
			setEnv(qodanaBranch, validateBranch(ci.Git.Branch, qEnv))
			setEnv(qodanaRevision, ci.Git.Revision)
		}
	}
	setEnv(qodanaEnv, fmt.Sprintf("%s:%s", qEnv, Prod.Version))
	setEnv(QodanaDistEnv, Prod.Home)
}

// bootstrap takes the given command (from CLI or qodana.yaml) and runs it.
func bootstrap(command string, project string) {
	if command != "" {
		log.Printf("Running %s...", command)
		var executor string
		var flag string
		switch runtime.GOOS {
		case "windows":
			executor = "cmd"
			flag = "/c"
		default:
			executor = "sh"
			flag = "-c"
		}

		if res := RunCmd(project, executor, flag, command); res > 0 {
			log.Printf("Provided bootstrap command finished with error: %d. Exiting...", res)
			os.Exit(res)
		}
	}
}

func setEnv(key string, value string) {
	if os.Getenv(key) == "" && value != "" {
		err := os.Setenv(key, value)
		if err != nil {
			return
		}
	}
}
