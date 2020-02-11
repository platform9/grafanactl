# Grafanactl

Grafanactl is yet another CLI utility to interact with Grafana's APIs.

This specific tool was created for the purpose of syncing hand-crafted dashboards across multiple grafana instances, or organizations.

This tool is inspired by, and overlaps heavily in functionality with the following:

- [overdrive3000/grafanactl](https://github.com/overdrive3000/grafanactl)
- [retzkek/grafanactl](https://github.com/retzkek/grafanactl)
- [grafana-tools/sdk](https://github.com/grafana-tools/sdk)

While it was desirable to set this apart from the others with a unique name, `grafanactl` succinctly describes the abilities of this tool, as paralleled by the use of `<tool>ctl` in many other projects.

## Usage

Current Command List

```bash
# Listing Dashboards
grafanactl dashboard search

# Downloading dashboards
grafanactl dashboard download --all
grafanactl dashboard download --all -t dashboards

# Uploading dashboards
grafanactl dashboard upload -f dashboards

# List folders
grafanactl folder search
```

## Configuration

Grafanactl supports a configuration file with the same input parameters as flags.

Multiple definitions of the same configuration item will overwrite each other.
The following is the order of precedence for config locations.

1. Flags set at runtime
2. Environment Variables
3. .grafanactl.yaml (current working dir)
4. $HOME/.grafanactl.yaml

It's important to note that any configuration option can be set via any method.

### Config File

The config file supports JSON, TOML, YAML, HCL, INI, envfile or Java properties formats. [1](https://github.com/spf13/viper#why-viper)

Here is a sample yaml file:

```yaml
apikey: DEFINITELYNOTYOURAPIKEY
url: https://grafana.your.domain
```

### Environment Variables

Environment variables should be set with a `GS_` prefix. This is to avoid collission with other programs.

Example:

```bash
# grafanactl.rc
export GRAFANA_APIKEY=DEFINITELYNOTYOURAPIKEY
export GRAFANA_URL=https://grafana.your.domain
```

### Flags

All flags can be discovered using --help on any of the subcommands.

Example:

```bash
# grafanactl dashboard download --help
Download dashboards from a grafana instance

Usage:
  grafanactl dashboard download [flags]

Flags:
  -a, --all             Download all dashboards
  -h, --help            help for download
  -t, --target string   Target directory to save dashboard files. (default ".")

Global Flags:
      --apikey string   A Grafana API Key
      --config string   config file - default in order of precedence:
                        - .grafanactl.yaml
                        - $HOME/.grafanactl.yaml
      --url string      The URL of a Grafana instance
```
