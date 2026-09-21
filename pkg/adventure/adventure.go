package adventure

import (
	"fmt"
	"math/rand"
	"projet-red/pkg/character"
	"projet-red/pkg/combat"
	"projet-red/pkg/ui"
)

// choix de la zone, les zones 2 et 3 demandent un niveau minimum
func Explorer(p *character.Personnage) {
	for {
		ui.ClearConsole()
		ui.AfficherTitre(fmt.Sprintf("Zones d'Aventure (Mode NG+%d)", p.NGPlus))
		ui.DisplayZoneMapArt()

		fmt.Println("1. Forêt des Gobelins (Niveau 1+)")
		fmt.Println("2. Caverne des Trolls (Niveau 3+)")
		fmt.Println("3. Donjon du Dragon (Niveau 5+)")
		fmt.Println("4. Antre du Dieu Démon (BOSS FINAL)")
		fmt.Println("0. Retour")

		choix := ui.LireEntree("\nDestination : ")

		switch choix {
		case "1":
			LancerEvenement(p, 1)
		case "2":
			if p.Niveau < 3 {
				fmt.Println(ui.Rouge + "Zone restreinte ! Niveau 3 requis." + ui.Reset)
				ui.LireEntree("Appuyez sur Entrée...")
				continue
			}
			LancerEvenement(p, 2)
		case "3":
			if p.Niveau < 5 {
				fmt.Println(ui.Rouge + "Zone restreinte ! Niveau 5 requis." + ui.Reset)
				ui.LireEntree("Appuyez sur Entrée...")
				continue
			}
			LancerEvenement(p, 3)
		case "4":
			// pas de niveau requis pour le boss
			CombattreBossFinal(p)
		case "0":
			return
		}
	}
}

// tire un monstre au hasard dans la zone et lance le combat
func LancerEvenement(p *character.Personnage, zone int) {
	ui.ClearConsole()

	// Facteurs multiplicatifs NG+ / Niveau du joueur
	multStat := 1.0 + (float64(p.NGPlus) * 0.5) + (float64(p.Niveau-1) * 0.1)
	multRecompense := 1.0 + (float64(p.NGPlus) * 0.8)

	// Bestiaire varié selon la zone
	var m *combat.Monstre

	// le taux de drop monte de 10% par NG+ dans toutes les zones
	if zone == 1 {
		mobs := []combat.Monstre{
			{
				Nom: "Gobelin Éclaireur",
				PV:  35, MaxPV: 35, Mana: 20, MaxMana: 20, Attaque: 8, Defense: 2, Dexterite: 9,
				GainXP: int(30 * multRecompense), GainOr: int(8 * multRecompense),
				Sort:      combat.SortMonstre{Nom: "Tir de Flèche Empoisonnée", Degats: 12, CoutMana: 10},
				TableDrop: []string{"Plume de corbeau", "Potion de soin"}, ChancDrop: 30 + (p.NGPlus * 10),
			},
			{
				Nom: "Loup Sauvage",
				PV:  45, MaxPV: 45, Mana: 0, MaxMana: 0, Attaque: 11, Defense: 3, Dexterite: 12,
				GainXP: int(40 * multRecompense), GainOr: int(12 * multRecompense),
				Sort:      combat.SortMonstre{Nom: "Morsure Féroce", Degats: 15, CoutMana: 0},
				TableDrop: []string{"Fourrure de loup", "Potion de Force (Buff Atk)"}, ChancDrop: 25 + (p.NGPlus * 10),
			},
		}
		// monstre au hasard dans la liste de la zone
		m = &mobs[rand.Intn(len(mobs))]
	} else if zone == 2 {
		mobs := []combat.Monstre{
			{
				Nom: "Troll des Cavernes",
				PV:  85, MaxPV: 85, Mana: 30, MaxMana: 30, Attaque: 16, Defense: 6, Dexterite: 5,
				GainXP: int(75 * multRecompense), GainOr: int(22 * multRecompense),
				Sort:      combat.SortMonstre{Nom: "Coup de Clavaire Brutal", Degats: 22, CoutMana: 15},
				TableDrop: []string{"Peau de troll", "Potion d'Armure (Buff Def)"}, ChancDrop: 35 + (p.NGPlus * 10),
			},
			{
				Nom: "Sorcier Squelette",
				PV:  65, MaxPV: 65, Mana: 60, MaxMana: 60, Attaque: 10, Defense: 4, Dexterite: 11,
				GainXP: int(85 * multRecompense), GainOr: int(30 * multRecompense),
				Sort:      combat.SortMonstre{Nom: "Orbe des Ténèbres", Degats: 28, CoutMana: 20},
				TableDrop: []string{"Lingot de Fer", "Potion de mana"}, ChancDrop: 40 + (p.NGPlus * 10),
			},
		}
		m = &mobs[rand.Intn(len(mobs))]
	} else {
		mobs := []combat.Monstre{
			{
				Nom: "Drake Infernale",
				PV:  130, MaxPV: 130, Mana: 50, MaxMana: 50, Attaque: 22, Defense: 9, Dexterite: 10,
				GainXP: int(150 * multRecompense), GainOr: int(60 * multRecompense),
				Sort:      combat.SortMonstre{Nom: "Souffle de Flammes", Degats: 35, CoutMana: 25},
				TableDrop: []string{"Minerai de Mithril", "Grande Potion de soin"}, ChancDrop: 45 + (p.NGPlus * 10),
			},
		}
		m = &mobs[rand.Intn(len(mobs))]
	}

	// Application du scaling NG+ / Niveau aux stats du monstre
	m.MaxPV = int(float64(m.MaxPV) * multStat)
	m.PV = m.MaxPV
	m.Attaque = int(float64(m.Attaque) * multStat)
	m.Defense = int(float64(m.Defense) * multStat)

	combat.LancerCombat(p, m)
}

// le boss scale plus vite que les mobs normaux (0.7 par NG+ contre 0.5)
func CombattreBossFinal(p *character.Personnage) {
	multStat := 1.0 + (float64(p.NGPlus) * 0.7) + (float64(p.Niveau-1) * 0.15)
	multRecompense := 1.0 + (float64(p.NGPlus) * 1.0)

	boss := &combat.Monstre{
		Nom:       fmt.Sprintf("Le Dieu Démon (NG+%d)", p.NGPlus),
		MaxPV:     int(220 * multStat),
		PV:        int(220 * multStat),
		MaxMana:   100,
		Mana:      100,
		Attaque:   int(26 * multStat),
		Defense:   int(8 * multStat),
		Dexterite: 12,
		GainXP:    int(500 * multRecompense),
		GainOr:    int(250 * multRecompense),
		Sort:      combat.SortMonstre{Nom: "Apocalypse Obscure", Degats: int(40 * multStat), CoutMana: 30},
		TableDrop: []string{"Minerai de Mithril", "Potion d'Elixir (Vie+Mana)"},
		// le boss drop à tous les coups
		ChancDrop: 100,
	}

	victoire := combat.LancerCombat(p, boss)

	// Passage obligatoire en NEW GAME+ (NG+) en cas de victoire contre le Boss
	if victoire {
		p.NGPlus++
		ui.ClearConsole()
		ui.AfficherTitre("🎉 VICTOIRE FINALE - NEW GAME+ DÉBLOQUÉ !")
		fmt.Printf(ui.Jaune+"Vous avez vaincu le Dieu Démon ! Le monde bascule en NEW GAME+ (NG+%d).\n"+ui.Reset, p.NGPlus)
		fmt.Println("----------------------------------------------------------------------")
		fmt.Println("|!| La difficulté générale des monstres augmente fortement |!|")
		fmt.Println("!!! Les récompenses d'XP, d'or et le taux de drop rare sont drastiquement augmentés !!!")
		fmt.Println(" Vous devez à nouveau explorer les zones depuis la première.")
		fmt.Println("----------------------------------------------------------------------")
		ui.LireEntree("Appuyez sur Entrée pour continuer votre légende...")
	}
}
