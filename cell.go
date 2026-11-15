// Copyright 2026 Stéphane Varin. All rights reserved.
// Use of this source code is governed by the MIT license.
// See the LICENSE file for details.
package evocell

// cell.go - the cell actor.

import (
	"github.com/vladopajic/go-actor/actor"
)

// Cell is the cellular agent: it carries a genome (codons), a fitness value,
// and a mutation rate modulated by received signals.
type Cell struct {
	metabolismMbx actor.Mailbox[IMessage]  // shared with its membrane
	codons        []int
	fitness       float64
	mutationRate  float64
	mutationBoost float64  // stress response — received from the tissue when constructing

	// diffFactory: injected at the initial seeding stage — it serves 
	// as the conduit for any re-grafting of the core (plasticity; cf. 
	// DifferentiationSignal). Nil indicates a fixed specialization 
	// (v0.1.0, without an injected factory).
	diffFactory   DifferentiationFactory

	// the nucleus is nil for a pure (totipotent) stem cell; once grafted, 
	// it specializes the metabolism (VowelCell, ConsonantCell...).
	nucleus      INucleus

	// profile is the current contextual sub-filtering (e.g., "suffix"),
	// updated by DifferentiationSignal — a refinement, not a new area.
	profile      string
	stressed     bool
}

// baselineMutationRate determines the resting metabolic rate. The 
// tick returns the cell to a calm state before mutation: stress is 
// a transient excitation, not a permanent state.
func (c *Cell) baselineMutationRate() float64 { return 0.02 }

// Within DoWork, the nucleus modulates behaviors:
// - mutate(): the domain's mutation rate (n.mutateRate)
// - AbsorbVerdict: dual filtering based on immunity and specialization
//   (if c.nucleus != nil && !c.nucleus.CanAbsorb(fragment) → rejection)
// - Tick: Express(codons) for excretion — a VowelCell only excretes
//   fragments useful for vowels (targeted horizontal transfer).
func (c *Cell) DoWork(ctx actor.Context) actor.WorkerStatus {

	select {
	case <-ctx.Done():
		return actor.WorkerEnd

	case msg := <-c.metabolismMbx.ReceiveC():

		switch v := msg.(type) {
		case Tick:
			c.mutationRate = c.baselineMutationRate()  // return to calm
			c.mutate()

		case AbsorbVerdict:

			// Fragment validated by the immunity: transgenesis. 
			// A segment of its genome is replaced by the fragment.
			if v.Verdict.Accepted {
				c.integrate(v.Fragment)
			}

		case BroadcastSignal:
	
			// Quorum sensing: the chemical signal modulates behavior.
			switch v.Signal.Kind {
			case SignalStress:

				// A neighbor is suffering → increase her mutation to seek
				// a solution FOR THE COLONY (collective intelligence).
				c.mutationRate *= c.mutationBoost

			case SignalStability:

				// A neighbor is solid → calm down, stabilize the area.
				c.mutationRate *= 0.5
			}			

		case GetGenome:

			// Response to a harvest request: defensive copy; the ecosystem 
			// never sees the live genome. The snapshot ALSO captures the 
			// consolidated knowledge of the plastic core (profile + evidence) 
			// — it is through this collection process that plasticity becomes 
			// exportable knowledge (LearnedProfiles), without ever reading a 
			// living core from the outside.
			snap := GenomeSnapshot{
				Codons: append([]int(nil), c.codons...),
				Fitness: c.fitness,
				Nucleus: c.nucleusDesc(),
			}


			if pn, ok := c.nucleus.(IPlasticNucleus); ok {
				snap.Profile, snap.Support = pn.LearnedProfile()
			}

			v.Reply.Send(ctx, snap)

		case DifferentiationSignal:

			// The membrane relays a differentiation command: the cell 
			// re-grafts its nucleus from the factory, configured
			// according to its specialization (based on data from the 
			// current BNF). It obeys without knowing the intention 
			// expressed by evoGo.
			if c.diffFactory != nil {
				c.GraftNucleus(c.diffFactory(v.Spec))
				c.profile = v.Spec.Profile  // refinement of the filtering
			}
			
		case WorldSampleMsg:

			// The membrane has perceived a sample of the world: the plastic 
			// core (if grafted) ingests the statistics and refines the
			// profile of its own zone. The cell does not understand the world
			// — it leaves it to its own interpretation.
			if pn, ok := c.nucleus.(IPlasticNucleus); ok {

				// Phase 1.
				pn.Observe(v.WorldSample)

				// Phase 2 immediately follows phase 1: the cell makes a 
				// proposal as soon as the thresholds are crossed — the 
				// observer conjectures, never the supervisor.
				pn.Consolidate()
			}

		}

		return actor.WorkerContinue
	}

}

// GraftNucleus grafts a specialized nucleus (differentiation of the 
// VowelCell, VerbCell, etc.) — mirror of the evoGo IOrganism mechanism.
func (c *Cell) GraftNucleus(n INucleus) { c.nucleus = n }

// integrate is transgenesis — the fragment validated by immunity replaces
// a segment of the genome. Defensive copy: the received fragment belongs
// to the message; it is never shared with the sender.
func (c *Cell) integrate(frag GeneFragment) {

	if len(frag.Codons) == 0 || len(frag.Codons) > len(c.codons) {
		return
	}

	// Genome head graft (anchor position: to be decided, e.g., by zone).
	copy(c.codons, frag.Codons)
	c.fitness = frag.Fitness  // donor fitness, until the next evaluation
}

// mutate performs a point mutation on the genome (at the current rate). 
// If a nucleus is grafted, it may restrict the mutation to its domain.
func (c *Cell) mutate() {

	if c.nucleus != nil {
		c.codons = c.nucleus.Express(c.codons)  // metabolic specialization
	}

	// ... occasional draw based on c.mutationRate (tissue rng, shared seed).
}

// nucleusDesc describes the plugin (zone + profile) for rendering (DOT 
// color), without exposing the nucleus itself.
func (c *Cell) nucleusDesc() string {

	if c.nucleus == nil {
		return "totipotent"
	}

	spec := c.nucleus.Zone()

	if spec.Profile != "" {
		return spec.NonTerminal + "/" + spec.Profile
	}

	return spec.NonTerminal
}