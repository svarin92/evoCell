// Copyright 2026 Stéphane Varin. All rights reserved.
// Use of this source code is governed by the MIT license.
// See the LICENSE file for details.
package evocell

// space.go - the tissue and the diffusion.

import (
	"github.com/vladopajic/go-actor/actor"
)

// Space is the tissue: registry of membranes (CellID → MailboxSender),
// collective statistics, and signal broadcasting engine (quorum sensing).
type Space struct {
	membranes []actor.Mailbox[IMessage]  // "ecosystem" mailboxes of each membrane
	cfg Config
	avg float64  // average tissue fitness (updated each generation)
}

// AverageFitness exposes the average fitness of the tissue (used by the 
// immune system).
func (s *Space) AverageFitness() float64 { return s.avg }

// Broadcast disseminates a generic message to all membranes (relaying signals 
// from the outside world: DifferentiationSignal; cf. plasticity).
func (s *Space) Broadcast(ctx actor.Context, msg IMessage) {

	for _, mbx := range s.membranes {
		_ = mbx.Send(ctx, msg)
	}

}

// BroadcastSignal broadcasts a chemical signal within the quorum radius.
func (s *Space) BroadcastSignal(ctx actor.Context, sig ChemicalSignal) {

	for _, mbx := range s.membranes {
		mbx.Send(ctx, sig)
	}

}

// Excrete releases a fragment to a neighboring cell (or cells).
func (s *Space) Excrete(ctx actor.Context, f GeneFragment) {

	// Here, the fragment is simply broadcast to everyone. However, a more
	// refined approach based on topology (Moore grid, etc.) is possible.
	for _, mbx := range s.membranes {
		mbx.Send(ctx, GeneFragmentMsg{GeneFragment: f})  // sealed envelope
	}

}

// PorosityAt calculates the permeability for the current generation.
func (s *Space) PorosityAt(gen int) float64 {
	return s.cfg.Porosity.Annealing(gen, s.cfg.MaxGenerations)
}

func (s *Space) Register(mbx actor.Mailbox[IMessage]) {
	s.membranes = append(s.membranes, mbx) 
}

// SetPorosity updates the porosity of all the tissue's membranes (called 
// by Tick: a single value, annealed according to the annealing profile).
func (s *Space) SetPorosity(p float64) {
	ctx := actor.ContextStarted()  // off-actor send (Tick relay)

	for _, mbx := range s.membranes {
		_ = mbx.Send(ctx, PorosityUpdate{Permeability: p})
	}

}
