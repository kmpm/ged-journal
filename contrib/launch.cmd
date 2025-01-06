::@echo off
SET GED_NATS=tls://connect.ngs.global
SET GED_CREDS=..\.secrets\ngs.creds
:: SET GED_CONTEXT=NGS-ged-dev-me_cli
SET BINDIR=..\dist  

:: Timestamp in YYYY-MM-DD_HHMMSS format
SET TS=%DATE%_%TIME:~0,2%%TIME:~3,2%%TIME:~6,2%


cd %BINDIR%
MKDIR ..\..\output\%TS%
@echo on
@REM start cmd.exe /k "ged-journal-windows-amd64.exe agent -l debug -f ..\var/agent-log.jsonl --nats %GED_NATS% --nats-creds %GED_CREDS%"
start cmd.exe /k "ged-journal-windows-amd64.exe collect -f ..\var\collect-log.jsonl --nats %GED_NATS% --nats-creds %GED_CREDS%"

del ..\var\sub-file-log.jsonl
ged-journal-windows-amd64.exe subscribe file -f ..\var\sub-file-log.jsonl --subject "ged.collector.>" --nats %GED_NATS% --nats-creds %GED_CREDS% ..\..\output\%TS%"
pause