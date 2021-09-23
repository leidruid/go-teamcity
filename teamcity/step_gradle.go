package teamcity

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

//StepCommandLine represents a a build step of type "CommandLine"
type StepGradle struct {
	ID       string
	Name     string
	stepType string
	stepJSON *stepJSON
	//CoverageEmmaIncludeSource	bool
	//CoverageEmmaInstrParamters string
	//CoverageIdeaIncludePatterns string
	//ToolJacoco	string
	GradleCmdParams  string
	GradleBuildFile  string
	GradleTasksNames string
	GradleWrapperUse bool
	ExecuteMode      StepExecuteMode
}

func NewStepGradle(name string, tasks string, gradleParams string, gradleBuildFile string) (*StepGradle, error) {
	if tasks == "" {
		return nil, errors.New("tasks is required")
	}
	return &StepGradle{
		Name:             name,
		stepType:         StepTypeGradle,
		GradleCmdParams:  gradleParams,
		GradleBuildFile:  gradleBuildFile,
		GradleTasksNames: tasks,
		GradleWrapperUse: true,
		ExecuteMode:      StepExecuteModeDefault,
	}, nil
}

//GetID is a wrapper implementation for ID field, to comply with Step interface
func (s *StepGradle) GetID() string {
	return s.ID
}

//GetName is a wrapper implementation for Name field, to comply with Step interface
func (s *StepGradle) GetName() string {
	return s.Name
}

//Type returns the step type, in this case "StepTypeCommandLine".
func (s *StepGradle) Type() BuildStepType {
	return StepTypeGradle
}

func (s *StepGradle) properties() *Properties {
	props := NewPropertiesEmpty()
	props.AddOrReplaceValue("teamcity.step.mode", string(s.ExecuteMode))
	props.AddOrReplaceValue("ui.gradleRunner.additional.gradle.cmd.params", s.GradleCmdParams)
	props.AddOrReplaceValue("ui.gradleRunner.gradle.tasks.names", s.GradleTasksNames)
	props.AddOrReplaceValue("ui.gradleRunner.gradle.build.file", s.GradleBuildFile)
	if s.GradleWrapperUse {
		props.AddOrReplaceValue("ui.gradleRunner.gradle.wrapper.useWrapper", strconv.FormatBool(s.GradleWrapperUse))
	}
	return props
}

func (s *StepGradle) serializable() *stepJSON {
	return &stepJSON{
		ID:         s.ID,
		Name:       s.Name,
		Type:       s.stepType,
		Properties: s.properties(),
	}
}

//MarshalJSON implements JSON serialization for StepCommandLine
func (s *StepGradle) MarshalJSON() ([]byte, error) {
	out := s.serializable()
	return json.Marshal(out)
}

//UnmarshalJSON implements JSON deserialization for StepCommandLine
func (s *StepGradle) UnmarshalJSON(data []byte) error {
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
	if v, ok := props.GetOk("ui.gradleRunner.gradle.wrapper.useWrapper"); ok {
		if v == "false" {
			s.GradleWrapperUse = false
		} else {
			s.GradleWrapperUse = true
		}
	}
	if v, ok := props.GetOk("ui.gradleRunner.additional.gradle.cmd.params"); ok {
		s.GradleCmdParams = v
	}
	if v, ok := props.GetOk("ui.gradleRunner.gradle.tasks.names"); ok {
		s.GradleTasksNames = v
	}
	if v, ok := props.GetOk("ui.gradleRunner.gradle.build.file"); ok {
		s.GradleBuildFile = v
	}
	if v, ok := props.GetOk("teamcity.step.mode"); ok {
		s.ExecuteMode = StepExecuteMode(v)
	}
	return nil
}
