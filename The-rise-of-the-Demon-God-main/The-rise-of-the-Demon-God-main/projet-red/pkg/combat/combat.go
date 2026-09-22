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
	// la plus grosse dextérité commence, égalité pour le joueur
	tourJoueur := p.Dexterite >= m.Dexterite
	poisonTours := 0
	dernierMessage := ""

	if tourJoueur {
		dernierMessage = "Grâce à votre Dextérité, vous attaquez en premier !"
	} else {
		dernierMessage = fmt.Sprintf("%s est plus rapide et prend l'initiative !", m.Nom)
	}

	for p.PV > 0 && m.PV > 0 {
		atkTotale := p.Attaque + p.BuffAttaque
		defTotale := p.Defense + p.BuffDefense

		// Construction de la liste des états pour l'affichage ui
		var etatsMonstre []string
		if poisonTours > 0 {
			etatsMonstre = append(etatsMonstre, fmt.Sprintf("Poison (%dt)", poisonTours))
		}
		var etatsJoueur []string

		// Effet du poison en début de tour si actif
		if poisonTours > 0 && m.PV > 0 {
			m.PV -= 5
			poisonTours--
			if m.PV < 0 {
				m.PV = 0
			}
			dernierMessage += fmt.Sprintf(" | 🧪 Poison: -5 PV au monstre (reste %dt)", poisonTours)
			if m.PV <= 0 {
				break
			}
		}

		// Utilisation de la nouvelle interface dynamique
		ui.RedessinerCombat(
			p.Nom, p.PV, p.MaxPV, p.Mana, p.MaxMana, etatsJoueur,
			m.Nom, m.PV, m.MaxPV, etatsMonstre,
			dernierMessage,
		)

		if tourJoueur {
			fmt.Println("1. Attaque physique")
			fmt.Println("2. Utiliser un sort")
			fmt.Println("3. Inventaire / Potions")
			fmt.Println("4. Tentative de Fuite")

			choix := ui.LireEntree("\nChoix : ")
			switch choix {
			case "1":
				degats := atkTotale - m.Defense
				if degats < 1 {
					degats = 1
				}
				m.PV -= degats
				if m.PV < 0 {
					m.PV = 0
				}
				dernierMessage = fmt.Sprintf("Vous frappez %s et infligez %d dégâts !", m.Nom, degats)

			case "2":
				if len(p.Sorts) == 0 {
					dernierMessage = "Vous ne connaissez aucun sort !"
					continue
				}
				fmt.Println("\n--- Vos sorts ---")
				for i, s := range p.Sorts {
					if s.Degats > 0 {
						fmt.Printf("%d. %s (Dégâts: %d, Coût: %d Mana)\n", i+1, s.Nom, s.Degats, s.CoutMana)
					} else {
						fmt.Printf("%d. %s (Soin: %d, Coût: %d Mana)\n", i+1, s.Nom, s.Soin, s.CoutMana)
					}
				}
				sIdxStr := ui.LireEntree("Sort à lancer (0 pour annuler) : ")
				var sIdx int
				fmt.Sscanf(sIdxStr, "%d", &sIdx)
				if sIdx > 0 && sIdx <= len(p.Sorts) {
					sort := p.Sorts[sIdx-1]
					if p.Mana >= sort.CoutMana {
						p.Mana -= sort.CoutMana
						if sort.Degats > 0 {
							m.PV -= sort.Degats
							if m.PV < 0 {
								m.PV = 0
							}
							dernierMessage = fmt.Sprintf("✨ Vous lancez %s et infligez %d dégâts !", sort.Nom, sort.Degats)
						} else if sort.Soin > 0 {
							p.PV += sort.Soin
							if p.PV > p.MaxPV {
								p.PV = p.MaxPV
							}
							dernierMessage = fmt.Sprintf("✨ Vous lancez %s et récupérez %d PV !", sort.Nom, sort.Soin)
						}
					} else {
						dernierMessage = "Pas assez de mana !"
						continue
					}
				} else {
					continue
				}

			case "3":
				if len(p.Inventaire) == 0 {
					dernierMessage = "Votre inventaire est vide !"
					continue
				}
				fmt.Println("\n--- Inventaire de combat ---")
				for i, it := range p.Inventaire {
					fmt.Printf("%d. %s\n", i+1, it)
				}
				itIdxStr := ui.LireEntree("Objet à utiliser (0 pour annuler) : ")
				var itIdx int
				fmt.Sscanf(itIdxStr, "%d", &itIdx)
				if itIdx > 0 && itIdx <= len(p.Inventaire) {
					objet := p.Inventaire[itIdx-1]
					if objet == "Potion de poison" {
						poisonTours = 3
						p.RetirerObjetIndex(itIdx - 1)
						dernierMessage = "🧪 Potion de poison lancée sur le monstre !"
					} else {
						p.UtiliserObjet(objet, itIdx-1)
						dernierMessage = fmt.Sprintf("Vous avez utilisé : %s", objet)
					}
				} else {
					continue
				}

			case "4":
				if rand.Float32() < 0.5 {
					ui.ClearConsole()
					fmt.Println(ui.Jaune + "Vous avez réussi à fuir le combat !" + ui.Reset)
					ui.LireEntree("\nAppuyez sur Entrée...")
					return false
				} else {
					dernierMessage = "Échec de la fuite ! Le monstre vous bloque."
				}

			default:
				dernierMessage = "Action invalide, tour passé."
			}
		} else {
			// Tour du monstre
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
				dernierMessage = fmt.Sprintf("🔥 %s utilise %s et vous inflige %d dégâts magiques !", m.Nom, m.Sort.Nom, degats)
			} else {
				degats := m.Attaque - defTotale
				if degats < 1 {
					degats = 1
				}
				p.PV -= degats
				if p.PV < 0 {
					p.PV = 0
				}
				dernierMessage = fmt.Sprintf("⚔️ %s vous attaque et vous inflige %d dégâts !", m.Nom, degats)
			}
		}

		tourJoueur = !tourJoueur
	}

	// Fin des buffs
	p.BuffAttaque = 0
	p.BuffDefense = 0

	// Résultat
	ui.ClearConsole()
	if p.PV > 0 {
		fmt.Printf(ui.Jaune+"\n!!! VICTOIRE !!! Vous avez vaincu %s !\n"+ui.Reset, m.Nom)
		p.Or += m.GainOr
		fmt.Printf("Vous ramassez %d PO.\n", m.GainOr)
		p.AjouterXP(m.GainXP)

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
		p.VerifierMort()
		return false
	}
}
