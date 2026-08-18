package qul

import (
	"fmt"
	"sort"
)

func GenerateInstances(phrase Phrase) ([]Instance, error) {
	instances := []Instance{}
	for verseKey, ranges := range phrase.Ayah {
		chapter, verse, err := DecodeVerseKey(verseKey)
		if err != nil {
			return []Instance{}, fmt.Errorf("invalid verse key %s: %w", verseKey, err)
		}

		for instanceInVerseIndex, instanceInVerse := range ranges {
			instances = append(instances, GenerateInstance(chapter, verse, instanceInVerseIndex, instanceInVerse))
		}
	}
	sort.Slice(instances, func(i, j int) bool {
		return compareInstances(instances[i], instances[j])
	})

	return instances, nil
}

func GenerateInstance(chapter, verse, instanceInVerseIndex int, instanceInVerse [2]int) (instance Instance) {
	from := instanceInVerse[0]
	to := instanceInVerse[1]
	instanceInVerseNumber := instanceInVerseIndex + 1

	instance = Instance{
		Chapter:         chapter,
		Verse:           verse,
		InstanceInVerse: instanceInVerseNumber,
		From:            from,
		To:              to,
	}
	return
}

func compareInstances(a, b Instance) bool {
	if a.Chapter != b.Chapter {
		return a.Chapter < b.Chapter
	}
	if a.Verse != b.Verse {
		return a.Verse < b.Verse
	}
	return a.InstanceInVerse < b.InstanceInVerse
}
