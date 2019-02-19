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
4. Run linters
5. Store streamer events in database instead memory
6. Use Redis or some other storage for managing user sessions

## Issues
1. Twitch can send you some events just in case your API can be reached via https, so in this case I decided to mock event subscription URL and just used ```https://twitch.free.beeceptor.com/subscription``` for this purpose.
It means that use can subscribe for some events, but can't get real notifications from Twitch.

## How would you deploy the above on AWS?
I would Dockerize my application run a cluster of EC2 instances with Docker in Swarm mode or Kubernetes. I also would configure Jenkins CI/CD or gitlab CI/CD to deliver new images of my application, run linters and tests.

## Where do you see bottlenecks in your proposed architecture and how would you approach scaling this app starting from 100 reqs/day to 900MM reqs/day over 6 months?

Bottlenecks:
1. Managing user sessions, updating Twitch token.
2. Managing and storing streamer events. For the streamer events, I would use some queue, for example, RabbitMQ and a bunch of workers that can handle events from RabbitMQ and store it in a database. 

Scaling:
In case we use Kubernetes it is easy to add new EC instance to the cluster and run additional instances of our application depends on CPU or Memory loading. 
