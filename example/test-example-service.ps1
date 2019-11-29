
$currentPwshExe = (Get-Process -id $pid | Get-Item).FullName

$serviceName = "ExampleService-" + (New-Guid)

$serviceLogDirecotry = "$PSScriptRoot/log"
$serviceCommand = $currentPwshExe
$serviceArguments = "$PSScriptRoot/test-service.ps1 -Message 'TestMessage with space!' -PathVar bind:PATH"

$serviceBinPath = 'P:\projects\service-wrapper\main.exe'
$serviceBinPath += ' -logdir "' + $serviceLogDirecotry + '"'
$serviceBinPath += ' -command "' + $serviceCommand + '"'
$serviceBinPath += ' -arguments "' + $serviceArguments + '"'

Write-Host "New-Service:"
Write-Host $serviceBinPath

New-Service -Name $serviceName -BinaryPathName $serviceBinPath -StartupType "Manual"
Start-Service -Name $serviceName

Start-Sleep(5)

Stop-Service -Name $serviceName
Remove-Service -Name $serviceName
