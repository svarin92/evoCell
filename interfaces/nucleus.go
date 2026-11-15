// Copyright 2026 Stéphane Varin. All rights reserved.
// Use of this source code is governed by the MIT license.
// See the LICENSE file for details.
package interfaces

// interfaces/nucleus.go — cell nucleus contracts.

/* INucleus */

// INucleus is the specialized nucleus: it interprets the genome according
// to its specific region. It acts as a cellular mirror for the IOrganism
// interfaces within evoGo—specifically regarding the MECHANISM: while an 
// evoGo organism is constructed by the engine for a given phenotype, the 
// cell nucleus is GRAFTABLE, configured by data from the current grammar.
type INucleus interface {

	// CanAbsorb determines whether a foreign fragment is relevant to this
	// core (specialization acts as an additional filtering layer, upstream of
	// the immune system). The predicate is CONSTRUCTED from the
	// terminals of the specialization. It does not have a hard-coded category.
	CanAbsorb(f GeneFragment) bool

	// It extracts from the genome the portion relevant to its zone.
	// E.g., a "vowel" nucleus retains only those codons that produce
	// vowel terminals.
	Express(codons []int) []int

	// MutateRate returns the mutation rate specific to the zone.
	MutateRate() float64

	// Zone returns the complete specialization of the area served (anchored
	// non-terminal in BNF + terminals + profile). This serves as the
	// identity of the nucleus — it is a piece of data, never an enumerated 
	// type.
	Zone() ZoneSpec
}

/* NucleusGrammatical */

// INucleusGrammatical adds the dimension of a "zone" to the graft:
// each cell knows which non-terminal it serves. This is the contract
// that carries the BNF ANCHOR — the immutable half of plasticity.
type INucleusGrammatical interface {
	INucleus

	// NonTerminal identifies the grammatical zone (e.g., "NP" or 
	// "noun_phrase"). It always corresponds to an actual non-terminal 
	// in the CURRENT BNF — never to an invented category (such as "root" 
	// or "suffix"). The addressing for the grafting process (HarvestByZone 
	// → Deposit) depends on this.
	NonTerminal() string
}

/* IPlasticNucleus */

// IPlasticNucleus: the nucleus that learns from PERCEPTION. Learning is 
// LOCAL: the nucleus refines the filtering within ITS own zone based on 
// the morphological statistics it observes — it never reconstructs an 
// intention, nor any phenotype either.
type IPlasticNucleus interface {

	// LearnedProfile: the current profile, along with its conjecture 
	// — the knowledge acquired, which is both inspectable and exportable 
	// (cf. Observer).
	LearnedProfile() (profile string, support float64)

	// Observe - phase 1, ingests a sample of the world (via the membrane 
	// interface) and updates its zone profile if it detects an emerging
	// regularity (unsupervised distributional heuristic).
	Observe(sample WorldSample)

	// Consolidate - phase 2, proposed after each 'Observe' (consolidation
	// is a consequence of perception, not a tissue-level event—in the sense 
	// that the core proposes as soon as thresholds are crossed, without the 
	// 'Tick' needing to be aware of it from the inside). If a pattern exceeds 
	// the thresholds, it becomes the profile for that zone.
	Consolidate()
}