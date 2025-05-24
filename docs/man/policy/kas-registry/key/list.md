---
title: List Keys
command:
  name: list
  aliases:
    - l
  flags:
    - name: limit
      shorthand: l
      description: Maximum number of keys to return
      required: true
    - name: offset
      shorthand: o
      description: Number of keys to skip before starting to return results
      required: true
    - name: alg
      shorthand: a
      description: Key Algorithm to filter for
    - name: kasId
      shorthand: i
      description: The id of the kas to filter for
    - name: kasName
      shorthand: n
      description: The name of the kas to filter for
    - name: kasUri
      shorthand: u
      description: The uri of the kas to filter for
---

List KAS Keys. You can filter based on algorithm, and kas id, name, or uri.
