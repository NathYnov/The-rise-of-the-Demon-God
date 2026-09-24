package economy

import (
	"fmt"
	"projet-red/pkg/character"
	"projet-red/pkg/ui"
)

// menu principal de la boutique
func OuvrirMarchand(p *character.Personnage) {
	for {
		ui.ClearConsole()
		ui.DisplayMerchantArt()
		ui.AfficherTitre("Boutique du Marchand")
		fmt.Printf(ui.Jaune+"Votre bourse : %d PO | Niveau : %d (Objets adaptés à votre niveau)\n\n"+ui.Reset, p.Or, p.Niveau)

		fmt.Println("--- CATÉGORIES ---")
		fmt.Println("1. Potions & Consommables")
		fmt.Println("2. Matériaux de Fabrication")
		fmt.Println("3. Livres de Sorts")
		fmt.Println("4. Extension d'Inventaire (+10 places, 30 PO)")
		fmt.Println("5. Vendre des ressources")
		fmt.Println("0. Quitter")

		choix := ui.LireEntree("\nFaites votre choix : ")

		switch choix {
		case "1":
			MenuPotions(p)
		case "2":
			MenuMateriaux(p)
		case "3":
			MenuLivresSorts(p)
		case "4":
			p.AugmenterInventaire()
			ui.LireEntree("Appuyez sur Entrée...")
		case "5":
			VendreRessources(p)
		case "0":
			return
		}
	}
}

// attention les prix et niveaux sont écrits deux fois (affichage + appel), à changer aux deux endroits
func MenuPotions(p *character.Personnage) {
	ui.ClearConsole()
	ui.AfficherTitre("Marchand - Potions")

	fmt.Println("1. Potion de soin (+35 PV) - 4 PO (Niv 1+)")
	fmt.Println("2. Potion de mana (+35 Mana) - 5 PO (Niv 1+)")
	fmt.Println("3. Potion de poison (Attaque ennemie) - 8 PO (Niv 1+)")
	fmt.Println("4. Potion de Force (+5 ATK combat) - 25 PO (Niv 2+)")
	fmt.Println("5. Potion d'Armure (+5 DEF combat) - 25 PO (Niv 2+)")
	fmt.Println("6. Grande Potion de soin (+100 PV) - 50 PO (Niv 3+)")
	fmt.Println("7. Potion d'Elixir (Vie + Mana) - 100 PO (Niv 4+)")
	fmt.Println("0. Retour")

	choix := ui.LireEntree("\nVotre choix : ")
	switch choix {
	case "1":
		AcheterObjet(p, "Potion de soin", 4, 1)
	case "2":
		AcheterObjet(p, "Potion de mana", 5, 1)
	case "3":
		AcheterObjet(p, "Potion de poison", 8, 1)
	case "4":
		AcheterObjet(p, "Potion de Force (Buff Atk)", 25, 2)
	case "5":
		AcheterObjet(p, "Potion d'Armure (Buff Def)", 25, 2)
	case "6":
		AcheterObjet(p, "Grande Potion de soin", 50, 3)
	case "7":
		AcheterObjet(p, "Potion d'Elixir (Vie+Mana)", 100, 4)
	}
}

func MenuMateriaux(p *character.Personnage) {
	ui.ClearConsole()
	ui.AfficherTitre("Marchand - Matériaux de Craft")

	fmt.Println("[Niveau 1]")
	fmt.Println("1. Plume de corbeau - 5 PO")
	fmt.Println("2. Cuir de sanglier - 7 PO")

	fmt.Println("\n[Niveau 2+]")
	fmt.Println("3. Fourrure de loup - 10 PO")
	fmt.Println("4. Lingot de Fer - 15 PO")

	fmt.Println("\n[Niveau 3+]")
	fmt.Println("5. Peau de troll - 20 PO")
	fmt.Println("6. Minerai de Mithril - 100 PO")
	fmt.Println("0. Retour")

	choix := ui.LireEntree("\nVotre choix : ")
	switch choix {
	case "1":
		AcheterObjet(p, "Plume de corbeau", 5, 1)
	case "2":
		AcheterObjet(p, "Cuir de sanglier", 7, 1)
	case "3":
		AcheterObjet(p, "Fourrure de loup", 10, 2)
	case "4":
		AcheterObjet(p, "Lingot de Fer", 15, 2)
	case "5":
		AcheterObjet(p, "Peau de troll", 20, 3)
	case "6":
		AcheterObjet(p, "Minerai de Mithril", 100, 3)
	}
}

func MenuLivresSorts(p *character.Personnage) {
	ui.ClearConsole()
	ui.AfficherTitre("Marchand - Livres de Sorts")

	fmt.Println("1. Livre de sort: Boule de Feu (100 PO) (Niv 1+)")
	fmt.Println("2. Livre de sort: Soin Magique (500 PO) (Niv 2+)")
	fmt.Println("3. Livre de sort: Éclair Céleste (1000 PO) (Niv 4+)")
	fmt.Println("0. Retour")

	choix := ui.LireEntree("\nVotre choix : ")
	switch choix {
	case "1":
		AcheterObjet(p, "Livre de sort: Boule de Feu", 100, 1)
	case "2":
		AcheterObjet(p, "Livre de sort: Soin Magique", 500, 2)
	case "3":
		AcheterObjet(p, "Livre de sort: Éclair Céleste", 1000, 4)
	}
}

// vérifie le niveau, la place dans l'inventaire puis l'or
func AcheterObjet(p *character.Personnage, nom string, prix int, nivRequis int) {
	if p.Niveau < nivRequis {
		fmt.Printf(ui.Rouge+"Niveau insuffisant ! (Niveau %d requis)\n"+ui.Reset, nivRequis)
		ui.LireEntree("Appuyez sur Entrée...")
		return
	}
	if len(p.Inventaire) >= p.MaxInventaire {
		fmt.Println(ui.Rouge+"Votre inventaire est plein !"+ui.Reset)
		ui.LireEntree("Appuyez sur Entrée...")
		return
	}
	if p.Or >= prix {
		p.Or -= prix
		p.Inventaire = append(p.Inventaire, nom)
		fmt.Printf(ui.Vert+"Achat réussi : %s pour %d PO.\n"+ui.Reset, nom, prix)
	} else {
		fmt.Println(ui.Rouge+"Vous n'avez pas assez d'or."+ui.Reset)
	}
	ui.LireEntree("Appuyez sur Entrée...")
}

// on ne peut revendre que les matériaux
func VendreRessources(p *character.Personnage) {
	ui.ClearConsole()
	ui.AfficherTitre("Revente de matériaux")

	// prix de revente, à peu près 50% du prix d'achat 
	prixItems := map[string]int{
		"Plume de corbeau":   3,
		"Cuir de sanglier":   4,
		"Fourrure de loup":   6,
		"Lingot de Fer":      8,
		"Peau de troll":      13,
		"Minerai de Mithril": 50,
	}

	if len(p.Inventaire) == 0 {
		fmt.Println("Rien à vendre.")
		ui.LireEntree("Appuyez sur Entrée...")
		return
	}

	for i, item := range p.Inventaire {
		val, existe := prixItems[item]
		if existe {
			fmt.Printf("%d. %s (Revente : %d PO)\n", i+1, item, val)
		} else {
			fmt.Printf("%d. %s (Non racheté)\n", i+1, item)
		}
	}

	idxStr := ui.LireEntree("\nNuméro de l'objet à vendre (0 pour annuler) : ")
	var idx int
	fmt.Sscanf(idxStr, "%d", &idx)

	if idx > 0 && idx <= len(p.Inventaire) {
		item := p.Inventaire[idx-1]
		valeur, existe := prixItems[item]
		if existe {
			p.Or += valeur
			p.RetirerObjetIndex(idx - 1)
			fmt.Printf(ui.Vert+"Vendu : %s pour %d PO !\n"+ui.Reset, item, valeur)
		} else {
			fmt.Println(ui.Rouge+"Le marchand ne veut pas cet objet."+ui.Reset)
		}
	}
	ui.LireEntree("Appuyez sur Entrée...")
}

// FORGERON STRUCTURÉ PAR NIVEAUX ET PAR TYPE
func OuvrirForgeron(p *character.Personnage) {
	for {
		ui.ClearConsole()
		ui.DisplayBlacksmithArt()
		ui.AfficherTitre("La Forge du Village")
		fmt.Printf(ui.Jaune+"Votre Or : %d PO | Les stats des objets augmentent selon votre niveau !\n\n"+ui.Reset, p.Or)

		fmt.Println("1. Équipement Tête (Casques / Chapeaux)")
		fmt.Println("2. Équipement Torse (Plastrons / Tuniques)")
		fmt.Println("3. Armes (Épées / Bâtons)")
		fmt.Println("0. Quitter")

		choix := ui.LireEntree("\nFaites votre choix : ")
		switch choix {
		case "1":
			MenuCraftTete(p)
		case "2":
			MenuCraftTorse(p)
		case "3":
			MenuCraftArmes(p)
		case "0":
			return
		}
	}
}

func MenuCraftTete(p *character.Personnage) {
	ui.ClearConsole()
	ui.AfficherTitre("Forge - Équipement Tête")

	// le bonus dépend du niveau au moment du craft
	bonusNiv := p.Niveau * 2

	fmt.Printf("1. [Niv 1] Casque en Cuir (10x Cuir de sanglier, 100 PO) -> +%d PV\n", 20+bonusNiv)
	fmt.Printf("2. [Niv 3] Casque en Fer (10x Lingot de Fer, 250 PO) -> +%d PV, +%d Def\n", 50+bonusNiv, 3+p.Niveau)
	fmt.Println("0. Retour")

	choix := ui.LireEntree("\nChoix : ")
	switch choix {
	case "1":
		CraftObjet(p, fmt.Sprintf("Casque en Cuir (Niv %d)", p.Niveau), map[string]int{"Cuir de sanglier": 10}, 100, "Tete", 20+bonusNiv, 0, 0)
	case "2":
		if p.Niveau < 3 {
			fmt.Println("Niveau 3 requis !")
			ui.LireEntree("Appuyez sur Entrée...")
			return
		}
		CraftObjet(p, fmt.Sprintf("Casque en Fer (Niv %d)", p.Niveau), map[string]int{"Lingot de Fer": 10}, 250, "Tete", 50+bonusNiv, 0, 3+p.Niveau)
	}
}

func MenuCraftTorse(p *character.Personnage) {
	ui.ClearConsole()
	ui.AfficherTitre("Forge - Équipement Torse")

	bonusNiv := p.Niveau * 3

	fmt.Printf("1. [Niv 1] Tunique Légère (10x Plume de corbeau, 5x Cuir de sanglier, 150 PO) -> +%d PV, +%d Def\n", 25+bonusNiv, 2)
	fmt.Printf("2. [Niv 2] Plastron en Peau (10x Peau de troll, 400 PO) -> +%d PV, +%d Def\n", 50+bonusNiv, 5+p.Niveau)
	fmt.Printf("3. [Niv 4] Armure en Mithril (10x Minerai de Mithril, 1000 PO) -> +%d PV, +%d Def\n", 90+bonusNiv, 10+p.Niveau)
	fmt.Println("0. Retour")

	choix := ui.LireEntree("\nChoix : ")
	switch choix {
	case "1":
		CraftObjet(p, fmt.Sprintf("Tunique Légère (Niv %d)", p.Niveau), map[string]int{"Plume de corbeau": 10, "Cuir de sanglier": 5}, 150, "Torse", 25+bonusNiv, 0, 2)
	case "2":
		if p.Niveau < 2 {
			fmt.Println("Niveau 2 requis !")
			ui.LireEntree("Appuyez sur Entrée...")
			return
		}
		CraftObjet(p, fmt.Sprintf("Plastron en Peau (Niv %d)", p.Niveau), map[string]int{"Peau de troll": 10}, 400, "Torse", 50+bonusNiv, 0, 5+p.Niveau)
	case "3":
		if p.Niveau < 4 {
			fmt.Println("Niveau 4 requis !")
			ui.LireEntree("Appuyez sur Entrée...")
			return
		}
		CraftObjet(p, fmt.Sprintf("Armure en Mithril (Niv %d)", p.Niveau), map[string]int{"Minerai de Mithril": 10}, 1000, "Torse", 90+bonusNiv, 0, 10+p.Niveau)
	}
}

func MenuCraftArmes(p *character.Personnage) {
	ui.ClearConsole()
	ui.AfficherTitre("Forge - Armes")

	bonusAtk := p.Niveau * 2

	fmt.Printf("1. [Niv 1] Épée de Bois Renforcée (10x Fourrure de loup, 150 PO) -> +%d Atk\n", 8+bonusAtk)
	fmt.Printf("2. [Niv 2] Épée en Fer (10x Lingot de Fer, 400 PO) -> +%d Atk\n", 20+bonusAtk)
	fmt.Printf("3. [Niv 4] Lame de Mithril (10x Minerai de Mithril, 1000 PO) -> +%d Atk\n", 50+bonusAtk)
	fmt.Println("0. Retour")

	choix := ui.LireEntree("\nChoix : ")
	switch choix {
	case "1":
		CraftObjet(p, fmt.Sprintf("Épée de Bois (Niv %d)", p.Niveau), map[string]int{"Fourrure de loup": 10}, 150, "Arme", 0, 8+bonusAtk, 0)
	case "2":
		if p.Niveau < 2 {
			fmt.Println("Niveau 2 requis !")
			ui.LireEntree("Appuyez sur Entrée...")
			return
		}
		CraftObjet(p, fmt.Sprintf("Épée en Fer (Niv %d)", p.Niveau), map[string]int{"Lingot de Fer": 10}, 400, "Arme", 0, 20+bonusAtk, 0)
	case "3":
		if p.Niveau < 4 {
			fmt.Println("Niveau 4 requis !")
			ui.LireEntree("Appuyez sur Entrée...")
			return
		}
		CraftObjet(p, fmt.Sprintf("Lame de Mithril (Niv %d)", p.Niveau), map[string]int{"Minerai de Mithril": 10}, 1000, "Arme", 0, 50+bonusAtk, 0)
	}
}

// fabrique l'objet : contrôle de l'or et des composants, on retire tout puis on applique les bonus
func CraftObjet(p *character.Personnage, itemCraft string, composants map[string]int, coutOr int, slot string, bonusPV int, bonusAtk int, bonusDef int) {
	// Vérification de l'or
	if p.Or < coutOr {
		fmt.Println(ui.Rouge+"Pas assez d'or !"+ui.Reset)
		ui.LireEntree("Appuyez sur Entrée...")
		return
	}

	// Vérification des composants
	for comp, qte := range composants {
		if p.CompterObjet(comp) < qte {
			fmt.Printf(ui.Rouge+"Composants insuffisants ! Requis: %d %s\n"+ui.Reset, qte, comp)
			ui.LireEntree("Appuyez sur Entrée...")
			return
		}
	}

	// Retrait composants et or
	p.Or -= coutOr
	for comp, qte := range composants {
		for i := 0; i < qte; i++ {
			p.RetirerObjetNom(comp)
		}
	}

	// Application des bonus
	p.MaxPV += bonusPV
	p.PV += bonusPV
	p.Attaque += bonusAtk
	p.Defense += bonusDef

	// attention on remplace juste le nom, les bonus de l'ancien équipement restent
	if slot == "Tete" {
		p.Stuff.Tete = itemCraft
	} else if slot == "Torse" {
		p.Stuff.Torse = itemCraft
	} else if slot == "Arme" {
		p.Stuff.Arme = itemCraft
	}

	fmt.Printf(ui.Vert+"🛠️ Objet fabriqué avec succès : %s !\n"+ui.Reset, itemCraft)
	fmt.Printf("Bonus appliqués -> PV: +%d | ATK: +%d | DEF: +%d\n", bonusPV, bonusAtk, bonusDef)
	ui.LireEntree("Appuyez sur Entrée...")
}
