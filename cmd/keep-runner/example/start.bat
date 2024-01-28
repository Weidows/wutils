@echo off
chcp 65001
tasklist | find /i "keep-runner" || powershell Start-Process -WindowStyle hidden keep-runner ol
echo 启动成功
pause
