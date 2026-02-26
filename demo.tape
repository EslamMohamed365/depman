Output demo.gif

# Settings
Set FontSize 14
Set Width 900
Set Height 600
Set Padding 20
Set WindowBar Colorful
Set TypingSpeed 40ms
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
Type "python3 -m venv venv"
Enter
Type "echo '[project]' > pyproject.toml"
Enter
Show

# Start depman
Type "/home/eslam/coding/projects/depman/depman"
Enter

# Wait for TUI to load
Sleep 2s

# Navigate in installed panel
Type "j"
Sleep 300ms

Type "j"
Sleep 300ms

Type "k"
Sleep 300ms

# Switch to outdated panel
Type "Tab"
Sleep 500ms

# Go back to installed
Type "Tab"
Sleep 500ms

# Search PyPI with /
Type "/"
Sleep 500ms

Type "flask"
Sleep 2s

# Select first result
Type "Enter"
Sleep 1s

# Back to dashboard
Escape
Sleep 500ms

# Show help
Type "?"
Sleep 2s

# Quit help
Escape
Sleep 500ms

# Quit app
Type "q"
Sleep 500ms

# Cleanup
Hide
Type "cd /tmp"
Enter
Type "rm -rf depman_demo"
Enter
