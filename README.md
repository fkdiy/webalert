# Webalert

Webalert is a small command-line application written in Go that monitors
specific elements on web pages and sends email notifications when their
content changes.

It can be used to monitor things like product availability, server offers,
status messages, headlines, or any other content that can be identified by
a CSS selector.

Webalert runs continuously, periodically checks all configured targets and
keeps their last known state in memory.

## Features

* Monitor multiple web pages
* Select monitored elements using CSS selectors
* Configurable check interval
* Add random jitter to check intervals
* Detect changes to the selected text content
* Send email notifications via SMTP
* Send one email per change or combine changes into a digest
* Continue monitoring if individual targets fail
* Load a custom configuration file using `--config`
* Pre-built binaries for Linux AMD64 and ARM64
* No database or persistent storage required

## Installation

### Download a pre-built binary

Pre-built Linux binaries are available on the
[GitHub Releases](https://github.com/fkdiy/webalert/releases) page.

Choose the binary matching your system:

* `webalert-linux-amd64` for x86-64 systems
* `webalert-linux-arm64` for ARM64 systems

After downloading the binary, make it executable:

```bash
chmod +x webalert-linux-amd64
```

Then run Webalert:

```bash
./webalert-linux-amd64
```

Webalert looks for `webalert.config.yaml` in the current working directory
by default. See the [Configuration](#configuration) section for details.

### Build from source

Alternatively, Webalert can be built from source.

Clone the repository:

```bash
git clone https://github.com/fkdiy/webalert.git
cd webalert
```

Download the dependencies:

```bash
go mod download
```

Build the application:

```bash
go build -o webalert ./cmd/webalert
```

Then start Webalert:

```bash
./webalert
```

During development, it can also be run directly:

```bash
go run ./cmd/webalert
```

## Configuration

Webalert is configured using a YAML file.

By default, it looks for `webalert.config.yaml` in the current working
directory.

Create your configuration based on the example configuration included in
the repository.

```yaml
interval: 60
jitter: 0

email:
  mode: digest
  from: "webalert@example.com"
  recipients:
    - notify@me.com
  smtp:
    host: smtp.example.com
    port: 587
    username: user@example.com
    password: change-me

targets:
  - url: https://example.com
    selector: ".price"
  - url: https://anotherexample.com
    selector: ".instock"
```

### Check interval

`interval` specifies the base interval between checks, in seconds.

For example:

```yaml
interval: 60
```

checks all configured targets every minute when no jitter is configured.

Webalert performs an initial check when it starts. The results of this
check become the initial state, so starting the application does not
generate notifications for existing content.

### Jitter

`jitter` adds a random variation to the configured check interval.

For every check, Webalert randomly selects a value between the negative
and positive jitter value and adds it to the base interval.

For example:

```yaml
interval: 60
jitter: 15
```

results in a new check interval between 45 and 75 seconds after each
check.

A new jitter value is generated for every interval, so checks do not
occur at a fixed frequency.

Set `jitter` to `0` to disable jitter:

```yaml
jitter: 0
```

The jitter value must be zero or a positive number and must be smaller
than the configured `interval`.

### Email notifications

SMTP settings and notification behavior are configured in the `email`
section.

Webalert supports two notification modes.

#### Per-change notifications

```yaml
mode: per_change
```

A separate email is sent for every detected change.

If five targets change during the same check, five emails are sent.

#### Digest notifications

```yaml
mode: digest
```

All changes detected during the same check are combined into a single
email.

If five targets change during the same check, one email containing all
five changes is sent.

### Targets

Each target consists of a URL and a CSS selector:

```yaml
targets:
  - url: https://example.com
    selector: ".price"
```

Webalert retrieves the page and stores the text content of the selected
element.

On subsequent checks, the current text is compared with the previously
stored value. If they differ, Webalert registers a change and sends a
notification according to the configured email mode.

## Custom configuration file

A different configuration file can be specified using the `--config`
option:

```bash
./webalert --config /path/to/config.yaml
```

The same option can be used with `go run`:

```bash
go run ./cmd/webalert --config /path/to/config.yaml
```

Command-line help is available with:

```bash
./webalert --help
```

## Error handling

An error while checking an individual target does not stop Webalert.

Unreachable websites, missing selectors, and similar errors are logged,
while the remaining targets continue to be monitored.

Email delivery errors are also logged without terminating the monitoring
process.

If the configuration cannot be loaded at startup, Webalert exits because
it cannot operate without a valid configuration.

## State

Webalert intentionally keeps its state only in memory.

No database or files are used to persist previously observed values. When
Webalert is restarted, all configured targets are therefore treated as new
targets and their current content becomes the new initial state.

This keeps the application lightweight and makes it suitable for simple
monitoring tasks where persistent history is not required.

## Running on a server

Webalert is designed to run as a long-lived process and can be deployed as
a single Go binary.

A simple installation might look like this:

```text
/opt/webalert/
├── webalert
└── webalert.config.yaml
```

For continuous operation on a Linux server, running Webalert as a systemd
service is recommended.

## Security

The configuration file contains SMTP credentials and should **not** be
committed to version control.

Add your actual configuration file to `.gitignore` and only commit a
sanitized example configuration.

If supported by your email provider, use an application-specific password
instead of your primary account password.

## Project structure

```text
webalert/
├── cmd/
│   └── webalert/
│       └── main.go
├── internal/
│   ├── config/
│   ├── email/
│   ├── jitter/
│   └── monitor/
├── webalert.config.example.yaml
├── go.mod
├── go.sum
└── README.md
```

* `cmd/webalert` contains the application entry point and coordinates the
  individual components.
* `internal/config` loads and parses the YAML configuration.
* `internal/email` creates and sends email notifications.
* `internal/jitter` generates random variations for check intervals.
* `internal/monitor` checks targets and detects changes.

## License

This project is licensed under the MIT License.
