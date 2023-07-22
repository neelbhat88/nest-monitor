CREATE TABLE "public"."weather" (
    "event_id" uuid,
    "temperature_f" decimal,
    "temperature_apparent_f" decimal,
    "humidity" integer,
    "pressure_surface_level" decimal,
    "wind_speed" decimal,
    "raw_data" jsonb,
    created timestamp with time zone default now(),
    PRIMARY KEY ("event_id"),
    FOREIGN KEY ("event_id") REFERENCES "public"."hvac_events"("event_id") ON DELETE CASCADE
);