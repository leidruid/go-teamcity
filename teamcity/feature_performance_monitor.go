package teamcity

import (
	"encoding/json"
)

//FeaturePerformanceMonitor represents a golang build feature. Implements BuildFeature interface
type FeaturePerformanceMonitor struct {
	id          string
	disabled    bool
	buildTypeID string

	properties *Properties
}

// NewPerformanceMonitor returns a new instance of the FeaturePerformanceMonitor struct
func NewPerformanceMonitor() (*FeaturePerformanceMonitor, error) {
	return &FeaturePerformanceMonitor{
		properties: NewProperties(),
	}, nil
}

//ID returns the ID for this instance.
func (f *FeaturePerformanceMonitor) ID() string {
	return f.id
}

//SetID sets the ID for this instance.
func (f *FeaturePerformanceMonitor) SetID(value string) {
	f.id = value
}

//Type returns the "commit-status-publisher", the keyed-type for this build feature instance
func (f *FeaturePerformanceMonitor) Type() string {
	return "perfmon"
}

//Disabled returns whether this build feature is disabled or not.
func (f *FeaturePerformanceMonitor) Disabled() bool {
	return f.disabled
}

//SetDisabled sets whether this build feature is disabled or not.
func (f *FeaturePerformanceMonitor) SetDisabled(value bool) {
	f.disabled = value
}

//BuildTypeID is a getter for the Build Type ID associated with this build feature.
func (f *FeaturePerformanceMonitor) BuildTypeID() string {
	return f.buildTypeID
}

//SetBuildTypeID is a setter for the Build Type ID associated with this build feature.
func (f *FeaturePerformanceMonitor) SetBuildTypeID(value string) {
	f.buildTypeID = value
}

//Properties returns a *Properties instance representing a serializable collection to be used.
func (f *FeaturePerformanceMonitor) Properties() *Properties {
	return f.properties
}

//MarshalJSON implements JSON serialization for FeatureCommitStatusPublisher
func (f *FeaturePerformanceMonitor) MarshalJSON() ([]byte, error) {
	out := &buildFeatureJSON{
		ID:         f.id,
		Disabled:   NewBool(f.disabled),
		Properties: f.properties,
		Inherited:  NewFalse(),
		Type:       f.Type(),
	}

	return json.Marshal(out)
}

//UnmarshalJSON implements JSON deserialization for FeatureCommitStatusPublisher
func (f *FeaturePerformanceMonitor) UnmarshalJSON(data []byte) error {
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
	//f.properties = NewProperties(aux.Properties.Items...)
	f.properties = nil

	return nil
}
