# Valoguard API

Valoguard API is a service that analyzes Valorant player statistics to detect potential cheaters. It uses the HenrikDev Valorant API to fetch player data and applies a scoring algorithm to determine the likelihood of a player using cheats.

## Features

- Player analysis based on multiple metrics:
  - Headshot percentage
  - Kill/Death ratio
  - Win rate
  - Account level vs. rank correlation
- Detailed player statistics
- Match summary information

## Prerequisites

- Go 1.21 or higher
- HenrikDev API key (get one at [https://docs.henrikdev.xyz/](https://docs.henrikdev.xyz/))

## Setup

1. Clone the repository:
   ```
   git clone https://github.com/edcnogueira/valoguard-api.git
   cd valoguard-api
   ```

2. Create a `.env` file based on the example:
   ```
   cp .env.example .env
   ```

3. Edit the `.env` file and add your HenrikDev API key:
   ```
   HENRIK_API_KEY=your_api_key_here
   ```

4. Install dependencies:
   ```
   go mod download
   ```

## Running the API

1. Set the environment variables:
   ```
   export HENRIK_API_KEY=your_api_key_here
   ```

2. Build and run the API:
   ```
   cd cmd
   go build -o valoguard-api
   ./valoguard-api
   ```

   The API will start on port 8080 by default. You can change the port by setting the `PORT` environment variable.

The script will build and run the API, then make a test request to analyze a player.

## API Endpoints

### Analyze Player

```
GET /analyze/:name/:tag?region=:region
```

Parameters:
- `name`: Player's in-game name
- `tag`: Player's tag (the numbers after the #)
- `region` (optional): Player's region (default: eu)
  - Valid regions: eu, na, ap, kr, latam, br

Example:
```
GET /analyze/TenZ/1337?region=na
```

Response:
```json
{
  "probability": 70,
  "stats": {
    "account_level": 125,
    "avg_damage": 156.8,
    "hs_percent": 38.5,
    "kd": 2.7,
    "rank": "Radiant",
    "win_rate": 75.0
  },
  "matches_summary": [
    {
      "match_id": "abc123",
      "score": 350
    },
    {
      "match_id": "def456",
      "score": 280
    }
  ]
}
```

The `probability` field represents the likelihood of the player using cheats, on a scale from 0 to 100.

## How It Works

The API calculates a "cheat score" based on several factors:
- Headshot percentage > 35%: +30 points
- K/D ratio > 2.5: +20 points
- Win rate > 70%: +20 points
- Low account level (< 50) with high rank (> Gold 1): +30 points

The maximum score is 100, indicating a high probability of cheating.

## Building the CLI with an embedded key (via Makefile)

You can compile the CLI embedding the Henrik API key directly into the binary using the Makefile.

Attention: replace HDEV-xxxx with your real key.

### How to build the CLI with the embedded key
- macOS/Linux/Windows (via Make):
    - make build-cli API_KEY=HDEV-xxxx

This will generate the binary at: bin/valoguard

### How to run the CLI
- bin/valoguard player cheating analyze --name "YourName" --tag "1234" [--region br]

You can still override at runtime:
- Via flag: bin/valoguard --api-key "your_key" player cheating analyze --name "YourName" --tag "1234"
- Or via env: HENRIK_API_KEY="your_key" bin/valoguard player cheating analyze --name "YourName" --tag "1234"

Important notes:
- Precedence (CLI): flag --api-key > environment variable HENRIK_API_KEY > embedded key (Makefile/ldflags).
- Using HENRIK_API_KEY=... go build does NOT embed the key; it only sets an env var for the build process. To embed, use the Makefile above.

Note: The HTTP server needs no changes and continues to work as before (provide HENRIK_API_KEY via environment variable when running it).


## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- [HenrikDev Valorant API](https://docs.henrikdev.xyz/) for providing the data
- [Fiber](https://github.com/gofiber/fiber) for the web framework
- [Cobra](https://github.com/spf13/cobra) for the CLI framework