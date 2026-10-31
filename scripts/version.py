#!/usr/bin/env python3
"""Single entry point for the app version (semantic version X.Y.Z).

Usage:
  version.py get                    print the current version
  version.py bump [patch|minor|major]
                                    increment it everywhere and print it
  version.py set X.Y.Z              set it everywhere and print it

The version lives in several files (Wails config, iOS and Android projects,
the Go User-Agent and the "Más" screen); this script keeps them in sync.
build/config.yml is the source of truth for `get`. Android's versionCode
(an integer that must grow with every Play upload) is incremented on every
change. Desktop scaffolding (build/darwin, build/windows) is not touched.
"""
import pathlib
import re
import sys

ROOT = pathlib.Path(__file__).resolve().parent.parent
SEMVER = r"\d+\.\d+\.\d+"

# (file, pattern with one group around the version, replacement template)
TARGETS = [
    ("build/config.yml", r'(?m)^(\s+version: ")(' + SEMVER + r')(")', r"\g<1>{v}\g<3>"),
    ("build/ios/Info.plist", r"(<key>CFBundleShortVersionString</key>\s*<string>)(" + SEMVER + r")(</string>)", r"\g<1>{v}\g<3>"),
    ("build/ios/Info.plist", r"(<key>CFBundleVersion</key>\s*<string>)(" + SEMVER + r")(</string>)", r"\g<1>{v}\g<3>"),
    ("build/ios/Info.dev.plist", r"(<key>CFBundleShortVersionString</key>\s*<string>)(" + SEMVER + r")(-dev</string>)", r"\g<1>{v}\g<3>"),
    ("build/ios/Info.dev.plist", r"(<key>CFBundleVersion</key>\s*<string>)(" + SEMVER + r")(</string>)", r"\g<1>{v}\g<3>"),
    ("build/android/app/build.gradle", r'(versionName ")(' + SEMVER + r')(")', r"\g<1>{v}\g<3>"),
    ("main.go", r'(const appVersion = ")(' + SEMVER + r')(")', r"\g<1>{v}\g<3>"),
    ("frontend/src/main.ts", r'(const APP_VERSION = ")(' + SEMVER + r')(")', r"\g<1>{v}\g<3>"),
]


def current():
    text = (ROOT / "build/config.yml").read_text()
    m = re.search(r'(?m)^\s+version: "(' + SEMVER + r')"', text)
    if not m:
        sys.exit("version not found in build/config.yml")
    return m.group(1)


def write(version):
    if not re.fullmatch(SEMVER, version):
        sys.exit(f"not a X.Y.Z version: {version}")
    for rel, pattern, repl in TARGETS:
        path = ROOT / rel
        text = path.read_text()
        new, n = re.subn(pattern, repl.format(v=version), text)
        if n == 0:
            sys.exit(f"version pattern not found in {rel}")
        path.write_text(new)
    gradle = ROOT / "build/android/app/build.gradle"
    text = gradle.read_text()
    text = re.sub(r"(versionCode )(\d+)", lambda m: f"{m.group(1)}{int(m.group(2)) + 1}", text, count=1)
    gradle.write_text(text)
    return version


def bump(part):
    major, minor, patch = (int(x) for x in current().split("."))
    if part == "major":
        major, minor, patch = major + 1, 0, 0
    elif part == "minor":
        minor, patch = minor + 1, 0
    elif part == "patch":
        patch += 1
    else:
        sys.exit(f"unknown part: {part}")
    return f"{major}.{minor}.{patch}"


if __name__ == "__main__":
    cmd = sys.argv[1] if len(sys.argv) > 1 else ""
    if cmd == "get":
        print(current())
    elif cmd == "bump":
        print(write(bump(sys.argv[2] if len(sys.argv) > 2 else "patch")))
    elif cmd == "set" and len(sys.argv) > 2:
        print(write(sys.argv[2]))
    else:
        sys.exit(__doc__)
