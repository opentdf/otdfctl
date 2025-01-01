---
title: List subject mappings
command:
  name: list
  aliases:
    - l
  flags:
    - name: limit
      shorthand: l
      description: Limit retrieved count
    - name: offset
      shorthand: o
      description: Offset (page) quantity from start of the list
---

For more information about subject mappings, see the `subject-mappings` subcommand.

## Example

```shell
otdfctl policy subject-mappings list
```

```shell
  SUCCESS   Found subject-mappings list                                                                                                                                                                                                                                                                                                                                                                                                              
╭─────────────────────────────────────────┬─────────────────────────────────┬─────────────────────────┬─────────────────┬─────────────────────────────────┬─────────────────────────┬─────────┬─────────┬──────────────────╮
│ID                                       │Subject AttrVal: Id              │Subject AttrVal: Value   │Actions          │Subject Condition Set: Id        │Subject Condition Set    │Labels   │Created …│Updated At        │
├─────────────────────────────────────────┼─────────────────────────────────┼─────────────────────────┼─────────────────┼─────────────────────────────────┼─────────────────────────┼─────────┼─────────┼──────────────────┤
│d71c4028-ce64-453b-8aa7-6edb45fbb848     │891cfe85-b381-4f85-9699-5f7dbfe2…│myvalue1                 │[{"Value":{"Stan…│738736ee-880d-40da-acae-672d1def…│[{"condition_groups":[{"…│[]       │Wed Dec …│Wed Dec 18 15:41:…│
│39866dd2-368b-41f6-b292-b4b68c01888b     │891cfe85-b381-4f85-9699-5f7dbfe2…│myvalue1                 │[{"Value":{"Stan…│8dc98f65-5f0a-4444-bfd1-6a818dc7…│[{"condition_groups":[{"…│[]       │Wed Dec …│Wed Dec 18 15:40:…│
│e6a3f940-e24f-4383-8763-718a1a304948     │2fe8dea1-3555-498c-afe9-99724f35…│value2                   │[{"Value":{"Stan…│798aacd2-abaf-4623-975e-3bb8ca43…│[{"condition_groups":[{"…│[]       │Fri Nov …│Fri Nov  1 14:46:…│
│9d06c757-06b9-4713-8fbd-5ef007b1afe2     │74babca6-016f-4f3e-a99b-4e46ea8d…│value1                   │[{"Value":{"Stan…│eaf866c0-327f-4826-846a-5041c3c2…│[{"condition_groups":[{"…│[]       │Fri Nov …│Fri Nov  1 14:46:…│
╰─────────────────────────────────────────┴─────────────────────────────────┴─────────────────────────┴─────────────────┴─────────────────────────────────┴─────────────────────────┴─────────┴─────────┴──────────────────╯
  NOTE   Use 'otdfctl policy subject-mappings get --id=<id> --json' to see all properties 
```
