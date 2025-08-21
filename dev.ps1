# GrafiSearch Development Server Script
# Equivalent to 'make dev' for Windows

Write-Host "🚀 Starting GrafiSearch development servers..." -ForegroundColor Green

# Start all processes in parallel using PowerShell jobs
$env:APP_ENV = "dev"

# Job 1: Go server with hot reload
$job1 = Start-Job -ScriptBlock {
    Set-Location $using:PWD
    $env:APP_ENV = "dev"
    air
} -Name "GoServer"

# Job 2: Vite dev server
$job2 = Start-Job -ScriptBlock {
    Set-Location $using:PWD
    npm run dev
} -Name "ViteServer"

# Job 3: Template watcher
$job3 = Start-Job -ScriptBlock {
    Set-Location $using:PWD
    templ generate --watch
} -Name "TemplateWatcher"

Write-Host "✅ Started all development servers in background jobs" -ForegroundColor Green
Write-Host "📊 Job Status:" -ForegroundColor Yellow
Get-Job

Write-Host "`n🌐 Access your app at: http://localhost:8042" -ForegroundColor Cyan
Write-Host "🛠️  Vite dev server: http://localhost:5173" -ForegroundColor Cyan
Write-Host "`n⚡ Press Ctrl+C to stop all servers" -ForegroundColor Yellow

# Wait for user to press Ctrl+C and cleanup
try {
    # Monitor jobs and display output
    while ($true) {
        $runningJobs = Get-Job | Where-Object { $_.State -eq "Running" }
        if ($runningJobs.Count -eq 0) {
            Write-Host "❌ All jobs have stopped" -ForegroundColor Red
            break
        }
        Start-Sleep -Seconds 2
    }
}
catch {
    Write-Host "`n🛑 Stopping all development servers..." -ForegroundColor Red
}
finally {
    # Clean up jobs
    Get-Job | Stop-Job
    Get-Job | Remove-Job
    Write-Host "✅ All servers stopped. Bye! (≧◡≦)" -ForegroundColor Green
}