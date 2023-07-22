package events_test

import (
	"context"
	"neelbhat88/nest-monitor/m/v2/internal/events"
	"strings"
	"testing"
)

func TestUnmarshalMessage(t *testing.T) {
	ctx := context.Background()

	nonHvac := `{"userId":"AVPHwEv_QPvqecQbRJymI5xY5ZkjHPOxKJuAWoEhz3As","eventId":"ed4da51d-ac48-4f26-8c4b-9672565dd2b6","timestamp":"2023-07-22T04:08:03.601595Z","resourceGroup":["enterprises/21dbe941-4e8e-4647-9772-82669f090fa3/devices/AVPHwEtiY5vtWQ1HuBgvXaMrUbkhYThz7kpGfU6C3rdsdVQxjgNFL48a9-VcsSxg6fKSCR6P4DMRZSQPe3lovUFrx0rg1A"],"resourceUpdate":{"name":"enterprises/21dbe941-4e8e-4647-9772-82669f090fa3/devices/AVPHwEtiY5vtWQ1HuBgvXaMrUbkhYThz7kpGfU6C3rdsdVQxjgNFL48a9-VcsSxg6fKSCR6P4DMRZSQPe3lovUFrx0rg1A","traits":{"sdm.devices.traits.Temperature":{"ambientTemperatureCelsius":23}}}}`

	event, err := events.UnmarshalMessage(ctx, []byte(nonHvac))
	if err != nil {
		t.Fatalf("failed with err: %v", err)
	}

	if event.ResourceUpdate.Traits.ThermostatHVAC.Status != "" {
		t.Fatal("nonHvac event has an hvac status")
	}

	hvac := `{"userId": "AVPHwEv_QPvqecQbRJymI5xY5ZkjHPOxKJuAWoEhz3As", "eventId": "2e747ab3-a564-4d2e-a368-28a735d0b6d7", "timestamp": "2023-07-22T19:23:21.296257Z", "resourceGroup": ["enterprises/21dbe941-4e8e-4647-9772-82669f090fa3/devices/AVPHwEtiY5vtWQ1HuBgvXaMrUbkhYThz7kpGfU6C3rdsdVQxjgNFL48a9-VcsSxg6fKSCR6P4DMRZSQPe3lovUFrx0rg1A"], "resourceUpdate": {"name": "enterprises/21dbe941-4e8e-4647-9772-82669f090fa3/devices/AVPHwEtiY5vtWQ1HuBgvXaMrUbkhYThz7kpGfU6C3rdsdVQxjgNFL48a9-VcsSxg6fKSCR6P4DMRZSQPe3lovUFrx0rg1A", "traits": {"sdm.devices.traits.ThermostatHvac": {"status": "OFF"}}}}`
	hvacEvent, err := events.UnmarshalMessage(ctx, []byte(hvac))
	if err != nil {
		t.Fatalf("failed with err: %v", err)
	}

	if hvacEvent.ResourceUpdate.Traits.ThermostatHVAC.Status == "" {
		t.Fatal("hvac event has no hvac status")
	}

	if strings.ToLower(hvacEvent.ResourceUpdate.Traits.ThermostatHVAC.Status) != "off" {
		t.Fatal("hvac event has no hvac status")
	}
}
