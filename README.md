# Qodana CLI [<img src="https://api.producthunt.com/widgets/embed-image/v1/top-post-badge.svg?post_id=304841&theme=dark&period=daily" alt="" align="right" width="190" height="41">](https://www.producthunt.com/posts/jetbrains-qodana)

[![JetBrains project](https://jb.gg/badges/official.svg)](https://confluence.jetbrains.com/display/ALL/JetBrains+on+GitHub)
[![Qodana](https://github.com/JetBrains/qodana-cli/actions/workflows/ci.yml/badge.svg)](https://github.com/JetBrains/qodana-cli/actions/workflows/ci.yml)
[![GoReport](https://goreportcard.com/badge/github.com/JetBrains/qodana-cli)][gh:goreport]
[![GitHub Discussions](https://img.shields.io/github/discussions/jetbrains/qodana)][jb:discussions]
[![Twitter Follow](https://img.shields.io/badge/follow-%40Qodana-1DA1F2?logo=twitter&style=social)][jb:twitter]

`qodana` is a simple cross-platform command-line tool to run [Qodana linters](https://www.jetbrains.com/help/qodana/docker-images.html) anywhere with minimum effort required.

#### tl;dr

[Install](https://github.com/JetBrains/qodana-cli/releases/latest) and run:

```console
qodana scan --show-report
```

You can also add the linter by its name with the `--linter` option (e.g. `--linter jetbrains/qodana-js`).

**Table of Contents**

<!-- toc -->
- [Installation](#Installation)
- [Usage](#usage)
- [Configuration](#configuration)
- [Why](#why)

<!-- tocstop -->

![qodana](https://user-images.githubusercontent.com/13538286/151153050-934c0f41-e059-480a-a89f-cd4b2ca7a930.gif)

## Installation

> 💡 The Qodana CLI is distributed and run as a binary. The Qodana linters with inspections are [Docker Images](https://www.jetbrains.com/help/qodana/docker-images.html) or, starting from version `2023.2`, your local/downloaded by CLI IDE installations (experimental support).
> - To run Qodana with a container (the default mode in CLI), you must have Docker or Podman installed and running locally to support this: https://www.docker.com/get-started, and, if you are using Linux, you should be able to run Docker from the current (non-root) user (https://docs.docker.com/engine/install/linux-postinstall/#manage-docker-as-a-non-root-user)
> - To run Qodana without a container, you must have the IDE installed locally to provide the IDE installation path to the CLI or specify the product code, and CLI will try to download the IDE automatically (experimental support).

#### macOS and Linux
##### Install with [Homebrew](https://brew.sh) (recommended)
```console
brew install jetbrains/utils/qodana
```
##### Install with our installer
```console
curl -fsSL https://jb.gg/qodana-cli/install | bash
```
Also, you can install `nightly` or any other version (e.g. `v2023.2.9`) the following way:
```
curl -fsSL https://jb.gg/qodana-cli/install | bash -s -- nightly
```

#### Windows
##### Install with [Windows Package Manager](https://learn.microsoft.com/en-us/windows/package-manager/winget/) (recommended)
```console
winget install -e --id JetBrains.QodanaCLI
```
##### Install with [Chocolatey](https://chocolatey.org)
```console
choco install qodana
```
##### Install with [Scoop](https://scoop.sh)
```console
scoop bucket add jetbrains https://github.com/JetBrains/scoop-utils
scoop install qodana
```

#### Anywhere else
Alternatively,
you can install the latest binary (or the apt/rpm/deb package)
from [this page](https://github.com/JetBrains/qodana-cli/releases/latest).

Or, if you have Go installed, you can install the latest version of the CLI with the following command:

```console
go install github.com/JetBrains/qodana-cli/v2023@main
```

## Usage

https://user-images.githubusercontent.com/13538286/233484685-b9225168-8379-41bf-b8c8-6149a324cea8.mp4

🎥 The "Get Started with Qodana CLI" video is [also available on YouTube](https://www.youtube.com/watch?v=RV1MFnURMP8).

### Prepare your project

Before you start using Qodana, you need to configure your project –
choose [a linter](https://www.jetbrains.com/help/qodana/linters.html) to use.
If you know what linter you want to use, you can skip this step.

Also, Qodana CLI can choose a linter for you. Just run the following command in your **project root**:

```console
qodana init
```

### Analyze your project

Right after you configured your project (or remember linter's name you want to run),
you can run Qodana inspections simply by invoking the following command in your project root:

```console
qodana scan
```

- After the first Qodana run, the following runs will be faster because of the saved Qodana cache in your project (defaults to `./<userCacheDir>/JetBrains/<linter>/cache`)
- The latest Qodana report will be saved to `./<userCacheDir>/JetBrains/<linter>/results` – you can find qodana.sarif.json and other Qodana artifacts (like logs) in this directory.

### View the report

After the analysis, the results are saved to `./<userCacheDir>/JetBrains/<linter>/results` by default.
Inside the directory `./<userCacheDir>/JetBrains/<linter>/results/report`, you can find a Qodana HTML report.
To view it in the browser, run the following command from your project root:

```console
qodana show
```

You can serve any Qodana HTML report regardless of the project if you provide the correct report path.

## Configuration

To find more CLI options run `qodana ...` commands with the `--help` flag.
If you want to configure Qodana or a check inside Qodana,
consider
using [`qodana.yaml` ](https://www.jetbrains.com/help/qodana/qodana-yaml.html) to have the same configuration on any CI you use and your machine.

> In some flags help texts you can notice that the default path contains `<userCacheDir>/JetBrains`. The `<userCacheDir>` differs from the OS you are running Qodana with.
> - macOS: `~/Library/Caches/`
> - Linux: `~/.cache/`
> - Windows: `%LOCALAPPDATA%\`
> Also, you can just run `qodana show -d` to open the directory with the latest Qodana report.

### init

Configure a project for Qodana

#### Synopsis

Configure a project for Qodana:
prepare Qodana configuration file by analyzing the project structure
and generating a default configuration qodana.yaml file.

```
qodana init [flags]
```

#### Options

```
  -f, --force                Force initialization (overwrite existing valid qodana.yaml)
  -h, --help                 help for init
  -i, --project-dir string   Root directory of the project to configure (default ".")
```

### scan

Scan project with Qodana

#### Synopsis

Scan a project with Qodana.
It runs one of Qodana Docker's images
(https://www.jetbrains.com/help/qodana/docker-images.html) and reports the results.

Note that most options can be configured via qodana.yaml (https://www.jetbrains.com/help/qodana/qodana-yaml.html) file.
But you can always override qodana.yaml options with the following command-line options.


```
qodana scan [flags]
```

#### Options

```
  -l, --linter string                   Use to run Qodana in a container (default). Choose linter (image) to use. Not compatible with --ide option. Available images are: jetbrains/qodana-jvm-community:2023.2, jetbrains/qodana-jvm:2023.2, jetbrains/qodana-jvm-android:2023.2, jetbrains/qodana-php:2023.2, jetbrains/qodana-python:2023.2, jetbrains/qodana-python-community:2023.2, jetbrains/qodana-js:2023.2, jetbrains/qodana-go:2023.2, jetbrains/qodana-dotnet:2023.2
      --ide string                      Use to run Qodana without a container. Path to the installed IDE, or a downloaded one: provide direct URL or a product code. Not compatible with --linter option. Available codes are QDNET, add -EAP part to obtain EAP versions
  -i, --project-dir string              Root directory of the inspected project (default ".")
  -o, --results-dir string              Override directory to save Qodana inspection results to (default <userCacheDir>/JetBrains/Qodana/<linter>/results)
      --cache-dir string                Override cache directory (default <userCacheDir>/JetBrains/Qodana/<linter>/cache)
      --report-dir string               Override directory to save Qodana HTML report to (default <userCacheDir>/JetBrains/<linter>/results/report)
      --print-problems                  Print all found problems by Qodana in the CLI output
      --clear-cache                     Clear the local Qodana cache before running the analysis
  -w, --show-report                     Serve HTML report on port
      --port int                        Port to serve the report on (default 8080)
      --yaml-name string                Override qodana.yaml name to use: 'qodana.yaml' or 'qodana.yml'
  -a, --analysis-id string              Unique report identifier (GUID) to be used by Qodana Cloud (default "<generated value>")
  -b, --baseline string                 Provide the path to an existing SARIF report to be used in the baseline state calculation
      --baseline-include-absent         Include in the output report the results from the baseline run that are absent in the current run
      --full-history --commit           Go through the full commit history and run the analysis on each commit. If combined with --commit, analysis will be started from the given commit. Could take a long time.
      --commit --script local-changes   Base changes commit to reset to, resets git and runs linter with --script local-changes: analysis will be run only on changed files since the given commit. If combined with `--full-history`, full history analysis will be started from the given commit.
      --fail-threshold string           Set the number of problems that will serve as a quality gate. If this number is reached, the inspection run is terminated with a non-zero exit code
      --disable-sanity                  Skip running the inspections configured by the sanity profile
  -d, --source-directory string         Directory inside the project-dir directory must be inspected. If not specified, the whole project is inspected
  -n, --profile-name string             Profile name defined in the project
  -p, --profile-path string             Path to the profile file
      --run-promo string                Set to 'true' to have the application run the inspections configured by the promo profile; set to 'false' otherwise (default: 'true' only if Qodana is executed with the default profile)
      --script string                   Override the run scenario (default "default")
      --stub-profile string             Absolute path to the fallback profile file. This option is applied in case the profile was not specified using any available options
      --apply-fixes                     Apply all available quick-fixes, including cleanup
      --cleanup                         Run project cleanup
      --property stringArray            Set a JVM property to be used while running Qodana using the --property property.name=value1,value2,...,valueN notation
  -s, --save-report                     Generate HTML report (default true)
  -e, --env stringArray                 Only for container runs. Define additional environment variables for the Qodana container (you can use the flag multiple times). CLI is not reading full host environment variables and does not pass it to the Qodana container for security reasons
  -v, --volume stringArray              Only for container runs. Define additional volumes for the Qodana container (you can use the flag multiple times)
  -u, --user string                     Only for container runs. User to run Qodana container as. Please specify user id – '$UID' or user id and group id $(id -u):$(id -g). Use 'root' to run as the root user (default: the current user)
      --skip-pull                       Only for container runs. Skip pulling the latest Qodana container
  -h, --help                            help for scan
```

### show

Show a Qodana report

#### Synopsis

Show (serve) the latest Qodana report.

Due to JavaScript security restrictions, the generated report cannot
be viewed via the file:// protocol (by double-clicking the index.html file).  
https://www.jetbrains.com/help/qodana/html-report.html 
This command serves the Qodana report locally and opens a browser to it.

```
qodana show [flags]
```

#### Options

```
  -d, --dir-only             Open report directory only, don't serve it
  -h, --help                 help for show
  -l, --linter string        Override linter to use
  -p, --port int             Specify port to serve report at (default 8080)
  -i, --project-dir string   Root directory of the inspected project (default ".")
  -r, --report-dir string    Specify HTML report path (the one with index.html inside) (default <userCacheDir>/JetBrains/<linter>/results/report)
```

### view

View SARIF files in CLI

#### Synopsis

Preview all problems found in SARIF files in CLI.

```
qodana view [flags]
```

#### Options

```
  -h, --help                help for view
  -f, --sarif-file string   Path to the SARIF file (default "./qodana.sarif.json")
```

### contributors

A command-line helper for Qodana pricing to calculate active contributors* in the given repository.

#### Synopsis

* An active contributor is anyone who has made a commit to any
  of the projects you’ve registered in Qodana Cloud within the last 90 days,
  regardless of when those commits were originally authored. The number of such
  contributors will be calculated using both the commit author information
  and the timestamp for when their contribution to the project was pushed.

** Ultimate Plus plan currently has a discount, more information can be found on https://www.jetbrains.com/qodana/buy/


```
qodana contributors [flags]
```

#### Options

```
  -d, --days int             Number of days since when to calculate the number of active contributors (default 30)
  -h, --help                 help for contributors
  -i, --project-dir string   Root directory of the inspected project (default ".")
```

### cloc

A command-line helper for project statistics: languages, lines of code. Powered by boyter/scc. For contributors, use "qodana contributors" command.

#### Synopsis

```
qodana cloc [flags]
```

#### Options

```
  -h, --help                      help for cloc
  -o, --output string             Output format, can be [tabular, wide, json, csv, csv-stream, cloc-yaml, html, html-table, sql, sql-insert, openmetrics] (default "tabular")
  -i, --project-dir stringArray   Project directory, can be specified multiple times to check multiple projects, if not specified, current directory will be used
```

## Why

![Comics by Irina Khromova](https://user-images.githubusercontent.com/13538286/151377284-28d845d3-a601-4512-9029-18f99d215ee1.png)

> 🖼 [Irina Khromova painted the illustration](https://www.instagram.com/irkin_sketch/)

Qodana linters are distributed via Docker images –
which become handy for developers (us) and users to run code inspections in CI.

But to set up Qodana in CI, one wants to try it locally first,
as there is some additional configuration tuning required that differs from project to project
(and we try to be as much user-friendly as possible).

It's easy to try Qodana locally by running a _simple_ command:

```console
docker run --rm -it -p 8080:8080 -v <source-directory>/:/data/project/ -v <output-directory>/:/data/results/ -v <caches-directory>/:/data/cache/ jetbrains/qodana-<linter> --show-report
```

**And that's not so simple**: you have to provide a few absolute paths, forward some ports, add a few Docker options...

- On Linux, you might want to set the proper permissions to the results produced after the container run – so you need to add an option like `-u $(id -u):$(id -g)`
- On Windows and macOS, when there is the default Docker Desktop RAM limit (2GB), your run might fail because of OOM (and this often happens on big Gradle projects on Gradle sync), and the only workaround, for now, is increasing the memory – but to find that out, one needs to look that up in the docs.
- That list could go on, but we've thought about these problems, experimented a bit, and created the CLI to simplify all of this.

**Isn't that a bit overhead to write a tool that runs Docker containers when we have Docker CLI already?** Our CLI, like Docker CLI, operates with Docker daemon via Docker Engine API using the official Docker SDK, so actually, our tool is our own tailored Docker CLI at the moment.

[gh:test]: https://github.com/JetBrains/qodana/actions/workflows/build-test.yml
[gh:goreport]: https://goreportcard.com/report/github.com/JetBrains/qodana-cli
[youtrack]: https://youtrack.jetbrains.com/issues/QD
[youtrack-new-issue]: https://youtrack.jetbrains.com/newIssue?project=QD&c=Platform%20GitHub%20Action
[jb:confluence-on-gh]: https://confluence.jetbrains.com/display/ALL/JetBrains+on+GitHub
[jb:discussions]: https://jb.gg/qodana-discussions
[jb:twitter]: https://twitter.com/Qodana
[jb:docker]: https://hub.docker.com/r/jetbrains/qodana
