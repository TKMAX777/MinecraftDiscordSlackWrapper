package main

type PermissionCode uint64

const (
	PermissionAdmin PermissionCode = 1 << iota
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
func GetPermissions(permissionCode PermissionCode) map[string]string {
	var permission = map[string]string{}
	var admin bool = false

	if permissionCode&PermissionAdmin > 0 {
		admin = true
	}

	if permissionCode&PermissionAdvancement > 0 || admin {
		permission["advancement"] = "/advancement"
	}

	if permissionCode&PermissionAttribute > 0 || admin {
		permission["attribute"] = "/attribute"
	}

	if permissionCode&PermissionBan > 0 || admin {
		permission["whitelist"] = "/whitelist"
		permission["ban"] = "/ban"
		permission["banlist"] = "/banlist"
		permission["ban-ip"] = "/ban-ip"
	}

	if permissionCode&PermissionBossbar > 0 || admin {
		permission["bossbar"] = "/bossbar"
	}

	if permissionCode&PermissionClear > 0 || admin {
		permission["clear"] = "/clear"
	}

	if permissionCode&PermissionClone > 0 || admin {
		permission["clone"] = "/clone"
	}

	if permissionCode&PermissionData > 0 || admin {
		permission["data"] = "/data"
	}

	if permissionCode&PermissionDataPack > 0 || admin {
		permission["datapack"] = "/datapack"
	}

	if permissionCode&PermissionDebug > 0 || admin {
		permission["debug"] = "/debug"
	}

	if permissionCode&PermissionDefaultGamemode > 0 || admin {
		permission["defaultgamemode"] = "/defaultgamemode"
	}

	if permissionCode&PermissionDeop > 0 || admin {
		permission["deop"] = "/deop"
	}

	if permissionCode&PermissionDifficulty > 0 || admin {
		permission["difficulty"] = "/difficulty"
	}

	if permissionCode&PermissionEffect > 0 || admin {
		permission["effect"] = "/effect"
	}

	if permissionCode&PermissionEnchant > 0 || admin {
		permission["enchant"] = "/enchant"
	}

	if permissionCode&PermissionExecute > 0 || admin {
		permission["execute"] = "/execute"
	}

	if permissionCode&PermissionExperience > 0 || admin {
		permission["experience"] = "/experience"
		permission["xp"] = "/xp"
	}

	if permissionCode&PermissionFill > 0 || admin {
		permission["fill"] = "/fill"
	}

	if permissionCode&PermissionForceLoad > 0 || admin {
		permission["forceload"] = "/forceload"
	}

	if permissionCode&PermissionFunction > 0 || admin {
		permission["function"] = "/function"
	}

	if permissionCode&PermissionGameMode > 0 || admin {
		permission["gamemode"] = "/gamemode"
	}

	if permissionCode&PermissionGameRule > 0 || admin {
		permission["gamerule"] = "/gamerule"
	}

	if permissionCode&PermissionGive > 0 || admin {
		permission["give"] = "/give"
	}

	if permissionCode&PermissionHelp > 0 || admin {
		permission["help"] = "/help"
	}

	if permissionCode&PermissionKick > 0 || admin {
		permission["kick"] = "/kick"
	}

	if permissionCode&PermissionKill > 0 || admin {
		permission["kill"] = "/kill"
	}

	if permissionCode&PermissionList > 0 || admin {
		permission["list"] = "/list"
	}

	if permissionCode&PermissionLocate > 0 || admin {
		permission["locate"] = "/locate"
		permission["locatebiome"] = "/locatebiome"
	}

	if permissionCode&PermissionLoot > 0 || admin {
		permission["loot"] = "/loot"
	}

	if permissionCode&PermissionMe > 0 || admin {
		permission["me"] = "/me"
	}

	if permissionCode&PermissionMsg > 0 || admin {
		permission["msg"] = "/msg"
		permission["tell"] = "/tell"
		permission["w"] = "/w"
	}

	if permissionCode&PermissionOP > 0 || admin {
		permission["op"] = "/op"
	}

	if permissionCode&PermissionPardon > 0 || admin {
		permission["pardon"] = "/pardon"
	}

	if permissionCode&PermissionParticle > 0 || admin {
		permission["particle"] = "/particle"
	}

	if permissionCode&PermissionPlaySound > 0 || admin {
		permission["playsound"] = "/playsound"
	}

	if permissionCode&PermissionRecipe > 0 || admin {
		permission["recipe"] = "/recipe"
	}

	if permissionCode&PermissionReload > 0 || admin {
		permission["reload"] = "/reload"
	}

	if permissionCode&PermissionReplaceItem > 0 || admin {
		permission["replaceitem"] = "/replaceitem"
	}

	if permissionCode&PermissionSave > 0 || admin {
		permission["save-all"] = "/save-all"
		permission["save-on"] = "/save-on"
		permission["save-off"] = "/save-off"
	}

	if permissionCode&PermissionSay > 0 || admin {
		permission["say"] = "/say"
	}

	if permissionCode&PermissionSchedule > 0 || admin {
		permission["schedule"] = "/schedule"
	}

	if permissionCode&PermissionScoreboard > 0 || admin {
		permission["scoreboard"] = "/scoreboard"
	}

	if permissionCode&PermissionSeed > 0 || admin {
		permission["seed"] = "/seed"
	}

	if permissionCode&PermissionSetBlock > 0 || admin {
		permission["setblock"] = "/setblock"
	}

	if permissionCode&PermissionSetIdleTimeOut > 0 || admin {
		permission["setidletimeout"] = "/setidletimeout"
	}

	if permissionCode&PermissionSetWorldSpawn > 0 || admin {
		permission["setworldspawn"] = "/setworldspawn"
	}

	if permissionCode&PermissionSpawnPoint > 0 || admin {
		permission["spawnpoint"] = "/spawnpoint"
	}

	if permissionCode&PermissionSpectate > 0 || admin {
		permission["spectate"] = "/spectate"
	}

	if permissionCode&PermissionSpreadPlayers > 0 || admin {
		permission["spreadplayers"] = "/spreadplayers"
	}

	if permissionCode&PermissionStop > 0 || admin {
		permission["stop"] = "/stop"
	}

	if permissionCode&PermissionStopSound > 0 || admin {
		permission["stopsound"] = "/stopsound"
	}

	if permissionCode&PermissionSummon > 0 || admin {
		permission["summon"] = "/summon"
	}

	if permissionCode&PermissionTag > 0 || admin {
		permission["tag"] = "/tag"
	}

	if permissionCode&PermissionTeam > 0 || admin {
		permission["team"] = "/team"
	}

	if permissionCode&PermissionTeammsg > 0 || admin {
		permission["teammsg"] = "/teammsg"
		permission["tm"] = "/tm"
	}

	if permissionCode&PermissionTeleport > 0 || admin {
		permission["teleport"] = "/teleport"
		permission["tp"] = "/tp"
	}

	if permissionCode&PermissionTellRow > 0 || admin {
		permission["tellraw"] = "/tellraw"
	}

	if permissionCode&PermissionTime > 0 || admin {
		permission["time"] = "/time"
	}

	if permissionCode&PermissionTitle > 0 || admin {
		permission["title"] = "/title"
	}

	if permissionCode&PermissionTrigger > 0 || admin {
		permission["trigger"] = "/trigger"
	}

	if permissionCode&PermissionWeather > 0 || admin {
		permission["weather"] = "/weather"
	}

	if permissionCode&PermissionWorldBorder > 0 || admin {
		permission["worldborder"] = "/worldborder"
	}

	return permission
}
