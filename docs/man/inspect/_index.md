---
title: Inspect a TDF file
command:
  name: inspect
  arguments:
    - name: file
      description: Path to the TDF file to inspect
      required: true
      type: string
---

# Inspect a TDF file

Prints the `manifest.json` of the specified TDF for inspection.

This is useful for development and administration.

## Example

```shell
$ otdfctl inspect example.tdf
```
