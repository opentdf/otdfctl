# Example: Adding a New Subcommand to an Existing Command

This guide demonstrates how to add a `hello` subcommand under the existing `dev` command in `otdfctl`, using the project's CLI helpers and conventions for arguments and flags.

---

## 1. Add Documentation

Create a new file at `docs/man/dev/hello.md` with:

```markdown
---
title: Print Hello, <name>!
command:
  name: dev/hello
  arguments:
    - name: name
      description: Name to greet (optional, defaults to 'World')
  flags:
    - name: excited
      description: Print the greeting in uppercase
      default: false
---

Prints a greeting to the terminal. If a name is provided, it greets that name. If the `--excited` flag is set, the greeting is printed in uppercase.

### Examples

```
otdfctl dev hello
Hello, World!

otdfctl dev hello Alice
Hello, Alice!

otdfctl dev hello Bob --excited
HELLO, BOB!
```
```

---

## 2. Create the Subcommand in `cmd/dev.go`

Add the following code to `cmd/dev.go`:

```go
// ...existing imports...

var helloCmd = man.Docs.GetCommand("dev/hello", man.WithRun(func(cmd *cobra.Command, args []string) {
    c := cli.New(cmd, args)
    name := c.Args.GetOptionalString(0, "World")
    excited := c.Flags.GetOptionalBool("excited")
    msg := fmt.Sprintf("Hello, %s!", name)
    if excited {
        msg = strings.ToUpper(msg)
    }
    c.Println(msg)
}))

func init() {
    // Register flags using the doc-driven helpers:
    helloCmd.Flags().StringP(
        helloCmd.GetDocFlag("name").Name,
        helloCmd.GetDocFlag("name").Shorthand,
        helloCmd.GetDocFlag("name").Default,
        helloCmd.GetDocFlag("name").Description,
    )
    helloCmd.Flags().Bool(
        helloCmd.GetDocFlag("excited").Name,
        helloCmd.GetDocFlag("excited").DefaultAsBool(),
        helloCmd.GetDocFlag("excited").Description,
    )
    devCmd.AddCommand(helloCmd)
}
```

- **Arguments:** `[name]` (optional) — the name to greet. Defaults to `World` if not provided.
- **Flags:** `--excited` — prints the greeting in uppercase if set.
- **Helpers Used:** `cli.New`, `c.Args.GetOptionalString`, `c.Flags.GetOptionalBool`, `c.Println`, and `GetDocFlag` for doc-driven flag registration.

---

## 3. Test the Subcommand

Run:

```sh
otdfctl dev hello
otdfctl dev hello Alice
otdfctl dev hello Bob --excited
```

You should see:

```
Hello, World!
Hello, Alice!
HELLO, BOB!
```

---

For more complex logic, implement the business logic in `pkg/handlers/` and call it from your subcommand.
