---
title: Decrypt a TDF file
command:
  name: decrypt [file]
  flags:
    - name: out
      shorthand: o
      description: 'The file destination for decrypted content to be written instead of stdout.'
      default: ''
    - name: tdf-type
      shorthand: t
      description: Deprecated. TDF type is now auto-detected.
      default: ''
---

Decrypt a Trusted Data Format (TDF) file and output the contents to stdout or a file in the current working directory.

The first argument is the TDF file with path from the current working directory being decrypted.

## Examples

Various ways to decrypt a TDF file

```shell
# decrypt file and write to standard output
otdfctl decrypt hello.txt.tdf

# decrypt file and write to hello.txt file
otdfctl decrypt hello.txt.tdf -o hello.txt

# decrypt piped TDF content and write to hello.txt file
cat hello.txt.tdf | otdfctl decrypt -o hello.txt
```

Advanced piping is supported

```shell
$ echo "hello world" | otdfctl encrypt | otdfctl decrypt | cat
hello world
```