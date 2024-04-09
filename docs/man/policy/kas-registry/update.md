---
title: Update a Key Access Server registration
command:
  name: update
  flags:
    - name: id
      shorthand: i
      description: ID of the Key Access Server registration
      required: true
    - name: uri
      shorthand: u
      description: URI of the Key Access Server
    - name: public-key-local
      shorthand: p
      description: Public key of the Key Access Server
    - name: public-key-remote
      shorthand: r
      description: URI of the public key of the Key Access Server
    - name: label
      description: "Optional metadata 'labels' in the format: key=value"
      shorthand: l
      default: ""
    - name: force-replace-labels
      description: Destructively replace entire set of existing metadata 'labels' with any provided to this command
      default: false
---
