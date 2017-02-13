<!DOCTYPE html>
<html>
{{ template "components/header.html" }}
    <body>
{{ .Header }}
{{ .LayoutContent }}
{{ template "components/footer.html" }}
{{ .Scripts }}
    </body>
</html>
