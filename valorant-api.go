package main

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// 01bb38e1-da47-4e6a-9b3d-945fe4655707	Agents
// f85cb6f7-33e5-4dc8-b609-ec7212301948	Contracts
// d5f120f8-ff8c-4aac-92ea-f2b5acbe9475	Sprays
// dd3bf334-87f3-40bd-b043-682a57a8dc3a	Gun Buddies
// 3f296c07-64c3-494c-923b-fe692a4fa1bd	Cards
// e7c63390-eda7-46e0-bb7a-a6abdacd2433	Skins
// 3ad1b2b2-acdb-4524-852f-954a76ddae0a	Skin Variants
// de7caa6b-adf7-4588-bbd1-143831e786c6	Titles

const (
	ItemTypeID_Agents       = "01bb38e1-da47-4e6a-9b3d-945fe4655707"
	ItemTypeID_Contracts    = "f85cb6f7-33e5-4dc8-b609-ec7212301948"
	ItemTypeID_Sprays       = "d5f120f8-ff8c-4aac-92ea-f2b5acbe9475"
	ItemTypeID_GunBuddies   = "dd3bf334-87f3-40bd-b043-682a57a8dc3a"
	ItemTypeID_Cards        = "3f296c07-64c3-494c-923b-fe692a4fa1bd"
	ItemTypeID_Skins        = "e7c63390-eda7-46e0-bb7a-a6abdacd2433"
	ItemTypeID_SkinVariants = "3ad1b2b2-acdb-4524-852f-954a76ddae0a"
	ItemTypeID_Titles       = "de7caa6b-adf7-4588-bbd1-143831e786c6"
)

type ValorantAPI struct {
	Shard       string
	Region      string
	Client      string
	Version     string
	Token       string
	AccessToken string
	PlayerUUID  string
	MatchID     string
}

type MainData struct {
	OwnedItems    []OwnedItemsResponse
	PlayerLoadout *PlayerLoadoutResponse
}

func NewValorantAPI() *ValorantAPI {
	return &ValorantAPI{
		Shard:       "",
		Region:      "",
		Client:      "",
		Version:     "",
		Token:       "",
		AccessToken: "",
		PlayerUUID:  "",
	}
}

func (v *ValorantAPI) GetMainData() (*MainData, error) {

	// %LocalAppData%\VALORANT\Saved\Logs\ShooterGame.log

	appData := os.Getenv("LOCALAPPDATA")
	if appData == "" {
		return nil, fmt.Errorf("LOCALAPPDATA environment variable is not set")
	}

	valorantPath := path.Join(appData, "VALORANT", "Saved", "Logs", "ShooterGame.log")
	if _, err := os.Stat(valorantPath); os.IsNotExist(err) {
		fmt.Println("Valorant shootergame file does not exist")
		return nil, fmt.Errorf("valorant shootergame file does not exist")
	}

	file, err := os.Open(valorantPath)
	if err != nil {
		fmt.Println("Error opening valorant shootergame file")
		return nil, fmt.Errorf("error opening valorant shootergame file: %w", err)
	}
	defer file.Close()

	// Read whole file

	content, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading valorant shootergame file")
		return nil, fmt.Errorf("error reading valorant shootergame file: %w", err)
	}

	re := regexp.MustCompile(`https://glz-(.+?)-1.(.+?).a.pvp.net`)
	matches := re.FindStringSubmatch(string(content))

	if len(matches) == 0 {
		fmt.Println("No matches found")
		return nil, fmt.Errorf("no matches found")
	}

	v.Region = matches[1]
	v.Shard = matches[2]

	if v.Region == "" || v.Shard == "" {
		fmt.Println("Region or shard is empty")
		return nil, fmt.Errorf("region or shard is empty")
	}

	v.Client = "ew0KCSJwbGF0Zm9ybVR5cGUiOiAiUEMiLA0KCSJwbGF0Zm9ybU9TIjogIldpbmRvd3MiLA0KCSJwbGF0Zm9ybU9TVmVyc2lvbiI6ICIxMC4wLjE5MDQyLjEuMjU2LjY0Yml0IiwNCgkicGxhdGZvcm1DaGlwc2V0IjogIlVua25vd24iDQp9"
	v.Version = "11.08.00.3880807"

	err = v.waitForLockfile()
	if err != nil {
		fmt.Println("Error waiting for lockfile:", err)
		return nil, fmt.Errorf("error waiting for lockfile: %w", err)
	}
	fmt.Println("Lockfile detected!")

	lockfileData, err := v.getLockfileDataWithRetry()
	if err != nil {
		fmt.Println("Error getting lockfile data:", err)
		return nil, fmt.Errorf("error getting lockfile data: %w", err)
	}

	fmt.Println("Lockfile data:", lockfileData)

	entitlementsToken, err := v.getEntitlementsTokenWithRetry(lockfileData.Port, lockfileData.Password)
	if err != nil {
		fmt.Println("Error getting session data:", err)
		return nil, fmt.Errorf("error getting session data: %w", err)
	}

	fmt.Println("Entitlements Token:", entitlementsToken)

	v.AccessToken = entitlementsToken.AccessToken
	v.Token = entitlementsToken.Token
	v.PlayerUUID = entitlementsToken.Subject

	ownedItems, err := v.GetOwnedItems(ItemTypeID_Skins)
	if err != nil {
		fmt.Println("Error getting owned items:", err)
		return nil, fmt.Errorf("error getting owned items: %w", err)
	}

	playerLoadout, err := v.GetPlayerLoadout()
	if err != nil {
		fmt.Println("Error getting player loadout:", err)
		return nil, fmt.Errorf("error getting player loadout: %w", err)
	}

	return &MainData{
		OwnedItems:    []OwnedItemsResponse{*ownedItems},
		PlayerLoadout: playerLoadout,
	}, nil
}

// waitForLockfile waits for the Riot Client lockfile to appear
func (v *ValorantAPI) waitForLockfile() error {
	appData := os.Getenv("LOCALAPPDATA")
	if appData == "" {
		return fmt.Errorf("LOCALAPPDATA environment variable is not set")
	}

	configPath := filepath.Join(appData, "Riot Games", "Riot Client", "Config")
	lockfilePath := filepath.Join(configPath, "lockfile")

	// Check if lockfile already exists
	if _, err := os.Stat(lockfilePath); err == nil {
		return nil
	}

	// Watch the directory for file changes
	watcher, err := os.Open(configPath)
	if err != nil {
		return fmt.Errorf("failed to open config directory: %v", err)
	}
	defer watcher.Close()

	// Poll for the lockfile every 100ms
	for {
		time.Sleep(100 * time.Millisecond)
		if _, err := os.Stat(lockfilePath); err == nil {
			return nil
		}
	}
}

// LockfileData represents the data from the Riot Client lockfile
type LockfileData struct {
	Name     string
	PID      string
	Port     string
	Password string
	Protocol string
}

type EntitlementsTokenResponse struct {
	/** Used as the token in requests */
	AccessToken  string
	Entitlements []interface{}
	Issuer       string
	Subject      string
	Token        string
}

// getLockfileData reads and parses the Riot Client lockfile
func (v *ValorantAPI) getLockfileData() (*LockfileData, error) {
	appData := os.Getenv("LOCALAPPDATA")
	if appData == "" {
		return nil, fmt.Errorf("LOCALAPPDATA environment variable is not set")
	}

	lockfilePath := filepath.Join(appData, "Riot Games", "Riot Client", "Config", "lockfile")
	contents, err := os.ReadFile(lockfilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read lockfile: %v", err)
	}

	parts := strings.Split(strings.TrimSpace(string(contents)), ":")
	if len(parts) != 5 {
		return nil, fmt.Errorf("invalid lockfile format")
	}

	return &LockfileData{
		Name:     parts[0],
		PID:      parts[1],
		Port:     parts[2],
		Password: parts[3],
		Protocol: parts[4],
	}, nil
}

// getSession fetches session data from the Riot Client API
func (v *ValorantAPI) getEntitlementsToken(port, password string) (*EntitlementsTokenResponse, error) {
	url := fmt.Sprintf("https://127.0.0.1:%s/entitlements/v1/token", port)

	// Create basic auth header
	auth := base64.StdEncoding.EncodeToString([]byte("riot:" + password))

	// Create HTTP client with TLS config to skip certificate verification
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Basic "+auth)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var result EntitlementsTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// getLockfileDataWithRetry attempts to get lockfile data with retry logic
func (v *ValorantAPI) getLockfileDataWithRetry() (*LockfileData, error) {
	var lockData *LockfileData
	var err error

	for {
		lockData, err = v.getLockfileData()
		if err == nil {
			break
		}

		fmt.Println("Waiting for lockfile...")
		if waitErr := v.waitForLockfile(); waitErr != nil {
			return nil, waitErr
		}
	}

	return lockData, nil
}

// getSessionDataWithRetry attempts to get session data with retry logic
func (v *ValorantAPI) getEntitlementsTokenWithRetry(port, password string) (*EntitlementsTokenResponse, error) {
	var entitlementsToken *EntitlementsTokenResponse
	var err error
	lastRetryMessage := time.Now()

	for {
		entitlementsToken, err = v.getEntitlementsToken(port, password)
		if err == nil {
			// Check if session is loaded
			if entitlementsToken.AccessToken == "" {
				time.Sleep(1500 * time.Millisecond)
				entitlementsToken = nil
				continue
			}
			break
		}

		currentTime := time.Now()
		if currentTime.Sub(lastRetryMessage) > time.Second {
			fmt.Println("Unable to get session data, retrying...")
			lastRetryMessage = currentTime
		}
	}

	return entitlementsToken, nil
}

type Gun struct {
	ID              string
	CharmInstanceID *string
	CharmID         *string
	CharmLevelID    *string
	SkinID          *string
	SkinLevelID     *string
	ChromaID        *string
	Attachments     []GunAttachment
}
type GunAttachment struct {
	ID string
}
type Spray struct {
	SprayID      string
	SprayLevelID *string
	EquipSlotID  string
}
type Identity struct {
	PlayerCardID           string
	PlayerTitleID          string
	AccountLevel           int
	PreferredLevelBorderID string
	HideAccountLevel       bool
}
type PlayerLoadoutResponse struct {
	Guns      []Gun
	Sprays    []Spray
	Identity  Identity
	Incognito bool
}

func (v *ValorantAPI) GetPlayerLoadout() (*PlayerLoadoutResponse, error) {

	headers := map[string]string{
		"X-Riot-Entitlements-Jwt": v.Token,
		"X-Riot-ClientPlatform":   v.Client,
		"X-Riot-ClientVersion":    v.Version,
		"Authorization":           "Bearer " + v.AccessToken,
	}

	req, err := http.NewRequest("GET", "https://pd."+v.Region+".a.pvp.net/personalization/v2/players/"+v.PlayerUUID+"/playerloadout", nil)
	if err != nil {
		fmt.Println("failed to create request: %w", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("failed to make request: %w", err)
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("failed to get player loadout: %w", err)
		return nil, fmt.Errorf("failed to get player loadout: %w", err)
	}

	fmt.Println(string(body))

	var playerLoadout PlayerLoadoutResponse
	if err := json.Unmarshal(body, &playerLoadout); err != nil {
		fmt.Println("failed to unmarshal player loadout: %w", err)
		return nil, fmt.Errorf("failed to unmarshal player loadout: %w", err)
	}

	return &playerLoadout, nil

}

func (v *ValorantAPI) SetPlayerLoadout(playerLoadout *PlayerLoadoutResponse) (*PlayerLoadoutResponse, error) {
	headers := map[string]string{
		"X-Riot-Entitlements-Jwt": v.Token,
		"X-Riot-ClientPlatform":   v.Client,
		"X-Riot-ClientVersion":    v.Version,
		"Authorization":           "Bearer " + v.AccessToken,
	}

	requestBody, err := json.Marshal(playerLoadout)
	if err != nil {
		fmt.Println("failed to marshal player loadout: %w", err)
		return nil, fmt.Errorf("failed to marshal player loadout: %w", err)
	}

	fmt.Println(" -------------------------------- REQUEST --------------------------------")
	fmt.Println(string(requestBody))

	req, err := http.NewRequest("PUT", "https://pd."+v.Region+".a.pvp.net/personalization/v2/players/"+v.PlayerUUID+"/playerloadout", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("failed to create request: %w", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("failed to make request: %w", err)
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("failed to get player loadout: %w", err)
		return nil, fmt.Errorf("failed to get player loadout: %w", err)
	}

	fmt.Println(" -------------------------------- RESPONSE --------------------------------")
	fmt.Println(string(body))
	fmt.Println("--------------------------------")

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to set player loadout: %w", err)
	}

	var playerLoadoutResponse PlayerLoadoutResponse
	if err := json.Unmarshal(body, &playerLoadoutResponse); err != nil {
		fmt.Println("failed to unmarshal player loadout: %w", err)
		return nil, fmt.Errorf("failed to unmarshal player loadout: %w", err)
	}

	return &playerLoadoutResponse, nil
}

type OwnedItemsResponse struct {
	ItemTypeID   string                          `json:"ItemTypeID"`
	Entitlements []OwnedItemsResponseEntitlement `json:"Entitlements"`
}
type OwnedItemsResponseEntitlement struct {
	TypeID     string `json:"TypeID"`
	ItemID     string `json:"ItemID"`
	InstanceID string `json:"InstanceID"`
}

func (v *ValorantAPI) GetOwnedItems(ItemTypeID string) (*OwnedItemsResponse, error) {
	headers := map[string]string{
		"X-Riot-Entitlements-Jwt": v.Token,
		"X-Riot-ClientPlatform":   v.Client,
		"X-Riot-ClientVersion":    v.Version,
		"Authorization":           "Bearer " + v.AccessToken,
	}
	// pd.{shard}.a.pvp.net/store/v1/entitlements/{puuid}/{ItemTypeID}
	url := fmt.Sprintf("https://pd.%s.a.pvp.net/store/v1/entitlements/%s/%s", v.Region, v.PlayerUUID, ItemTypeID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("failed to create request: %w", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("failed to make request: %w", err)
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("failed to get owned items: %w", err)
		return nil, fmt.Errorf("failed to get owned items: %w", err)
	}

	var ownedItems OwnedItemsResponse
	if err := json.Unmarshal(body, &ownedItems); err != nil {
		fmt.Println("failed to unmarshal owned items: %w", err)
		return nil, fmt.Errorf("failed to unmarshal owned items: %w", err)
	}

	return &ownedItems, nil
}

type PreGamePlayerResponse struct {
	Subject string
	MatchID string
	Version int
}

type PreGamePlayerResponseError struct {
	HTTPStatus int    `json:"httpStatus"`
	ErrorCode  string `json:"errorCode"`
	Message    string `json:"message"`
}

func (v *ValorantAPI) GetPreGamePlayer() (*PreGamePlayerResponse, error) {
	headers := map[string]string{
		"X-Riot-Entitlements-Jwt": v.Token,
		"X-Riot-ClientPlatform":   v.Client,
		"X-Riot-ClientVersion":    v.Version,
		"Authorization":           "Bearer " + v.AccessToken,
	}
	// https://glz-{region}-1.{shard}.a.pvp.net/pregame/v1/players/{puuid}
	url := fmt.Sprintf("https://glz-%s-1.%s.a.pvp.net/pregame/v1/players/%s", v.Region, v.Shard, v.PlayerUUID)
	fmt.Println(" -------------------------------- URL --------------------------------")
	fmt.Println(url)
	fmt.Println("--------------------------------")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("failed to create request: %w", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("failed to make request: %w", err)
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("failed to get pre game player: %w", err)
		return nil, fmt.Errorf("failed to get pre game player: %w", err)
	}

	fmt.Println(" -------------------------------- RESPONSE PLAYERS --------------------------------")
	fmt.Println(string(body))
	fmt.Println("--------------------------------")

	var preGamePlayer PreGamePlayerResponse
	if err := json.Unmarshal(body, &preGamePlayer); err != nil {
		fmt.Println("failed to unmarshal pre game player: %w", err)
		return nil, fmt.Errorf("failed to unmarshal pre game player: %w", err)
	}

	if preGamePlayer.Subject == "" {
		return nil, fmt.Errorf("failed to get pre game player: %w", err)
	}

	v.MatchID = preGamePlayer.MatchID

	return &preGamePlayer, nil
}

// PreGameMatchResponse represents the pre-game match data
type PreGameMatchResponse struct {
	// Pre-Game Match ID
	ID                 string        `json:"ID"`
	Version            int           `json:"Version"`
	Teams              []Team        `json:"Teams"`
	AllyTeam           *Team         `json:"AllyTeam"`
	EnemyTeam          *Team         `json:"EnemyTeam"`
	ObserverSubjects   []interface{} `json:"ObserverSubjects"`
	MatchCoaches       []interface{} `json:"MatchCoaches"`
	EnemyTeamSize      int           `json:"EnemyTeamSize"`
	EnemyTeamLockCount int           `json:"EnemyTeamLockCount"`
	PregameState       string        `json:"PregameState"` // "character_select_active" | "provisioned"
	// Date in ISO 8601 format
	LastUpdated string `json:"LastUpdated"`
	// Map ID
	MapID          string        `json:"MapID"`
	MapSelectPool  []interface{} `json:"MapSelectPool"`
	BannedMapIDs   []interface{} `json:"BannedMapIDs"`
	CastedVotes    *interface{}  `json:"CastedVotes,omitempty"`
	MapSelectSteps []interface{} `json:"MapSelectSteps"`
	MapSelectStep  int           `json:"MapSelectStep"`
	Team1          string        `json:"Team1"` // ("Blue" | "Red") | string
	GamePodID      string        `json:"GamePodID"`
	// Game Mode
	Mode           string `json:"Mode"`
	VoiceSessionID string `json:"VoiceSessionID"`
	MUCName        string `json:"MUCName"`
	// JWT containing match ID and player IDs
	TeamMatchToken string `json:"TeamMatchToken"`
	// Queue ID
	QueueID              string      `json:"QueueID"`
	ProvisioningFlowID   string      `json:"ProvisioningFlowID"` // "Matchmaking" | "CustomGame"
	IsRanked             bool        `json:"IsRanked"`
	PhaseTimeRemainingNS int64       `json:"PhaseTimeRemainingNS"`
	StepTimeRemainingNS  int64       `json:"StepTimeRemainingNS"`
	AltModesFlagADA      bool        `json:"altModesFlagADA"`
	TournamentMetadata   interface{} `json:"TournamentMetadata"`
	RosterMetadata       interface{} `json:"RosterMetadata"`
}

// Team represents a team in the pre-game match
type Team struct {
	TeamID  string   `json:"TeamID"` // ("Blue" | "Red") | string
	Players []Player `json:"Players"`
}

// Player represents a player in the pre-game match
type Player struct {
	// Player UUID
	Subject string `json:"Subject"`
	// Character ID
	CharacterID             string            `json:"CharacterID"`
	CharacterSelectionState string            `json:"CharacterSelectionState"` // "" | "selected" | "locked"
	PregamePlayerState      string            `json:"PregamePlayerState"`      // "joined"
	CompetitiveTier         int               `json:"CompetitiveTier"`
	PlayerIdentity          PlayerIdentity    `json:"PlayerIdentity"`
	SeasonalBadgeInfo       SeasonalBadgeInfo `json:"SeasonalBadgeInfo"`
	IsCaptain               bool              `json:"IsCaptain"`
}

// PlayerIdentity represents player identity information
type PlayerIdentity struct {
	// Player UUID
	Subject string `json:"Subject"`
	// Card ID
	PlayerCardID string `json:"PlayerCardID"`
	// Title ID
	PlayerTitleID string `json:"PlayerTitleID"`
	AccountLevel  int    `json:"AccountLevel"`
	// Preferred Level Border ID
	PreferredLevelBorderID string `json:"PreferredLevelBorderID"`
	Incognito              bool   `json:"Incognito"`
	HideAccountLevel       bool   `json:"HideAccountLevel"`
}

// SeasonalBadgeInfo represents seasonal badge information
type SeasonalBadgeInfo struct {
	// Season ID
	SeasonID        string      `json:"SeasonID"`
	NumberOfWins    int         `json:"NumberOfWins"`
	WinsByTier      interface{} `json:"WinsByTier"` // null
	Rank            int         `json:"Rank"`
	LeaderboardRank int         `json:"LeaderboardRank"`
}

func (v *ValorantAPI) GetPreGameMatch() (*PreGameMatchResponse, error) {
	headers := map[string]string{
		"X-Riot-Entitlements-Jwt": v.Token,
		"X-Riot-ClientPlatform":   v.Client,
		"X-Riot-ClientVersion":    v.Version,
		"Authorization":           "Bearer " + v.AccessToken,
	}
	// https://glz-{region}-1.{shard}.a.pvp.net/pregame/v1/matches/{pre-game match id}
	url := fmt.Sprintf("https://glz-%s-1.%s.a.pvp.net/pregame/v1/matches/%s", v.Region, v.Shard, v.MatchID)
	fmt.Println(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("failed to create request: %w", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("failed to make request: %w", err)
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("failed to get pre game match: %w", err)
		return nil, fmt.Errorf("failed to get pre game match: %w", err)
	}

	fmt.Println(" -------------------------------- RESPONSE --------------------------------")
	fmt.Println(string(body))
	fmt.Println("--------------------------------")

	var preGameMatch PreGameMatchResponse
	if err := json.Unmarshal(body, &preGameMatch); err != nil {
		fmt.Println("failed to unmarshal pre game match: %w", err)
		return nil, fmt.Errorf("failed to unmarshal pre game match: %w", err)
	}

	if preGameMatch.ID == "" {
		return nil, fmt.Errorf("failed to get pre game match: %w", err)
	}

	return &preGameMatch, nil
}
