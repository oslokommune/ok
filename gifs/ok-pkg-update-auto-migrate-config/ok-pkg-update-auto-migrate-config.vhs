Require ok

Output ok-pkg-update-auto-migrate-config.webm

Set FontSize 23
Set Width 1600
Set Height 900
Set Framerate 24

Set TypingSpeed 100ms

# Source update.tape

Sleep 3s

Type "ok pkg update"
Sleep 3s
Enter
Sleep 10s

Type "git diff packages.yml"
Sleep 5s
Enter
Sleep 10s

Type "git diff config/app-hello.yml"
Sleep 5s
Enter
Sleep 5s

Type "# Done!"

Sleep 3s
