package qul

import (
	"log"
	"sort"
)

func GenerateInstances(phrase Phrase) (instances []Instance) {
	for verseKey, ranges := range phrase.Ayah {
		for instanceInVerseIndex, instanceInVerse := range ranges {
			instances = append(instances, GenerateInstance(verseKey, instanceInVerseIndex, instanceInVerse))
		}
	}
	sort.Slice(instances, func(i, j int) bool {
		return compareInstances(instances[i], instances[j])
	})

	return
}

func GenerateInstance(verseKey string, instanceInVerseIndex int, instanceInVerse [2]int) (instance Instance) {
	from := instanceInVerse[0]
	to := instanceInVerse[1]
	instanceInVerseNumber := instanceInVerseIndex + 1

	instance = Instance{
		VerseKey:        verseKey,
		InstanceInVerse: instanceInVerseNumber,
		From:            from,
		To:              to,
	}
	return
}

func compareInstances(a, b Instance) bool {
	compareInstancesErrorMessage := "error while compare instances %+v and %+v: decode verse key %s: %v"

	chapterA, verseA, err := DecodeVerseKey(a.VerseKey)
	if err != nil {
		log.Println(compareInstancesErrorMessage, a, b, a.VerseKey, err)
		return false
	}

	chapterB, verseB, err := DecodeVerseKey(b.VerseKey)
	if err != nil {
		log.Println(compareInstancesErrorMessage, a, b, b.VerseKey, err)
		return false
	}

	if chapterA != chapterB {
		return chapterA < chapterB
	}
	if verseA != verseB {
		return verseA < verseB
	}
	return a.InstanceInVerse < b.InstanceInVerse
}
