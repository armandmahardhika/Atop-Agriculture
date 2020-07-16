# AtopIOTServer
The Atop backend project

# Dependent services install guide
## MongoDB 
### Install
#### MacOS

```bash
brew tap mongodb/brew
brew install mongodb-community@4.2
brew services start mongodb-community@4.2
```
### Run startup script
#### MacOS & Linux
In the folder `AtopIOTServer/asset/mongodb` 
```bash
mongo localhost:27010/atop init.js
```
### Stop

#### MacOS

```bash
brew services stop mongodb-community@4.2
```

### Check mongo running
#### MacOS & Linux
```bash
ps aux | grep -v grep | grep mongod
```

## Mosquitto
### Install
#### MacOS
```bash
brew install mosquitto
brew services start mosquitto
```
## Run 
### In develop
Just type command ` go run main.go` in the root folder of project.
