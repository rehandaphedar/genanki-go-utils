package qul

type Word struct {
	ID       int    `json:"id"`
	Surah    string `json:"surah"`
	Ayah     string `json:"ayah"`
	Word     string `json:"word"`
	Location string `json:"location"`
	Text     string `json:"text"`
}

type Source struct {
	Key  string `json:"key"`
	From int    `json:"from"`
	To   int    `json:"to"`
}

type Phrase struct {
	Surahs int                 `json:"surahs"`
	Ayahs  int                 `json:"ayahs"`
	Count  int                 `json:"count"`
	Source Source              `json:"source"`
	Ayah   map[string][][2]int `json:"ayah"`
}

type MetadataAyah struct {
	ID          int    `json:"id"`
	SurahNumber int    `json:"surah_number"`
	AyahNumber  int    `json:"ayah_number"`
	VerseKey    string `json:"verse_key"`
	WordsCount  int    `json:"words_count"`
	Text        string `json:"text"`
}

type MetadataJuz struct {
	JuzNumber     int               `json:"juz_number"`
	VersesCount   int               `json:"verses_count"`
	FirstVerseKey string            `json:"first_verse_key"`
	LastVerseKey  string            `json:"last_verse_key"`
	VerseMapping  map[string]string `json:"verse_mapping"`
}

type MetadataHizb struct {
	HizbNumber    int               `json:"hizb_number"`
	VersesCount   int               `json:"verses_count"`
	FirstVerseKey string            `json:"first_verse_key"`
	LastVerseKey  string            `json:"last_verse_key"`
	VerseMapping  map[string]string `json:"verse_mapping"`
}

type MetadataRub struct {
	RubNumber     int               `json:"rub_number"`
	VersesCount   int               `json:"verses_count"`
	FirstVerseKey string            `json:"first_verse_key"`
	LastVerseKey  string            `json:"last_verse_key"`
	VerseMapping  map[string]string `json:"verse_mapping"`
}

type MetadataManzil struct {
	ManzilNumber  int               `json:"manzil_number"`
	VersesCount   int               `json:"verses_count"`
	FirstVerseKey string            `json:"first_verse_key"`
	LastVerseKey  string            `json:"last_verse_key"`
	VerseMapping  map[string]string `json:"verse_mapping"`
}

type MetadataRuku struct {
	RukuNumber      int               `json:"ruku_number"`
	SurahRukuNumber int               `json:"surah_ruku_number"`
	VersesCount     int               `json:"verses_count"`
	FirstVerseKey   string            `json:"first_verse_key"`
	LastVerseKey    string            `json:"last_verse_key"`
	VerseMapping    map[string]string `json:"verse_mapping"`
}

type Index struct {
	Word   WordIndex
	Page   map[string]int
	Juz    map[string]int
	Hizb   map[string]int
	Rub    map[string]int
	Manzil map[string]int
	Ruku   map[string]int
	Tag    TagIndex
}

type WordIndex struct {
	Words      map[string]string
	VerseWords map[string][]string
}

type TagIndex struct {
	Verse map[string][]string
	Page  map[int][]string
}

type MetadataDivision struct {
	Juz    map[string]MetadataJuz
	Hizb   map[string]MetadataHizb
	Rub    map[string]MetadataRub
	Manzil map[string]MetadataManzil
	Ruku   map[string]MetadataRuku
}

type TagFormat struct {
	Chapter *string
	Verse   *string
	Page    *string
	Juz     *string
	Hizb    *string
	Rub     *string
	Manzil  *string
	Ruku    *string
}
