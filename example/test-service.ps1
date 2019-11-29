param(
    [string]$Message
)

$i = 0
while($true) {
    Start-Sleep(1)
    Write-Host "Loop $i : $Message"

    if($i % 2 -eq 0) {
        #todo stderr
        Write-Host "Error $i : $Message"
    }

    $i++
}
