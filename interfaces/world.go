// Copyright 2026 Stéphane Varin. All rights reserved.
// Use of this source code is governed by the MIT license.
// See the LICENSE file for details.
package interfaces

// WorldSample is a sample of the WORLD—it expresses no specific GOAL. The 
// hook releases the vocabulary corpus drop by drop (a morphogen gradient); 
// the membrane perceives it, and the nucleus refines its profile through 
// UNSUPERVISED distributional heuristics (prefixes/suffixes shared across 
// samples—with no one labeling anything as a "root"). 
// Contract DATA: IPlasticNucleus.Observe consumes it; the sealed WorldSampleMsg 
// (root) envelope carries it.
//
// The sample carries NEITHER the evaluator’s target NOR the fitness—only the 
// morphological regularity of the world. The cell learns that "this language 
// forms adjectives using '-en'," rather than "I must produce 'golden'." Intent 
// resides on the evoGo side; morphology belongs to the world.
type WorldSample struct {
	Text string     // a term from the world (e.g., "golden", "wooden", "goldfish")
	Weight float64  // frequency/support within the corpus (weights the heuristic)
}