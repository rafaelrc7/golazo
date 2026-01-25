# Notification

Goal notifications require one-time setup depending on your operating system.

## macOS

Notifications use AppleScript, which requires enabling notifications for Script Editor:

1. Open **Script Editor** (`/Applications/Utilities/Script Editor.app`)
2. Paste and run: `display notification "test" with title "test"`
3. Open **System Settings → Notifications → Script Editor**
4. Enable/Allow notifications and set alert style to "Banners"

## Linux

Notifications require `libnotify`. Install if not present:

```bash
# Debian/Ubuntu
sudo apt install libnotify-bin

# Fedora
sudo dnf install libnotify

# Arch
sudo pacman -S libnotify
```

## Windows

Notifications should work out-of-box on Windows 10/11.