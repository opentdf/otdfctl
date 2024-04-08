---
title: Create a Subject Condition Set

command:
  name: create
  flags:
    - name: subject-sets
      description: A JSON array of subject sets, containing a list of condition groups, each with one or more conditions
      shorthand: s
      required: true
      default: ""
    - name: subject-sets-file-json
      description: A JSON file with path from $HOME containing an array of subject sets
      shorthand: j
      default: ""
      required: false
    - name: label
      description: "Optional metadata 'labels' in the format: key=value"
      shorthand: l
      type: string-slice
      default: ""
    - name: force-replace-labels
      description: Destructively replace entire set of existing metadata 'labels' with any provided to this command
      type: bool
      default: false
---
