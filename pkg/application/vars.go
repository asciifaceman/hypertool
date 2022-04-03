package application

import "regexp"

var openingTemplate string = `
<html>
<head/>
<body><center>
<h2>This page has been lovingly re-rendered for your convenience.
This tool is not associated with Patrick Leach,
so please don't bother him about it and respect his bandwidth.
</h2>
<br><br>
`
var closingTemplate string = `
</center></body>

</html>
`

var reFindProductId = regexp.MustCompile(`(\s[a-zA-Z]{1,5}\d{1,3}\s)`)
var reFindImage = regexp.MustCompile(`(http:\/\/www\.supertool\.com\/forsale\S*\.jpg)`)
var reFindCurrency = regexp.MustCompile(`(\$(?:(?:0|[1-9]\d*)(?:\.\d*)?|\.\d+))`)
