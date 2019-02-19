## How to run streamer_event_server?

If Golang  in not installed on your computer (```/bin``` folder has two binaries files, one for Mac and another one for Linux, in case you don't have Golang on your computer and can't build by your own)

Run in the project root folder (Mac)

```
  make run-mac
```

Run in the project root folder (Linux)

```
 make run-linux
```

If Golang installed on your computer

Run in the project root folder (Mac)

```
 make build-mac && make run-mac
```

Run in the project root folder (Linux)

```
 make build-linux && make run-linux
```

## How to test streamer_event_server?

1. You need to run streamer_event_server (```make run-mac``` or ```make run-linux```)
2. streamer_event_server will be started on port :7001
3. Open your browser and visit ``` http://localhost:7001/ ```
4. Press "Login using Twitch"
5. Press "Subscribe" to subscribe on your favorite streamer
6. Visit ``` http://localhost:7001/livestream ``` to watch livestream, read chat and watch events
7. Press "Logout" to logout from the application

## TODO
1. Move clientID, clientSecret to ENV variables
2. Add Dockerfile to have a change to run application as a container
3. Cover code by Unit test
4. Store streamer events in database instead memory
5. Use Redis or some other storage for managing user sessions

## Issues
1. Twitch can send you some events just in case your API can be reached via https, so in this case I decided to mock event subscription URL and just used ```https://twitch.free.beeceptor.com/subscription``` for this purpose.
It means that use can subscribe for some events, but can't get real notifications from Twitch.