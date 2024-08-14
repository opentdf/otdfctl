---
title: Create a Subject Condition Set

command:
  name: create
  aliases:
    - c
    - add
    - new
  flags:
    - name: subject-sets
      description: A JSON array of subject sets, containing a list of condition groups, each with one or more conditions
      shorthand: s
      required: true
      default: ''
    - name: subject-sets-file-json
      description: A JSON file with path from the current working directory containing an array of subject sets
      shorthand: j
      default: ''
      required: false
    - name: label
      description: "Optional metadata 'labels' in the format: key=value"
      shorthand: l
      default: ''
    - name: force-replace-labels
      description: Destructively replace entire set of existing metadata 'labels' with any provided to this command
      default: false
---

### Example Subject Condition Sets

`--subject-sets` example input:

```json
[
  {
    "condition_groups": [
      {
        "conditions": [
          {
            "operator": "SUBJECT_MAPPING_OPERATOR_ENUM_IN",
            "subject_external_values": ["CoolTool", "RadService", "ShinyThing"],
            "subject_external_selector_value": ".team.name"
          },
          {
            "operator": "SUBJECT_MAPPING_OPERATOR_ENUM_IN",
            "subject_external_values": ["marketing"],
            "subject_external_selector_value": ".org.name"
          }
        ],
        "boolean_operator": 1
      }
    ]
  }
]
```

ConditionGroup `boolean_operator` is driven through the API `CONDITION_BOOLEAN_TYPE_ENUM` definition:

| CONDITION_BOOLEAN_TYPE_ENUM | index value | meaning               |
| --------------------------- | ----------- | --------------------- |
| AND                         | 1           | all conditions met    |
| OR                          | 2           | any one condition met |

Condition `operator` is driven through the API `SUBJECT_MAPPING_OPERATOR_ENUM` definition,
and is evaluated by applying the `subject_external_selector_value` to the Subject entity
representation (token or Entity Resolution Service response) and comparing the logical operator
against the list of `subject_external_values`:

| SUBJECT_MAPPING_OPERATOR_ENUM | index value | meaning                      |
| ----------------------------- | ----------- | ---------------------------- |
| IN                            | 1           | any of the values found      |
| NOT_IN                        | 2           | none of the values found     |
| IN_CONTAINS                   | 3           | contains one of these values |

In the example SCS above, the Subject entity MUST BE represented with a token claim or ERS response
containing a field at `.team.name` identifying them as team name "CoolTool", "RadService", or "ShinyThing", AND THEY MUST ALSO have a field `org.name` of "marketing".

This structure if their team name was "CoolTool" might look like:

```json
{
  "team": {
    "name": "CoolTool" // could alternatively be RadService or ShinyThing
  },
  "org": {
    "name": "marketing"
  }
}
```

If the `.org.name` were `sales` instead, the condition set would not be met, and the Subject would
not be found to be entitled to the Attribute Value applicable to this Subject Condition Set via Subject Mapping between.
