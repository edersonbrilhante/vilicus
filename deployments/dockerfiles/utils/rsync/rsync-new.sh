#!/usr/bin/env bash

# Based in https://github.com/neam/docker-diff-based-layers

# Cleanup output
rm $OUT/*

# Use rsync to figure out what changes are necessary to go from old image contents to new
echo "List of changes necessary to go from old image contents to new: "
rsync -a -x --human-readable --delete-after --checksum --dry-run --itemize-changes $RESTRICT_DIFF_TO_PATH rsync://old:873/vilicus_conf/ | tee $OUT/changes.rsync.log

# Add files to add to a tar archive
cat $OUT/changes.rsync.log | grep '^<f' | while read -a cols; do echo "$RESTRICT_DIFF_TO_PATH/"${cols[@]:1}; done > $OUT/files-to-add.list
tar -cf $OUT/files-to-add.tar -T $OUT/files-to-add.list 2>&1 | grep -v  "Removing leading"

# Add files to remove to a list
cat $OUT/changes.rsync.log | grep '^*deleting' | while read -a cols; do echo "$RESTRICT_DIFF_TO_PATH/"${cols[@]:1}; done > $OUT/files-to-remove.list

echo "done" > $OUT/job.rsync.log