@echo off

if exist run.exe (
    del run.exe
    echo run.exe deletado.
)

air --build.cmd "go build -o run.exe cmd/main.go" --build.exclude_dir cli --build.bin "run.exe"