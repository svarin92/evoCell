// Copyright 2026 Stéphane Varin. All rights reserved.
// Use of this source code is governed by the MIT license.
// See the LICENSE file for details.
package evocell

// tissue.go - assemble the superorganism.

import (
	"sync"
	"time"

	"github.com/vladopajic/go-actor/actor"
)

/* ProfileView */

type ProfileView struct {
	Profile string
	Support float64
}

/* Tissue */

// Assemble tissue: N cells + N membranes + 1 immune system into a single 
// combined actor (atomic Start/Stop). The tissue rests upon four operations
// — seedling, life, harvest, death.
type Tissue struct {
	cfg         Config
	space       *Space
	immune      *ImmuneSystem
	combined    actor.Actor

	// The "metabolism" mailboxes (cell ↔ membrane), retained so that 
	// the tissue can trigger Tick and GetGenome FROM the non-actor world 
	// (the evoGo hook is not an actor).
	metabolisms []actor.Mailbox[IMessage]  // to communicate with the cells

	// The organs assembled by the tissue (one for each non-terminal parent 
	// in the BNF) — these are lightweight objects, not actors.
	organs []*Organ

	// It is populated during the seeding phase by GrammarZoning (each 
	// grafted cell serves a host non-terminal).
	zoneOf map[CellID]string

	// The hierarchy from a BNF parent to its child zones is inferred by 
	// GrammarZoning (the organs read it).
	zonesOfParent map[string][]string
}

// AssignZone assigns a cell to its grammatical zone. This occurs during
// DIFFERENTIATION (the grafting of the NucleusGrammatical), rather than 
// during the initial seeding phase: at the time of seeding, cells are 
// totipotent (having a null nucleus), and their zone is only determined 
// once a grammatical nucleus is grafted. GrammarZoning—or specialized 
// factories such as NewVerbCell—invokes AssignZone, passing the host 
// non-terminal. If this is not called, the cell remains "unZoned" during 
// the harvesting process (HarvestByZone).
func (t *Tissue) AssignZone(id CellID, nonTerminal string) {
	t.zoneOf[id] = nonTerminal
}

// Broadcast propagates a message from the outside world (hook) to all
// membranes via the Space.
func (t *Tissue) Broadcast(msg IMessage) {
	t.space.Broadcast(actor.ContextStarted(), msg)
}

// Harvest collects the genetic snapshots of cells after generation. 
// Asynchronous request-reply harvesting: for each cell, a temporary 
// response mailbox receives the GenomeSnapshot. Returns the available 
// fragments (cells respond at their own pace — non-blocking).
func (t *Tissue) Harvest() []GeneFragment {
	ctx := actor.ContextStarted()  // off-actor dispatch (see Tick)
	fragments := make([]GeneFragment, 0, len(t.metabolisms))
	var mu sync.Mutex  // serializes the appends from the collection goroutines
	var wg sync.WaitGroup

	for i, met := range t.metabolisms {
		reply := actor.NewMailbox[GenomeSnapshot]()
		reply.Start()
		_ = met.Send(ctx, GetGenome{Reply: reply})
		wg.Add(1)

		go func(id CellID) {
			defer wg.Done()

			select {
			case snap := <-reply.ReceiveC():

				// A live snapshot becomes a fragment available to the pool,
				// tracked back to its origin cell (From = CellID).
				// Append operation under a mutex: since collection goroutines
				// are concurrent, the shared slice must be protected against
				// data races.
				mu.Lock()

				// A live snapshot becomes a fragment available in the pool.
				fragments = append(fragments, GeneFragment{
					From: id, 
					Codons: snap.Codons,
					Fitness: snap.Fitness,
				})

				mu.Unlock()

			case <-time.After(50 * time.Millisecond):
			
				// Slow or dead cell: we harvest without it.

			}

		}(CellID(i))

	}

	wg.Wait()

	return fragments
}

// HarvestByOrgan coordinates the harvest. The Deposit receives organ 
// fragments (coherent multi-zones), targeted at the non-terminal parent 
// of weak individuals.
func (t *Tissue) HarvestByOrgan() map[string][]GeneFragment { 
	organs := make(map[string][]GeneFragment)

	for _, organ := range t.organs {

		if organ.Coherence() < 0.5 {
			continue  // organe disharmonique : on ne livre pas un
					  // assemblage qu'un Validate rejetterait
		}

		perZone := t.HarvestByZone()
		fragments := make(map[string]GeneFragment)

		for _, zone := range organ.zones {

			if frags, ok := perZone[zone]; ok && len(frags) > 0 {
				fragments[zone] = BestFragmentOf(frags)
			}

		}

		if organ.Validate(fragments) {
			
			// L'assemblage validé devient un fragment d'organe :
			// concaténation ordonnée des contributions des zones.
			organs[organ.NonTerminal()] = assemble(fragments, organ.zones)

		}

	}

	return organs 
}

// HarvestByZone harvests the fragments and groups them by grammatical zone: 
// the key is the non-terminal host of the donor cell. It is this structuring 
// that allows the Deposit to target an individual weak GE precisely on the 
// non-terminal where it derives poorly.
func (t *Tissue) HarvestByZone() map[string][]GeneFragment {
	perZone := make(map[string][]GeneFragment)

	for _, frag := range t.Harvest() {
		zone, known := t.zoneOf[frag.From]

		if !known {
			zone = "unzoned"  // non-grafted totipotent cell
		}

		perZone[zone] = append(perZone[zone], frag)
	}

	return perZone
}

// LearnedProfiles reveals knowledge of the tissue — the evidence of 
// emerging profiles by zone. It is the VALORIZATION of plasticity: 
// the  learned heuristic becomes a first-class data point, inspectable 
// and exportable (Observer), rather than a hidden setting.
//
// ACTOR SAFETY: knowledge must NOT be read directly from live nuclei 
// (to avoid data races from outside the actor realm) — instead, it is 
// accessed via an existing snapshot (GetGenome → GenomeSnapshot, using 
// a defensive copy delivered via the mailbox) and then aggregated by 
// zone using the `zoneOf` registry.
func (t *Tissue) LearnedProfiles() map[string]ProfileView {
	out := make(map[string]ProfileView)

	for _, frag := range t.Harvest() {

		if frag.Profile == "" {
			continue
		}

		zone, known := t.zoneOf[frag.From]

		if !known {
			zone = "unzoned"
		}

		out[zone] = ProfileView{Profile: frag.Profile, Support: frag.Support}
	}

	return out
}

func (t *Tissue) Start() { 
	t.combined.Start()
}

func (t *Tissue) Stop() { 
	t.combined.Stop() 
}

// Tick advances the generation clock for all cells and updates the 
// porosity of all membranes. Called via the hook from the non-actor 
// world (OnGenerationStart): communication occurs through mailboxes, 
// never via direct method calls — the outside world speaks to the 
// tissue ONLY through messages.
func (t *Tissue) Tick(gen int) {

	// Non-actor sending context: ContextStarted() is the canonical
	// go-actor context for writing to a mailbox from the outside
	// world (see FanOut in the library). Idiomatic and safe: never 
	// cancelled, it simply signals that "the system is alive."
	ctx := actor.ContextStarted()

	tick := Tick{Generation: gen}
	porosity := t.space.PorosityAt(gen)

	for _, met := range t.metabolisms {

		// 1. The cell receives the clock signal (return to quiescence, 
		//    mutation).
		_ = met.Send(ctx, tick)  // ErrMailboxStopped ignored: dead cell

	}

	// 2. The membrane updates its porosity via its ecosystem channel:
	//    the Space relays the porosity tick to the membranes.
	t.space.SetPorosity(porosity)

	// 3. The immunity recalculates its reference (average tissue 
	//    fitness).
	t.immune.UpdateAverage(t.space.AverageFitness())
}

// ZoneOf returns the child zones of a parent non-terminal (a BNF hierarchy 
// inferred by GrammarZoning) — the organs assemble within them.
func (t *Tissue) ZonesOf(parentNT string) []string {
	return t.zonesOfParent[parentNT]
}