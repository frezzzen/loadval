export namespace main {
	
	export class GunAttachment {
	    ID: string;
	
	    static createFrom(source: any = {}) {
	        return new GunAttachment(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	    }
	}
	export class Gun {
	    ID: string;
	    CharmInstanceID?: string;
	    CharmID?: string;
	    CharmLevelID?: string;
	    SkinID?: string;
	    SkinLevelID?: string;
	    ChromaID?: string;
	    Attachments: GunAttachment[];
	
	    static createFrom(source: any = {}) {
	        return new Gun(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.CharmInstanceID = source["CharmInstanceID"];
	        this.CharmID = source["CharmID"];
	        this.CharmLevelID = source["CharmLevelID"];
	        this.SkinID = source["SkinID"];
	        this.SkinLevelID = source["SkinLevelID"];
	        this.ChromaID = source["ChromaID"];
	        this.Attachments = this.convertValues(source["Attachments"], GunAttachment);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class Identity {
	    PlayerCardID: string;
	    PlayerTitleID: string;
	    AccountLevel: number;
	    PreferredLevelBorderID: string;
	    HideAccountLevel: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Identity(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.PlayerCardID = source["PlayerCardID"];
	        this.PlayerTitleID = source["PlayerTitleID"];
	        this.AccountLevel = source["AccountLevel"];
	        this.PreferredLevelBorderID = source["PreferredLevelBorderID"];
	        this.HideAccountLevel = source["HideAccountLevel"];
	    }
	}
	export class Spray {
	    SprayID: string;
	    SprayLevelID?: string;
	    EquipSlotID: string;
	
	    static createFrom(source: any = {}) {
	        return new Spray(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.SprayID = source["SprayID"];
	        this.SprayLevelID = source["SprayLevelID"];
	        this.EquipSlotID = source["EquipSlotID"];
	    }
	}
	export class PlayerLoadoutResponse {
	    Guns: Gun[];
	    Sprays: Spray[];
	    Identity: Identity;
	    Incognito: boolean;
	
	    static createFrom(source: any = {}) {
	        return new PlayerLoadoutResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Guns = this.convertValues(source["Guns"], Gun);
	        this.Sprays = this.convertValues(source["Sprays"], Spray);
	        this.Identity = this.convertValues(source["Identity"], Identity);
	        this.Incognito = source["Incognito"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class OwnedItemsResponseEntitlement {
	    TypeID: string;
	    ItemID: string;
	    InstanceID: string;
	
	    static createFrom(source: any = {}) {
	        return new OwnedItemsResponseEntitlement(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.TypeID = source["TypeID"];
	        this.ItemID = source["ItemID"];
	        this.InstanceID = source["InstanceID"];
	    }
	}
	export class OwnedItemsResponse {
	    ItemTypeID: string;
	    Entitlements: OwnedItemsResponseEntitlement[];
	
	    static createFrom(source: any = {}) {
	        return new OwnedItemsResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ItemTypeID = source["ItemTypeID"];
	        this.Entitlements = this.convertValues(source["Entitlements"], OwnedItemsResponseEntitlement);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class MainData {
	    OwnedSkins: OwnedItemsResponse[];
	    OwnedSkinVariants: OwnedItemsResponse[];
	    OwnedAgents: OwnedItemsResponse[];
	    OwnedCards: OwnedItemsResponse[];
	    PlayerLoadout?: PlayerLoadoutResponse;
	
	    static createFrom(source: any = {}) {
	        return new MainData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.OwnedSkins = this.convertValues(source["OwnedSkins"], OwnedItemsResponse);
	        this.OwnedSkinVariants = this.convertValues(source["OwnedSkinVariants"], OwnedItemsResponse);
	        this.OwnedAgents = this.convertValues(source["OwnedAgents"], OwnedItemsResponse);
	        this.OwnedCards = this.convertValues(source["OwnedCards"], OwnedItemsResponse);
	        this.PlayerLoadout = this.convertValues(source["PlayerLoadout"], PlayerLoadoutResponse);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	export class SeasonalBadgeInfo {
	    SeasonID: string;
	    NumberOfWins: number;
	    WinsByTier: any;
	    Rank: number;
	    LeaderboardRank: number;
	
	    static createFrom(source: any = {}) {
	        return new SeasonalBadgeInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.SeasonID = source["SeasonID"];
	        this.NumberOfWins = source["NumberOfWins"];
	        this.WinsByTier = source["WinsByTier"];
	        this.Rank = source["Rank"];
	        this.LeaderboardRank = source["LeaderboardRank"];
	    }
	}
	export class PlayerIdentity {
	    Subject: string;
	    PlayerCardID: string;
	    PlayerTitleID: string;
	    AccountLevel: number;
	    PreferredLevelBorderID: string;
	    Incognito: boolean;
	    HideAccountLevel: boolean;
	
	    static createFrom(source: any = {}) {
	        return new PlayerIdentity(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Subject = source["Subject"];
	        this.PlayerCardID = source["PlayerCardID"];
	        this.PlayerTitleID = source["PlayerTitleID"];
	        this.AccountLevel = source["AccountLevel"];
	        this.PreferredLevelBorderID = source["PreferredLevelBorderID"];
	        this.Incognito = source["Incognito"];
	        this.HideAccountLevel = source["HideAccountLevel"];
	    }
	}
	export class Player {
	    Subject: string;
	    CharacterID: string;
	    CharacterSelectionState: string;
	    PregamePlayerState: string;
	    CompetitiveTier: number;
	    PlayerIdentity: PlayerIdentity;
	    SeasonalBadgeInfo: SeasonalBadgeInfo;
	    IsCaptain: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Player(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Subject = source["Subject"];
	        this.CharacterID = source["CharacterID"];
	        this.CharacterSelectionState = source["CharacterSelectionState"];
	        this.PregamePlayerState = source["PregamePlayerState"];
	        this.CompetitiveTier = source["CompetitiveTier"];
	        this.PlayerIdentity = this.convertValues(source["PlayerIdentity"], PlayerIdentity);
	        this.SeasonalBadgeInfo = this.convertValues(source["SeasonalBadgeInfo"], SeasonalBadgeInfo);
	        this.IsCaptain = source["IsCaptain"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	export class Team {
	    TeamID: string;
	    Players: Player[];
	
	    static createFrom(source: any = {}) {
	        return new Team(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.TeamID = source["TeamID"];
	        this.Players = this.convertValues(source["Players"], Player);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PreGameMatchResponse {
	    ID: string;
	    Version: number;
	    Teams: Team[];
	    AllyTeam?: Team;
	    EnemyTeam?: Team;
	    ObserverSubjects: any[];
	    MatchCoaches: any[];
	    EnemyTeamSize: number;
	    EnemyTeamLockCount: number;
	    PregameState: string;
	    LastUpdated: string;
	    MapID: string;
	    MapSelectPool: any[];
	    BannedMapIDs: any[];
	    CastedVotes?: any;
	    MapSelectSteps: any[];
	    MapSelectStep: number;
	    Team1: string;
	    GamePodID: string;
	    Mode: string;
	    VoiceSessionID: string;
	    MUCName: string;
	    TeamMatchToken: string;
	    QueueID: string;
	    ProvisioningFlowID: string;
	    IsRanked: boolean;
	    PhaseTimeRemainingNS: number;
	    StepTimeRemainingNS: number;
	    altModesFlagADA: boolean;
	    TournamentMetadata: any;
	    RosterMetadata: any;
	
	    static createFrom(source: any = {}) {
	        return new PreGameMatchResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.Version = source["Version"];
	        this.Teams = this.convertValues(source["Teams"], Team);
	        this.AllyTeam = this.convertValues(source["AllyTeam"], Team);
	        this.EnemyTeam = this.convertValues(source["EnemyTeam"], Team);
	        this.ObserverSubjects = source["ObserverSubjects"];
	        this.MatchCoaches = source["MatchCoaches"];
	        this.EnemyTeamSize = source["EnemyTeamSize"];
	        this.EnemyTeamLockCount = source["EnemyTeamLockCount"];
	        this.PregameState = source["PregameState"];
	        this.LastUpdated = source["LastUpdated"];
	        this.MapID = source["MapID"];
	        this.MapSelectPool = source["MapSelectPool"];
	        this.BannedMapIDs = source["BannedMapIDs"];
	        this.CastedVotes = source["CastedVotes"];
	        this.MapSelectSteps = source["MapSelectSteps"];
	        this.MapSelectStep = source["MapSelectStep"];
	        this.Team1 = source["Team1"];
	        this.GamePodID = source["GamePodID"];
	        this.Mode = source["Mode"];
	        this.VoiceSessionID = source["VoiceSessionID"];
	        this.MUCName = source["MUCName"];
	        this.TeamMatchToken = source["TeamMatchToken"];
	        this.QueueID = source["QueueID"];
	        this.ProvisioningFlowID = source["ProvisioningFlowID"];
	        this.IsRanked = source["IsRanked"];
	        this.PhaseTimeRemainingNS = source["PhaseTimeRemainingNS"];
	        this.StepTimeRemainingNS = source["StepTimeRemainingNS"];
	        this.altModesFlagADA = source["altModesFlagADA"];
	        this.TournamentMetadata = source["TournamentMetadata"];
	        this.RosterMetadata = source["RosterMetadata"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PreGamePlayerResponse {
	    Subject: string;
	    MatchID: string;
	    Version: number;
	
	    static createFrom(source: any = {}) {
	        return new PreGamePlayerResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Subject = source["Subject"];
	        this.MatchID = source["MatchID"];
	        this.Version = source["Version"];
	    }
	}
	
	

}

