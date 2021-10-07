package teamcity

import "fmt"

// StatusPublisherBitbucketServerOptions represents parameters used to create Github Commit Status Publisher Feature
type StatusPublisherBitbucketServerOptions struct {
	//Host is the Bitbucket Server URL
	Host string
	//Username
	Username string
	//Password
	Password string
}

//NewCommitStatusPublisherBitbucketServerOptionsPassword returns options created for AuthenticationType = 'password'. No validation is performed, parameters indicate mandatory fields.
func NewCommitStatusPublisherBitbucketServerOptionsPassword(host string, username string, password string) StatusPublisherBitbucketServerOptions {
	return StatusPublisherBitbucketServerOptions{
		Host:     host,
		Username: username,
		Password: password,
	}
}

//NewFeatureCommitStatusPublisherBitbucketServer creates a Build Feature Commit status Publisher to Bitbucket Server with the given options and validates the required properties.
//VcsRootID is optional - if empty, it will apply the commit publisher feature to all VCS roots.
func NewFeatureCommitStatusPublisherBitbucketServer(opt StatusPublisherBitbucketServerOptions, vcsRootID string) (*FeatureCommitStatusPublisher, error) {
	if opt.Host == "" {
		return nil, fmt.Errorf("Host is required")
	}

	if opt.Username == "" || opt.Password == "" {
		return nil, fmt.Errorf("username/password required for auth type 'password'")
	}

	out := &FeatureCommitStatusPublisher{
		Options:    opt,
		properties: opt.Properties(),
	}

	if vcsRootID != "" {
		out.vcsRootID = vcsRootID
	}

	return out, nil
}

//Properties returns a *Properties collection with properties filled related to this commit publisher parameters to be used in build features
func (s StatusPublisherBitbucketServerOptions) Properties() *Properties {
	props := NewPropertiesEmpty()

	props.AddOrReplaceValue("publisherId", "atlassianStashPublisher")
	props.AddOrReplaceValue("stashBaseUrl", s.Host)
	props.AddOrReplaceValue("stashUsername", s.Username)
	props.AddOrReplaceValue("secure:stashPassword", s.Password)

	return props
}

//CommitStatusPublisherBitbucketServerOptionsFromProperties grabs a Properties collection and transforms back to a StatusPublisherBitbucketServerOptions
func CommitStatusPublisherBitbucketServerOptionsFromProperties(p *Properties) (*StatusPublisherBitbucketServerOptions, error) {
	var out StatusPublisherBitbucketServerOptions

	if host, ok := p.GetOk("stashBaseUrl"); ok {
		out.Host = host
	} else {
		return nil, fmt.Errorf("properties do not have 'stashBaseUrl' key")
	}

	if user, ok := p.GetOk("stashUsername"); ok {
		out.Username = user
	} else {
		return nil, fmt.Errorf("properties do not have 'stashUsername' key")
	}

	return &out, nil
}
