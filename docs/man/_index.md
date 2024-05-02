---
title: otdfctl - OpenTDF Control Tool

command:
  name: otdfctl
  aliases: []
  flags:
    - name: host
      description: host:port of the Virtru Data Security Platform gRPC server
      default: localhost:8080
    - name: log-level
      description: log level
      enum:
        - debug
        - info
        - warn
        - error
        - fatal
        - panic
      default: info
    - name: with-client-creds-file
      description: path to a JSON file containing a `clientId` and `clientSecret` for authentication via client-credentials flow
---
