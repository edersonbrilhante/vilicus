FROM vilicus/postgres-presets:latest
ADD local-volumes/output/files-to-add.tar /
ADD local-volumes/output/files-to-remove.list /.files-to-remove.list
RUN if [ -s /.files-to-remove.list ]; then xargs -d '\n' -a /.files-to-remove.list rm && rm /.files-to-remove.list; fi
