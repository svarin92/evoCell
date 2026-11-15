// Copyright 2026 Stéphane Varin. All rights reserved.
// Use of this source code is governed by the MIT license.
// See the LICENSE file for details.
package evocell

// organ.go — the organ: emerging cooperation between zones.

// An organ is a functional contract: a collection of zones (tissues) 
// bound by a common commitment—to produce coherent genetic material 
// for a higher-order non-terminal entity. It is not an actor, but 
// rather a lightweight object that makes use of the cells within 
// those zones.
type Organ struct {

	// nonTerminal: The contract for the component—always a non-terminal that 
	// serves as the parent in the current BNF (e.g., "noun_phrase"). This
	// syntactic anchor ensures that the extracted data is placed exactly at 
	// the syntactic location expected by the ProductionHistory.
	nonTerminal string  // the contract (e.g. “GN”)

	// zones: member zones (child non-terminals of the parent).
	zones []string  // member zones (e.g. determiner, noun, adjective)

	// coherence: an indicator of the quality of the current assembly,
	// driven by inter-zone signals (the "nerves" of the organ).
	coherence float64

	// tissue: the host tissue — the organ borrows its channels (collection,
	// assignment), but never its actors.
	tissue *Tissue
}

var _ IOrgan = (*Organ)(nil)

// Coherence returns the current quality index ([0, 1]).
func (o *Organ) Coherence() float64 { return o.coherence }

// NonTerminal defines the contract for the component (parent BNF anchor).
func (o *Organ) NonTerminal() string { return o.nonTerminal }

// Mutually absorbable: each fragment passes the kernel’s `Can Absorb` check
// against *every* other zone — the agreement is bilateral.
func (o *Organ) mutuallyAbsorbable(fragments map[string]GeneFragment) bool {

	// Implementation v0.1.0: default (full assembly validation = v0.2.0, 
	// depends on profiled kernels).
	return true
}

// Validate verifies that a set of fragments—one per member zone — correctly 
// forms a valid subtree for its corresponding non-terminal. This is emergent 
// validation: each zone performs local filtering, while the organ evaluates 
// the assembly (e.g., "la chevaux" [the horses—using the wrong grammatical 
// gender] → rejection).
func (o *Organ) Validate(fragments map[string]GeneFragment) bool { 
	
	// 1. Completeness: one contribution per member area.
	if len(fragments) != len(o.zones) {
		return false
	}

	for _, zone := range o.zones {
		frag, ok := fragments[zone]

		if !ok || len(frag.Codons) == 0 {
			return false
		}

	}

	// 2. Coherence: the fragments must be mutuallyabsorbable — the fragment 
	// from the "noun" zone must pass the `CanAbsorb` check of the "adjective" 
	// zone's core (agreement in gender, number, etc.). The precise resolution 
	// is delegated to the grafted cores.
	return o.mutuallyAbsorbable(fragments)
}

// UpdateCoherence recomposes the index from recent signals from limb areas: 
// the “nerve” of the organ. Called by the tissue to each Tick — the organ 
// ages with its zones.
func (o *Organ) UpdateCoherence(perZone map[string]float64) {

	if len(o.zones) == 0 {
		return
	}

	sum := 0.0

	for _, zone := range o.zones {

		if c, ok := perZone[zone]; ok {
			sum += c
		}

	}

	o.coherence = sum / float64(len(o.zones))
}