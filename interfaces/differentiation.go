// Copyright 2026 Stéphane Varin. All rights reserved.
// Use of this source code is governed by the MIT license.
// See the LICENSE file for details.
package interfaces

// DifferentiationFactory builds a nucleus for a specialization RELATED TO
// EXECUTION within the evoGo symbol table. Mirror of the `evoGo` framework’s 
// `OrganismFactory` interface; it is configured by grammatical data, without 
// any a priori knowledge of the possible grammars. This is the prerequisite 
// for plasticity — the membrane never assembles a nucleus itself; instead, it 
// requests one from the factory (which is injected when the tissue is seeded).
type DifferentiationFactory func(spec ZoneSpec) INucleus