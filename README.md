# kagami

[//]: # ([![]&#40;https://img.shields.io/github/actions/workflow/status/ditschedev/kagami/test.yml?branch=main&longCache=true&label=Test&logo=github%20actions&logoColor=fff&#41;]&#40;https://github.com/ditschedev/swag-ts/actions?query=workflow%3ATest&#41;)
[![Go Report Card](https://goreportcard.com/badge/github.com/ditschedev/kagami)](https://goreportcard.com/report/github.com/ditschedev/kagami)

Kagami is a simple CLI tool for mirroring git repositories from one provider to another.
It is versatile and easy to use.

## Installation

```bash
go install github.com/ditschedev/kagami@latest
```

## Configuration
The configuration is done using a yaml file. The file should contain a list of repositories 
that should be mirrored. Here is an example configuration file:

```yaml
repositories:
  - name: my-repo
    remote_uri: "git@github.com:ditschedev/kagami.git"
    mirror_uri: "git@gitlab.com:ditschedev/kagami.git"
```

The configuration file can be passed to the CLI using the `-c` flag.

### Authentication
Currently only SSH authentication is supported. 
Make sure you have the correct SSH keys set up for the repositories you want to mirror and have 
your ssh agent up and running.

## Cron

Kagami can be used in a cron job to mirror repositories periodically. 
Simply provide a path to the configuration file and add an entry to your crontab.

```bash
crontab -e
```

Add the following line to make kagami run every 5 minutes:

```bash
*/5 * * * * kagami mirror -c /path/to/repositories.yml
```

