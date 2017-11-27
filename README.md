# matterbot

**matterbot** forwards mattermost messages per mail, if their contains a configurable prefix marker.


_We use it to forward messages with the **@ml** prefix from our mattermost instance to our mailling list._

## Example

if you run **matterbot** with the following forward mappings:

```
./matterbot ... -forward user1=user1@mail.com,user2=abc@example.com ...
```

and you write a message in mattermost with any of the configured marker prefixes:

```
@user1, @mattermost-user, @user2 we meet us at 4pm
```

an email to `user1@mail.com` and `abc@example.com` with the body: `we meet us at 4pm` are send.


## Usage

FIXME


### Arguments / Flags

You can set the parameters per command line, or over environment variables.

| flag           | environment     | description (default)                      |
|----------------|-----------------|--------------------------------------------|
|-forward        | FORWARD         | mapping from marker to receiver address    |
|-mattermost-url | MATTERMOST_URL  | mattermost host _(http://127.0.0.1:8065)_  |
|-mattermost-user| MATTERMOST_USER | mattermost user _(matterbot@example.com)_  |
|-mattermost-pass| MATTERMOST_PASS | mattermost password _(tobrettam)_          |
|-mail-host      | MAIL_HOST       | mail host with port _(127.0.0.1:25)_       |
|-mail-user      | MAIL_USER       | mail user _(matterbot@localhost)_          |
|-mail-pass      | MAIL_PASS       | mail password _(tobrettam)_                |
|-mail-use-tls   | MAIL_USE_TLS    | use TLS instead of STARTTLS _(false -> use STARTTLS)_    |
|-mail-subject   | MAIL_SUBJECT    | _(mattermost: {{.User}} writes in channel {{.Channel}})_ |
|-mail-body      | MAIL_BODY       | _({{.Body}})_                              |
|-quiet          | QUIET           | be quiet _(false)_                         |
|-verbose        | VERBOSE         | enable verbose output _(false)_            |




## Run it

**matterbot** can run as an native application or in an docker container.

### Native app

- clone this repository: `git clone https://github.com/section77/matterbot.git`

- build **matterbot**: `go build -ldflags "-X main.version=$(git describe --tags)"`

- run it
```
./matterbot -forward user1=user1@mail.com,user2=abc@example.com \
  -mattermost-user muser \
  -mattermost-pass mpass \
  -mail-host smtp.gmail.com:465
  -mail-user my-gmail-account@gmail.com
  -mail-pass my-gmail-pass
  -mail-use-tls
```


### Docker Container

You can use a prebuild container (a) or create you own container (b)

#### a.) Use a prebuild container from hub.docker.com

- pull the container: `docker pull jkeck/matterbot`

#### b.) Build your own container

- clone this repository: `git clone https://github.com/section77/matterbot.git`

- build **matterbot**: `go build -ldflags "-X main.version=$(git describe --tags)"`

- create an container: `docker build -t matterbot .`

#### Run the container

```
docker run --rm \
  --env FORWARD=user1=user1@mail.com,user2=abc@example.com
  --env MATTERMOST_URL=http://ip:port \
  --env MATTERMOST_USER=muser \
  --env MATTERMOST_PASS=mpass \
  --env MAIL_HOST=smtp.gmail.com:465 \
  --env MAIL_USER=my-gmail-account@gmail.com \
  --env MAIL_PASS=my-gmail-pass \
  --env MAIL_USE_TLS=true \
  -t matterbot jkeck/matterbot
```
