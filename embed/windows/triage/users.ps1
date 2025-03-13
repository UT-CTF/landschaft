function Write-Divider {
    Write-Host
    Write-Host ("-" * 50)
    Write-Host
}

function Get-DomainInfo {
    $domain = (Get-WmiObject -Class Win32_ComputerSystem).Domain
    $domainJoined = $false
    if ($null -ne $domain -and $domain -ne "WORKGROUP") {
        $domainJoined = $true
    }

    if ($domainJoined) {
        Write-Host "This machine is joined to a domain."
        Write-Host "Domain: $domain"
        if ((Get-WmiObject -Class Win32_ComputerSystem).DomainRole -eq 5) {
            Write-Host "This machine is a domain controller."
            return $true
        }
    }
    else {
        Write-Host "This machine is not joined to a domain."
    }
    return $false
}

function Get-ADUserInfo {
    Write-Host "Users:"
    $users = Get-ADUser -Filter {Enabled -eq $true}
    Write-Host "Enabled AD Users ($($users.Length)):"
    foreach ($user in $users) {
        Write-Host "`t$($user.name)"
    }

    $users = Get-ADUser -Filter {Enabled -eq $false}
    Write-Host "Disabled AD Users ($($users.Length)):"
    foreach ($user in $users) {
        Write-Host "`t$($user.name)"
    }
}

function Get-UserInfo {
    $enabledUsers = @()
    $disabledUsers = @()

    $allUsers = Get-LocalUser
    foreach ($user in $allUsers) {
        if ($user.enabled) {
            $enabledUsers += $user
        }
        else {
            $disabledUsers += $user
        }
    }

    Write-Host "Users:"
    Write-Host "Enabled Local Users ($($enabledUsers.Length)):"
    foreach ($user in $enabledUsers) {
        Write-Host "`t$($user.name)"
    }

    Write-Host "Disabled Local Users: ($($disabledUsers.Length))"
    foreach ($user in $disabledUsers) {
        Write-Host "`t$($user.name)"
    }
}

function Get-ADGroupInfo {
    Write-Host "Groups:"

    $groupList = Get-ADGroup -Filter * -Properties *

    foreach ($group in $groupList) {
        
        $members = @()
        $group | Get-ADGroupMember | ? { $_.objectClass -eq "User" } | % {$members += $_.name}
        $memberof = $group | Select-Object -expandproperty memberof | Get-ADObject | % {$_.Name}
        
        if ($members.Length -gt 0) {
            Write-Host "$($group.name) ($($members.Length)): "
            $printlist = $members | Select-Object -First 10
            $printlist | % { Write-Host "`t$($_)" }
            if ($members.Length -gt 10) {
                Write-Host "`t... and $($members.Length - 10) more"
            }
            Write-Host "Member Of:"
            $memberof | % { Write-Host "`t$($_)" }
            Write-Host
        }
    }
}

function Get-GroupInfo {
    Write-Host "Groups:"

    $groupList = Get-LocalGroup

    foreach ($group in $groupList) {
        # Write-Host "Triaging Group: $($group.name)"
        $groupMembers = Get-LocalGroupMember -Group $group.name
        $members = @()
        $groupMembers | % {
            if ($_.objectClass -eq "User") {
                $members += $_.name
            }
        }
        if ($members.Length -gt 0) {
            Write-Host "$($group.name) ($($members.Length)): "
            foreach ($member in $members) {
                Write-Host "`t$member"
            }
        }
    }
}


$isdc = Get-DomainInfo

Write-Divider

if ($isdc) { Get-ADUserInfo }
else { Get-UserInfo }

Write-Divider

if($isdc) { Get-ADGroupInfo }
else { Get-GroupInfo }
