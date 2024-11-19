@echo off
title Stellar By Vin505

echo Welcome to StellarV4
echo Note: This tool is still in beta version.
echo Update: Low CPU usage, more powerfull.
echo.
set /p target_url=Enter target URL: 
set /p workers=Enter number of workers: 
set /p duration=Enter duration (in seconds): 
echo.
rem
go run main.go %target_url% %workers% %duration%