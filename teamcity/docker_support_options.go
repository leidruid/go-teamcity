package teamcity

import (
	"fmt"
	"strconv"
)

type DockerSupportOptions struct {
	//Login2registry string of external id's docker connections
	Login2registry string
	LoginCheckbox  string
	//On server clean-up, delete pushed Docker images from registry
	CleanupPushed bool
}

//Properties returns a *Properties collection with properties filled related to this commit publisher parameters to be used in build features
func (s DockerSupportOptions) Properties() *Properties {
	props := NewPropertiesEmpty()

	props.AddOrReplaceValue("login2registry", s.Login2registry)
	props.AddOrReplaceValue("loginCheckbox", "on")
	props.AddOrReplaceValue("cleanupPushed", strconv.FormatBool(s.CleanupPushed))

	return props
}

func NewDockerSupportOptions(registry string, cleanup bool) DockerSupportOptions {
	return DockerSupportOptions{
		Login2registry: registry,
		CleanupPushed:  cleanup,
		LoginCheckbox:  "on",
	}
}

func NewFeatureDockerSupport(opt DockerSupportOptions) (*FeatureDockerSupport, error) {
	if opt.Login2registry == "" {
		return nil, fmt.Errorf("registry is required")
	}

	out := &FeatureDockerSupport{
		Options:    opt,
		properties: opt.Properties(),
	}

	return out, nil
}

//DockerSupportOptionsFromProperties grabs a Properties collection and transforms back to a DockerSupportOptions
func DockerSupportOptionsFromProperties(p *Properties) (*DockerSupportOptions, error) {
	var out DockerSupportOptions
	if registries, ok := p.GetOk("login2registry"); ok {
		out.Login2registry = registries
	} else {
		return nil, fmt.Errorf("Properties do not have 'login2registry' key")
	}

	if v, ok := p.GetOk("cleanupPushed"); ok {
		if v == "true" {
			out.CleanupPushed = true
		} else {
			out.CleanupPushed = false
		}
	} else {
		return nil, fmt.Errorf("Properties do not have 'cleanupPushed' key")
	}

	return &out, nil
}
