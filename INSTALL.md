# Installing otdfctl

This guide will help you install **otdfctl** on your computer. Follow the steps for your operating system.

## What is otdfctl?

`otdfctl` is a command-line tool for working with OpenTDF. Once installed, you'll be able to run it from any folder on your computer.

---

## 📥 Installation Instructions

### For macOS Users

#### Step 1: Download the Right Version

First, you need to know which type of Mac you have:
- **Apple Silicon (M1, M2, M3, etc.)**: Download the `arm64` version
- **Intel Mac**: Download the `amd64` version

**Don't know which one you have?** Open Terminal and type:
```bash
uname -m
```
- If it says `arm64`, you have Apple Silicon
- If it says `x86_64`, you have an Intel Mac

#### Step 2: Download and Install

**For Apple Silicon Macs (M1/M2/M3):**
```bash
# Download the latest version
curl -LO https://github.com/opentdf/otdfctl/releases/download/v0.28.0/otdfctl-0.28.0-darwin-arm64.tar.gz

# Extract the file
tar -xzf otdfctl-0.28.0-darwin-arm64.tar.gz

# Move it to a folder in your PATH
sudo mv otdfctl /usr/local/bin/

# Make it executable
sudo chmod +x /usr/local/bin/otdfctl

# Clean up the downloaded file
rm otdfctl-0.28.0-darwin-arm64.tar.gz
```

**For Intel Macs:**
```bash
# Download the latest version
curl -LO https://github.com/opentdf/otdfctl/releases/download/v0.28.0/otdfctl-0.28.0-darwin-amd64.tar.gz

# Extract the file
tar -xzf otdfctl-0.28.0-darwin-amd64.tar.gz

# Move it to a folder in your PATH
sudo mv otdfctl /usr/local/bin/

# Make it executable
sudo chmod +x /usr/local/bin/otdfctl

# Clean up the downloaded file
rm otdfctl-0.28.0-darwin-amd64.tar.gz
```

> **Note:** You'll be asked for your password when using `sudo`. This is normal and required to install the tool.

---

### For Linux Users

#### Step 1: Download the Right Version

Choose based on your system:
- **Most modern PCs**: `amd64`
- **Raspberry Pi or ARM devices**: `arm64` or `arm`

**Not sure?** Run this command:
```bash
uname -m
```

#### Step 2: Download and Install

**For amd64 (most common):**
```bash
# Download the latest version
curl -LO https://github.com/opentdf/otdfctl/releases/download/v0.28.0/otdfctl-0.28.0-linux-amd64.tar.gz

# Extract the file
tar -xzf otdfctl-0.28.0-linux-amd64.tar.gz

# Move it to a folder in your PATH
sudo mv otdfctl /usr/local/bin/

# Make it executable
sudo chmod +x /usr/local/bin/otdfctl

# Clean up the downloaded file
rm otdfctl-0.28.0-linux-amd64.tar.gz
```

**For ARM64:**
```bash
# Download the latest version
curl -LO https://github.com/opentdf/otdfctl/releases/download/v0.28.0/otdfctl-0.28.0-linux-arm64.tar.gz

# Extract the file
tar -xzf otdfctl-0.28.0-linux-arm64.tar.gz

# Move it to a folder in your PATH
sudo mv otdfctl /usr/local/bin/

# Make it executable
sudo chmod +x /usr/local/bin/otdfctl

# Clean up the downloaded file
rm otdfctl-0.28.0-linux-arm64.tar.gz
```

**For ARM (32-bit):**
```bash
# Download the latest version
curl -LO https://github.com/opentdf/otdfctl/releases/download/v0.28.0/otdfctl-0.28.0-linux-arm.tar.gz

# Extract the file
tar -xzf otdfctl-0.28.0-linux-arm.tar.gz

# Move it to a folder in your PATH
sudo mv otdfctl /usr/local/bin/

# Make it executable
sudo chmod +x /usr/local/bin/otdfctl

# Clean up the downloaded file
rm otdfctl-0.28.0-linux-arm.tar.gz
```

---

### For Windows Users

#### Step 1: Download the Right Version

Most Windows computers use `amd64`. ARM versions are for Surface Pro X or other ARM-based Windows devices.

#### Step 2: Download and Install

1. **Download the file** for your system:
   - For most PCs: [otdfctl-0.28.0-windows-amd64.zip](https://github.com/opentdf/otdfctl/releases/download/v0.28.0/otdfctl-0.28.0-windows-amd64.zip)
   - For ARM devices: [otdfctl-0.28.0-windows-arm64.zip](https://github.com/opentdf/otdfctl/releases/download/v0.28.0/otdfctl-0.28.0-windows-arm64.zip)

2. **Extract the ZIP file:**
   - Right-click the downloaded file
   - Select "Extract All..."
   - Choose a location (e.g., `C:\Program Files\otdfctl`)

3. **Add to PATH:**
   - Press `Windows key` and search for "Environment Variables"
   - Click "Edit the system environment variables"
   - Click "Environment Variables..." button
   - Under "System variables", find and select "Path", then click "Edit..."
   - Click "New" and add the folder where you extracted `otdfctl.exe` (e.g., `C:\Program Files\otdfctl`)
   - Click "OK" on all windows

4. **Restart your command prompt or PowerShell** for the changes to take effect.

---

## ✅ Verify Installation

After installation, open a **new** terminal window and run:

```bash
otdfctl --version
```

You should see output like:
```
otdfctl v0.28.0
```

If you see this, congratulations! 🎉 You've successfully installed otdfctl.

---

## 🔧 Troubleshooting

### "command not found" or "otdfctl is not recognized"

**On macOS/Linux:**
- Make sure you opened a **new** terminal window after installation
- Verify the file is in `/usr/local/bin/` by running: `ls -l /usr/local/bin/otdfctl`
- Check if `/usr/local/bin` is in your PATH: `echo $PATH`

**On Windows:**
- Make sure you opened a **new** Command Prompt or PowerShell window
- Verify the PATH was added correctly in Environment Variables
- Try searching for `otdfctl.exe` in File Explorer to confirm where it's located

### "Permission denied" on macOS/Linux

If you get a permission error, make sure the file is executable:
```bash
sudo chmod +x /usr/local/bin/otdfctl
```

### macOS: "otdfctl cannot be opened because the developer cannot be verified"

This is a security feature on macOS. To allow it:
```bash
sudo xattr -d com.apple.quarantine /usr/local/bin/otdfctl
```

---

## 📚 Next Steps

Now that otdfctl is installed, you can:
- Run `otdfctl --help` to see available commands
- Visit the [documentation](https://github.com/opentdf/otdfctl) for usage examples
- Configure your first profile with `otdfctl auth login`

---

## 🔄 Updating otdfctl

To update to a newer version, simply repeat the installation steps with the new version number. The new version will replace the old one.

To check for new releases, visit: https://github.com/opentdf/otdfctl/releases
