#!/usr/bin/env python3
import os
import sys
import argparse
import subprocess
import tempfile
import shutil
import re
import difflib

def parse_args():
    parser = argparse.ArgumentParser(description="Compare target repository with golang-starter template.")
    parser.add_argument("--template-dir", help="Path to local template directory (if omitted, clones from GitHub)")
    parser.add_argument("--target-dir", default=os.getcwd(), help="Path to target directory (default: current working directory)")
    return parser.parse_args()

def get_git_remote_info(target_dir):
    try:
        url = subprocess.check_output(
            ["git", "remote", "get-url", "origin"],
            cwd=target_dir,
            stderr=subprocess.DEVNULL
        ).decode("utf-8").strip()
        # Parse github.com/owner/project or git@github.com:owner/project
        match = re.search(r"github\.com[:/]([^/]+)/([^/.]+)(?:\.git)?", url)
        if match:
            return match.group(1), match.group(2)
    except Exception:
        pass
    return None, None

def get_go_module_info(target_dir):
    go_mod_path = os.path.join(target_dir, "go.mod")
    if os.path.exists(go_mod_path):
        try:
            with open(go_mod_path, "r", encoding="utf-8") as f:
                content = f.read()
            match = re.search(r"module\s+([^\s\n]+)", content)
            if match:
                module_path = match.group(1)
                parts = module_path.split("/")
                if len(parts) >= 3:
                    # github.com/owner/project
                    return parts[-2], parts[-1]
                elif len(parts) == 2:
                    return parts[0], parts[1]
                else:
                    return None, parts[0]
        except Exception:
            pass
    return None, None

def detect_project_info(target_dir):
    # Try go.mod first, then git remote, then fallback to directory name
    owner, project = get_go_module_info(target_dir)
    if not owner or not project:
        r_owner, r_project = get_git_remote_info(target_dir)
        if r_owner:
            owner = owner or r_owner
        if r_project:
            project = project or r_project

    if not owner:
        owner = "toozej"
    if not project:
        project = os.path.basename(os.path.abspath(target_dir))

    return owner, project

def clone_template():
    temp_dir = tempfile.mkdtemp(prefix="golang-starter-template-")
    print(f"Cloning template repository from GitHub to temporary folder: {temp_dir}...")
    try:
        subprocess.check_call(
            ["git", "clone", "--depth", "1", "https://github.com/toozej/golang-starter.git", temp_dir],
            stdout=subprocess.DEVNULL,
            stderr=subprocess.DEVNULL
        )
        return temp_dir
    except subprocess.CalledProcessError as e:
        shutil.rmtree(temp_dir)
        print(f"Error: Failed to clone template repo: {e}", file=sys.stderr)
        sys.exit(1)

def customize_content(content, target_owner, target_project):
    # Replace golang-starter with target project name
    content = content.replace("golang-starter", target_project)
    # Replace toozej with target owner
    content = content.replace("toozej", target_owner)
    return content

def get_scaffolding_files():
    return [
        ".github/workflows/ci.yaml",
        ".github/workflows/release.yaml",
        ".github/workflows/weekly-docker-refresh.yaml",
        "Makefile",
        ".goreleaser.yml",
        ".gitignore",
        ".dockerignore",
        ".pre-commit-config.yaml",
        "checkmake.ini",
        ".air.toml",
        "docker-compose.yml",
        "pkg/version/version.go",
        "pkg/version/version_test.go",
        "pkg/man/man.go",
        "pkg/man/man_test.go",
        "pkg/config/config.go",
        "pkg/config/config_test.go"
    ]

def find_target_dockerfiles(target_dir):
    dockerfiles = []
    exclude_dirs = {".git", "vendor", "node_modules", "dist", "out", "tmp"}
    for root, dirs, files in os.walk(target_dir):
        # Filter directories in-place
        dirs[:] = [d for d in dirs if d not in exclude_dirs]
        for f in files:
            if "dockerfile" in f.lower():
                full_path = os.path.join(root, f)
                rel_path = os.path.relpath(full_path, target_dir)
                dockerfiles.append(rel_path)
    return dockerfiles

def get_template_dockerfiles(template_dir):
    files = ["Dockerfile", "Dockerfile.distroless", "Dockerfile.goreleaser", "Dockerfile.goreleaser.distroless"]
    valid_files = []
    for f in files:
        if os.path.exists(os.path.join(template_dir, f)):
            valid_files.append(f)
    return valid_files

def main():
    args = parse_args()

    target_dir = os.path.abspath(args.target_dir)
    if not os.path.isdir(target_dir):
        print(f"Error: Target directory does not exist: {target_dir}", file=sys.stderr)
        sys.exit(1)

    owner, project = detect_project_info(target_dir)
    print(f"Detected target project configuration:")
    print(f"  Owner/Username: {owner}")
    print(f"  Project Name:   {project}")
    print(f"  Target Dir:     {target_dir}\n")

    temp_template_dir = None
    if args.template_dir:
        template_dir = os.path.abspath(args.template_dir)
        if not os.path.isdir(template_dir):
            print(f"Error: Local template directory does not exist: {template_dir}", file=sys.stderr)
            sys.exit(1)
    else:
        temp_template_dir = clone_template()
        template_dir = temp_template_dir

    try:
        # 1. Compare standard scaffolding files
        print("=" * 60)
        print("Comparing Scaffolding & Configuration Files")
        print("=" * 60)

        scaffolding = get_scaffolding_files()
        for rel_file in scaffolding:
            template_path = os.path.join(template_dir, rel_file)
            target_path = os.path.join(target_dir, rel_file)

            if not os.path.exists(template_path):
                continue

            # Read template
            with open(template_path, "r", encoding="utf-8", errors="ignore") as f:
                template_raw = f.read()
            template_cust = customize_content(template_raw, owner, project)

            if not os.path.exists(target_path):
                print(f"\n[ADD] {rel_file}")
                print(f"  File exists in template but is missing in target.")
                print(f"  You may want to create it with the following customized content:")
                print("-" * 40)
                # Print first 20 lines as preview
                lines = template_cust.splitlines()
                for line in lines[:25]:
                    print(line)
                if len(lines) > 25:
                    print(f"... ({len(lines) - 25} more lines)")
                print("-" * 40)
                continue

            # Read target
            with open(target_path, "r", encoding="utf-8", errors="ignore") as f:
                target_content = f.read()

            if template_cust.strip() == target_content.strip():
                print(f"[MATCHED] {rel_file}")
            else:
                print(f"\n[MODIFIED] {rel_file}")
                diff = difflib.unified_diff(
                    target_content.splitlines(),
                    template_cust.splitlines(),
                    fromfile=f"target/{rel_file}",
                    tofile=f"template/{rel_file} (customized)",
                    lineterm=""
                )
                diff_list = list(diff)
                if diff_list:
                    print("\n".join(diff_list))
                else:
                    print("  (Only whitespace/ending differences)")

        # 2. Compare and match Dockerfiles by content similarity
        print("\n" + "=" * 60)
        print("Matching and Comparing Dockerfiles (by content similarity)")
        print("=" * 60)

        target_dockerfiles = find_target_dockerfiles(target_dir)
        template_dockerfiles = get_template_dockerfiles(template_dir)

        if not target_dockerfiles:
            print("No Dockerfiles found in target repo.")
        else:
            for target_rel in target_dockerfiles:
                target_path = os.path.join(target_dir, target_rel)
                with open(target_path, "r", encoding="utf-8", errors="ignore") as f:
                    target_content = f.read()
                
                best_match_file = None
                best_match_ratio = -1.0
                best_match_cust_content = None

                # Check similarity against each customized template Dockerfile
                for temp_file in template_dockerfiles:
                    temp_path = os.path.join(template_dir, temp_file)
                    with open(temp_path, "r", encoding="utf-8", errors="ignore") as f:
                        temp_raw = f.read()
                    temp_cust = customize_content(temp_raw, owner, project)

                    # Calculate similarity ratio
                    ratio = difflib.SequenceMatcher(None, target_content, temp_cust).ratio()
                    if ratio > best_match_ratio:
                        best_match_ratio = ratio
                        best_match_file = temp_file
                        best_match_cust_content = temp_cust

                print(f"\nTarget File:   {target_rel}")
                if best_match_file:
                    print(f"Closest Match: {best_match_file} (Similarity: {best_match_ratio:.2%})")
                    if best_match_ratio >= 0.999:
                        print(f"[MATCHED] {target_rel} is identical to {best_match_file}")
                    else:
                        print(f"Diff against {best_match_file} (customized):")
                        diff = difflib.unified_diff(
                            target_content.splitlines(),
                            best_match_cust_content.splitlines(),
                            fromfile=f"target/{target_rel}",
                            tofile=f"template/{best_match_file} (customized)",
                            lineterm=""
                        )
                        diff_list = list(diff)
                        if diff_list:
                            print("\n".join(diff_list))
                        else:
                            print("  (Only whitespace/ending differences)")
                else:
                    print("Could not find any matching template Dockerfile.")

        print("\n" + "=" * 60)
        print("Comparison Completed.")
        print("=" * 60)
        print("Instructions: Please use these diffs to carefully back-fill updates.")
        print("Remember to NOT overwrite target Dockerfiles in their entirety, but merge updates.")

    finally:
        if temp_template_dir:
            shutil.rmtree(temp_template_dir)

if __name__ == "__main__":
    main()
