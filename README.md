# notification-file

[![CodeQL](https://github.com/zbalkan/notification-file/actions/workflows/codeql.yml/badge.svg?branch=master)](https://github.com/zbalkan/notification-file/actions/workflows/codeql.yml)

CrowdSec plugin that writes the alerts to a file so that any SIEM agent can consume.

## Summary

Alerts are saved to the path defined as `log_path` in `file.yaml`. Modify the settings based on your requirements, especially OS.

The log format is `ndjson`.

## Usage

Refer to CrowdSec documentation for [registering a plugin][def1].

### Example: Collect CrowdSec alerts with Wazuh

Create a group called `CrowdSec` and use the configuration below. Then add the servers, on which CrowdSec is installed, as a member to this group.

```xml
<agent_config>
  <localfile>
    <log_format>json</log_format>
    <location>/tmp/crowdsec_alerts.json</location>
  </localfile>
</agent_config>
```

Create a rule, give a it reasonable name and ID. Use the root rule below, then start adding more rules matching scenarios.

```xml
<group name="crowdsec,">
  <rule id="99999" level="3">
    <decoded_as>json</decoded_as>
    <field name="program">crowdsec</field>
    <description>Crowdsec: Messages grouped.</description>
  </rule>
 </group>
```

## Development

Refer to CrowdSec documentation for [writing a plugin][def2].

Build using `go build -o bin/notification-file` command on Linux, `go build -o bin/notification-file.exe` on Windows.

[def1]: https://docs.crowdsec.net/docs/notification_plugins/intro#registering-plugin-to-profile
[def2]: https://docs.crowdsec.net/docs/notification_plugins/writing_your_own_plugin
