package teamcity

import "encoding/json"

type FeaturePullRequests struct {
	id          string
	vcsRootID   string
	disabled    bool
	buildTypeID string

	Options    FeaturePullRequestsOptions
	properties *Properties
}

//FeaturePullRequestsOptions represents options needed to create a docker support build feature
type FeaturePullRequestsOptions interface {
	Properties() *Properties
}

//ID returns the ID for this instance.
func (f *FeaturePullRequests) ID() string {
	return f.id
}

//SetID sets the ID for this instance.
func (f *FeaturePullRequests) SetID(value string) {
	f.id = value
}

//Type returns the "PullRequests", the keyed-type for this build feature instance
func (f *FeaturePullRequests) Type() string {
	return "pullRequests"
}

//Disabled returns whether this build feature is disabled or not.
func (f *FeaturePullRequests) Disabled() bool {
	return f.disabled
}

//SetDisabled sets whether this build feature is disabled or not.
func (f *FeaturePullRequests) SetDisabled(value bool) {
	f.disabled = value
}

//BuildTypeID is a getter for the Build Type ID associated with this build feature.
func (f *FeaturePullRequests) BuildTypeID() string {
	return f.buildTypeID
}

//SetBuildTypeID is a setter for the Build Type ID associated with this build feature.
func (f *FeaturePullRequests) SetBuildTypeID(value string) {
	f.buildTypeID = value
}

//Properties returns a *Properties instance representing a serializable collection to be used.
func (f *FeaturePullRequests) Properties() *Properties {
	return f.properties
}

//MarshalJSON implements JSON serialization for FeaturePullRequests
func (f *FeaturePullRequests) MarshalJSON() ([]byte, error) {
	out := &buildFeatureJSON{
		ID:         f.id,
		Disabled:   NewBool(f.disabled),
		Properties: f.properties,
		Inherited:  NewFalse(),
		Type:       f.Type(),
	}

	return json.Marshal(out)
}

//UnmarshalJSON implements JSON deserialization for FeaturePullRequests
func (f *FeaturePullRequests) UnmarshalJSON(data []byte) error {
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

	opt, err := PullRequestsOptionsFromProperties(f.properties)
	if err != nil {
		return err
	}

	f.Options = opt

	return nil
}
