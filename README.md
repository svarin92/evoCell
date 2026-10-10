# evoCell: A Cellular Agent Layer for Evolutionary Symbolic Engines

---

## **Description**

**evoCell** is a demonstration module (**proof of concept**) of **cellular agents** developed in **Go**, built on the **Actor model** (via [go-actor](https://github.com/vladopajic/go-actor)) and designed to hybridize with evolutionary engines such as **[evoGo](https://github.com/svarin92/evoGo)**. Unlike a grammatical evolution engine that evolves complete genomes, evoCell operates by maintaining a **gene pool**: a living population of cellular agents that cultivate, filter, and circulate *fragments* of genetic material between individuals.

In this bio-inspired architecture, the **Membrane** is not merely a physical barrier: it is the communication interface that defines how each agent interacts with its ecosystem. Each cell is an autonomous actor whose membrane performs three critical functions:

- **Excretion** (horizontal gene transfer): a high-performing cell extracts a fragment of its genome and sends it to its neighbours, propagating its success as a *beneficial infection*.
- **Absorption** (immune filtering): no incoming material is ever integrated directly; the membrane routes it to an **immune system** that renders a self/non-self verdict before any transgenesis.
- **Signalling** (quorum sensing): cells diffuse chemical signals to influence their neighbours — a **stress signal** raises the colony's mutation rate (searching for a solution *for the colony*), a **stability signal** calms aggressive mutations and stabilizes the zone.

### **Key Features**

- **Stem cells and differentiation**: every `Cell` is totipotent; grafting a specialized `Nucleus` turns it into a `VowelCell`, `ConsonantCell`, `VerbCell`, `DigitCell`, `OntologyCell`... mirroring evoGo's dynamic `IOrganism` specialization mechanism.
- **Grammar zoning**: the tissue topology is derived automatically from the non-terminals of the host engine's grammar — each rule becomes a zone, each production a neighbourhood contract. The tissue orchestrates the circulation and expression of genetic fragments, while the grammar determines how these exchanges are articulated.
- **Annealed porosity**: membranes are porous early in the cycle, favouring numerous exchanges between cells to maximize exploration. Late in the cycle, they become sealed, which stabilizes and exploits the best genomes. A single parameter drives both inter-cellular exchange and asymmetric division.
- **Collective robustness**: if an agent mutates into a "cancerous" form, membranes learn to block its messages — a collective immune memory intervenes, like antibodies against a pathogen.

---

## **Foundations and inspirations**

**evoCell** draws inspiration from:

- **[go-actor](https://github.com/vladopajic/go-actor)**: a lightweight Go library for the Actor model, whose *mailboxes* and *workers* provide the concurrency backbone of the tissue.
- **Horizontal Gene Transfer** and **Quorum Sensing** in microbiology: bacteria, which exchange plasmids and synchronize their behaviour through chemical signals, directly inspire the model of excretion and signalling.
- **Stem cell biology**: asymmetric division (a stem cell producing a differentiated cell) is preferred during organogenesis over the symmetric division of mature tissues. This principle guides cell division in evoCell.
- **[evoGo](https://github.com/svarin92/evoGo)**: the evolutionary symbolic engine this module hybridizes with, whose immune system (`IImmune`) and dynamic organism specialization (`IOrganism`) prefigured this cellular layer.

---

## **Architecture**

### **What is a cellular system?**

A cellular system is a population of agents that cultivate, filter, and circulate *partial* genetic material (fragments), organized by domain, **never producing or evaluating a complete phenotype**. Formally, it is the quintuple **⟨Z, M, I, P, T⟩**:

1. **Z — Zoning**: spatial organization by domain (grammatical zones derived from the host grammar): *where* fragments live.
2. **M — Membranes**: communication interfaces (excretion, absorption, signalling): *how* fragments circulate.
3. **I — Immunity**: toxic-pattern memory and self/non-self verdicts: *which* fragments are refused.
4. **P — Porosity**: annealed permeability profile: *how much* circulates, depending on the phase of the cycle.
5. **T — Clock**: asynchronous *ticks* between the host engine's generations — the system lives *while* the host population waits.

The unit of the host population is the complete genome (a phenotype candidate), while that of the cellular system is the fragment. The link between these two scales is not a one-to-one association between an individual and its component, but that of a flow which, at the start of each generation, fragments the whole genomes, only to retain and graft back the best-performing elements.

### **The biological hierarchy**

The module embraces the four levels of biological organization, from the simplest to the most complex:

1. **Cell** (`cell.go`): an actor holding a codon genome, a *fitness*, and a mutation rate modulated by received signals. Totipotent by default; a graftable `Nucleus` specializes its metabolism.
2. **Tissue** (`space.go`): a grammatical zone — a population of specialized cells sharing a domain. It manages the registry of membrane PIDs, the collective statistics, and diffuses quorum sensing.
3. **Organ** (`organ.go`): an assembly of zones bound by a **functional contract**. Its role is to produce, through their quorum dialogue alone, coherent genetic material for a higher-level non-terminal (e.g. the "noun phrase" organ coordinates the determiner, noun, and adjective zones; its coherence index is the tissue's top-level signal). Organs thus validate *assemblies*: even two perfect fragments may still form an ailing whole.
4. **Organism**: the host engine's individual — the complete genome and its phenotype, which no cellular layer ever produces.

### **The supporting actors**
- **Membrane** (`membrane.go`): the cell's only exchange interface. It manages two mailboxes: one private for internal metabolism (cell ↔ its membrane), and one shared for the tissue's ecosystem (membrane ↔ tissue).
- **Immune system** (`immune.go`): the security filter. It evaluates every incoming fragment in light of the memorized toxic patterns and the tissue's average *fitness*.
- **Space** (`space.go`): the tissue's host structure. It maintains the registry of membrane PIDs, aggregates collective statistics, and orchestrates the diffusion of quorum sensing.
- **Tissue coordinator** (`tissue.go`): the *endocrine interface* between the host engine and the cells. It regulates the division policy, where the probability of asymmetric differentiation is driven by the current porosity.

### **Hybridization with evoGo**

evoGo exposes a minimal, nil-safe contract (`interfaces/hybridization.go`):

```go
type HybridizationHook interface {
	OnGenerationStart(generation int, population []IIndividual)
	OnGenerationEnd(generation int, bestEver IIndividual)
}
```

A single line of code is enough to hybridize the module. The population is created inside SearchLoop and already owns its Genomizer: the attachment therefore happens through an option, with no reinstantiation.

```go
best, err := ge.SearchLoop(..., ge.WithHybridizationHook(evogo.NewHook(cfg, obs)))
```

Without hybridization, evoGo works like a nucleus stripped of its membrane: it never depends on evoCell. The adapter (`evogo/adapter.go`) remains the only file bridging the two worlds.

### **The tissue life cycle**

The coordinator and the observer generate explicit traces between two generations of evoGo. Starting from the `letter.bnf` grammar and for the target `"golden"`, the unfolding of a cycle can produce the following sequence:

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
	...
	"github.com/svarin92/evoGo/ge"

	"github.com/svarin92/evoCell/config"
	"github.com/svarin92/evoCell/evogo"
	"github.com/svarin92/evoCell/renderer"
)

cfg := evocell.Config{
	Cells:             50,
	MaxGenerations:    100,
	MutationBoost:     2.5,
	DivisionThreshold: 0.8,
	DeathThreshold:    0.2,
	Porosity: evocell.PorosityProfile{
		StartPermeability: 0.9,  // porous: exploration
		EndPermeability:   0.05, // sealed: exploitation
		Annealing:         evocell.LinearAnnealing(0.9, 0.05),
	},
}

obs := renderer.NewObserver(cfg)

// The tissue seeds itself from the population and lives between generations.
bestEver, err := ge.SearchLoop(
	100, 50, grammar, "golden",       
	replacementFunc, selectionFunc, fitnessFunc,    
	
	// Hybridization through an option: the population is created inside 
	// SearchLoop and already owns Genomizer. L'observateur feeds the 
	// rendering (tissue.dot, pool.csv).
	ge.WithHybridizationHook(evogo.NewHook(cfg, obs)),          
)
```
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

- **A strictly biological definition of the organ.** Function precedes structure (the contract): an organ is not a cell type, but a functional commitment binding differentiated tissues. These tissues keep their domains (the zones) — populations of specialized cells sharing a grammatical territory. The assembly is validated in an emergent manner: the zone filters locally, the organ judges the whole — for the association of a perfect determiner fragment and a perfect noun fragment may form an ailing noun phrase. Finally, the coherence index is the organ's own vital sign. Every biological property has its formal homologue here; the correspondence is **operational, not purely metaphorical**.

- **A complete hierarchy.** From the sequence of elements to the global entity: fragment → cell → tissue → organ → organism. Language appears as the binding of the intermediate levels: cells "know" letters, zones "know" domains, organs "know" syntactic groups. The sentence emerges only from their coordination, never from an isolated level. This is a **causal dynamics of coupling**, not a mere theoretical classification.

- **An internal empirical validation.** When the host organism (the GE population) stagnates on a fitness plateau, it is the activity of the organs — through coordinated GN/GV grafts — that restarts the dynamics. Within the closed world of the system, the thesis of a language operating as a physiology of organs produces measurable effects: it acts as a **verifiable optimization heuristic**, endowed with its own resolving power.

**Scope of the experimentation** — evoCell does not claim to settle a philosophical debate beyond our reach. Where the Chomskyan language-organ denotes an innate, individual module, the organ here is collective and distributed, closer to an enunciative approach of language as a shared activity. Its main contribution lies in its nature as an **executable hypothesis**: it is falsifiable within its own universe — offering a physiology that can be executed, observed, and refined.

In short, this language endowed with a physiology — where cells carry the letters, tissues the domains, and organs the syntactic groups — *works*: it repairs, unblocks, and enriches. This hypothesis turned into operational code grounds a symbolic language model distinct from conventional neural approaches: regularities emerge there from a grammar and a physiology of organs, rather than from statistics based on corpora.