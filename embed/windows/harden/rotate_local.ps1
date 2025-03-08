param(
    [Parameter(Mandatory = $true)]
    [string]$Path,
    [int]$Length,
    [switch]$Apply
)

Write-Host "File path: $Path"
Write-Host "Length: $Length"