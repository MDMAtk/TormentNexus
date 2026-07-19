; TormentNexus NSIS Installer Script
; Creates a professional Windows installer

!include "MUI2.nsh"
!include "FileFunc.nsh"

; General
Name "TormentNexus"
OutFile "tormentnexus-setup.exe"
InstallDir "$PROGRAMFILES\TormentNexus"
InstallDirRegKey HKCU "Software\TormentNexus" ""
RequestExecutionLevel admin

; Version Info
VIProductVersion "1.0.0.0"
VIAddVersionKey "ProductName" "TormentNexus"
VIAddVersionKey "CompanyName" "TormentNexus Team"
VIAddVersionKey "FileVersion" "1.0.0"
VIAddVersionKey "FileDescription" "TormentNexus AI Control Plane Installer"
VIAddVersionKey "LegalCopyright" "Copyright 2026 TormentNexus Team"

; Interface Settings
!define MUI_ABORTWARNING
!define MUI_ICON "${NSISDIR}\Contrib\Graphics\Icons\modern-install.ico"
!define MUI_UNICON "${NSISDIR}\Contrib\Graphics\Icons\modern-uninstall.ico"

; Pages
!insertmacro MUI_PAGE_LICENSE "LICENSE.txt"
!insertmacro MUI_PAGE_COMPONENTS
!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_PAGE_FINISH

!insertmacro MUI_UNPAGE_CONFIRM
!insertmacro MUI_UNPAGE_INSTFILES

; Languages
!insertmacro MUI_LANGUAGE "English"

; Installer Sections
Section "TormentNexus Core" SecCore
  SectionIn RO
  
  SetOutPath "$INSTDIR\bin"
  File "bin\tormentnexus.exe"
  
  SetOutPath "$INSTDIR"
  File "README.md"
  File "LICENSE.txt"
  
  ; Create uninstaller
  WriteUninstaller "$INSTDIR\uninstall.exe"
  
  ; Create start menu items
  CreateDirectory "$SMPROGRAMS\TormentNexus"
  CreateShortCut "$SMPROGRAMS\TormentNexus\TormentNexus.lnk" "$INSTDIR\bin\tormentnexus.exe" "serve"
  CreateShortCut "$SMPROGRAMS\TormentNexus\Uninstall.lnk" "$INSTDIR\uninstall.exe"
  
  ; Create desktop shortcut
  CreateShortCut "$DESKTOP\TormentNexus.lnk" "$INSTDIR\bin\tormentnexus.exe" "serve"
  
  ; Write registry keys
  WriteRegStr HKCU "Software\TormentNexus" "" $INSTDIR
  WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\TormentNexus" \
    "DisplayName" "TormentNexus"
  WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\TormentNexus" \
    "UninstallString" "$\"$INSTDIR\uninstall.exe$\""
  WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\TormentNexus" \
    "DisplayIcon" "$\"$INSTDIR\bin\tormentnexus.exe$\""
  WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\TormentNexus" \
    "Publisher" "TormentNexus Team"
  WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\TormentNexus" \
    "DisplayVersion" "1.0.0"
  
  ; Get installed size
  ${GetSize} "$INSTDIR" "/S=0K" $0 $1 $2
  IntFmt $0 "0x%08X" $0
  WriteRegDWORD HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\TormentNexus" \
    "EstimatedSize" "$0"
SectionEnd

Section "Configuration" SecConfig
  SetOutPath "$PROFILE\.tormentnexus"
  
  ; Create default config if it doesn't exist
  IfFileExists "$PROFILE\.tormentnexus\config.yaml" config_exists
    FileOpen $0 "$PROFILE\.tormentnexus\config.yaml" w
    FileWrite $0 "# TormentNexus Configuration$\r$\n"
    FileWrite $0 "host: 127.0.0.1$\r$\n"
    FileWrite $0 "port: 7778$\r$\n"
    FileWrite $0 "$\r$\n"
    FileWrite $0 "# Memory Configuration$\r$\n"
    FileWrite $0 "memory:$\r$\n"
    FileWrite $0 "  l2_enabled: true$\r$\n"
    FileWrite $0 "  l3_enabled: true$\r$\n"
    FileWrite $0 "  l4_enabled: false$\r$\n"
    FileWrite $0 "$\r$\n"
    FileWrite $0 "# Provider Configuration$\r$\n"
    FileWrite $0 "providers:$\r$\n"
    FileWrite $0 "  deepseek:$\r$\n"
    FileWrite $0 "    enabled: true$\r$\n"
    FileWrite $0 "    api_key: `$\"$\"$\r$\n"
    FileWrite $0 "  lmstudio:$\r$\n"
    FileWrite $0 "    enabled: true$\r$\n"
    FileWrite $0 "    url: http://127.0.0.1:1234$\r$\n"
    FileClose $0
  config_exists:
  
  CreateDirectory "$PROFILE\.tormentnexus\memory"
SectionEnd

Section "Add to PATH" SecPath
  ; Add to user PATH
  ReadRegStr $0 HKCU "Environment" "PATH"
  WriteRegStr HKCU "Environment" "PATH" "$0;$INSTDIR\bin"
  
  ; Broadcast environment change
  SendMessage ${HWND_BROADCAST} ${WM_WININICHG} 0 "STR:Environment" /TIMEOUT=5000
SectionEnd

Section "Start Menu Shortcuts" SecStartMenu
  CreateDirectory "$SMPROGRAMS\TormentNexus"
  CreateShortCut "$SMPROGRAMS\TormentNexus\TormentNexus.lnk" "$INSTDIR\bin\tormentnexus.exe" "serve"
  CreateShortCut "$SMPROGRAMS\TormentNexus\Dashboard.lnk" "http://127.0.0.1:7778"
  CreateShortCut "$SMPROGRAMS\TormentNexus\Uninstall.lnk" "$INSTDIR\uninstall.exe"
SectionEnd

; Descriptions
!insertmacro MUI_FUNCTION_DESCRIPTION_BEGIN
  !insertmacro MUI_DESCRIPTION_TEXT ${SecCore} "Core TormentNexus files (required)"
  !insertmacro MUI_DESCRIPTION_TEXT ${SecConfig} "Create default configuration files"
  !insertmacro MUI_DESCRIPTION_TEXT ${SecPath} "Add TormentNexus to system PATH"
  !insertmacro MUI_DESCRIPTION_TEXT ${SecStartMenu} "Create Start Menu shortcuts"
!insertmacro MUI_FUNCTION_DESCRIPTION_END

; Uninstaller Section
Section "Uninstall"
  ; Remove files
  RMDir /r "$INSTDIR"
  
  ; Remove shortcuts
  RMDir /r "$SMPROGRAMS\TormentNexus"
  Delete "$DESKTOP\TormentNexus.lnk"
  
  ; Remove registry keys
  DeleteRegKey HKCU "Software\TormentNexus"
  DeleteRegKey HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\TormentNexus"
  
  ; Remove from PATH (simplified - in real installer would parse and remove)
  ReadRegStr $0 HKCU "Environment" "PATH"
  ; Note: Proper PATH removal would require string parsing
  
  ; Note: We don't remove user data ($PROFILE\.tormentnexus) by default
  ; Users can manually delete it if they want
SectionEnd
