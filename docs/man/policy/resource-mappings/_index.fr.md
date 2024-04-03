---
title: Gérer les mappages de ressources
command:
  name: resource-mappings
---

# Gérer les mappages de ressources

Les mappages de ressources sont utilisés pour mapper les ressources à leurs valeurs d'attribut respectives en fonction des termes qui sont liés aux données. Seul, ce service n'est pas très utile, mais lorsqu'il est combiné avec un PEP ou un PDP qui peut utiliser les mappages de ressources, il devient un outil puissant pour automatiser le contrôle d'accès.

Par exemple, le PDP de marquage utilise des mappages de ressources pour mapper les ressources en fonction des termes trouvés dans les métadonnées et les documents qui lui sont envoyés. Combiné avec les mappages de ressources, il peut alors déterminer quels attributs doivent être appliqués au TDF et renvoyer ces attributs au PEP.
