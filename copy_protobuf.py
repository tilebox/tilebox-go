#!/usr/bin/env python

"""
A small utility script to copy generated protobuf messages from the core repository into this repository

This will put them in the right place, and automatically fix the imports as well

Usage:
./copy_protobuf.py <path-to-core-repo> <path-to-clients-repo>
"""

import sys
from pathlib import Path

_SERVICE_MODULES = {
    "workflows-service": {
        "workflows": {
            "core": Path("workflows-service/protogen/go/workflows/v1"),
            "diagram": Path("workflows-service/protogen/go/workflows/v1"),
            "job": Path("workflows-service/protogen/go/workflows/v1"),
            "lease": Path("workflows-service/protogen/go/workflows/v1"),
            "task": Path("workflows-service/protogen/go/workflows/v1"),
            "workflows": Path("workflows-service/protogen/go/workflows/v1"),
        },
    },
}


def main() -> None:
    core_repo, clients_repo = Path(__file__).parent.parent / "core", Path(__file__).parent
    if len(sys.argv) == 3:  # manual arg parsing, don't want any dependencies for this simple script
        core_repo = Path(sys.argv[1])
        clients_repo = Path(sys.argv[2])
    elif len(sys.argv) != 1:
        print("Usage: ./copy_protobuf.py <core-repo> <clients-repo>")  # noqa: T201
        sys.exit(1)

    for service_folder, packages in _SERVICE_MODULES.items():
        for package, modules in packages.items():
            for module in modules:
                copy_module(core_repo, clients_repo, service_folder, package, module)


def copy_module(core_repo: Path, clients_repo: Path, service_folder: str, package: str, module: str) -> None:
    target_path = clients_repo / _SERVICE_MODULES[service_folder][package][module]
    search_path = core_repo / service_folder / "protogen" / "go" / package / "v1"
    for file in search_path.glob(f"{module}*"):
        if file.is_dir():
            for subfile in file.glob("*"):
                target_file = target_path / file.name / subfile.name
                target_file.parent.mkdir(parents=True, exist_ok=True)
                target_file.write_text(fix_imports(subfile.read_text()))
                print(f"A Copied and fixed file: {target_file}")
            continue

        target_file = target_path / file.name
        target_file.parent.mkdir(parents=True, exist_ok=True)
        target_file.write_text(fix_imports(file.read_text()))
        print(f"Copied and fixed file: {target_file}")  # noqa: T201


def fix_imports(content: str) -> str:
    for service_folder, packages in _SERVICE_MODULES.items():
        for package in packages:
            content = content.replace(
                f'v1 "github.com/tilebox/core/{service_folder}/protogen/go/{package}/v1"',
                f'v1 "github.com/tilebox/tilebox-go/{service_folder}/protogen/go/{package}/v1"',
            )
    return content.strip() + "\n"  # end of file -> new line fix


if __name__ == "__main__":
    main()
