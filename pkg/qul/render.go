package qul

import (
	"fmt"
	"log"
)

func RenderInstances(index Index, instances []Instance, metadataAyahByVerseKey map[string]MetadataAyah, imageFormat string) {
	for i := range instances {
		verseKey := EncodeVerseKey(instances[i].Chapter, instances[i].Verse)

		instances[i].Phrase = RenderPhrase(
			index.Word,
			verseKey,
			instances[i].From,
			instances[i].To,
		)
		instances[i].Context = RenderContext(
			index.Word,
			metadataAyahByVerseKey,
			verseKey,
			instances[i].From,
		)
		instances[i].Continuation = RenderContinuation(
			index.Word,
			metadataAyahByVerseKey,
			verseKey,
			instances[i].To,
		)
		instances[i].Image = fmt.Sprintf(imageFormat, index.Page[verseKey])
	}
}

func RenderPhrase(wordIndex WordIndex, verseKey string, from, to int) []string {
	return RenderRange(wordIndex, Source{Key: verseKey, From: from, To: to})
}

func RenderContext(wordIndex WordIndex, metadataAyahByVersekey map[string]MetadataAyah, verseKey string, from int) []string {
	if from == 1 {
		previousVerseKey, found := GetPreviousVerseKey(metadataAyahByVersekey, verseKey)
		if found {
			return RenderVerseFrom(wordIndex, metadataAyahByVersekey, previousVerseKey, 1)
		}
		// TODO: Better Context, Continuation edge cases?
		return []string{}
	}
	return RenderRange(wordIndex, Source{Key: verseKey, From: 1, To: from - 1})
}

func RenderContinuation(wordIndex WordIndex, metadataAyahByVerseKey map[string]MetadataAyah, verseKey string, to int) []string {
	words := wordIndex.VerseWords[verseKey]
	if to+1 == len(words) {
		nextVerseKey, found := GetNextVerseKey(metadataAyahByVerseKey, verseKey)
		if found {
			return RenderVerseFrom(wordIndex, metadataAyahByVerseKey, nextVerseKey, 1)
		}
		return []string{}
	}
	return RenderVerseFrom(wordIndex, metadataAyahByVerseKey, verseKey, to+1)
}

func RenderVerseFrom(wordIndex WordIndex, metadataAyahByVerseKey map[string]MetadataAyah, verseKey string, from int) []string {
	return RenderRange(wordIndex, Source{Key: verseKey, From: from, To: metadataAyahByVerseKey[verseKey].WordsCount})
}

func RenderRange(wordIndex WordIndex, source Source) []string {
	words := wordIndex.VerseWords[source.Key]

	if source.To >= len(words) {
		log.Printf("invalid range %+v, silently fixing", source)
		source.To = len(words)
	}

	if (source.To + 1) == len(words) {
		source.To++
	}

	var parts []string
	for i := source.From - 1; i < source.To; i++ {
		parts = append(parts, words[i])
	}
	return parts
}
