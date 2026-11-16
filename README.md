# evoCell: A Cellular Agent Layer for Evolutionary Symbolic Engines

---

## **Description**

**evoCell** is a demonstration module (**proof of concept**) of **cellular agents** developed in **Go**, built on the **Actor model** (via [go-actor](https://github.com/vladopajic/go-actor)) and designed to hybridize with evolutionary engines such as **[evoGo](https://github.com/svarin92/evoGo)**. Unlike a grammatical evolution engine that evolves complete genomes, evoCell operates by maintaining a **gene pool**: a living population of cellular agents that cultivate, filter, and circulate *fragments* of genetic material between individuals.

In this bio-inspired architecture, the **Membrane** is not merely a physical barrier: it is the communication interface that defines how each agent interacts with its ecosystem. Each cell is an autonomous actor whose membrane performs three critical functions:

- **Excretion** (horizontal gene transfer): a high-performing cell extracts a fragment of its genome and sends it to its neighbours, propagating its success as a *beneficial infection*.
- **Absorption** (immune filtering): no incoming material is ever integrated directly; the membrane routes it to an **immune system** that renders a self/non-self verdict before any transgenesis.
- **Signalling** (quorum sensing): cells diffuse chemical signals to influence their neighbours — a **stress signal** raises the colony's mutation rate (searching for a solution *for the colony*), a **stability signal** calms aggressive mutations and stabilizes the zone.

### **Key Features**

- **Stem cells and data-driven differentiation**: every `Cell` is totipotent. By grafting a generic `Nucleus` — configured via a `ZoneSpec` read from the symbol table of the current BNF (Backus-Naur form) grammar — it specializes for a given grammatical zone. No category is fixed a priori in the module (`VowelCell`, `ConsonantCell`, or again `VerbCell` are absent from the code): changing the grammar does not require changing a single line of evoCell.
- **Grammar zoning**: the tissue topology is derived automatically from the non-terminals of the host engine's grammar — each rule becomes a zone, each production a neighbourhood contract. The tissue orchestrates the circulation and expression of genetic fragments, while the grammar determines how these exchanges are articulated.
- **Annealed porosity**: membranes are porous early in the cycle, favouring numerous exchanges between cells to maximize exploration. Late in the cycle, they become sealed, which stabilizes and exploits the best genomes. A single parameter drives both inter-cellular exchange and asymmetric division.
- **Collective robustness**: if an agent mutates into a "cancerous" form, membranes learn to block its messages — a collective immune memory intervenes, like antibodies against a pathogen.

---

## **Foundations and Inspirations**

**evoCell** draws inspiration from:

- **[go-actor](https://github.com/vladopajic/go-actor)**: a lightweight Go library for the Actor model, whose *mailboxes* and *workers* provide the concurrency backbone of the tissue.
- **Horizontal Gene Transfer** and **Quorum Sensing** in microbiology: bacteria, which exchange plasmids and synchronize their behaviour through chemical signals, directly inspire the model of excretion and signalling.
- **Stem cell biology**: asymmetric division (a stem cell producing a differentiated cell) is preferred during organogenesis over the symmetric division of mature tissues. This principle guides cell division in evoCell.
- ***The Phenomenon of Life* by [Hans Jonas](https://www.puf.com/le-phenomene-de-la-vie-0)** (PUF): Hans Jonas's organic philosophy directly irrigates evoCell's ontology. For Jonas, the organism is an exception: its true identity resides in a process rather than in a fixed substance, its form persisting at the very heart of the ceaseless exchanges of matter. evoCell borrows identity as a cycle rather than a content (the Harvest → Deposit → Recombine circuit offers a computational translation of metabolism); the membrane as an act of delimitation constitutes the self (an inside exists only by opposing a filtered outside); needful freedom and immanent teleology: just as openness to the environment engenders specialization in the living (emergent plasticity), evoCell's regularities emerge from the world's samples rather than being decreed a priori.
- **[evoGo](https://github.com/svarin92/evoGo)**: the evolutionary symbolic engine this module hybridizes with, whose immune system (`IImmune`) and dynamic organism specialization (`IOrganism`) prefigured this cellular layer.

---

## **Architecture**

### **What is a cellular system?**

A cellular system is a population of agents that cultivate, filter, and circulate *partial* genetic material (fragments), organized by domain, **never producing or evaluating a complete phenotype**. Formally, it is the quintuple **⟨Z, M, I, P, T⟩**:

1. **Z — Zoning**: spatial organization by domain (grammatical zones derived from the host grammar): *where* fragments live.
2. **M — Membranes**: communication interfaces (excretion, absorption, signalling): *how* fragments circulate.
3. **I — Immunity**: toxic-pattern memory and self/non-self verdicts: *which* fragments are refused.
4. **P — Porosity**: annealed permeability profile: *how much* circulates, depending on the phase of the cycle.
5. **T(ick) — Clock**: asynchronous *ticks* between the host engine's generations — the system lives *while* the host population waits.

The unit of the host population is the complete genome (a phenotype candidate), while that of the cellular system is the fragment. The link between these two scales is not a one-to-one association between an individual and its component, but that of a flux which, at the start of each generation, fragments the whole genomes, only to retain and graft back the best-performing elements.

### **The biological hierarchy**

The module embraces the four levels of biological organization, from the simplest to the most complex:

1. **Cell** (`cell.go`): an actor holding a codon genome, a *fitness*, and a mutation rate modulated by received signals. Totipotent by default; a graftable `Nucleus` specializes its metabolism.
2. **Tissue** (`space.go`): a grammatical zone — a population of specialized cells sharing a domain. It manages the registry of membrane addresses, the collective statistics, and diffuses quorum sensing.
3. **Organ** (`organ.go`): an assembly of zones bound by a **functional contract**. Its role is to produce, through their quorum dialogue alone, coherent genetic material for a higher-level non-terminal (e.g. the "noun phrase" organ coordinates the determiner, noun, and adjective zones; its coherence index is the tissue's top-level signal). Organs thus validate *assemblies*: even two perfect fragments may still form an incoherent whole.
4. **Organism**: the host engine's individual — the complete genome and its phenotype, which no cellular layer ever produces.

### **The supporting actors**
- **Membrane** (`membrane.go`): the cell's only exchange interface. It manages two mailboxes: one private for internal metabolism (cell ↔ its membrane), and one shared for the tissue's ecosystem (membrane ↔ tissue).
- **Immune system** (`immune.go`): the security filter. It evaluates every incoming fragment in light of the memorized toxic patterns and the tissue's average *fitness*.
- **Space** (`space.go`): the tissue's host structure. It maintains the registry of membrane addresses, aggregates collective statistics, and orchestrates the diffusion of quorum sensing.
- **Tissue coordinator** (`tissue.go`): the *endocrine interface* between the host engine and the cells. It regulates the division policy, where the probability of asymmetric differentiation is driven by the current porosity.

### **Hybridization with evoGo**

evoGo exposes a minimal, nil-safe contract (`interfaces/hybridization.go`):

```go
type HybridizationHook interface {
	OnGenerationStart(generation int, population []IIndividual)
	OnGenerationEnd(generation int, bestEver IIndividual)
}
```

A single line of code is enough to hybridize the module. The population is born inside `SearchLoop` and already owns its Genomizer: the attachment therefore happens through an option, with no re-instantiation.

```go
best, err := ge.SearchLoop(..., ge.WithHybridizationHook(evogo.NewHook(cfg, obs)))
```

Without hybridization, evoGo works like a nucleus stripped of its membrane: it never depends on evoCell. The adapter (`evogo/adapter.go`) remains the only file to ensure the bridge between the two worlds.

### **The tissue life cycle**

The coordinator and the observer generate explicit traces between two generations of evoGo. Starting from the `letter.bnf` grammar and for the target `"golden"`, the unfolding of a cycle could produce the following sequence:

**Seeding — `OnGenerationStart(1)`**: the 60 GE genomes are disassembled into fragments and sown into the zones derived from the grammar's non-terminals:

```
[tissue] seeding zones from letter.bnf: {string, letter, string_tail, vowel, consonant}
[tissue] 60 genomes disassembled -> 214 fragments (avg 3.6/genome)
[zone:vowel]        41 fragments absorbed by 23 cells
[zone:consonant]    87 fragments absorbed by 25 cells
[zone:string_tail]  86 fragments absorbed by 12 cells
```

**Metabolism — inter-generation *ticks***: cells excrete fragments (HGT), apply immune verdicts, and exchange quorum signals, under the control of the annealed porosity:

```
[tick g=12] porosity=0.79  stress_signals=0  stability_signals=3
[cell #14 zone:vowel]      excrete: codons[3:9] fitness=0.83 -> 3 neighbours (HGT)
[cell #7  zone:consonant]  absorb: from #14 -> IMMUNE: low fitness (0.31 < avg 0.57) REJECTED
[cell #22 zone:consonant]  absorb: from #14 -> IMMUNE: accepted (permeable, self-compatible)
[cell #22 zone:consonant]  transgenesis: codons[3:9] <- fragment of #14
[cell #31 zone:vowel]      SIGNAL stress (fitness 0.17 dropping) -> neighbours: mutation x2.5
[cell #4  zone:vowel]      SIGNAL stability (fitness 0.92) -> neighbours: calmed (mutation x0.5)
[immune] toxic pattern learned sig=a3f9 (sender #31 saturated 40% of mailbox quota)
```

Note that cell #31 was flooding its neighbours; its pattern is now present in the immune memory, and its messages are blocked throughout the tissue.

**Harvest — `OnGenerationEnd(12)`**: organ-coordinated fragments are grafted onto weak GE (Grammatical Evolution) individuals:

```
[organ string] coherence=0.71  zones={letter,string_tail} assemblies validated 6/8
[deposit] individual #37 weak on <string_tail> <- grafted fragment (zone:string_tail, fitness 0.78)
[deposit] individual #52 weak on <letter> <- grafted organ fragment (letters coherent)
[pool] zone vowel: 41 fragments, diversity 0.63 | consonant: 87, diversity 0.58
```

**The rescue regime** — the clearest demonstration of evoCell's value. Here, while evoGo has reached a plateau, a *hook* is grafted:

```
[generation 19] GE alone: best fitness 0.50 (plateau for 7 generations)
[generation 20] TISSUE GRAFTED — seeding from the stagnant population
[generation 21] [deposit] individual #8 <- organ fragment (HGT origin: cell #14, g=12)
[generation 23] best fitness 0.67  (stagnation broken)
[generation 27] best fitness 0.83
[generation 31] best fitness 1.00  -> "golden"
```

The *fitness* curve restarts precisely when the gene pool begins to be explored and exploited: this is the work of the *beneficial infection*.

**The Observer's outputs**: the tissue graph (`tissue.dot`) and the pool metrics (`pool.csv`) complement evoGo's `best_ever.dot`:

```
tissue.dot: clusters = grammatical zones | node colour = grafted Nucleus
            edges = HGT transfers, thickness ∝ current porosity
pool.csv:  generation, zone, fragments, diversity, average donor fitness
```

---

## **Getting Started**

```go
import (
	"github.com/svarin92/evoGo/ge"

	"github.com/svarin92/evoCell/config"
	evogo "github.com/svarin92/evoCell/evogo"
	"github.com/svarin92/evoCell/renderer"
)

// Tissue configuration.
cfg := config.Config{
	Cells:             50,
	MaxGenerations:    100,
	MutationBoost:     2.5,
	DivisionThreshold: 0.8,
	DeathThreshold:    0.2,
	Porosity: config.PorosityProfile{
		StartPermeability: 0.9,  // porous: exploration
		EndPermeability:   0.05, // sealed: exploitation
		Annealing:         config.LinearAnnealing(0.9, 0.05),
	},
}

obs := renderer.NewObserver(cfg)

// The tissue seeds itself from the population and lives between generations.
bestEver, err := ge.SearchLoop(
	100, 50, grammar, "golden",
	replacementFunc, selectionFunc, fitnessFunc,

	// Hybridization through an option: the population is born inside
	// SearchLoop, which already owns its own Genomizer. The observer
	// feeds the rendering (tissue.dot, pool.csv).
	ge.WithHybridizationHook(evogo.NewHook(cfg, obs)),
)
```
---

## **Plasticity: specialization as a heuristic**

The evolution of the architecture between the two engines reflects a fundamental paradigm shift in the way the processing of forms and symbols is conceived: **specialization is not a fixed taxonomy, but an emergent heuristic**.

In evoGo, organisms (`IOrganism`) are predetermined: the target phenotype belongs to fixed categories (alphabetical, numeric, ontological) anchored in the structure of the engine. This is a robust model, but a rigid one. In evoCell, on the contrary, the cell does not know in advance which categories it will have to handle. Non-terminal symbols are dynamic (the table of grammatical symbols is rebuilt on the fly by the `builder`), and the grafted nucleus continuously interprets the active grammar.

This mechanism induces two fundamental ruptures:

- **The membrane is where discovery happens.** The only interface with the outside world, it converts context into specialization. The host engine spots a morphological regularity among its high-performing phenotypes — such as the pattern "golden" decomposed into root and suffix. The differentiation signal then transmits its specification and re-grafts the nuclei on the fly, dispensing with any recompilation. Like a stem cell specializing under the effect of the morphogens of its environment, the system finds its equivalent here: the symbol table acts as a morphogen.
- **Filtering is a heuristic.** When a nucleus consults the collective immune memory through its absorption predicate, it operates like a B lymphocyte: the memorized antibody — a winning pattern or a toxic imprint — steers its specialization. Zones no longer filter according to arbitrary rules, but according to what experience has already put to the test. By capitalizing on its past successes, the system's plasticity becomes an engine of discovery.

In practice: a cellular population can reorient itself during a cycle toward entirely unforeseen sub-categories (such as "suffixes", "roots", or "long vowels"), as soon as the tissue observes that these classifications nourish the most vulnerable individuals. This is the operational equivalent of biological metaplasia — the process by which a differentiated cell re-specializes under the constraint of its environment. Through this mechanism, *the tissue learns its own anatomy*.

### Perception: the membrane as a sensory organ

Plasticity draws on two complementary sources: the command (the host engine identifies a regularity and orders the re-graft) and perception — the latter renewing more faithfully with the logic of the living.

For the membrane is not reduced to a mere passive interface: it acts as a genuine sensory organ. The tissue thus has a perception channel where the raw corpus (`golden`, `wooden`, `goldfish`, etc.) diffuses like a morphogen gradient. Every plastic nucleus extracts regularities from it through an unsupervised heuristic. No label is ever applied here ("this is a root, that a suffix"): it is the statistical distribution alone that decides.

The architectural boundary thus remains inviolate. The corpus transmits to the system the structures of the world — its morphological recurrences — not its final goal (the evaluator's target never travels through this channel). The cell simply learns to decode the regularities of the language, not to imitate a predetermined result.

### Knowledge as the fruit of plasticity

The heuristic thus forged is not an opaque parameter buried in the workings of the system: it becomes first-class data. Captured along the harvests — through capture mechanisms that preserve the integrity of living nuclei — the emergent profiles are aggregated per zone and then externalized.

This material proves both readable by a human and consumable by the host engine, which can thus enrich its own pattern library from the tissue's discoveries.

This is why plasticity here goes beyond a mere adaptation mechanism to establish itself as the very condition of self-learning. A system with a fixed specialization can optimize itself, but it cannot learn. evoCell brings together the three pillars of a genuine evolutionary cognition:
- Encoding: specialization fixes the regularities it has discovered.
- Explication: knowledge is readable, inspectable, and exportable.
- Feedback: filtering redefines what the zone retains, in return bending the trajectory of the individuals.

*The tissue does not merely adapt: it produces morphological knowledge.*
---

## **Roadmap**

- [ ] Coordinator actor: integrate the command mailbox and asynchronous management of `Harvest` responses.
- [ ] Division policy: apply asymmetric division governed by `P(asymmetric) = porosity`.
- [ ] Harvest and transfer mechanism: implement zone-aware harvest (`HarvestByZone`) and organ grafting (`HarvestByOrgan → Deposit`) targeting weak individuals.
- [ ] Tissue topologies: structure communication, first by *broadcast*, then by Moore grids combined with a quorum radius.
- [ ] Dynamic validation: put the actors through their paces with [go-super-actor](https://github.com/vladopajic/go-super-actor).
- [ ] Language modelling: design a symbolic, physiological language model, free from the statistical approach of LLMs — let regularities emerge from a grammar coupled with a physiology of organs.

---

## **License**

Distributed under the MIT License. See [LICENSE](LICENSE) for details.

Copyright (c) 2026 Stéphane Varin. All rights reserved.

---

## **Coda: language as an organ, made executable**

**evoCell** is the proof of concept that language can be treated operationally as an organ. The hypothesis "language is an organ" ceases to be a mere descriptive metaphor to become an executable, coherent, and fruitful architecture.

Three arguments ground this claim:

- **A strictly biological definition of the organ.** Function precedes structure (the contract): an organ is not a cell type, but a functional commitment binding differentiated tissues. These tissues keep their domains (the zones) — populations of specialized cells sharing a grammatical territory. The assembly is validated in an emergent manner: the zone filters locally, the organ judges the whole — for the association of a perfect determiner fragment and a perfect noun fragment may still form an ailing noun phrase. Finally, the coherence index is the organ's own vital sign. Every biological property has its formal homologue here; the correspondence is **operational, not purely metaphorical**.

- **A complete hierarchy.** From the sequence of elements to the global entity: fragment → cell → tissue → organ → organism. Language appears as the binding of the intermediate levels: cells "know" letters, zones "know" domains, organs "know" syntactic groups. The sentence emerges only from their coordination, never from an isolated level. This is a **causal dynamics of coupling**, not a mere theoretical classification.

- **An internal empirical validation.** When the host organism (the GE population) stagnates on a fitness plateau, it is the activity of the organs — through coordinated NP/VP grafts — that restarts the dynamics. Within the closed world of the system, the thesis of a language operating as a physiology of organs produces measurable effects: it acts as a **verifiable optimization heuristic**, endowed with its own resolving power.

**Scope of the experimentation** — evoCell does not claim to settle a philosophical debate beyond our reach. Where the Chomskyan language-organ denotes an innate, individual module, the organ here is collective and distributed, closer to an enunciative approach of language as a shared activity. Its main contribution lies in its nature as an **executable hypothesis**: it is falsifiable within its own universe — offering a physiology that can be executed, observed, and refined.

In short, this language endowed with a physiology — where cells carry the letters, tissues the domains, and organs the syntactic groups — *works*: it repairs, unblocks, and enriches. This hypothesis turned into operational code grounds a symbolic language model distinct from conventional neural approaches: regularities emerge there from a grammar and a physiology of organs, not from statistics over corpora.