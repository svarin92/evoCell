// Copyright 2026 Stéphane Varin. All rights reserved.
// Use of this source code is governed by the MIT license.
// See the LICENSE file for details.
package evogo

import "github.com/svarin92/evoCell/renderer"

// NewHook builds the adapter, configuring the tissue and the observer.
func NewHook(cfg Config, obs *renderer.Observer) *Hook {
    return &Hook{cfg: cfg, obs: obs}
}

// LazyHook is a nil-safe decorator: it ignores calls until the activation 
// generation is reached ("rescue" mode — the tissue kicks in when evoGo 
// stalls).
func NewLazyHook(inner IHybridizationHook, activateAt int) *LazyHook {
	return &LazyHook{inner: inner, activateAt: activateAt}
}
