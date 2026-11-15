// Copyright 2026 Stéphane Varin. All rights reserved.
// Use of this source code is governed by the MIT license.
// See the LICENSE file for details.
package interfaces

// interfaces/zonespec.go — specialization as data.

// ZoneSpec describes a zone containing the information the tissue needs to 
// know — everything else is read from the evoGo grammar AT THE MOMENT of 
// differentiation. evocell contains NO grammar: the non-terminal is data, 
// not a type. If the BNF changes, the next generation re-zones based on the 
// new symbol table — the tissue follows the grammar, never the other way 
// around.
type ZoneSpec struct {

	// NonTerminal is the host area ANCHORED within the current BNF. It serves 
	// as the attachment point where a collected fragment can be integrated, 
	// in accordance with the derivation (evoGo.ProductionHistory).
	NonTerminal string

	// The terminals that this zone can produce serve as the raw material for 
	// filtering (CanAbsorb).
	Terminals []string

	// The potential sub-zones (BNF hierarchy) where the organs will assemble.
	Children []string

	// The contextual subcategory (e.g., "root", "suffix"), discovered on the 
	// fly by evoGo (cf. WeakestNT). An empty profile corresponds to the basic 
	// filtering of the non-terminal.
	Profile string
}