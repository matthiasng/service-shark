
$errorActionPreference = "Stop"

$currentPwshExe = (Get-Process -id $pid | Get-Item).FullName
$psTestScript = "$PSScriptRoot/test-service.ps1"

$serviceName = "ExampleService-" + (New-Guid)

$serviceLogDirecotry = "$PSScriptRoot/log"
$serviceCommand = $currentPwshExe

$serviceSharkPath = "$PSScriptRoot/../service-shark.exe"

if(!(Test-Path $serviceSharkPath)) {
    Write-Host "Run: go build -o service-shark.exe .\main.go"
    exit
}

$serviceBinPath = $serviceSharkPath
$serviceBinPath += ' -n "' + $serviceName + '"'
$serviceBinPath += ' -l "' + $serviceLogDirecotry + '"'
$serviceBinPath += ' -c "' + $serviceCommand + '"'
$serviceBinPath += ' -a "' + $psTestScript + '"'
$serviceBinPath += ' -a "-Message" -a "TestMessage with space!"'
$serviceBinPath += ' -a "-PathVar" -a "bind:PATH"'


Write-Host "New-Service:"
Write-Host $serviceBinPath

try {
    New-Service -Name $serviceName -BinaryPathName $serviceBinPath -StartupType "Manual"
    Start-Service -Name $serviceName

    Start-Sleep(5)
}
finally {
    Stop-Service -Name $serviceName
    Remove-Service -Name $serviceName
}
