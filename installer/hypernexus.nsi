; HyperNexus / TormentNexus NSIS Installer Script
; Supports both Open Source (TormentNexus) and Corporate (HyperNexus) editions

!include "MUI2.nsh"
!include "FileFunc.nsh"
!include "LogicLib.nsh"

; ─── Edition Detection ───────────────────────────────────────────────
; The edition is determined by the presence of a file or command-line parameter
; To build corporate edition: makensis /DCORPORATE=1 hypernexus.nsi
; To build open source edition: makensis hypernexus.nsi

!ifdef CORPORATE
  !define EDITION "hypernexus"
  !define PRODUCT_NAME "HyperNexus"
  !define COMPANY_NAME "HyperNexus Corp"
  !define PRODUCT_DESC "HyperNexus Corporate AI Control Plane"
  !define CONFIG_DIR ".hypernexus"
  !define REG_KEY "HyperNexus"
  !define TRAY_TIP "HyperNexus (Running)"
  !define EXE_NAME "hypernexus.exe"
  !define OUT_NAME "hypernexus-setup.exe"
  !define CLOUD_MODE "true"
!else
  !define EDITION "tormentnexus"
  !define PRODUCT_NAME "TormentNexus"
  !define COMPANY_NAME "TormentNexus Team"
  !define PRODUCT_DESC "TormentNexus Open Source AI Control Plane"
  !define CONFIG_DIR ".tormentnexus"
  !define REG_KEY "TormentNexus"
  !define TRAY_TIP "TormentNexus (Running)"
  !define EXE_NAME "tormentnexus.exe"
  !define OUT_NAME "tormentnexus-setup.exe"
  !define CLOUD_MODE "false"
!endif

; General
Name "${PRODUCT_NAME}"
OutFile "${OUT_NAME}"
InstallDir "$PROGRAMFILES\${PRODUCT_NAME}"
InstallDirRegKey HKCU "Software\${REG_KEY}" ""
RequestExecutionLevel admin

; Version Info
VIProductVersion "1.0.0.0"
VIAddVersionKey "ProductName" "${PRODUCT_NAME}"
VIAddVersionKey "CompanyName" "${COMPANY_NAME}"
VIAddVersionKey "FileVersion" "1.0.0"
VIAddVersionKey "FileDescription" "${PRODUCT_DESC} Installer"
VIAddVersionKey "LegalCopyright" "Copyright 2026 ${COMPANY_NAME}"

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

; ─── Installer Sections ──────────────────────────────────────────────

Section "${PRODUCT_NAME} Core" SecCore
  SectionIn RO
  
  SetOutPath "$INSTDIR\bin"
  File "bin\${EXE_NAME}"
  
  SetOutPath "$INSTDIR"
  File "README.md"
  File "LICENSE.txt"
  
  ; Write edition marker file
  FileOpen $0 "$INSTDIR\edition.txt" w
  FileWrite $0 "${EDITION}"
  FileClose $0
  
  ; Create uninstaller
  WriteUninstaller "$INSTDIR\uninstall.exe"
  
  ; Create start menu items
  CreateDirectory "$SMPROGRAMS\${PRODUCT_NAME}"
  CreateShortCut "$SMPROGRAMS\${PRODUCT_NAME}\${PRODUCT_NAME}.lnk" "$INSTDIR\bin\${EXE_NAME}" "serve"
  CreateShortCut "$SMPROGRAMS\${PRODUCT_NAME}\Uninstall.lnk" "$INSTDIR\uninstall.exe"
  
  ; Create desktop shortcut
  CreateShortCut "$DESKTOP\${PRODUCT_NAME}.lnk" "$INSTDIR\bin\${EXE_NAME}" "serve"
  
  ; Write registry keys
  WriteRegStr HKCU "Software\${REG_KEY}" "" $INSTDIR
  WriteRegStr HKCU "Software\${REG_KEY}" "Edition" "${EDITION}"
  WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\${REG_KEY}" \
    "DisplayName" "${PRODUCT_NAME}"
  WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\${REG_KEY}" \
    "UninstallString" "$\"$INSTDIR\uninstall.exe$\""
  WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\${REG_KEY}" \
    "DisplayIcon" "$\"$INSTDIR\bin\${EXE_NAME}$\""
  WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\${REG_KEY}" \
    "Publisher" "${COMPANY_NAME}"
  WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\${REG_KEY}" \
    "DisplayVersion" "1.0.0"
  
  ; Get installed size
  ${GetSize} "$INSTDIR" "/S=0K" $0 $1 $2
  IntFmt $0 "0x%08X" $0
  WriteRegDWORD HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\${REG_KEY}" \
    "EstimatedSize" "$0"
SectionEnd

Section "Configuration" SecConfig
  SetOutPath "$PROFILE\${CONFIG_DIR}"
  
  ; Create branding configuration file
  FileOpen $0 "$PROFILE\${CONFIG_DIR}\branding.json" w
  FileWrite $0 '{$\r$\n'
  FileWrite $0 '  "edition": "${EDITION}",$\r$\n'
  FileWrite $0 '  "product_name": "${PRODUCT_NAME}",$\r$\n'
  FileWrite $0 '  "company_name": "${COMPANY_NAME}",$\r$\n'
  FileWrite $0 '  "tray_tooltip": "${TRAY_TIP}",$\r$\n'
  FileWrite $0 '  "dashboard_title": "${PRODUCT_NAME} Dashboard",$\r$\n'
  FileWrite $0 '  "config_dir": "${CONFIG_DIR}",$\r$\n'
  FileWrite $0 '  "registry_key": "${REG_KEY}"$\r$\n'
  FileWrite $0 '}$\r$\n'
  FileClose $0
  
  ; Create default config if it doesn't exist
  IfFileExists "$PROFILE\${CONFIG_DIR}\config.yaml" config_exists
    FileOpen $0 "$PROFILE\${CONFIG_DIR}\config.yaml" w
    FileWrite $0 "# ${PRODUCT_NAME} Configuration$\r$\n"
    FileWrite $0 "host: 127.0.0.1$\r$\n"
    FileWrite $0 "port: 7778$\r$\n"
    FileWrite $0 "$\r$\n"
    
    !ifdef CORPORATE
      ; Corporate mode cloud configuration
      FileWrite $0 "# Cloud Configuration (HyperNexus Corporate)$\r$\n"
      FileWrite $0 "cloud:$\r$\n"
      FileWrite $0 "  enabled: true$\r$\n"
      FileWrite $0 "  endpoint: https://api.hypernexus.io$\r$\n"
      FileWrite $0 "  auth_token: `$\r$\n"
      FileWrite $0 "$\r$\n"
    !endif
    
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
  
  CreateDirectory "$PROFILE\${CONFIG_DIR}\memory"
SectionEnd

Section "Add to PATH" SecPath
  ; Add to user PATH
  ReadRegStr $0 HKCU "Environment" "PATH"
  WriteRegStr HKCU "Environment" "PATH" "$0;$INSTDIR\bin"
  
  ; Broadcast environment change
  SendMessage ${HWND_BROADCAST} ${WM_WININICHG} 0 "STR:Environment" /TIMEOUT=5000
SectionEnd

Section "Start Menu Shortcuts" SecStartMenu
  CreateDirectory "$SMPROGRAMS\${PRODUCT_NAME}"
  CreateShortCut "$SMPROGRAMS\${PRODUCT_NAME}\${PRODUCT_NAME}.lnk" "$INSTDIR\bin\${EXE_NAME}" "serve"
  CreateShortCut "$SMPROGRAMS\${PRODUCT_NAME}\Dashboard.lnk" "http://127.0.0.1:7779"
  CreateShortCut "$SMPROGRAMS\${PRODUCT_NAME}\Uninstall.lnk" "$INSTDIR\uninstall.exe"
SectionEnd

!ifdef CORPORATE
Section "Cloud Connection" SecCloud
  ; Additional cloud connection configuration for corporate edition
  SetOutPath "$PROFILE\${CONFIG_DIR}"
  
  ; Create cloud connection script
  FileOpen $0 "$PROFILE\${CONFIG_DIR}\connect-cloud.bat" w
  FileWrite $0 "@echo off$\r$\n"
  FileWrite $0 "echo Connecting to HyperNexus Cloud...$\r$\n"
  FileWrite $0 "echo.$\r$\n"
  FileWrite $0 "echo Please enter your HyperNexus Cloud credentials:$\r$\n"
  FileWrite $0 "set /p CLOUD_AUTH="Authentication Token: "$\r$\n"
  FileWrite $0 "echo.$\r$\n"
  FileWrite $0 "echo Saving configuration...$\r$\n"
  FileWrite $0 "setx HN_CLOUD_AUTH "%CLOUD_AUTH%"$\r$\n"
  FileWrite $0 "setx HN_CLOUD_ENDPOINT "https://api.hypernexus.io"$\r$\n"
  FileWrite $0 "echo.$\r$\n"
  FileWrite $0 "echo Configuration saved! Please restart ${PRODUCT_NAME}.$\r$\n"
  FileWrite $0 "pause$\r$\n"
  FileClose $0
  
  CreateShortCut "$SMPROGRAMS\${PRODUCT_NAME}\Connect to Cloud.lnk" "$PROFILE\${CONFIG_DIR}\connect-cloud.bat"
SectionEnd
!endif

; ─── Descriptions ────────────────────────────────────────────────────

!insertmacro MUI_FUNCTION_DESCRIPTION_BEGIN
  !insertmacro MUI_DESCRIPTION_TEXT ${SecCore} "Core ${PRODUCT_NAME} files (required)"
  !insertmacro MUI_DESCRIPTION_TEXT ${SecConfig} "Create default configuration files"
  !insertmacro MUI_DESCRIPTION_TEXT ${SecPath} "Add ${PRODUCT_NAME} to system PATH"
  !insertmacro MUI_DESCRIPTION_TEXT ${SecStartMenu} "Create Start Menu shortcuts"
  !ifdef CORPORATE
    !insertmacro MUI_DESCRIPTION_TEXT ${SecCloud} "Configure cloud connection to HyperNexus Cloud"
  !endif
!insertmacro MUI_FUNCTION_DESCRIPTION_END

; ─── Uninstaller Section ─────────────────────────────────────────────

Section "Uninstall"
  ; Remove files
  RMDir /r "$INSTDIR"
  
  ; Remove shortcuts
  RMDir /r "$SMPROGRAMS\${PRODUCT_NAME}"
  Delete "$DESKTOP\${PRODUCT_NAME}.lnk"
  
  ; Remove registry keys
  DeleteRegKey HKCU "Software\${REG_KEY}"
  DeleteRegKey HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\${REG_KEY}"
  
  ; Remove from PATH (simplified - in real installer would parse and remove)
  ReadRegStr $0 HKCU "Environment" "PATH"
  ; Note: Proper PATH removal would require string parsing
  
  ; Note: We don't remove user data ($PROFILE\${CONFIG_DIR}) by default
  ; Users can manually delete it if they want
SectionEnd
