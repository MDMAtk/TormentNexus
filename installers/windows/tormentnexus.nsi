; TormentNexus NSIS Installer
; Modern GUI installer for Windows

!include "MUI2.nsh"
!include "FileFunc.nsh"

; ── General ──
Name "TormentNexus"
OutFile "tormentnexus-setup.exe"
InstallDir "$LOCALAPPDATA\TormentNexus"
InstallDirRegKey HKCU "Software\TormentNexus" ""
RequestExecutionLevel user
Unicode True

; ── Version Info ──
VIProductVersion "1.0.0.0"
VIAddVersionKey "ProductName" "TormentNexus"
VIAddVersionKey "CompanyName" "TormentNexus"
VIAddVersionKey "FileDescription" "AI Control Plane with Persistent Memory"
VIAddVersionKey "FileVersion" "1.0.0"
VIAddVersionKey "ProductVersion" "1.0.0"
VIAddVersionKey "LegalCopyright" "© 2026 TormentNexus"

; ── MUI Settings ──
!define MUI_ABORTWARNING
!define MUI_ICON "${NSISDIR}\Contrib\Graphics\Icons\modern-install.ico"
!define MUI_UNICON "${NSISDIR}\Contrib\Graphics\Icons\modern-uninstall.ico"

; ── Welcome Page ──
!define MUI_WELCOMEPAGE_TITLE "Welcome to TormentNexus Setup"
!define MUI_WELCOMEPAGE_TEXT "TormentNexus is an AI Control Plane with Persistent Memory.$\r$\n$\r$\nFeatures:$\r$\n  - 26,000+ MCP Tools$\r$\n  - 4-Tier Memory System$\r$\n  - Multi-Agent Orchestration$\r$\n  - Local-First Architecture$\r$\n$\r$\nThis wizard will guide you through the installation.$\r$\n$\r$\nClick Next to continue."

; ── Pages ──
!insertmacro MUI_PAGE_WELCOME
!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_INSTFILES
!define MUI_FINISHPAGE_RUN "$INSTDIR\tormentnexus.exe"
!define MUI_FINISHPAGE_RUN_PARAMETERS "serve"
!define MUI_FINISHPAGE_RUN_TEXT "Start TormentNexus server"
!define MUI_FINISHPAGE_LINK "Visit TormentNexus website"
!define MUI_FINISHPAGE_LINK_LOCATION "https://tormentnexus.site"
!insertmacro MUI_PAGE_FINISH

!insertmacro MUI_UNPAGE_CONFIRM
!insertmacro MUI_UNPAGE_INSTFILES

; ── Languages ──
!insertmacro MUI_LANGUAGE "English"

; ── Installer Sections ──
Section "TormentNexus (Required)" SecMain
    SectionIn RO
    
    ; Set output path
    SetOutPath "$INSTDIR"
    
    ; Install files
    File "tormentnexus.exe"
    
    ; Create config directory
    CreateDirectory "$INSTDIR\.tormentnexus"
    CreateDirectory "$INSTDIR\logs"
    
    ; Create default config
    FileOpen $0 "$INSTDIR\.tormentnexus\config.json" w
    FileWrite $0 '{$\r$\n'
    FileWrite $0 '  "version": "1.0.0",$\r$\n'
    FileWrite $0 '  "server": {$\r$\n'
    FileWrite $0 '    "host": "127.0.0.1",$\r$\n'
    FileWrite $0 '    "port": 7778$\r$\n'
    FileWrite $0 '  },$\r$\n'
    FileWrite $0 '  "memory": {$\r$\n'
    FileWrite $0 '    "enabled": true,$\r$\n'
    FileWrite $0 '    "tiers": ["L1", "L2", "L3", "L4"]$\r$\n'
    FileWrite $0 '  },$\r$\n'
    FileWrite $0 '  "mcp": {$\r$\n'
    FileWrite $0 '    "catalog": true,$\r$\n'
    FileWrite $0 '    "autoInstall": false$\r$\n'
    FileWrite $0 '  }$\r$\n'
    FileWrite $0 '}$\r$\n'
    FileClose $0
    
    ; Store installation folder
    WriteRegStr HKCU "Software\TormentNexus" "" $INSTDIR
    
    ; Create uninstaller
    WriteUninstaller "$INSTDIR\Uninstall.exe"
    
    ; Add to Add/Remove Programs
    WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\TormentNexus" \
        "DisplayName" "TormentNexus"
    WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\TormentNexus" \
        "UninstallString" "$\"$INSTDIR\Uninstall.exe$\""
    WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\TormentNexus" \
        "QuietUninstallString" "$\"$INSTDIR\Uninstall.exe$\" /S"
    WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\TormentNexus" \
        "InstallLocation" "$INSTDIR"
    WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\TormentNexus" \
        "DisplayVersion" "1.0.0"
    WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\TormentNexus" \
        "Publisher" "TormentNexus"
    WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\TormentNexus" \
        "URLInfoAbout" "https://tormentnexus.site"
    WriteRegDWORD HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\TormentNexus" \
        "NoModify" 1
    WriteRegDWORD HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\TormentNexus" \
        "NoRepair" 1
    
    ; Get installed size
    ${GetSize} "$INSTDIR" "/S=0K" $0 $1 $2
    IntFmt $0 "0x%08X" $0
    WriteRegDWORD HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\TormentNexus" \
        "EstimatedSize" "$0"
SectionEnd

Section "Desktop Shortcut" SecDesktop
    CreateShortcut "$DESKTOP\TormentNexus.lnk" "$INSTDIR\tormentnexus.exe" "serve" "$INSTDIR\tormentnexus.exe" 0
SectionEnd

Section "Start Menu Shortcuts" SecStartMenu
    CreateDirectory "$SMPROGRAMS\TormentNexus"
    CreateShortcut "$SMPROGRAMS\TormentNexus\Start TormentNexus.lnk" "$INSTDIR\tormentnexus.exe" "serve" "$INSTDIR\tormentnexus.exe" 0
    CreateShortcut "$SMPROGRAMS\TormentNexus\Dashboard.lnk" "http://localhost:7778" "" "" 0
    CreateShortcut "$SMPROGRAMS\TormentNexus\Uninstall.lnk" "$INSTDIR\Uninstall.exe" "" "$INSTDIR\Uninstall.exe" 0
SectionEnd

; ── Uninstaller Section ──
Section "Uninstall"
    ; Remove files
    Delete "$INSTDIR\tormentnexus.exe"
    Delete "$INSTDIR\Uninstall.exe"
    
    ; Remove config (ask first)
    MessageBox MB_YESNO "Remove configuration and data?" IDNO NoDeleteConfig
        RMDir /r "$INSTDIR\.tormentnexus"
        RMDir /r "$INSTDIR\logs"
    NoDeleteConfig:
    
    ; Remove directories
    RMDir "$INSTDIR"
    
    ; Remove shortcuts
    Delete "$DESKTOP\TormentNexus.lnk"
    RMDir /r "$SMPROGRAMS\TormentNexus"
    
    ; Remove registry keys
    DeleteRegKey HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\TormentNexus"
    DeleteRegKey HKCU "Software\TormentNexus"
SectionEnd
