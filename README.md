# notification-file

CrowdSec plugin that writes the alerts to a file so that any SIEM agent can consume.

Alerts are saved to the path defined as `log_path` in `file.yaml`. Modify the settings based on your requirements, especially OS.

The log format is `ndjson`.

## Usage

Refer to CrowdSec documentation for [registering a plugin][def].

## Development

Refer to CrowdSec documentation for [writing a plugin][def2].

[def]: https://docs.crowdsec.net/docs/notification_plugins/writing_your_own_plugin
[def2]: https://docs.crowdsec.net/docs/notification_plugins/writing_your_own_plugin
