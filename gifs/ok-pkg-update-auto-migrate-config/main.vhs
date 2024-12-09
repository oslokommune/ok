Require ok

Output main.gif
Output main.webm

Set Width 1500
Set Height 600

Set FontSize 15

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
Sleep 10s

Type "# Done!"

Sleep 3s
