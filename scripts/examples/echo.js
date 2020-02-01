#!/usr/bin/env node

const os = require('os')

// Usage: http POST :8080/examples/echo.js msg==hello foo=bar

console.log("Echo script:")

console.log("Hostname=", os.hostname())

console.log("User-Agent=", process.env["user_agent"])

console.log("msg=", process.env["msg"])

console.log("body=", process.argv)
