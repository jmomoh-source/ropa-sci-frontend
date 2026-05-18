package models

import (
    "encoding/json"
    "fmt"
    "os"
    "strings"
    "math/rand"
)

const dataFile = "data/players.json"

// SavePlayer writes a new player to disk
// Returns an error if the username already exists
func SavePlayer(p Player) error {
    // Create data/ folder first if it doesn't exist
    if err := os.MkdirAll("data", 0755); err != nil {
        return fmt.Errorf("could not create data folder: %w", err)
    }

    // Load existing players
    players, _ := LoadPlayers()

    // Check for duplicate username before appending
    for _, existing := range players {
        if strings.EqualFold(existing.Username, p.Username) {
            return fmt.Errorf("username '%s' already exists", p.Username)
        }
    }

    // Add the new player
    players = append(players, p)

    // Write back to file with readable formatting
    data, err := json.MarshalIndent(players, "", "  ")
    if err != nil {
        return err
    }

    return os.WriteFile(dataFile, data, 0644)
}

// LoadPlayers reads all players from disk
func LoadPlayers() ([]Player, error) {
    data, err := os.ReadFile(dataFile)
    if err != nil {
        return []Player{}, nil // file doesn't exist yet — that's fine
    }

    var players []Player
    err = json.Unmarshal(data, &players)
    return players, err
}

// UsernameExists checks if a Gitea username is already registered
func UsernameExists(username string) bool {
    players, err := LoadPlayers()
    if err != nil {
        return false
    }
    for _, p := range players {
        if strings.EqualFold(p.Username, username) {
            return true
        }
    }
    return false
}

// FindPlayerByUsername returns a player if they exist
// Returns the player, a found bool, and any error
func FindPlayerByUsername(username string) (Player, bool, error) {
    players, err := LoadPlayers()
    if err != nil {
        return Player{}, false, err
    }
    for _, p := range players {
        if strings.EqualFold(p.Username, username) {
            return p, true, nil
        }
    }
    return Player{}, false, nil
}

// UpdatePlayer finds an existing player by username and overwrites their record
// Used to save lifetime stats after every match ends
func UpdatePlayer(p Player) error {
    players, err := LoadPlayers()
    if err != nil {
        return fmt.Errorf("could not load players: %w", err)
    }

    found := false
    for i, existing := range players {
        if strings.EqualFold(existing.Username, p.Username) {
            players[i] = p // overwrite with updated data
            found = true
            break
        }
    }

    if !found {
        return fmt.Errorf("player '%s' not found", p.Username)
    }

    data, err := json.MarshalIndent(players, "", "  ")
    if err != nil {
        return err
    }

    return os.WriteFile(dataFile, data, 0644)
}

// GenerateRoomCode creates a random 4-digit room code
func GenerateRoomCode() string {
    digits := rand.Intn(9000) + 1000 // always 4 digits: 1000-9999
    return fmt.Sprintf("RPS-%d", digits)
}