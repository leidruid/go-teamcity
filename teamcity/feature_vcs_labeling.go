package teamcity

import "encoding/json"

type FeatureVcsLabeling struct {
	id          string
	disabled    bool
	buildTypeID string
	Options     FeatureVcsLabelingOptions

	properties *Properties
}

//FeatureVcsLabelingOptions represents options needed to create a vcs labeling build feature
type FeatureVcsLabelingOptions interface {
	Properties() *Properties
}

//ID returns the ID for this instance.
func (f *FeatureVcsLabeling) ID() string {
	return f.id
}

//SetID sets the ID for this instance.
func (f *FeatureVcsLabeling) SetID(value string) {
	f.id = value
}

//Type returns the "VcsLabeling", the keyed-type for this build feature instance
func (f *FeatureVcsLabeling) Type() string {
	return "VcsLabeling"
}

//Disabled returns whether this build feature is disabled or not.
func (f *FeatureVcsLabeling) Disabled() bool {
	return f.disabled
}

//SetDisabled sets whether this build feature is disabled or not.
func (f *FeatureVcsLabeling) SetDisabled(value bool) {
	f.disabled = value
}

//BuildTypeID is a getter for the Build Type ID associated with this build feature.
func (f *FeatureVcsLabeling) BuildTypeID() string {
	return f.buildTypeID
}

//SetBuildTypeID is a setter for the Build Type ID associated with this build feature.
func (f *FeatureVcsLabeling) SetBuildTypeID(value string) {
	f.buildTypeID = value
}

//Properties returns a *Properties instance representing a serializable collection to be used.
func (f *FeatureVcsLabeling) Properties() *Properties {
	return f.properties
}

//MarshalJSON implements JSON serialization for FeatureVcsLabeling
func (f *FeatureVcsLabeling) MarshalJSON() ([]byte, error) {
	out := &buildFeatureJSON{
		ID:         f.id,
		Disabled:   NewBool(f.disabled),
		Properties: f.properties,
		Inherited:  NewFalse(),
		Type:       f.Type(),
	}

	return json.Marshal(out)
}

//UnmarshalJSON implements JSON deserialization for FeatureVcsLabeling
func (f *FeatureVcsLabeling) UnmarshalJSON(data []byte) error {
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

	opt, err := VcsLabelingOptionsFromProperties(f.properties)
	if err != nil {
		return err
	}

	f.Options = opt

	return nil
}
