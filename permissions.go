package main

const (
	PermissionAdmin = 1 << iota
	PermissionAdvancement
	PermissionAttribute
	PermissionBan
	PermissionBossbar
	PermissionClear
	PermissionClone
	PermissionData
	PermissionDataPack
	PermissionDebug
	PermissionDefaultGamemode
	PermissionDeop
	PermissionDifficulty
	PermissionEffect
	PermissionEnchant
	PermissionExecute
	PermissionExperience
	PermissionFill
	PermissionForceLoad
	PermissionFunction
	PermissionGameMode
	PermissionGameRule
	PermissionGive
	PermissionHelp
	PermissionKick
	PermissionKill
	PermissionList
	PermissionLocate
	PermissionLoot
	PermissionMe
	PermissionMsg
	PermissionOP
	PermissionPardon
	PermissionPardonIP
	PermissionParticle
	PermissionPlaySound
	PermissionRecipe
	PermissionReload
	PermissionReplaceItem
	PermissionSave
	PermissionSay
	PermissionSchedule
	PermissionScoreboard
	PermissionSeed
	PermissionSetBlock
	PermissionSetIdleTimeOut
	PermissionSetWorldSpawn
	PermissionSpawnPoint
	PermissionSpectate
	PermissionSpreadPlayers
	PermissionStop
	PermissionStopSound
	PermissionSummon
	PermissionTag
	PermissionTeam
	PermissionTeammsg
	PermissionTeleport
	PermissionTellRow
	PermissionTime
	PermissionTitle
	PermissionTrigger
	PermissionWeather
	PermissionWhiteList
	PermissionWorldBorder
)

// GetPermissions return permissions map
func GetPermissions(PermissionCode uint64) map[string]string {
	var permission = map[string]string{}
	var admin bool = false

	if PermissionCode&PermissionAdmin > 0 {
		admin = true
	}

	if PermissionCode&PermissionAdvancement > 0 || admin {
		permission["advancement"] = "/advancement"
	}

	if PermissionCode&PermissionAttribute > 0 || admin {
		permission["attribute"] = "/attribute"
	}

	if PermissionCode&PermissionBan > 0 || admin {
		permission["whitelist"] = "/whitelist"
		permission["ban"] = "/ban"
		permission["banlist"] = "/banlist"
		permission["ban-ip"] = "/ban-ip"
	}

	if PermissionCode&PermissionBossbar > 0 || admin {
		permission["bossbar"] = "/bossbar"
	}

	if PermissionCode&PermissionClear > 0 || admin {
		permission["clear"] = "/clear"
	}

	if PermissionCode&PermissionClone > 0 || admin {
		permission["clone"] = "/clone"
	}

	if PermissionCode&PermissionData > 0 || admin {
		permission["data"] = "/data"
	}

	if PermissionCode&PermissionDataPack > 0 || admin {
		permission["datapack"] = "/datapack"
	}

	if PermissionCode&PermissionDebug > 0 || admin {
		permission["debug"] = "/debug"
	}

	if PermissionCode&PermissionDefaultGamemode > 0 || admin {
		permission["defaultgamemode"] = "/defaultgamemode"
	}

	if PermissionCode&PermissionDeop > 0 || admin {
		permission["deop"] = "/deop"
	}

	if PermissionCode&PermissionDifficulty > 0 || admin {
		permission["difficulty"] = "/difficulty"
	}

	if PermissionCode&PermissionEffect > 0 || admin {
		permission["effect"] = "/effect"
	}

	if PermissionCode&PermissionEnchant > 0 || admin {
		permission["enchant"] = "/enchant"
	}

	if PermissionCode&PermissionExecute > 0 || admin {
		permission["execute"] = "/execute"
	}

	if PermissionCode&PermissionExperience > 0 || admin {
		permission["experience"] = "/experience"
		permission["xp"] = "/xp"
	}

	if PermissionCode&PermissionFill > 0 || admin {
		permission["fill"] = "/fill"
	}

	if PermissionCode&PermissionForceLoad > 0 || admin {
		permission["forceload"] = "/forceload"
	}

	if PermissionCode&PermissionFunction > 0 || admin {
		permission["function"] = "/function"
	}

	if PermissionCode&PermissionGameMode > 0 || admin {
		permission["gamemode"] = "/gamemode"
	}

	if PermissionCode&PermissionGameRule > 0 || admin {
		permission["gamerule"] = "/gamerule"
	}

	if PermissionCode&PermissionGive > 0 || admin {
		permission["give"] = "/give"
	}

	if PermissionCode&PermissionHelp > 0 || admin {
		permission["help"] = "/help"
	}

	if PermissionCode&PermissionKick > 0 || admin {
		permission["kick"] = "/kick"
	}

	if PermissionCode&PermissionKill > 0 || admin {
		permission["kill"] = "/kill"
	}

	if PermissionCode&PermissionList > 0 || admin {
		permission["list"] = "/list"
	}

	if PermissionCode&PermissionLocate > 0 || admin {
		permission["locate"] = "/locate"
		permission["locatebiome"] = "/locatebiome"
	}

	if PermissionCode&PermissionLoot > 0 || admin {
		permission["loot"] = "/loot"
	}

	if PermissionCode&PermissionMe > 0 || admin {
		permission["me"] = "/me"
	}

	if PermissionCode&PermissionMsg > 0 || admin {
		permission["msg"] = "/msg"
		permission["tell"] = "/tell"
		permission["w"] = "/w"
	}

	if PermissionCode&PermissionOP > 0 || admin {
		permission["op"] = "/op"
	}

	if PermissionCode&PermissionPardon > 0 || admin {
		permission["pardon"] = "/pardon"
	}

	if PermissionCode&PermissionParticle > 0 || admin {
		permission["particle"] = "/particle"
	}

	if PermissionCode&PermissionPlaySound > 0 || admin {
		permission["playsound"] = "/playsound"
	}

	if PermissionCode&PermissionRecipe > 0 || admin {
		permission["recipe"] = "/recipe"
	}

	if PermissionCode&PermissionReload > 0 || admin {
		permission["reload"] = "/reload"
	}

	if PermissionCode&PermissionReplaceItem > 0 || admin {
		permission["replaceitem"] = "/replaceitem"
	}

	if PermissionCode&PermissionSave > 0 || admin {
		permission["save-all"] = "/save-all"
		permission["save-on"] = "/save-on"
		permission["save-off"] = "/save-off"
	}

	if PermissionCode&PermissionSay > 0 || admin {
		permission["say"] = "/say"
	}

	if PermissionCode&PermissionSchedule > 0 || admin {
		permission["schedule"] = "/schedule"
	}

	if PermissionCode&PermissionScoreboard > 0 || admin {
		permission["scoreboard"] = "/scoreboard"
	}

	if PermissionCode&PermissionSeed > 0 || admin {
		permission["seed"] = "/seed"
	}

	if PermissionCode&PermissionSetBlock > 0 || admin {
		permission["setblock"] = "/setblock"
	}

	if PermissionCode&PermissionSetIdleTimeOut > 0 || admin {
		permission["setidletimeout"] = "/setidletimeout"
	}

	if PermissionCode&PermissionSetWorldSpawn > 0 || admin {
		permission["setworldspawn"] = "/setworldspawn"
	}

	if PermissionCode&PermissionSpawnPoint > 0 || admin {
		permission["spawnpoint"] = "/spawnpoint"
	}

	if PermissionCode&PermissionSpectate > 0 || admin {
		permission["spectate"] = "/spectate"
	}

	if PermissionCode&PermissionSpreadPlayers > 0 || admin {
		permission["spreadplayers"] = "/spreadplayers"
	}

	if PermissionCode&PermissionStop > 0 || admin {
		permission["stop"] = "/stop"
	}

	if PermissionCode&PermissionStopSound > 0 || admin {
		permission["stopsound"] = "/stopsound"
	}

	if PermissionCode&PermissionSummon > 0 || admin {
		permission["summon"] = "/summon"
	}

	if PermissionCode&PermissionTag > 0 || admin {
		permission["tag"] = "/tag"
	}

	if PermissionCode&PermissionTeam > 0 || admin {
		permission["team"] = "/team"
	}

	if PermissionCode&PermissionTeammsg > 0 || admin {
		permission["teammsg"] = "/teammsg"
		permission["tm"] = "/tm"
	}

	if PermissionCode&PermissionTeleport > 0 || admin {
		permission["teleport"] = "/teleport"
		permission["tp"] = "/tp"
	}

	if PermissionCode&PermissionTellRow > 0 || admin {
		permission["tellraw"] = "/tellraw"
	}

	if PermissionCode&PermissionTime > 0 || admin {
		permission["time"] = "/time"
	}

	if PermissionCode&PermissionTitle > 0 || admin {
		permission["title"] = "/title"
	}

	if PermissionCode&PermissionTrigger > 0 || admin {
		permission["trigger"] = "/trigger"
	}

	if PermissionCode&PermissionWeather > 0 || admin {
		permission["weather"] = "/weather"
	}

	if PermissionCode&PermissionWorldBorder > 0 || admin {
		permission["worldborder"] = "/worldborder"
	}

	return permission
}
