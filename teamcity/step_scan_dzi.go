package teamcity

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type StepScanDzi struct {
	ID                      string
	Name                    string
	stepType                string
	stepJSON                *stepJSON
	ArtifactsUrls           []string
	FailBuildByTagsStatuses bool
	ExecuteMode             StepExecuteMode
	ExecuteCondition        [][]string
}

//NewStepScanDzi creates a step scanning
func NewStepScanDzi(name string, urls []string) (*StepScanDzi, error) {

	return &StepScanDzi{
		Name:          name,
		stepType:      StepTypeScanDzi,
		ArtifactsUrls: urls,
	}, nil
}

//GetID is a wrapper implementation for ID field, to comply with Step interface
func (s *StepScanDzi) GetID() string {
	return s.ID
}

//GetName is a wrapper implementation for Name field, to comply with Step interface
func (s *StepScanDzi) GetName() string {
	return s.Name
}

//Type returns the step type, in this case "StepTypeCommandLine".
func (s *StepScanDzi) Type() BuildStepType {
	return StepTypeScanDzi
}

func (s *StepScanDzi) properties() *Properties {
	props := NewPropertiesEmpty()
	props.AddOrReplaceValue("teamcity.step.mode", s.ExecuteMode)
	props.AddOrReplaceValue("artifacts_urls", strings.Join(s.ArtifactsUrls, "\\n"))
	props.AddOrReplaceValue("fail_build_by_tags_statuses", strconv.FormatBool(s.FailBuildByTagsStatuses))

	ecs, _ := json.Marshal(s.ExecuteCondition)
	props.AddOrReplaceValue("teamcity.step.conditions", string(ecs))

	return props
}

func (s *StepScanDzi) serializable() *stepJSON {
	return &stepJSON{
		ID:         s.ID,
		Name:       s.Name,
		Type:       s.stepType,
		Properties: s.properties(),
	}
}

//MarshalJSON implements JSON serialization for StepScanDzi
func (s *StepScanDzi) MarshalJSON() ([]byte, error) {
	out := s.serializable()
	return json.Marshal(out)
}

//UnmarshalJSON implements JSON deserialization for StepScanDzi
func (s *StepScanDzi) UnmarshalJSON(data []byte) error {
	var aux stepJSON
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.Type != StepTypeScanDzi {
		return fmt.Errorf("invalid type %s trying to deserialize into StepTypeScanDzi entity", aux.Type)
	}
	s.Name = aux.Name
	s.ID = aux.ID
	s.stepType = StepTypeScanDzi

	props := aux.Properties
	if v, ok := props.GetOk("artifacts_urls"); ok {
		s.ArtifactsUrls = strings.Split(v, "\\n")
	}

	if v, ok := props.GetOk("fail_build_by_tags_statuses"); ok {
		if v == "False" {
			s.FailBuildByTagsStatuses = false
		} else {
			s.FailBuildByTagsStatuses = true
		}
	}

	if v, ok := props.GetOk("teamcity.step.mode"); ok {
		s.ExecuteMode = v
	}

	if v, ok := props.GetOk("teamcity.step.conditions"); ok {
		var ecJson [][]string
		_ = json.Unmarshal([]byte(v), &ecJson)
		s.ExecuteCondition = ecJson
	}

	return nil
}
