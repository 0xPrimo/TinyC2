# TinyC2
TinyC2 is a small C2 Framework developed with the goal of reimplementing Havoc Pro Runtime Channel Switching and Cobalt Strike UDC2 features.

# Disclaimer

> This repository is intended for authorized, educational, or ethical security research purposes only. The tools and techniques within are not to be used against systems without explicit permission. The author assumes no liability for any unauthorized or illegal misuse of this project.

# Install

1. Install dependencies
```bash
sudo apt update
sudo apt install golang cmake make mingw-w64 openjdk-11-jdk zip
```
2. Clone the project
```bash
git clone --recurse-submodules git@github.com:0xPrimo/TinyC2.git
cd tinyc2
```

3. Download crystal palace project from [here](https://tradecraftgarden.org/crystalpalace.html) then run Install.sh:
```bash
sudo ./install.sh /path/to/dist
```

4. Build project
```
make
```


# Usage

```
./tinyc2 server/config.yaml
```

1. After you build the plugin you can register it with:
```
tinyc2> plugin [plugin name] [/path/to/plugin.so]
```

2. Create listener, and make sure the config is correct:
```
tinyc2> listener start [plugin name] [listener name] [/path/to/config.yaml]
```

3. Now generate the implant payload: 
```
tinyc2> listener generate [listener name]
```


https://github.com/user-attachments/assets/28a366fc-15a3-4a91-81ea-e947e300ff35

# Credit
- [A PIC Security Research Adventure - Tradecraft Garden](https://tradecraftgarden.org/) 
- [Havoc Professional Release](https://www.infinitycurve.org/blog/release)
- [@RastaMouse Crystal Palace series & crystal-kit repo](https://github.com/rasta-mouse/Crystal-Kit)
- [Cobalt Strike 4.12: Fix Up, Look Sharp!](https://www.cobaltstrike.com/blog/cobalt-strike-412-fix-up-look-sharp)
