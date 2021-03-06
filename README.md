# Botto-LFG
### This is a lightweight but intelligent Discord Bot that keeps a certain number of free channels in a given category on your server with built-in optimizations for large servers.
It will be a great addition to your server if you would like to declutter by removing unused channels dynamically, while being flexible in case more users want to join.

How to invite the hosted bot - there is no guaranteed uptime:
* add the bot to your server using the [invite link](https://discord.com/api/oauth2/authorize?client_id=809953817562972160&permissions=1040&scope=bot)
* make sure the bot has permissions to view and manage channels
* use b-watchcat categoryID to enable for a category, b-lookaway categoryID to disable

How to self-host:
* download the latest release from https://github.com/Cufee/botto-lfg/releases or clone and run go build
* run the executable to generate a config file
* place your Discord Bot token into the config.json
* adjust any other settings as needed
* invite the bot to your server by generating and using OAuth2 Bot scope link
* run the executable
* use b-watchcat categoryID to enable for a category, b-lookaway categoryID to disable

Config File:
 * "prefix": string - Prefix for bot commands
 * "channels_buffer": int - How many free channels to keep around
 * "channel_user_limit": int - User limit for dynamically created channels, leaving this at 0 will set the limit based on existing channels in the category
 * "event_spacing": int - How often should the bot create/delete channels in seconds. I recommend increasing the timer for large servers to avoid spam.
 * "token": string - Your Discord Bot token
