---
title: otdfctl - OpenTDF Control Tool

command:
  name: otdfctl
  aliases: []
  flags:
    - name: host
      description: host:port of the Virtru Data Security Platform gRPC server
      default: localhost:8080
    - name: config-file
      description: config file (default is $HOME/.otdfctl.yaml)
      default: ''
    - name: log-level
      description: log level (debug, info, warn, error, fatal, panic)
      default: info
---
