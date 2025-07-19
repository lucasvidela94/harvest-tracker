#!/bin/bash
# Harvest Scripts - Universal Installation
# This script sets up the Harvest task tracking system for any user

echo "🌾 Harvest Task Tracker - Universal Installation"
echo "================================================"

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
echo "📁 Scripts location: $SCRIPT_DIR"

# Check if Python3 is available
if ! command -v python3 &> /dev/null; then
    echo "❌ Error: Python3 is required but not installed"
    echo "Please install Python3 before continuing"
    exit 1
fi

echo "✅ Python3 found: $(python3 --version)"

# Create data directory in user's home
DATA_DIR="$HOME/.harvest"
mkdir -p "$DATA_DIR"
echo "📁 Data directory: $DATA_DIR"

# Create initial configuration
CONFIG_FILE="$DATA_DIR/config.json"
if [ ! -f "$CONFIG_FILE" ]; then
    cat > "$CONFIG_FILE" << EOF
{
    "daily_hours_target": 8.0,
    "daily_standup_hours": 0.25,
    "data_file": "$DATA_DIR/tasks.json",
    "user_name": "$USER",
    "company": "",
    "timezone": "$(timedatectl show --property=Timezone --value 2>/dev/null || echo 'UTC')"
}
EOF
    echo "✅ Configuration created: $CONFIG_FILE"
else
    echo "ℹ️  Configuration already exists: $CONFIG_FILE"
fi

# Create initial tasks file if it doesn't exist
TASKS_FILE="$DATA_DIR/tasks.json"
if [ ! -f "$TASKS_FILE" ]; then
    echo "[]" > "$TASKS_FILE"
    echo "✅ Tasks file created: $TASKS_FILE"
else
    echo "ℹ️  Tasks file already exists: $TASKS_FILE"
fi

# Install pyperclip if available
if command -v pip3 &> /dev/null; then
    echo "📦 Installing pyperclip for automatic clipboard copying..."
    pip3 install --user pyperclip > /dev/null 2>&1
    if [ $? -eq 0 ]; then
        echo "✅ pyperclip installed"
    else
        echo "⚠️  Could not install pyperclip. Manual copying will be required."
    fi
else
    echo "⚠️  pip3 not found. Install pyperclip manually: pip3 install pyperclip"
fi

# Detect shell and configure aliases
echo ""
echo "🔧 Configuring aliases..."

# Detect current shell
CURRENT_SHELL=""
SHELL_CONFIG=""

# Check if we're in zsh (more reliable detection)
if [ -n "$ZSH_VERSION" ] || [ "$SHELL" = "/bin/zsh" ] || [ "$SHELL" = "/usr/bin/zsh" ]; then
    CURRENT_SHELL="zsh"
    SHELL_CONFIG="$HOME/.zshrc"
elif [ -n "$BASH_VERSION" ] || [ "$SHELL" = "/bin/bash" ] || [ "$SHELL" = "/usr/bin/bash" ]; then
    CURRENT_SHELL="bash"
    SHELL_CONFIG="$HOME/.bashrc"
else
    CURRENT_SHELL="unknown"
    SHELL_CONFIG="$HOME/.profile"
fi

echo "🐚 Detected shell: $CURRENT_SHELL"
echo "📝 Config file: $SHELL_CONFIG"

# Create aliases for current session
alias harvest="$SCRIPT_DIR/harvest"
alias finish="$SCRIPT_DIR/finish"
alias week="$SCRIPT_DIR/week"

echo "✅ Aliases created for current session:"
echo "  harvest - Add tasks quickly"
echo "  finish  - Complete your day"
echo "  week    - Generate weekly report"

# Add aliases to shell config file
if [ -f "$SHELL_CONFIG" ]; then
    # Check if aliases already exist
    if ! grep -q "alias harvest=" "$SHELL_CONFIG"; then
        echo ""
        echo "📝 Adding aliases to $SHELL_CONFIG..."
        
        # Add a section header if it doesn't exist
        if ! grep -q "Harvest Task Tracker" "$SHELL_CONFIG"; then
            echo "" >> "$SHELL_CONFIG"
            echo "# Harvest Task Tracker" >> "$SHELL_CONFIG"
        fi
        
        # Add aliases
        echo "alias harvest='$SCRIPT_DIR/harvest'" >> "$SHELL_CONFIG"
        echo "alias finish='$SCRIPT_DIR/finish'" >> "$SHELL_CONFIG"
        echo "alias week='$SCRIPT_DIR/week'" >> "$SHELL_CONFIG"
        
        echo "✅ Aliases added to $SHELL_CONFIG"
        echo "🔄 Reloading shell configuration..."
        
        # Reload the shell config
        if [ "$CURRENT_SHELL" = "zsh" ]; then
            source "$SHELL_CONFIG"
        elif [ "$CURRENT_SHELL" = "bash" ]; then
            source "$SHELL_CONFIG"
        fi
        
        echo "✅ Shell configuration reloaded"
    else
        echo "ℹ️  Aliases already exist in $SHELL_CONFIG"
    fi
else
    echo "⚠️  Shell config file not found: $SHELL_CONFIG"
    echo "💡 Please manually add these lines to your shell config:"
    echo ""
    echo "# Harvest Task Tracker"
    echo "alias harvest='$SCRIPT_DIR/harvest'"
    echo "alias finish='$SCRIPT_DIR/finish'"
    echo "alias week='$SCRIPT_DIR/week'"
fi

# Create a quick start guide
QUICK_START="$DATA_DIR/quick_start.txt"
cat > "$QUICK_START" << 'EOF'
🌾 HARVEST QUICK START GUIDE

First time setup:
1. Run: harvest daily
2. Add tasks: harvest tech "Task description" 2.0
3. Check status: harvest status
4. Generate report: harvest report

Common commands:
- harvest daily                    # Add daily standup
- harvest tech "Fix bug" 2.0       # Add technical task
- harvest meeting "Sync" 1.0       # Add meeting
- harvest qa "Testing" 1.5         # Add QA task
- harvest status                   # Show today's status
- harvest report                   # Generate report for Harvest
- finish                          # Complete your day
- week                            # Show weekly report

Tips:
- Use quotes for task descriptions with spaces
- Hours can be: 2.0, "two hours", "half", "1.5 hours"
- Reports are automatically copied to clipboard
- Data is stored in ~/.harvest/tasks.json

Need help? Run: harvest
EOF

echo ""
echo "📚 Quick start guide created: $QUICK_START"
echo ""
echo "🎉 Installation complete! You can now use:"
echo "  harvest daily"
echo "  harvest tech 'Your first task' 2.0"
echo ""
echo "📖 For more info, see: $QUICK_START" 