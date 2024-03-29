# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]


## [0.3.0] - 2021-09-30
### Added
 * **Default Branch Discovery** - if `--target` option isn't provided, glmr makes a
   Gitlab API call to figure out the name of the project's default branch.

## [0.2.1] - 2020-11-11
### Fixed
 * `glmr create` searches ancestor directories for a git project root (#4)

## [0.2.0] - 2020-11-03
### Added
 * Minor improvements
   - Improve README getting-started
   - Prompt wording/formatting
   - Example configuration file
   - `install` and `build-local` make targets for easy installation
 * `version` subcommand
 * `create` subcommand flags for squashing commits & deleing source branch on merge
### Changed
 * `create` uses available project MR templates over the user MR template.

## [0.1.0] - 2020-10-22
First real release
### Added
 * Helpers for user input and repository name processing
 * `create` command group for POSTing new MRs

