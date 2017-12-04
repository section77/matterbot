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


## Arguments / Flags

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

- build **matterbot**: `go build -ldflags "-X main.version=$(git describe --tags --always)"`

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

_Note 1: this build is a multi-stage build so you need docker >= 17.05_

_Note 2: the build image has ~650MB_

- clone this repository: `git clone https://github.com/section77/matterbot.git`

- create an container: `docker build -t matterbot .`

#### Run the container

```
docker run --rm \
  --env FORWARD=user1=user1@mail.com,user2=abc@example.com \
  --env MATTERMOST_URL=http://ip:port \
  --env MATTERMOST_USER=muser \
  --env MATTERMOST_PASS=mpass \
  --env MAIL_HOST=smtp.gmail.com:465 \
  --env MAIL_USER=my-gmail-account@gmail.com \
  --env MAIL_PASS=my-gmail-pass \
  --env MAIL_USE_TLS=true \
  --name matterbot \
  jkeck/matterbot
```

## Usage

```
⟩ ./matterbot -h

matterbot forwards mattermost messages per mail, if their contains a configurable prefix marker.


Usage:

  To forward messages to 'user1@mail.com' and 'abc@example.com' call matterbot
  with the '-forward' flag:

    ./matterbot ... -forward user1=user1@mail.com,user2=abc@example.com ...

  If the chat-message contains any of the given prefix marker ('@user1', '@user2'),
  the message are send to the given mail address.

Flags:

  -forward string
        mapping from marker to receiver mail address. example: 'user1=user1@gmail.com,user2=abc@mail.com'
  -mail-body string
        mail body (default "{{.Content}}")
  -mail-host string
        mail-server host (default "127.0.0.1:25")
  -mail-pass string
        mail login pass (default "tobrettam")
  -mail-subject string
        mail subject (default "mattermost: {{.User}} writes in channel {{.Channel}}")
  -mail-use-tls
        use TLS instead of STARTTLS
  -mail-user string
        mail login user (default "matterbot@localhost")
  -mattermost-pass string
        mattermost password (default "tobrettam")
  -mattermost-url string
        mattermost url (default "http://127.0.0.1:8065")
  -mattermost-user string
        mattermost user (default "matterbot")
  -quiet
        disable logging / be quiet
  -v	show version and exit
  -verbose
        enable verbose / debug output
```
