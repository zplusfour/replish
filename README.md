# replish [![Run on Replit](https://replit.com/badge/github/ReplDepot/replish)](https://replit.com/github/ReplDepot/replish)

A shell-based toolchain to extend the functionality on a repl

## Goals:

* Single binary (before integrating with replit) for convinence and compatibility. 
* Compatiblity with OS-native webdav mount & ssh client
* Easy to use, self-explainatory, beginner friendly.

## Functionality goals:

* ✅ Basic webdav mounting of a repl with basic auth
* ⚠️ Git repo on replit
* ❌ Command runner with helpful error messages similar to [rebound](https://github.com/shobrook/rebound)
* ❌ Log capture
* ❌ Health check route, ping route 
* ⚠️ WebSocket Proxy for port forwarding similar to ngrok, with a local client. 

Problems: 

* File changes may not persist
* Git push does not work with webdav
