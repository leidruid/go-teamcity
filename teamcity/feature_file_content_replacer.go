package teamcity

import "encoding/json"

type FeatureFileContentReplacer struct {
	id          string
	disabled    bool
	buildTypeID string
	Options     FeatureFileContentReplacerOptions

	properties *Properties
}

//FeatureFileContentReplacerOptions represents options needed to create a docker support build feature
type FeatureFileContentReplacerOptions interface {
	Properties() *Properties
}

//ID returns the ID for this instance.
func (f *FeatureFileContentReplacer) ID() string {
	return f.id
}

//SetID sets the ID for this instance.
func (f *FeatureFileContentReplacer) SetID(value string) {
	f.id = value
}

//Type returns the "FileContentReplacer", the keyed-type for this build feature instance
func (f *FeatureFileContentReplacer) Type() string {
	return "JetBrains.FileContentReplacer"
}

//Disabled returns whether this build feature is disabled or not.
func (f *FeatureFileContentReplacer) Disabled() bool {
	return f.disabled
}

//SetDisabled sets whether this build feature is disabled or not.
func (f *FeatureFileContentReplacer) SetDisabled(value bool) {
	f.disabled = value
}

//BuildTypeID is a getter for the Build Type ID associated with this build feature.
func (f *FeatureFileContentReplacer) BuildTypeID() string {
	return f.buildTypeID
}

//SetBuildTypeID is a setter for the Build Type ID associated with this build feature.
func (f *FeatureFileContentReplacer) SetBuildTypeID(value string) {
	f.buildTypeID = value
}

//Properties returns a *Properties instance representing a serializable collection to be used.
func (f *FeatureFileContentReplacer) Properties() *Properties {
	return f.properties
}

//MarshalJSON implements JSON serialization for FeatureFileContentReplacer
func (f *FeatureFileContentReplacer) MarshalJSON() ([]byte, error) {
	out := &buildFeatureJSON{
		ID:         f.id,
		Disabled:   NewBool(f.disabled),
		Properties: f.properties,
		Inherited:  NewFalse(),
		Type:       f.Type(),
	}

	return json.Marshal(out)
}

//UnmarshalJSON implements JSON deserialization for FeatureFileContentReplacer
func (f *FeatureFileContentReplacer) UnmarshalJSON(data []byte) error {
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

	opt, err := FileContentReplacerOptionsFromProperties(f.properties)
	if err != nil {
		return err
	}

	f.Options = opt

	return nil
}
