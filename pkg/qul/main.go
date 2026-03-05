package qul

import (
	"database/sql"
	"fmt"
	"log"
	"sort"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func BuildIndex(
	layoutPath string, words map[string]Word,
	metadataDivision MetadataDivision,
	tagFormat TagFormat,
) (Index, error) {
	wordIndex, err := buildWordIndex(words)
	if err != nil {
		return Index{}, fmt.Errorf("build word index: %w", err)
	}
	pageIndex, err := buildPageIndex(layoutPath, words)
	if err != nil {
		return Index{}, fmt.Errorf("build page index: %w", err)
	}

	index := Index{
		Word:   wordIndex,
		Page:   pageIndex,
		Juz:    make(map[string][]int),
		Hizb:   make(map[string][]int),
		Rub:    make(map[string][]int),
		Manzil: make(map[string][]int),
		Ruku:   make(map[string][]int),
		Tag: TagIndex{
			Verse: make(map[string][]string),
			Page:  make(map[int][]string),
		},
	}

	for _, juz := range metadataDivision.Juz {
		addVerseEntries(index.Juz, juz.JuzNumber, juz.VerseMapping)
	}
	for _, hizb := range metadataDivision.Hizb {
		addVerseEntries(index.Hizb, hizb.HizbNumber, hizb.VerseMapping)
	}
	for _, rub := range metadataDivision.Rub {
		addVerseEntries(index.Rub, rub.RubNumber, rub.VerseMapping)
	}
	for _, manzil := range metadataDivision.Manzil {
		addVerseEntries(index.Manzil, manzil.ManzilNumber, manzil.VerseMapping)
	}
	for _, ruku := range metadataDivision.Ruku {
		addVerseEntries(index.Ruku, ruku.RukuNumber, ruku.VerseMapping)
	}

	pageTagSets := make(map[int]map[string]bool)
	for verseKey, page := range index.Page {

		tags, err := buildTagsForVerse(index, tagFormat, verseKey)
		if err != nil {
			return Index{}, fmt.Errorf("build tags for verse %s: %w", verseKey, err)
		}
		index.Tag.Verse[verseKey] = tags

		if _, exists := pageTagSets[page]; !exists {
			pageTagSets[page] = make(map[string]bool)
		}
		for _, tag := range tags {
			pageTagSets[page][tag] = true
		}
	}

	for page, tagSet := range pageTagSets {
		tags := make([]string, 0, len(tagSet))
		for tag := range tagSet {
			tags = append(tags, tag)
		}
		index.Tag.Page[page] = tags
	}

	return index, nil
}

func BuildTagsForPhrase(index Index, phrase Phrase) []string {
	tagSet := make(map[string]bool)

	for verseKey := range phrase.Ayah {
		if tags, ok := index.Tag.Verse[verseKey]; ok {
			for _, tag := range tags {
				tagSet[tag] = true
			}
		}
	}

	tags := make([]string, 0, len(tagSet))
	for t := range tagSet {
		tags = append(tags, t)
	}
	return tags
}

func buildWordIndex(words map[string]Word) (WordIndex, error) {
	type entry struct {
		position int
		text     string
	}

	ayahMap := make(map[string][]entry, len(words))
	wordIndex := WordIndex{
		Words:      make(map[string]string, len(words)),
		VerseWords: make(map[string][]string),
	}

	for wordKey, word := range words {
		wordIndex.Words[wordKey] = word.Text

		chapterNumber, verseNumber, wordNumber, err := parseWordLocation(word, wordKey)
		if err != nil {
			return WordIndex{}, err
		}

		verseKey := EncodeVerseKey(chapterNumber, verseNumber)
		ayahMap[verseKey] = append(ayahMap[verseKey], entry{wordNumber, word.Text})
	}

	for verseKey, entries := range ayahMap {
		sort.Slice(entries, func(i, j int) bool { return entries[i].position < entries[j].position })
		texts := make([]string, len(entries))
		for i, e := range entries {
			texts[i] = e.text
		}
		wordIndex.VerseWords[verseKey] = texts
	}

	return wordIndex, nil
}

func parseWordLocation(word Word, wordKey string) (chapterNumber, verseNumber, wordNumber int, err error) {
	chapterNumber, err = strconv.Atoi(word.Surah)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid chapter index %q at loc %q: %w", word.Surah, wordKey, err)
	}
	verseNumber, err = strconv.Atoi(word.Ayah)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid verse index %q at loc %q: %w", word.Ayah, wordKey, err)
	}
	wordNumber, err = strconv.Atoi(word.Word)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid word index %q at loc %q: %w", word.Word, wordKey, err)
	}
	return chapterNumber, verseNumber, wordNumber, nil
}

func buildPageIndex(dbPath string, words map[string]Word) (map[string]int, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open layout db: %w", err)
	}
	defer db.Close()

	wordPageMap, err := buildWordPageMap(db)
	if err != nil {
		return nil, fmt.Errorf("build word page map: %w", err)
	}

	pageIndex := make(map[string]int)
	for wordKey, word := range words {
		if word.Word != "1" {
			continue
		}
		if page, ok := wordPageMap[word.ID]; ok {
			chapterNumber, err := strconv.Atoi(word.Surah)
			if err != nil {
				return nil, fmt.Errorf("invalid chapter index %q at loc %q: %w", word.Surah, wordKey, err)
			}
			verseNumber, err := strconv.Atoi(word.Ayah)
			if err != nil {
				return nil, fmt.Errorf("invalid verse index %q at loc %q: %w", word.Ayah, wordKey, err)
			}
			pageIndex[EncodeVerseKey(chapterNumber, verseNumber)] = page
		}
	}
	return pageIndex, nil
}

func buildWordPageMap(db *sql.DB) (map[int]int, error) {
	rows, err := db.Query(`SELECT first_word_id, last_word_id, page_number FROM pages WHERE line_type = 'ayah'`)

	if err != nil {
		return nil, fmt.Errorf("query pages table: %w", err)
	}
	defer rows.Close()

	wordPageMap := make(map[int]int)
	for rows.Next() {
		var firstWord, lastWord, page int
		if err := rows.Scan(&firstWord, &lastWord, &page); err != nil {
			return nil, fmt.Errorf("scan pages row: %w", err)
		}
		for wordId := firstWord; wordId <= lastWord; wordId++ {
			wordPageMap[wordId] = page
		}
	}

	return wordPageMap, rows.Err()
}

func addVerseEntries(m map[string][]int, num int, verseMapping map[string]string) {
	for chapterStr, rangeStr := range verseMapping {
		chapter, err := strconv.Atoi(chapterStr)
		if err != nil {
			log.Printf("invalid chapter %q in verse mapping: %v", chapterStr, err)
			continue
		}
		from, to, err := DecodeRange(rangeStr)
		if err != nil {
			log.Printf("error while decoding range %s: %v", rangeStr, err)
			continue
		}
		for verse := from; verse <= to; verse++ {
			key := EncodeVerseKey(chapter, verse)
			m[key] = append(m[key], num)
		}
	}
}

func buildTagsForVerse(index Index,
	tagFormat TagFormat,
	verseKey string) ([]string, error) {

	chapter, _, err := DecodeVerseKey(verseKey)
	if err != nil {
		return nil, fmt.Errorf("decode verse key %s: %w", verseKey, err)
	}

	paddedVerseKey, err := PadVerseKey(verseKey)
	if err != nil {
		return nil, fmt.Errorf("pad verse key %s: %w", verseKey, err)
	}

	tagSet := map[string]bool{
		fmt.Sprintf(*tagFormat.Chapter, chapter):      true,
		fmt.Sprintf(*tagFormat.Verse, paddedVerseKey): true,
	}

	if page, ok := index.Page[verseKey]; ok {
		tagSet[fmt.Sprintf(*tagFormat.Page, page)] = true
	}

	for _, juz := range DeduplicateInts(index.Juz[verseKey]) {
		tagSet[fmt.Sprintf(*tagFormat.Juz, juz)] = true
	}

	for _, hizb := range DeduplicateInts(index.Hizb[verseKey]) {
		tagSet[fmt.Sprintf(*tagFormat.Hizb, hizb)] = true
	}

	for _, rub := range DeduplicateInts(index.Rub[verseKey]) {
		tagSet[fmt.Sprintf(*tagFormat.Rub, rub)] = true
	}

	for _, manzil := range DeduplicateInts(index.Manzil[verseKey]) {
		tagSet[fmt.Sprintf(*tagFormat.Manzil, manzil)] = true
	}

	for _, ruku := range DeduplicateInts(index.Ruku[verseKey]) {
		tagSet[fmt.Sprintf(*tagFormat.Ruku, ruku)] = true
	}

	tags := make([]string, 0, len(tagSet))
	for t := range tagSet {
		tags = append(tags, t)
	}
	return tags, nil
}
