@REM Watch out! Chinese word will cause error in this file

@echo off
cd ..\..\

tasklist | find /i "wutils" || powershell Start-Process -WindowStyle hidden ./wutils ol

echo Start Successfully
pause
