package main

import (
	"fmt"

	"projet-red/pkg/character"
	"projet-red/pkg/engine"
	"projet-red/pkg/ui"
)

func main() {
	// menu principal, on reboucle tant que le joueur ne quitte pas
	for {
		ui.ClearConsole()
		ui.AfficherBanniere()

		fmt.Println(ui.Vert + "1. Nouvelle Partie" + ui.Reset)
		fmt.Println(ui.Vert + "2. Charger une Partie" + ui.Reset)
		fmt.Println(ui.Violet + "3. secret" + ui.Reset)
		fmt.Println(ui.Rouge + "4. Quitter" + ui.Reset)

		choix := ui.LireEntree("\nFaites votre choix : ")

		switch choix {
		case "1":
			// nouvelle partie : on crée le perso puis on lance le jeu
			perso := character.CreationPersonnage()
			engine.BouclePrincipale(perso)
		case "2":
			// ChargerPartie renvoie nil si rien n'a été chargé
			perso := engine.ChargerPartie()
			if perso != nil {
				engine.BouclePrincipale(perso)
			}
		case "3":
			// easter egg avec une phrase cachée
			secret := ui.LireEntree("Entrez la phrase cachée : ")
			secret = fmt.Sprintln(secret)
			if secret == "Qui sont-ils" {
				fmt.Println("Nathan, Soumaila, Shilsa !!!")
				return
			}
			fmt.Println("Nathan, Soumaila, Shilsa !!!")
			ui.LireEntree("\nAppuyez sur Entrée pour revenir au menu...")
			continue
		case "4":
			fmt.Println(ui.Cyan + "Merci d'avoir joué !" + ui.Reset)
			return
		}
	}
}
