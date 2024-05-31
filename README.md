# AOSS - IaC scanning GCB Plugin

## Description

scc-iac-scan-report-utils provide scripts for optional steps for AOSS's IAC Scanning GCB Plugin that is validator and SARIFConverter.

## SARIFConverter

sarif converter converters the IaC scanning result format to the more popular SARIF format
and uploads it to the Google Cloud Storage artifact registry bucket.

## Validator

Validator provides the status of the build based on the violation limitation parameter provided in the voilation expression as an input.
In case the voilation expression is not passed as an input by the user default parameter is taken

- Example expression 

``` 'Critical:10,High:10,Medium:10,Low:10,Operator:or' ```
