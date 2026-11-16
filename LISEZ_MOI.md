# evoCell : Une couche d'agents cellulaires pour les moteurs évolutifs symboliques

---

## **Description**

**evoCell** est un module de démonstration (**preuve de concept**) d'**agents cellulaires** développé en **Go**, fondé sur le **modèle d'acteurs** (via [go-actor](https://github.com/vladopajic/go-actor)) et conçu pour s'hybrider avec des moteurs évolutifs tels qu'**[evoGo](https://github.com/svarin92/evoGo)**. À la différence d'un moteur d'évolution grammaticale qui fait évoluer des génomes complets, evoCell opère en maintenant un **pool génique** : une population vivante d'agents cellulaires qui cultivent, filtrent et font circuler des *fragments* de matériel génétique entre les individus.

Dans cette architecture bio-inspirée, la **Membrane** n'est pas une simple barrière physique : c'est l'interface de communication qui définit la manière dont chaque agent interagit avec son écosystème. Chaque cellule est un acteur autonome dont la membrane remplit trois fonctions critiques :

- **L'excrétion** (transfert horizontal de gènes) : une cellule performante extrait un fragment de son génome et l'envoie à ses voisines, propageant son succès comme une *infection bénéfique*.
- **L'absorption** (filtrage immunitaire) : aucun matériel entrant n'est intégré directement ; la membrane le route vers un **système immunitaire** qui rend un verdict Soi/non-Soi avant toute transgénèse.
- **La signalisation** (quorum sensing) : les cellules diffusent des signaux chimiques pour influencer leurs voisines — un **signal de stress** élève le taux de mutation de la colonie (chercher une solution *pour la colonie*), un **signal de stabilité** apaise les mutations agressives et stabilise la zone.

### **Caractéristiques principales**

- **Cellules souches et différenciation par donnée** : chaque `Cell` est totipotente. Par la greffe d'un `Nucleus` générique - configuré via une `ZoneSpec` lue dans la table des symboles de la grammaire BNF (Backus-Naur form) courante -, elle se spécialise pour une zone grammaticale donnée. Aucune catégorie n'est figée *a priori* dans le module (`VowelCell`, `ConsonantCell`, ou encore `VerbCell` sont absentes du code) : modifier la grammaire ne requiert pas de modifier une seule ligne d'evoCell.
- **Zonage grammatical** : la topologie du tissu est déduite automatiquement des non-terminaux de la grammaire du moteur hôte — chaque règle devient une zone, chaque production un contrat de voisinage. Le tissu orchestre la circulation et l'expression des fragments génétiques, tandis que la grammaire détermine la façon dont ces échanges s'articulent.
- **Porosité de recuit** (annealing) : les membranes sont poreuses en début de cycle, favorisant de nombreux échanges entre cellules pour maximiser l'exploration. En fin de cycle, elles deviennent étanches, ce qui permet de stabiliser et d'exploiter les meilleurs génomes. Un seul paramètre pilote à la fois les échanges inter-cellulaires et la division asymétrique.
- **Robustesse collective** : si un agent mute vers une forme « cancéreuse », les membranes apprennent à bloquer ses messages — une mémoire immunitaire collective intervient à la manière des anticorps contre un pathogène.

---

## **Fondements et inspirations**

**evoCell** s'inspire de :

- **[go-actor](https://github.com/vladopajic/go-actor)** : une bibliothèque Go légère pour le modèle d'acteurs, dont les *mailboxes* et les *workers* fournissent l'ossature concurrente du tissu.
- **Le transfert horizontal de gènes** et le **quorum sensing** en microbiologie : les bactéries, qui échangent des plasmides et synchronisent leur comportement via des signaux chimiques, inspirent directement le modèle de l'excrétion et de la signalisation.
- **La biologie des cellules souches** : la division asymétrique (production d'une cellule différenciée à partir d'une cellule souche) est préférée lors de l'organogenèse, à la division symétrique des tissus matures. Ce principe guide la division cellulaire dans evoCell.
- ***Le Phénomène de la vie* de [Hans Jonas](https://www.puf.com/le-phenomene-de-la-vie-0)** (PUF) : La philosophie organique de Hans Jonas irrigue directement l'ontologie d'evoCell. Pour Jonas, l'organisme est une exception : sa véritable identité réside dans un processus plutôt que dans une substance figée, sa forme persistant au cœur même des échanges incessants de matière. evoCell en emprunte l'identité comme cycle et non comme contenu (le circuit Harvest → Deposit → Recombine offre une traduction computationnelle du métabolisme) ; la membrane comme acte de délimitation constitue le soi (un dedans n'existe qu'en s'opposant à un dehors filtré) ; la liberté nécessiteuse et la téléologie immanente : de même que l'ouverture au milieu engendre la spécialisation chez le vivant (plasticité émergente), les régularités d'evoCell émergent des échantillons du monde au lieu d'être décrétées *a priori*.
- **[evoGo](https://github.com/svarin92/evoGo)** : le moteur évolutif symbolique avec lequel ce module s'hybride, dont le système immunitaire (`IImmune`) et la spécialisation dynamique des organismes (`IOrganism`) préfiguraient cette couche cellulaire.

---

## **Architecture**

### **Qu'est-ce qu'un système cellulaire ?**

Un système cellulaire est une population d'agents qui cultivent, filtrent et font circuler du matériel génétique *partiel* (fragments), organisée par domaine, **sans jamais produire ni évaluer un phénotype complet**. Formellement, c'est le quintuplet **⟨Z, M, I, P, T⟩** :

1. **Z — Zonage** : organisation spatiale par domaine (zones grammaticales déduites de la grammaire hôte) qui détermine *où* vivent les fragments.
2. **M — Membranes** : interfaces de communication (excrétion, absorption, signalisation) qui régissent *comment* circulent les fragments.
3. **I — Immunité** : mémoire des motifs toxiques et verdicts Soi/non-Soi qui définissent *quels* fragments sont refusés.
4. **P — Porosité** : profil de perméabilité recuit qui fixe *combien* de fragments circulent selon la phase du cycle.
5. **T(ick) — Horloge** : les *ticks* asynchrones entre les générations du moteur hôte — le système vit *pendant* que la population hôte attend.

L'unité de la population hôte est le génome complet (un candidat phénotype), tandis que celle du système cellulaire est le fragment. Le lien entre ces deux échelles n'est pas une association terme à terme entre un individu et sa composante, mais celui d'un flux qui, au début de chaque génération, fragmente les génomes entiers, pour ne retenir et greffer en retour que les éléments les plus performants.

### **La hiérarchie biologique**

Le module épouse les quatre niveaux de l'organisation biologique, du plus simple au plus complexe :

1. **Cellule** (`cell.go`) : un acteur portant un génome de codons, une *fitness* et un taux de mutation modulé par les signaux reçus. Totipotente par défaut ; un `Nucleus` greffable spécialise son métabolisme.
2. **Tissu** (`space.go`) : une zone grammaticale — une population de cellules spécialisées partageant un domaine. Il gère l'annuaire des adresses de membranes, les statistiques collectives et diffuse le quorum sensing.
3. **Organe** (`organ.go`) : une assemblée de zones liées par un **contrat fonctionnel**. Son rôle est de produire, par leur seul dialogue de quorum, un matériel génétique cohérent pour un non-terminal de rang supérieur (ex. l'organe « groupe nominal » coordonne les zones déterminant, nom et adjectif ; son indice de cohérence constitue le signal de haut niveau du tissu). Les organes valident ainsi les *assemblages* : même deux fragments parfaits peuvent encore former un ensemble incohérent.
4. **Organisme** : l'individu du moteur hôte — le génome complet et son phénotype, qu'aucune couche cellulaire ne produit jamais.

### **Les acteurs supports**
- **Membrane** (`membrane.go`) : l'unique interface d'échange de la cellule. Elle gère deux boîtes de réception (mailboxes) : une privée pour le métabolisme interne (cellule ↔ sa membrane) et une partagée pour l'écosystème du tissu (membrane ↔ tissu).
- **Système immunitaire** (`immune.go`) : le filtre de sécurité. Il évalue chaque fragment entrant à la lumière des motifs toxiques mémorisés et de la *fitness* moyenne du tissu.
- **Space** (`space.go`) : la structure d'accueil du tissu. Il maintient l'annuaire des adresses de membranes, agrège les statistiques collectives et orchestre la diffusion du quorum sensing.
- **Coordinateur du tissu** (`tissue.go`) : l'interface *endocrinienne* entre le moteur hôte et les cellules. Il régule la politique de division, où la probabilité de différenciation asymétrique est pilotée par la porosité courante.

### **Hybridation avec evoGo**

evoGo expose un contrat minimal, nil-safe (`interfaces/hybridization.go`) :

```go
type HybridizationHook interface {
	OnGenerationStart(generation int, population []IIndividual)
	OnGenerationEnd(generation int, bestEver IIndividual)
}
```

Une seule ligne de code suffit à hybrider le module. La population naît à l'intérieur de SearchLoop et possède déjà son propre Genomizer : l'attachement se fait donc par option, sans réinstanciation.

```go
best, err := ge.SearchLoop(..., ge.WithHybridizationHook(evogo.NewHook(cfg, obs)))
```

Sans hybridation, evoGo fonctionne comme un noyau dépourvu de membrane : il ne dépend jamais d'evoCell. L'adaptateur (`evogo/adapter.go`) reste le seul fichier à assurer le pont entre les deux mondes.

### **Le cycle de vie du tissu**

Le coordinateur et l'observateur génèrent des traces explicites entre deux générations d'evoGo. À partir de la grammaire `letter.bnf` et pour la cible `"golden"`, le déroulement d'un cycle pourrait produire la séquence suivante :

**Semis — `OnGenerationStart(1)`** : les 60 génomes GE sont décomposés en fragments et semés dans les zones déduites des non-terminaux de la grammaire :

```
[tissue] semis des zones depuis letter.bnf : {string, letter, string_tail, vowel, consonant}
[tissue] 60 génomes démontés -> 214 fragments (moy. 3,6/génome)
[zone:vowel]        41 fragments absorbés par 23 cellules
[zone:consonant]    87 fragments absorbés par 25 cellules
[zone:string_tail]  86 fragments absorbés par 12 cellules
```

**Métabolisme — *ticks* inter-générations** : les cellules excrètent des fragments (HGT), appliquent les verdicts immunitaires et échangent des signaux de quorum, sous le contrôle de la porosité recuite :

```
[tick g=12] porosité=0,79  signaux_stress=0  signaux_stabilité=3
[cell #14 zone:vowel]      excrétion : codons[3:9] fitness=0,83 -> 3 voisines (HGT)
[cell #7  zone:consonant]  absorption : de #14 -> IMMUNITÉ : fitness faible (0,31 < moy 0,57) REJETÉ
[cell #22 zone:consonant]  absorption : de #14 -> IMMUNITÉ : accepté (perméable, self-compatible)
[cell #22 zone:consonant]  transgénèse : codons[3:9] <- fragment de #14
[cell #31 zone:vowel]      SIGNAL stress (fitness 0,17 en chute) -> voisines : mutation x2,5
[cell #4  zone:vowel]      SIGNAL stabilité (fitness 0,92) -> voisines : apaisées (mutation x0,5)
[immune] motif toxique appris sig=a3f9 (l'émetteur #31 saturait 40 % du quota de mailbox)
```

Notons que la cellule #31 inondait ses voisines ; son motif est désormais présent dans la mémoire immunitaire, et ses messages sont bloqués dans tout le tissu.

**Récolte — `OnGenerationEnd(12)`** : des fragments coordonnés par organe sont greffés sur les individus GE (Grammatical Evolution) faibles :

```
[organe string] cohérence=0,71  zones={letter,string_tail} assemblages validés 6/8
[deposit] individu #37 faible sur <string_tail> <- fragment greffé (zone:string_tail, fitness 0,78)
[deposit] individu #52 faible sur <letter> <- fragment d'organe greffé (lettres cohérentes)
[pool] zone vowel : 41 fragments, diversité 0,63 | consonant : 87, diversité 0,58
```

**Le régime de secours (rescue)** — démonstration la plus claire de la valeur d'evoCell. Ici, alors qu'evoGo a atteint un plateau, un *hook* est greffé :

```
[génération 19] GE seul : meilleure fitness 0,50 (plateau depuis 7 générations)
[génération 20] TISSU GREFFÉ — semis depuis la population stagnante
[génération 21] [deposit] individu #8 <- fragment d'organe (origine HGT : cellule #14, g=12)
[génération 23] meilleure fitness 0,67  (stagnation brisée)
[génération 27] meilleure fitness 0,83
[génération 31] meilleure fitness 1,00  -> "golden"
```

La courbe de *fitness* repart précisément quand le pool génique commence à être exploré et exploité: c'est l'œuvre de l'*infection bénéfique*.

**Sorties de l'Observateur** : le graphe du tissu (`tissue.dot`) et les métriques du pool (`pool.csv`) complètent le `best_ever.dot` d'evoGo :

```
tissue.dot : clusters = zones grammaticales | couleur des nœuds = Nucleus greffé
             arêtes = transferts HGT, épaisseur ∝ porosité courante
pool.csv   : génération, zone, fragments, diversité, fitness moyenne des donneurs
```

---

## **Démarrage**

```go
import (
	"github.com/svarin92/evoGo/ge"

	"github.com/svarin92/evoCell/config"
	evogo "github.com/svarin92/evoCell/evogo"
	"github.com/svarin92/evoCell/renderer"
)

// Configuration du tissu.
cfg := config.Config{
	Cells:             50,
	MaxGenerations:    100,
	MutationBoost:     2.5,
	DivisionThreshold: 0.8,
	DeathThreshold:    0.2,
	Porosity: config.PorosityProfile{
		StartPermeability: 0.9,  // poreux : exploration
		EndPermeability:   0.05, // étanche : exploitation
		Annealing:         config.LinearAnnealing(0.9, 0.05),
	},
}

obs := renderer.NewObserver(cfg)

// Le tissu s'auto-sème depuis la population et vit entre les générations.
bestEver, err := ge.SearchLoop(
	100, 50, grammar, "golden",       
	replacementFunc, selectionFunc, fitnessFunc,    
	
	// Hybridation par option : la population naît dans SearchLoop, laquelle
	// possède déjà son propre Genomizer. L'observateur alimente le rendu 
	// (tissue.dot, pool.csv).
	ge.WithHybridizationHook(evogo.NewHook(cfg, obs)),          
)
```
---

## **Plasticité : la spécialisation comme heuristique**

L'évolution de l'architecture entre les deux moteurs traduit un changement de paradigme fondamental dans la manière de concevoir le traitement des formes et des symboles : **la spécialisation n'est pas une taxonomie figée, mais une heuristique émergente**.

Dans evoGo, les organismes (`IOrganism`) sont prédéterminés : le phénotype cible appartient à des catégories fixes (alphabétique, numérique, ontologique) ancrées dans la structure du moteur. C'est un modèle robuste, mais rigide. Dans evoCell, au contraire, la cellule ignore à l'avance les catégories qu'elle devra traiter. Les symboles non-terminaux sont dynamiques (la table des symboles grammaticaux est reconstruite à la volée par le `builder`), et le noyau greffé interprète en continu la grammaire active.

Ce mécanisme induit deux ruptures fondamentales :

- **La membrane est le lieu de la découverte.** Seule interface avec l'extérieur, elle convertit le contexte en spécialisation. Le moteur hôte repère une régularité morphologique parmi ses phénotypes performants — tel le motif « golden » décomposé en racine et en suffixe. Le signal de différenciation transmet alors sa spécification et regreffe les noyaux à la volée, s'affranchissant de toute recompilation. À l'image d'une cellule souche qui se spécialise sous l'effet des morphogènes de son milieu, le système trouve ici son équivalent : la table des symboles fait office de morphogène.
- **Le filtrage est une heuristique.** Lorsqu'un noyau consulte la mémoire immunitaire collective via son prédicat d'absorption, il opère comme un lymphocyte B : l'anticorps mémorisé — motif gagnant ou empreinte toxique — oriente sa spécialisation. Les zones ne filtrent plus selon des règles arbitraires, mais selon ce que l'expérience a déjà éprouvé. En capitalisant sur ses succès passés, la plasticité du système se fait moteur de découverte.

En pratique : une population cellulaire peut se réorienter en cours de cycle vers des sous-catégories totalement imprévues (telles que « suffixes », « racines » ou « voyelles longues »), dès lors que le tissu constate que ces classifications nourrissent les individus les plus vulnérables. C'est l'équivalent opérationnel de la métaplasie biologique — ce processus où une cellule différenciée se re-spécialise sous la contrainte de son environnement. Par ce mécanisme, *le tissu apprend sa propre anatomie*. 

### La perception : la membrane comme organe sensoriel

La plasticité puise à deux sources complémentaires : la commande (le moteur hôte identifie une régularité et ordonne la regreffe) et la perception — cette dernière renouant plus fidèlement avec la logique du vivant.

Car la membrane ne se réduit pas à une simple interface passive : elle agit comme un véritable organe sensoriel. Le tissu dispose ainsi d'un canal perceptif où le corpus brut (`golden`, `wooden`, `goldfish`, etc.) se diffuse à la manière d'un gradient de morphogène. Chaque noyau plastique y extrait des régularités par une heuristique non supervisée. Nulle étiquette n'est ici apposée (« ceci est une racine, cela un suffixe ») : c'est la distribution statistique qui en décide seule.

La frontière architecturale demeure ainsi inviolée. Le corpus transmet au système les structures du monde — ses récurrences morphologiques —, et non son but final (la cible de l'évaluateur ne transite jamais par ce canal). La cellule apprend simplement à décoder les régularités du langage, et non à imiter un résultat prédéterminé.

### La connaissance comme fruit de la plasticité

L'heuristique ainsi forgée n'est pas un paramètre opaque enfoui dans les rouages du système : elle devient une donnée de premier rang. Capturés au fil des récoltes — par des mécanismes de capture préservant l'intégrité des noyaux vivants —, les profils émergents sont agrégés par zone puis externalisés.

Ce matériau s'avère à la fois lisible par l'humain et consommable par le moteur hôte, qui peut ainsi enrichir sa propre bibliothèque de motifs à partir des découvertes du tissu.

C'est pourquoi la plasticité dépasse ici le simple mécanisme d'adaptation pour s'ériger en condition de l'auto-apprentissage. Un système à spécialisation figée peut s'optimiser, mais il ne peut pas apprendre. evoCell réunit les trois piliers d'une véritable cognition évolutive :
- L'encodage : la spécialisation fixe les régularités découvertes.
- L'explicitation : la connaissance est lisible, inspectable et exportable.
- La rétroaction : le filtrage redéfinit ce que la zone retient, infléchissant en retour la trajectoire des individus.

---

## **Feuille de route**

- [ ] Acteur coordinateur : intégrer la boîte aux lettres de commande et la gestion asynchrone des réponses Harvest.
- [ ] Politique de division : appliquer une division asymétrique régie par `P(asymétrique) = porosité`.
- [ ] Mécanisme de récolte et de transfert : implémenter la récolte par zone (`HarvestByZone`) et la greffe d'organes (`HarvestByOrgan → Deposit`) ciblant les individus faibles.
- [ ] Topologies du tissu : structurer la communication, d'abord par *broadcast*, puis par grilles de Moore associées à un rayon de quorum.
- [ ] Validation dynamique : mettre les acteurs au banc d'essais de [go-super-actor](https://github.com/vladopajic/go-super-actor).
- [ ] Modélisation du langage : concevoir un modèle de langage symbolique et physiologique, affranchi de l'approche statistique des LLM - faire émerger les régularités par une grammaire couplée à une physiologie d'organes.

---

## **Licence**

Distribué sous la licence MIT. Voir [LICENSE](LICENSE) pour les détails.

Copyright (c) 2026 Stéphane Varin. All rights reserved.

---

## **Coda : le langage comme organe, rendu exécutable**

**evoCell** constitue la preuve de concept que le langage peut être traité opérationnellement comme un organe. L'hypothèse « le langage est un organe » cesse d'être une simple métaphore descriptive pour devenir une architecture exécutable, cohérente et féconde.

Trois arguments fondent cette affirmation :

- **Une définition strictement biologique de l'organe.** La fonction y précède la structure (le contrat) : un organe n'est pas un type cellulaire, mais un engagement fonctionnel liant des tissus différenciés. Ces derniers conservent leurs domaines (les zones), à savoir des populations de cellules spécialisées partageant un territoire grammatical. L'assemblage est validé de manière émergente : la zone filtre localement, l'organe juge l'ensemble — car l'association d'un fragment déterminant parfait et d'un fragment nom parfait peut former un groupe nominal bancal. Enfin, l'indice de cohérence constitue le signe vital propre à l'organe. Chaque propriété biologique possède ici son homologue formel ; la correspondance est **opérationnelle, et non purement métaphorique**.

- **Une hiérarchie complète.** De la suite d'éléments jusqu'à l'entité globale : fragment → cellule → tissu → organe → organisme. Le langage apparaît comme le liant des niveaux intermédiaires : les cellules « connaissent » les lettres, les zones « connaissent » les domaines, les organes « connaissent » les groupes syntaxiques. La phrase n'émerge que de leur coordination, jamais d'un niveau isolé. Il s'agit d'une **dynamique causale de couplage**, et non d'une simple classification théorique.

- **Une validation empirique interne.** Lorsque l'organisme hôte (la population GE) stagne sur un plateau de *fitness*, c'est l'activité des organes — via les greffes coordonnées de GN/GV — qui relance la dynamique. Dans le monde clos du système, la thèse d'un langage fonctionnant comme une physiologie d'organes produit des effets mesurables : elle agit comme une **heuristique d'optimisation vérifiable**, dotée d'un pouvoir de résolution propre.

**Portée de l'expérimentation** — evoCell ne prétend pas trancher un débat philosophique hors de notre portée. Là où le langage-organe chomskyen désigne un module inné et individuel, l'organe est ici collectif et distribué, plus proche d'une approche énonciative du langage comme activité partagée. Son apport principal réside dans sa nature d'**hypothèse exécutable** : elle est falsifiable au sein de son propre univers — offrant une physiologie que l'on peut exécuter, observer et raffiner.

En somme, ce langage doté d'une physiologie — où les cellules portent les lettres, les tissus les domaines et les organes les groupes syntaxiques — *fonctionne* : il répare, débloque et enrichit. Cette hypothèse, transformée en code opérationnel, fonde un modèle de langage symbolique distinct des approches neuronales conventionnelles : les régularités y émergent d'une grammaire et d'une physiologie d'organes, et non de corpus statistiques.