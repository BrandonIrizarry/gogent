#!/usr/bin/env bash

# Print the list of all untracked files.
#
# There are other ways to achieve this, but this appears to be the
# cleanest way, since it both lists actual directories, and doesn't
# insist on listing any files under ignored directories.
#
# See
#
# https://stackoverflow.com/a/2196755/4570292
#
git clean -ndX | awk '{print $3}'
