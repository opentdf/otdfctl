---
title: Update a Subject Condition Set

command:
  name: update
  flags:
    - name: id
      description: The ID of the subject condition set to update
      shorthand: i
      required: true
    - name: subject-sets
      description: A JSON array of subject sets, containing a list of condition groups, each with one or more conditions
      shorthand: s
      default: ""
    - name: label
      description: "Optional metadata 'labels' in the format: key=value"
      shorthand: l
      default: ""
    - name: force-replace-labels
      description: Destructively replace entire set of existing metadata 'labels' with any provided to this command
      default: false
---
