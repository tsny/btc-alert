## Usage

`go get github.com/tsny/btc-alert`

Follow the configuration section below, go build, and run it

## Configuration

Create a `config.json` file in the root of the project.

`config.json` houses each 'interval' which represents an interval that is checked every minute.
If an interval lapses, then a notification is sent if the absolute value of the overall percentage change in the asset is 
more than the `percentThreshold` field.

`occurence`s are minutes.

Config file will look like this: 

```json
{
    "bootNotification": false,
    "discord": {
        "enabled": false,
        "token": "token.WgvCew.ObvImI8RA0HAn-kf2u8omFQGINI",
        "channelId": "715135158184218",
    },
    "yahooTickers":[
        "MSFT",
        "GM",
        "GME"
    ],
    "thresholds": [
        {
            "threshold": 500
        }
    ],
    "intervals": [
        {
            "maxOccurences": 5,
            "percentThreshold": 1
        },
        {
            "maxOccurences": 20,
            "percentThreshold": 2
        },
        {
            "maxOccurences": 60,
            "percentThreshold": 3
        }
    ]
}
```

### API

Checkout `routes.go` for a few of the endpoints, port defaults to 8080

`GET /crypto`
Returns the publishers associated with all tracked crypto coins
![Example API](https://i.imgur.com/eMDdj3S.png)


### Desktop Notifications

Example:

![Example Screenshot](https://i.imgur.com/VZ6bHZU.png)

Notifications work on both Mac and Windows

### Discord

If you have a Discord bot, you can redirect to there!

![Example Screenshot](https://i.imgur.com/zbDfI9B.png)

#### Discord Commands

{ticker}'s can be either be something like BTC or BTC-USD or btc

`!get {ticker}`

Returns the last minute candlestick for the provided ticker

`!trade {ticker}`

Returns the current price with a 1% and 2% difference for helping calculate trade costs

`!sub {ticker} {price}`

Subscribes to the event of a certin price point

`!stat {ticker}`

Returns the last 24 hour candlestick

`!gainers`
`!losers`

Returns the top market movers (gainers/losers) from Yahoo Finance

![Movers](https://i.imgur.com/MM9NXqE.png)

`!whois {ticker}`

Returns the company summary from MarketWatch 

![Summary](https://i.imgur.com/YY8wHq2.png)

## Logs Example

![Example Screenshot](https://i.imgur.com/lKS8kzG.png)

## Dependencies

- github.com/gen2brain/beeep     -  Notifications
- github.com/bwmarrin/discordgo  -  Discord API
- github.com/PuerkitoBio/goquery - Tables

- Yahoo Finance                  -  Ticker prices
- CoindeskAPI                    -  Coin Price Data
- CoinbaseAPI                    -  Coins
