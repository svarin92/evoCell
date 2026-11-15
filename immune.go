// Copyright 2026 Stéphane Varin. All rights reserved.
// Use of this source code is governed by the MIT license.
// See the LICENSE file for details.
package evocell

// immune.go, or immune filtering – every absorbed fragment passes before 
// the immune system: nothing is integrated without a verdict. It is also 
// the immune system that learns to block "cancerous" agents.

import "math/rand"

/* ImmuneSystem */

// Immune System filters all incoming material. Bio-inspired by the evoGo 
// Immune contract (self/non-self + pathogen memory).
type ImmuneSystem struct {

    // Immune memory of fragments already deemed toxic (the equivalent of 
    // antibodies) — an agent that saturates resources ends up here and will 
    // be collectively blocked.
    failedPatterns map[string]float64  // imprint → observed fitness

    // Current average tissue fitness (relative self/non-self).
    populationAverage float64

    // Local generator (testable: injectable seed).
    rng *rand.Rand
}

// LearnToxic marks an imprint as pathogenic (collective response).
func (i *ImmuneSystem) LearnToxic(f GeneFragment, fitness float64) {
    i.failedPatterns[signature(f)] = fitness
}

// Screen renders a verdict on an incoming fragment.
func (i *ImmuneSystem) Screen(
    f GeneFragment, 
    avgFitness,
    permeability float64,
) AcceptVerdict {

    // 1. Recognition: blocking fragment already cataloged as toxic.
    if w, known := i.failedPatterns[signature(f)]; known && w < avgFitness {
        return AcceptVerdict{false, "non-self: known toxic pattern"}
    }

    // 2. Poor selection: emitter rejection worse than average
    if f.Fitness < avgFitness {
        return AcceptVerdict{false, "low fitness"}
    }

    // 3. Probability of random rejection linked to current porosity.
    if i.rng.Float64() > permeability {
        return AcceptVerdict{false, "membrane sealed"}
    }

    return AcceptVerdict{true, "accepted"}
}

// UpdateAverage is called by the tissue at each generation.
func (i *ImmuneSystem) UpdateAverage(avg float64) { 
    i.populationAverage = avg 
}

// Immunity vouches for a fragment — a criterion for nucleus specialization 
// (ImmuneNucleus.CanAbsorb). It rejects known toxic imprints and ensures 
// those that align with patterns that have proven successful (translated 
// from 'evoGo' patterns via the hook). This functions like a B-lymphocyte: 
// the learned antibody becomes a discriminating factor.
func (i *ImmuneSystem) VouchesFor(f GeneFragment) bool {

    if _, toxic := i.failedPatterns[signature(f)]; toxic {
        return false  // antibodies: preemptive neutralization
    }

    return true  // default: no guarantee, but no veto either
}