package configs

import (
	"github.com/pkg/errors"
	"github.com/statping-ng/statping-ng/utils"
	"net/http"
	"strconv"
)

func LoadConfigForm(r *http.Request) (*DbConfig, error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}
	g := r.PostForm.Get
	dbHost := g("db_host")
	dbUser := g("db_user")
	dbPass := g("db_password")
	dbDatabase := g("db_database")
	dbConn := g("db_connection")
	dbPort := utils.ToInt(g("db_port"))
	project := g("project")
	username := g("username")
	password := g("password")
	description := g("description")
	domain := g("domain")
	email := g("email")
	language := g("language")
	reports, _ := strconv.ParseBool(g("send_reports"))

	// New fields for Intermediate Statuses
	statusMinorOutageName := g("status_minor_outage_name")
	statusMinorOutageColor := g("status_minor_outage_color")
	statusMajorOutageName := g("status_major_outage_name")
	statusMajorOutageColor := g("status_major_outage_color")
	enableIntermediateStatuses, _ := strconv.ParseBool(g("enable_intermediate_statuses"))

	if project == "" || username == "" || password == "" {
		err := errors.New("Missing required elements on setup form")
		return nil, err
	}

	p := utils.Params
	p.Set("DB_CONN", dbConn)
	p.Set("DB_HOST", dbHost)
	p.Set("DB_USER", dbUser)
	p.Set("DB_PORT", dbPort)
	p.Set("DB_PASS", dbPass)
	p.Set("DB_DATABASE", dbDatabase)
	p.Set("NAME", project)
	p.Set("DESCRIPTION", description)
	p.Set("LANGUAGE", language)
	p.Set("ALLOW_REPORTS", reports)
	p.Set("ADMIN_USER", username)
	p.Set("ADMIN_PASSWORD", password)
	p.Set("ADMIN_EMAIL", email)
	p.Set("ENABLE_INTERMEDIATE_STATUSES", enableIntermediateStatuses)
	p.Set("STATUS_MINOR_OUTAGE_NAME", statusMinorOutageName)
	p.Set("STATUS_MINOR_OUTAGE_COLOR", statusMinorOutageColor)
	p.Set("STATUS_MAJOR_OUTAGE_NAME", statusMajorOutageName)
	p.Set("STATUS_MAJOR_OUTAGE_COLOR", statusMajorOutageColor)

	confg := &DbConfig{
		DbConn:       dbConn,
		DbHost:       dbHost,
		DbUser:       dbUser,
		DbPass:       dbPass,
		DbData:       dbDatabase,
		DbPort:       int(dbPort),
		Project:      project,
		Description:  description,
		Domain:       domain,
		Username:     username,
		Password:     password,
		Email:        email,
		Location:     utils.Directory,
		Language:     language,
		AllowReports: reports,

		// Add the new config fields for status
		StatusMinorOutageName:   	statusMinorOutageName,
		StatusMinorOutageColor:  	statusMinorOutageColor,
		StatusMajorOutageName:   	statusMajorOutageName,
		StatusMajorOutageColor:  	statusMajorOutageColor,
		EnableIntermediateStatuses: enableIntermediateStatuses,
	}

	return confg, nil
}
