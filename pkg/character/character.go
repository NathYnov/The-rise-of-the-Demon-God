package character

import (
	"fmt"
	"projet-red/pkg/ui"
	"strings"
)

// les 3 emplacements d'équipement
type Equipement struct {
	Tete  string
	Torse string
	Arme  string
}

// un sort du joueur, soit des dégâts soit du soin
type Sort struct {
	Nom      string
	Degats   int
	Soin     int
	CoutMana int
}

// tout ce qui est sauvegardé dans le json passe par cette struct
type Personnage struct {
	Nom             string
	Race            string
	Niveau          int
	XP              int
	MaxXP           int
	PV              int
	MaxPV           int
	Mana            int
	MaxMana         int
	Attaque         int
	Defense         int
	Dexterite       int
	Or              int
	Inventaire      []string
	MaxInventaire   int
	NbAmeliorations int
	Stuff           Equipement
	Sorts           []Sort
	BuffAttaque     int // Bonus temporaire d'attaque
	BuffDefense     int // Bonus temporaire de défense
	NGPlus          int // Niveau de New Game+ (0 = normal, 1 = NG+1, etc.)
}

// première lettre en majuscule, le reste en minuscule
func FormaterNom(s string) string {
	// nom vide donc nom par défaut
	if len(s) == 0 {
		return "Héros"
	}
	s = strings.ToLower(s)
	return strings.ToUpper(string(s[0])) + s[1:]
}

func CreationPersonnage() *Personnage {
	ui.ClearConsole()
	ui.AfficherTitre("Création de Personnage")

	nom := ui.LireEntree("Entrez le nom de votre héros : ")
	nom = FormaterNom(nom)

	perso := &Personnage{
		Nom:             nom,
		Niveau:          1,
		XP:              0,
		MaxXP:           100,
		Or:              100, // 100 PO de départ
		MaxInventaire:   10,
		NbAmeliorations: 0,
		Inventaire:      []string{"Potion de soin"},
		Sorts:           []Sort{},
		NGPlus:          0,
	}

	fmt.Println("\nChoisissez votre race :")
	fmt.Println("1. Humain (PV: 100, Mana: 50, Attaque: 12, Def: 5, Dex: 10)")
	fmt.Println("2. Elfe (PV: 80, Mana: 100, Attaque: 8, Def: 3, Dex: 15)")
	fmt.Println("3. Nain (PV: 130, Mana: 30, Attaque: 14, Def: 8, Dex: 6)")
	fmt.Println("4. Mage (PV: 70, Mana: 120, Attaque: 6, Def: 2, Dex: 12)")
	fmt.Println("5. Dragonoïde (PV: 140, Mana: 40, Attaque: 16, Def: 7, Dex: 5)")

	choixRace := ui.LireEntree("Choix (1-5) : ")
	// si le joueur tape autre chose que 2 à 5 on part sur Humain (default)
	switch choixRace {
	case "2":
		perso.Race = "Elfe"
		perso.MaxPV, perso.PV = 80, 80
		perso.MaxMana, perso.Mana = 100, 100
		perso.Attaque, perso.Defense, perso.Dexterite = 8, 3, 15
	case "3":
		perso.Race = "Nain"
		perso.MaxPV, perso.PV = 130, 130
		perso.MaxMana, perso.Mana = 30, 30
		perso.Attaque, perso.Defense, perso.Dexterite = 14, 8, 6
	case "4":
		perso.Race = "Mage"
		perso.MaxPV, perso.PV = 70, 70
		perso.MaxMana, perso.Mana = 120, 120
		perso.Attaque, perso.Defense, perso.Dexterite = 6, 2, 12
	case "5":
		perso.Race = "Dragonoïde"
		perso.MaxPV, perso.PV = 140, 140
		perso.MaxMana, perso.Mana = 40, 40
		perso.Attaque, perso.Defense, perso.Dexterite = 16, 7, 5
	default:
		perso.Race = "Humain"
		perso.MaxPV, perso.PV = 100, 100
		perso.MaxMana, perso.Mana = 50, 50
		perso.Attaque, perso.Defense, perso.Dexterite = 12, 5, 10
	}

	// tout le monde commence avec coup de poing
	perso.Sorts = append(perso.Sorts, Sort{
		Nom:      "Coup de Poing",
		Degats:   10,
		CoutMana: 5,
	})

	return perso
}

func (p *Personnage) AfficherFiche() {
	ui.ClearConsole()
	ui.AfficherTitre("Fiche de Personnage")
	fmt.Printf("Nom : %s | Race : %s | Mode : NG+%d\n", p.Nom, p.Race, p.NGPlus)
	fmt.Printf("Niveau : %d | XP : %d/%d\n", p.Niveau, p.XP, p.MaxXP)
	fmt.Printf("PV : %d/%d | Mana : %d/%d\n", p.PV, p.MaxPV, p.Mana, p.MaxMana)
	fmt.Printf("Attaque : %d (+%d buff) | Défense : %d (+%d buff) | Dextérité : %d\n",
		p.Attaque, p.BuffAttaque, p.Defense, p.BuffDefense, p.Dexterite)
	fmt.Printf("Or : %d PO\n", p.Or)
	fmt.Printf("Équipement -> Tête: %s | Torse: %s | Arme: %s\n", p.Stuff.Tete, p.Stuff.Torse, p.Stuff.Arme)
	fmt.Printf("Inventaire : %d/%d objets\n", len(p.Inventaire), p.MaxInventaire)

	ui.LireEntree("\nAppuyez sur Entrée pour revenir au menu...")
}

// gère le gain d'xp, on peut monter plusieurs niveaux d'un coup
func (p *Personnage) AjouterXP(gain int) {
	p.XP += gain
	fmt.Printf(ui.Cyan+"Vous gagnez %d d'XP !\n"+ui.Reset, gain)

	for p.XP >= p.MaxXP {
		p.XP -= p.MaxXP
		p.Niveau++
		// il faut 40% d'xp en plus à chaque niveau
		p.MaxXP = int(float64(p.MaxXP) * 1.4)
		// gains fixes par niveau, pv et mana remis au max
		p.MaxPV += 20
		p.PV = p.MaxPV
		p.MaxMana += 10
		p.Mana = p.MaxMana
		p.Attaque += 3
		p.Defense += 2
		p.Dexterite += 2
		fmt.Printf(ui.Jaune+"🎉 NIVEAU SUPÉRIEUR ! Vous êtes maintenant niveau %d !\n"+ui.Reset, p.Niveau)
	}
}

// pas de game over, le joueur revient avec la moitié de sa vie
func (p *Personnage) VerifierMort() bool {
	if p.PV <= 0 {
		fmt.Println(ui.Rouge + "\n💀 Vous avez été vaincu..." + ui.Reset)
		p.PV = p.MaxPV / 2
		p.BuffAttaque = 0
		p.BuffDefense = 0
		fmt.Printf(ui.Vert+"Un miracle produit son effet : vous ressuscitez avec %d PV !\n"+ui.Reset, p.PV)
		ui.LireEntree("Appuyez sur Entrée...")
		return true
	}
	return false
}

func (p *Personnage) OuvrirInventaire() {
	ui.ClearConsole()
	ui.AfficherTitre("Inventaire")

	if len(p.Inventaire) == 0 {
		fmt.Println("Votre inventaire est vide.")
		ui.LireEntree("\nAppuyez sur Entrée...")
		return
	}

	for i, item := range p.Inventaire {
		fmt.Printf("%d. %s\n", i+1, item)
	}

	fmt.Println("\nActions :")
	fmt.Println("1. Utiliser un objet / Consommable / Livre")
	fmt.Println("2. Quitter")
	choix := ui.LireEntree("Votre choix : ")

	if choix == "1" {
		idxStr := ui.LireEntree("Numéro de l'objet : ")
		var idx int
		fmt.Sscanf(idxStr, "%d", &idx)
		if idx > 0 && idx <= len(p.Inventaire) {
			item := p.Inventaire[idx-1]
			p.UtiliserObjet(item, idx-1)
		}
	}
}

// applique l'effet de l'objet, index sert à le retirer ensuite
func (p *Personnage) UtiliserObjet(nom string, index int) {
	switch nom {
	case "Potion de soin":
		p.PV += 35
		if p.PV > p.MaxPV {
			p.PV = p.MaxPV
		}
		fmt.Println(ui.Vert + "Vous buvez une Potion de soin (+35 PV)." + ui.Reset)
		p.RetirerObjetIndex(index)
	case "Grande Potion de soin":
		p.PV += 80
		if p.PV > p.MaxPV {
			p.PV = p.MaxPV
		}
		fmt.Println(ui.Vert + "Vous buvez une Grande Potion de soin (+80 PV)." + ui.Reset)
		p.RetirerObjetIndex(index)
	case "Potion de mana":
		p.Mana += 35
		if p.Mana > p.MaxMana {
			p.Mana = p.MaxMana
		}
		fmt.Println(ui.Vert + "Vous buvez une Potion de mana (+35 Mana)." + ui.Reset)
		p.RetirerObjetIndex(index)
	case "Potion d'Elixir (Vie+Mana)":
		p.PV += 60
		p.Mana += 40
		if p.PV > p.MaxPV {
			p.PV = p.MaxPV
		}
		if p.Mana > p.MaxMana {
			p.Mana = p.MaxMana
		}
		fmt.Println(ui.Vert + "Vous buvez un Elixir (+60 PV / +40 Mana)." + ui.Reset)
		p.RetirerObjetIndex(index)
	// les buffs ne durent que pour le prochain combat
	case "Potion de Force (Buff Atk)":
		p.BuffAttaque += 5
		fmt.Println(ui.Vert + "Vous buvez une Potion de Force (+5 ATK pour le prochain combat)." + ui.Reset)
		p.RetirerObjetIndex(index)
	case "Potion d'Armure (Buff Def)":
		p.BuffDefense += 5
		fmt.Println(ui.Vert + "Vous buvez une Potion d'Armure (+5 DEF pour le prochain combat)." + ui.Reset)
		p.RetirerObjetIndex(index)
	// les livres apprennent un sort puis disparaissent
	case "Livre de sort: Boule de Feu":
		p.AjouterSort(Sort{Nom: "Boule de Feu", Degats: 50, CoutMana: 50})
		p.RetirerObjetIndex(index)
	case "Livre de sort: Soin Magique":
		p.AjouterSort(Sort{Nom: "Soin Magique", Soin: 100, CoutMana: 100})
		p.RetirerObjetIndex(index)
	case "Livre de sort: Éclair Céleste":
		p.AjouterSort(Sort{Nom: "Éclair Céleste", Degats: 500, CoutMana: 500})
		p.RetirerObjetIndex(index)
	default:
		fmt.Println("Cet objet ne peut pas être consommé ici ou est réservé au combat / forgeron.")
	}
	ui.LireEntree("Appuyez sur Entrée...")
}

// on vérifie qu'on connaît pas déjà le sort
func (p *Personnage) AjouterSort(s Sort) {
	for _, sort := range p.Sorts {
		if sort.Nom == s.Nom {
			fmt.Println(ui.Rouge + "Vous connaissez déjà ce sort !" + ui.Reset)
			return
		}
	}
	p.Sorts = append(p.Sorts, s)
	fmt.Printf(ui.Vert+"Vous avez appris le sort : %s !\n"+ui.Reset, s.Nom)
}

// enlève l'élément i du slice
func (p *Personnage) RetirerObjetIndex(i int) {
	if i >= 0 && i < len(p.Inventaire) {
		p.Inventaire = append(p.Inventaire[:i], p.Inventaire[i+1:]...)
	}
}

// retire la première occurrence, false si l'objet n'y est pas
func (p *Personnage) RetirerObjetNom(nom string) bool {
	for i, item := range p.Inventaire {
		if item == nom {
			p.RetirerObjetIndex(i)
			return true
		}
	}
	return false
}

// nombre d'exemplaires d'un objet dans l'inventaire
func (p *Personnage) CompterObjet(nom string) int {
	c := 0
	for _, item := range p.Inventaire {
		if item == nom {
			c++
		}
	}
	return c
}

// 30 PO pour +10 places, 3 achats max
func (p *Personnage) AugmenterInventaire() {
	if p.NbAmeliorations >= 3 {
		fmt.Println(ui.Rouge + "Capacité d'inventaire déjà au maximum !" + ui.Reset)
		return
	}
	if p.Or >= 30 {
		p.Or -= 30
		p.MaxInventaire += 10
		p.NbAmeliorations++
		fmt.Printf(ui.Vert+"Inventaire agrandi ! Capacité actuelle : %d objets.\n"+ui.Reset, p.MaxInventaire)
	} else {
		fmt.Println(ui.Rouge + "Pas assez d'or (30 PO requis)." + ui.Reset)
	}
}
