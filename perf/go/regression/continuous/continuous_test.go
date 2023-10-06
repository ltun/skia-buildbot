package continuous

import (
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.skia.org/infra/go/paramtools"
	"go.skia.org/infra/go/testutils"
	"go.skia.org/infra/perf/go/alerts"
	"go.skia.org/infra/perf/go/clustering2"
	"go.skia.org/infra/perf/go/config"
	"go.skia.org/infra/perf/go/dataframe"
	"go.skia.org/infra/perf/go/dataframe/mocks"
	gitmocks "go.skia.org/infra/perf/go/git/mocks"
	"go.skia.org/infra/perf/go/git/provider"
	notifymocks "go.skia.org/infra/perf/go/notify/mocks"
	"go.skia.org/infra/perf/go/regression"
	regressionmocks "go.skia.org/infra/perf/go/regression/mocks"
	shortcutmocks "go.skia.org/infra/perf/go/shortcut/mocks"
	"go.skia.org/infra/perf/go/stepfit"
	"go.skia.org/infra/perf/go/types"
	"go.skia.org/infra/perf/go/ui/frame"
)

func TestBuildConfigsAndParamSet(t *testing.T) {
	c := Continuous{
		provider: func(_ context.Context) ([]*alerts.Alert, error) {
			// Only fill in ID since we are just testing if ch channel returns
			// what we set here.
			return []*alerts.Alert{
				{
					IDAsString: "1",
				},
				{
					IDAsString: "3",
				},
			}, nil
		},
		paramsProvider: func() paramtools.ReadOnlyParamSet {
			return paramtools.ReadOnlyParamSet{
				"config": []string{"8888", "565"},
			}
		},
		pollingDelay: time.Nanosecond,
		instanceConfig: &config.InstanceConfig{
			DataStoreConfig: config.DataStoreConfig{},
			GitRepoConfig:   config.GitRepoConfig{},
			IngestionConfig: config.IngestionConfig{},
		},
		flags: &config.FrontendFlags{},
	}

	// Build channel.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch := c.buildConfigAndParamsetChannel(ctx)

	// Read value.
	cnp := <-ch

	// Confirm it conforms to expectations.
	assert.Equal(t, c.paramsProvider(), cnp.paramset)
	assert.Len(t, cnp.configs, 2)
	ids := []string{}
	for _, cfg := range cnp.configs {
		ids = append(ids, cfg.IDAsString)
	}
	assert.Subset(t, []string{"1", "3"}, ids)

	// Confirm we continue to get items from the channel.
	cnp = <-ch
	assert.Equal(t, c.paramsProvider(), cnp.paramset)
}

func TestMatchingConfigsFromTraceIDs_TraceIDSliceIsEmpty_ReturnsEmptySlice(t *testing.T) {
	config := alerts.NewConfig()
	config.Query = "foo=bar"
	traceIDs := []string{}
	matchingConfigs := matchingConfigsFromTraceIDs(traceIDs, []*alerts.Alert{config})
	require.Empty(t, matchingConfigs)
}

func TestMatchingConfigsFromTraceIDs_OneConfigThatMatchesZeroTraces_ReturnsEmptySlice(t *testing.T) {
	config := alerts.NewConfig()
	config.Query = "arch=some-unknown-arch"
	traceIDs := []string{
		",arch=x86,config=8888,",
		",arch=arm,config=8888,",
	}
	matchingConfigs := matchingConfigsFromTraceIDs(traceIDs, []*alerts.Alert{config})
	require.Empty(t, matchingConfigs)
}

func TestMatchingConfigsFromTraceIDs_OneConfigThatMatchesOneTrace_ReturnsTheOneConfig(t *testing.T) {
	config := alerts.NewConfig()
	config.Query = "arch=x86"
	traceIDs := []string{
		",arch=x86,config=8888,",
		",arch=arm,config=8888,",
	}
	matchingConfigs := matchingConfigsFromTraceIDs(traceIDs, []*alerts.Alert{config})
	require.Len(t, matchingConfigs, 1)
}

func TestMatchingConfigsFromTraceIDs_TwoConfigsThatMatchesOneTrace_ReturnsBothConfigs(t *testing.T) {
	config1 := alerts.NewConfig()
	config1.Query = "arch=x86"
	config2 := alerts.NewConfig()
	config2.Query = "arch=arm"
	traceIDs := []string{
		",arch=x86,config=8888,",
		",arch=arm,config=8888,",
	}
	matchingConfigs := matchingConfigsFromTraceIDs(traceIDs, []*alerts.Alert{config1, config2})
	require.Len(t, matchingConfigs, 2)
}

func TestMatchingConfigsFromTraceIDs_GroupByMatchesTrace_ReturnsConfigWithRestrictedQuery(t *testing.T) {
	config1 := alerts.NewConfig()
	config1.Query = "arch=x86"
	config1.GroupBy = "config"
	traceIDs := []string{
		",arch=x86,config=8888,",
		",arch=arm,config=8888,",
	}
	matchingConfigs := matchingConfigsFromTraceIDs(traceIDs, []*alerts.Alert{config1})
	require.Len(t, matchingConfigs, 1)
	require.Equal(t, "arch=x86&config=8888", matchingConfigs[0].Query)
	_, err := url.ParseQuery(matchingConfigs[0].Query)
	require.NoError(t, err)
}

func TestMatchingConfigsFromTraceIDs_MultipleGroupByPartsMatchTrace_ReturnsConfigWithRestrictedQueryUsingAllMatchingGroupByKeys(t *testing.T) {
	config := alerts.NewConfig()
	config.Query = "arch=x86"
	config.GroupBy = "config,device"
	traceIDs := []string{
		",arch=x86,config=8888,device=Pixel4,",
		",arch=arm,config=8888,device=Pixel4,",
	}
	matchingConfigs := matchingConfigsFromTraceIDs(traceIDs, []*alerts.Alert{config})
	require.Len(t, matchingConfigs, 1)
	require.Equal(t, "arch=x86&config=8888&device=Pixel4", matchingConfigs[0].Query)
	_, err := url.ParseQuery(matchingConfigs[0].Query)
	require.NoError(t, err)
}

type allMocks struct {
	perfGit          *gitmocks.Git
	shortcutStore    *shortcutmocks.Store
	regressionStore  *regressionmocks.Store
	notifier         *notifymocks.Notifier
	dataFrameBuilder *mocks.DataFrameBuilder
}

func createArgsForReportRegressions(t *testing.T) (*Continuous, *regression.RegressionDetectionRequest, []*regression.RegressionDetectionResponse, *alerts.Alert, allMocks) {
	pg := gitmocks.NewGit(t)
	ss := shortcutmocks.NewStore(t)
	rs := regressionmocks.NewStore(t)
	cp := func(ctx context.Context) ([]*alerts.Alert, error) {
		return nil, nil
	}
	n := notifymocks.NewNotifier(t)
	pp := func() paramtools.ReadOnlyParamSet {
		return nil
	}
	dfb := mocks.NewDataFrameBuilder(t)
	i := &config.InstanceConfig{}
	f := &config.FrontendFlags{}

	req := &regression.RegressionDetectionRequest{}
	resp := []*regression.RegressionDetectionResponse{}
	cfg := &alerts.Alert{}

	c := &Continuous{
		perfGit:        pg,
		shortcutStore:  ss,
		store:          rs,
		provider:       cp,
		notifier:       n,
		paramsProvider: pp,
		dfBuilder:      dfb,
		pollingDelay:   time.Microsecond,
		instanceConfig: i,
		flags:          f,
		current: &Current{
			Alert: &alerts.Alert{},
		},
	}

	allMocks := allMocks{
		perfGit:          pg,
		shortcutStore:    ss,
		regressionStore:  rs,
		notifier:         n,
		dataFrameBuilder: dfb,
	}

	return c, req, resp, cfg, allMocks

}

func TestReportRegressions_EmptyRegressionDetectionResponse_NoRegressionsReported(t *testing.T) {
	c, req, resp, cfg, _ := createArgsForReportRegressions(t)
	// We know this works since we didn't need to supply any implementations for any of the mocks.
	c.reportRegressions(context.Background(), req, resp, cfg)
}

func TestReportRegressions_OneNewStepDownRegressionFound_OneRegressionStoredAndNotified(t *testing.T) {
	ctx := context.Background()
	c, req, resp, cfg, allMocks := createArgsForReportRegressions(t)

	const regressionCommitNumber = types.CommitNumber(2)
	resp = append(resp, &regression.RegressionDetectionResponse{
		Frame: &frame.FrameResponse{
			DataFrame: &dataframe.DataFrame{
				Header: []*dataframe.ColumnHeader{
					{Offset: 1},
					{Offset: regressionCommitNumber},
				},
				ParamSet: paramtools.ReadOnlyParamSet{
					"device_name": []string{"sailfish", "sargo", "wembley"},
				},
			},
		},
		Summary: &clustering2.ClusterSummaries{
			Clusters: []*clustering2.ClusterSummary{
				{
					Keys: []string{
						",device_name=sailfish",
						",device_name=sargo",
						",device_name=wembley",
					},
					Shortcut: "some-shortcut-id",
					StepFit: &stepfit.StepFit{
						Status: stepfit.LOW,
					},
					StepPoint: &dataframe.ColumnHeader{
						Offset: regressionCommitNumber,
					},
				},
			},
		},
	})

	commitAtStep := provider.Commit{
		Subject: "The subject of the commit where a regression occurred.",
	}
	previousCommit := provider.Commit{
		Subject: "The subject of the commit right before where a regression occurred.",
	}

	const notificationID = "some-notification-id"

	// First call to CommitFromCommitNumber.
	allMocks.perfGit.On("CommitFromCommitNumber", testutils.AnyContext, types.CommitNumber(2)).Return(commitAtStep, nil)

	// First call to CommitFromCommitNumber is for the previous commit.
	allMocks.perfGit.On("CommitFromCommitNumber", testutils.AnyContext, types.CommitNumber(1)).Return(previousCommit, nil)
	cfg.DirectionAsString = alerts.DOWN

	// Returns true to indicate that this is a newly found regression. Note that
	// this is called twice, first to store the regression since it's new, then
	// called again to store the notification ID.
	allMocks.regressionStore.On("SetLow", testutils.AnyContext, regressionCommitNumber, cfg.IDAsString, resp[0].Frame, resp[0].Summary.Clusters[0]).Return(true, nil).Twice()

	allMocks.notifier.On("RegressionFound", ctx, commitAtStep, previousCommit, cfg, resp[0].Summary.Clusters[0], resp[0].Frame).Return(notificationID, nil)

	c.reportRegressions(ctx, req, resp, cfg)

	require.Equal(t, notificationID, resp[0].Summary.Clusters[0].NotificationID)
}
