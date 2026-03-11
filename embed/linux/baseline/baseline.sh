#!/bin/bash
# baseline.sh <dir>
# Creates a numbered snapshot in <dir>/baseline/N/ and diffs with the previous run.

set -euo pipefail

baseline_dir="${1:-.}/baseline"

create_snapshot_dir() {
    mkdir -p "$baseline_dir"
    chmod 700 "$baseline_dir"
    local latest
    latest=$(ls -1 "$baseline_dir" 2>/dev/null | grep -E '^[0-9]+$' | sort -n | tail -1)
    if [[ -z "$latest" ]]; then
        snap_num=1
    else
        snap_num=$((latest + 1))
    fi
    snap_dir="$baseline_dir/$snap_num"
    mkdir -p "$snap_dir"
    echo "Creating snapshot #$snap_num in $snap_dir"
}

gather_system_info() {
    {
        echo "=== System Information ==="
        uname -a
        echo
        echo "=== OS Release ==="
        cat /etc/os-release 2>/dev/null || true
        echo
        echo "=== Uptime ==="
        uptime
    } > "$snap_dir/system_info.txt"
}

gather_network_info() {
    {
        echo "=== Network Interfaces ==="
        ip a
        echo
        echo "=== Routes ==="
        ip route show all
        echo
        echo "=== Listening Ports ==="
        ss -nltup
        echo
        echo "=== Active Sessions ==="
        w
        echo
        echo "=== Firewall (iptables) ==="
        iptables-save 2>/dev/null || echo "iptables not available"
        echo
        ip6tables-save 2>/dev/null || true
        echo
        echo "=== DNS ==="
        cat /etc/resolv.conf 2>/dev/null || true
    } > "$snap_dir/network_info.txt"
}

gather_user_info() {
    {
        echo "=== /etc/passwd ==="
        cat /etc/passwd
        echo
        echo "=== /etc/group ==="
        cat /etc/group
        echo
        echo "=== Sudoers ==="
        cat /etc/sudoers 2>/dev/null || echo "Cannot read sudoers"
        echo
        echo "=== Authorized Keys ==="
        while IFS=: read -r user _ _ _ _ home _; do
            key_file="$home/.ssh/authorized_keys"
            if [[ -f "$key_file" ]]; then
                echo "--- $user ---"
                cat "$key_file"
                echo
            fi
        done < /etc/passwd
        echo
        echo "=== Last Logins ==="
        last -n 20 2>/dev/null || true
    } > "$snap_dir/user_info.txt"
}

gather_packages() {
    {
        echo "=== Installed Packages ==="
        if command -v dpkg-query &>/dev/null; then
            dpkg-query -W -f='${Package}\t${Version}\n'
        elif command -v rpm &>/dev/null; then
            rpm -qa --qf '%{NAME}\t%{VERSION}-%{RELEASE}\n'
        elif command -v apk &>/dev/null; then
            apk info -v
        else
            echo "No supported package manager found"
        fi
    } > "$snap_dir/packages.txt"
}

gather_services() {
    {
        echo "=== Services ==="
        if command -v systemctl &>/dev/null; then
            systemctl list-units --type=service --all --no-pager
            echo
            echo "=== Enabled Services ==="
            systemctl list-unit-files --type=service --state=enabled --no-pager
        else
            service --status-all 2>&1 || true
        fi
    } > "$snap_dir/services.txt"
}

gather_processes() {
    {
        echo "=== Running Processes ==="
        ps aux
        echo
        echo "=== Crontabs ==="
        crontab -l 2>/dev/null && echo "(root crontab above)" || echo "No root crontab"
        echo
        for dir in /var/spool/cron/crontabs /var/spool/cron; do
            if [[ -d "$dir" ]]; then
                for f in "$dir"/*; do
                    [[ -f "$f" ]] || continue
                    echo "--- $(basename "$f") ---"
                    cat "$f"
                    echo
                done
            fi
        done
        echo
        echo "=== System Cron ==="
        ls /etc/cron.d/ 2>/dev/null && cat /etc/cron.d/* 2>/dev/null || true
    } > "$snap_dir/processes.txt"
}

gather_filesystem() {
    {
        echo "=== Disk Usage ==="
        df -h
        echo
        echo "=== Mounts ==="
        cat /proc/mounts 2>/dev/null || mount
        echo
        echo "=== SUID/SGID Files ==="
        find / -xdev \( -perm -4000 -o -perm -2000 \) -type f 2>/dev/null | sort
    } > "$snap_dir/filesystem.txt"
}

diff_with_previous() {
    local prev_num=$((snap_num - 1))
    local prev_dir="$baseline_dir/$prev_num"
    if [[ ! -d "$prev_dir" ]]; then
        echo "No previous snapshot to compare with."
        return
    fi
    echo
    echo "=== Changes from snapshot #$prev_num to #$snap_num ==="
    for f in system_info network_info user_info packages services processes filesystem; do
        local cur="$snap_dir/${f}.txt"
        local prev="$prev_dir/${f}.txt"
        if [[ -f "$cur" && -f "$prev" ]]; then
            local d
            d=$(diff -u "$prev" "$cur" 2>/dev/null || true)
            if [[ -n "$d" ]]; then
                echo
                echo "--- $f ---"
                echo "$d"
            fi
        fi
    done
}

main() {
    create_snapshot_dir
    echo "Gathering system info..."
    gather_system_info
    echo "Gathering network info..."
    gather_network_info
    echo "Gathering user info..."
    gather_user_info
    echo "Gathering packages..."
    gather_packages
    echo "Gathering services..."
    gather_services
    echo "Gathering processes..."
    gather_processes
    echo "Gathering filesystem info..."
    gather_filesystem
    echo
    echo "Snapshot #$snap_num complete: $snap_dir"
    diff_with_previous
}

main
