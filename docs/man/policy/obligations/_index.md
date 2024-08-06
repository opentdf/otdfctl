---
title: Manage obligations
command:
  name: obligations
---

Commands to manage obligations within the platform.

Obligations are requirements that should be fulfilled by a Subject or Environment requesting access to Resource data. Obligations
are distinct from Entitlements because they are a post-entitlement decision.

For example, if a Subject user is entitled through a Subject Mapping to access data that is considered confidential intellectual
property, an obligation can be utilized to ensure the access path is through a PEP or flow that supports watermarking to enhance
security of the data post-decryption. It is a precondition of an obligation that a Subject requesting access must already be entitled.

Obligations are represented as platform attributes. They can be either added to TDFs directly, or derived via mapping.

**Example:**

As a platform admin, Alice already has the hierarchical attribute `https://namespace.io/attr/card_suits` with high-to-low values `ace, king, queen, jack`
that are successfully entitled to various Subject users at each intended level of access.

Alice wants to ensure data of the hierarchical attribute `https://namespace.io/attr/card_suits/value/queen`
is accessible to Subjects that are entitled to Queen or King data, but Alice wants an additional protection of `watermarking` when
the TDF'd data is decrypted.

To add a `watermarking` assurance when a PEP handles decrypted TDF data a `digital rights management` , she can utilize obligations.

**Option 1: Derived obligations**

To dynamically tie all TDF data of that attribute value to this new obligation to watermark, she will create an obligation
mapped in platform policy to the `https://namespace.io/attr/card_suits/value/queen` attribute value:

`policy obligations create --name drm --value watermarking --attr <attribute value id>`

From this point forward, any time an access request is made on a TDF containing the `/attr/card_suits/value/queen` attribute value, the obligation
will be derived and considered in the decision.

The same obligation can also be assigned multiple attribute values when derived, so `--attr` can take multiple attribute value IDs, and the
`assign` and `remove` subcommands exist on `obligations` as well to assign attribute values to already-created obligations.

**Option 2: TDF-only obligations**

To add this obligation to any TDF data without permanently associating every TDF containing an attribute value to the obligation (see derived above),
Alice can create the obligation without associating it to an attribute value:

`policy obligations create --name drm --value watermarking`

Without the derived assocation, the obligation can be added on TDF encrypt as any other attribute FQN. There is nothing stopping a PEP or admin from utilizing obligations in both manners in tandem as appropriate.
