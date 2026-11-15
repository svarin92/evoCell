// Copyright 2026 Stéphane Varin. All rights reserved.
// Use of this source code is governed by the MIT license.
// See the LICENSE file for details.
package evocell

// message.go - the three vocabularies of the membrane: EXCRETION, ABSORPTION,
// SIGNALING.

// Note: if the module is "github.com/vladopajic/go-actor" (go.mod), the
// package itself is imported via its subdirectory: ".../go-actor/actor".
// Note: go-actor has no concept of a PID; an "address" is a MailboxSender
// reference held by the sender. CellID serves as a logical identifier.
import "github.com/vladopajic/go-actor/actor"

// CellID identifies a cell within the tissue (for tracing, or as an index 
// within the Space). Declared in interfaces/fragment.go; re-exported via 
// alias.go: it is a contract data element (required by GeneFragment.From).

// IMessage is the basic envelope: everything passes through the membrane. 
// The Message interface is maintained within the package that declares the
// concrete types, so as not to expose the message hierarchy to the outside 
// world. The scope of the package - private method isMessage() is limited 
// to the package boundary.
type IMessage interface{ isMessage() }

// --- 1. EXCRETION: horizontal gene transfer (HGT) ---

// GeneFragment is a fragment of codons extracted from a high-performing 
// genome. "The beneficial infection" is the success of a cell spreads.
type GeneFragmentMsg struct {
	GeneFragment  // the data (interface contract), embedded
}

// --- 2. ABSORPTION: incoming genetic material, to be filtered ---

// AbsorptionRequest reaches the membrane; nothing is integrated prior to 
// immune validation.
type AbsorptionRequest struct {
	Fragment GeneFragment
	Reply    actor.MailboxSender[AcceptVerdict]  // verdict rendered to the cell
}

type AcceptVerdict struct {
	Accepted bool
	Reason   string  // "non-self" | "low fitness" | "membrane sealed"
}

// --- 3. SIGNALING: quorum sensing ---

// Chemical signal diffusing through the tissue.
type SignalKind int

const (
	SignalStress SignalKind = iota  // fitness in decline → neighbors: mutate more
	SignalStability  // strong fitness → related: calm down, stabilize
)

// String provides a readable description of the "stress" and "stability" 
// signals.
func (k SignalKind) String() string {

	switch k {
	case SignalStress:
		return "stress"

	case SignalStability:
		return "stability"
	}

	return "unknown"
}

type ChemicalSignal struct {
	From    CellID
	Kind    SignalKind
	Fitness float64
}

// --- Internal messages: cell ↔ membrane ---

type ExcreteRequest struct{ Fragment GeneFragment }  // cell → membrane : excrete
type AbsorbVerdict struct{  // membrane → cell: fragment + immune verdict
	Fragment GeneFragment
	Verdict AcceptVerdict
}  
type BroadcastSignal struct{ Signal ChemicalSignal }  // membrane → tissue: quorum
type Tick struct{ Generation int }  // tissue clock (driven by the adapter)
type Division struct{ ChildCodons []int }  // result of a successful division

// --- Harvest: the tissue requests cellular genomes ---

// GetGenome snapshot request sent to a cell (collection). The response is 
// sent to the requester's response mailbox.
type GetGenome struct {
	Reply actor.MailboxSender[GenomeSnapshot]
}

// PorosityUpdate passes the Tick to the membranes by propagating the annealed 
// porosity of the current generation (a global value for the entire tissue).
type PorosityUpdate struct {
	Permeability float64
}

// --- Plasticity & perception (hook → Space → membrane → cell) ---
// Following the same path as PorosityUpdate: these messages originate outside 
// the tissue (the hook — a translation from the host engine), pass through 
// the Space, and are perceived by the membrane — which RELAYS them to its 
// cell via metabolism. The membrane is not merely a sensory organ; it is also 
// the sole gateway through which the external context is transformed into 
// specialization.

// DifferentiationSignal triggers a re-grafting of the nucleus. It is issued 
// by the hook (a discovery from evoGo), relayed by the Space to the membranes, 
// and passed on by each membrane to its respective cell. The cell receives a 
// zone specification — never the intention that motivates it.
type DifferentiationSignal struct {

	// The zone as read from the BNF symbol table CURRENT (non-terminal 
	// + terminals + profile).
	Spec ZoneSpec
}

// WorldSampleMsg is a sealed envelope containing a sample of the world
// (perception).
type WorldSampleMsg struct {
	WorldSample  // the data (interface contract), embedded
}

// v is a defensive copy of a cell's genetic material at the time of harvest
// — the ecosystem never sees the living genome.
type GenomeSnapshot struct {
	Codons  []int
	Fitness float64
	Nucleus string  // graft description (DOT node color)

	// Consolidated knowledge of the plastic core at the time of harvest
	// (extension of the perceptual channel) — empty if not consolidated.
	Profile string   // emerging profile of the zone (e.g., suffix: "en")
	Support float64  // relative prominence of the profile
}

func (GeneFragmentMsg) isMessage() {}
func (AbsorptionRequest) isMessage() {}
func (AcceptVerdict) isMessage() {}
func (ChemicalSignal) isMessage() {}
func (ExcreteRequest) isMessage() {}
func (AbsorbVerdict) isMessage() {}
func (BroadcastSignal) isMessage() {}
func (Tick) isMessage() {}
func (Division) isMessage() {}
func (GetGenome) isMessage() {}
func (PorosityUpdate) isMessage() {}
func (DifferentiationSignal) isMessage() {}  // plasticity: Space → membrane → cell
func (WorldSampleMsg) isMessage() {}  // perception: Space → membrane → cell