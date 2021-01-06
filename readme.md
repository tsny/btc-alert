## Usage

`go get github.com/tsny/btc-alert`

Follow the configuration section below, go build, and run it

## Configuration

`config.json` houses each 'interval' which represents an interval that is checked every minute.
If an interval lapses, then a notification is sent if the absolute value of the overall percentage change in the asset is 
more than the `percentThreshold` field.

`occurence`s are minutes.

### Discord

If you have a Discord bot, you can redirect to there!

![Example Screenshot](https://i.imgur.com/zbDfI9B.png)

Config file will look like this: 

```json
{
    "token": "Wgaisdg.OOLSIDG-kf2u8omFQGINI",
    "channelId": "815218518838123",
    "clearOnBoot": true
}
```

## Example

![Example Screenshot](https://i.imgur.com/lKS8kzG.png)

## Dependencies

- github.com/gen2brain/beeep     -  Notifications
- github.com/bwmarrin/discordgo  -  Discord API
- Yahoo Finance                  -  BTC Price Data
