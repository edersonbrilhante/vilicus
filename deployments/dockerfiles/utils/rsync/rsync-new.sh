#!/usr/bin/env bash

# Based in https://github.com/neam/docker-diff-based-layers

echo "Cleanup output"
rm -f $OUT/*

echo "List of changes necessary to go from old image contents to new: "
rsync -a -x --human-readable --delete-after --itemize-changes --dry-run /postgres/ rsync://old:873/vilicus_conf/ | tee $OUT/changes.rsync.log

echo "Creating list of new files"
cat $OUT/changes.rsync.log | grep '^<f' | while read -a cols; do echo "${RESTRICT_DIFF_TO_PATH}/"${cols[@]:1}; done > $OUT/files-to-add.list

echo "Creating taz file with new files"
tar -cvf $OUT/files-to-add.tar -T $OUT/files-to-add.list 2>&1

\echo "Creating list of files to delete"
cat $OUT/changes.rsync.log | grep '^*deleting' | while read -a cols; do echo "${RESTRICT_DIFF_TO_PATH}/"${cols[@]:1}; done > $OUT/files-to-remove.list