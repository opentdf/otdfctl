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

`policy obligations create --name drm --value watermarking`

From this point forward, any time an access request is made on a TDF containing the `/attr/card_suits/value/queen` attribute value, the obligation
will be derived and considered in the decision.

The same obligation can also be assigned multiple attribute values when derived, so `--attr` can take multiple attribute value IDs, and the
`assign` and `remove` subcommands exist on `obligations` as well to assign attribute values to already-created obligations.

**Option 2: TDF-only obligations**

To add this obligation to any TDF data without permanently associating every TDF containing an attribute value to the obligation (see derived above),
Alice can create the obligation without associating it to an attribute value:

`policy obligations create --name drm --value watermarking`

Without the derived assocation, the obligation can be added on TDF encrypt as any other attribute FQN. There is nothing stopping a PEP or admin from utilizing obligations in both manners in tandem as appropriate.

**Evaluation: Fulfillment Conditions**

In either scenario above, resolution of obligation satisfaction is similar to an `anyOf` rule on an attribute definition. If the obligation for `drm` contains
several values, and only one of them is `watermark`, a PEP or environmental entity that successfully meets the admin-defined obligation fulfillment conditions
for specifically that obligation or any of the other child values of the `drm` obligation parent for the data (derived via attributes or on the TDF)
would result in a permitted access decision.

The fulfillment conditions of an obligation can be thought of as loosely similar to Condition Sets within Subject Mappings that drive entitlements.

To allow access to data with obligations (derived by assignment to attribute values or directly added to a TDF), an admin must define the conditions
an entity must meet as provided by the Entity Resolution Service (ERS) or OIDC token claims if running the platform without an ERS.

For example, if `drm:watermark` is a required obligation contextualizing a TDF access decision, a user must first be entitled to the data attributes via subject mappings,
but they must also meet the conditions that they're attempting to access through `Some_Cool_PEP`, which is known to the admin to support watermarking
as a feature and respect it as an obligation on decrypt. The admin should define through the `fulfillments` subcommand the conditions where the obligation
is fulfilled. In plain English, that would be: `if a user's access token indicates they are accessing through Some_Cool_PEP, the obligation has been fulfilled`.

TODO: fulfillment conditions instruction here that comes from protos

As you can see in the example above, the user's entity chain indicated they came through `Some_Cool_PEP` and therefore they were granted access alongside the obligation
for `Some_Cool_PEP` to receive in response and drive the watermarking behavior.
