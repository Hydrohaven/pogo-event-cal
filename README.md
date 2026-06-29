# About
A Pokemon Go CLI tool used to automatically add new game events on a Google Calendar anyone can subscribe to

# How
```mermaid
graph TD
    %% Trigger Mechanism
    Cron[Cron Job: Daily/Weekly] -->|Triggers| App[Go CLI Binary]

    %% Ingestion Layer
    App -->|1. HTTP GET| NetHTTP[net/http Package]
    NetHTTP -->|Fetches HTML| LeekDuck[leekduck.com]
    LeekDuck -->|Returns HTML Data| GoQuery[goquery Parser]
    
    %% Processing Layer
    GoQuery -->|2. Traverses DOM| Extract[Extract: Title, Date, Time]
    Extract -->|3. Compares State| Cache{Local Cache File}
    
    %% Output Layer
    Cache -->|New Event| GCal[Google Calendar API]
    Cache -->|Duplicate| Skip[Skip Event]
    GCal -->|4. POST Request| Calendar[Pokemon Go Event Calendar]
    
    %% Subtle/Professional Styling
    classDef default fill:#f4f5f6,stroke:#c5c7c9,stroke-width:1px,color:#1e293b;
    classDef trigger fill:#e0f2fe,stroke:#38bdf8,stroke-width:1px,color:#0369a1;
    classDef core fill:#f0fdf4,stroke:#4ade80,stroke-width:1px,color:#14532d;
    classDef logic fill:#fff7ed,stroke:#fb923c,stroke-width:1px,color:#7c2d12;
    
    class Cron trigger;
    class App core;
    class Cache logic;
```


# Why
I wanted to learn Go in a fun way and really didn't wanna go through Tour of Go cuz it looked REALLY long and boring, so here I am. I also really love Pokemon Go so this seem liked a perfect way to learn the language!

# Future Additions
- Web-based dashboard to query all logged events to view featured pokemon, what region they're from, etc. for collection/research tasks