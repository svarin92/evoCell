// Copyright 2026 Stéphane Varin. All rights reserved.
// Use of this source code is governed by the MIT license.
// See the LICENSE file for details.
package renderer

// NewObserver construit l'Observateur. Le hook (adapter.go) le détient
// et le nourrit aux charnières de génération ; le binaire applicatif
// l'exporte en fin de run — jamais les cellules.
func NewObserver(cfg Config) *Observer {
	return &Observer{Zones: make(map[string]*ZoneView)}
}
