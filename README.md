# sshhoneypot

SSH Honeypot, printing source IP, username and password used.
This is an attempt to collect some info on SSH brute force attacks from the evil people on the internet.

## Lib used

It uses [github.com/gliderlabs/ssh](https://github.com/gliderlabs/ssh)

## Output

It's in the logs

## Build

```
docker build -t sshhoneypot .
```

## Usage

```
docker run -d -n sshhoneypot -p 2222:2222 sshhoneypot
docker logs -f sshhoneypot
```

You can now test this with your local connection (assuming it's running on your machine):
```
ssh localhost -p 2222
```

type some random password and see what happens.

You can then redirect port 22 from your firewall/router to port 2222 on your local machine and wait.
You won't have to wait long to realise that the internet is a scary place.