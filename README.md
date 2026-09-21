# Projet RED : The Rise of the Demon God

RPG en console écrit en **Go**. On crée un héros, on s'équipe, on explore des zones de plus en plus dangereuses et on affronte le Dieu Démon. Une fois le boss vaincu, la partie repart en **New Game+** avec des monstres plus forts.

Équipe : Nathan, Soumaila, Shilsa.

## Sommaire

- [Lancer le jeu](#lancer-le-jeu)
- [Comment jouer](#comment-jouer)
- [Les races](#les-races)
- [Les zones et les monstres](#les-zones-et-les-monstres)
- [Combat](#combat)
- [Marchand et forgeron](#marchand-et-forgeron)
- [New Game+](#new-game)
- [Sauvegardes](#sauvegardes)
- [Structure du projet](#structure-du-projet)

## Lancer le jeu

Prérequis : Go 1.20 ou plus, et un terminal qui gère l'UTF-8 et les couleurs ANSI (Windows Terminal, PowerShell récent, terminal Linux ou macOS).

```bash
# lancer directement
go run .

# ou compiler puis lancer
go build -o projet-red.exe
./projet-red.exe          # projet-red.exe sous Windows

# ou si problème d'autorisation
mkdir -p tmp
GOTMPDIR=$(pwd)/tmp go run main.go
```

Le jeu doit être lancé depuis la racine du projet, car le dossier `saves/` est relatif à l'endroit où on l'exécute.

## Comment jouer

Tout se joue avec des menus numérotés : on tape le chiffre puis Entrée.

**Menu principal**

1. Nouvelle Partie
2. Charger une Partie
3. secret (un petit easter egg)
4. Quitter

**Menu du village**

1. Info Personnage
2. Inventaire & Consommables
3. Aller chez le Marchand
4. Aller chez le Forgeron
5. Partir en Aventure
6. Mannequin d'Entraînement
7. Sauvegarder et Quitter

Au départ, le héros a 100 PO, un inventaire de 10 places, une Potion de soin et le sort Coup de Poing (10 dégâts, 5 mana).

## Les races

| Race | PV | Mana | Attaque | Défense | Dextérité |
|---|---|---|---|---|---|
| Humain (par défaut) | 100 | 50 | 12 | 5 | 10 |
| Elfe | 80 | 100 | 8 | 3 | 15 |
| Nain | 130 | 30 | 14 | 8 | 6 |
| Mage | 70 | 120 | 6 | 2 | 12 |
| Dragonoïde | 140 | 40 | 16 | 7 | 5 |

À chaque niveau : +20 PV max, +10 mana max, +3 attaque, +2 défense, +2 dextérité, et PV/mana remis au maximum. L'XP nécessaire pour le niveau suivant augmente de 40 % à chaque fois.

## Les zones et les monstres

| Zone | Niveau requis | Monstres |
|---|---|---|
| Forêt des Gobelins | 1 | Gobelin Éclaireur, Loup Sauvage |
| Caverne des Trolls | 3 | Troll des Cavernes, Sorcier Squelette |
| Donjon du Dragon | 5 | Drake Infernale |
| Antre du Dieu Démon | aucun | Le Dieu Démon (boss final) |

Les monstres lâchent parfois des matériaux ou des potions. Le boss drop à coup sûr.

## Combat

- Le combat est au tour par tour. Celui qui a la plus grande dextérité commence (le joueur en cas d'égalité).
- Dégâts d'une attaque : attaque moins défense de la cible, avec un minimum de 1.
- À son tour, le joueur peut attaquer, lancer un sort (si assez de mana) ou utiliser un objet.
- La Potion de poison inflige 5 dégâts au monstre pendant 3 tours.
- Les monstres utilisent leur sort une fois sur deux quand ils ont assez de mana.
- Les buffs des potions de Force et d'Armure ne durent que jusqu'à la fin du prochain combat.
- Il n'y a pas de game over : un héros vaincu revient avec la moitié de ses PV.

## Marchand et forgeron

**Marchand**

- Potions : soin (4 PO), mana (5), poison (8), Force (12), Armure (12), Grande potion de soin (15), Élixir (25). Certaines demandent un niveau minimum.
- Matériaux : plume de corbeau, cuir de sanglier, fourrure de loup, lingot de fer, peau de troll, minerai de Mithril.
- Livres de sorts : Boule de Feu, Soin Magique, Éclair Céleste.
- Extension d'inventaire : +10 places pour 30 PO, trois fois maximum.
- Revente des matériaux.

**Forgeron**

Il fabrique de l'équipement pour trois emplacements (tête, torse, arme) contre de l'or et des matériaux. Les bonus dépendent du niveau du héros au moment de la fabrication.

| Emplacement | Objets |
|---|---|
| Tête | Casque en Cuir (niv. 1), Casque en Fer (niv. 3) |
| Torse | Tunique Légère (niv. 1), Plastron en Peau (niv. 2), Armure en Mithril (niv. 4) |
| Arme | Épée de Bois (niv. 1), Épée en Fer (niv. 2), Lame de Mithril (niv. 4) |

## New Game+

Battre le Dieu Démon fait passer le héros en NG+1 (puis NG+2, etc.). Pour chaque niveau de NG+ :

- les stats des monstres augmentent de 50 % (70 % pour le boss) ;
- l'or et l'XP gagnés augmentent de 80 % (100 % pour le boss) ;
- la chance de drop rare gagne 10 points.

Le niveau du héros fait aussi monter les stats des monstres.

## Sauvegardes

« Sauvegarder et Quitter » écrit le héros dans `saves/<Nom>.json`. Depuis le menu principal, « Charger une Partie » liste les fichiers du dossier `saves/` et recharge celui qu'on choisit.

## Structure du projet

```
projet-red/
├── main.go                  menu principal
├── go.mod
├── saves/                   sauvegardes JSON
└── pkg/
    ├── character/           héros : races, stats, inventaire, sorts, XP
    ├── combat/              combat au tour par tour, structure des monstres
    ├── adventure/           zones, bestiaire, boss final, difficulté NG+
    ├── economy/             marchand, forgeron, fabrication d'objets
    ├── engine/              menu du village, sauvegarde et chargement
    └── ui/                  affichage, couleurs ANSI, dessins ASCII
```

Un point à savoir si vous modifiez le code : les noms des monstres servent à choisir leur dessin ASCII dans `ui.DisplayMonsterArt`, ils doivent donc rester identiques des deux côtés.
