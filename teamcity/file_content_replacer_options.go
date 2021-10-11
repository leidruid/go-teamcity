package teamcity

import (
	"fmt"
	"strconv"
	"strings"
)

type FileContentReplacerOptions struct {
	TeamcityFileContentReplacerFailBuild            bool
	TeamcityFileContentReplacerFileEncoding         string
	TeamcityFileContentReplacerFileEncodingCustom   string
	TeamcityFileContentReplacerPattern              string
	TeamcityFileContentReplacerPatternCaseSensitive bool
	TeamcityFileContentReplacerRegexMode            string
	TeamcityFileContentReplacerReplacement          string
	TeamcityFileContentReplacerWildcards            []string
}

//Properties returns a *Properties collection with properties filled related to this commit publisher parameters to be used in build features
func (s FileContentReplacerOptions) Properties() *Properties {
	props := NewPropertiesEmpty()

	props.AddOrReplaceValue("teamcity.file.content.replacer.failBuild", strconv.FormatBool(s.TeamcityFileContentReplacerFailBuild))

	props.AddOrReplaceValue("teamcity.file.content.replacer.file.encoding", s.TeamcityFileContentReplacerFileEncoding)

	if s.TeamcityFileContentReplacerFileEncoding != "custom" {
		props.AddOrReplaceValue("teamcity.file.content.replacer.file.encoding.custom", s.TeamcityFileContentReplacerFileEncoding)
	} else {
		props.AddOrReplaceValue("teamcity.file.content.replacer.file.encoding.custom", s.TeamcityFileContentReplacerFileEncodingCustom)
	}

	props.AddOrReplaceValue("teamcity.file.content.replacer.pattern", s.TeamcityFileContentReplacerPattern)
	props.AddOrReplaceValue("teamcity.file.content.replacer.pattern.case.sensitive", strconv.FormatBool(s.TeamcityFileContentReplacerPatternCaseSensitive))
	props.AddOrReplaceValue("teamcity.file.content.replacer.regexMode", s.TeamcityFileContentReplacerRegexMode)
	props.AddOrReplaceValue("teamcity.file.content.replacer.replacement", s.TeamcityFileContentReplacerReplacement)
	props.AddOrReplaceValue("teamcity.file.content.replacer.wildcards", strings.Join(s.TeamcityFileContentReplacerWildcards, "\n"))

	return props
}

func NewFileContentReplacerOptions(encoding, encodingCustom, pattern, regexMode, replacement string, wildcards []string, failBuild, caseSensitive bool) FileContentReplacerOptions {
	return FileContentReplacerOptions{
		TeamcityFileContentReplacerFailBuild:            failBuild,
		TeamcityFileContentReplacerFileEncoding:         encoding,
		TeamcityFileContentReplacerFileEncodingCustom:   encodingCustom,
		TeamcityFileContentReplacerPattern:              pattern,
		TeamcityFileContentReplacerPatternCaseSensitive: caseSensitive,
		TeamcityFileContentReplacerRegexMode:            regexMode,
		TeamcityFileContentReplacerReplacement:          replacement,
		TeamcityFileContentReplacerWildcards:            wildcards,
	}
}

func NewFeatureFileContentReplacer(opt FileContentReplacerOptions) (*FeatureFileContentReplacer, error) {
	if len(opt.TeamcityFileContentReplacerWildcards) == 0 {
		return nil, fmt.Errorf("teamcity.file.content.replacer.wildcards is required")
	}

	if opt.TeamcityFileContentReplacerPattern == "" {
		return nil, fmt.Errorf("teamcity.file.content.replacer.pattern is required")
	}

	out := &FeatureFileContentReplacer{
		Options:    opt,
		properties: opt.Properties(),
	}

	return out, nil
}

//FileContentReplacerOptionsFromProperties grabs a Properties collection and transforms back to a FileContentReplacerOptions
func FileContentReplacerOptionsFromProperties(p *Properties) (*FileContentReplacerOptions, error) {
	var out FileContentReplacerOptions
	if v, ok := p.GetOk("teamcity.file.content.replacer.failBuild"); ok {
		if v == "true" {
			out.TeamcityFileContentReplacerFailBuild = true
		} else {
			out.TeamcityFileContentReplacerFailBuild = false
		}
	} else {
		return nil, fmt.Errorf("Properties do not have 'teamcity.file.content.replacer.failBuild' key")
	}

	if v, ok := p.GetOk("teamcity.file.content.replacer.file.encoding"); ok {
		out.TeamcityFileContentReplacerFileEncoding = v
	} else {
		return nil, fmt.Errorf("Properties do not have 'teamcity.file.content.replacer.file.encoding' key")
	}

	if v, ok := p.GetOk("teamcity.file.content.replacer.file.encoding.custom"); ok {
		out.TeamcityFileContentReplacerFileEncodingCustom = v
	} else {
		return nil, fmt.Errorf("Properties do not have 'teamcity.file.content.replacer.file.encoding.custom' key")
	}

	if v, ok := p.GetOk("teamcity.file.content.replacer.pattern"); ok {
		out.TeamcityFileContentReplacerPattern = v
	} else {
		return nil, fmt.Errorf("Properties do not have 'teamcity.file.content.replacer.pattern' key")
	}

	if v, ok := p.GetOk("teamcity.file.content.replacer.pattern.case.sensitive"); ok {
		if v == "true" {
			out.TeamcityFileContentReplacerPatternCaseSensitive = true
		} else {
			out.TeamcityFileContentReplacerPatternCaseSensitive = false
		}
	} else {
		return nil, fmt.Errorf("Properties do not have 'teamcity.file.content.replacer.pattern.case.sensitive' key")
	}

	if v, ok := p.GetOk("teamcity.file.content.replacer.regexMode"); ok {
		out.TeamcityFileContentReplacerRegexMode = v
	} else {
		return nil, fmt.Errorf("Properties do not have 'teamcity.file.content.replacer.regexMode' key")
	}

	if v, ok := p.GetOk("teamcity.file.content.replacer.replacement"); ok {
		out.TeamcityFileContentReplacerReplacement = v
	} else {
		return nil, fmt.Errorf("Properties do not have 'teamcity.file.content.replacer.replacement' key")
	}

	if v, ok := p.GetOk("teamcity.file.content.replacer.wildcards"); ok {
		out.TeamcityFileContentReplacerWildcards = strings.Split(v, "\n")
	} else {
		return nil, fmt.Errorf("Properties do not have 'teamcity.file.content.replacer.wildcards' key")
	}

	return &out, nil
}
