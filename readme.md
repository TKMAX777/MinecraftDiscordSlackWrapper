# MinecraftDiscordWrapper
## 概要
DiscordにMinecraftの出力内容を垂れ流すのが目的なプログラム。ついでに許可したユーザに関してはsayを用いることが可能。

## 目次
<!-- TOC -->

- [MinecraftDiscordWrapper](#minecraftdiscordwrapper)
    - [概要](#概要)
    - [目次](#目次)
    - [基本設定](#基本設定)
        - [DiscordAPI](#discordapi)
        - [設定ファイル](#設定ファイル)
    - [コマンド](#コマンド)
        - [利用例](#利用例)
        - [設定](#設定)
        - [PermissionCode](#permissioncode)
        - [標準設定](#標準設定)
    - [免責事項](#免責事項)

<!-- /TOC -->

## 基本設定
### DiscordAPI
- `Manage Message` を持ったBotTokenが必要です。

### 設定ファイル

このプログラムをserver.jarと同階層に配置。settings.jsonも次のように記述したものを配置する。

```json
{
    "Discord":{
        "GuildID": "用いるDiscordサーバのGuildID",
        "Token": "DiscordToken",
        "Default":{
            "AvaterURI": "BotのアイコンのURI。お好みで。",
            "UserName":"MinecraftWrapper",
            "HookURI": "https://discord.com/api/webhooks/***"
        },
        "Error":{
            "AvaterURI": "Error時のBotのアイコンのURI。お好みで。",
            "UserName":"MinecraftWrapper - Error -",
            "HookURI": "https://discord.com/api/webhooks/***(省略可。その場合Defaultが適用される)"
        }        
    },
    "Minecraft": {
        "Options": [
            "実行時のオプション",
            "例",
            "-Xms1G",
            "-Xmx4G"
        ]
    }
}
```


## コマンド
### 利用例

Discordで

```
say Message
```

と入力。

→ Minecraft内に

```
[Server] [USER_Name]Message
```

というコメントが流れる。

### 設定
次のような
`name_dict.json`
ファイルを同層に作成すると、Discordで許可したコマンドを叩くことが可能

```json
[
    {
	"DiscordID":"DiscordUserID",
	"Name": "USER_NAME",
    "PermissionCode": 00000,
    }
]
```

### PermissionCode
各々のPermissionCodeの和を入力。Adminは全ての権限を許可します。

| Permission | Code |
| --- | --- |
| Admin | 1 |
| Advancement | 2 |
| Attribute | 4 |
| Ban | 8 |
| Bossbar | 16 |
| Clear | 32 |
| Clone | 64 |
| Data | 128 |
| DataPack | 256 |
| Debug | 512 |
| DefaultGamemode | 1024 |
| Deop | 2048 |
| Difficulty | 4096 |
| Effect | 8192 |
| Enchant | 16384 |
| Execute | 32768 |
| Experience | 65536 |
| Fill | 131072 |
| ForceLoad | 262144 |
| Function | 524288 |
| GameMode | 1048576 |
| GameRule | 2097152 |
| Give | 4194304 |
| Help | 8388608 |
| Kick | 16777216 |
| Kill | 33554432 |
| List | 67108864 |
| Locate | 134217728 |
| Loot | 268435456 |
| Me | 536870912 |
| Msg | 1073741824 |
| OP | 2147483648 |
| Pardon | 4294967296 |
| PardonIP | 8589934592 |
| Particle | 17179869184 |
| PlaySound | 34359738368 |
| Recipe | 68719476736 |
| Reload | 137438953472 |
| ReplaceItem | 274877906944 |
| Save | 549755813888 |
| Say | 1099511627776 |
| Schedule | 2199023255552 |
| Scoreboard | 4398046511104 |
| Seed | 8796093022208 |
| SetBlock | 17592186044416 |
| SetIdleTimeOut | 35184372088832 |
| SetWorldSpawn | 70368744177664 |
| SpawnPoint | 140737488355328 |
| Spectate | 281474976710656 |
| SpreadPlayers | 562949953421312 |
| Stop | 1125899906842620 |
| StopSound | 2251799813685250 |
| Summon | 4503599627370500 |
| Tag | 9007199254740990 |
| Team | 18014398509482000 |
| Teammsg | 36028797018964000 |
| Teleport | 72057594037927900 |
| TellRow | 144115188075856000 |
| Time | 288230376151712000 |
| Title | 576460752303423000 |
| Trigger | 1152921504606850000 |
| Weather | 2305843009213690000 |
| WhiteList | 4611686018427390000 |
| WorldBorder | 9223372036854780000 |

### 標準設定
`name_dict.json` 
に次オブジェクトを追加することで、全ユーザの設定を一括に行うこともできます。ただし、sayコマンドや、msgコマンドのユーザ名が `Unknown` として扱われるため、通常は以上の設定を記述することをおすすめします。

```
{
	"DiscordID":"Default",
    "PermissionCode": 00000,
}
```

## 免責事項
このプログラムはMinecraftサーバプログラムの標準入出力を完全にラップしたものとなっております。

悪意のある第三者や、その他このプログラムに存在する未知のバグにより、サーバ運営に於いて深刻なセキュリティリスクを負う可能性があります。それらにより生じた損害について、製作者は一切の責任は持ちません。あくまでご利用は自己責任でお願いします。

