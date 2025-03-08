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
    }
    else {
        Write-Host "This machine is not joined to a domain."
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

function Get-GroupInfo {
    Write-Host "Groups:"

    $groupList = Get-LocalGroup

    foreach ($group in $groupList) {
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


Get-DomainInfo

Write-Divider

Get-UserInfo

Write-Divider

Get-GroupInfo