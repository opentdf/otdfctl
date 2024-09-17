---
title: Update a Subject Condition Set

command:
  name: update
  aliases:
    - u
  flags:
    - name: id
      description: The ID of the subject condition set to update
      shorthand: i
      required: true
    - name: subject-sets
      description: A JSON array of subject sets, containing a list of condition groups, each with one or more conditions
      shorthand: s
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

Replace the existing conditional logic within an SCS with new conditional logic, passing either JSON directly or a JSON file.

For more information about subject condition sets, see the `subject-condition-sets` subcommand.
