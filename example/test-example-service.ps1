
$errorActionPreference = "Stop"

$currentPwshExe = (Get-Process -id $pid | Get-Item).FullName
$psTestScript = "$PSScriptRoot/test-service.ps1"

$serviceName = "ExampleService-" + (New-Guid)

$serviceWorkingDir = $PSScriptRoot
$serviceCommand = $currentPwshExe

$serviceSharkPath = "$PSScriptRoot/../service-shark.exe"

if(!(Test-Path $serviceSharkPath)) {
    Write-Host "service-shark.exe not found. Run 'go build -o service-shark.exe .\main.go'"
    exit
}

$serviceBinPath = $serviceSharkPath
$serviceBinPath += ' -name "' + $serviceName + '"'
$serviceBinPath += ' -workdir "' + $serviceWorkingDir + '"'
$serviceBinPath += ' -cmd "' + $serviceCommand + '"'
$serviceBinPath += ' --'
$serviceBinPath += ' "' + $psTestScript + '"'
$serviceBinPath += ' -Message "TestMessage with space!"'
$serviceBinPath += ' -PathVar "env:PATH"'


Write-Host "New-Service:"
Write-Host $serviceBinPath

try {
    New-Service -Name $serviceName -BinaryPathName $serviceBinPath -StartupType "Manual"
    Start-Service -Name $serviceName

    Start-Sleep(5)
} finally {
    Stop-Service -Name $serviceName
    Remove-Service -Name $serviceName
}
