---
command: subject-mappings
title: Manage subject mappings
---

# Subject Mappings

Relations between Attribute Values and Subject Condition Sets that define the allowed Actions.

If a User's properties match a Subject Condition Set, the corresponding Subject Mapping provides them a set of allowed Actions
on any Resource (data) containing the mapped Attribute Value. 

	Attribute Value  <------  Subject Mapping ------->  Subject Condition Set

	Subject Mapping: 
		- Attribute Value: associated Attribute Value that the Subject Mapping Actions are relevant to
		- Actions: permitted Actions a Subject can take on Resources containing the Attribute Value
		- Subject Condition Set: associated logical structure of external fields and values to match a Subject

Platform consumption flow:
Subject/User -> IdP/LDAP's External Fields & Values -> SubjectConditionSet -> SubjectMapping w/ Actions -> AttributeValue

Note: SubjectConditionSets are reusable among SubjectMappings and are available under separate 'policy' commands.