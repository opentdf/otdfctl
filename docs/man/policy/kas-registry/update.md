---
title: Update a Key Access Server registration
command:
  name: update
  aliases:
    - u
  flags:
    - name: id
      shorthand: i
      description: ID of the Key Access Server registration
      required: true
    - name: uri
      shorthand: u
      description: URI of the Key Access Server
    - name: public-keys
      shorthand: c
      description: "(Deprecated: Use otdfctl policy kas-registry key) One or more 'cached' public keys saved for the KAS"
    - name: public-key-remote
      shorthand: r
      description: "(Deprecated: Use otdfctl policy kas-registry key) URI of the 'remote' public key of the Key Access Server"
    - name: name
      shorthand: n
      description: Optional name of the registered KAS (must be unique within Policy)
    - name: label
      description: "Optional metadata 'labels' in the format: key=value"
      shorthand: l
      default: ''
    - name: force-replace-labels
      description: Destructively replace entire set of existing metadata 'labels' with any provided to this command
      default: false
---

Update the `uri`, `metadata`, or name for a KAS registered to the platform.

:::warning Deprecated flags
`--public-keys` and `--public-key-remote` are deprecated and will be removed in an upcoming release.
Use `otdfctl policy kas-registry key` commands to manage KAS keys instead.
:::

If resource data has been TDFd utilizing key splits from the registered KAS, deletion from
the registry (and therefore any associated grants) may prevent decryption depending on the
type of grants and relevant key splits.

Make sure you know what you are doing.

For more information about registration of Key Access Servers, see the manual for `kas-registry`.

## Example

```shell
otdfctl policy kas-registry update --id 3c39618a-cd8c-48cf-a60c-e8a2f4be4dd5 --name example-kas2-newname
```
