    #!/bin/bash
    # Check for a given path
    if [ -z "$1" ]; then
        echo "Please provide the path to the backup file."
        exit 1
    fi
    read -p "Are you sure you want to restore /etc from $1? (y/n) " response

    if [ "$response" != "y" ] && [ "$response" != "Y" ]; then
        echo "Restore cancelled. Maybe you're not sudo?"
        exit 0
    fi

    # Verify backup file exists
    if [ ! -f "$1" ]; then
        echo "Backup file not found!"
        exit 1
    fi

    # Check if file is a gzipped tar and contains /etc
    if ! tar -tzf "$1" | grep -q "^etc/"; then
        echo "Invalid backup file format!"
        exit 1
    fi
    # Restore the /etc/ folder from the backup
    tar -xzf "$1" -C /

    # status code
    if [ $? -eq 0 ]; then
        echo "Restore successful!"
    else
        echo "Restore failed!"
        exit 1
    fi