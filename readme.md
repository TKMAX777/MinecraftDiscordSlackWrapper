# MinecraftDiscordWrapper
## 概要
DiscordにMinecraftの出力内容を垂れ流すのが目的なプログラム。ついでに許可したユーザに関して、コマンドを用いることが可能。

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
        - [sayコマンドの自動付加](#sayコマンドの自動付加)
        - [標準設定](#標準設定)

<!-- /TOC -->

## 基本設定
### DiscordAPI
- Discord側からの操作を受け付ける場合は`Manage Message` を持ったBotTokenが必要です。
  - メッセージ類をDiscordに流す目的など、サーバ→Discord方向の機能しか用いない場合、これは不要です。

### 設定ファイル

このプログラムをserver.jarと同階層に配置。settings.jsonも次のように記述したものを配置する。

```json
{
    "Discord":{
        "UseDiscord": true/false, (Discord機能を利用するか)
        "UseDiscord2Minecraft": true/false(DiscordからMinecraftのメッセージ送信を含む操作を行うか),
        
        "UserName": "MinecraftWrapper",
        "SendOption": 後述の数値,
        "AddOnlineNumber": true,

	    "GuildID": "用いるDiscordサーバのGuildID",
        "ChannelID": "ChannelID",

	    "Token": "DiscordToken(UseDiscord2Minecraftがfalseの場合不要)",
	
        "InfoOnly": true,(これによって、サーバ出力のthread/INFOのみ取り出される)
        "JoinAndLeftOnly": true,(これによって、join/leftのみ取り出される)
        "AddOnlineNumber": true,(これをすると、Join/Leftの後にオンライン人数が表示される)
	
        "Reaction": {
                "Join": "Joinメッセージの頭に付けるリアクション(省略可)",
                "Left": "Leftメッセージの頭に付けるリアクション(省略可)"
        },
	
        "AvaterURI": "BotのアイコンのURI。お好みで。",
        "UserName":"MinecraftWrapper",
        "HookURI": "https://discord.com/api/webhooks/***",(DiscordTokenにManageWebhookが含まれる場合不要)
    },
    "Slack": {
        "UseSlack": true/false, (Slack機能を利用するか)
        "UseSlack2Minecraft": true/false(SlackからMinecraftのメッセージ送信を含む操作を行うか),

        "AvaterURI": "BotのアイコンのURI。お好みで。", 

        "UserName": "MinecraftWrapper",
        "ChannelID": "ChannelID",

        "SendOption": 後述の数値,

        "SendAllMessages": true, (Slackに投稿されたメッセージをすべてMinecraftに転送する, Say権限が必要)
        "Token": "xoxb-****",
        "EventToken": "xapp-1-***", (UseDiscord2Minecraftがfalseの場合不要)

        "Permissions": 後述のPermissionCode,
        "AddOnlineNumber": true,(これをすると、Join/Leftの後にオンライン人数が表示される)
        
        "Reaction": {
            "Join": ":revolving_hearts:",
            "Left": ":wave:"
        }
    },
    "Minecraft": {
        "JAVA":"JAVA Path 省略可",
        "ThreadInfoRegExp":"Thread/INFOの特定の為の正規表現(標準設定なら不要)",
        "Options": [
            "実行時のオプション",
            "例",
            "-Xms1G",
            "-Xmx4G"
        ]
    }
}
```

### SendOption

- 次の数値の和

|SendOption|数値|
| --- | --- |
|1|すべての出力を転送|
|2|Thread/INFO（入退室を含まない）|
|4|入退室|



## コマンド
### 利用例

Discord / Slackで

```
say Message
```

と入力。

→ Minecraft内に

```
[Server] [USER_Name]Message
```

というコメントが流れる。

### 設定 - Discordの場合 -
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

### 設定 - Slackの場合 -

`settings.json` に該当するPermissionCodeを入力

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
| Stop | 1125899906842624 |
| StopSound | 2251799813685248 |
| Summon | 4503599627370496 |
| Tag | 9007199254740992 |
| Team | 18014398509481984 |
| Teammsg | 36028797018963968 |
| Teleport | 72057594037927936 |
| TellRow | 144115188075855872 |
| Time | 288230376151711744 |
| Title | 576460752303423488 |
| Trigger | 1152921504606846976 |
| Weather | 2305843009213693952 |
| WhiteList | 4611686018427387904 |
| WorldBorder | 9223372036854775808 |

### sayコマンドの自動付加 - Discord -

次のように設定に 
`SendAllMessages`
項目を増やすことで、全てのメッセージを転送することができます。

```json
[{
    "DiscordID": "DiscordID",
    "Name": "USER_NAME",
    "PermissionCode": 0000,
    "SendAllMessages": true
}]
```

ただし、
**PermissionCodeでSayが有効になっている必要があります**
。

### sayコマンドの自動付加 - Slack -

`settings.json` 内に `"SendAllMessages": true` を追加

### 標準設定 - Discord -
`name_dict.json` 
に次オブジェクトを追加することで、全ユーザの設定を一括に行うこともできます。ただし、sayコマンドや、msgコマンドのユーザ名が `Unknown` として扱われるため、通常は以上の設定を記述することをおすすめします。

```
{
	"DiscordID":"Default",
    "PermissionCode": 00000
}
```
