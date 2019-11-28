
$i = 0
while($true) {
    sleep(1)
    Write-Host "Loop $i"

    if($i % 2 -eq 0) {
        Write-Error "Error $i"
    }

    $i++
}
