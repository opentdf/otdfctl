---
title: Selectors
command:
  name: selectors
  aliases:
    - sel
---

Commands to generate and test selectors on Subject Entity Representations. For more information, see the help manual for each subcommand
or additional context within Subject Condition Sets.

## Flattening Syntax

The platform maintains a very simple flattening library such that the below structure flattens into the key/value pairs beneath.

Subject input:

```json
{
  "key": "abc",
  "something": {
    "nested": "nested_value",
    "list": ["item_1", "item_2"]
  }
}
```

Flattened Selectors:

| Selector             | Value          | Significance              |
| -------------------- | -------------- | ------------------------- |
| ".key"               | "abc"          | specified field           |
| ".something.nested"  | "nested_value" | nested field              |
| ".something.list[0]" | "item_1"       | first index specifically  |
| ".something.list[]"  | "item_1"       | any index in the list     |
| ".something.list[1]" | "item_2"       | second index specifically |
| ".something.list[]"  | "item_2"       | any index in the list     |
