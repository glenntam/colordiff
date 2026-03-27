# colordiff

colordiff is an updated fork of [artyom](https://github.com/artyom/colordiff)'s 2022 colordiff Go port.

## updates

- Introduces stricter linter checks.
- Passes context to the subprocess.
- Reverts coloring to match the original kimmel 1991 [colordiff](https://github.com/kimmel/colordiff)
- Adds Makefile for convenience

## installation

> go install github.com/glenntam/colordiff@latest

## usage

Exactly the same as the original:

> colordiff file1 file2

or

> diff -u file1 file2 | colordiff

## motivation

This Go version is distro/arch agnostic and can be installed without root. I needed to create this the moment I was on a remote machine without root access. My eyes would have gone blind without it.
