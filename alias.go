// Copyright 2026 Stéphane Varin. All rights reserved.
// Use of this source code is governed by the MIT license.
// See the LICENSE file for details.
package evocell

import (
	"github.com/svarin92/evoCell/config"
	"github.com/svarin92/evoCell/interfaces"
)

// Interfaces imported to ensure architectural consistency.
type (
	CellID          = interfaces.CellID
	GeneFragment    = interfaces.GeneFragment
	ZoneSpec        = interfaces.ZoneSpec
	WorldSample     = interfaces.WorldSample

	INucleus        = interfaces.INucleus
	IOrgan          = interfaces.IOrgan
	IPlasticNucleus = interfaces.IPlasticNucleus

	DifferentiationFactory = interfaces.DifferentiationFactory
)

// Alias ​​for concrete types.
type (
	Config          = config.Config
)
