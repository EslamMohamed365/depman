Output demo.gif

# Settings
Set FontSize 14
Set Width 800
Set Height 600
Set Padding 20
Set WindowBar Colorful
Set TypingSpeed 50ms
Set Theme "TokyoNight"

# Start fresh terminal
Hide
Type "cd /tmp"
Enter
Type "rm -rf depman_demo"
Enter
Type "mkdir depman_demo"
Enter
Type "cd depman_demo"
Enter
Type "echo '[project]' > pyproject.toml"
Enter
Show

# Start depman
Type "/home/eslam/coding/projects/depman/depman"
Enter

# Wait for TUI to load
Sleep 2s

# Try to see help
Type "?"
Sleep 2s

# Quit
Type "q"
Sleep 500ms

# Cleanup
Hide
Type "cd /tmp"
Enter
Type "rm -rf depman_demo"
Enter
