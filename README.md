# Grafana Sync

Grafana Sync is a tool that enables replication of dashboards across
multiple grafana instances, or organizations.

You can download dashboards for a specific org, or folder.

You can upload dashboards to a specific org, preserving folder structure.

## Usage

Current Command List

```bash
# Listing Dashboards
grafana-sync dashboard list
grafana-sync dashboard search

# Downloading dashboards
grafana-sync dashboard download --all
grafana-sync dashboard download --all -t dashboards

# Uploading dashboards
grafana-sync dashboard upload -f dashboards

# List folders
grafana-sync folder list
grafana-sync folder search
```

## Configuration

Grafana supports a configuration file with the same input parameters as flags.

Multiple definitions of the same configuration item will overwrite each other.
The following is the order of precedence for config locations.

1. Flags set at runtime
2. Environment Variables
3. .grafana-sync.yaml (current working dir)
4. $HOME/.grafana-sync.yaml

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
# grafana-sync.rc
export GS_APIKEY=DEFINITELYNOTYOURAPIKEY
export GS_URL=https://grafana.your.domain
```

### Flags

All flags can be discovered using --help on any of the subcommands.

Example:

```bash
# grafana-sync dashboard download --help
Download dashboards from a grafana instance

Usage:
  grafana-sync dashboard download [flags]

Flags:
  -a, --all             Download all dashboards
  -h, --help            help for download
  -t, --target string   Target directory to save dashboard files. (default ".")

Global Flags:
      --apikey string   A Grafana API Key
      --config string   config file - default in order of precedence:
                        - .grafana-sync.yaml
                        - $HOME/.grafana-sync.yaml
      --url string      The URL of a Grafana instance
```
