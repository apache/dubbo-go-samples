# Static router configuration

This directory contains static router samples configured directly in code.
Here, `static_config` means the router rules are injected statically on the consumer side, instead of being delivered dynamically from a config center.

English | [中文](README_CN.md)

## Prerequisites

- Go 1.25+.

## Sub-samples

- `condition`: service-scope static condition router sample with direct URLs
- `tag`: application-scope static tag router sample with direct URLs

## How to use

Choose one sub-sample and follow its own README:

- `condition/README.md`
- `tag/README.md`

The `condition` sample uses direct URLs only.
The `tag` sample uses direct URLs only.
Neither sample requires a config center.
Each sub-sample can be run locally after following its own prerequisites and setup steps.
