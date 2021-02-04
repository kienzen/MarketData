function Get-PriceData {
param(
[Parameter(Mandatory=$true)]
[String]$RequestFileLocation,
[Parameter(Mandatory=$true)]
[System.Collections.ArrayList]$GruppenDataNoTimestamp,
[Parameter(Mandatory=$true)]
[System.Collections.ArrayList]$GruppenDataWithTimestamp
)

<#Import the request File for a Shard and format it accordingly. Create and add a custom property "ID---Currency" used for filtering#>
$shardlist = Import-Csv -Delimiter "	" -Path $RequestFileLocation -Encoding UTF7 -Header "ID","Währung","Kursquelle" 
$shardlist = $shardlist[1..($shardlist.count -1)] |Select-Object -Property @{Name="IDWährung";Expression={ $_.ID +"---"+ $_.Währung}}



<##############################Get the matching entries for a list WITH NO timestamps#################################>

<#A list for results#>
$shardListNoTimestamp = New-Object System.Collections.ArrayList

foreach ($gruppenmasterWP in $GruppenDataNoTimestamp){
    if($gruppenmasterWP.IDWährung -in $shardlist.IDWährung){
    
    $null = $shardListNoTimestamp.Add($gruppenmasterWp)

    }
}

<#Output the result list#>
$outputNoTimestamp = "C:\Users\J420429\Desktop\Marketdata\" + "Shard"+".txt"
$shardListNoTimestamp | Select-Object -Property "ID","WP-SOA","LegNr.","Kursdatum", "Kursart","Währung","Notierung","Kursquelle","Marktbereichs-Identifikation","Freikommentar"|ConvertTo-Csv -NoTypeInformation -Delimiter "	" | % {$_ -replace '"',''} |Out-File $outputNoTimestamp -Force

<##############################Get the matching entries for a list CONTAINING timestamps#################################>

<#A list for results#>
$shardListWithTimestamps = New-Object System.Collections.ArrayList

foreach ($gruppenmasterWP in $GruppenDataWithTimestamp){
    if($gruppenmasterWP.IDWährung -in $shardlist.IDWährung){
    
    $null = $shardListWithTimestamps.Add($gruppenmasterWp)

    }
}

<#Output the result list#>
$outputWithTimestamp = "C:\Users\J420429\Desktop\Marketdata\" + "Shard"+"Timestamp"+".txt"
$shardListWithTimestamps | Select-Object -Property "ID","WP-SOA","LegNr.","Kursdatum","Zeit", "Kursart","Kurs","Währung","Notierung","Kursquelle","Marktbereichs-Identifikation","Datenquelle"|ConvertTo-Csv -NoTypeInformation -Delimiter "	"  | % {$_ -replace '"',''}|Out-File $outputWithTimestamp -Force

}
