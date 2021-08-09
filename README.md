# ENSAP CLI

**Attention**: Cet outil est encore en développement. L'outil `ensap` n'étant accessible que depuis très peu de temps, les commandes présentées ci-dessous peuvent évoluer.

## Description

Outil permettat d'intéragir avec l'API de l'**E**space **N**umérique **S**écurisé de l'**A**gent **P**ublique français (`ensap`).

Cet outil permet pour le moment de:
* lister les différentes fiches de paie présentes dans son espace `ensap`.
* télécharger une fiche de paie identifiée par son _document ID_.

## Installation

En attendant que des paquets ou binaires soient construits automatiquement, vous devez vous-même _compiler_ ce programme.

Pour cela vous aurez besoin d'avoir le compilateur `go` installé (voir [site officiel](https://golang.org/dl/)).

```
$ CGO_ENABLED=0 go build -ldflags "-s -w -extldflags '-static'" -o ensap .
```

## Configuration

La commande `ensap` utilise un fichier de configuration pour faciliter l'utilisation du _CLI_.

Le fichier de configuration à utiliser peut-être spécifié via l'option `--config`.
L'emplacement par défaut du fichier de configuration à créer est `~/.ensap.yaml`.

### Exemple de configuration

```
$ cat > ~/.ensap.yaml <<EOF
---
api:
  endpoint: 'ensap.gouv.fr'
  username: 'METTRE_ICI_VOTRE_IDENTIFIANT' # mettre ici l'identifiant vous servant à vous connecter au site ensap.gouv.fr
  password: 'METTRE_ICI_VOTRE_MOT_DE_PASSE' # mettre ici le mot de passe vous servant à vous connecter au site ensap.gouv.fr
EOF
```

## Utilisation

### Opérations sur les fiches de paie

#### Lister les fichers de paie présentes dans mon espace `ensap`

```
$ ensap list-remuneration-paie
Remuneration paie documents:
* 2021_07_BP_juillet.pdf:
  - Document UUID: adc2cf49-6c5f-4765-1e74-7bacc1f2be68
  - Document date: 01 Jul 21 00:00 CEST
  - Labels:
    - Juillet 2021
    - 2021_07_BP_juillet.pdf (PDF, 490 Ko)
    - 42 BTC
```

#### Télécharger une fiche de paie en connaissant son identifiant (`UUID`)

```
$ ensap download-remuneration-paie adc2cf49-6c5f-4765-1e74-7bacc1f2be68
```

La commande ci-dessus téléchargera le document identifié par l'`UUID` `adc2cf49-6c5f-4765-1e74-7bacc1f2be68` (obtenu via la commande `ensap list-remuneration-paie`) et
enregistrera le document sous le nom `2021_07_BP_juillet.pdf`.

Vous pouvez utiliser l'option `-o` ou `--out` pour spécifier un fichier de destination différent.