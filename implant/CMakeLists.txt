cmake_minimum_required(VERSION 3.14)

project(Implant)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_COMPILER /usr/bin/x86_64-w64-mingw32-c++)

include_directories(
        Implant/include
        Implant/external/
)

add_executable(Implant
        Implant/src/Main.cpp

        # implant
        #
        Implant/src/Implant/Implant.cpp
        Implant/src/Implant/ImplantExecute.cpp
        Implant/src/Implant/ImplantJob.cpp
        Implant/src/Implant/ImplantTransport.cpp

        # peer-to-peer
        Implant/src/Peer/Peer.cpp
        Implant/src/Peer/Smb.cpp

        # channel
        #
        Implant/src/Channel/Channel.cpp

        # commands
        #
        Implant/src/Command/Command.cpp
        Implant/src/Command/CommandChannel.cpp
        Implant/src/Command/CommandExit.cpp
        Implant/src/Command/CommandPs.cpp
        Implant/src/Command/CommandCd.cpp
        Implant/src/Command/CommandCp.cpp
        Implant/src/Command/CommandShell.cpp
        Implant/src/Command/CommandDownload.cpp
        Implant/src/Command/CommandUpload.cpp
        Implant/src/Command/CommandRun.cpp
        Implant/src/Command/CommandExecuteAssembly.cpp
        Implant/src/Command/CommandJob.cpp
        Implant/src/Command/CommandInlineExecute.cpp
        Implant/src/Command/CommandInjectShellcode.cpp
        Implant/src/Command/CommandToken.cpp
        Implant/src/Command/CommandPivot.cpp

        # coff loader
        #
        Implant/src/Coff/BofApis.cpp
        Implant/src/Coff/CoffLoader.cpp
        Implant/src/Coff/CoffRelocations.cpp

        #
        # stdlib
        Implant/src/Rtlib/Base64.cpp
        Implant/src/Rtlib/Fs.cpp
        Implant/src/Rtlib/String.cpp
        Implant/src/Rtlib/Pipe.cpp
)

target_link_options(Implant PRIVATE
        "-Wl,--gc-sections"
        "-s"
        "-static"
        "-static-libgcc"
        "-static-libstdc++"
)

set(DEFAULT_CHANNEL "{0x0}" CACHE STRING "default channel")
set(DEFAULT_CHANNEL_CONFIG "{0x0}" CACHE STRING "default channel config")

target_compile_definitions(Implant PRIVATE "DEFAULT_CHANNEL_CONFIG=${DEFAULT_CHANNEL_CONFIG}")
target_compile_definitions(Implant PRIVATE "DEFAULT_CHANNEL=${DEFAULT_CHANNEL}")

target_link_libraries(Implant PRIVATE
        ntdll
        crypt32
        secur32
)

