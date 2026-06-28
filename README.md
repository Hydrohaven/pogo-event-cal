# What
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
    
    %% Styling
    style Cron fill:#f9f,stroke:#333,stroke-width:2px
    style App fill:#bbf,stroke:#333,stroke-width:2px
    style Cache fill:#ffb,stroke:#333,stroke-width:2px
```


# Why
I wanted to learn Go in a fun way and really didn't wanna go through Tour of Go cuz it looked REALLY long and boring, so here I am. I also really love Pokemon Go so this seem liked a perfect way to learn the language!

# Future Additions
- Online dashboard to query events and see what pokemon they have, what region they're from, etc.