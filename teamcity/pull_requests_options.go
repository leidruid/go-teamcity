package teamcity

import (
	"fmt"
	"strings"
)

type PullRequestsOptions struct {
	AuthenticationType string
	FilterSourceBranch []string
	FilterTargetBranch []string
	ServerUrl          string
	Username           string
	Password           string
}

//Properties returns a *Properties collection with properties filled related to this commit publisher parameters to be used in build features
func (s PullRequestsOptions) Properties() *Properties {
	props := NewPropertiesEmpty()
	props.AddOrReplaceValue("providerType", "bitbucketServer")
	props.AddOrReplaceValue("authenticationType", s.AuthenticationType)
	props.AddOrReplaceValue("filterSourceBranch", strings.Join(s.FilterSourceBranch, "\n"))
	props.AddOrReplaceValue("filterTargetBranch", strings.Join(s.FilterTargetBranch, "\n"))
	props.AddOrReplaceValue("serverUrl", s.ServerUrl)

	if s.AuthenticationType == "password" {
		props.AddOrReplaceValue("username", s.Username)
		props.AddOrReplaceValue("secure:password", s.Password)
	}

	return props
}

func NewPullRequestsOptionsPassword(username, password, url string, source, target []string) PullRequestsOptions {
	return PullRequestsOptions{
		AuthenticationType: "password",
		FilterSourceBranch: source,
		FilterTargetBranch: target,
		ServerUrl:          url,
		Username:           username,
		Password:           password,
	}
}

func NewPullRequestsOptionsVcs(url string, source, target []string) PullRequestsOptions {
	return PullRequestsOptions{
		AuthenticationType: "vcsRoot",
		FilterSourceBranch: source,
		FilterTargetBranch: target,
		ServerUrl:          url,
	}
}

func NewFeaturePullRequests(opt PullRequestsOptions, vcsRootID string) (*FeaturePullRequests, error) {
	if opt.AuthenticationType == "" {
		return nil, fmt.Errorf("AuthenticationType is required")
	}
	if opt.AuthenticationType != "password" && opt.AuthenticationType != "vcsRoot" {
		return nil, fmt.Errorf("invalid AuthenticationType, must be 'password' or 'vcsRoot'")
	}
	if opt.ServerUrl == "" {
		return nil, fmt.Errorf("ServerUrl is required")
	}

	if opt.AuthenticationType == "password" {
		if opt.Username == "" || opt.Password == "" {
			return nil, fmt.Errorf("username/password required for auth type 'password'")
		}
	}

	out := &FeaturePullRequests{
		Options:    opt,
		properties: opt.Properties(),
	}

	if vcsRootID != "" {
		out.vcsRootID = vcsRootID
	}

	return out, nil
}

//PullRequestsOptionsFromProperties grabs a Properties collection and transforms back to a PullRequestsOptions
func PullRequestsOptionsFromProperties(p *Properties) (*PullRequestsOptions, error) {
	var out PullRequestsOptions
	if url, ok := p.GetOk("serverUrl"); ok {
		out.ServerUrl = url
	} else {
		return nil, fmt.Errorf("Properties do not have 'serverUrl' key")
	}

	if authType, ok := p.GetOk("authenticationType"); ok {
		out.AuthenticationType = authType
		switch authType {
		case "password":
			u, _ := p.GetOk("username")
			out.Username = u
			//Password or AccessToken is never returned from properties, because it is secure. Once set, we cannot read it back
		}
	} else {
		return nil, fmt.Errorf("Properties do not have 'authenticationType' key")
	}

	return &out, nil
}
