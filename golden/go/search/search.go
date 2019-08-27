// search contains the core functionality for searching for digests across a tile.
package search

import (
	"fmt"
	"net/url"
	"sort"
	"strings"
	"sync"

	"go.skia.org/infra/go/paramtools"
	"go.skia.org/infra/go/sklog"
	"go.skia.org/infra/go/tiling"
	"go.skia.org/infra/go/util"
	"go.skia.org/infra/golden/go/blame"
	"go.skia.org/infra/golden/go/diff"
	"go.skia.org/infra/golden/go/digest_counter"
	"go.skia.org/infra/golden/go/digesttools"
	"go.skia.org/infra/golden/go/indexer"
	"go.skia.org/infra/golden/go/summary"
	"go.skia.org/infra/golden/go/types"
)

const (
	// SORT_ASC indicates that we want to sort in ascending order.
	SORT_ASC = "asc"

	// SORT_DESC indicates that we want to sort in descending order.
	SORT_DESC = "desc"

	// MAX_ROW_DIGESTS is the maximum number of digests we'll compare against
	// before limiting the result to avoid overload.
	MAX_ROW_DIGESTS = 200

	// MAX_LIMIT is the maximum limit we will return.
	MAX_LIMIT = 200
)

// Point is a single point. Used in Trace.
type Point struct {
	X int `json:"x"` // The commit index [0-49].
	Y int `json:"y"`
	S int `json:"s"` // Status of the digest: 0 if the digest matches our search, 1-8 otherwise.
}

// Trace describes a single trace, used in Traces.
type Trace struct {
	Data   []Point           `json:"data"`  // One Point for each test result.
	ID     tiling.TraceId    `json:"label"` // The id of the trace. Keep the json as label to be compatible with dots-sk.
	Params map[string]string `json:"params"`
}

// DigestStatus is a digest and its status, used in Traces.
type DigestStatus struct {
	Digest types.Digest `json:"digest"`
	Status string       `json:"status"`
}

// Traces is info about a group of traces. Used in Digest.
type Traces struct {
	TileSize int            `json:"tileSize"`
	Traces   []Trace        `json:"traces"`  // The traces where this digest appears.
	Digests  []DigestStatus `json:"digests"` // The other digests that appear in Traces.
}

// DiffDigest is information about a digest different from the one in Digest.
type DiffDigest struct {
	Closest  *digesttools.Closest `json:"closest"`
	ParamSet map[string][]string  `json:"paramset"`
}

// Diff is only populated for digests that are untriaged?
// Might still be useful to find diffs to closest pos for a neg, and vice-versa.
// Will also be useful if we ever get a canonical trace or centroid.
type Diff struct {
	Diff float32 `json:"diff"` // The smaller of the Pos and Neg diff.

	// Either may be nil if there's no positive or negative to compare against.
	Pos *DiffDigest `json:"pos"`
	Neg *DiffDigest `json:"neg"`
	//Centroid *DiffDigest

	Blame *blame.BlameDistribution `json:"blame"`
}

// Digest's are returned from Search, one for each match to Query.
type Digest struct {
	Test     string              `json:"test"`
	Digest   string              `json:"digest"`
	Status   string              `json:"status"`
	ParamSet map[string][]string `json:"paramset"`
	Traces   *Traces             `json:"traces"`
	Diff     *Diff               `json:"diff"`
}

// CommitRange is a range of commits, starting at the git hash Begin and ending at End, inclusive.
//
// Currently unimplemented in search.
type CommitRange struct {
}

// TODO: filter within tests.
const (
	GROUP_TEST_MAX_COUNT = "count"
)

// Query is the query that Search understands.
type Query struct {
	// Diff metric to use.
	Metric string   `json:"metric"`
	Sort   string   `json:"sort"`
	Match  []string `json:"match"`

	// Blaming
	BlameGroupID string `json:"blame"`

	// Image classification
	Pos            bool `json:"pos"`
	Neg            bool `json:"neg"`
	Head           bool `json:"head"`
	Unt            bool `json:"unt"`
	IncludeIgnores bool `json:"include"`

	// URL encoded query string
	QueryStr string     `json:"query"`
	Query    url.Values `json:"-"`

	// URL encoded query string to select the right hand side of comparisons.
	RQueryStr string              `json:"rquery"`
	RQuery    paramtools.ParamSet `json:"-"`

	// Trybot support.
	IssueStr      string  `json:"issue"`
	Issue         int64   `json:"-"`
	PatchsetsStr  string  `json:"patchsets"` // Comma-separated list of patchsets.
	Patchsets     []int64 `json:"-"`
	IncludeMaster bool    `json:"master"` // Include digests also contained in master when searching code review issues.

	// Filtering.
	FCommitBegin string  `json:"fbegin"`     // Start commit
	FCommitEnd   string  `json:"fend"`       // End commit
	FRGBAMin     int32   `json:"frgbamin"`   // Min RGBA delta
	FRGBAMax     int32   `json:"frgbamax"`   // Max RGBA delta
	FDiffMax     float32 `json:"fdiffmax"`   // Max diff according to metric
	FGroupTest   string  `json:"fgrouptest"` // Op within grouped by test.
	FRef         bool    `json:"fref"`       // Only digests with reference.

	// Pagination.
	Offset int32 `json:"offset"`
	Limit  int32 `json:"limit"`

	// Do not include diffs in search.
	NoDiff bool `json:"nodiff"`
}

func (q *Query) IgnoreState() types.IgnoreState {
	is := types.ExcludeIgnoredTraces
	if q.IncludeIgnores {
		is = types.IncludeIgnoredTraces
	}
	return is
}

// SearchResponse is the standard search response. Depending on the query some fields
// might be empty, i.e. IssueDetails only makes sense if a trybot isssue was given in the query.
type SearchResponse struct {
	Digests       []*Digest
	Total         int32
	Commits       []*tiling.Commit
	IssueResponse *IssueResponse
}

// IssueResponse contains specific query responses when we search for a trybot issue. Currently
// it extends trybot.IssueDetails.
type IssueResponse struct {
	QueryPatchsets []int64
}

// excludeClassification returns true if the given label/status for a digest
// should be excluded based on the values in the query.
func (q *Query) excludeClassification(cl types.Label) bool {
	return ((cl == types.NEGATIVE) && !q.Neg) ||
		((cl == types.POSITIVE) && !q.Pos) ||
		((cl == types.UNTRIAGED) && !q.Unt)
}

// DigestSlice is a utility type for sorting slices of Digest by their max diff.
type DigestSlice []*Digest

func (p DigestSlice) Len() int           { return len(p) }
func (p DigestSlice) Less(i, j int) bool { return p[i].Diff.Diff > p[j].Diff.Diff }
func (p DigestSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// digestIndex returns the index of the digest d in digestInfo, or -1 if not found.
func digestIndex(d types.Digest, digestInfo []DigestStatus) int {
	for i, di := range digestInfo {
		if di.Digest == d {
			return i
		}
	}
	return -1
}

// blameGroupID takes a blame distribution with just indices of commits and
// returns an id for the blame group, which is just a string, the concatenated
// git hashes in commit time order.
func blameGroupID(b blame.BlameDistribution, commits []*tiling.Commit) string {
	ret := []string{}
	for _, index := range b.Freq {
		ret = append(ret, commits[index].Hash)
	}
	return strings.Join(ret, ":")
}

// digestsFromTrace returns all the digests in the given trace, controlled by
// 'head', and being robust to tallies not having been calculated for the
// trace.
func digestsFromTrace(id tiling.TraceId, tr *types.GoldenTrace, head bool, lastCommitIndex int, digestsByTrace map[tiling.TraceId]digest_counter.DigestCount) types.DigestSlice {
	digests := types.DigestSet{}
	if head {
		// Find the last non-missing value in the trace.
		for i := lastCommitIndex; i >= 0; i-- {
			if tr.IsMissing(i) {
				continue
			} else {
				digests[tr.Digests[i]] = true
				break
			}
		}
	} else {
		// Use the digestsByTrace if available, otherwise just inspect the trace.
		if t, ok := digestsByTrace[id]; ok {
			for k := range t {
				digests[k] = true
			}
		} else {
			for i := lastCommitIndex; i >= 0; i-- {
				if !tr.IsMissing(i) {
					digests[tr.Digests[i]] = true
				}
			}
		}
	}

	return digests.Keys()
}

// TODO(stephana): Replace digesttools.Closest here and above with an overall
// consolidated structure to measure the distance between two digests.

// SRDigestDiff contains the result of comparing two digests.
type SRDigestDiff struct {
	Left  *SRDigest     `json:"left"`  // The left hand digest and its params.
	Right *SRDiffDigest `json:"right"` // The right hand digest, its params and the diff result.
}

// CompareDigests compares two digests that were generated by the given test. It returns
// an instance of DigestDiff.
func (c *SearchAPI) CompareDigests(test types.TestName, left, right types.Digest) (*SRDigestDiff, error) {
	// Get the diff between the two digests
	diffResult, err := c.DiffStore.Get(diff.PRIORITY_NOW, left, types.DigestSlice{right})
	if err != nil {
		return nil, err
	}

	// Return an error if we could not find the diff.
	if len(diffResult) != 1 {
		return nil, fmt.Errorf("Could not find diff between %s and %s", left, right)
	}

	exp, err := c.ExpectationsStore.Get()
	if err != nil {
		return nil, err
	}

	idx := c.Indexer.GetIndex()

	return &SRDigestDiff{
		Left: &SRDigest{
			Test:     test,
			Digest:   left,
			Status:   exp.Classification(test, left).String(),
			ParamSet: idx.GetParamsetSummary(test, left, types.IncludeIgnoredTraces),
		},
		Right: &SRDiffDigest{
			Test:        test,
			Digest:      right,
			Status:      exp.Classification(test, right).String(),
			ParamSet:    idx.GetParamsetSummary(test, right, types.IncludeIgnoredTraces),
			DiffMetrics: diffResult[right].(*diff.DiffMetrics),
		},
	}, nil
}

// CTQuery is the input structure to the CompareTest function.
type CTQuery struct {
	// RowQuery is the query to select the row digests.
	RowQuery *Query `json:"rowQuery"`

	// ColumnQuery is the query to select the column digests.
	ColumnQuery *Query `json:"columnQuery"`

	// Match is the list of parameter fields where the column digests have to match
	// the value of the row digests. That means column digests will only be included
	// if the corresponding parameter values match the corresponding row digest.
	Match []string `json:"match"`

	// SortRows defines by what to sort the rows.
	SortRows string `json:"sortRows"`

	// SortColumns defines by what to sort the digest.
	SortColumns string `json:"sortColumns"`

	// RowsDir defines the sort direction for rows.
	RowsDir string `json:"rowsDir"`

	// ColumnsDir defines the sort direction for columns.
	ColumnsDir string `json:"columnsDir"`

	// Metric is the diff metric to use for sorting.
	Metric string `json:"metric"`
}

// CTResponse is the structure returned by the CompareTest.
type CTResponse struct {
	Grid      *CTGrid                       `json:"grid"`
	Name      string                        `json:"name"`
	Corpus    string                        `json:"source_type"`
	Summaries map[types.TestName]*CTSummary `json:"summaries"`
	Positive  int                           `json:"pos"`
	Negative  int                           `json:"neg"`
	Untriaged int                           `json:"unt"`
}

// CTGrid contains the grid of diff values returned by CompareTest.
type CTGrid struct {
	// Rows contains the row digest and the number of times it occurs.
	Rows []*CTRow `json:"rows"`

	// RowsTotal contains the total number of rows for the given query.
	RowsTotal int `json:"rowTotal"`

	// Columns contains the reference points calculated for each row digests.
	Columns []string `json:"columns"` // Contains the column types.

	// ColumnsTotal contains the total number of column digests.
	ColumnsTotal int `json:"columnsTotal"`
}

// CTDigestCount captures the digest and how often it appears in the tile.
type CTDigestCount struct {
	Digest types.Digest `json:"digest"`
	N      int          `json:"n"`
}

// CTRow is used by CTGrid to encode row digest information.
type CTRow struct {
	CTDigestCount
	TestName types.TestName   `json:"test"`
	Values   []*CTDiffMetrics `json:"values"`
}

// CTDiffMetrics contains diff metric between the contain digest and the
// corresponding row digest.
type CTDiffMetrics struct {
	*diff.DiffMetrics
	CTDigestCount
}

type CTSummary struct {
	Pos       int `json:"pos"`
	Neg       int `json:"neg"`
	Untriaged int `json:"untriaged"`
}

func ctSummaryFromSummary(sum *summary.Summary) *CTSummary {
	return &CTSummary{
		Pos:       sum.Pos,
		Neg:       sum.Neg,
		Untriaged: sum.Untriaged,
	}
}

// CompareTest allows to compare the digests within one test. It assumes that
// the provided instance of CTQuery is consistent in that the row query and
// column query contain the same test names and the same corpus field.
func (c *SearchAPI) CompareTest(ctq *CTQuery) (*CTResponse, error) {
	// Retrieve the row digests.
	rowDigests, err := c.filterTileCompare(ctq.RowQuery)
	if err != nil {
		return nil, err
	}
	totalRowDigests := len(rowDigests)

	// Build the rows output.
	rows := getCTRows(rowDigests, ctq.SortRows, ctq.RowsDir, ctq.RowQuery.Limit, ctq.RowQuery.IgnoreState(), c.Indexer.GetIndex())

	// If the number exceeds the maximum we always sort and trim by frequency.
	if len(rows) > MAX_ROW_DIGESTS {
		ctq.SortRows = SORT_FIELD_COUNT
	}

	// If we sort by image frequency then we can sort and limit now, reducing the
	// number of diffs we need to make.
	sortEarly := (ctq.SortRows == SORT_FIELD_COUNT)
	var uniqueTests types.TestNameSet = nil
	if sortEarly {
		uniqueTests = sortAndLimitRows(&rows, rowDigests, ctq.SortRows, ctq.RowsDir, ctq.Metric, ctq.RowQuery.Limit)
	}

	// Get the column digests conditioned on the result of the row digests.
	columnDigests, err := c.filterTileWithMatch(ctq.ColumnQuery, ctq.Match, rowDigests)
	if err != nil {
		return nil, err
	}

	// Compare the rows in parallel.
	var wg sync.WaitGroup
	wg.Add(len(rows))
	rowLenCh := make(chan int, len(rows))
	for idx, rowElement := range rows {
		go func(idx int, digest types.Digest) {
			defer wg.Done()
			var total int
			var err error
			rows[idx].Values, total, err = getDiffs(c.DiffStore, digest, columnDigests[digest].Keys(), ctq.ColumnsDir, ctq.Metric, ctq.ColumnQuery.Limit)
			if err != nil {
				sklog.Errorf("Unable to calculate diff of row for digest %s. Got error: %s", digest, err)
			}
			rowLenCh <- total
		}(idx, rowElement.Digest)
	}
	wg.Wait()

	// TODO(stephana): Add reference points (i.e. closest positive/negative, in trace)
	// to columns. Without these reference points the result only contains the
	// diff values.

	// Find the max length of rows and trim them if necessary.
	columns := []string{}
	columnsTotal := 0
	close(rowLenCh)
	for t := range rowLenCh {
		if t > columnsTotal {
			columnsTotal = t
		}
	}

	if !sortEarly {
		uniqueTests = sortAndLimitRows(&rows, rowDigests, ctq.SortRows, ctq.RowsDir, ctq.Metric, ctq.RowQuery.Limit)
	}

	// Get the summaries of all tests in the result.
	testSummaries := c.Indexer.GetIndex().GetSummaries(types.ExcludeIgnoredTraces)
	ctSummaries := make(map[types.TestName]*CTSummary, len(uniqueTests))
	for testName := range uniqueTests {
		ctSummaries[testName] = ctSummaryFromSummary(testSummaries[testName])
	}

	ret := &CTResponse{
		Grid: &CTGrid{
			Rows:         rows,
			RowsTotal:    totalRowDigests,
			Columns:      columns,
			ColumnsTotal: columnsTotal,
		},
		Corpus:    ctq.RowQuery.Query.Get(types.CORPUS_FIELD),
		Summaries: ctSummaries,
	}

	return ret, nil
}

// filterTileCompare iterates over the tile and finds digests that match the given query.
// It returns a map[digest]ParamSet which contains all the found digests and
// the paramsets that generated them.
func (c *SearchAPI) filterTileCompare(query *Query) (map[types.Digest]paramtools.ParamSet, error) {
	ret := map[types.Digest]paramtools.ParamSet{}

	// Add digest/trace to the result.
	addFn := func(test types.TestName, digest types.Digest, traceID tiling.TraceId, trace *types.GoldenTrace, acceptRet interface{}) {
		if found, ok := ret[digest]; ok {
			found.AddParams(trace.Params())
		} else {
			ret[digest] = paramtools.NewParamSet(trace.Params())
		}
	}

	exp, err := c.ExpectationsStore.Get()
	if err != nil {
		return nil, err
	}

	if err := iterTile(query, addFn, nil, ExpSlice{exp}, c.Indexer.GetIndex()); err != nil {
		return nil, err
	}
	return ret, nil
}

// paramsMatch Returns true if all the parameters listed in matchFields have matching values
// in condParamSets and params.
func paramsMatch(matchFields []string, condParamSets paramtools.ParamSet, params paramtools.Params) bool {
	for _, field := range matchFields {
		val, valOk := params[field]
		condVals, condValsOk := condParamSets[field]
		if !(valOk && condValsOk && util.In(val, condVals)) {
			return false
		}
	}
	return true
}

// filterTileWithMatch iterates over the tile and finds the digests that match
// the query and satisfy the condition of matching parameter values for the
// fields listed in matchFields. condDigests contains the digests their
// parameter sets for which we would like to find a set of digests for
// comparison. It returns a set of digests for each digest in condDigests.
func (c *SearchAPI) filterTileWithMatch(query *Query, matchFields []string, condDigests map[types.Digest]paramtools.ParamSet) (map[types.Digest]types.DigestSet, error) {
	if len(condDigests) == 0 {
		return map[types.Digest]types.DigestSet{}, nil
	}

	ret := make(map[types.Digest]types.DigestSet, len(condDigests))
	for d := range condDigests {
		ret[d] = types.DigestSet{}
	}

	// Define the acceptFn and addFn.
	var acceptFn AcceptFn = nil
	var addFn AddFn = nil
	if len(matchFields) >= 0 {
		matching := make(types.DigestSlice, 0, len(condDigests))
		acceptFn = func(params paramtools.Params, digests types.DigestSlice) (bool, interface{}) {
			matching = matching[:0]
			for digest, paramSet := range condDigests {
				if paramsMatch(matchFields, paramSet, params) {
					matching = append(matching, digest)
				}
			}
			return len(matching) > 0, matching
		}
		addFn = func(test types.TestName, digest types.Digest, traceID tiling.TraceId, trace *types.GoldenTrace, acceptRet interface{}) {
			for _, d := range acceptRet.(types.DigestSlice) {
				ret[d][digest] = true
			}
		}
	} else {
		addFn = func(test types.TestName, digest types.Digest, traceID tiling.TraceId, trace *types.GoldenTrace, acceptRet interface{}) {
			for d := range condDigests {
				ret[d][digest] = true
			}
		}
	}

	exp, err := c.ExpectationsStore.Get()
	if err != nil {
		return nil, err
	}

	if err := iterTile(query, addFn, acceptFn, ExpSlice{exp}, c.Indexer.GetIndex()); err != nil {
		return nil, err
	}
	return ret, nil
}

// getCTRows returns the instance of CTRow that correspond to the given set of row digests.
func getCTRows(entries map[types.Digest]paramtools.ParamSet, sortField, sortDir string, limit int32, is types.IgnoreState, idx indexer.IndexSearcher) []*CTRow {
	talliesByTest := idx.DigestCountsByTest(is)
	ret := make([]*CTRow, 0, len(entries))
	for digest, paramSet := range entries {
		testName := types.TestName(paramSet[types.PRIMARY_KEY_FIELD][0])
		ret = append(ret, &CTRow{
			TestName: testName,
			CTDigestCount: CTDigestCount{
				Digest: digest,
				N:      talliesByTest[testName][digest],
			},
		})
	}
	return ret
}

// sortAndLimitRows sorts the given rows based on field, direction and diffMetric (if sorted by
// by diff). After the sort it will slice the result to be not larger than limit.
func sortAndLimitRows(rows *[]*CTRow, rowDigests map[types.Digest]paramtools.ParamSet, field, direction string, diffMetric string, limit int32) types.TestNameSet {
	// Determine the less function used for sorting the rows.
	var lessFn ctRowSliceLessFn
	if field == SORT_FIELD_COUNT {
		lessFn = func(c *ctRowSlice, i, j int) bool { return c.data[i].N < c.data[j].N }
	} else {
		lessFn = func(c *ctRowSlice, i, j int) bool {
			return (len(c.data[i].Values) > 0) && (len(c.data[j].Values) > 0) && (c.data[i].Values[0].Diffs[diffMetric] < c.data[j].Values[0].Diffs[diffMetric])
		}
	}

	sortSlice := sort.Interface(newCTRowSlice(*rows, lessFn))
	if direction == SORT_DESC {
		sortSlice = sort.Reverse(sortSlice)
	}

	sort.Sort(sortSlice)
	lastIdx := util.MinInt32(limit, int32(len(*rows)))
	discarded := (*rows)[lastIdx:]
	for _, row := range discarded {
		delete(rowDigests, row.Digest)
	}
	*rows = (*rows)[:lastIdx]

	uniqueTests := types.TestNameSet{}
	for _, paramSets := range rowDigests {
		for _, t := range paramSets[types.PRIMARY_KEY_FIELD] {
			uniqueTests[types.TestName(t)] = true
		}
	}
	return uniqueTests
}

// Sort adapter to allow sorting rows by supplying a less function.
type ctRowSliceLessFn func(c *ctRowSlice, i, j int) bool
type ctRowSlice struct {
	lessFn ctRowSliceLessFn
	data   []*CTRow
}

func newCTRowSlice(data []*CTRow, lessFn ctRowSliceLessFn) *ctRowSlice {
	return &ctRowSlice{lessFn: lessFn, data: data}
}
func (c *ctRowSlice) Len() int           { return len(c.data) }
func (c *ctRowSlice) Less(i, j int) bool { return c.lessFn(c, i, j) }
func (c *ctRowSlice) Swap(i, j int)      { c.data[i], c.data[j] = c.data[j], c.data[i] }

// getDiffs gets the sorted and limited comparison of one digest against the list of digests.
// Arguments:
//    digest: primary digest
//    colDigests: the digests to compare against
//    sortDir: sort direction of the resulting list
//    diffMetric: id of the diffmetric to use (assumed to be defined in the diff package).
//    limit: is the maximum number of diffs to return after the sort.
func getDiffs(diffStore diff.DiffStore, digest types.Digest, colDigests types.DigestSlice, sortDir, diffMetric string, limit int32) ([]*CTDiffMetrics, int, error) {
	diffMap, err := diffStore.Get(diff.PRIORITY_NOW, digest, colDigests)
	if err != nil {
		return nil, 0, err
	}

	ret := make([]*CTDiffMetrics, 0, len(diffMap))
	for colDigest, diffMetrics := range diffMap {
		ret = append(ret, &CTDiffMetrics{
			DiffMetrics:   diffMetrics.(*diff.DiffMetrics),
			CTDigestCount: CTDigestCount{Digest: colDigest, N: 0},
		})
	}

	// TODO(stephana): Add the reference points for each row.

	lessFn := func(c *ctDiffMetricsSlice, i, j int) bool {
		return c.data[i].Diffs[diffMetric] < c.data[j].Diffs[diffMetric]
	}
	sortSlice := sort.Interface(newCTDiffMetricsSlice(ret, lessFn))
	if sortDir == SORT_DESC {
		sortSlice = sort.Reverse(sortSlice)
	}
	sort.Sort(sortSlice)
	return ret[:util.MinInt(int(limit), len(ret))], len(ret), nil
}

// Sort adapter to allow sorting lists of diff metrics via a less function.
type ctDiffMetricsSliceLessFn func(c *ctDiffMetricsSlice, i, j int) bool
type ctDiffMetricsSlice struct {
	lessFn ctDiffMetricsSliceLessFn
	data   []*CTDiffMetrics
}

func newCTDiffMetricsSlice(data []*CTDiffMetrics, lessFn ctDiffMetricsSliceLessFn) *ctDiffMetricsSlice {
	return &ctDiffMetricsSlice{lessFn: lessFn, data: data}
}
func (c *ctDiffMetricsSlice) Len() int           { return len(c.data) }
func (c *ctDiffMetricsSlice) Less(i, j int) bool { return c.lessFn(c, i, j) }
func (c *ctDiffMetricsSlice) Swap(i, j int)      { c.data[i], c.data[j] = c.data[j], c.data[i] }
