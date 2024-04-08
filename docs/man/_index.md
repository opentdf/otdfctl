---
title: tructl - OpenTDF Control Tool

command:
  name: tructl
  aliases: []
  flags:
    - name: json
      description: output single command in JSON (overrides configured output format)
      default: "false"
    - name: host
      description: host:port of the Virtru Data Security Platform gRPC server
      default: localhost:8080
    - name: config-file
      description: config file (default is $HOME/.tructl.yaml)
      default: ""
    - name: log-level
      description: log level (debug, info, warn, error, fatal, panic)
      default: info
---
