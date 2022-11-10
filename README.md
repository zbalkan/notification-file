# notification-file

[![CodeQL](https://github.com/zbalkan/notification-file/actions/workflows/codeql.yml/badge.svg?branch=master)](https://github.com/zbalkan/notification-file/actions/workflows/codeql.yml)

CrowdSec plugin that writes the alerts to a file so that any SIEM agent can consume.

## Summary

Alerts are saved to the path defined as `log_path` in `file.yaml`. Modify the settings based on your requirements, especially OS.

The log format is `ndjson`.

## Usage

Refer to CrowdSec documentation for [registering a plugin][def].

## Development

Refer to CrowdSec documentation for [writing a plugin][def2].

Build using `go build -o bin/notification-file` command on Linux, `go build -o bin/notification-file.exe` on Windows.

[def]: https://docs.crowdsec.net/docs/notification_plugins/intro#registering-plugin-to-profile
[def2]: https://docs.crowdsec.net/docs/notification_plugins/writing_your_own_plugin
