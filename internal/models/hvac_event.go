package models

type HVACEvent struct {
	EventID        string         `json:"eventId"`
	Timestamp      string         `json:"timestamp"`
	ResourceUpdate ResourceUpdate `json:"resourceUpdate"`
}

type ResourceUpdate struct {
	Traits Traits `json:"traits"`
}

type Traits struct {
	ThermostatHVAC ThermostatHVAC `json:"sdm.devices.traits.ThermostatHvac"`
}

type ThermostatHVAC struct {
	Status string `json:"status"`
}
