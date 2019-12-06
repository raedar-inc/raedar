package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

var db *sql.DB

const (
	dbhost = "DBHOST"
	dbport = "DBPORT"
	dbuser = "DBUSER"
	dbpass = "DBPASS"
	dbname = "DBNAME"
)

func initDb() {
	config := dbConfig()
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config[dbhost], config[dbport],
		config[dbuser], config[dbpass], config[dbname])

	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected!")
}

func dbConfig() map[string]string {
	conf := make(map[string]string)
	host, ok := os.LookupEnv(dbhost)
	if !ok {
		panic("DBHOST environment variable required but not set")
	}
	port, ok := os.LookupEnv(dbport)
	if !ok {
		panic("DBPORT environment variable required but not set")
	}
	user, ok := os.LookupEnv(dbuser)
	if !ok {
		panic("DBUSER environment variable required but not set")
	}
	password, ok := os.LookupEnv(dbpass)
	if !ok {
		panic("DBPASS environment variable required but not set")
	}
	name, ok := os.LookupEnv(dbname)
	if !ok {
		panic("DBNAME environment variable required but not set")
	}
	conf[dbhost] = host
	conf[dbport] = port
	conf[dbuser] = user
	conf[dbpass] = password
	conf[dbname] = name
	return conf
}

// week represents the summary of a week of activity
// on a repository
type week struct {
	ID           int
	RepositoryID int
	WeekNumber   int
	Year         int
	CreatedOn    string
	IssuesClosed int
	IssuesOpened int
	Stars        int
	Commits      int
	WeekStart    string
	WeekEnd      string
	PrOpened     int
	PrMerged     int
	PrClosed     int
}

// timeline represents important events happened on a
// repository, which will be displayed on the repo timeline
type timeline struct {
	ID           int
	RepositoryID int
	Title        string
	Description  string
	Emoji        string
	Date         string
}

// repository contains the details of a repository
type repository struct {
	ID              int
	Name            string
	Owner           string
	RepoAge         int
	Initialized     bool
	CommitsPerMonth string
	StarsPerMonth   string
	TotalStars      int
}

// owner contains the details of an owner or a repo
type owner struct {
	ID                  int
	Name                string
	Description         string
	Avatar              string
	GitHubID            string
	AddedBy             string
	Enabled             bool
	InstallationID      string
	RepositorySelection string
}

// repoData contains the aggregate repository data returned
// by the API call
type repoData struct {
	MonthlyData monthlyData
	WeeklyData  []week
	Years       map[int]bool
	Timeline    []timeline
	Repository  repository
	Owner       owner
}

// monthlyData contains the monthly activity of a repo
type monthlyData struct {
	CommitsPerMonth string
	StarsPerMonth   string
}

// Error handling types

type errRepoNotInitialized string

func (e errRepoNotInitialized) Error() string {
	return string(e)
}

type errRepoNotFound string

func (e errRepoNotFound) Error() string {
	return string(e)
}

// parseParams accepts a req and returns the `num` path tokens found after the `prefix`.
// returns an error if the number of tokens are less or more than expected
func parseParams(req *http.Request, prefix string, num int) ([]string, error) {
	url := strings.TrimPrefix(req.URL.Path, prefix)
	params := strings.Split(url, "/")
	if len(params) != num || len(params[0]) == 0 || len(params[1]) == 0 {
		return nil, fmt.Errorf("Bad format. Expecting exactly %d params", num)
	}
	return params, nil
}

// repoHandler processes the response by parsing the params, then calling
// `query()`, and marshaling the result in JSON format, sending it to
// `http.ResponseWriter`.
func repoHandler(w http.ResponseWriter, req *http.Request) {
	repo := repository{}
	params, err := parseParams(req, "/api/repo/", 2)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	repo.Owner = params[0]
	repo.Name = params[1]

	data, err := queryRepo(&repo)
	if err != nil {
		switch err.(type) {
		case errRepoNotFound:
			http.Error(w, err.Error(), 404)
		case errRepoNotInitialized:
			http.Error(w, err.Error(), 401)
		default:
			http.Error(w, err.Error(), 500)
		}
		return
	}

	out, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Fprintf(w, string(out))
}

// queryRepo first fetches the repository, and if nothing is wrong
// it returns the result of fetchData()
func queryRepo(repo *repository) (*repoData, error) {
	err := fetchRepo(repo)
	if err != nil {
		return nil, err
	}

	return fetchData(repo)
}

// fetchData calls utility functions to collect data from
// the database, builds and returns the `RepoData` value
func fetchData(repo *repository) (*repoData, error) {
	data := repoData{}
	err := fetchMonthlyData(repo, &data)
	if err != nil {
		return nil, err
	}
	err = fetchWeeklyData(repo, &data)
	if err != nil {
		return nil, err
	}
	err = fetchYearlyData(repo, &data)
	if err != nil {
		return nil, err
	}
	err = fetchTimelineData(repo, &data)
	if err != nil {
		return nil, err
	}
	err = fetchOwnerData(repo, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// fetchRepo given a Repository value with name and owner of the repo
// fetches more details from the database and fills the value with more
// data
func fetchRepo(repo *repository) error {
	if len(repo.Name) == 0 {
		return fmt.Errorf("Repository name not correctly set")
	}
	if len(repo.Owner) == 0 {
		return fmt.Errorf("Repository owner not correctly set")
	}
	sqlStatement := `
        SELECT
            id,
            initialized,
            repository_created_months_ago
        FROM repositories
        WHERE repository_owner=$1 and repository_name=$2
        LIMIT 1;`
	row := db.QueryRow(sqlStatement, repo.Owner, repo.Name)
	err := row.Scan(&repo.ID, &repo.Initialized, &repo.RepoAge)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			//locally handle SQL error, abstract for caller
			return errRepoNotFound("Repository not found")
		default:
			return err
		}
	}
	if !repo.Initialized {
		return errRepoNotInitialized("Repository not initialized")
	}
	if repo.RepoAge < 3 {
		return errRepoNotInitialized("Repository not initialized")
	}
	return nil
}

// fetchOwnerData given a Repository object with the `Owner` value
// it fetches information about it from the database
func fetchOwnerData(repo *repository, data *repoData) error {
	if len(repo.Owner) == 0 {
		return fmt.Errorf("Repository owner not correctly set")
	}
	sqlStatement := `
        SELECT
            id,
            name,
            COALESCE(description, ''),
            COALESCE(avatar_url, ''),
            COALESCE(github_id, ''),
            added_by,
            enabled,
            COALESCE(installation_id, ''),
            repository_selection
        FROM organizations
        WHERE name=$1
        ORDER BY id DESC LIMIT 1;`
	row := db.QueryRow(sqlStatement, repo.Owner)
	err := row.Scan(&data.Owner.ID,
		&data.Owner.Name,
		&data.Owner.Description,
		&data.Owner.Avatar,
		&data.Owner.GitHubID,
		&data.Owner.AddedBy,
		&data.Owner.Enabled,
		&data.Owner.InstallationID,
		&data.Owner.RepositorySelection)
	if err != nil {
		return err
	}
	return nil
}

// fetchMonthlyData given a repository ID, it fetches the monthly
// data information
func fetchMonthlyData(repo *repository, data *repoData) error {
	if repo.ID == 0 {
		return fmt.Errorf("Repository ID not correctly set")
	}
	data.MonthlyData = monthlyData{}
	sqlStatement := `
        SELECT
            commits_per_month,
            stars_per_month
        FROM repositories_historic_data
        WHERE repository_id=$1
        ORDER BY id DESC LIMIT 1;`
	row := db.QueryRow(sqlStatement, repo.ID)
	err := row.Scan(
		&data.MonthlyData.CommitsPerMonth,
		&data.MonthlyData.StarsPerMonth)
	if err != nil {
		return err
	}

	return nil
}

// fetchWeeklyData given a repository ID, it fetches the weekly
// data information
func fetchWeeklyData(repo *repository, data *repoData) error {
	if repo.ID == 0 {
		return fmt.Errorf("Repository ID not correctly set")
	}
	rows, err := db.Query(`
        SELECT
            id,
            repository_id,
            week_number,
            year,
            created_on,
            issues_closed,
            issues_opened,
            stars,
            commits,
            week_start,
            week_end,
            pr_opened,
            pr_merged,
            pr_closed
        FROM repositories_weekly_data
        WHERE repository_id=$1
        ORDER BY id ASC`, repo.ID)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		week := week{}
		err = rows.Scan(
			&week.ID,
			&week.RepositoryID,
			&week.WeekNumber,
			&week.Year,
			&week.CreatedOn,
			&week.IssuesClosed,
			&week.IssuesOpened,
			&week.Stars,
			&week.Commits,
			&week.WeekStart,
			&week.WeekEnd,
			&week.PrOpened,
			&week.PrMerged,
			&week.PrClosed)
		if err != nil {
			return err
		}
		data.WeeklyData = append(data.WeeklyData, week)
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
}

// fetchYearlyData returns the list of years for which we have weekly data
// available
func fetchYearlyData(repo *repository, data *repoData) error {
	if data.WeeklyData == nil {
		return fmt.Errorf("Repository weekly data not correctly set")
	}
	data.Years = make(map[int]bool)
	for i := 0; i < len(data.WeeklyData); i++ {
		year := data.WeeklyData[i].Year
		data.Years[year] = true
	}
	return nil
}

// fetchTimelineData returns all the timeline data we have in the db about
// the repo
func fetchTimelineData(repo *repository, data *repoData) error {
	if repo.ID == 0 {
		return fmt.Errorf("Repository ID not correctly set")
	}
	rows, err := db.Query(`
        SELECT
            id,
            repository_id,
            title,
            description,
            emoji,
            date
        FROM repositories_timelines
        WHERE repository_id=$1
        ORDER BY date ASC`, repo.ID)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		timeline := timeline{}
		err = rows.Scan(
			&timeline.ID,
			&timeline.RepositoryID,
			&timeline.Title,
			&timeline.Description,
			&timeline.Emoji,
			&timeline.Date)
		if err != nil {
			return err
		}
		data.Timeline = append(data.Timeline, timeline)
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
}
