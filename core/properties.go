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
	"github.com/JetBrains/qodana-cli/v2023/cloud"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func getPropertiesMap(
	prefix string,
	eap bool,
	appInfoXml string,
	systemDir string,
	logDir string,
	confDir string,
	pluginsDir string,
	dotNet DotNet,
	deviceIdSalt []string,
	plugins []string,
	analysisId string,
	coverageDir string,
) map[string]string {
	properties := map[string]string{
		"-Dfus.internal.reduce.initial.delay":          "true",
		"-Didea.headless.enable.statistics":            strconv.FormatBool(cloud.Token.IsAllowedToSendFUS()),
		"-Didea.headless.statistics.max.files.to.send": "5000",
		"-Dinspect.save.project.settings":              "true",
		"-Djava.awt.headless":                          "true",
		"-Djava.net.useSystemProxies":                  "true",
		"-Djdk.attach.allowAttachSelf":                 "true",
		"-Djdk.module.illegalAccess.silent":            "true",
		"-Dkotlinx.coroutines.debug":                   "off",
		"-Dsun.io.useCanonCaches":                      "false",
		"-Dsun.tools.attach.tmp.only":                  "true",

		"-Didea.headless.statistics.device.id":   deviceIdSalt[0],
		"-Didea.headless.statistics.salt":        deviceIdSalt[1],
		"-Didea.platform.prefix":                 "Qodana",
		"-Didea.parent.prefix":                   prefix,
		"-Didea.config.path":                     quoteIfSpace(confDir),
		"-Didea.system.path":                     quoteIfSpace(systemDir),
		"-Didea.plugins.path":                    quoteIfSpace(pluginsDir),
		"-Didea.application.info.value":          quoteIfSpace(appInfoXml),
		"-Didea.log.path":                        quoteIfSpace(logDir),
		"-Didea.qodana.thirdpartyplugins.accept": "true",
		"-Dqodana.automation.guid":               quoteIfSpace(analysisId),

		"-XX:SoftRefLRUPolicyMSPerMB": "50",
		"-XX:MaxJavaStackTraceDepth":  "10000",
		"-XX:ReservedCodeCacheSize":   "512m",
		"-XX:CICompilerCount":         "2",
		"-XX:MaxRAMPercentage":        "70",

		"-Didea.job.launcher.without.timeout": "true",
	}
	if coverageDir != "" {
		properties["-Dqodana.coverage.input"] = quoteIfSpace(coverageDir)
	}
	if eap {
		properties["-Deap.login.enabled"] = "false"
	}
	if len(plugins) > 0 {
		properties["-Didea.required.plugins.id"] = strings.Join(plugins, ",")
	}
	if prefix == "WebStorm" {
		properties["-Dqodana.recommended.profile.resource"] = "qodana-js.recommended.yaml"
		properties["-Dqodana.starter.profile.resource"] = "qodana-js.starter.yaml"
	}
	if prefix == "Rider" {
		if Prod.is233orNewer() {
			properties["-Dqodana.recommended.profile.resource"] = "qodana-dotnet.recommended.yaml"
			properties["-Dqodana.starter.profile.resource"] = "qodana-dotnet.starter.yaml"
		}
		properties["-Didea.class.before.app"] = "com.jetbrains.rider.protocol.EarlyBackendStarter"
		properties["-Drider.collect.full.container.statistics"] = "true"
		properties["-Drider.suppress.std.redirect"] = "true"
		if dotNet.Project != "" {
			properties["-Dqodana.net.project"] = dotNet.Project
		} else if dotNet.Solution != "" {
			properties["-Dqodana.net.solution"] = dotNet.Solution
		}
		if dotNet.Configuration != "" {
			properties["-Dqodana.net.configuration"] = dotNet.Configuration
		}
		if dotNet.Platform != "" {
			properties["-Dqodana.net.platform"] = dotNet.Platform
		}
	}

	return properties
}

// GetProperties writes key=value `props` to file `f` having later key occurrence win
func GetProperties(opts *QodanaOptions, yamlProps map[string]string, dotNetOptions DotNet, plugins []string) []string {
	lines := []string{
		fmt.Sprintf("-Xlog:gc*:%s", quoteIfSpace(filepath.Join(opts.logDirPath(), "gc.log"))),
		`-Djdk.http.auth.tunneling.disabledSchemes=""`,
		"-XX:+HeapDumpOnOutOfMemoryError",
		"-XX:+UseG1GC",
		"-XX:-OmitStackTraceInFastThrow",
		"-ea",
	}
	treatAsRelease := os.Getenv(QodanaTreatAsRelease)
	if treatAsRelease == "true" {
		lines = append(lines, "-Deap.require.license=release")
	}

	cliProps, flags := opts.properties()
	for _, f := range flags {
		if f != "" && !Contains(lines, f) {
			lines = append(lines, f)
		}
	}

	props := getPropertiesMap(
		Prod.parentPrefix(),
		Prod.EAP,
		opts.appInfoXmlPath(Prod.IdeBin()),
		filepath.Join(opts.CacheDir, "idea", Prod.getVersionBranch()),
		opts.logDirPath(),
		opts.ConfDirPath(),
		filepath.Join(opts.CacheDir, "plugins", Prod.getVersionBranch()),
		dotNetOptions,
		getDeviceIdSalt(),
		plugins,
		opts.AnalysisId,
		opts.CoverageDir,
	)
	for k, v := range yamlProps { // qodana.yaml – overrides vmoptions
		if !strings.HasPrefix(k, "-") {
			k = fmt.Sprintf("-D%s", k)
		}
		props[k] = v
	}
	for k, v := range cliProps { // CLI – overrides anything
		if !strings.HasPrefix(k, "-") {
			k = fmt.Sprintf("-D%s", k)
		}
		props[k] = v
	}

	for k, v := range props {
		lines = append(lines, fmt.Sprintf("%s=%s", k, v))
	}

	sort.Strings(lines)

	return lines
}

// writeProperties writes the given key=value `props` to file `f` (sets the environment variable)
func writeProperties(opts *QodanaOptions) { // opts.confDirPath(Prod.Version)  opts.vmOptionsPath(Prod.Version)
	properties := GetProperties(opts, Config.Properties, Config.DotNet, getPluginIds(Config.Plugins))
	err := os.WriteFile(opts.vmOptionsPath(), []byte(strings.Join(properties, "\n")), 0o644)
	if err != nil {
		log.Fatal(err)
	}
	err = os.Setenv(Prod.vmOptionsEnv(), opts.vmOptionsPath())
	if err != nil {
		log.Fatal(err)
	}
}
