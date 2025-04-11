#!/bin/bash
# Check for a given path
if [ -z "$1" ]; then
    mkdir -p "$1"
    BACKUP_FILE="$1/etc_backup_$(date +%Y%m%d_%H%M%S).tar.gz"
else
    BACKUP_DEST="/backup"
    # Create the backup destination folder if it doesn't exist
    mkdir -p "$BACKUP_DEST"
    BACKUP_FILE="$BACKUP_DEST/etc_backup_$(date +%Y%m%d_%H%M%S).tar.gz"
fi

echo "attempting to backup to $BACKUP_FILE"
# Backup the /etc/ folder
tar -czf "$BACKUP_FILE" /etc

# Check if the backup was successful
if [ $? -eq 0 ]; then
    echo "Backup successful! File saved to: $BACKUP_FILE"
else
    echo "Backup failed!"
    exit 1
fi
