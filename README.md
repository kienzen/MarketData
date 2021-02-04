# MarketData
Two Proof-of-Concepts for a File-based solution for sharing market data between DataMaster SimCorp Dimension installation and 50+ individual ones.
Both scripts prepare the source data in a common process and then launch separate threads for each installation.

First one was written in PowerShell, but a time of 3 minutes was found to be too slow.

The second was written in GO - the benchmark process showed average processing time of 2,8 seconds.

