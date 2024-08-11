@echo off
chcp 65001
tasklist | find /i "wutils" || powershell Start-Process -WindowStyle hidden ../wutils ol
echo 启动成功
pause
