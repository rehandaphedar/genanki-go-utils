# Introduction

Utilities for [genanki-go](https://github.com/npcnixel/genanki-go/).

# QUL

Installation:
```sh
go get git.sr.ht/~rehandaphedar/genanki-go-utils/pkg/qul
```

The package helps interact with the [Quranic Universal Library (QUL)](https://qul.tarteel.ai/resources/quran-metadata).

# DB Fix

Installation:
```sh
go get git.sr.ht/~rehandaphedar/genanki-go-utils/pkg/dbfix
```

The package serves as a workaround until the following PRs are merged:
1. [Fix notes.tags serialization to Anki-compatible tag format (space-delimited) by shepilov · Pull Request #5 · npcnixel/genanki-go · GitHub](https://github.com/npcnixel/genanki-go/pull/5)
2. [make GenerateIntID always return a positive int64 by rehandaphedar · Pull Request #6 · npcnixel/genanki-go · GitHub](https://github.com/npcnixel/genanki-go/pull/6)

To use it run:
```golang
...
dbfix.FixDb(inputPath, outputPath)
```

Note:
The fix for 2. only regenerates the `id` of the `cards` table. As for the `id` and `guid` of the `notes` table, you should manually change the ID of the note after creation like so:
```golang
note := genanki.NewNote(...)
note.ID = dbfix.GenerateIntID()
deck.AddNote(note)
```

This is because of two reasons:
1. After the initial db is generated, the `id` of the `notes` table has already been used for generated `guid` and linked to the `nid` of the `cards` table. Modifying it after the fact means manually modifying multiple downstream side effects.
2. There is a way to pass a custom ID for a note, however, there is no way to pass a custom ID for a card.
