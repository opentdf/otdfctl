---
title: Encrypt file or stdin as a TDF
command:
  name: encrypt [file]
  flags:
    - name: out
      shorthand: o
      description: The output file TDF in the current working directory instead of stdout ('-o file.txt' and '-o file.txt.tdf' both write the TDF as file.txt.tdf).
      default: ''
    - name: attr
      shorthand: a
      description: Attribute value Fully Qualified Names (FQNs, i.e. 'https://example.com/attr/attr1/value/value1') to apply to the encrypted data.
    - name: mime-type
      description: The MIME type of the input data. If not provided, the MIME type is inferred from the input data.
    - name: tdf-type
      shorthand: t
      description: The type of tdf to encrypt as. TDF3 supports structured manifests and larger payloads. Nano has a smaller footprint and more performant, but does not support structured manifests or large payloads.
      enum:
        - tdf3
        - ztdf
        - nano
      default: tdf3
    - name: ecdsa-binding
      description: For nano type containers only, enables ECDSA policy binding
    - name: kas-url-path
      description: URL path to the KAS service at the platform endpoint domain. Leading slash is required if needed.
      default: /kas
    - name: with-assertions
      description: >
        EXPERIMENTAL: JSON string containing list of assertions to be applied during encryption. example - '[{"id":"assertion1","type":"handling","scope":"tdo","appliesToState":"encrypted","statement":{"format":"json+stanag5636","schema":"urn:nato:stanag:5636:A:1:elements:json","value":"{\"ocl\":\"2024-10-21T20:47:36Z\"}"}}]'
---

Build a Trusted Data Format (TDF) with encrypted content from a specified file or input from stdin utilizing OpenTDF platform.

## Examples

Various ways to encrypt a file

```shell
# output to stdout
otdfctl encrypt hello.txt

# output to hello.txt.tdf
otdfctl encrypt hello.txt --out hello.txt.tdf

# encrypt piped content and write to hello.txt.tdf
cat hello.txt | otdfctl encrypt --out hello.txt.tdf
```

Automatically append .tdf to the output file name

```shell
$ cat hello.txt | otdfctl encrypt --out hello.txt; ls
hello.txt  hello.txt.tdf

$ cat hello.txt | otdfctl encrypt --out hello.txt.tdf; ls
hello.txt  hello.txt.tdf
```

Advanced piping is supported

```shell
$ echo "hello world" | otdfctl encrypt | otdfctl decrypt | cat
hello world
```

## Attributes

Attributes can be added to the encrypted data. The attribute value is a Fully Qualified Name (FQN) that is used to
restrict access to the data based on entity entitlements.

```shell
# output to hello.txt.tdf with attribute
otdfctl encrypt hello.txt --out hello.txt.tdf --attr https://example.com/attr/attr1/value/value1
```

## NanoTDF

NanoTDF is a lightweight TDF format that is more performant and has a smaller footprint than TDF3. NanoTDF does not
support structured manifests or large payloads.

```shell
# output to nano.tdf
otdfctl encrypt hello.txt --tdf-type nano --out hello.txt.tdf
```
