create table public.hvac_events (
    event_id uuid primary key default gen_random_uuid(),
    event_timestamp timestamp default NULL,
    hvac_status text not null default ''
);