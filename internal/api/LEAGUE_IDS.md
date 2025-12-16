# League ID Reference

This document tracks the league IDs used by the FotMob API in this application.

## Supported Leagues

The application currently supports **14 leagues/competitions**:

### Top 5 European Leagues
| League | Country | FotMob ID |
|--------|---------|-----------|
| Premier League | England | 47 |
| La Liga | Spain | 87 |
| Bundesliga | Germany | 54 |
| Serie A | Italy | 55 |
| Ligue 1 | France | 53 |

### European Competitions
| Competition | Type | FotMob ID |
|-------------|------|-----------|
| UEFA Champions League | Club | 42 |
| UEFA Europa League | Club | 73 |
| UEFA Euro | International | 50 |

### South America
| League/Competition | Country | FotMob ID |
|--------------------|---------|-----------|
| Brasileirão Série A | Brazil | 268 |
| Liga Profesional | Argentina | 112 |
| Copa Libertadores | International | 14 |
| Copa America | International | 44 |

### Other
| League/Competition | Country | FotMob ID |
|--------------------|---------|-----------|
| MLS | USA | 130 |
| FIFA World Cup | International | 77 |

## FotMob API League IDs

**Location:** `internal/fotmob/client.go`

```go
SupportedLeagues = []int{
    // Top 5 European Leagues
    47,  // Premier League
    87,  // La Liga
    54,  // Bundesliga
    55,  // Serie A (Italy)
    53,  // Ligue 1
    // European Competitions
    42,  // UEFA Champions League
    73,  // UEFA Europa League
    50,  // UEFA Euro
    // South America
    268, // Brasileirão Série A
    112, // Liga Profesional Argentina
    14,  // Copa Libertadores
    44,  // Copa America
    // Other
    130, // MLS
    77,  // FIFA World Cup
}
```

**API Endpoint:** `https://www.fotmob.com/api/leagues?id={leagueID}&tab={tab}`

Where `tab` can be:
- `fixtures` - Upcoming matches
- `results` - Finished matches

## Notes

- **FotMob** is used for both the **Live Matches** and **Stats** views
- When adding new leagues, update `internal/fotmob/client.go` and this document
- Tournament data (World Cup, Euro, Copa America) only available during competition periods

