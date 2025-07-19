#!/bin/bash
# Harvest Scripts - Uninstall
# This script removes Harvest aliases and configuration

echo "🗑️  Harvest Task Tracker - Uninstall"
echo "===================================="

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Detect current shell
CURRENT_SHELL=""
SHELL_CONFIG=""

if [ -n "$ZSH_VERSION" ]; then
    CURRENT_SHELL="zsh"
    SHELL_CONFIG="$HOME/.zshrc"
elif [ -n "$BASH_VERSION" ]; then
    CURRENT_SHELL="bash"
    SHELL_CONFIG="$HOME/.bashrc"
else
    CURRENT_SHELL="unknown"
    SHELL_CONFIG="$HOME/.profile"
fi

echo "🐚 Detected shell: $CURRENT_SHELL"
echo "📝 Config file: $SHELL_CONFIG"

# Remove aliases from shell config
if [ -f "$SHELL_CONFIG" ]; then
    echo ""
    echo "🧹 Removing aliases from $SHELL_CONFIG..."
    
    # Create a temporary file without the Harvest aliases
    TEMP_FILE=$(mktemp)
    
    # Remove the Harvest section and aliases
    awk '
    /^# Harvest Task Tracker$/ {
        in_harvest_section = 1
        next
    }
    /^alias harvest=/ || /^alias finish=/ || /^alias week=/ {
        if (in_harvest_section) {
            next
        }
    }
    /^[^#]/ && in_harvest_section {
        in_harvest_section = 0
    }
    {
        if (!in_harvest_section) {
            print
        }
    }
    ' "$SHELL_CONFIG" > "$TEMP_FILE"
    
    # Replace the original file
    mv "$TEMP_FILE" "$SHELL_CONFIG"
    
    echo "✅ Aliases removed from $SHELL_CONFIG"
else
    echo "⚠️  Shell config file not found: $SHELL_CONFIG"
fi

# Remove aliases from current session
echo ""
echo "🔄 Removing aliases from current session..."
unalias harvest 2>/dev/null || true
unalias finish 2>/dev/null || true
unalias week 2>/dev/null || true

echo "✅ Aliases removed from current session"

# Ask if user wants to remove data
echo ""
read -p "🗑️  Do you want to remove all Harvest data? (y/N): " -n 1 -r
echo

if [[ $REPLY =~ ^[Yy]$ ]]; then
    DATA_DIR="$HOME/.harvest"
    if [ -d "$DATA_DIR" ]; then
        echo "🗑️  Removing data directory: $DATA_DIR"
        rm -rf "$DATA_DIR"
        echo "✅ Data directory removed"
    else
        echo "ℹ️  Data directory not found"
    fi
else
    echo "ℹ️  Data directory preserved at ~/.harvest"
fi

echo ""
echo "🎉 Uninstall complete!"
echo "💡 You may need to restart your terminal or run 'source $SHELL_CONFIG'" 