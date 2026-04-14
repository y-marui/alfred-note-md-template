.PHONY: help install hooks lint format typecheck test test-cov vendor build deploy release run clean update-charter

# Default target
help:
	@echo "Alfred Workflow Template - Development Commands"
	@echo ""
	@echo "  make install     Install dev dependencies (uv sync + pre-commit)"
	@echo "  make hooks       Install pre-commit hooks"
	@echo "  make lint        Run ruff linter"
	@echo "  make format      Auto-format with ruff format + ruff --fix"
	@echo "  make typecheck   Run mypy type checker"
	@echo "  make test        Run tests"
	@echo "  make test-cov    Run tests with coverage report"
	@echo "  make vendor      Install runtime deps into workflow/vendor/"
	@echo "  make build       Build .alfredworkflow package"
	@echo "  make deploy      Build and install into Alfred"
	@echo "  make release     Create GitHub Release (requires git tag)"
	@echo "  make run Q=''    Simulate Alfred with query Q"
	@echo "  make clean       Remove build artifacts"
	@echo "  make update-charter  Pull latest dev-charter via git subtree"
	@echo ""

# ---------------------------------------------------------------------------
# Setup
# ---------------------------------------------------------------------------
install:
	uv sync
	pre-commit install

hooks:
	pre-commit install

# ---------------------------------------------------------------------------
# Code quality
# ---------------------------------------------------------------------------
lint:
	ruff check src/ tests/

format:
	ruff format src/ tests/
	ruff check --fix src/ tests/

typecheck:
	mypy src/

# ---------------------------------------------------------------------------
# Tests
# ---------------------------------------------------------------------------
test:
	pytest

test-cov:
	pytest --cov=src --cov-report=term-missing --cov-report=html

# ---------------------------------------------------------------------------
# Build
# ---------------------------------------------------------------------------
vendor:
	./scripts/vendor.sh

build: vendor
	./scripts/build.sh

deploy: build
	@open dist/*.alfredworkflow

release:
	./scripts/release.sh

# ---------------------------------------------------------------------------
# Local dev
# ---------------------------------------------------------------------------
run:
	@./scripts/dev.sh "$(Q)"

# ---------------------------------------------------------------------------
# Maintenance
# ---------------------------------------------------------------------------
clean:
	rm -rf .build dist/ .coverage htmlcov/ .mypy_cache/ .pytest_cache/
	find . -name "__pycache__" -type d -exec rm -rf {} + 2>/dev/null || true
	find . -name "*.pyc" -delete

# ---------------------------------------------------------------------------
# Dev Charter
# ---------------------------------------------------------------------------
update-charter:
	git subtree pull --prefix=docs/dev-charter dev-charter main --squash
