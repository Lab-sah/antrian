@echo off
set "BROWSER=%ProgramFiles%\Google\Chrome\Application\chrome.exe"
if not exist "%BROWSER%" set "BROWSER=%ProgramFiles(x86)%\Google\Chrome\Application\chrome.exe"
if not exist "%BROWSER%" set "BROWSER=%LocalAppData%\Google\Chrome\Application\chrome.exe"
if not exist "%BROWSER%" set "BROWSER=%ProgramFiles(x86)%\Microsoft\Edge\Application\msedge.exe"
if not exist "%BROWSER%" set "BROWSER=%ProgramFiles%\Microsoft\Edge\Application\msedge.exe"
if not exist "%BROWSER%" (
  echo Chrome atau Microsoft Edge tidak ditemukan di PC ini.
  pause
  exit /b 1
)

set "URL="

rem Kalau link.md ada, pakai baris pertama yang diawali http.
if exist "%~dp0link.md" (
  for /f "usebackq delims=" %%L in (`findstr /b /i "http" "%~dp0link.md"`) do (
    set "URL=%%L"
    goto :launch
  )
)

rem link.md belum ada atau belum berisi link: minta link, lalu simpan ke link.md.
:ask
cls
echo ==============================================
echo   Link kiosk antrean belum diatur.
echo   Paste link web kiosk di bawah ini, lalu Enter.
echo   Tips: klik kanan di jendela ini untuk paste.
echo ==============================================
echo.
set "URL="
set /p "URL=Link: "
if /i not "%URL:~0,4%"=="http" (
  echo.
  echo Link harus diawali http atau https. Coba lagi.
  pause
  goto :ask
)

setlocal EnableDelayedExpansion
>"%~dp0link.md" echo(!URL!
endlocal

:launch
start "" "%BROWSER%" --app="%URL%" --start-fullscreen --kiosk-printing --user-data-dir="%LocalAppData%\ErzapKioskAntrean"
