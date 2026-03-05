package qul

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

func DecodeRange(rangeStr string) (int, int, error) {
	parts := strings.Split(rangeStr, "-")

	from, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid first part in range %s: %w", rangeStr, err)
	}
	if len(parts) == 1 {
		return from, from, nil
	}
	to, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid second part in range %s: %w", rangeStr, err)
	}
	return from, to, nil
}

func EncodeVerseKey(chapter, verse int) string {
	return fmt.Sprintf("%d:%d", chapter, verse)
}

func DecodeVerseKey(verseKey string) (int, int, error) {
	parts := strings.Split(verseKey, ":")
	chapter, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid first part in verse key %s: %w", verseKey, err)
	}
	verse, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid second part in verse key %s: %w", verseKey, err)
	}
	return chapter, verse, nil
}

func PadVerseKey(verseKey string) (string, error) {
	chapter, verse, err := DecodeVerseKey(verseKey)
	if err != nil {
		return "", fmt.Errorf("invalid verse key %s: %w", verseKey, err)

	}
	return fmt.Sprintf("%03d:%03d", chapter, verse), nil
}

func DeduplicateInts(ints []int) []int {
	seen := make(map[int]bool)
	var deduplicated []int
	for _, value := range ints {
		if !seen[value] {
			seen[value] = true
			deduplicated = append(deduplicated, value)
		}
	}
	return deduplicated
}

func GetPreviousVerseKey(metadataAyahByVerseKey map[string]MetadataAyah, verseKey string) (string, bool) {
	chapter, verse, err := DecodeVerseKey(verseKey)
	if err != nil {
		log.Printf("error while decoding verse key %s: %v", verseKey, err)
		return "", false
	}
	if verse <= 1 {
		if chapter <= 1 {
			return "", false
		}
		previousChapter := chapter - 1
		maxVerse := 0
		for _, metadataAyahEntry := range metadataAyahByVerseKey {
			if metadataAyahEntry.SurahNumber == previousChapter {
				maxVerse = max(maxVerse, metadataAyahEntry.AyahNumber)
			}
		}
		if maxVerse == 0 {
			return "", false
		}
		return fmt.Sprintf("%d:%d", previousChapter, maxVerse), true
	}
	return fmt.Sprintf("%d:%d", chapter, verse-1), true
}

func GetNextVerseKey(metadataAyahByVerseKey map[string]MetadataAyah, verseKey string) (string, bool) {
	chapter, verse, err := DecodeVerseKey(verseKey)
	if err != nil {
		log.Printf("error while decoding verse key %s: %v", verseKey, err)
		return "", false
	}

	maxVerse := 0
	for _, metadataAyahEntry := range metadataAyahByVerseKey {
		if metadataAyahEntry.SurahNumber == chapter {
			maxVerse = max(maxVerse, metadataAyahEntry.AyahNumber)
		}
	}
	if verse < maxVerse {
		return fmt.Sprintf("%d:%d", chapter, verse+1), true
	}

	nextChapter := chapter + 1
	for _, metadataAyahEntry := range metadataAyahByVerseKey {
		if metadataAyahEntry.SurahNumber == nextChapter && metadataAyahEntry.AyahNumber == 1 {
			return fmt.Sprintf("%d:1", nextChapter), true
		}
	}
	return "", false
}
