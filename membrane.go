// Copyright 2026 Stéphane Varin. All rights reserved.
// Use of this source code is governed by the MIT license.
// See the LICENSE file for details.
package evocell

// membrane.go - the central actor.

import (
	"math/rand"

	"github.com/vladopajic/go-actor/actor"
)

// Membrane is an actor: the cell's sole entry and exit point. It has TWO 
// mailboxes: one for exchanges with its cell (metabolism), and one for 
// exchanges with the tissue (ecosystem).
type Membrane struct {

	// metabolism: private channel between cell and membrane.
	metabolismMbx actor.Mailbox[IMessage]

	// Ecosystem: channel with neighboring membranes (via Space).
	tissueMbx     actor.Mailbox[IMessage]
	immune        *ImmuneSystem
	permeability  float64  // updated on every tick
	space         *Space   // neighbor directory (CellID → MailboxSender)
	self          CellID

	// Random membrane decisions (whether or not to excrete a fragment,
	// with a probability equal to the current porosity).
	rng           *rand.Rand  // Random Number Generator
}

// DoWork is the membrane actor's loop.
func (m *Membrane) DoWork(ctx actor.Context) actor.WorkerStatus {

	select {
	case <-ctx.Done():
		return actor.WorkerEnd
	
	case msg := <-m.metabolismMbx.ReceiveC():
		m.handleMetabolism(ctx, msg)
		return actor.WorkerContinue
	
	case msg := <-m.tissueMbx.ReceiveC():
		m.handleTissue(ctx, msg)
		return actor.WorkerContinue
	}

}

// handleMetabolism processes requests for its cell.
func (m *Membrane) handleMetabolism(ctx actor.Context, msg IMessage) {

	switch v := msg.(type) {
	case ExcreteRequest:

		// --- EXCRETION — Horizontal Gene Transfer (HGT) ---
		// Send if the probability corresponds to the current porosity: 
		// porous at the start of the cycle (exploration), impermeable 
		// at the end (exploitation).
		if m.rng.Float64() <= m.permeability {
			m.space.Excrete(ctx, v.Fragment)
		}

	case Tick:

		// The tissue clock updates the porosity (annealing).
		m.permeability = m.space.PorosityAt(v.Generation)

	}

}

// handleTissue processes what comes from the outside.
func (m *Membrane) handleTissue(ctx actor.Context, msg IMessage) {

	switch v := msg.(type) {
	case GeneFragmentMsg:

		// --- ABSORPTION —--
		// Nothing is integrated directly: immune filtering first, verdict 
		// afterwards. (The v.GeneFragment data is unfolded from the envelope 
		// for analysis; immunity and the nucleus pass judgment on a piece of 
		// DATA.)
		frag := v.GeneFragment
		verdict := m.immune.Screen(frag, m.space.AverageFitness(), m.permeability)
	
		if verdict.Accepted {

			// Transfer validated → relayed to the cell via metabolism.
			m.metabolismMbx.Send(ctx, AbsorbVerdict{Verdict: verdict})

		}

	case ChemicalSignal:

		// --- SIGNALING — quorum sensing--- 
		// The signal crosses the membrane and influences the cell's behavior.
		m.metabolismMbx.Send(ctx, BroadcastSignal{Signal: v})
	

	case PorosityUpdate:
		
		// TICK — the annealed porosity of the current generation.
		m.permeability = v.Permeability

	case DifferentiationSignal:

		// PLASTICITY — a command for re-grafting originating from the 
		// outside world (via the Space). The membrane RELAYS this to its 
		// cell: it is the sole gateway to the outside world, converting 
		// context into specialization — the cell itself never hears the 
		// rationale.
		m.metabolismMbx.Send(ctx, v)

	case WorldSampleMsg:

		// PERCEPTION — a sample of the world (a corpus delivered by the
		// hook, drop by drop). The membrane is the cell’s sensory organ: 
		// it perceives, it relays — the nucleus will learn what these 
		// morphological statistics signify for ITS own zone.
		m.metabolismMbx.Send(ctx, v)

	}

}