package teamcity

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

//StepDocker represents a build step of type "Docker"
type StepDocker struct {
	ID       string
	Name     string
	stepType string
	stepJSON *stepJSON
	// docker.command.type
	CommandType string
	//dockerfile.source
	CommandSource string
	// command.args
	Args string
	// docker.push.remove.image
	PushRemoveImage *bool
	// dockerfile.content
	// dockerfile.path
	// docker.sub.command
	Content string
	// docker.image.namesAndTags
	Tag string
	// teamcity.build.workingDir
	WorkingDir string
	// teamcity.step.mode
	ExecuteMode StepExecuteMode
}

func NewStepDocker(name string, commandType string) (*StepDocker, error) {
	if commandType == "" {
		return nil, errors.New("need to set command_type parameter")
	}
	return &StepDocker{
		Name:        name,
		stepType:    StepTypeDocker,
		CommandType: commandType,
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

//Type returns the step type, in this case "StepTypeDocker".
func (s *StepDocker) Type() BuildStepType {
	return StepTypeDocker
}

func (s *StepDocker) properties() *Properties {
	props := NewPropertiesEmpty()
	props.AddOrReplaceValue("teamcity.step.mode", string(s.ExecuteMode))
	props.AddOrReplaceValue("command.args", s.Args)
	props.AddOrReplaceValue("docker.command.type", s.CommandType)
	if s.CommandType == "push" {
		props.AddOrReplaceValue("docker.push.remove.image", strconv.FormatBool(*s.PushRemoveImage))
		props.AddOrReplaceValue("dockerfile.source", "PATH")
		props.AddOrReplaceValue("docker.image.namesAndTags", s.Tag)
	}
	if s.CommandType == "build" {
		props.AddOrReplaceValue("dockerfile.source", s.CommandSource)
		props.AddOrReplaceValue("docker.image.namesAndTags", s.Tag)
		if s.CommandSource == "CONTENT" {
			props.AddOrReplaceValue("dockerfile.content", s.Content)
		} else if s.CommandSource == "PATH" {
			props.AddOrReplaceValue("dockerfile.path", s.Content)
		}
	}
	if s.CommandType == "other" {
		props.AddOrReplaceValue("dockerfile.source", "PATH")
		props.AddOrReplaceValue("docker.sub.command", s.Content)
		props.AddOrReplaceValue("teamcity.build.workingDir", s.WorkingDir)
	}
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

//MarshalJSON implements JSON serialization for StepDocker
func (s *StepDocker) MarshalJSON() ([]byte, error) {
	out := s.serializable()
	return json.Marshal(out)
}

//UnmarshalJSON implements JSON deserialization for StepDocker
func (s *StepDocker) UnmarshalJSON(data []byte) error {
	var aux stepJSON
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.Type != string(StepTypeDocker) {
		return fmt.Errorf("invalid type %s trying to deserialize into StepDocker entity", aux.Type)
	}
	s.Name = aux.Name
	s.ID = aux.ID
	s.stepType = StepTypeDocker

	props := aux.Properties
	if v, ok := props.GetOk("docker.command.type"); ok {
		s.CommandType = v
	}
	if v, ok := props.GetOk("docker.sub.command"); ok {
		if s.CommandType == "other" {
			s.Content = v
		}
	}
	if v, ok := props.GetOk("teamcity.build.workingDir"); ok {
		s.WorkingDir = v
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
			s.PushRemoveImage = NewBool(true)
		} else {
			s.PushRemoveImage = nil
		}
	}
	if v, ok := props.GetOk("docker.image.namesAndTags"); ok {
		s.Tag = v
	}
	if v, ok := props.GetOk("teamcity.step.mode"); ok {
		s.ExecuteMode = v
	}
	return nil
}
