package teamcity

import (
	"encoding/json"
	"fmt"
)

// BuildStepType represents most common step types for build steps
type BuildStepType = string

const (
	//StepTypePowershell step type
	StepTypePowershell BuildStepType = "jetbrains_powershell"
	//StepTypeDotnetCli step type
	StepTypeDotnetCli BuildStepType = "dotnet.cli"
	//StepTypeCommandLine (shell/cmd) step type
	StepTypeCommandLine          BuildStepType = "simpleRunner"
	StepTypeGradle               BuildStepType = "gradle-runner"
	StepTypeDocker               BuildStepType = "DockerCommand"
	StepTypeOctopusPushPackage   BuildStepType = "octopus.push.package"
	StepTypeOctopusCreateRelease BuildStepType = "octopus.create.release"
)

//StepExecuteMode represents how a build configuration step will execute regarding others.
type StepExecuteMode = string

// Step interface represents a build configuration/template build step. To interact with concrete step types, see the Step* types.
type Step interface {
	GetID() string
	GetName() string
	Type() string

	serializable() *stepJSON
}

type stepJSON struct {
	Disabled   *bool       `json:"disabled,omitempty" xml:"disabled"`
	Href       string      `json:"href,omitempty" xml:"href"`
	ID         string      `json:"id,omitempty" xml:"id"`
	Inherited  *bool       `json:"inherited,omitempty" xml:"inherited"`
	Name       string      `json:"name,omitempty" xml:"name"`
	Properties *Properties `json:"properties,omitempty"`
	Type       string      `json:"type,omitempty" xml:"type"`
}

type stepsJSON struct {
	Count int32       `json:"count,omitempty" xml:"count"`
	Items []*stepJSON `json:"step"`
}

var stepsReadingFunc = func(dt []byte, out interface{}) error {
	var payload stepsJSON
	if err := json.Unmarshal(dt, &payload); err != nil {
		return err
	}

	var steps = make([]Step, payload.Count)
	for i := 0; i < int(payload.Count); i++ {
		sdt, err := json.Marshal(payload.Items[i])
		if err != nil {
			return err
		}
		err = stepReadingFunc(sdt, &steps[i])
		if err != nil {
			return err
		}
	}
	replaceValue(out, &steps)
	return nil
}

var stepReadingFunc = func(dt []byte, out interface{}) error {
	var payload stepJSON
	if err := json.Unmarshal(dt, &payload); err != nil {
		return err
	}

	var step Step
	var err error
	switch payload.Type {
	case StepTypePowershell:
		var ps StepPowershell
		err = ps.UnmarshalJSON(dt)
		step = &ps
	case StepTypeCommandLine:
		var cmd StepCommandLine
		err = cmd.UnmarshalJSON(dt)
		step = &cmd
	case StepTypeOctopusPushPackage:
		var opp StepOctopusPushPackage
		err = opp.UnmarshalJSON(dt)
		step = &opp
	case StepTypeOctopusCreateRelease:
		var ocr StepOctopusCreateRelease
		err = ocr.UnmarshalJSON(dt)
		step = &ocr
	case StepTypeGradle:
		var grd StepGradle
		err = grd.UnmarshalJSON(dt)
		step = &grd
	case StepTypeDocker:
		var dck StepDocker
		err = dck.UnmarshalJSON(dt)
		step = &dck
	default:
		return fmt.Errorf("Unsupported step type: '%s' (id:'%s')", payload.Type, payload.ID)
	}
	if err != nil {
		return err
	}

	replaceValue(out, &step)
	return nil
}
