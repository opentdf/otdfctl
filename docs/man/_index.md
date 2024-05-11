---
title: otdfctl - OpenTDF Control Tool

command:
  name: otdfctl
  aliases: []
  flags:
    - name: host
      description: host:port of the OpenTDF Platform gRPC server
      default: localhost:8080
    - name: insecure
      description: use insecure connection
      default: false
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
    - name: plaintext
      description: use plaintext connection
      default: false
    - name: with-client-creds-file
      description: path to a JSON file containing a 'clientId' and 'clientSecret' for auth via client-credentials flow
    - name: with-client-creds
      description: JSON string containing a 'clientId' and 'clientSecret' for auth via client-credentials flow
      default: ''
---
