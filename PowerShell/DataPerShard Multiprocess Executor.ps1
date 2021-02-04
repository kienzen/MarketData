
$RequestFileLocation = "C:\Users\J420429\Desktop\Marketdata\INST01_MDA_20210114_1244.txt"
$GruppenDataNoTimestampLocation = "C:\Users\J420429\Desktop\Marketdata\MDL-WP_20210113.txt" 
$GruppenDataTimestampLocation = "C:\Users\J420429\Desktop\Marketdata\MDL-WPTS_20210113.txt"


<#Import the File with data WITHOUT timestamps from a Gruppenmasterand format it accordingly. Create and add a custom property "ID---Currency" used for filtering#>
$gruppenData = Import-Csv -path $GruppenDataNoTimestampLocation -Delimiter "	" -Encoding UTF7 -Header "ID","WP-SOA","LegNr.","Kursdatum", "Kursart","Kurs","Währung","Notierung","Kursquelle","Marktbereichs-Identifikation","Freikommentar"  
$gruppenData = $gruppenData[1..($gruppenData.Count -1)] | Select-Object -Property "ID","WP-SOA","LegNr.","Kursdatum", "Kursart","Währung","Notierung","Kursquelle","Marktbereichs-Identifikation","Freikommentar", @{Name="IDWährung";Expression={ $_.ID +"---"+ $_.Währung}} 

<#Import the File containing data WITH timestamps from a Gruppenmaster and format it accordingly. Create and add a custom property "ID---Currency" used for filtering#>
$gruppenDataTimestamp = Import-Csv -path  $GruppenDataTimestampLocation -Delimiter "	" -Encoding UTF7 -Header "ID","WP-SOA","LegNr.","Kursdatum","Zeit", "Kursart","Kurs","Währung","Notierung","Kursquelle","Marktbereichs-Identifikation","Datenquelle"  
$gruppenDataTimestamp = $gruppenDataTimestamp[1..($gruppenDataTimestamp.Count -1)] | Select-Object "ID","WP-SOA","LegNr.","Kursdatum","Zeit", "Kursart","Kurs","Währung","Notierung","Kursquelle","Marktbereichs-Identifikation","Datenquelle",@{Name="IDWährung";Expression={ $_.ID +"---"+ $_.Währung}}  


Start-Job -InitializationScript {import-module -name C:\Users\J420429\Desktop\Marketdata\DataPerShard.ps1} -ScriptBlock {Get-PriceData -RequestFileLocation C:\Users\J420429\Desktop\Marketdata\INST01_MDA_20210114_1244.txt -GruppenDataNoTimestamp $Using:gruppenData -GruppenDataWithTimestamp $Using:gruppenDataTimestamp }