package teamcity

import (
	"fmt"
	"strconv"
	"strings"
)

type VcsLabelingOptions struct {
	BranchFilter    []string
	LabelingPattern string
	SuccessfulOnly  bool
	VcsRootId       string
}

//Properties returns a *Properties collection with properties filled related to this commit publisher parameters to be used in build features
func (s VcsLabelingOptions) Properties() *Properties {
	props := NewPropertiesEmpty()

	props.AddOrReplaceValue("branchFilter", strings.Join(s.BranchFilter, "\n"))
	props.AddOrReplaceValue("labelingPattern", s.LabelingPattern)
	props.AddOrReplaceValue("successfulOnly", strconv.FormatBool(s.SuccessfulOnly))
	props.AddOrReplaceValue("vcsRootId", s.VcsRootId)

	return props
}

func NewVcsLabelingOptions(filter []string, label, vcs string, success bool) VcsLabelingOptions {
	return VcsLabelingOptions{
		BranchFilter:    filter,
		LabelingPattern: label,
		SuccessfulOnly:  success,
		VcsRootId:       vcs,
	}
}

func NewFeatureVcsLabeling(opt VcsLabelingOptions) (*FeatureVcsLabeling, error) {
	if opt.LabelingPattern == "" {
		return nil, fmt.Errorf("pattern is required")
	}

	out := &FeatureVcsLabeling{
		Options:    opt,
		properties: opt.Properties(),
	}

	return out, nil
}

//VcsLabelingOptionsFromProperties grabs a Properties collection and transforms back to a VcsLabelingOptions
func VcsLabelingOptionsFromProperties(p *Properties) (*VcsLabelingOptions, error) {
	var out VcsLabelingOptions
	if label, ok := p.GetOk("labelingPattern"); ok {
		out.LabelingPattern = label
	} else {
		return nil, fmt.Errorf("Properties do not have 'labelingPattern' key")
	}

	if vcs, ok := p.GetOk("vcsRootId"); ok {
		out.VcsRootId = vcs
	} else {
		return nil, fmt.Errorf("Properties do not have 'vcsRootId' key")
	}

	//if branch, ok := p.GetOk("branchFilter"); ok {
	//	out.BranchFilter = branch
	//}

	if v, ok := p.GetOk("successfulOnly"); ok {
		if v == "true" {
			out.SuccessfulOnly = true
		} else {
			out.SuccessfulOnly = false
		}
	}

	return &out, nil
}
