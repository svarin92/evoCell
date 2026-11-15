// Copyright 2026 Stéphane Varin. All rights reserved.
// Use of this source code is governed by the MIT license.
// See the LICENSE file for details.
package renderer

import (
	"github.com/vladopajic/go-actor/actor"
)

// Observer est un acteur qui reçoit des COPIES des messages du tissu
// (via le coordinateur) et maintient l'état visuel du système cellulaire.
// Pattern Spy/Observer : rendu non intrusif, zéro impact sur le métabolisme.
type Observer struct {
	Mbx actor.Mailbox[IMessage]  // mailbox d'abonnement (branches par le coordinateur)

	// État visuel agrégé par record() — jamais lu par le métabolisme.
	Zones  map[string]*ZoneView   // Z : cellules par zone, Nucleus greffé
	HGT    []TransferView         // M : excrétions (origine, cible, génération)
	Immune []VerdictView          // I : verdicts self/non-self + motifs appris
	Quorum map[string][]TickView  // signaux stress/stabilité par zone et tick
	Pool   []PoolRow              // T : diversité/qualité du bassin par génération
}

// --- Types de vue : l'état visuel agrégé par record() ---
// Chaque vue est une projection plate et immuable d'un message du tissu :
// l'Observateur ne garde jamais de référence sur un objet vivant.

// ZoneView (Z) : une zone grammaticale et ses cellules.
type ZoneView struct {
	Name string  // non-terminal hôte (ex. "vowel")
	Cells []CellID  // cellules présentes
	Nuclei map[CellID]string  // Nucleus greffé (couleur du noeud DOT)
}

// TransferView (M) : un transfert HGT, pour les arêtes du graphe.
type TransferView struct {
	From, To CellID
	Gen, Tick int
	Codons int  // taille du fragment (épaisseur ∝ porosité)
}

// VerdictView (I) : un verdict immunitaire ou un motif toxique appris.
type VerdictView struct {
	Cell CellID
	Accepted bool
	Reason string // "non-self" | "low fitness" | "membrane sealed"
	Toxic bool    // true : LearnToxic (motif appris)
}

// TickView (quorum) : l'état des signaux d'une zone à un tick.
type TickView struct {
	Gen, Tick int
	Stress, Stable int // signaux reçus dans la zone
}

// PoolRow (T) : une ligne de pool.csv, par génération.
type PoolRow struct {
	Gen int
	Zone string
	Fragments int
	Diversity float64
	AvgDonor float64 // fitness moyenne des donneurs
}

// Start lance l'acteur (branché sur le coordinateur, cf. D1).
func (o *Observer) Start() { /* actor.Start() — cf. cycle de vie du tissu */ }

// Stop arrête l'acteur et draine la mailbox (exports déjà sûrs après).
func (o *Observer) Stop() { /* actor.Stop() */ }

// Ce que l'Observateur collecte et rend :
// - Tissu & zones (Z) : graphe DOT — clusters = zones, noeuds = cellules
// (couleur = Nucleus greffé : VowelCell, VerbCell...)
// - Circulation (M) : arêtes = excrétions HGT, épaisseur ∝ porosité (P)
// - Verdicts immunitaires (I): journal self/non-self + carte de chaleur
// des motifs toxiques appris (LearnToxic)
// - Quorum sensing : traces temporelles des signaux stress/stabilité
// par zone (accord GN↔GV visible !)
// - Flux GE↔tissu (T) : séries temporelles par zone : diversité et
// qualité du pool, par génération

// Record agrège une copie de message dans l'état visuel. Fonction pure
// vis-à-vis du tissu : aucun envoi, aucun blocage, aucun retour.
func (o *Observer) Record(msg IMessage) { /* ... */ }

// RecordGenerationEnd : point d'entrée public pour le hook — consigne
// une charnière de génération (greffes effectuées, pool par zone).
// Délègue à record() après enveloppage : l'Observateur ne connaît du
// hook que des messages, jamais de références vivantes.
func (o *Observer) RecordGenerationEnd(
	gen int, 
	grafted int,
	perZone map[string][]GeneFragment,
) {
	// ... enveloppe la charnière en message interne et agrège (pool.csv)
}

func (o *Observer) DoWork(ctx actor.Context) actor.WorkerStatus {

	select {
	case <-ctx.Done():
		return actor.WorkerEnd

	case msg := <-o.Mbx.ReceiveC():
		o.Record(msg) // agrège dans son état visuel (aucun envoi, aucun blocage)
		return actor.WorkerContinue
	}

}

// ExportTissue : graphe du tissu en DOT (mêmes conventions que best_ever.dot).
func (o *Observer) ExportTissue(path string) error { /* ... */ return nil }

// ExportPoolMetrics : CSV des métriques du bassin par zone et génération —
// la « courbe d'enrichissement du pool génique », mesure de ce que le tissu
// APPORTE (aucun équivalent côté evoGo).
func (o *Observer) ExportPoolMetrics(path string) error { /* ... */ return nil }
