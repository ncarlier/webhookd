#!/usr/bin/env node

const os = require('os')

// Usage: http POST :8080/examples/echo.js msg==hello foo=bar

console.log("This is a simple echo hook using NodeJS.")

const { hook_name, hook_id , user_agent, msg} = process.env

console.log(`Hook information: name=${hook_name}, id=${hook_id}`)

console.log("Hostname=", os.hostname())

console.log(`User-Agent=${user_agent}`)

console.log(`msg=${msg}`)

console.log("body=", process.argv)
