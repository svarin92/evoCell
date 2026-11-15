// Copyright 2026 Stéphane Varin. All rights reserved.
// Use of this source code is governed by the MIT license.
// See the LICENSE file for details.
package interfaces

// interfaces/fragment.go — the contract’s data model.

// CellID identifie une cellule dans le tissu (traçage, annuaire du
// Space). Donnée de contrat : GeneFragment.From l'exige, donc elle
// vit dans le paquet des contrats.
type CellID int

// GeneFragment is a piece of genetic DATA that passes through all 
// contracts (INucleus.CanAbsorb, ImmuneSystem.Screen, Harvest, hook 
// evoGo).
type GeneFragment struct {
	From    CellID   // transmitting cell
	Codons  []int    // genome fragment (codon sequence)
	Fitness float64  // transmitter fitness when sending

	// Perceptual channel: knowledge is consolidated by the originating cell 
	// at the end of the collection phase. Empty in the case of an excretion 
	// fragment (HGT) — only the COLLECTION snapshot carries the acquired 
	// knowledge (LearnedProfiles).
	Profile string   // emerging profile for the zone (e.g., "suffix:en")
	Support float64  // relative evidence of the profile (heuristic)
}