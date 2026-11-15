// Copyright 2026 Stéphane Varin. All rights reserved.
// Use of this source code is governed by the MIT license.
// See the LICENSE file for details.
package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"

	"github.com/svarin92/evoGo/evaluator"
	"github.com/svarin92/evoGo/ge"
	"github.com/svarin92/evoGo/grammar"
	"github.com/svarin92/evoGo/operators"

	"github.com/svarin92/evoCell/config"
	"github.com/svarin92/evoCell/evogo"
	"github.com/svarin92/evoCell/renderer"
)

// examples/hybridation/main.go - Hybridization of evoGo with evoCell on 
// letter.bnf and the "golden" target. "golden" is a mechanical example 
// that does not, on its own, truly showcase the value of evoCell, since 
// evoGo is capable of rapidly reaching the target without the contribution 
// of the tissue. In contrast, a simple target scenario more effectively 
// illustrates the cellular system's mechanics: seeding, ticks, excretions, 
// signals, porosity, harvesting, and rendering. To fully demonstrate the 
// value of evoCell, the "recovery-after-stagnation" regime described below 
// should be used: evoGo runs in isolation until a fitness plateau is reached, 
// at which point the graft (hook) causes the progression curve to rise again.

// Two selectable modes:
// - tissue at-start: the tissue is active from generation 1 (mechanics demo)
// - tissue rescue-at=20: evoGo runs independently until generation 20, then
//   the hook is grafted: the tissue "rescues" a stagnant population (value
//   demo).
var CLI struct {
	File string `arg:"" type:"existingfile" help:"Grammar file (e.g. letter.bnf)."`
	Target string `short:"t" default:"golden" help:"Target phenotype."`
	Tissue string `short:"c" enum:"off,at-start,rescue" default:"at-start" help:"Tissue regime."`
	RescueAt int `default:"20" help:"Generation at which the rescue hook is grafted."`
	Out string `short:"o" default:"." help:"Output directory for tissue.dot and pool.csv."`
}

func main() {
	ctx := kong.Parse(&CLI)

	// Open the grammar file provided as an argument.
	r, err := os.Open(CLI.File)
	ctx.FatalIfErrorf(err, "failed to open grammar")
	defer r.Close()
	
	// Create a new grammar based on the file.
	grammar, err := grammar.NewSpeciesGrammar(ctx, r)
	ctx.FatalIfErrorf(err, "failed to create species grammar")

	// Create an evaluation function that measures the proximity of 
	// individuals to the target.
	fitnessFunc := evaluator.NewFitness(CLI.Target)

	// replacementFunc defines the generational replacement strategy.
	replacementFunc := func(newIndividuals, oldIndividuals []IIndividual, eliteSize int) []IIndividual { 
		return operators.GenerationalReplacement(newIndividuals, oldIndividuals, eliteSize) 
	}
	
	// selectionFunc defines the strategy for selecting parents for 
	// reproduction.
	selectionFunc := func(individuals []IIndividual, tournamentSize int) []IIndividual { 
		return operators.TournamentSelection(individuals, tournamentSize) 
	}

	// Tissue configuration.
	cfg := config.Config{
		Cells: 50,  // Cells: free cardinality (cf. flow, not identity)
		MaxGenerations: 100,
		MutationBoost: 2.5,  // Response to stress signal
		DivisionThreshold: 0.8,
		DeathThreshold: 0.2,
		Porosity: config.PorosityProfile{
			StartPermeability: 0.9,  // porous: exploration (extensive exchange)
			EndPermeability: 0.05,   // sealed: exploitation
			Annealing: config.LinearAnnealing(0.9, 0.05),
		},

	}

	// The Tissue and Pool Observer.
	obs := renderer.NewObserver(cfg)
	obs.Start()
	defer obs.Stop()

	// The hook: the ONLY hybridization line. It attaches to the Genomizer
	// OF the population (nil-safe pass-through on the evoGo side).
	hook := evogo.NewHook(cfg, obs)

	var opts []ge.SearchOption

	switch CLI.Tissue {
	case "at-start":
		opts = append(opts, ge.WithHybridizationHook(hook))  // Hooked since Gen 1
	
	case "rescue":
	
		// The tissue only comes into play from RescueAt onwards.
		opts = append(opts, ge.WithHybridizationHook(evogo.NewLazyHook(hook, CLI.RescueAt)))

	case "off":
		
		// No hook: evoGo without membrane, baseline behavior.
	}

	bestEver, err := ge.SearchLoop(
		Generations,      // Total number of generations
		PopulationSize,   // Population size
		grammar,          // Grammar used to generate individuals
		CLI.Target,       // Target to be reached: "golden"
		replacementFunc,  // Generational replacement function (here, elitist replacement)
		selectionFunc,    // Parent selection function (here, tournament selection)
		fitnessFunc,      // Evaluation function
		opts...,          // options: WithHybridizationHook (or nothing)
	)
	ctx.FatalIfErrorf(err, "evolution failed: %v\n")

	// Cellular renderer.
	_ = obs.ExportTissue(CLI.Out + "/tissue.dot")     // zones, cells, HGT
	_ = obs.ExportPoolMetrics(CLI.Out + "/pool.csv")  // pool enrichment
	fmt.Printf("Best individual: %s\n", bestEver)
}