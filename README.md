# ENSAP CLI

**Attention**: Cet outil est encore en développement. L'outil `ensap` n'étant accessible que depuis très peu de temps, les commandes présentées ci-dessous peuvent évoluer.

## Description

Outil permettat d'intéragir avec l'API de l'**E**space **N**umérique **S**écurisé de l'**A**gent **P**ublique français (`ensap`).

Cet outil permet pour le moment de:
* lister les différentes fiches de paie présentes dans son espace `ensap`.
* télécharger une fiche de paie identifiée par son _document ID_.

## Installation

### Paquets / binaires pré construits

Vous pourrez trouver des paquets Windows, Linux et MacOS sur [la page de release](https://github.com/riton/ensap-cli/releases) de ce projet.

### Construire à partir des sources

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

* [Opérations sur les fiches de paie](README.payrolls.md)

