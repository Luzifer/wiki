[![Go Report Card](https://goreportcard.com/badge/github.com/Luzifer/wiki)](https://goreportcard.com/report/github.com/Luzifer/wiki)
![](https://badges.fyi/github/license/Luzifer/wiki)
![](https://badges.fyi/github/downloads/Luzifer/wiki)
![](https://badges.fyi/github/latest-release/Luzifer/wiki)
![](https://knut.in/project-status/wiki)

# Luzifer / wiki

`wiki` is a small file-based Wiki implementation with web-editing capabilities and a Git backed storage for history of pages.

The goal of this project was to have a small application to be deployed without any dependencies to open a Wiki for note taking or documentation purpose. 

The software itself has no concept of users or authentication and is held as simple as possible. Saved pages are stored as plain Markdown file onto the local disk inside a Git repository which on the one hand can be used to backup the state (just add a remote and set up a cron to push changes) and on the other hand to recover contents if someone deleted contents from a page.

## Usage

```console
# wiki --help
Usage of wiki:
      --data-dir string    Directory to store data to (default "./data/")
      --listen string      Port/IP to listen on (default ":3000")
      --log-level string   Log level (debug, info, warn, error, fatal) (default "info")
      --version            Prints current version and exits
```

To use this you can
- download pre-build binaries from the [releases](https://github.com/Luzifer/wiki/releases)
- pull the [Docker image](https://hub.docker.com/r/luzifer/wiki)
- or `go get -u github.com/Luzifer/wiki` the project

Given you've used the binary you can now just execute `./wiki` and go to `http://localhost:3000`. Everything you save will be stored in the `./data` directory.
