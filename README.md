# TinyC2
TinyC2 is a small C2 Framework developed with the goal of reimplementing Havoc Pro Runtime Channel Switching and Cobalt Strike UDC2 features.

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

# Technical Details

In the recent releases of Havoc Pro and Cobalt Strike, a few new features caught my attention and found them really cool. The first one was Havoc Pro's Runtime Channel Switching. It's built into the Kaine-kit extension and allows the operator to reconfigure implant listener configuration or switch between different listeners at runtime.

> *A feature has been introduced to kaine that allows operators to dynamically update a session's listener configuration at runtime. Additionally, **sessions can switch between entirely different listener protocols at runtime**, such as transitioning from HTTP to SMB or vice versa, regardless of whether the listener is configured as P2P or Direct.*
> — *[Havoc Professional Release](https://infinitycurve.org/blog/release)*

The second one was Cobalt Strike's UDC2 (User-Defined Command and Control). It allows the operator to develop custom C2 channels by implementing their communication logic directly into a BOF (Beacon Object File). This BOF is then patched into beacon payload.

> … UDC2 addresses these issues by enabling Cobalt Strike users to **develop custom C2 channels as BOFs**. The UDC2 BOF is patched in on payload creation and is invoked by Beacon to proxy out all its traffic over the custom channel. This makes it possible to **combine custom C2 channels with custom UDRLs/transformations**. 
> — *[Cobalt Strike 4.12: Fix Up, Look Sharp!](https://www.cobaltstrike.com/blog/cobalt-strike-412-fix-up-look-sharp)*

## Architecture
To see how TinyC2 works, we first need to look at its core. TinyC2 is made up of three main components:

- **Plugin Management:** TinyC2 supports dynamically loading Go plugins at runtime to extend TinyC2 functionality. Currently, TinyC2 only support listener plugins.
- **Listener Management:** This component handles the listeners, allowing the user to start or stop them, as well as generate implant payloads from specific listener.
- **Implant Management:** This component  responsible for managing implant sessions for example request implant to execute command or terminating the session.

<img width="1220" height="773" alt="image" src="https://github.com/user-attachments/assets/c42e02c5-bff4-4ff1-b2c5-df4c18b7673e" />

## Custom Listeners
In order to add a custom listener to TinyC2 the user should create a go plugin, think of it as a `DLL` that the server can load at runtime. this plugin should have one export with name `NewListener` and must return an object that implements `IListener` interface.

After writing your go plugin you will also need to implement the PIC which will be built with Crystal Palace linker.

This PIC will be used in two cases:

- When listener manager wants to generate implant payload.
- When implant manager wants to register a new channel to target implant payload at runtime.

See example.

<img width="1488" height="1016" alt="image" src="https://github.com/user-attachments/assets/de33b941-5e07-4258-813f-3d481cef9ddd" />

# Todo
- more listener plugins
- more implant commands
- support x86
- support agent plugin

# Credit
- [A PIC Security Research Adventure - Tradecraft Garden](https://tradecraftgarden.org/) 
- [Havoc Professional Release](https://www.infinitycurve.org/blog/release)
- [@RastaMouse Crystal Palace series & crystal-kit repo](https://github.com/rasta-mouse/Crystal-Kit)
- [Cobalt Strike 4.12: Fix Up, Look Sharp!](https://www.cobaltstrike.com/blog/cobalt-strike-412-fix-up-look-sharp)
