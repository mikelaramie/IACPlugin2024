# AOSS - IaC scanning GCB Plugin

## Description

scc-iac-scan-report-utils provide script for handling the result from gcloud scc iac-validation-reports create.
Currently it offers the following utility files

## SARIFConverter

SARIFConverter converters the report generated by "gcloud scc iac-validation-reports create" command to a more
popular SARIF format.

## Validator

It checks the scc iac-validation-report against limits set by failure criteria and returns the validation outcome.

- These validation failure criteria could be passed as an expression in a form of input to the script.

    ``` 'Critical:2,Low:5,Operator:or' ```

- If no expression is passed to the scipt. default criteria is used to perform these validation

``` 'Critical:1,High:1,Medium:1,Low:1,Operator:or' ```

> [!NOTE]
> - For Operator only AND and OR operators are supported.
> - Each expression should have an operator only once.
> - All Severity: Critical, High, Medium, Low can be present in the expression at most once.




