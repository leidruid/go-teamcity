package teamcity

import (
	"fmt"
)

type SshAgentOptions struct {
	TeamcitySshKey string
}

//Properties returns a *Properties collection with properties filled related to this commit publisher parameters to be used in build features
func (s SshAgentOptions) Properties() *Properties {
	props := NewPropertiesEmpty()
	props.AddOrReplaceValue("teamcitySshKey", s.TeamcitySshKey)

	return props
}

func NewSshAgentOptions(sshKey string) SshAgentOptions {
	return SshAgentOptions{
		TeamcitySshKey: sshKey,
	}
}

func NewFeatureSshAgent(opt SshAgentOptions) (*FeatureSshAgent, error) {
	if opt.TeamcitySshKey == "" {
		return nil, fmt.Errorf("teamcitySshKey is required")
	}

	out := &FeatureSshAgent{
		Options:    opt,
		properties: opt.Properties(),
	}

	return out, nil
}

//SshAgentOptionsFromProperties grabs a Properties collection and transforms back to a SshAgentOptions
func SshAgentOptionsFromProperties(p *Properties) (*SshAgentOptions, error) {
	var out SshAgentOptions
	if sshKey, ok := p.GetOk("teamcitySshKey"); ok {
		out.TeamcitySshKey = sshKey
	} else {
		return nil, fmt.Errorf("Properties do not have 'teamcitySshKey' key")
	}

	return &out, nil
}
