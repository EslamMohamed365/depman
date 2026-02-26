# Output settings
Output demo.gif
Set Shell zsh
Set FontSize 18
Set Width 1200
Set Height 800
Set Padding 40
Set WindowBar Colorful
Set TypingSpeed 100ms
Set Theme "TokyoNight"

# --- SETUP (Hidden) ---
Hide
Type "mkdir -p demo_vhs"
Enter
Type "cd demo_vhs"
Enter
Type "echo 'requests==2.28.1' > requirements.txt"
Enter
Type "echo 'flask==2.2.2' >> requirements.txt"
Enter
Type "python3 -m venv .venv"
Enter
Type "../depman --help > /dev/null" # Ensure it's built
Enter
Type "clear"
Enter
Show

# --- DEMO ---
# Start depman
Type "../depman"
Sleep 1s
Enter
Sleep 3s

# Navigate installed packages
Type "j"
Sleep 500ms
Type "j"
Sleep 500ms
Type "k"
Sleep 1s

# Search PyPI
Type "/"
Sleep 1s
Type "httpx"
Sleep 2s
Enter
Sleep 4s

# Show Help
Type "?"
Sleep 3s

# Quit
Type "q"
Sleep 1s

# --- CLEANUP (Hidden) ---
Hide
Type "cd .."
Enter
Type "rm -rf demo_vhs"
Enter
