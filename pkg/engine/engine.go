package engine

import (
	"encoding/json"
	"fmt"
	"os"

	"projet-red/pkg/adventure"
	"projet-red/pkg/character"
	"projet-red/pkg/combat"
	"projet-red/pkg/economy"
	"projet-red/pkg/ui"
)

// écrit le perso dans saves/<nom>.json
func Sauvegarder(p *character.Personnage) {
	// on crée le dossier au cas où il n'existe pas encore
	os.MkdirAll("saves", 0755)

	fichier := fmt.Sprintf("saves/%s.json", p.Nom)
	// indent de 2 espaces pour que le json reste lisible
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		fmt.Println("Erreur lors de la sauvegarde :", err)
		return
	}

	err = os.WriteFile(fichier, data, 0644)
	if err == nil {
		fmt.Println(ui.Vert+"Partie sauvegardée dans"+ui.Reset, fichier)
	}
}

// affiche les fichiers de saves/ et charge celui choisi (nil si ça rate)
func ChargerPartie() *character.Personnage {
	ui.ClearConsole()
	ui.AfficherTitre("Charger une partie")

	fichiers, err := os.ReadDir("saves")
	if err != nil || len(fichiers) == 0 {
		fmt.Println(ui.Rouge+"Aucune sauvegarde trouvée."+ui.Reset)
		ui.LireEntree("Appuyez sur Entrée...")
		return nil
	}

	for i, f := range fichiers {
		fmt.Printf("%d. %s\n", i+1, f.Name())
	}

	choixStr := ui.LireEntree("\nNuméro de la sauvegarde : ")
	// si ce n'est pas un nombre idx reste à 0 donc on ressort plus bas
	var idx int
	fmt.Sscanf(choixStr, "%d", &idx)

	if idx > 0 && idx <= len(fichiers) {
		nomFichier := "saves/" + fichiers[idx-1].Name()
		data, err := os.ReadFile(nomFichier)
		if err != nil {
			fmt.Println("Erreur lors de la lecture du fichier.")
			return nil
		}

		// le json remplit directement p
		var p character.Personnage
		err = json.Unmarshal(data, &p)
		if err != nil {
			fmt.Println("Fichier de sauvegarde corrompu.")
			return nil
		}

		fmt.Println(ui.Vert+"Sauvegarde chargée avec succès !"+ui.Reset)
		ui.LireEntree("Appuyez sur Entrée...")
		return &p
	}
	return nil
}

// menu du village, le joueur revient toujours ici entre deux actions
func BouclePrincipale(p *character.Personnage) {
	for {
		ui.ClearConsole()
		ui.AfficherBanniere()
		fmt.Printf(ui.Cyan+" Joueur : %s | Niv : %d (NG+%d) | PV : %d/%d | Or : %d PO\n\n"+ui.Reset,
			p.Nom, p.Niveau, p.NGPlus, p.PV, p.MaxPV, p.Or)

		fmt.Println("1. Info Personnage")
		fmt.Println("2. Inventaire & Consommables")
		fmt.Println("3. Aller chez le Marchand")
		fmt.Println("4. Aller chez le Forgeron")
		fmt.Println("5. Partir en Aventure")
		fmt.Println("6. Mannequin d'Entraînement")
		fmt.Println("7. Sauvegarder et Quitter")

		choix := ui.LireEntree("\nChoix : ")

		switch choix {
		case "1":
			p.AfficherFiche()
		case "2":
			p.OuvrirInventaire()
		case "3":
			economy.OuvrirMarchand(p)
		case "4":
			economy.OuvrirForgeron(p)
		case "5":
			adventure.Explorer(p)
		case "6":
			// monstre bidon pour tester les combats sans risque
			entrainement := &combat.Monstre{
				Nom: "Mannequin de Paille",
				PV:  50, MaxPV: 50, Mana: 0, MaxMana: 0, Attaque: 2, Defense: 0, Dexterite: 1,
				GainXP: 15, GainOr: 5,
				Sort:      combat.SortMonstre{Nom: "Rien", Degats: 0, CoutMana: 0},
				TableDrop: []string{"Potion de soin"}, ChancDrop: 20,
			}
			combat.LancerCombat(p, entrainement)
		case "7":
			Sauvegarder(p)
			fmt.Println(ui.Cyan+"Au revoir !"+ui.Reset)
			return
		}
	}
}
