# Opérations sur les fiches de paie

## Lister les fichers de paie présentes dans mon espace `ensap`

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

## Télécharger une fiche de paie en connaissant son identifiant (`UUID`)

```
$ ensap download-remuneration-paie adc2cf49-6c5f-4765-1e74-7bacc1f2be68
```

La commande ci-dessus téléchargera le document identifié par l'`UUID` `adc2cf49-6c5f-4765-1e74-7bacc1f2be68` (obtenu via la commande `ensap list-remuneration-paie`) et
enregistrera le document sous le nom `2021_07_BP_juillet.pdf`.

Vous pouvez utiliser l'option `-o` ou `--out` pour spécifier un fichier de destination différent.

## Utilisation avancée

### Modifier le format d'affichage de la commande `list-remuneration-paie`

La commande `list-remuneration-paie` affiche par défault ses données au format présenté dans la section [Lister les fichers de paie présentes dans mon espace `ensap`](#lister-les-fichers-de-paie-présentes-dans-mon-espace-ensap).

Le _template_ par défaut utilisé pour formatter la sortie de cette commande est le suivant:

```go
Remuneration paie documents:
{{- range $index, $document := . }}
{{- with $document }}
* {{ .NomDocument }}:
  - Document UUID: {{ .DocumentUUID }}
  - Document date: {{ .DateDocument }}
  - Labels:
    - {{ .Libelle1 }}
    - {{ .Libelle2 }}
    - {{ .Libelle3 }}
{{- end }}
{{- end }}
```

L'utilisateur peut spécifier un _template_ alternatif à utiliser dans le fichier de configuration du CLI `ensap`.

Exemple:

```yaml
rendering:
  templates:
    list-remuneration-paie-documents: |
      Only document UUIDs:
      {{- range $index, $document := . }}
      {{- with $document }}
      Document UUID = {{ .DocumentUUID }}
      {{- end }}
      {{- end }}
```

est un template qui permet un affichage comme:

```
$ ensap list-remuneration-paie
Only document UUIDs:
Document UUID = zfc9cc49-4c5f-1265-8d78-7bcaa2e2be68
Document UUID = 41a17871-ae47-4749-809f-2d1e2bcd1e26
```