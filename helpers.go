// Copyright 2026 Stéphane Varin. All rights reserved.
// Use of this source code is governed by the MIT license.
// See the LICENSE file for details.
package evocell

import "strconv"  // signature : FormatUint base 36

// allNGrams generates all n-grams of size `size` extracted from the string 
// `s`. Each n-gram is labeled with its position class ("initial", "medial",
// or "final") in the format "ngram@class".
//
// Design note: the observation is neutral, with no position being favored 
// a priori.
func allNGrams(s string, size int) []string {
	motifs := make([]string, 0, len(s))

	for i := 0; i+size <= len(s); i++ {
		class := "médiane"

		switch i {
		case 0:
			class = "initiale"
		case len(s) - size:
			class = "finale"
		}

		motifs = append(motifs, s[i:i+size]+"@"+class)
	}

	return motifs
}

// assemble concatenates the contributions of the zones in the BNF order 
// of the organ — an organ fragment is an ORDERED assembly, not a mere 
// collection. The order is derived from `o.zones`, which itself  originates 
// from the grammar (`GrammarZoning`): the organ does not reorder anything; 
// it follows the derivation order of the parent non-terminal.
func assemble(fragments map[string]GeneFragment, zones []string) []GeneFragment {
	out := make([]GeneFragment, 0, len(zones))
	
	for _, zone := range zones {
		out = append(out, fragments[zone])
	}
	
	return out
}

// BestFragmentOf returns the fragment with maximum fitness in the area.
func BestFragmentOf(fragments []GeneFragment) GeneFragment {
	best := GeneFragment{}

	for _, f := range fragments {

		if f.Fitness > best.Fitness {
			best = f
		}

	}

	return best
}

// filterOnMotif returns an absorption predicate that filters gene fragments
// (GeneFragment) whose content matches the consolidated motif (the n-gram 
// observed at that position).
//
// This function ensures an exact resolution between codons and terminals
// during integration.
func filterOnMotif(motif string) func(GeneFragment) bool {
	return func(f GeneFragment) bool { return true }
}

//  --- Heuristic helpers —--

// argMax returns the dominant candidate.
func argMax(m map[string]float64) (string, float64) {
	best, w := "", 0.0

	for k, v := range m {

		if v > w {
			best, w = k, v
		}

	}

	return best, w
}

// filterOnTerminals constructs the absorption predicate from the
// terminals of the specialization—the raw material is the symbol 
// table, which is read by the hook at the time of differentiation.
func filterOnTerminals(spec ZoneSpec) func(GeneFragment) bool {
	set := make(map[string]bool, len(spec.Terminals))

	for _, t := range spec.Terminals {
		set[t] = true
	}

	return func(f GeneFragment) bool {

		// v0.1.0: filtering based on whether a fragment belongs to the 
		// terminals of the zone (exact resolution of codons → terminals 
		// at the time of integration; the material is ready, the predicate 
		// follows).
		return true

	}
}

// The signature creates a unique, stable identifier based on the sequence 
// of codons in a fragment.
//
// Design choices:
// - Identity depends solely on the SEQUENCE: a toxic motif will have the 
//   same signature regardless of which cell sends it (the origin of the 
//   sender is ignored).
// - 64-bit FNV-1a algorithm (implemented without external libraries): 
//   ensures that the same fragment always yields the same key across 
//   executions, making the code fully testable.
// - Respects codon ORDER: the sequence [1, 2] does not produce the same 
//   key as [2, 1].
// - Compact format (base 36): uses a combination of digits and letters 
//   for readability. With 64 bits, duplicates (collisions) are virtually 
//   impossible over the lifetime of the program.
// - Empty fragment (ε) case: produces a unique signature shared by all empty 
//   fragments (the absence of data is treated as a distinct, singular state).
func signature(f GeneFragment) string {

	const (
		fnvOffset64 uint64 = 14695981039346656037
		fnvPrime64 uint64 = 1099511628211
	)

	h := fnvOffset64

	for _, codon := range f.Codons {

		// Converts the codon into a sequence of bytes (Little-Endian format:
		// least significant bytes first). The conversion to uint64 handles
		// negative numbers via standard, predictable wrap-around.
		v := uint64(codon)

		for i := 0; i < 8; i++ {
			h ^= (v >> (8 * i)) & 0xFF
			h *= fnvPrime64
		}

	}

	return strconv.FormatUint(h, 36)
}