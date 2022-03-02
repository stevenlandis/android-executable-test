# Arm Test

This repo is for testing running executables on my kindle fire. Eventually I want to use the kindle as a server that runs in the background to run periodic tasks.

# Setup

```sh
cd site-js && yarn
cp env.example.json env.json # and fill in real values
```

# Running

```sh
node script.json run
```

This should compile everything and upload to the attached android device at `/data/local/tmp/dist`. It will then kick off `./serve` in the background and pass in `env.json` to stdin. Server logs will stored in `/data/local/tmp/dist/some.log`
