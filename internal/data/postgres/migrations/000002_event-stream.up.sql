create table public.event_stream (
    id uuid primary key default gen_random_uuid(),
    created timestamp with time zone default now(),
    event_message jsonb default NULL    
);
