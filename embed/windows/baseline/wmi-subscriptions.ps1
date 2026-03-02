param(
    [Parameter(Mandatory = $true)]
    [string]$BaselinePath
)

if (-not (Test-Path $BaselinePath)) {
    New-Item -ItemType Directory -Path $BaselinePath | Out-Null
}

$Namespace = "root\subscription"

# Write-Host "[*] Enumerating WMI Event Filters..."
$filters = Get-CimInstance -Namespace $Namespace -ClassName __EventFilter -ErrorAction SilentlyContinue |
    ForEach-Object {
        [PSCustomObject]@{
            Name           = $_.Name
            Query          = $_.Query
            QueryLanguage  = $_.QueryLanguage
            EventNamespace = $_.EventNamespace
            CreatorSID     = ($_.CreatorSID -join ",")
        }
    } | Sort-Object Name

$filters | Export-Csv "$BaselinePath\wmi-eventfilters.csv" -NoTypeInformation


# Write-Host "[*] Enumerating CommandLineEventConsumers..."
$cmdConsumers = Get-CimInstance -Namespace $Namespace -ClassName CommandLineEventConsumer -ErrorAction SilentlyContinue |
    ForEach-Object {
        [PSCustomObject]@{
            Name                = $_.Name
            CommandLineTemplate = $_.CommandLineTemplate
            RunInteractively    = $_.RunInteractively
            WorkingDirectory    = $_.WorkingDirectory
            CreatorSID          = ($_.CreatorSID -join ",")
        }
    } | Sort-Object Name

$cmdConsumers | Export-Csv "$BaselinePath\wmi-commandlineconsumers.csv" -NoTypeInformation


# Write-Host "[*] Enumerating ActiveScriptEventConsumers..."
$scriptConsumers = Get-CimInstance -Namespace $Namespace -ClassName ActiveScriptEventConsumer -ErrorAction SilentlyContinue |
    ForEach-Object {
        [PSCustomObject]@{
            Name         = $_.Name
            ScriptingEngine = $_.ScriptingEngine
            ScriptText   = $_.ScriptText
            CreatorSID   = ($_.CreatorSID -join ",")
        }
    } | Sort-Object Name

$scriptConsumers | Export-Csv "$BaselinePath\wmi-activescriptconsumers.csv" -NoTypeInformation


# Write-Host "[*] Enumerating Other EventConsumer types..."
$allConsumers = Get-CimInstance -Namespace $Namespace -ClassName __EventConsumer -ErrorAction SilentlyContinue |
    Where-Object {
        $_.CimClass.CimClassName -notin @("CommandLineEventConsumer","ActiveScriptEventConsumer")
    } |
    ForEach-Object {
        [PSCustomObject]@{
            Name       = $_.Name
            ClassName  = $_.CimClass.CimClassName
            CreatorSID = ($_.CreatorSID -join ",")
        }
    } | Sort-Object Name

$allConsumers | Export-Csv "$BaselinePath\wmi-allconsumers.csv" -NoTypeInformation


# Write-Host "[*] Enumerating Filter-To-Consumer Bindings..."
$bindings = Get-CimInstance -Namespace $Namespace -ClassName __FilterToConsumerBinding -ErrorAction SilentlyContinue |
    ForEach-Object {
        [PSCustomObject]@{
            Filter   = ($_.Filter).Name
            Consumer = ($_.Consumer).Name
        }
    } | Sort-Object Filter, Consumer

$bindings | Export-Csv "$BaselinePath\wmi-bindings.csv" -NoTypeInformation


# Write-Host "[+] WMI subscription baseline completed."
