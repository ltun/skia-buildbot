package perfresults

import "encoding/json"

// PerfResults represents the contenst of a perf_results.json file generated by a
// telemetry-based benchmark. The full format is not formally defined, but some
// documnentation for it exists in various places.  The most comprehensive doc is
// https://chromium.googlesource.com/external/github.com/catapult-project/catapult/+/HEAD/docs/Histogram-set-json-format.md
type PerfResults struct {
	Histograms      []Histogram
	GenericSets     []GenericSet
	DateRanges      []DateRange
	RelatedNameMaps []RelatedNameMap
}

// NonEmptyHistogramNames returns a list of names of histograms whose SampleValues arrays are non-empty.
func (pr *PerfResults) NonEmptyHistogramNames() []string {
	ret := []string{}
	for _, h := range pr.Histograms {
		if len(h.SampleValues) > 0 {
			ret = append(ret, h.Name)
		}
	}
	return ret
}

// Histogram is an individual benchmark measurement.
type Histogram struct {
	Name string `json:"name"`
	Unit string `json:"unit"`

	// optional fields
	Description  string    `json:"description"`
	SampleValues []float64 `json:"sampleValues"`
	// Diagnostics maps a diagnostic key to a guid, which points to e.g. a genericSet.
	Diagnostics map[string]string `json:"diagnostics"`
}

// GenericSet is a normalized value that other parts of the json file can reference by guid.
type GenericSet struct {
	GUID   string `json:"guid"`
	Values []any  `json:"values"` // Can be string or number. sigh.
}

// DateRange is a range of dates.
type DateRange struct {
	GUID string  `json:"guid"`
	Min  float64 `json:"min"`
	Max  float64 `json:"max"`
}

// RelatedNameMap is a map from short names to full histogram names.
type RelatedNameMap struct {
	GUID  string            `json:"guid"`
	Names map[string]string `json:"names"`
}

// UnmarshalJSON parses a byte slice into a PerfResults instance.
func (pr *PerfResults) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	for _, m := range raw {
		h := Histogram{}
		if err := json.Unmarshal(m, &h); err != nil {
			return err
		}
		if h.Name != "" {
			pr.Histograms = append(pr.Histograms, h)
			continue
		}

		type diagnostic struct {
			Type string `json:"type"`
			GUID string `json:"guid"`
		}

		d := diagnostic{}
		if err := json.Unmarshal(m, &d); err != nil {
			return err
		}
		if d.Type == "GenericSet" {
			gs := GenericSet{}
			if err := json.Unmarshal(m, &gs); err != nil {
				return err
			}
			pr.GenericSets = append(pr.GenericSets, gs)
			continue
		}
		if d.Type == "DateRange" {
			dr := DateRange{}
			if err := json.Unmarshal(m, &dr); err != nil {
				return err
			}
			pr.DateRanges = append(pr.DateRanges, dr)
		}
		if d.Type == "RelatedNameMap" {
			rnm := RelatedNameMap{}
			if err := json.Unmarshal(m, &rnm); err != nil {
				return err
			}
			pr.RelatedNameMaps = append(pr.RelatedNameMaps, rnm)
		}
	}
	return nil
}