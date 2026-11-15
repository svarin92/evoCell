// Copyright 2026 Stéphane Varin. All rights reserved.
// Use of this source code is governed by the MIT license.
// See the LICENSE file for details.
package evogo

import (
	"github.com/svarin92/evoCell"
	"github.com/svarin92/evoCell/config"

	"github.com/svarin92/evoCell/renderer"

	"github.com/svarin92/evoGo/interfaces"
)

// Interfaces imported to ensure architectural consistency.
type (
	IGenomeCarrier         = interfaces.IGenomeCarrier
	IHybridizationHook     = interfaces.IHybridizationHook
	IMorphologicalPatterns = interfaces.IMorphologicalPatterns
	IIndividual            = interfaces.IIndividual
	IRuleModel             = interfaces.IRuleModel
)

// Alias ​​for concrete types.
type (
	Config                = config.Config

	DifferentiationSignal = evocell.DifferentiationSignal
	GeneFragment          = evocell.GeneFragment
	Tissue                = *evocell.Tissue

	Observer              = *renderer.Observer 
)