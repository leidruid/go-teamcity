package teamcity

import "encoding/json"

type FeatureDockerSupport struct {
	id          string
	disabled    bool
	buildTypeID string
	Options     FeatureDockerSupportOptions

	properties *Properties
}

//FeatureDockerSupportOptions represents options needed to create a docker support build feature
type FeatureDockerSupportOptions interface {
	Properties() *Properties
}

//ID returns the ID for this instance.
func (f *FeatureDockerSupport) ID() string {
	return f.id
}

//SetID sets the ID for this instance.
func (f *FeatureDockerSupport) SetID(value string) {
	f.id = value
}

//Type returns the "DockerSupport", the keyed-type for this build feature instance
func (f *FeatureDockerSupport) Type() string {
	return "DockerSupport"
}

//Disabled returns whether this build feature is disabled or not.
func (f *FeatureDockerSupport) Disabled() bool {
	return f.disabled
}

//SetDisabled sets whether this build feature is disabled or not.
func (f *FeatureDockerSupport) SetDisabled(value bool) {
	f.disabled = value
}

//BuildTypeID is a getter for the Build Type ID associated with this build feature.
func (f *FeatureDockerSupport) BuildTypeID() string {
	return f.buildTypeID
}

//SetBuildTypeID is a setter for the Build Type ID associated with this build feature.
func (f *FeatureDockerSupport) SetBuildTypeID(value string) {
	f.buildTypeID = value
}

//Properties returns a *Properties instance representing a serializable collection to be used.
func (f *FeatureDockerSupport) Properties() *Properties {
	return f.properties
}

//MarshalJSON implements JSON serialization for FeatureDockerSupport
func (f *FeatureDockerSupport) MarshalJSON() ([]byte, error) {
	out := &buildFeatureJSON{
		ID:         f.id,
		Disabled:   NewBool(f.disabled),
		Properties: f.properties,
		Inherited:  NewFalse(),
		Type:       f.Type(),
	}

	return json.Marshal(out)
}

//UnmarshalJSON implements JSON deserialization for FeatureDockerSupport
func (f *FeatureDockerSupport) UnmarshalJSON(data []byte) error {
	var aux buildFeatureJSON
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	f.id = aux.ID

	disabled := aux.Disabled
	if disabled == nil {
		disabled = NewFalse()
	}
	f.disabled = *disabled
	f.properties = NewProperties(aux.Properties.Items...)

	opt, err := DockerSupportOptionsFromProperties(f.properties)
	if err != nil {
		return err
	}

	f.Options = opt

	return nil
}
