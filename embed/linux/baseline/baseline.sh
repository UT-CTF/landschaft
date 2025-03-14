#!/bin/bash

# Directory for baselining information
baseline_dir="./baseline"

# Function to create numbered directories
create_numbered_directory() {
    mkdir -p "$baseline_dir"
    chmod 600 "$baseline_dir"
    latest_dir=$(ls -1 "$baseline_dir" | sort -n | tail -n 1)
    if [[ -z $latest_dir ]]; then
        next_dir="1"
    else
        next_dir=$((latest_dir + 1))
    fi
    mkdir -p "$baseline_dir/$next_dir"
}

# Function to gather system information
gather_system_info() {
    echo "=== System Information ===" > "$baseline_dir/$next_dir/system_info.txt"
    uname -a >> "$baseline_dir/$next_dir/system_info.txt"
}

# Function to gather network information
gather_network_info() {
    echo "=== Network Configuration ===" > "$baseline_dir/$next_dir/network_info.txt"
    ip a >> "$baseline_dir/$next_dir/network_info.txt"
    echo >> "$baseline_dir/$next_dir/network_info.txt"
    echo "=== Routes ===" >> "$baseline_dir/$next_dir/network_info.txt"
    ip route show all >> "$baseline_dir/$next_dir/network_info.txt"
    echo >> "$baseline_dir/$next_dir/network_info.txt"
    echo "=== Netstat ===" >> "$baseline_dir/$next_dir/network_info.txt"
    netstat -tuln >> "$baseline_dir/$next_dir/network_info.txt"
    echo >> "$baseline_dir/$next_dir/network_info.txt"
    echo "=== IpTables ===" >> "$baseline_dir/$next_dir/network_info.txt"
    iptables-save >> "$baseline_dir/$next_dir/network_info.txt"
    echo >> "$baseline_dir/$next_dir/network_info.txt"
    ip6tables-save >> "$baseline_dir/$next_dir/network_info.txt"
    echo >> "$baseline_dir/$next_dir/network_info.txt"
    echo "=== Sessions ===" >> "$baseline_dir/$next_dir/network_info.txt"
    w >> "$baseline_dir/$next_dir/network_info.txt"
    echo >> "$baseline_dir/$next_dir/network_info.txt"
    echo "=== Ports ===" >> "$baseline_dir/$next_dir/network_info.txt"
    ss -nltup >> "$baseline_dir/$next_dir/network_info.txt"
    echo "=== resolv.conf ===" >> "$baseline_dir/$next_dir/network_info.txt"
    cat "/etc/resolv.conf" >> "$baseline_dir/$next_dir/network_info.txt"
}

# Function to gather user information
gather_user_info() {
    echo "=== User Information ===" > "$baseline_dir/$next_dir/user_info.txt"
    cat /etc/passwd >> "$baseline_dir/$next_dir/user_info.txt"
    echo >> "$baseline_dir/$next_dir/user_info.txt"
    echo "=== Sudoers ===" >> "$baseline_dir/$next_dir/user_info.txt"
    cat /etc/sudoers >> "$baseline_dir/$next_dir/user_info.txt"
    echo >> "$baseline_dir/$next_dir/user_info.txt"
    echo "=== Authorized Keys ===" >> "$baseline_dir/$next_dir/user_info.txt"
    all_users=$(getent passwd | cut -d: -f1)
    for user in $all_users; do
        local authorized_keys_file="~$user/.ssh/authorized_keys"
        if [ -f "$authorized_keys_file" ]; then
            echo "=== Authorized keys for user: $user ===" >> "$baseline_dir/$next_dir/user_info.txt"
            cat "$authorized_keys_file" >> "$baseline_dir/$next_dir/user_info.txt"
            echo >> "$baseline_dir/$next_dir/user_info.txt"
        fi
    done
}

# Function to gather installed packages information
gather_installed_packages() {
    echo "=== Installed Packages ===" > "$baseline_dir/$next_dir/installed_packages.txt"
    dpkg -l >> "$baseline_dir/$next_dir/installed_packages.txt"  # For Debian-based systems
    # rpm -qa >> "$baseline_dir/$next_dir/installed_packages.txt"  # For RPM-based systems like Red Hat, CentOS
}

# Function to gather filesystem information
gather_filesystem_info() {
    echo "=== Filesystem Information ===" > "$baseline_dir/$next_dir/filesystem_info.txt"
    df -h >> "$baseline_dir/$next_dir/filesystem_info.txt"
}

gather_process_info() {
    echo "=== Process Information ===" > "$baseline_dir/$next_dir/process_info.txt"
    ps aux >> "$baseline_dir/$next_dir/process_info.txt"
    echo >> "$baseline_dir/$next_dir/process_info.txt"
    all_users=$(getent passwd | cut -d: -f1)
    for user in $all_users; do
        local user_crontab=$(crontab -u "$user" -l 2>/dev/null)
        if [ -n "$user_crontab" ]; then
            echo "=== Crontab for user: $user ===" >> "$baseline_dir/$next_dir/process_info.txt"
            echo "$user_crontab" >> "$baseline_dir/$next_dir/process_info.txt"
        fi
    done
}

# Function to diff files from the current run with the previous run
diff_with_previous_run() {
    prev_dir=$(($next_dir - 1))
    if [ -d "$baseline_dir/$prev_dir" ]; then
        echo "=== Differences from previous run ==="
        echo "=== System ==="
        diff -u "$baseline_dir/$prev_dir/system_info.txt" "$baseline_dir/$next_dir/system_info.txt"
        echo "=== Network ==="
        diff -u "$baseline_dir/$prev_dir/network_info.txt" "$baseline_dir/$next_dir/network_info.txt"
        echo "=== User ==="
        diff -u "$baseline_dir/$prev_dir/user_info.txt" "$baseline_dir/$next_dir/user_info.txt"
        echo "=== Packages ==="
        diff -u "$baseline_dir/$prev_dir/installed_packages.txt" "$baseline_dir/$next_dir/installed_packages.txt"
        echo "=== Filesystem ==="
        diff -u "$baseline_dir/$prev_dir/filesystem_info.txt" "$baseline_dir/$next_dir/filesystem_info.txt"
        echo "=== Process ==="
        diff -u "$baseline_dir/$prev_dir/process_info.txt" "$baseline_dir/$next_dir/process_info.txt"
    else
        echo "No previous run to compare with."
    fi
}

# Main function to execute all the gathering functions and diff
main() {
    echo "Creating system baseline..."
    create_numbered_directory
    gather_system_info
    gather_network_info
    gather_user_info
    gather_installed_packages
    gather_filesystem_info
    gather_process_info
    echo "System baseline created in directory: $baseline_dir/$next_dir"
    diff_with_previous_run
}

# Call the main function
main
