---
title: Agent skill
weight: 9
---

`skills/inkpress/` in the repository is an
[Agent Skill](https://docs.claude.com/en/docs/agents-and-tools/agent-skills).
It teaches Claude to plan a deck, write the slides with the right layouts and
notes, build the PDF and check the pages, and change templates.

## Install

To use it in Claude Code, copy or symlink it into `~/.claude/skills/` for all
projects, or into a project's `.claude/skills/`:

```sh
git clone https://github.com/inkcheck/inkpress
ln -s "$PWD/inkpress/skills/inkpress" ~/.claude/skills/inkpress
```

Then ask for a deck:

```text
Make a 10-minute talk from notes.md, with speaker notes.
```

## Diagrams

The skill works with the
[inkline skill](https://inkcheck.github.io/inkline/) for diagrams. Install
both to get decks with diagrams in them.
