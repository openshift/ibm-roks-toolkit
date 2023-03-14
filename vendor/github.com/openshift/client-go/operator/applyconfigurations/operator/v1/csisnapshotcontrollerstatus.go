// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1

// CSISnapshotControllerStatusApplyConfiguration represents an declarative configuration of the CSISnapshotControllerStatus type for use
// with apply.
type CSISnapshotControllerStatusApplyConfiguration struct {
	OperatorStatusApplyConfiguration `json:",inline"`
}

// CSISnapshotControllerStatusApplyConfiguration constructs an declarative configuration of the CSISnapshotControllerStatus type for use with
// apply.
func CSISnapshotControllerStatus() *CSISnapshotControllerStatusApplyConfiguration {
	return &CSISnapshotControllerStatusApplyConfiguration{}
}

// WithObservedGeneration sets the ObservedGeneration field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ObservedGeneration field is set to the value of the last call.
func (b *CSISnapshotControllerStatusApplyConfiguration) WithObservedGeneration(value int64) *CSISnapshotControllerStatusApplyConfiguration {
	b.ObservedGeneration = &value
	return b
}

// WithConditions adds the given value to the Conditions field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Conditions field.
func (b *CSISnapshotControllerStatusApplyConfiguration) WithConditions(values ...*OperatorConditionApplyConfiguration) *CSISnapshotControllerStatusApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithConditions")
		}
		b.Conditions = append(b.Conditions, *values[i])
	}
	return b
}

// WithVersion sets the Version field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Version field is set to the value of the last call.
func (b *CSISnapshotControllerStatusApplyConfiguration) WithVersion(value string) *CSISnapshotControllerStatusApplyConfiguration {
	b.Version = &value
	return b
}

// WithReadyReplicas sets the ReadyReplicas field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ReadyReplicas field is set to the value of the last call.
func (b *CSISnapshotControllerStatusApplyConfiguration) WithReadyReplicas(value int32) *CSISnapshotControllerStatusApplyConfiguration {
	b.ReadyReplicas = &value
	return b
}

// WithGenerations adds the given value to the Generations field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Generations field.
func (b *CSISnapshotControllerStatusApplyConfiguration) WithGenerations(values ...*GenerationStatusApplyConfiguration) *CSISnapshotControllerStatusApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithGenerations")
		}
		b.Generations = append(b.Generations, *values[i])
	}
	return b
}
