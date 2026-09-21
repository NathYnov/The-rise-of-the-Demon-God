package combat

import (
	"fmt"
	"math/rand"
	"projet-red/pkg/character"
	"projet-red/pkg/ui"
)

// sort du monstre, il n'a pas de soin contrairement au joueur
type SortMonstre struct {
	Nom      string
	Degats   int
	CoutMana int
}

type Monstre struct {
	Nom       string
	PV        int
	MaxPV     int
	Mana      int
	MaxMana   int
	Attaque   int
	Defense   int
	Dexterite int
	GainXP    int
	GainOr    int
	Sort      SortMonstre
	TableDrop []string // Objets rares pouvant être droppés
	ChancDrop int      // Pourcentage de chance (ex: 25%)
}

// combat au tour par tour, renvoie true si le joueur gagne
func LancerCombat(p *character.Personnage, m *Monstre) bool {
	ui.ClearConsole()
	ui.AfficherTitre("COMBAT : " + p.Nom + " vs " + m.Nom)
	ui.DisplayMonsterArt(m.Nom)

	// la plus grosse dextérité commence, égalité pour le joueur
	tourJoueur := p.Dexterite >= m.Dexterite
	if tourJoueur {
		fmt.Println(ui.Vert+"Grâce à votre Dextérité, vous attaquez en premier !"+ui.Reset)
	} else {
		fmt.Printf(ui.Rouge+"%s est plus rapide et prend l'initiative !\n"+ui.Reset, m.Nom)
	}

	// tours de poison qu'il reste sur le monstre
	poisonTours := 0

	for p.PV > 0 && m.PV > 0 {
		// recalculé à chaque tour à cause des potions
		atkTotale := p.Attaque + p.BuffAttaque
		defTotale := p.Defense + p.BuffDefense

		fmt.Printf(ui.Cyan+"\n--- Status : %s (%d/%d PV | %d/%d Mana) | %s (%d/%d PV | %d/%d Mana) ---\n"+ui.Reset,
			p.Nom, p.PV, p.MaxPV, p.Mana, p.MaxMana, m.Nom, m.PV, m.MaxPV, m.Mana, m.MaxMana)

		// Effet du poison sur le monstre
		if poisonTours > 0 && m.PV > 0 {
			m.PV -= 5
			poisonTours--
			fmt.Printf(ui.Violet+"🧪 Le poison inflige 5 dégâts au monstre ! (Reste %d tours, PV Monstre: %d/%d)\n"+ui.Reset, poisonTours, m.PV, m.MaxPV)
			if m.PV <= 0 {
				break
			}
		}

		if tourJoueur {
			fmt.Println("\nVos Actions :")
			fmt.Println("1. Attaque physique")
			fmt.Println("2. Utiliser un sort")
			fmt.Println("3. Inventaire / Potions")

			choix := ui.LireEntree("Choix : ")
			switch choix {
			case "1":
				// dégâts = attaque - défense, jamais moins de 1
				degats := atkTotale - m.Defense
				if degats < 1 {
					degats = 1
				}
				m.PV -= degats
				if m.PV < 0 {
					m.PV = 0
				}
				fmt.Printf(ui.Vert+"Vous frappez %s et infligez %d dégâts ! (PV Monstre: %d/%d)\n"+ui.Reset, m.Nom, degats, m.PV, m.MaxPV)
			case "2":
				if len(p.Sorts) == 0 {
					fmt.Println(ui.Rouge+"Vous ne connaissez aucun sort !"+ui.Reset)
					continue
				}
				fmt.Println("Vos sorts :")
				for i, s := range p.Sorts {
					if s.Degats > 0 {
						fmt.Printf("%d. %s (Dégâts: %d, Coût: %d Mana)\n", i+1, s.Nom, s.Degats, s.CoutMana)
					} else {
						fmt.Printf("%d. %s (Soin: %d, Coût: %d Mana)\n", i+1, s.Nom, s.Soin, s.CoutMana)
					}
				}
				sIdxStr := ui.LireEntree("Sort à lancer : ")
				var sIdx int
				fmt.Sscanf(sIdxStr, "%d", &sIdx)
				if sIdx > 0 && sIdx <= len(p.Sorts) {
					sort := p.Sorts[sIdx-1]
					// on regarde le mana avant de lancer
					if p.Mana >= sort.CoutMana {
						p.Mana -= sort.CoutMana
						if sort.Degats > 0 {
							m.PV -= sort.Degats
							if m.PV < 0 {
								m.PV = 0
							}
							fmt.Printf(ui.Bleu+"✨ Vous lancez %s ! %d dégâts infligés. (PV Monstre: %d/%d)\n"+ui.Reset, sort.Nom, sort.Degats, m.PV, m.MaxPV)
						} else if sort.Soin > 0 {
							p.PV += sort.Soin
							if p.PV > p.MaxPV {
								p.PV = p.MaxPV
							}
							fmt.Printf(ui.Vert+"✨ Vous lancez %s et récupérez %d PV !\n"+ui.Reset, sort.Nom, sort.Soin)
						}
					} else {
						fmt.Println(ui.Rouge+"Pas assez de mana !"+ui.Reset)
						continue
					}
				}
			case "3":
				fmt.Println("Inventaire de combat :")
				for i, it := range p.Inventaire {
					fmt.Printf("%d. %s\n", i+1, it)
				}
				itIdxStr := ui.LireEntree("Objet à utiliser (0 pour annuler) : ")
				var itIdx int
				fmt.Sscanf(itIdxStr, "%d", &itIdx)
				if itIdx > 0 && itIdx <= len(p.Inventaire) {
					objet := p.Inventaire[itIdx-1]
					// le poison est géré ici parce que UtiliserObjet n'a pas accès au monstre
					if objet == "Potion de poison" {
						poisonTours = 3
						p.RetirerObjetIndex(itIdx - 1)
						fmt.Println(ui.Violet+"🧪 Potion de poison lancée sur le monstre ! (5 dégâts par tour pendant 3 tours)"+ui.Reset)
					} else {
						p.UtiliserObjet(objet, itIdx-1)
					}
				} else {
					continue
				}
			default:
				fmt.Println(ui.Rouge+"Action invalide, tour passé."+ui.Reset)
			}
		} else {
			// Tour du monstre : choix entre attaque physique ou compétence magique
			// 1 chance sur 2 d'utiliser son sort s'il a assez de mana
			if m.Mana >= m.Sort.CoutMana && m.Sort.CoutMana > 0 && rand.Intn(100) < 50 {
				m.Mana -= m.Sort.CoutMana
				degats := m.Sort.Degats - defTotale
				if degats < 1 {
					degats = 1
				}
				p.PV -= degats
				if p.PV < 0 {
					p.PV = 0
				}
				fmt.Printf(ui.Rouge+"🔥 %s utilise sa compétence [%s] et inflige %d dégâts magiques !\n"+ui.Reset, m.Nom, m.Sort.Nom, degats)
			} else {
				degats := m.Attaque - defTotale
				if degats < 1 {
					degats = 1
				}
				p.PV -= degats
				if p.PV < 0 {
					p.PV = 0
				}
				fmt.Printf(ui.Rouge+"⚔️ %s attaque et inflige %d dégâts ! (Vos PV: %d/%d)\n"+ui.Reset, m.Nom, degats, p.PV, p.MaxPV)
			}
		}

		// on passe la main à l'autre
		tourJoueur = !tourJoueur
	}

	// Fin des buffs de combat temporaires
	p.BuffAttaque = 0
	p.BuffDefense = 0

	// Résultat
	if p.PV > 0 {
		fmt.Printf(ui.Jaune+"\n!!! VICTOIRE !!! Vous avez vaincu %s !\n"+ui.Reset, m.Nom)
		p.Or += m.GainOr
		fmt.Printf("Vous ramassez %d PO.\n", m.GainOr)
		p.AjouterXP(m.GainXP)

		// Système de Drop Rare
		if len(m.TableDrop) > 0 {
			chance := rand.Intn(100)
			if chance < m.ChancDrop {
				itemDrop := m.TableDrop[rand.Intn(len(m.TableDrop))]
				fmt.Printf(ui.Jaune+"!!! DROP RARE ! Le monstre a laissé tomber : %s !!!\n"+ui.Reset, itemDrop)
				if len(p.Inventaire) < p.MaxInventaire {
					p.Inventaire = append(p.Inventaire, itemDrop)
				} else {
					fmt.Println("Mais votre inventaire est plein !")
				}
			}
		}
		ui.LireEntree("\nAppuyez sur Entrée...")
		return true
	} else {
		// défaite : le joueur est ressuscité avec la moitié de ses pv
		p.VerifierMort()
		return false
	}
}
