// Copyright 2026 Stéphane Varin. All rights reserved.
// Use of this source code is governed by the MIT license.
// See the LICENSE file for details.
package evocell

// grammarzoning.go — zones derived from evoGo's BNF.

import "github.com/vladopajic/go-actor/actor"

// GrammarZoning constructs the topology of the tissue from the grammar 
// itself: each non-terminal in the BNF becomes a zone, and each production 
// rule becomes a neighborhood contract. There is no duplication: the tissue 
// is structured starting from evoGo.
type GrammarZoning struct {
	zones map[string]*Zone  // non-terminal → specialized zone
}

type Zone struct {
	NonTerminal string
	Metabolisms []actor.Mailbox[IMessage]  // cells in the zone
	Neighbors []string  // adjacent non-terminals (rules)
	SignalBudget float64  // local quorum intensity
}

// PropagateQuorum propagates the quorum across zones: a noun phrase (NP) 
// zone that has reached syntactic agreement emits a stability signal that 
// dampens mutations in the neighboring verb phrase (VP) zone (an emergent,
// not imposed, subject-verb agreement).
func (z *GrammarZoning) PropagateQuorum(ctx actor.Context, from string, sig ChemicalSignal) {

	for _, nt := range z.zones[from].Neighbors {

		if zone, ok := z.zones[nt]; ok {

			for _, mbx := range zone.Metabolisms {
				mbx.Send(ctx, sig)
			}

		}

	}

}