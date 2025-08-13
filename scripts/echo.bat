@echo off
REM Usage: http POST :8080/echo msg==hello foo=bar

echo This is a simple echo hook.

echo Hook information: name=%hook_name%, id=%hook_id%, method=%hook_method%

for /f "tokens=*" %%a in ('hostname') do (
    echo Command result: hostname=%%a
)

echo Header variable: User-Agent=%user_agent%

echo Query parameter: msg=%msg%

echo Body payload: %1