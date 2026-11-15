// Copyright 2026 Stéphane Varin. All rights reserved.
// Use of this source code is governed by the MIT license.
// See the LICENSE file for details.
package evocell

// nucleus_grammatical.go — the generic implementation of the nucleus.

/* ImmuneNucleus */

// ImmuneNucleus: a core that delegates its filtering process to an "immune 
// memory." Directly bio-inspired: a B-lymphocyte only retains an antigen 
// recognized by its own antibodies—in this context, a fragment is only 
// absorbed if it resembles a pattern that has "succeeded" (a library of 
// patterns) and does not resemble any toxic pattern (failedPatterns).
type ImmuneNucleus struct {
    delegate *NucleusGrammatical
    immune *ImmuneSystem
}

func (n *ImmuneNucleus) Zone() ZoneSpec { return n.delegate.Zone() }
func (n *ImmuneNucleus) Express(c []int) []int { return n.delegate.Express(c) }
func (n *ImmuneNucleus) NonTerminal() string { return n.delegate.NonTerminal() }
func (n *ImmuneNucleus) MutateRate() float64 { return n.delegate.MutateRate() }

func (n *ImmuneNucleus) CanAbsorb(f GeneFragment) bool {

    // 1. The delegate filters by non-terminal and profile.
    if !n.delegate.CanAbsorb(f) {
        return false
    }

    // 2. Immunity makes the call: the rejection of a known toxic motif is 
    // anticipated (the antibody neutralizes BEFORE infection). Conversely, 
    // a motif that has succeeded facilitates uptake (success begets success).
    return n.immune.VouchesFor(f)
}

/* NucleusGrammatical */

type NucleusGrammatical struct {
	spec ZoneSpec  // data from the current grammar

	// Profile predicate, constructed from the specialization endpoints 
	// (filtering is a derived piece of data, not a typed behavior). 
	// Nil = basic filtering.
	profileFilter func(GeneFragment) bool

	// --- Learning Core (IPlasticNucleus) ---
	
	// Stores form (morphological) statistics derived from the perception 
	// of the environment.  Raw, internal data remains private (variables 
	// in lowercase). Only consolidated knowledge is shared externally via 
	// the profile (LearnedProfile).
	samples int  // samples ingested
	totalWeight float64  // observed cumulative weight
	motifFreq map[string]float64  // motif « seq@position » → fréquence pondérée
}

// CanAbsorb derives the predicate from the TERMINALS of the specialization —
// a fragment is relevant if it produces at least one terminal that the zone 
// can produce. No categories are hard-coded: with letter.bnf, a "vowel" cell 
// retains only material that derives into a, e, i, o, u — because the TABLE 
// says so, and not because an AlphabetNucleus type enumerates them.
func (n *NucleusGrammatical) CanAbsorb(f GeneFragment) bool {

	if n.profileFilter != nil {
		return n.profileFilter(f)
	}

	return true  // basic filtering: relevant to the non-terminal
}

func (n *NucleusGrammatical) Express(codons []int) []int { return codons }

// MutateRate defines a moderate mutation of the grammatical domain.
func (n *NucleusGrammatical) MutateRate() float64 { return 0.02 }

func (n *NucleusGrammatical) NonTerminal() string { return n.spec.NonTerminal }

func (n *NucleusGrammatical) Zone() ZoneSpec { return n.spec }

// --- Perception: IPlasticNucleus (Local, Unsupervised Learning) ---
//
// Role: The 'Observe' method ingests a sample from the environment and
// updates the heuristic for its local zone.
//
// Bio-inspiration: Much like B-lymphocytes that refine their antibodies
// through clonal selection upon exposure to antigens, the nucleus adjusts 
// its profile through iterative exposure to morphological regularities.
//
// Unsupervised Mechanism: No external labeling is required; the statistical
// distribution guides the emergence of structures. Recurring terminal n-grams
// become candidate "suffixes," while frequent prefixes become "roots."

// Sufficient heuristic thresholds, without which profiles cannot emerge. 
// A rare regularity does not constitute knowledge.
const (
	MinSamplesForProfile = 30    // minimum samples before consolidation
	MinProfileSupport    = 0.15  // minimum relative support of the dominant n-gram
)

// --- EPISTEMOLOGY: The heuristic does NOT KNOW what it is looking for.
// Prohibit a priori categories ("suffix", "root") during the observation 
// phase: searching for final n-grams is already asserting that the 
// morphological world is made of suffixes — a category established BEFORE 
// observation.
//
// The discipline is twofold:
//  1. IDENTIFY: accumulate n-grams across ALL positions (initial, medial, 
//     final), weighted.
//  2. PROPOSE: a regularity = (n-gram, dominant position) above thresholds. 
//     The category emerges from this (final → final regularity; initial → 
//     initial...); it is never posited a priori.
//
// Immunological analogy:
// The immune system recognizes patterns without knowing the concept of 
// "virus"; we label its discoveries after the fact. The linguistic name 
// ("suffix") is a gloss provided by the host or the human, never by the 
// tissue. The consolidated profile is the PATTERN ("en@final"), not a 
// label.

// Observe - phase 1: accumulate material (no decision here).
func (n *NucleusGrammatical) Observe(sample WorldSample) {

	if sample.Text == "" {
		return
	}

	n.samples++
	n.totalWeight += sample.Weight

	// All n-grams of length 2, at ALL their positions: the position is part 
	// of the observation, not the hypothesis.
	for _, m := range allNGrams(sample.Text, 2) {
		n.motifFreq[m] += sample.Weight  // key « seq@position »
	}

}

// Consolidate - phase 2: proposes a category. It is called after each
// observation: if a pattern exceeds the thresholds, it becomes the profile 
// for the zone. The proposed category is DERIVED from the dominant position 
// of the observed pattern.
func (n *NucleusGrammatical) Consolidate() {

	if n.samples < MinSamplesForProfile {
		return  // insufficient material: do not speculate based on nothing
	}

	best, freq := argMax(n.motifFreq)

	if n.totalWeight == 0 || freq/n.totalWeight < MinProfileSupport {
		return  // regularity too rare: this does not constitute knowledge
	}

	// The profile is the MOTIF (form + observed position), not a name.
	n.spec.Profile = best  // e.g., "en@finale"
	n.profileFilter = filterOnMotif(best)
}

// LearnedProfile: acquired knowledge — a profile and a conjecture.
func (n *NucleusGrammatical) LearnedProfile() (string, float64) {
	return n.spec.Profile, n.support()
}

// support returns the weighted relative frequency within the observed 
// corpus.
func (n *NucleusGrammatical) support() float64 {

	if n.totalWeight == 0 {
		return 0
	}

	_, w := argMax(n.motifFreq)
	return w / n.totalWeight
}