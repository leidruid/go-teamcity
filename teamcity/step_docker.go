package teamcity

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

//StepCommandLine represents a a build step of type "CommandLine"
type StepDocker struct {
	ID            string
	Name          string
	stepType      string
	stepJSON      *stepJSON
	CommandSource string
	// command.args
	Args string
	// docker.command.type
	CommandType string
	// docker.push.remove.image
	PushRemoveImage bool
	// dockerfile.content
	// dockerfile.path
	Content string
	// docker.image.namesAndTags
	Tag string
	// teamcity.step.mode
	ExecuteMode StepExecuteMode
	//docker.sub.command
}

func NewStepDocker(name string, fromSource string, dockerCommandType string, dockerContent string, dockerArgs string, dockerTag string) (*StepDocker, error) {
	if dockerContent == "" {
		return nil, errors.New("tasks is required")
	}
	return &StepDocker{
		Name:          name,
		stepType:      StepTypeDocker,
		CommandType:   dockerCommandType,
		Content:       dockerContent,
		Args:          dockerArgs,
		Tag:           dockerTag,
		CommandSource: fromSource,
		ExecuteMode:   StepExecuteModeDefault,
	}, nil
}

//GetID is a wrapper implementation for ID field, to comply with Step interface
func (s *StepDocker) GetID() string {
	return s.ID
}

//GetName is a wrapper implementation for Name field, to comply with Step interface
func (s *StepDocker) GetName() string {
	return s.Name
}

//Type returns the step type, in this case "StepTypeCommandLine".
func (s *StepDocker) Type() BuildStepType {
	return StepTypeGradle
}

func (s *StepDocker) properties() *Properties {
	props := NewPropertiesEmpty()
	props.AddOrReplaceValue("teamcity.step.mode", string(s.ExecuteMode))
	props.AddOrReplaceValue("command.args", s.Args)
	props.AddOrReplaceValue("docker.command.type", s.CommandType)
	if s.CommandType == "push" {
		props.AddOrReplaceValue("docker.push.remove.image", strconv.FormatBool(s.PushRemoveImage))
	}
	if s.CommandType == "build" {
		if s.CommandSource == "CONTENT" {
			props.AddOrReplaceValue("dockerfile.content", s.Content)
		} else {
			props.AddOrReplaceValue("dockerfile.path", s.Content)
		}
	}
	props.AddOrReplaceValue("docker.image.namesAndTags", s.Tag)
	return props
}

func (s *StepDocker) serializable() *stepJSON {
	return &stepJSON{
		ID:         s.ID,
		Name:       s.Name,
		Type:       s.stepType,
		Properties: s.properties(),
	}
}

//MarshalJSON implements JSON serialization for StepCommandLine
func (s *StepDocker) MarshalJSON() ([]byte, error) {
	out := s.serializable()
	return json.Marshal(out)
}

//UnmarshalJSON implements JSON deserialization for StepCommandLine
func (s *StepDocker) UnmarshalJSON(data []byte) error {
	var aux stepJSON
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.Type != string(StepTypeGradle) {
		return fmt.Errorf("invalid type %s trying to deserialize into StepGradle entity", aux.Type)
	}
	s.Name = aux.Name
	s.ID = aux.ID
	s.stepType = StepTypeGradle

	props := aux.Properties
	if v, ok := props.GetOk("docker.command.type"); ok {
		s.CommandType = v
	}
	if v, ok := props.GetOk("dockerfile.source"); ok {
		s.CommandSource = v
		if v == "PATH" {
			if c, cok := props.GetOk("dockerfile.path"); cok {
				s.Content = c
			}
		} else if v == "CONTENT" {
			if c, cok := props.GetOk("dockerfile.content"); cok {
				s.Content = c
			}
		}
	}
	if v, ok := props.GetOk("command.args"); ok {
		s.Args = v
	}
	if v, ok := props.GetOk("docker.push.remove.image"); ok {
		if v == "true" {
			s.PushRemoveImage = true
		} else {
			s.PushRemoveImage = false
		}
	}
	if v, ok := props.GetOk("teamcity.step.mode"); ok {
		s.ExecuteMode = StepExecuteMode(v)
	}
	return nil
}
