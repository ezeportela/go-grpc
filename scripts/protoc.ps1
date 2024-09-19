$baseDir = Get-Location
$protoFiles = Get-ChildItem -Recurse -Path $baseDir -Filter *.proto

foreach ($proto in $protoFiles) {
    $relativePath = $proto.FullName.Substring($baseDir.Path.Length + 1)
    $relativePath = "./$relativePath"
    
    Write-Host "Processing $relativePath"
    & protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative $relativePath
}