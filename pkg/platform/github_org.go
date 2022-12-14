package platform

import (
	"fmt"
	"strings"

	"github.com/gocarina/gocsv"
)

// GithubOrgMembers parses the CSV file generated by the Github Members page.
type GithubOrgMembers struct{}

func (p *GithubOrgMembers) Description() ProcessorDescription {
	return ProcessorDescription{
		Kind: "github-org",
		Name: "Github Organization Members",
		Steps: []string{
			"Open https://github.com/orgs/<org>/people",
			"Click Export",
			"Select 'CSV'",
			"Download resulting CSV file for analysis",
			"Execute 'acls-in-yaml --kind={{.Kind}} --input={{.Path}}'",
		},
	}
}

type githubMemberRecord struct {
	Login string `csv:"login"`
	Name  string `csv:"name"`
	Role  string `csv:"role"`
}

func (p *GithubOrgMembers) Process(c Config) (*Artifact, error) {
	src, err := NewSourceFromConfig(c, p)
	if err != nil {
		return nil, fmt.Errorf("source: %w", err)
	}
	a := &Artifact{Metadata: src}

	records := []githubMemberRecord{}
	if err := gocsv.UnmarshalBytes(src.content, &records); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	for _, r := range records {
		u := User{
			Account: r.Login,
			Name:    strings.TrimSpace(r.Name),
			Role:    r.Role,
		}

		if strings.HasSuffix(u.Name, "Bot") || strings.HasSuffix(u.Account, "Bot") {
			a.Bots = append(a.Bots, u)
			continue
		}

		a.Users = append(a.Users, u)
	}

	return a, nil
}
