package main

type Generic struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Chaos struct {
	Trigger []string `json:"trigger"`
	Target  []string `json:"target"`
}

type Mutation struct {
	Name        string `json:"name"`
	Target      string `json:"target"`
	Description string `json:"description"`
}

type Affix struct {
	Description string   `json:"description"`
	PointValue  string   `json:"pointValue"`
	Upgrade     string   `json:"upgrade"`
	Affinities  []string `json:"affinities"`
}

type Amulet struct {
	Name     string  `json:"name"`
	Base     Affix   `json:"base"`
	Upgrades []Affix `json:"upgrades"`
}

type GlyphPath struct {
	Theme string   `json:"theme"`
	Tiers []string `json:"tiers"`
}

type Relic struct {
	Name            string  `json:"name"`
	StartingAffixes []Affix `json:"startingAffixes"`
}

type Crystal struct {
	Name     string   `json:"name"`
	Passives []string `json:"passives"`
	Powers   []string `json:"powers"`
}

type Augment struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Augments struct {
	Quincy    []Augment `json:"Quincy"`
	Viktor    []Augment `json:"Viktor"`
	Arthur    []Augment `json:"Arthur"`
	Nathaniel []Augment `json:"Nathaniel"`
}

// Party is the campaign's roster. SpecialisationsTaken records which sessions
// each specialisation was picked at; its length drives how often an enchanted
// item is replaced by a specialised loot type.
type Party struct {
	SpecialisationsTaken []string      `json:"specialisationsTaken"`
	Members              []PartyMember `json:"members"`
}

// PartyMember is one roster entry. Non-players (the sidekick) roll insight but
// have no crystal and take no specialisations.
type PartyMember struct {
	Name            string   `json:"name"`
	Player          bool     `json:"player"`
	FlareChance     int      `json:"flareChance"`
	Insight         int      `json:"insight"`
	Specialisations []string `json:"specialisations"`
}

type Boon struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Children    []Boon `json:"children"`
}
