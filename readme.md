# MinecraftDiscordSlackWrapper

## 概要

DiscordやSlackにMinecraftの出力内容を垂れ流すプログラム。ついでに許可したユーザに関して、コマンドを用いることが可能。

## 目次

<!-- TOC -->

- [MinecraftDiscordSlackWrapper](#minecraftdiscordslackwrapper)
    - [概要](#概要)
    - [目次](#目次)
    - [基本設定](#基本設定)
        - [DiscordAPI](#discordapi)
        - [SlackAPI](#slackapi)
        - [設定ファイル](#設定ファイル)
    - [参考 - SendJoinStateMessages -](#参考---sendjoinstatemessages--)

<!-- /TOC -->

## 基本設定

### DiscordAPI

- Discord側からの操作を受け付ける場合は `Manage Message` を持ったBotTokenが必要です。
  - メッセージ類をDiscordに流す目的など、サーバ→Discord方向の機能しか用いない場合、これは不要です。

### SlackAPI

次の権限を持ったBotトークンと、AppLevelトークンが必要です。

- API Permissions
[channels:history](https://api.slack.com/scopes/channels:history)
[chat:write](https://api.slack.com/scopes/chat:write)
[chat:write.customize](https://api.slack.com/scopes/chat:write.customize)
[groups:history](https://api.slack.com/scopes/groups:history)
[groups:write](https://api.slack.com/scopes/groups:write)
[users:read](https://api.slack.com/scopes/users:read)


- EventAPI Permissions
[message.channels](https://api.slack.com/events/message.channels) 
[message.groups](https://api.slack.com/events/message.groups)    

### 設定ファイル

- `settings`以下のファイルを編集することで、設定が変更できます。

## 参考 - SendJoinStateMessages -

Discordでは `SendJoinStateMessage: true` を指定すると、顔アイコン・名前・ログイン人数・今回のログインからの経過時間を画像で表示します。入退室時には即時、経過時間は1分ごとに同じメッセージを編集して更新します。全員がログアウトすると画像は削除されます。

経過時間はWrapperが検出したログイン時刻から計測し、再ログインでリセットされます。Wrapper再起動前の時刻は引き継ぎません。

Slackでも `SendJoinStateMessage: true` を指定すると、顔アイコン・名前・ログイン人数・経過時間を表示します。入退室時と1分ごとに既存のメッセージを更新します。

Slackの設定に関して、`SendJoinStateMessages`のオプションを指定する場合、次のようなフォーマットのユーザ一覧が表示されるようになります。

![SendJoinStateMessage](https://gyazo.com/b7327b2af897355258819807b4692be2.png)
