# live-streamer
Mp3 streamer.

## The why
I need a service that can stream music to another device into my intranet.
For example a youtube uri that my old internet radio cannot play. Or a service for an internet radio receiver 
that I want to build using a ST32 controller and the VS1053 mp3 decoder.   

Since I have a couple of raspberry running, I can setup a streaming server on it.
There are a lot of streaming services that could be setup, but looking for those information is not funny enought.
I also want an easy access using my iphone without app store, also a web one.

I figure out how to start cvlc as streamer and built a very basic web  interface with vuetify
to start and stop it. Enought to build this repository and have some fun in 
developing a golang application for the arm6l processor like a raspberry (3 or 4).

I have tried to develop directly on my pihole device (Raspberry 3 with 1 Gb Ram),
but with Visual Code and go it was to slow for me. So I use it only for target device.
On the Raspberry 4 with 4G the developement works perfectly.

## Setup
An empty sqlite3 database is needed before starting the service. 
The empty database could be generated with:  
``` cat ./db/ref/player-repo-ref.db.sql | sqlite3 /home/pi/test-data.db ```

The file config.toml shold be also manually changed.