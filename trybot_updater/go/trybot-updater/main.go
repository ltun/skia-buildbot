// trybot_updater is an application that updates a repo's buildbucket.config file.
package main

import (
	"bytes"
	"context"
	"flag"
	"sort"
	"strings"
	"text/template"
	"time"

	"go.skia.org/infra/go/auth"
	"go.skia.org/infra/go/common"
	"go.skia.org/infra/go/gerrit"
	"go.skia.org/infra/go/gitiles"
	"go.skia.org/infra/go/httputils"
	"go.skia.org/infra/go/skerr"
	"go.skia.org/infra/go/sklog"
	"go.skia.org/infra/go/util"
	"go.skia.org/infra/task_scheduler/go/specs"
)

const (
	// The format of this file is that of a gerrit extension config (not a proto).
	// The buildbucket extension parses the config like this:
	// https://chromium.googlesource.com/infra/gerrit-plugins/buildbucket/+/refs/heads/master/src/main/java/com/googlesource/chromium/plugins/buildbucket/GetConfig.java
	bbCfgFileName = "buildbucket.config"
	// Branch buildbucket.config lies in. Hopefully this will change one day, see b/38258213.
	bbCfgBranch = "refs/meta/config"
	// Template used to create buildbucket.config content.
	bbCfgTemplate = `
{{- range .EmptyBuckets}}[bucket "{{.}}"]{{end}}
[bucket "{{.BucketName}}"]
{{- range .Jobs}}
	builder = {{.}}
{{- end}}
`
)

var (
	// Flags.
	repoUrl       = flag.String("repo_url", common.REPO_SKIA, "Repo that needs buildbucket.config updated from it's tasks.json file.")
	bucketName    = flag.String("bucket_name", "luci.skia.skia.primary", "Name of the bucket to update in buildbucket.config.")
	emptyBuckets  = common.NewMultiStringFlag("empty_bucket", nil, "Empty buckets to specify in buildbucket.config. Eg: luci.chromium.try. See skbug.com/9639 for why these buckets are empty.")
	pollingPeriod = flag.Duration("polling_period", 10*time.Minute, "How often to poll tasks.json.")
	submit        = flag.Bool("submit", false, "If set, automatically submit the Gerrit change to update buildbucket.config")
	local         = flag.Bool("local", false, "Running locally if true. As opposed to in production.")
	promPort      = flag.String("prom_port", ":20000", "Metrics service address (e.g., ':20000')")

	bbCfgTemplateParsed = template.Must(template.New("buildbucket_config").Parse(bbCfgTemplate))
)

// getBuildbucketCfgFromJobs reads tasks.json from the specified repository and returns
// contents of what the new buildbucket.config file should be.
func getBuildbucketCfgFromJobs(ctx context.Context, repo *gitiles.Repo) (string, error) {
	// Read tasks.json from the specified repository.
	tasksContents, err := repo.ReadFileAtRef(ctx, specs.TASKS_CFG_FILE, "master")
	if err != nil {
		return "", skerr.Fmt("Could not read %s: %s", specs.TASKS_CFG_FILE, err)
	}
	tasksCfg, err := specs.ParseTasksCfg(string(tasksContents))
	if err != nil {
		return "", skerr.Fmt("Could not parse %s: %s", specs.TASKS_CFG_FILE, err)
	}

	// Create a sorted slice of jobs.
	jobs := make([]string, 0, len(tasksCfg.Jobs))
	for j := range tasksCfg.Jobs {
		jobs = append(jobs, j)
	}
	sort.Strings(jobs)

	// Use jobs to create content of buildbucket.config.
	bbCfg := new(bytes.Buffer)
	if err := bbCfgTemplateParsed.Execute(bbCfg, struct {
		EmptyBuckets []string
		BucketName   string
		Jobs         []string
	}{
		EmptyBuckets: *emptyBuckets,
		BucketName:   *bucketName,
		Jobs:         jobs,
	}); err != nil {
		return "", skerr.Fmt("Failed to execute bbCfg template: %s", err)
	}
	return bbCfg.String(), nil
}

// getCurrentBuildbucketCfg returns the current contents of buildbucket.config for the
// specified repository.
func getCurrentBuildbucketCfg(ctx context.Context, repo *gitiles.Repo) (string, error) {
	contents, err := repo.ReadFileAtRef(ctx, bbCfgFileName, bbCfgBranch)
	if err != nil {
		return "", skerr.Fmt("Could not read %s: %s", bbCfgFileName, err)
	}
	return string(contents), nil
}

// updateBuildbucketCfg creates a Gerrit CL to update buildbucket.config. If submit flag is true then that CL
// is automatically self-approved and submitted.
func updateBuildbucketCfg(ctx context.Context, g *gerrit.Gerrit, repo *gitiles.Repo, cfgContents string) error {
	commitMsg := "Update buildbucket.config"
	repoSplit := strings.Split(*repoUrl, "/")
	project := strings.TrimSuffix(repoSplit[len(repoSplit)-1], ".git")
	baseCommitInfo, err := repo.Details(ctx, bbCfgBranch)
	if err != nil {
		return skerr.Fmt("Could not get details of %s: %s", bbCfgBranch, err)
	}
	baseCommit := baseCommitInfo.Hash
	ci, err := gerrit.CreateAndEditChange(ctx, g, project, bbCfgBranch, commitMsg, baseCommit, func(ctx context.Context, g gerrit.GerritInterface, ci *gerrit.ChangeInfo) error {
		if err := g.EditFile(ctx, ci, bbCfgFileName, cfgContents); err != nil {
			return skerr.Fmt("Could not edit %s: %s", bbCfgFileName, err)
		}
		return nil
	})
	if err != nil {
		return skerr.Fmt("Could not create Gerrit change: %s", err)
	}
	sklog.Infof("Uploaded change https://skia-review.googlesource.com/c/%s/+/%d", project, ci.Issue)

	if *submit {
		// TODO(rmistry): Change reviewer to be the trooper after verifying that things work.
		reviewers := []string{"rmistry@google.com"}
		if err := g.SetReview(ctx, ci, "", gerrit.CONFIG_CHROMIUM.SelfApproveLabels, reviewers); err != nil {
			return abandonGerritChange(ctx, g, ci, err)
		}
		if err := g.Submit(ctx, ci); err != nil {
			return abandonGerritChange(ctx, g, ci, err)
		}
		sklog.Infof("Submitted change https://skia-review.googlesource.com/c/%s/+/%d", project, ci.Issue)
	}

	return nil
}

// abandonGerritChange abandons the specified CL and returns the specified err after abandoning.
// If abandoning fails then that error is wrapped with the specified err.
func abandonGerritChange(ctx context.Context, g *gerrit.Gerrit, issue *gerrit.ChangeInfo, err error) error {
	if abandonErr := g.Abandon(ctx, issue, ""); abandonErr != nil {
		return skerr.Wrapf(err, "failed to abandon CL with %s", abandonErr)
	}
	return err
}

func main() {
	common.InitWithMust("trybot_updater", common.PrometheusOpt(promPort))
	defer sklog.Flush()
	ctx := context.Background()

	// OAuth2.0 TokenSource.
	ts, err := auth.NewDefaultTokenSource(false, auth.SCOPE_USERINFO_EMAIL, auth.SCOPE_GERRIT)
	if err != nil {
		sklog.Fatal(err)
	}
	// Authenticated HTTP client.
	httpClient := httputils.DefaultClientConfig().WithTokenSource(ts).With2xxOnly().Client()

	// Instantiate Gerrit.
	gUrl := strings.Split(*repoUrl, ".googlesource.com")[0] + "-review.googlesource.com"
	g, err := gerrit.NewGerrit(gUrl, httpClient)
	if err != nil {
		sklog.Fatal(err)
	}

	// Instantiate Gitiles using the specified repo URL.
	repo := gitiles.NewRepo(*repoUrl, httpClient)

	// TODO(rmistry): Use pubsub instead of polling.
	go util.RepeatCtx(ctx, *pollingPeriod, func(ctx context.Context) {
		existingCfg, err := getCurrentBuildbucketCfg(ctx, repo)
		if err != nil {
			sklog.Errorf("Could not get contents of buildbucket.config from %s: %s", repo, err)
		}

		newCfg, err := getBuildbucketCfgFromJobs(ctx, repo)
		if err != nil {
			sklog.Errorf("Could not get list of jobs from %s: %s", repo, err)
		}

		// Only update buildbucket.config if the config is different.
		if newCfg != existingCfg {
			if err := updateBuildbucketCfg(ctx, g, repo, newCfg); err != nil {
				sklog.Errorf("Could not update buildbucket.config: %s", err)
				sklog.Info("Sleep for 10 mins since there was a error with Gerrit to give it time to recover.")
				time.Sleep(10 * time.Minute)
			}
		} else {
			sklog.Info("Config has not changed")
		}
	})

	select {}
}
