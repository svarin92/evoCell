// Copyright 2026 Stéphane Varin. All rights reserved.
// Use of this source code is governed by the MIT license.
// See the LICENSE file for details.
package evocell

import (
	"math/rand"

	"github.com/vladopajic/go-actor/actor"
)

/* Exports */

func NewCell(codons []int, mbx actor.Mailbox[IMessage], mutationBoost float64) *Cell {
	return &Cell{
		metabolismMbx: mbx, 
		codons: codons, 
		mutationRate: 0.02,
		mutationBoost: mutationBoost,  // provided by the tissue from Cfg.MutationBoost
	}
}

// NewGrammarZoning iterates through the non-terminals of evoGo (via the 
// adapter, which is the only component that imports interfaces.IRuleModel) 
// and creates one zone per rule.
func NewGrammarZoning(
	nonTerminals []string,
	rules map[string][]string,
) *GrammarZoning {
	z := &GrammarZoning{zones: make(map[string]*Zone)}

	for nt, neighbors := range rules {
		z.zones[nt] = &Zone{NonTerminal: nt, Neighbors: neighbors}
	}

	return z
}

func NewImmune(seed int64) *ImmuneSystem {
	return &ImmuneSystem{
		failedPatterns: make(map[string]float64),
		rng: rand.New(rand.NewSource(seed)),
	}

}

func NewMembrane(self CellID, space *Space, immune *ImmuneSystem, seed int64) *Membrane {
	return &Membrane{
		metabolismMbx: actor.NewMailbox[IMessage](),
		tissueMbx:     actor.NewMailbox[IMessage](),
		immune:        immune,
		space:         space,
		self:          self,
		rng:           rand.New(rand.NewSource(seed)),  // seed specific to the membrane
	}
}

// NewNucleusFor is the mandatory generic factory for any differentiation based on 
// specialization.
func NewNucleusFor(spec ZoneSpec) INucleus {
	return &NucleusGrammatical{
		spec: spec,
		profileFilter: filterOnTerminals(spec),
	}
}

// NewOrgan assembles the organ based on the BNF: the parent non-terminal 
// and  its child zones are derived from the grammar (GrammarZoning) — 
// never  invented. An organ without a BNF anchor is stillborn.
func NewOrgan(tissue *Tissue, parentNT string) *Organ {
	o := &Organ{
		nonTerminal: parentNT,
		tissue: tissue,
	}

	// Member zones: the non-terminals that are children of the parent, 
	// already registered by GrammarZoning (no duplication).
	for _, zone := range tissue.ZonesOf(parentNT) {
		o.zones = append(o.zones, zone)
	}

	return o
}

func NewTissue(cfg Config, genomes [][]int) *Tissue {
	space := &Space{cfg: cfg}
	immune := NewImmune(cfg.Seed)
	t := &Tissue{cfg: cfg, 
		space: space, 
		immune: immune, 
		zoneOf: make(map[CellID]string),
		zonesOfParent: make(map[string][]string),  // filled in by GrammarZoning
	}
	actors := make([]actor.Actor, 0, 2*cfg.Cells)

	for i := 0; i < cfg.Cells; i++ {

		// "met" = metabolismMbx: the private channel between the cell and 
		// its membrane.
		met := actor.NewMailbox[IMessage]()
		
		mem := NewMembrane(CellID(i), space, immune, cfg.Seed+int64(i))
		cell := NewCell(genomes[i], met, cfg.MutationBoost)
		space.Register(mem.tissueMbx)
		t.metabolisms = append(t.metabolisms, met)

		// Each mailbox is an actor; cell and membrane share "met" 
		// (metabolism).
		actors = append(actors, met, mem.tissueMbx, actor.New(mem), actor.New(cell))
	}

	t.combined = actor.Combine(actors...).Build()
	return t
}