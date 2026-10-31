#!/usr/bin/env python3
"""Development signing helper for `wails3 task device:ios`.

Usage:
  ios_signing.py identity <profile.mobileprovision>
      Prints the SHA-1 of the keychain code-signing identity whose
      certificate is included in the provisioning profile.
  ios_signing.py entitlements <profile.mobileprovision> <bundle id> <out.plist>
      Writes development entitlements derived from the profile (App ID
      prefix, team id, keychain group, get-task-allow).
"""
import hashlib
import plistlib
import re
import subprocess
import sys


def load_profile(path):
    raw = subprocess.run(["security", "cms", "-D", "-i", path], capture_output=True, check=True).stdout
    return plistlib.loads(raw)


def identity(profile_path):
    profile = load_profile(profile_path)
    allowed = {hashlib.sha1(c).hexdigest().upper() for c in profile.get("DeveloperCertificates", [])}
    out = subprocess.run(["security", "find-identity", "-v", "-p", "codesigning"], capture_output=True, text=True).stdout
    for sha in re.findall(r"\b([0-9A-F]{40})\b", out):
        if sha in allowed:
            return sha
    return ""


def entitlements(profile_path, bundle_id, out_path):
    profile = load_profile(profile_path)
    prefix = profile["ApplicationIdentifierPrefix"][0]
    team = profile["TeamIdentifier"][0]
    app_id = f"{prefix}.{bundle_id}"
    ents = {
        "application-identifier": app_id,
        "com.apple.developer.team-identifier": team,
        "keychain-access-groups": [app_id],
        "get-task-allow": True,
    }
    with open(out_path, "wb") as f:
        plistlib.dump(ents, f)


if __name__ == "__main__":
    cmd = sys.argv[1] if len(sys.argv) > 1 else ""
    try:
        if cmd == "identity":
            print(identity(sys.argv[2]))
        elif cmd == "entitlements":
            entitlements(sys.argv[2], sys.argv[3], sys.argv[4])
        else:
            sys.exit(__doc__)
    except (FileNotFoundError, subprocess.CalledProcessError, KeyError, IndexError):
        # Missing/invalid profile: print nothing, the task explains what to do.
        pass
