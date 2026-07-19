Set WshShell = CreateObject("WScript.Shell")
WshShell.CurrentDirectory = "C:\Users\hyper\workspace\tormentnexus\apps\web"
WshShell.Environment("Process").Item("PORT") = "7779"
WshShell.Run "cmd /c node .next-build\standalone\apps\web\server.js", 0, False
