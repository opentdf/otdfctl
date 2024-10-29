---
title: Match a subject or set of selectors to relevant subject mappings
command:
  name: match
  flags:
    - name: subject
      shorthand: s
      description: A Subject Entity Representation string (JSON or JWT, auto-detected)
      default: ''
    - name: selector
      shorthand: x
      description: "Individual selectors (i.e. '.department' or '.realm_access.roles[]') that may be found in SubjectConditionSets"
---

Given that Subject Mappings contain Subject Condition Sets (see either relevant command for more information), this tool can consume an Entity Representation
or a set of provided selectors to query the platform Policy for any relevant Subject Mappings.

Given an Entity Representation of a Subject via `--subject` (an OIDC Access Token JWT, or a JSON object such as from an Entity Resolution Service response),
this command will parse all possible valid selectors and check those for presence in any Subject Condition Set referenced on a Subject Mapping to an Attribute Value.

Given a set of selectors (`--selector`), this command will look for any Subject Mappings with Subject Condition Sets containing those same selectors.

> [!NOTE]
> The values of the selectors and any IN/NOT_IN/IN_CONTAINS logic of Subject Condition Sets is irrelevant to this command. Evaluation of any matched conditions
> is handled by the Authorization Service to determine entitlements. This command is specifically for management of policy - to facilitate lookup of current
> conditions driven by known selectors as a precondition for administration of entitlement given the logical *operators* of the matched conditions and their relations.
