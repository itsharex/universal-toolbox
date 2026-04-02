; ========================================================
; XTool NSIS 安装脚本
; 用于创建 Windows 安装包
; ========================================================

!define APP_NAME "XTool"
!define APP_NAME_EN "XTool"
!define VERSION "1.0.0"
!define PUBLISHER "MasterPick"
!define WEB_SITE "https://github.com/MasterPick/universal-toolbox"
!define UNINST_KEY "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APP_NAME}"
!define INSTDIR_KEY "Software\${APP_NAME}"

; 包含现代 UI
!include "MUI2.nsh"
!include "FileFunc.nsh"

; 安装程序属性
Name "${APP_NAME}"
OutFile "XTool_Setup_${VERSION}.exe"
InstallDir "$PROGRAMFILES64\${APP_NAME}"
InstallDirRegKey HKLM "${INSTDIR_KEY}" "Install_Dir"
RequestExecutionLevel admin
SetCompressor /SOLID lzma

; 界面设置
!define MUI_ICON "build\appicon.ico"
!define MUI_UNICON "build\appicon.ico"
!define MUI_WELCOMEFINISHPAGE_BITMAP "build\installer-welcome.bmp"
!define MUI_ABORTWARNING

; 页面
!insertmacro MUI_PAGE_WELCOME
!insertmacro MUI_PAGE_LICENSE "LICENSE"
!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_PAGE_FINISH

!insertmacro MUI_UNPAGE_CONFIRM
!insertmacro MUI_UNPAGE_INSTFILES

; 语言
!insertmacro MUI_LANGUAGE "SimpChinese"
!insertmacro MUI_LANGUAGE "English"

; 安装段
Section "Install" SecInstall
  SectionIn RO

  ; 设置输出路径
  SetOutPath "$INSTDIR"

  ; 清理旧版本
  RMDir /r "$INSTDIR\*.*"

  ; 复制文件
  File "build\bin\XTool.exe"
  File "README.md"
  File "LICENSE"

  ; 创建卸载程序
  WriteUninstaller "$INSTDIR\Uninstall.exe"

  ; 写入注册表
  WriteRegStr HKLM "${INSTDIR_KEY}" "Install_Dir" "$INSTDIR"
  WriteRegStr HKLM "${INSTDIR_KEY}" "Version" "${VERSION}"
  WriteRegStr HKLM "${UNINST_KEY}" "DisplayName" "${APP_NAME}"
  WriteRegStr HKLM "${UNINST_KEY}" "DisplayVersion" "${VERSION}"
  WriteRegStr HKLM "${UNINST_KEY}" "Publisher" "${PUBLISHER}"
  WriteRegStr HKLM "${UNINST_KEY}" "UninstallString" '"$INSTDIR\Uninstall.exe"'
  WriteRegStr HKLM "${UNINST_KEY}" "DisplayIcon" '"$INSTDIR\XTool.exe"'
  WriteRegStr HKLM "${UNINST_KEY}" "URLInfoAbout" "${WEB_SITE}"
  WriteRegDWORD HKLM "${UNINST_KEY}" "NoModify" 1
  WriteRegDWORD HKLM "${UNINST_KEY}" "NoRepair" 1

  ; 计算安装大小
  ${GetSize} "$INSTDIR" "/S=0K" $0 $1 $2
  IntFmt $0 "0x%08X" $0
  WriteRegDWORD HKLM "${UNINST_KEY}" "EstimatedSize" "$0"

  ; 创建开始菜单快捷方式
  CreateDirectory "$SMPROGRAMS\${APP_NAME}"
  CreateShortCut "$SMPROGRAMS\${APP_NAME}\${APP_NAME}.lnk" "$INSTDIR\XTool.exe"
  CreateShortCut "$SMPROGRAMS\${APP_NAME}\卸载.lnk" "$INSTDIR\Uninstall.exe"

  ; 创建桌面快捷方式（重要：使用 .lnk 而不是 .link）
  CreateShortCut "$DESKTOP\${APP_NAME}.lnk" "$INSTDIR\XTool.exe"

SectionEnd

; 卸载段
Section "Uninstall"
  ; 删除文件
  Delete "$INSTDIR\XTool.exe"
  Delete "$INSTDIR\README.md"
  Delete "$INSTDIR\LICENSE"
  Delete "$INSTDIR\Uninstall.exe"

  ; 删除目录
  RMDir "$INSTDIR"

  ; 删除快捷方式
  Delete "$SMPROGRAMS\${APP_NAME}\${APP_NAME}.lnk"
  Delete "$SMPROGRAMS\${APP_NAME}\卸载.lnk"
  RMDir "$SMPROGRAMS\${APP_NAME}"
  Delete "$DESKTOP\${APP_NAME}.lnk"

  ; 删除注册表
  DeleteRegKey HKLM "${UNINST_KEY}"
  DeleteRegKey HKLM "${INSTDIR_KEY}"
SectionEnd

; 安装前检查
Function .onInit
  ; 检查是否已安装
  ReadRegStr $0 HKLM "${INSTDIR_KEY}" "Install_Dir"
  ${If} $0 != ""
    MessageBox MB_YESNO|MB_ICONQUESTION \
      "检测到已安装 ${APP_NAME}，是否覆盖安装？" \
      IDYES proceed
    Abort
  ${EndIf}
proceed:
FunctionEnd

; 安装完成后运行
Function .onInstSuccess
  MessageBox MB_YESNO|MB_ICONQUESTION \
    "${APP_NAME} 安装完成！是否立即运行？" \
    IDYES launch
  Return
launch:
  Exec '"$INSTDIR\XTool.exe"'
FunctionEnd
