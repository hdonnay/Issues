// Copyright 2015 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Issue is a client for reading and updating issues in a GitHub project issue tracker.

	usage: issue [-a] [-e] [-p owner/repo] <query>

Issue runs the query against the given project's issue tracker and
prints a table of matching issues, sorted by issue summary.
The default owner/repo is golang/go.

If multiple arguments are given as the query, issue joins them by
spaces to form a single issue search. These two commands are equivalent:

	issue assignee:rsc author:robpike
	issue "assignee:rsc author:robpike"

Searches are always limited to open issues.

If the query is a single number, issue prints that issue in detail,
including all comments.

Authentication

Issue expects to find a GitHub "personal access token" in
$HOME/.github-issue-token and will use that token to authenticate
to GitHub when reading or writing issue data.
A token can be created by visiting https://github.com/settings/tokens/new.
The token only needs the 'repo' scope checkbox, and optionally 'private_repo'
if you want to work with issue trackers for private repositories.
It does not need any other permissions.
The -token flag specifies an alternate file from which to read the token.

Acme Editor Integration

If the -a flag is specified, issue runs as a collection of acme windows
instead of a command-line tool. In this mode, the query is optional.
If no query is given, issue uses "state:open".

There are three kinds of acme windows: issue, issue creation, issue list,
search result, and milestone list.

The following text forms can be looked for (right clicked on)
and open a window (or navigate to an existing one).

	nnnn			issue #nnnn
	#nnnn			issue #nnnn
	all			the issue list
	milestone(s)		the milestone list
	<milestone-name>	the named milestone (e.g., Go1.5)

Executing "New" opens an issue creation window.

Executing "Search <query>" opens a new window showing the
results of that search.

Issue Window

An issue window, opened by loading an issue number,
displays full detail about an issue, a header followed by each comment.
For example:

	Title: time: Duration should implement fmt.Formatter
	State: closed
	Assignee: robpike
	Closed: 2015-01-08 05:20:00
	Labels: release-none repo-main size-m
	Milestone:
	URL: https://github.com/golang/go/issues/8786

	Reported by dsymonds (2014-09-21 23:02:50)

		It'd be nice if http://play.golang.org/p/KCnUQOPyol
		printed "[+3us]", which would require time.Duration
		implementing fmt.Formatter to get the '+' flag.

	Comment by rsc (2015-01-08 05:17:06)

		time must not depend on fmt.

Executing "Get" reloads the issue data.

Executing "Put" updates an issue. It saves any changes to the issue header
and, if any text has been entered between the header and the "Reported by" line,
posts that text as a new comment. If both succeed, Put then reloads the issue data.
The "Closed" and "URL" headers cannot be changed.

Issue Creation Window

An issue creation window, opened by executing "New", is like an issue window
but displays only an empty issue template:

	Title:
	Assignee:
	Labels:
	Milestone:

	<describe issue here>

Once the template has been completed (only the title is required), executing "Put"
creates the issue and converts the window into a issue window for the new issue.

Issue List Window

An issue list window displays a list of all open issue numbers and titles.
If the project has any open milestones, they are listed in a header line.
For example:

	Milestones: Go1.4.1 Go1.5 Go1.5Maybe

	9027	archive/tar: round-trip of Header misses values
	8669	archive/zip: not possible to a start writing zip at offset other than zero
	8359	archive/zip: not possible to specify deflate compression level
	...

As in any window, right clicking on an issue number opens a window for that issue.

Search Result Window

A search result window, opened by executing "Search <query>", displays a list of issues
matching a search query. It shows the query in a header line. For example:

	Search author:rsc

	9131	bench: no documentation
	599	cmd/5c, 5g, 8c, 8g: make 64-bit fields 64-bit aligned
	6699	cmd/5l: use m to store div/mod denominator
	4997	cmd/6a, cmd/8a: MOVL $x-8(SP) and LEAL x-8(SP) are different
	...

Executing "Sort" in a search result window toggles between sorting by title
and sorting by decreasing issue number.

Bulk Edit Window

Executing "Bulk" in an issue list or search result window opens a new
bulk edit window applying to the displayed issues. If there is a non-empty
text selection in the issue list or search result list, the bulk edit window
is restricted to issues in the selection.

The bulk edit window consists of a metadata header followed by a list of issues, like:

	State: open
	Assignee:
	Labels:
	Milestone: Go1.4.3

	10219	cmd/gc: internal compiler error: agen: unknown op
	9711	net/http: Testing timeout on Go1.4.1
	9576	runtime: crash in checkdead
	9954	runtime: invalid heap pointer found in bss on openbsd/386

The metadata header shows only metadata shared by all the issues.
In the above example, all four issues are open and have milestone Go1.4.3,
but they have no common labels nor a common assignee.

The bulk edit applies to the issues listed in the window text; adding or removing
issue lines changes the set of issues affected by Get or Put operations.

Executing "Get" refreshes the metadata header and issue summaries.

Executing "Put" updates all the listed issues. It applies any changes made to
the metadata header and, if any text has been entered between the header
and the first issue line, posts that text as a comment. If all operations succeed,
Put then refreshes the window as Get does.

Milestone List Window

The milestone list window, opened by loading any of the names
"milestone", "Milestone", or "Milestones", displays the open project
milestones, sorted by due date, along with the number of open issues in each.
For example:

	2015-01-15	Go1.4.1		1
	2015-07-31	Go1.5		215
	2015-07-31	Go1.5Maybe	5

Loading one of the listed milestone names opens a search for issues
in that milestone.

Alternate Editor Integration

The -e flag enables basic editing of issues with editors other than acme.
The editor invoked is $VISUAL if set, $EDITOR if set, or else ed.
Issue prepares a textual representation of issue data in a temporary file,
opens that file in the editor, waits for the editor to exit, and then applies any
changes from the file to the actual issues.

When <query> is a single number, issue -e edits a single issue.
See the ``Issue Window'' section above.

If the <query> is the text "new", issue -e creates a new issue.
See the ``Issue Creation Window'' section above.

Otherwise, for general queries, issue -e edits multiple issues in bulk.
See the ``Bulk Edit Window'' section above.

JSON Output

The -json flag causes issue to print the results in JSON format
using these data structures:

	type Issue struct {
		Number    int
		Ref       string
		Title     string
		State     string
		Assignee  string
		Closed    time.Time
		Labels    []string
		Milestone string
		URL       string
		Reporter  string
		Created   time.Time
		Text      string
		Comments  []*Comment
	}

	type Comment struct {
		Author string
		Time   time.Time
		Text   string
	}

If asked for a specific issue, the output is an Issue with Comments.
Otherwise, the result is an array of Issues without Comments.
*/
package main // import "rsc.io/github/issue"

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var (
	acmeFlag     = flag.Bool("a", false, "open in new acme window")
	acmeWrap     = flag.Int("w", 120, "wrap width for printing issues in acme windows")
	editFlag     = flag.Bool("e", false, "edit in system editor")
	jsonFlag     = flag.Bool("json", false, "write JSON output")
	project      = flag.String("p", "golang/go", "GitHub owner/repo name")
	rawFlag      = flag.Bool("raw", false, "do no processing of markdown")
	tokenFile    = flag.String("token", "", "read GitHub token personal access token from `file` (default $HOME/.github-issue-token)")
	apiRootArg   = flag.String("api", "", "base url for github instance (default github.com)")
	projectOwner = ""
	projectRepo  = ""

	apiRoot *url.URL
)

func usage() {
	fmt.Fprintf(os.Stderr, `usage: issue [-a] [-e] [-p owner/repo] <query>

If query is a single number, prints the full history for the issue.
Otherwise, prints a table of matching results.
`)
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	flag.Parse()
	log.SetFlags(0)
	log.SetPrefix("issue: ")

	if flag.NArg() == 0 && !*acmeFlag {
		usage()
	}

	if *jsonFlag && *acmeFlag {
		log.Fatal("cannot use -a with -json")
	}
	if *jsonFlag && *editFlag {
		log.Fatal("cannot use -e with -acme")
	}

	f := strings.Split(*project, "/")
	if len(f) != 2 {
		log.Fatal("invalid form for -p argument: must be owner/repo, like golang/go")
	}
	projectOwner = f[0]
	projectRepo = f[1]
	if *apiRootArg != "" {
		u, err := url.Parse(*apiRootArg)
		if err != nil {
			log.Fatal(err)
		}
		apiRoot = u
	}

	loadAuth()

	if *acmeFlag {
		acmeMode()
	}

	q := strings.Join(flag.Args(), " ")

	if *editFlag && q == "new" {
		editIssue([]byte(createTemplate), new(github.Issue))
		return
	}

	n, _ := strconv.Atoi(q)
	if n != 0 {
		if *editFlag {
			var buf bytes.Buffer
			issue, err := showIssue(&buf, n)
			if err != nil {
				log.Fatal(err)
			}
			editIssue(buf.Bytes(), issue)
			return
		}
		if _, err := showIssue(os.Stdout, n); err != nil {
			log.Fatal(err)
		}
		return
	}

	if *editFlag {
		all, err := searchIssues(q)
		if err != nil {
			log.Fatal(err)
		}
		if len(all) == 0 {
			log.Fatal("no issues matched search")
		}
		sort.Sort(issuesByTitle(all))
		bulkEditIssues(all)
		return
	}

	if err := showQuery(os.Stdout, q); err != nil {
		log.Fatal(err)
	}
}

func showIssue(w io.Writer, n int) (*github.Issue, error) {
	issue, _, err := client.Issues.Get(context.Background(), projectOwner, projectRepo, n)
	if err != nil {
		return nil, err
	}
	updateIssueCache(issue)
	return issue, printIssue(w, issue)
}

const timeFormat = "2006-01-02 15:04:05"

func printIssue(w io.Writer, issue *github.Issue) error {
	if *jsonFlag {
		showJSONIssue(w, issue)
		return nil
	}

	fmt.Fprintf(w, "Title: %s\n", getString(issue.Title))
	fmt.Fprintf(w, "State: %s\n", getString(issue.State))
	fmt.Fprintf(w, "Assignee: %s\n", getUserLogin(issue.Assignee))
	if issue.ClosedAt != nil {
		fmt.Fprintf(w, "Closed: %s\n", getTime(issue.ClosedAt).Format(timeFormat))
	}
	fmt.Fprintf(w, "Labels: %s\n", strings.Join(getLabelNames(issue.Labels), " "))
	fmt.Fprintf(w, "Milestone: %s\n", getMilestoneTitle(issue.Milestone))
	fmt.Fprintf(w, "URL: %s\n", issue.GetHTMLURL())

	fmt.Fprintf(w, "\nReported by %s (%s)\n", getUserLogin(issue.User), getTime(issue.CreatedAt).Format(timeFormat))
	if issue.Body != nil {
		if *rawFlag {
			fmt.Fprintf(w, "\n%s\n\n", *issue.Body)
		} else {
			text := strings.TrimSpace(*issue.Body)
			if text != "" {
				fmt.Fprintf(w, "\n\t%s\n", wrap(text, "\t"))
			}
		}
	}

	var output []string

	for page := 1; ; {
		list, resp, err := client.Issues.ListComments(context.Background(), projectOwner, projectRepo, getInt(issue.Number), &github.IssueListCommentsOptions{
			ListOptions: github.ListOptions{
				Page:    page,
				PerPage: 100,
			},
		})
		for _, com := range list {
			var buf bytes.Buffer
			w := &buf
			fmt.Fprintf(w, "%s\n", getTime(com.CreatedAt).Format(time.RFC3339))
			fmt.Fprintf(w, "\nComment by %s (%s)\n", getUserLogin(com.User), getTime(com.CreatedAt).Format(timeFormat))
			if com.Body != nil {
				if *rawFlag {
					fmt.Fprintf(w, "\n%s\n\n", *com.Body)
				} else {
					text := strings.TrimSpace(*com.Body)
					if text != "" {
						fmt.Fprintf(w, "\n\t%s\n", wrap(text, "\t"))
					}
				}
			}
			output = append(output, buf.String())
		}
		if err != nil {
			return err
		}
		if resp.NextPage < page {
			break
		}
		page = resp.NextPage
	}

	for page := 1; ; {
		list, resp, err := client.Issues.ListIssueEvents(context.Background(), projectOwner, projectRepo, getInt(issue.Number), &github.ListOptions{
			Page:    page,
			PerPage: 100,
		})
		for _, ev := range list {
			var buf bytes.Buffer
			w := &buf
			fmt.Fprintf(w, "%s\n", getTime(ev.CreatedAt).Format(time.RFC3339))
			switch event := getString(ev.Event); event {
			case "mentioned", "subscribed", "unsubscribed":
				// ignore
			case "added_to_project", "moved_columns_in_project", "removed_from_project":
				event = strings.Replace(event, "_", " ", -1)
				fallthrough
			default:
				fmt.Fprintf(w, "\n* %s %s (%s)\n", getUserLogin(ev.Actor), event, getTime(ev.CreatedAt).Format(timeFormat))
			case "closed", "referenced", "merged":
				id := getString(ev.CommitID)
				if id != "" {
					if len(id) > 7 {
						id = id[:7]
					}
					id = " in commit " + id
				}
				fmt.Fprintf(w, "\n* %s %s%s (%s)\n", getUserLogin(ev.Actor), event, id, getTime(ev.CreatedAt).Format(timeFormat))
				if id != "" {
					commit, _, err := client.Git.GetCommit(context.Background(), projectOwner, projectRepo, *ev.CommitID)
					if err == nil {
						fmt.Fprintf(w, "\n\tAuthor: %s <%s> %s\n\tCommitter: %s <%s> %s\n\n\t%s\n",
							getString(commit.Author.Name), getString(commit.Author.Email), getTime(commit.Author.Date).Format(timeFormat),
							getString(commit.Committer.Name), getString(commit.Committer.Email), getTime(commit.Committer.Date).Format(timeFormat),
							wrap(getString(commit.Message), "\t"))
					}
				}
			case "assigned", "unassigned":
				fmt.Fprintf(w, "\n* %s %s %s (%s)\n", getUserLogin(ev.Actor), event, getUserLogin(ev.Assignee), getTime(ev.CreatedAt).Format(timeFormat))
			case "labeled", "unlabeled":
				fmt.Fprintf(w, "\n* %s %s %s (%s)\n", getUserLogin(ev.Actor), event, getString(ev.Label.Name), getTime(ev.CreatedAt).Format(timeFormat))
			case "milestoned", "demilestoned":
				if event == "milestoned" {
					event = "added to milestone"
				} else {
					event = "removed from milestone"
				}
				fmt.Fprintf(w, "\n* %s %s %s (%s)\n", getUserLogin(ev.Actor), event, getString(ev.Milestone.Title), getTime(ev.CreatedAt).Format(timeFormat))
			case "renamed":
				fmt.Fprintf(w, "\n* %s changed title (%s)\n  - %s\n  + %s\n", getUserLogin(ev.Actor), getTime(ev.CreatedAt).Format(timeFormat), getString(ev.Rename.From), getString(ev.Rename.To))
			}
			output = append(output, buf.String())
		}
		if err != nil {
			return err
		}
		if resp.NextPage < page {
			break
		}
		page = resp.NextPage
	}

	sort.Strings(output)
	for _, s := range output {
		i := strings.Index(s, "\n")
		fmt.Fprintf(w, "%s", s[i+1:])
	}

	return nil
}

func showQuery(w io.Writer, q string) error {
	all, err := searchIssues(q)
	if err != nil {
		return err
	}
	sort.Sort(issuesByTitle(all))
	if *jsonFlag {
		showJSONList(all)
		return nil
	}
	for _, issue := range all {
		fmt.Fprintf(w, "%v\t%v\n", getInt(issue.Number), getString(issue.Title))
	}
	return nil
}

type issuesByTitle []*github.Issue

func (x issuesByTitle) Len() int      { return len(x) }
func (x issuesByTitle) Swap(i, j int) { x[i], x[j] = x[j], x[i] }
func (x issuesByTitle) Less(i, j int) bool {
	if getString(x[i].Title) != getString(x[j].Title) {
		return getString(x[i].Title) < getString(x[j].Title)
	}
	return getInt(x[i].Number) < getInt(x[j].Number)
}

func searchIssues(q string) ([]*github.Issue, error) {
	if opt, ok := queryToListOptions(q); ok {
		return listRepoIssues(opt)
	}

	var all []*github.Issue
	for page := 1; ; {
		// TODO(rsc): Rethink excluding pull requests.
		x, resp, err := client.Search.Issues(context.Background(), "type:issue state:open repo:"+*project+" "+q, &github.SearchOptions{
			ListOptions: github.ListOptions{
				Page:    page,
				PerPage: 100,
			},
		})
		for i := range x.Issues {
			updateIssueCache(&x.Issues[i])
			all = append(all, &x.Issues[i])
		}
		if err != nil {
			return all, err
		}
		if resp.NextPage < page {
			break
		}
		page = resp.NextPage
	}
	return all, nil
}

func queryToListOptions(q string) (opt github.IssueListByRepoOptions, ok bool) {
	if strings.ContainsAny(q, `"'`) {
		return
	}
	for _, f := range strings.Fields(q) {
		i := strings.Index(f, ":")
		if i < 0 {
			return
		}
		key, val := f[:i], f[i+1:]
		switch key {
		default:
			return
		case "milestone":
			if opt.Milestone != "" || val == "" {
				return
			}
			id := findMilestone(ioutil.Discard, &val)
			if id == nil {
				return
			}
			opt.Milestone = fmt.Sprint(*id)
		case "state":
			if opt.State != "" || val == "" {
				return
			}
			opt.State = val
		case "assignee":
			if opt.Assignee != "" || val == "" {
				return
			}
			opt.Assignee = val
		case "author":
			if opt.Creator != "" || val == "" {
				return
			}
			opt.Creator = val
		case "mentions":
			if opt.Mentioned != "" || val == "" {
				return
			}
			opt.Mentioned = val
		case "label":
			if opt.Labels != nil || val == "" {
				return
			}
			opt.Labels = strings.Split(val, ",")
		case "sort":
			if opt.Sort != "" || val == "" {
				return
			}
			opt.Sort = val
		case "updated":
			if !opt.Since.IsZero() || !strings.HasPrefix(val, ">=") {
				return
			}
			// TODO: Can set Since if we parse val[2:].
			return
		case "no":
			switch val {
			default:
				return
			case "milestone":
				if opt.Milestone != "" {
					return
				}
				opt.Milestone = "none"
			}
		}
	}
	return opt, true
}

func listRepoIssues(opt github.IssueListByRepoOptions) ([]*github.Issue, error) {
	var all []*github.Issue
	for page := 1; ; {
		xopt := opt
		xopt.ListOptions = github.ListOptions{
			Page:    page,
			PerPage: 100,
		}
		issues, resp, err := client.Issues.ListByRepo(context.Background(), projectOwner, projectRepo, &xopt)
		for i := range issues {
			updateIssueCache(issues[i])
			all = append(all, issues[i])
		}
		if err != nil {
			return all, err
		}
		if resp.NextPage < page {
			break
		}
		page = resp.NextPage
	}

	// Filter out pull requests, since we cannot say type:issue like in searchIssues.
	// TODO(rsc): Rethink excluding pull requests.
	save := all[:0]
	for _, issue := range all {
		if issue.PullRequestLinks == nil {
			save = append(save, issue)
		}
	}
	return save, nil
}

func loadMilestones() ([]*github.Milestone, error) {
	// NOTE(rsc): There appears to be no paging possible.
	all, _, err := client.Issues.ListMilestones(context.Background(), projectOwner, projectRepo, &github.MilestoneListOptions{
		State: "open",
	})
	if err != nil {
		return nil, err
	}
	if all == nil {
		all = []*github.Milestone{}
	}
	return all, nil
}

func wrap(t string, prefix string) string {
	out := ""
	t = strings.Replace(t, "\r\n", "\n", -1)
	max := 70
	if *acmeFlag {
		max = *acmeWrap
	}
	lines := strings.Split(t, "\n")
	for i, line := range lines {
		if i > 0 {
			out += "\n" + prefix
		}
		s := line
		for len(s) > max {
			i := strings.LastIndex(s[:max], " ")
			if i < 0 {
				i = max - 1
			}
			i++
			out += s[:i] + "\n" + prefix
			s = s[i:]
		}
		out += s
	}
	return out
}

var client *github.Client

// GitHub personal access token, from https://github.com/settings/applications.
var authToken string

func loadAuth() {
	const short = ".github-issue-token"
	filename := filepath.Clean(os.Getenv("HOME") + "/" + short)
	shortFilename := filepath.Clean("$HOME/" + short)
	if *tokenFile != "" {
		filename = *tokenFile
		shortFilename = *tokenFile
	}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal("reading token: ", err, "\n\n"+
			"Please create a personal access token at https://github.com/settings/tokens/new\n"+
			"and write it to ", shortFilename, " to use this program.\n"+
			"The token only needs the repo scope, or private_repo if you want to\n"+
			"view or edit issues for private repositories.\n"+
			"The benefit of using a personal access token over using your GitHub\n"+
			"password directly is that you can limit its use and revoke it at any time.\n\n")
	}
	fi, err := os.Stat(filename)
	if fi.Mode()&0077 != 0 {
		log.Fatalf("reading token: %s mode is %#o, want %#o", shortFilename, fi.Mode()&0777, fi.Mode()&0700)
	}
	authToken = strings.TrimSpace(string(data))
	t := &oauth2.Transport{
		Source: &tokenSource{AccessToken: authToken},
	}
	client = github.NewClient(&http.Client{Transport: t})
	client.BaseURL = apiRoot
	client.UploadURL = apiRoot
}

type tokenSource oauth2.Token

func (t *tokenSource) Token() (*oauth2.Token, error) {
	return (*oauth2.Token)(t), nil
}

func getInt(x *int) int {
	if x == nil {
		return 0
	}
	return *x
}

func getString(x *string) string {
	if x == nil {
		return ""
	}
	return *x
}

func getUserLogin(x *github.User) string {
	if x == nil || x.Login == nil {
		return ""
	}
	return *x.Login
}

func getTime(x *time.Time) time.Time {
	if x == nil {
		return time.Time{}
	}
	return (*x).Local()
}

func getMilestoneTitle(x *github.Milestone) string {
	if x == nil || x.Title == nil {
		return ""
	}
	return *x.Title
}

func getLabelNames(x []github.Label) []string {
	var out []string
	for _, lab := range x {
		out = append(out, getString(lab.Name))
	}
	sort.Strings(out)
	return out
}

var issueCache struct {
	sync.Mutex
	m map[int]*github.Issue
}

func updateIssueCache(issue *github.Issue) {
	n := getInt(issue.Number)
	if n == 0 {
		return
	}
	issueCache.Lock()
	if issueCache.m == nil {
		issueCache.m = make(map[int]*github.Issue)
	}
	issueCache.m[n] = issue
	issueCache.Unlock()
}

func bulkReadIssuesCached(ids []int) ([]*github.Issue, error) {
	var all []*github.Issue
	issueCache.Lock()
	for _, id := range ids {
		all = append(all, issueCache.m[id])
	}
	issueCache.Unlock()

	var errbuf bytes.Buffer
	for i, id := range ids {
		if all[i] == nil {
			issue, _, err := client.Issues.Get(context.Background(), projectOwner, projectRepo, id)
			if err != nil {
				fmt.Fprintf(&errbuf, "reading #%d: %v\n", id, err)
				continue
			}
			updateIssueCache(issue)
			all[i] = issue
		}
	}
	var err error
	if errbuf.Len() > 0 {
		err = fmt.Errorf("%s", strings.TrimSpace(errbuf.String()))
	}
	return all, err
}

// JSON output
// If you make changes to the structs, copy them back into the doc comment.

type Issue struct {
	Number    int
	Ref       string
	Title     string
	State     string
	Assignee  string
	Closed    time.Time
	Labels    []string
	Milestone string
	URL       string
	Reporter  string
	Created   time.Time
	Text      string
	Comments  []*Comment
}

type Comment struct {
	Author string
	Time   time.Time
	Text   string
}

func showJSONIssue(w io.Writer, issue *github.Issue) {
	data, err := json.MarshalIndent(toJSONWithComments(issue), "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	data = append(data, '\n')
	w.Write(data)
}

func showJSONList(all []*github.Issue) {
	j := []*Issue{} // non-nil for json
	for _, issue := range all {
		j = append(j, toJSON(issue))
	}
	data, err := json.MarshalIndent(j, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	data = append(data, '\n')
	os.Stdout.Write(data)
}

func toJSON(issue *github.Issue) *Issue {
	j := &Issue{
		Number:    getInt(issue.Number),
		Ref:       fmt.Sprintf("%s/%s#%d\n", projectOwner, projectRepo, getInt(issue.Number)),
		Title:     getString(issue.Title),
		State:     getString(issue.State),
		Assignee:  getUserLogin(issue.Assignee),
		Closed:    getTime(issue.ClosedAt),
		Labels:    getLabelNames(issue.Labels),
		Milestone: getMilestoneTitle(issue.Milestone),
		URL:       issue.GetHTMLURL(),
		Reporter:  getUserLogin(issue.User),
		Created:   getTime(issue.CreatedAt),
		Text:      getString(issue.Body),
		Comments:  []*Comment{},
	}
	if j.Labels == nil {
		j.Labels = []string{}
	}
	return j
}

func toJSONWithComments(issue *github.Issue) *Issue {
	j := toJSON(issue)
	for page := 1; ; {
		list, resp, err := client.Issues.ListComments(context.Background(), projectOwner, projectRepo, getInt(issue.Number), &github.IssueListCommentsOptions{
			ListOptions: github.ListOptions{
				Page:    page,
				PerPage: 100,
			},
		})
		if err != nil {
			log.Fatal(err)
		}
		for _, com := range list {
			j.Comments = append(j.Comments, &Comment{
				Author: getUserLogin(com.User),
				Time:   getTime(com.CreatedAt),
				Text:   getString(com.Body),
			})
		}
		if resp.NextPage < page {
			break
		}
		page = resp.NextPage
	}
	return j
}
