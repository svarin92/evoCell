// Copyright 2026 Stéphane Varin. All rights reserved.
// Use of this source code is governed by the MIT license.
// See the LICENSE file for details.
package evogo

import (
	evocell "github.com/svarin92/evoCell"
	"github.com/vladopajic/go-actor/actor"
)

// Hook adapts a Tissue to the HybridizationHook contract from evoGo.
var _ IHybridizationHook = (*Hook)(nil)

/* LazyHook */

// LazyHook is a nil-safe decorator: it ignores calls until the activation 
// generation is reached ("rescue" mode — the mechanism kicks in when evoGo 
// stalls). 100% on the evocell side: evoGo remains unmodified; its contract 
// still consists of two methods.
type LazyHook struct {
	activateAt int
	inner      IHybridizationHook  // the real hook (tissue)
}

func (h *LazyHook) OnGenerationStart(gen int, pop []IIndividual) {

	if gen >= h.activateAt {
		h.inner.OnGenerationStart(gen, pop)
	}

}

func (h *LazyHook) OnGenerationEnd(gen int, pop []IIndividual) {

	if gen >= h.activateAt {
		h.inner.OnGenerationEnd(gen, pop)
	}

}

/* Hook */

type Hook struct {
	cfg    Config
	obs    Observer 
	tissue Tissue

	// --- Morphological discovery ---

	// The hook acts as a translator, converting evoGo patterns into 
	// differentiation commands.
	patterns IMorphologicalPatterns  // discovery engine (controller/morphology.go, v1.3.0)
	
	// The phenotypes of the latest generation are populated at each 
	// OnGenerationEnd, providing discovery elements to Discover, 
	// before being consumed by pushDifferentiation at the following 
	// OnGenerationStart.
	recentPhenotypes []string

	// evoGo serializer: the ONLY way to access the terminals of a
	// non-terminal (GetTerminals) — the hook does not re-derive 
	// anything, nor does it re-type the symbol table. The Genomizer 
	// maintains its own table (map[string]IRuleModel) for its own use; 
	// the hook does not duplicate it, but instead queries the evoGo 
	// component responsible for serializing the rules.
	serializer ISerializer
}

// OnGenerationStart synchronizes the tissue with the GE population.
func (h *Hook) OnGenerationStart(gen int, pop []IIndividual) {

	if h.tissue == nil {
		
		// First generation: seed the tissue from the GE genomes.
		h.tissue = evocell.NewTissue(h.cfg, GenomesOf(pop))
		h.tissue.Start()
		return
	
	}

	h.tissue.Tick(gen)
}

// OnGenerationEnd harvests the best fragments PER ZONE and reinject them into 
// the GE population via GenomeCarrier.Deposit — Symptom-based medicine: each 
// weak individual receives fragments from THEIR area of ​​weakness, not from 
// all zonss.
func (h *Hook) OnGenerationEnd(gen int, pop []IIndividual) {

	if h.tissue == nil {
		return  // the tissue has never been sown (hook nil-safe)
	}

	

	// 1. HarvestByZone: the best fragments, grouped by non-terminal.
	perZone := h.tissue.HarvestByZone()  // map[nonTerminal][]GeneFragment
	                                     // keys = non-terminals from the BNF

	// 2. Deposit: graft onto weak individuals, targeted by zone.
	grafted := 0

	for _, w := range weakOn(pop) {
		fragments, ok := perZone[w.On]

		if !ok || len(fragments) == 0 {
			continue  // no therapy available for this non-terminal
		}

		if carrier := carrierOf(w.Ind); carrier != nil {
			codons, anchor := recombine(w.Ind, fragments)

			if anchor >= 0 {

				// Atomic deposit (evoGo contract v1.2.0): either the entire
				// tranche is added, or nothing changes.
				_ = carrier.Deposit(anchor, codons)
				grafted++
			}

		}

	}

	// 3. Traces: the Observer records the graft (pool.csv, tissue.dot).
	h.obs.RecordGenerationEnd(gen, grafted, perZone)
}

// pushDifferentiation requires referring to the WeakestNT section — 
// it involves translating the discovered patterns into differentiation 
// commands.
func (h *Hook) pushDifferentiation(gen int) {

	for _, pattern := range h.patterns.Discover(h.recentPhenotypes, 0.3) {

		if pattern.Support < 0.3 {
			continue  // pattern too rare: no general re-grafting
		}

		h.tissue.Broadcast(DifferentiationSignal{
			Spec: evocell.ZoneSpec{
				NonTerminal: pattern.NonTerminal,  // BNF ANCHOR (symbol table)
				Profile: pattern.Profile,  // refinement of the filtering

				// Terminals: read from the table for this non-terminal — this 
				// constitutes the basis for filtering (filterOnTerminals).
				Terminals: h.terminalsOf(pattern.NonTerminal),
			},
		})
	}

}

// StreamWorld diffuses the world’s corpus into the fabric, drop by drop 
// — the morphogen gradient. The corpus is loaded by the HOST (main.go or 
// evoGo): evocell reads no files; the world communicates solely through 
// messages. Each sample is weighted by its frequency within the corpus; 
// plastic kernels deduce regularities via UNSUPERVISED heuristics (e.g., 
// recurring "-en" or "-ish" → a "suffix" profile). The corpus constitutes 
// the WORLD (vocabulary), never the GOAL (the evaluator’s target never 
// passes through this channel).
func (h *Hook) StreamWorld(corpus []string, freq map[string]float64) {
	ctx := actor.ContextStarted()

	for _, word := range corpus {
		w := 1.0

		if f, ok := freq[word]; ok {
			w = f
		}

		h.tissue.Broadcast(evocell.WorldSampleMsg{
			WorldSample: evocell.WorldSample{Text: word, Weight: w},
		})
	}

	_ = ctx
}

// terminalsOf: the terminals of a non-terminal, read from the evoGo symbol 
// table (the builder reconstructs it for each BNF — the hook reads the 
// CURRENT table, never a static, fixed list).
func (h *Hook) terminalsOf(nonTerminal string) []string {
	set := h.serializer.GetTerminals() 
	out := make([]string, 0, len(set)) 
	
	for t := range set { 
		out = append(out, t) 
	} 
	
	return out
}