---
title: otdfctl - OpenTDF Control Tool

command:
  name: otdfctl
  flags:
    - name: version
      description: show version
      default: false
    - name: host
      description: Hostname of the platform (i.e. https://localhost)
      default:
    - name: tls-no-verify
      description: disable verification of the server's TLS certificate
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
    - name: with-client-creds-file
      description: path to a JSON file containing a 'clientId' and 'clientSecret' for auth via client-credentials flow
    - name: with-client-creds
      description: JSON string containing a 'clientId' and 'clientSecret' for auth via client-credentials flow
      default: ''
---
