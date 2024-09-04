---
title: Subject mappings
command:
  name: subject-mappings
  aliases:
    - subm
    - sm
    - submap
    - subject-mapping
---

# Manage subject mappings

As data is bound to fully qualified Attribute Values when encrypted within a TDF, Entities are entitled to Attribute Values through a mechanism called Subject Mappings.

A Subject Mapping (SM) is the relation of a Subject Condition Set (SCS, see `subject-condition-sets` command)
to an Attribute Value to determine a Subject's Entitlement to an Attribute Value.

Entities (Subjects, Users, Machines, etc.) are defined by a representation (Entity Representation) of their identity from an identity provider (idP).
The OpenTDF Platform is not itself an idP, and it utilizes the OpenID Connect (OIDC) protocol as well as idP pluggability to rely upon an Entity store
of truth outside the platform to represent Entity identities.
