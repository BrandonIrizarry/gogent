# Introduction

Gogent is a backend for implementing LLM agent frontend UIs, be they
CLI, TUI, or GUI based.

It is currently a work in progress.

# Installation

`go get github/BrandonIrizarry/gogent`

Note that specific versions are available, viz.,

`go get github/BrandonIrizarry/gogent@v0.1.0`


# Motivation

Gogent originally started out as a CLI REPL-based coding agent. I then
got the idea to upgrade the user interface to use the [BubbleTea](https://github.com/charmbracelet/bubbletea)
TUI framework. I first added the file picker [bubble](https://github.com/charmbracelet/bubbles), and had
plans to introduce configuration screens and such.

At some point, I started getting confused over where to insert the
BubbleTea stuff - the UI - on top of the LLM client logic; and so I
decided to decouple them.

Now, I can write *various* frontends for the LLM agent, which now
exports its own API for such frontends to consume. For example, I've
already written a proof-of-concept CLI client UI, which more or less
corresponds to the original Gogent project as it existed before the
decoupling into frontend and backend. This frontend is what I
currently use to test the backend, which is still in an early phase of
development.


