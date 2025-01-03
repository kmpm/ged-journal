package messages

type StatisticsEvent struct {
	Event
	BankAccount struct {
		CurrentWealth          int64 `json:"Current_Wealth"`
		SpentOnShips           int64 `json:"Spent_On_Ships"`
		SpentOnOutfitting      int64 `json:"Spent_On_Outfitting"`
		SpentOnRepairs         int64 `json:"Spent_On_Repairs"`
		SpentOnFuel            int64 `json:"Spent_On_Fuel"`
		SpentOnAmmoConsumables int64 `json:"Spent_On_Ammo_Consumables"`
		InsuranceClaims        int64 `json:"Insurance_Claims"`
		SpentOnInsurance       int64 `json:"Spent_On_Insurance"`
		OwnedShipCount         int64 `json:"Owned_Ship_Count"`
		SpentOnSuits           int64 `json:"Spent_On_Suits"`
		SpentOnWeapons         int64 `json:"Spent_On_Weapons"`
		SpentOnSuitConsumables int64 `json:"Spent_On_Suit_Consumables"`
		SuitsOwned             int64 `json:"Suits_Owned"`
		WeaponsOwned           int64 `json:"Weapons_Owned"`
		SpentOnPremiumStock    int64 `json:"Spent_On_Premium_Stock"`
		PremiumStockBought     int64 `json:"Premium_Stock_Bought"`
	} `json:"Bank_Account"`
	Combat struct {
		Assasination_Profits     int64 `json:"Assasination_Profits"`
		Assassinations           int64 `json:"Assassinations"`
		BountiesClaimed          int64 `json:"Bounties_Claimed"`
		BountyHuntingProfit      int64 `json:"Bounty_Hunting_Profit"`
		CombatBonds              int64 `json:"Combat_Bonds"`
		ConflictZoneHigh         int64 `json:"ConflictZone_High"`
		ConflictZoneLow          int64 `json:"ConflictZone_Low"`
		ConflictZoneMedium       int64 `json:"ConflictZone_Medium"`
		ConflictZoneTotal        int64 `json:"ConflictZone_Total"`
		ConflictZoneHighWins     int64 `json:"ConflictZone_High_Wins"`
		ConflictZoneLowWins      int64 `json:"ConflictZone_Low_Wins"`
		ConflictZoneMediumWins   int64 `json:"ConflictZone_Medium_Wins"`
		ConflictZoneTotalWins    int64 `json:"ConflictZone_Total_Wins"`
		CombatBondProfits        int64 `json:"Combat_Bond_Profits"`
		DropshipsBooked          int64 `json:"Dropships_Booked"`
		DropshipsCancelled       int64 `json:"Dropships_Cancelled"`
		DropshipsTaken           int64 `json:"Dropships_Taken"`
		HighestSingleReward      int64 `json:"Highest_Single_Reward"`
		OnFootCombatBonds        int64 `json:"OnFoot_Combat_Bonds"`
		OnFootCombatBondsProfits int64 `json:"OnFoot_Combat_Bonds_Profits"`
		OnFootShipsDestroyed     int64 `json:"OnFoot_Ships_Destroyed"`
		OnFootSkimmersKilled     int64 `json:"OnFoot_Skimmers_Killed"`
		OnFootScavsKilled        int64 `json:"OnFoot_Scavs_Killed"`
		SettlementConquered      int64 `json:"Settlement_Conquered"`
		SettlementDefended       int64 `json:"Settlement_Defended"`
	} `json:"Combat"`

	Crime struct {
		Fines                    int64 `json:"Fines"`
		Notoriety                int64 `json:"Notoriety"`
		TotalFines               int64 `json:"Total_Fines"`
		BountiesReceived         int64 `json:"Bounties_Received"`
		TotalBounties            int64 `json:"Total_Bounties"`
		HighestBounty            int64 `json:"Highest_Bounty"`
		MalwareUploaded          int64 `json:"Malware_Uploaded"`
		SettlementsStateShutdown int64 `json:"Settlements_State_Shutdown"`
		ProductionSabotage       int64 `json:"Production_Sabotage"`
		ProductionTheft          int64 `json:"Production_Theft"`
		TotalMurders             int64 `json:"Total_Murders"`
		CitizensMurdered         int64 `json:"Citizens_Murdered"`
		OmnipolMurdered          int64 `json:"Omnipol_Murdered"`
		GuardsMurdered           int64 `json:"Guards_Murdered"`
		DataStolen               int64 `json:"Data_Stolen"`
		GoodsStolen              int64 `json:"Goods_Stolen"`
		SampleStolen             int64 `json:"Sample_Stolen"`
		TotalStolen              int64 `json:"Total_Stolen"`
		TurretsDestroyed         int64 `json:"Turrets_Destroyed"`
		TurretsOverloaded        int64 `json:"Turrets_Overloaded"`
		TurretsTotal             int64 `json:"Turrets_Total"`
		ValueStolenStateChange   int64 `json:"Value_Stolen_State_Change"`
		ProfilesCloned           int64 `json:"Profiles_Cloned"`
	} `json:"Crime"`
	Smuggling struct {
		BlackMarketsTradedWith   int64   `json:"Black_Markets_Traded_With"`
		BlackMarketsProfits      int64   `json:"Black_Markets_Profits"`
		ResourcesSmuggled        int64   `json:"Resources_Smuggled"`
		AverageProfit            float64 `json:"Average_Profit"`
		HighestSingleTransaction int64   `json:"Highest_Single_Transaction"`
	} `json:"Smuggling"`
	Trading struct {
		MarketsTradedWith        int64   `json:"Markets_Traded_With"`
		MarketProfits            int64   `json:"Market_Profits"`
		ResourcesTraded          int64   `json:"Resources_Traded"`
		AverageProfit            float64 `json:"Average_Profit"`
		HighestSingleTransaction int64   `json:"Highest_Single_Transaction"`
		DataSold                 int64   `json:"Data_Sold"`
		GoodsSold                int64   `json:"Goods_Sold"`
		AssetsSold               int64   `json:"Assets_Sold"`
	} `json:"Trading"`
	Mining struct {
		MiningProfits     int64 `json:"Mining_Profits"`
		QuantityMined     int64 `json:"Quantity_Mined"`
		MaterialsCollcted int64 `json:"Materials_Collected"`
	} `json:"Mining"`
	Exploration struct {
		SystemsVisited            int64   `json:"Systems_Visited"`
		ExplorationProfits        int64   `json:"Exploration_Profits"`
		PlanetsScannedToLevel2    int64   `json:"Planets_Scanned_To_Level_2"`
		PlanetsScannedToLevel3    int64   `json:"Planets_Scanned_To_Level_3"`
		EfficientScans            int64   `json:"Efficient_Scans"`
		HighestPayout             int64   `json:"Highest_Payout"`
		TotalHyperspaceDistance   int64   `json:"Total_Hyperspace_Distance"`
		TotalHyperspaceJumps      int64   `json:"Total_Hyperspace_Jumps"`
		GreatestDistanceFromStart float64 `json:"Greatest_Distance_From_Start"`
		TimePlayed                int64   `json:"Time_Played"`
		OnFootDistanceTravelled   int64   `json:"OnFoot_Distance_Travelled"`
		ShuttleJourneys           int64   `json:"Shuttle_Journeys"`
		ShuttleDistanceTravelled  int64   `json:"Shuttle_Distance_Travelled"`
		SpentOnShuttles           int64   `json:"Spent_On_Shuttles"`
		FirstFootfalls            int64   `json:"First_Footfalls"`
		PlanetFootfalls           int64   `json:"Planet_Footfalls"`
		SettlementsVisited        int64   `json:"Settlements_Visited"`
	} `json:"Exploration"`
	Passengers struct {
		PassengersMissionsAccepted  int64 `json:"Passengers_Missions_Accepted"`
		PassengersMissionsBulk      int64 `json:"Passengers_Missions_Bulk"`
		PassengersMissionsDelivered int64 `json:"Passengers_Missions_Delivered"`
		PassengersMissionsVIP       int64 `json:"Passengers_Missions_VIP"`
		PassengersMissionsEjected   int64 `json:"Passengers_Missions_Ejected"`
	} `json:"Passengers"`
	SearchAndRescue struct {
		SearchRescueTraded        int64 `json:"SearchRescue_Traded"`
		SearchRescueProfit        int64 `json:"SearchRescue_Profit"`
		SearchRescueCount         int64 `json:"SearchRescue_Count"`
		SalvageLegalPOI           int64 `json:"Salvage_Legal_POI"`
		SalvageLegalSettlements   int64 `json:"Salvage_Legal_Settlements"`
		SalvageIllegalPOI         int64 `json:"Salvage_Illegal_POI"`
		SalvageIllegalSettlements int64 `json:"Salvage_Illegal_Settlements"`
		MaglocksOpened            int64 `json:"Maglocks_Opened"`
		PanelsOpened              int64 `json:"Panels_Opened"`
		SettlementsStateFireOut   int64 `json:"Settlements_State_FireOut"`
		SettlementsStateReboot    int64 `json:"Settlements_State_Reboot"`
	} `json:"Search_And_Rescue"`
	Crafting struct {
		CountOfUsedEngineers  int64 `json:"Count_Of_Used_Engineers"`
		RecipesGenerated      int64 `json:"Recipes_Generated"`
		RecipesGeneratedRank1 int64 `json:"Recipes_Generated_Rank_1"`
		RecipesGeneratedRank2 int64 `json:"Recipes_Generated_Rank_2"`
		RecipesGeneratedRank3 int64 `json:"Recipes_Generated_Rank_3"`
		RecipesGeneratedRank4 int64 `json:"Recipes_Generated_Rank_4"`
		RecipesGeneratedRank5 int64 `json:"Recipes_Generated_Rank_5"`
		SuitModsApplied       int64 `json:"Suit_Mods_Applied"`
		WeaponModsApplied     int64 `json:"Weapon_Mods_Applied"`
		SuitsUpgraded         int64 `json:"Suits_Upgraded"`
		WeaponsUpgraded       int64 `json:"Weapons_Upgraded"`
		SuitsUpgradedFull     int64 `json:"Suits_Upgraded_Full"`
		WeaponsUpgradedFull   int64 `json:"Weapons_Upgraded_Full"`
		SuitModsAppliedFull   int64 `json:"Suit_Mods_Applied_Full"`
		WeaponModsAppliedFull int64 `json:"Weapon_Mods_Applied_Full"`
	} `json:"Crafting"`
	Crew struct {
		NpcCrewHired      int64 `json:"NpcCrew_Hired"`
		NpcCrewFired      int64 `json:"NpcCrew_Fired"`
		NpcCrewDied       int64 `json:"NpcCrew_Died"`
		NpcCrewTotalWages int64 `json:"NpcCrew_Total_Wages"`
	} `json:"Crew"`
	Multicrew struct {
		MulticrewTimeTotal        int64 `json:"Multicrew_Time_Total"`
		MulticrewGunnerTimeTotal  int64 `json:"Multicrew_Gunner_Time_Total"`
		MulticrewFighterTimeTotal int64 `json:"Multicrew_Fighter_Time_Total"`
		MulticrewCreditsTotal     int64 `json:"Multicrew_Credits_Total"`
		MulticrewFinesTotal       int64 `json:"Multicrew_Fines_Total"`
	} `json:"Multicrew"`
	MaterialTraderStats struct {
		TradesCompleted int64 `json:"Trades_Completed"`
		MaterialsTraded int64 `json:"Materials_Traded"`
		AssetsTradedIn  int64 `json:"Assets_Traded_In"`
		AssetsTradedOut int64 `json:"Assets_Traded_Out"`
	} `json:"Material_Trader_Stats"`
	Exobiology struct {
		OrganicGenusEncountered   int64 `json:"Organic_Genus_Encountered"`
		OrganicSpeciesEncountered int64 `json:"Organic_Species_Encountered"`
		OrganicVariantEncountered int64 `json:"Organic_Variant_Encountered"`
		OrganicDataProfits        int64 `json:"Organic_Data_Profits"`
		OrganicData               int64 `json:"Organic_Data"`
		FirstLoggedProfits        int64 `json:"First_Logged_Profits"`
		FirstLogged               int64 `json:"First_Logged"`
		OrganicSystems            int64 `json:"Organic_Systems"`
		OrganicPlanets            int64 `json:"Organic_Planets"`
		OrganicGenus              int64 `json:"Organic_Genus"`
		OrganicSpecies            int64 `json:"Organic_Species"`
	} `json:"Exobiology"`
}
