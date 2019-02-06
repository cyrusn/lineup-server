# A simple line up server for parents day

## Build

`GOOS=linux GOARCH=amd64 go build -o ./lineup main.go`

## setup

- create a config file `./config.yaml`, please see below for the default key and value.

```yaml
# root command
config: "./config.yaml"
key: "skill-vein-planet-neigh-envoi"
dsn: "root@/lineupTestDB"

# server command
port: ":5000"
static: "./public"
# jwt token expire time in minutes
time: 300

# import command
import: ./data/user.json

# create command
overwrite: false
```

## import

```json
// import users by `import` cmd, the schema of the json file is as below
[
  {
    "userAlias": "teacher1",
    "password": "password1",
    "role": "teacher"
  },
  {
    "userAlias": "teacher2",
    "password": "password2",
    "role": "teacher"
  },
  {
    "userAlias": "student1",
    "password": "password1",
    "role": "student"
  }
]
```
