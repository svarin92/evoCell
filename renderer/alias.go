// Copyright 2026 Stéphane Varin. All rights reserved.
// Use of this source code is governed by the MIT license.
// See the LICENSE file for details.
package renderer

import (
	evocell "github.com/svarin92/evoCell"
	"github.com/svarin92/evoCell/config"
)

// Interfaces imported to ensure architectural consistency.
type (
	IMessage      = evocell.IMessage
)

// Alias ​​for concrete types.
type (
	Config       = config.Config

	CellID       = evocell.CellID
	GeneFragment = evocell.GeneFragment
)
