import os
import subprocess
from contextlib import contextmanager
from pathlib import Path
from typing import Union
import glob
import yaml

from cookiecutter.utils import rmtree


@contextmanager
def cd(p: Path):
    now = Path.cwd()
    try:
        os.chdir(str(p))
        yield
    finally:
        os.chdir(str(now))


@contextmanager
def bake(cookies, *args, **kwrags):
    r = cookies.bake(*args, **kwrags)
    try:
        yield r
    finally:
        rmtree(str(r.project_path))


def run(
    cmd: Union[str, list[str]], dir: Path, *args, **kwargs
) -> subprocess.CompletedProcess:
    with cd(dir):
        print(f"[run] {cmd}")
        return subprocess.run(cmd, check=True, *args, **kwargs)


def check_result(result):
    assert result.exit_code == 0
    assert result.exception is None
    assert result.project_path.is_dir()
    project_path = result.project_path

    def do(cmd: Union[str, list[str]]):
        run(cmd, project_path)

    sequnce = [
        ["git", "init"],
        ["git", "config", "init.defaultBranch", "main"],
        ["git", "config", "user.email", "test-bake@example.com"],
        ["git", "config", "user.name", "test-bake"],
        ["make", "init"],
        ["git", "add", "-A"],
        ["git", "commit", "-m", "Init"],
        ["git", "tag", "v0.1.0"],
        ["make", "test"],
        ["make", "vet"],
        ["make", "vuln"],
        ["make"],
    ]

    for seq in sequnce:
        do(seq)

    with cd(project_path):
        print("[check] all cookiecutter has been replaced")
        r = subprocess.run(["git", "grep", "cookiecutter"], capture_output=True, text=True)
        assert len(r.stdout) == 0
        print("[check] valid yaml")
        for yaml_file in glob.glob(".github/**/*.yml", recursive=True):
            with open(yaml_file) as f:
                print(f"[yaml] {yaml_file}")
                yaml.safe_load(f)


def test_bake_and_make_default(cookies):
    with bake(cookies) as result:
        check_result(result)


def test_bake_and_make_go_command(cookies):
    context = {
        "project_category": "Command",
        "project_name": "command 120",
    }
    with bake(cookies, extra_context=context) as result:
        check_result(result)


def test_bake_and_make_go_code_generator(cookies):
    context = {
        "project_category": "Code-Generator",
        "project_name": "generator 120",
    }
    with bake(cookies, extra_context=context) as result:
        check_result(result)
