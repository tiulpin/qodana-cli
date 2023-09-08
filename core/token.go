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
	"github.com/JetBrains/qodana-cli/cloud"
	log "github.com/sirupsen/logrus"
	"os"
)

func (o *QodanaOptions) loadToken(refresh bool) string {
	tokenFetchers := []func(bool) string{
		func(_ bool) string { return o.getTokenFromCliArgs() },
		func(_ bool) string { return o.getTokenFromEnv() },
		o.getTokenFromKeychain,
		func(_ bool) string { return o.getTokenFromUserInput() },
	}

	for _, fetcher := range tokenFetchers {
		if token := fetcher(refresh); token != "" {
			return token
		}
	}
	return ""
}

func (o *QodanaOptions) getTokenFromCliArgs() string {
	tokenFromCliArgs := o.getenv(QodanaToken)
	if tokenFromCliArgs != "" {
		log.Debug("Loaded token from CLI args environment")
		return tokenFromCliArgs
	}
	return ""
}

func (o *QodanaOptions) getTokenFromEnv() string {
	tokenFromEnv := os.Getenv(QodanaToken)
	if tokenFromEnv != "" {
		log.Debug("Loaded token from the environment variable")
		return tokenFromEnv
	}
	return ""
}

func (o *QodanaOptions) getTokenFromKeychain(refresh bool) string {
	log.Debugf("project id: %s", o.id())
	tokenFromKeychain, err := cloud.GetCloudToken(o.id())
	if err == nil && tokenFromKeychain != "" {
		WarningMessage(
			"Got %s from the system keyring, declare %s env variable or run %s to override it",
			PrimaryBold(QodanaToken),
			PrimaryBold(QodanaToken),
			PrimaryBold("qodana init -f"),
		)
		o.setenv(QodanaToken, tokenFromKeychain)
		log.Debugf("Loaded token from the system keyring with id %s", o.id())
		if !refresh {
			return tokenFromKeychain
		}
	}
	return ""
}

func (o *QodanaOptions) getTokenFromUserInput() string {
	if IsInteractive() {
		WarningMessage(cloud.EmptyTokenMessage)
		token := setupToken(o.ProjectDir, o.id())
		if token != "" {
			log.Debugf("Loaded token from the user input, saved to the system keyring with id %s", o.id())
			return token
		}
	}
	return ""
}

// ValidateToken checks if QODANA_TOKEN is set in CLI args, or environment or the system keyring, returns it's value.
func (o *QodanaOptions) ValidateToken(refresh bool) string {
	token := o.loadToken(refresh)
	client := cloud.NewQodanaClient()
	if projectName := client.ValidateToken(token); projectName == "" {
		if token != "" {
			ErrorMessage(cloud.InvalidTokenMessage)
			os.Exit(1)
		}
	} else {
		SuccessMessage("Linked project name: %s", projectName)
		o.setenv(QodanaToken, token)
		return token
	}
	return token
}
