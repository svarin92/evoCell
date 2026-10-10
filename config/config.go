// Copyright 2026 Stéphane Varin. All rights reserved.
// Use of this source code is governed by the MIT license.
// See the LICENSE file for details.
package config

// PorosityProfile manages the membrane's opening policy over successive 
// generations (early exploration, late exploitation).
type PorosityProfile struct {

	// Maximum exchange rate at the beginning of cycles (e.g., 0.9).
	StartPermeability float64

	// Minimum exchange rate at the end of cycles (e.g., 0.05).
	EndPermeability float64

	// Linear or sigmoid (the 'true' value for the final cycle corresponds 
	// to the "last generation").
	Annealing func(generation, maxGenerations int) float64
}

func LinearAnnealing(start, end float64) func(int, int) float64 {
	return func(gen, maxGen int) float64 {

		if maxGen <= 0 {
			return start
		}

		t := float64(gen) / float64(maxGen)
		return start + (end-start)*t
	}

}

type Config struct {
	
	// Number of cells in the tissue.
	Cells int  

	MaxGenerations int // for porosity annealing

	// Minimum fitness required for a cell to divide.
	DivisionThreshold float64

	// Maximum fitness for a cell to die.
	DeathThreshold float64

	// Mutation rate increase factor when a stress signal is received 
	// (collective response).
	MutationBoost float64

	// Signal propagation range (0 = entire tissue).
	QuorumRadius int

	// Membrane permeability profile.
	Porosity PorosityProfile
}