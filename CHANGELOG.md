# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Changed
- Improved overall speed by ~40% by removing fmt package.
- Fixed package comments

## [0.0.1] - 2020-09-21
### Added
- Functions to wrap [github.com/shopspring/decimal v1.2.0](https://github.com/shopspring/decimal/releases/tag/v1.2.0) Decimal and NullDecimal types.
- Helper functions GetName(), SetName(), Resolve(), Decimal(), Math()
- Zero var with name zero representing zero
- New methods duplicated WithName for easier drop in nature.
- NewFromDecimal() and NewFromDecimalWithName() for easier integration.
- ResolveTo() method wraps SetName() and Resolve() for easier resolution.
