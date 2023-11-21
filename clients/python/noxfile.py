"""Nox sessions."""
import os
import shutil
import sys
from pathlib import Path
from textwrap import dedent

import nox

try:
    from nox_poetry import Session, session
except ImportError:
    message = f"""\
    Nox failed to import the 'nox-poetry' package.

    Please install it using the following command:

    {sys.executable} -m pip install nox-poetry"""
    raise SystemExit(dedent(message)) from None


package = "model_registry"
python_versions = ["3.11", "3.10", "3.9"]
nox.needs_version = ">= 2021.6.6"
nox.options.sessions = (
    "lint",
    "mypy",
    "tests",
    "docs-build",
)


TFX_NIGHTLY_URL = "https://pypi-nightly.tensorflow.org/simple"


@session(python=python_versions[0])
def lint(session: Session) -> None:
    """Lint using ruff."""
    session.install("ruff")

    session.run("ruff", "check", "src")


@session(python=python_versions[0])
def mypy(session: Session) -> None:
    """Type check with mypy."""
    session.install("mypy")

    session.run("mypy", "src")


@session(python=python_versions)
def tests(session: Session) -> None:
    """Run the test suite."""
    session.install(
        "--extra-index-url",
        TFX_NIGHTLY_URL,
        ".",
        "coverage[toml]",
        "pytest",
        "pytest-cov",
    )

    try:
        session.run(
            "pytest",
            "--cov",
            "--cov-config=pyproject.toml",
            *session.posargs,
            env={"COVERAGE_FILE": f".coverage.{session.python}"},
        )
    finally:
        if session.interactive:
            session.notify("coverage", posargs=[])


@session(python=python_versions[0])
def coverage(session: Session) -> None:
    """Produce the coverage report."""
    args = session.posargs or ["report"]

    session.install("coverage[toml]")

    if not session.posargs and any(Path().glob(".coverage.*")):
        session.run("coverage", "combine")

    session.run("coverage", *args)


@session(name="docs-build", python=python_versions[0])
def docs_build(session: Session) -> None:
    """Build the documentation."""
    args = session.posargs or ["docs", "docs/_build"]
    if not session.posargs and "FORCE_COLOR" in os.environ:
        args.insert(0, "--color")

    session.install(
        "--extra-index-url",
        TFX_NIGHTLY_URL,
        ".",
        "sphinx",
        "furo",
        "myst-parser[linkify]",
    )

    build_dir = Path("docs", "_build")
    if build_dir.exists():
        shutil.rmtree(build_dir)

    session.run("sphinx-build", *args)


@session(python=python_versions[0])
def docs(session: Session) -> None:
    """Build and serve the documentation with live reloading on file changes."""
    args = session.posargs or ["--open-browser", "docs", "docs/_build"]
    session.install(
        "--extra-index-url",
        TFX_NIGHTLY_URL,
        ".",
        "sphinx",
        "furo",
        "myst-parser[linkify]",
        "sphinx-autobuild",
    )

    build_dir = Path("docs", "_build")
    if build_dir.exists():
        shutil.rmtree(build_dir)

    session.run("sphinx-autobuild", *args)
