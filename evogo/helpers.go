// Copyright 2026 Stéphane Varin. All rights reserved.
// Use of this source code is governed by the MIT license.
// See the LICENSE file for details.
package evogo

import evocell "github.com/svarin92/evoCell"

// evogo/helpers.go — translation between the two worlds.

// carrierOf returns an individual's GenomeCarrier — the evoGo contract
// (v1.2.0) by which a host genome accepts a codon deposit. Nil if the
// individual does not expose the contract: the graft is simply ignored
// (carrier OPTIONAL, see contract semantics).
func carrierOf(ind IIndividual) IGenomeCarrier {
	carrier, _ := ind.(IGenomeCarrier)
	return carrier
}

// CodonsOf extracts the slice of codons from an evoGo individual. The
// concrete type of IIndividual determines the access method; the adapter
// makes no assumptions about the internal structure of the genome.
func CodonsOf(ind IIndividual) []int {

	// Implementation committed to integration: based on the actual
	// IIndividual interface (e.g., ind.GetGenome()), including
	// defensive copying — the tissue must never mutate the host genome.
	return append([]int(nil), ind.GetGenome()...)
}

// GenomesOf extracts the genomes (codon segments) from a population of evoGo 
// individuals. It is used to seed the tissue (OnGenerationStart, gen 0): each 
// complete genome will be broken down into fragments by the zones.
func GenomesOf(pop []IIndividual) [][]int {
	out := make([][]int, 0, len(pop))

	for _, ind := range pop {
		out = append(out, CodonsOf(ind))
	}

	return out
}

// recombination resolves the graft for a weak individual: the best 
// fragment from its zone is anchored into its genome as a graft, 
// thereby reconstituting the genome. The fragment overwrites only its
// corresponding segment — the rest of the host genome remains intact.
func recombine(ind IIndividual, fragments []GeneFragment) (genome []int, anchor int) {
	genome = CodonsOf(ind)
	best := evocell.BestFragmentOf(fragments)
	
	if len(best.Codons) == 0 || len(best.Codons) > len(genome) {
		return genome, -1  // Nothing graftable: intact genome, invalid anchor
	}

	// Non-terminal anchor: position resolved via ProductionHistory
	// (evoGo contract); anchor at the top pending.
	anchor = 0
	copy(genome[anchor:], best.Codons)

	return genome, anchor
}

// --- The transplant circuit (Harvest → Deposit) ---

// A weak GE individual and the non-terminal on which its derivation 
// is limping — the “symptom” that the transplant must address.
type Weakness struct {
	Ind IIndividual
	On  string  // non-terminal host of the weakness zone
}

// weakOn returns individuals below the average fitness of the population, 
// each labeled by their weakest non-terminal. Reading the derivation history 
// is delegated to the ProductionHistory of the Genomizer (evoGo contract): 
// it's the one that knows which rule has derived — the adapter cannot.
func weakOn(population []IIndividual) []Weakness {
	avg := 0.0

	for _, ind := range population {
		avg += ind.GetFitness()
	}

	avg /= float64(len(population))
	weak := make([]Weakness, 0)

	for _, ind := range population {

		if ind.GetFitness() < avg {

			// Implementation resolved at integration: reading the
			// ProductionHistory (evoGo contract) → non-terminal the
			// least well derived from this individual.
			if nt, ok := weakestNonTerminalOf(ind); ok {
				weak = append(weak, Weakness{Ind: ind, On: nt})
			}

		}

	}

	return weak
}

// weakestNonTerminalOf returns the weakest non-terminal of the derivation 
// of an individual (delegated to the ProductionHistory of evoGo).
func weakestNonTerminalOf(ind IIndividual) (string, bool) {

	// Resolution at integration: ind.ProductionHistory().WeakestNT()
	return "", false
}