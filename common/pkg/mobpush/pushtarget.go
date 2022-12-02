package mobpush

type TargetType int

const (
	TARGET_ALL   TargetType = 1
	TARGET_ALIAS TargetType = 2
	TARGET_TAGS  TargetType = 3
	TARGET_RIDS  TargetType = 4
	TARGET_AREAS TargetType = 9
)

type PushTarget struct {
	Target   int      `json:"target"`
	Tags     []string `json:"tags"`
	TagsType string   `json:"tagsType"`
	Alias    []string `json:"alias"`
	Rids     []string `json:"rids"`

	Block     string     `json:"block"`
	City      string     `json:"city"`
	Country   string     `json:"country"`
	Province  string     `json:"province"`
	PushAreas *PushAreas `json:"pushAreas,omitempty"`
}

type PushLabel struct {
	LabelIds []string `json:"labelIds"`
	MobId    string   `json:"mobId"`
	Type     int      `json:"type"`
}

type PushAreas struct {
	Countries []PushCountry `json:"countries"`
}

type PushCountry struct {
	Country          string         `json:"country"`
	Provinces        []PushProvince `json:"provinces"`
	ExcludeProvinces []string       `json:"excludeProvinces"`
}

type PushProvince struct {
	Province      string   `json:"province"`
	Cities        []string `json:"cities"`
	ExcludeCities []string `json:"excludeCities"`
}

func (p *Push) setTarget(targetType TargetType) *Push {
	p.PushTarget.Target = int(targetType)
	return p
}

func (p *Push) setAlias(alias []string) *Push {
	p.PushTarget.Alias = alias
	return p
}

func (p *Push) setTags(tags []string) *Push {
	p.PushTarget.Tags = tags
	return p
}

func (p *Push) setRids(rids []string) *Push {
	p.PushTarget.Rids = rids
	return p
}

func (p *Push) setPushAreas(pushAreas PushAreas) *Push {
	p.PushTarget.PushAreas = &pushAreas
	return p
}
