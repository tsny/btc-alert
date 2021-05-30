- [ ] Add beginTime -> endTime in movement alerts for like $1000

- [ ] Save queues into csv so they can be reloaded
    Is this one really necessary?

- [x] Streaks should be tracked both directions and mention the current price and percentage change

- [x] Basic graphs

- [ ] Disable Discord notifications/parsing via hostname? Something like `!toggle <hostname>`

- [ ] Maybe there should be thresholds that have messages get sent out but then we have addl. thresholds for alerts that hit @everyone

- [x] Alert when streak is over

- [ ] Steal the info section that CoinGecko from the discord channel returns, it would be cool

- [ ] !get command should return the latest price movements as well as current summary

- [ ] Price alerts shouldd be continous not in chunks. I.E, 5 minute intervals should overlap, but maybe we avoid alerts if the interval has already alerted recently unless the next interval is higher than the last like if it goes from 3% to 3.5%

- [ ] Have timed job that refreshes the movers watchlist on weekday mornings
